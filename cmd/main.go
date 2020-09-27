package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/Barty-Uruk/kfmstest/pkg/logger"
	httpTransport "github.com/Barty-Uruk/kfmstest/pkg/transport/http"
)

var AppIsReady = false
var AppIsReadyMutex sync.RWMutex

func main() {

	ctx := context.Background()

	var (
		logger = logger.NewLogger("debug", "2006-01-02T15:04:05.999999999Z07:00")
	)

	srv := NewService("hello", logger.With("svc", "hello service"))

	var m http.Handler
	{
		m = NewHttpHandler(logger.With("component", "server"), srv)
	}

	quit := make(chan bool)
	httpTransport.RunServer(logger.With("component", "RunServer"), quit, m, cfg.Servers.Main.Http, cfg.Opentracing)

	logger.Info("msg", "Application is ready")

	<-quit

}
