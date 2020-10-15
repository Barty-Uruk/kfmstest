package health

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

type loggingService struct {
	logger log.Logger
	Service
}

func (s *loggingService) GetLiveness(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", codeFrom(err),
			"method", "GetLiveness",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetLiveness(ctx)
}

func (s *loggingService) GetReadiness(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", codeFrom(err),
			"method", "GetReadiness",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetReadiness(ctx)
}

func (s *loggingService) GetVersion(ctx context.Context) (v Version, err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", codeFrom(err),
			"method", "GetVersion",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetVersion(ctx)
}

func codeFrom(err error) int {
	switch err {
	case ErrServiceNotReady:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
