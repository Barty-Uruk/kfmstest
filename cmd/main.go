package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Barty-Uruk/kfmstest/cmd/service"
	"github.com/Barty-Uruk/kfmstest/configs"
	"github.com/Barty-Uruk/kfmstest/pkg/logger"
	httpTransport "github.com/Barty-Uruk/kfmstest/transport/http"
)

func main() {

	ctx := context.Background()

	var (
		logger = logger.NewLogger("debug", "2006-01-02T15:04:05.999999999Z07:00")
	)
	// get configuration
	cfg := configs.NewConfig()
	if err := cfg.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}

	// print config
	if err := cfg.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "print config: %s", err)
		os.Exit(1)
	}
	srv, err := service.NewService(logger.With("svc", "amazon service"), &cfg.S3)
	if err != nil {
		logger.Fatal("err", fmt.Sprintf("s3 service creation error: %s", err.Error()))
		os.Exit(1)
	}
	var m http.Handler
	{
		m = service.NewHttpHandler(logger.With("component", "server"), srv)
	}

	quit := make(chan bool)
	httpTransport.RunServer(ctx, logger.With("component", "RunServer"), quit, m, cfg.Servers.Main.Http, cfg.Opentracing)

	logger.Info("msg", "Application is ready")

	<-quit

}
