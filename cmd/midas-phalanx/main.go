// https://github.com/enricofoltran/simple-go-server
// https://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/sirupsen/logrus"
	"github.com/utopia-planitia/exocomp/middleware"
)

var (
	Version      string = ""
	GitTag       string = ""
	GitCommit    string = ""
	GitTreeState string = ""
	listen       string
	healthy      int32
)

func main() {
	flag.StringVar(&listen, "listen", ":8080", "address to listen on")
	flag.Parse()

	// https://github.com/sirupsen/logrus/issues/436
	//	loggerOrg := log.New(os.Stdout, "http: ", log.LstdFlags)
	loggerOrg := log.New(os.Stderr, "", log.LstdFlags)
	logger := logrus.New()
	loggerOrg.SetOutput(logger.Writer())

	//	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter = new(logrus.TextFormatter)                     //default
	logger.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	logger.Println("Simple go server")
	logger.Println("Version:", Version)
	logger.Println("GitTag:", GitTag)
	logger.Println("GitCommit:", GitCommit)
	logger.Println("GitTreeState:", GitTreeState)

	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", index())

	middleware.Run(listen, loggerOrg, router)

	os.Stderr.Sync()
	os.Stdout.Sync()
	//	log.Printf("main exit\n")
	//	time.Sleep(2 * time.Second)
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		time.Sleep(2 * time.Second)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}
