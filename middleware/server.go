package middleware

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func Run(listen string, logger *log.Logger, routes *http.ServeMux) {

	running := int32(0)

	routes.Handle("/healthz", healthz(&running))

	access := logAccess(logger)(routes)
	delay := logDelay(logger)(access)
	requestID := requestID(rand.Int63)(delay)

	handler := requestID

	server := &http.Server{
		Addr:         listen,
		Handler:      handler,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Printf("server is shutting down\n")
		atomic.StoreInt32(&running, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Printf("server is ready to handle requests at %s\n", listen)
	atomic.StoreInt32(&running, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("could not listen on %s: %v\n", listen, err)
	}
	<-done
	logger.Printf("server stopped\n")
}
