package middleware

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

// https://github.com/yurishkuro/opentracing-tutorial

func jaeger() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tracer, closer := tracing.Init("hello-world")
			defer closer.Close()
			opentracing.SetGlobalTracer(tracer)

			span := tracer.StartSpan("say-hello")
			//			span.SetTag("hello-to", helloTo)
			defer span.Finish()

			ctx := r.Context()
			ctx = opentracing.ContextWithSpan(ctx, span)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
