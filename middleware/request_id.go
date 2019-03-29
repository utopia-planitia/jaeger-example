package middleware

import (
	"context"
	"net/http"
	"unsafe"
)

type contextKey string

const requestIDKey = contextKey("requestId")

func requestID(randInt63 func() int64) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ID := r.Header.Get("X-Request-Id")
			if ID == "" {
				ID = randString(randInt63, 32)
			}

			r.Header.Set("X-Request-Id", ID)
			w.Header().Set("X-Request-Id", ID)

			ctx := context.WithValue(r.Context(), requestIDKey, ID)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

// from https://stackoverflow.com/a/31832326
const (
	letterBytes   = "0123456789abcde"
	letterIdxBits = 4                    // 4 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(r func() int64, n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, r(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
