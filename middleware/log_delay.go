package middleware

import (
	"log"
	"net/http"
	"time"
)

func logDelay(l *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			start := time.Now()
			h.ServeHTTP(w, r)
			elapsed := time.Since(start)

			l.Println(requestID, "took", elapsed)

		})
	}
}
