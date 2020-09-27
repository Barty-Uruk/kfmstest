package chi

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	log "github.com/Barty-Uruk/kfmstest/pkg/logger"
)

var DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.Log, NoColor: false})

func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}
