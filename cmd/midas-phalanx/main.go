// https://github.com/enricofoltran/simple-go-server
// https://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/utopia-planitia/exocomp/middleware"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	GitTag    string = ""
	GitCommit string = ""
	listen    string
)

func main() {
	flag.StringVar(&listen, "listen", ":8080", "address to listen on")
	flag.Parse()

	logger := logrus.New()

	if os.Getenv("LOG_TEXT") == "true" || terminal.IsTerminal(int(os.Stdout.Fd())) {
		logger.Formatter = &logrus.TextFormatter{ForceColors: os.Getenv("LOG_COLORS") != "false"}
	} else {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	if os.Getenv("LOG_LEVEL") == "trace" {
		logger.Level = logrus.TraceLevel
	}
	if os.Getenv("LOG_LEVEL") == "info" {
		logger.Level = logrus.InfoLevel
	}
	if os.Getenv("LOG_LEVEL") == "warn" {
		logger.Level = logrus.WarnLevel
	}

	if GitTag != "" {
		logger.Println("GitTag:", GitTag)
	}
	if GitCommit != "" {
		logger.Println("GitCommit:", GitCommit)
	}

	logger.Println("server is starting...")

	router := http.NewServeMux()
	router.Handle("/", index())

	middleware.Run(listen, logger, router)
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		span, ctx := opentracing.StartSpanFromContext(ctx, "formatString")
		defer span.Finish()

		helloTo := "david"

		span.SetTag("hello-to", helloTo)

		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		span.LogFields(
			log.String("event", "string-format"),
			log.String("value", helloStr),
		)

		println(helloStr)
		span.LogKV("event", "println")

		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fakeWorkSpan, ctx := opentracing.StartSpanFromContext(ctx, "some work")
		time.Sleep(0 * time.Second)
		fakeWorkSpan.Finish()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, world!")
	})
}
