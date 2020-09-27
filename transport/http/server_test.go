package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/Krew-Guru/kas/configs"
)

func TestRunServer(t *testing.T) {
	type args struct {
		ctx             context.Context
		quit            chan bool
		router          *gin.Engine
		confServer      configs.Server
		confOpentracing configs.Opentracing
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunServer(tt.args.ctx, tt.args.quit, tt.args.router, tt.args.confServer, tt.args.confOpentracing)
		})
	}
}

func Test_shutDown(t *testing.T) {
	type args struct {
		ctx context.Context
		srv *http.Server
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutDown(tt.args.ctx, tt.args.srv)
		})
	}
}
