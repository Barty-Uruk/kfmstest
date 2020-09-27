package middlewares

import (
	"context"
	"fmt"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/Barty-Uruk/kfmstest/configs"
	logger "github.com/Barty-Uruk/kfmstest/pkg/logger"
)

func OpentracingMiddleware(confOpentracing configs.Opentracing, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//var serverSpan opentracing.Span
		if confOpentracing.Enabled {
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				// Optionally record something about err here
				logger.Debug(fmt.Sprint("no RootSpan find", r.Header))
			}

			// Create the span referring to the RPC client if available.
			// If wireContext == nil, a root span will be created.
			RootSpan := opentracing.StartSpan(
				r.URL.Path,
				opentracing.ChildOf(wireContext))

			defer RootSpan.Finish()

			ctx := opentracing.ContextWithSpan(context.Background(), RootSpan)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
