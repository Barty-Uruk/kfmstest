package chi

import (
	"fmt"
	"net/http"

	log "github.com/Barty-Uruk/kfmstest/pkg/logger"
)

func RequestLogger(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("request", fmt.Sprintf("%+v", r.Header))
			next.ServeHTTP(w, r)
			logger.Debug("response", fmt.Sprintf("%+v", w.Header()))
		}
		return http.HandlerFunc(fn)
	}
}
