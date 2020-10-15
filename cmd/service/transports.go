package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	r.Method("GET", "/hello", kithttp.NewServer(
		makeHelloEndpoint(s),
		decodeHelloRequest,
		encodeHelloResponse,
		options...,
	))
	r.Method("POST", "/upload/", kithttp.NewServer(
		makeUploadFileEndpoint(s),
		decodeUploadFileRequest,
		encodeUploadFileResponse,
		options...,
	))

	logger.Info("msg", "handler started")
	return r
}

func encodeTokenValidationResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// login
func decodeHelloRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request HelloRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrInvalidArgument
	}
	if request.Name == "" {
		return nil, ErrInvalidArgument
	}

	return request, nil
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
func encodeHelloResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	resp, ok := response.(HelloResponse)
	if !ok {
		encodeErrorResponse(ctx, ErrTokenConvertationError, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	//log.Debug("response", fmt.Sprintf("%#v", response))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
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

func decodeUploadFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	file, info, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fmt.Println(file, "=========", info.Size)

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
