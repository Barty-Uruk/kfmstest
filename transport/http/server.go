package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Krew-Guru/kas/pkg/logger"
	"github.com/Krew-Guru/kfms/configs"
)

func RunServer(ctx context.Context, logger *log.Logger, quit chan bool, handler http.Handler, confServer configs.Http, confOpentracing configs.Opentracing) {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", confServer.Port),
		Handler: handler,
	}

	go func() {
		logger.Info("msg", "Starting HTTP/REST gateway server... listen: "+fmt.Sprintf(":%v", confServer.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("err", fmt.Sprintf("Listen HTTP/REST gateway server %s", err.Error()))
		}

		time.Sleep(time.Second * 1)
		logger.Info("msg", "Gracefull stoped HTTP/REST gateway server")
		quit <- true
	}()

	go shutDown(ctx, logger, srv)
}

// shutDown shutdown HTTP server
func shutDown(ctx context.Context, logger *log.Logger, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("msg", "Shutting down HTTP/REST gateway server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("err", fmt.Sprintf("Shutdown HTTP/REST gateway server: %s", err.Error()))
	}

	logger.Info("msg", "Shutdown done HTTP/REST gateway server")
}
