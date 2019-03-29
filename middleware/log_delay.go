package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func logDelay(l *log.Logger, err time.Duration, warn time.Duration, info time.Duration) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			start := time.Now()
			h.ServeHTTP(w, r)
			elapsed := time.Since(start)

			m := l.WithFields(log.Fields{
				"requestID":    requestID,
				"duration":     elapsed,
				"microseconds": int64(elapsed),
			})
			if elapsed > err {
				m.Error("request duration")
				return
			}
			if elapsed > warn {
				m.Warn("request duration")
				return
			}
			if elapsed > info {
				m.Info("request duration")
				return
			}
		})
	}
}
