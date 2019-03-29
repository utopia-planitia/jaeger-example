package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// log2LogrusWriter exploits the documented fact that the standard
// log pkg sends each log entry as a single io.Writer.Write call:
// https://golang.org/pkg/log/#Logger
// https://github.com/sirupsen/logrus/issues/436
type log2LogrusWriter struct {
	f func(args ...interface{})
}

func (w *log2LogrusWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	w.f(string(b))
	return n, nil
}

func Run(listen string, log2 *logrus.Logger, routes *http.ServeMux) {

	running := int32(0)

	logger := log.New(&log2LogrusWriter{f: log2.Warn}, "", 0)

	routes.Handle("/healthz", healthz(&running))

	handler := http.Handler(routes)

	handler = logAccess(log2)(handler)
	handler = logDelay(log2, 200*time.Microsecond, 150*time.Microsecond, 100*time.Microsecond)(handler)
	handler = requestID(rand.Int63)(handler)

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
