package middleware

import (
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func logAccess(l *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span, _ := opentracing.StartSpanFromContext(r.Context(), "logAccess")

			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}
			l.WithFields(log.Fields{
				"requestID": requestID,
				"method":    r.Method,
				"path":      r.URL.Path,
				"remote":    r.RemoteAddr,
				"userAgent": r.UserAgent(),
			}).Info("accessed")

			span.Finish()

			h.ServeHTTP(w, r)
		})
	}
}
