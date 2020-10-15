package health

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeGetLivenessEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.GetLiveness(ctx)
		return GetLivenessResponse{"ok", err}, nil
	}
}

func MakeGetReadinessEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.GetReadiness(ctx)
		return GetReadinessResponse{"ok", err}, nil
	}
}

func MakeGetVersionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := s.GetVersion(ctx)
		return GetVersionResponse{resp, err}, nil
	}
}

type GetLivenessRequest struct{}

type GetLivenessResponse struct {
	Status string `json:"status"`
	Err    error  `json:"error,omitempty"`
}

func (r GetLivenessResponse) Error() error { return r.Err }

type GetReadinessRequest struct{}

type GetReadinessResponse struct {
	Status string `json:"status"`
	Err    error  `json:"error,omitempty"`
}

func (r GetReadinessResponse) Error() error { return r.Err }

type GetVersionRequest struct{}

type GetVersionResponse struct {
	Version Version `json:"version"`
	Err     error   `json:"error,omitempty"`
}

func (r GetVersionResponse) Error() error { return r.Err }
