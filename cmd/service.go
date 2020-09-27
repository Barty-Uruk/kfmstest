package service

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kit/kit/endpoint"
)

type Service interface {
	SayHello(ctx context.Context, name string) (string, error)
}
type service struct {
	HelloWord string
	Logger    *log.Logger
}

func (s *service) SayHello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("%s, %s", s.HelloWord, name), nil
}
func NewService(word string, log *log.Logger) Service {
	return &service{
		Logger:    log,
		HelloWord: word,
	}
}

type HelloRequest struct {
	Name string `json:"name"`
}
type HelloResponse struct {
	HelloText string `json:"hello_text"`
}

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Hello endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Hello: makeHelloEndpoint(s),
	}
}

func makeHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HelloRequest)
		text, err := s.SayHello(ctx, req.Name)
		return HelloResponse{HelloText: text}, err
	}
}
