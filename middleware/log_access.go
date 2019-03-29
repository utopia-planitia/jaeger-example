package middleware

import (
	"log"
	"net/http"
)

func logAccess(l *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}
			l.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

			h.ServeHTTP(w, r)
		})
	}
}
