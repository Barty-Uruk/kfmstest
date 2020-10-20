package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Barty-Uruk/kfmstest/pkg/health"
	log "github.com/Barty-Uruk/kfmstest/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	kithttp "github.com/go-kit/kit/transport/http"
)

var (
	ErrPageNotFound           = errors.New("page not found")
	ErrRequestInvalid         = errors.New("request invalid")
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrTokenInvalidError      = errors.New("Token invalid error")
	ErrTokenConvertationError = errors.New("Token convertation error")
	ErrUploadFile             = errors.New("error upload file")
	ErrUploadMaxFileSize      = errors.New("error upload file: max file size")
)

// NewService wires Go kit endpoints to the HTTP transport.
func NewHttpHandler(logger *log.Logger, s Service) http.Handler {
	// set-up router and initialize http endpoints
	r := chi.NewRouter()
	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
	}
	endpoints := MakeEndpoints(s)
	r.Method("POST", "/upload", kithttp.NewServer(
		endpoints.UploadFile,
		decodeUploadFileRequest,
		encodeUploadFileResponse,
		options...,
	))
	r.Method("GET", "/download", kithttp.NewServer(
		endpoints.DownloadFile,
		decodeDownloadFileRequest,
		encodeDownloadResponse,
		options...,
	))

	logger.Info("msg", "handler started")
	return r
}

func encodeDownloadResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}

	resp, ok := response.(DownloadResponse)
	if !ok {
		encodeErrorResponse(ctx, errors.New("error download file"), w)
		return nil
	}
	_, err := io.Copy(w, resp.File)
	return errors.Wrap(err, "copy file error")
}

func encodeUploadFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	resp, ok := response.(UploadResponse)
	if !ok {
		encodeErrorResponse(ctx, errors.New("error upload file"), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrUploadFile:
		return http.StatusBadRequest
	case health.ErrServiceNotReady:
		return http.StatusBadRequest
	case ErrInvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusUnauthorized
	}
}

func decodeDownloadFileRequest(_ context.Context, r *http.Request) (interface{}, error) {

	req := DownloadRequest{
		FileName: r.URL.Query().Get("filename"),
	}
	return req, req.validate()
}
func decodeUploadFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	file, info, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	req := UploadRequest{
		File:       file,
		FolderName: "test",
		FileName:   info.Filename,
	}
	return req, req.validate()
}

func decodeGetReadinessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return health.GetLivenessRequest{}, nil
}

func decodeGetVersionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return health.GetLivenessRequest{}, nil
}
