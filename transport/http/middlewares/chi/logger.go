package chi

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	log "github.com/Krew-Guru/kas/pkg/logger"
)

var DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.Log, NoColor: false})

func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}
