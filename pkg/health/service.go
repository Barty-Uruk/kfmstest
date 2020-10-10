package health

import (
	"context"
	"errors"

	log "github.com/Krew-Guru/kas/pkg/logger"
	"github.com/Krew-Guru/kas/pkg/version"
)

var ErrServiceNotReady = errors.New("service not ready")

type Service interface {
	GetLiveness(ctx context.Context) error
	GetReadiness(ctx context.Context) error
	GetVersion(ctx context.Context) (Version, error)
}

func NewService(logger *log.Logger) Service {
	svc := &service{logger, false}
	go svc.load()
	return svc
}

type service struct {
	logger *log.Logger
	// app is ready if cache loaded
	ready bool
}

type Version struct {
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	Version   string `json:"version"`
}

func (s *service) GetLiveness(ctx context.Context) error {
	return nil
}

func (s *service) GetReadiness(ctx context.Context) error {
	if !s.ready {
		return ErrServiceNotReady
	}
	return nil
}

func (s *service) GetVersion(ctx context.Context) (Version, error) {
	return Version{
		version.BuildTime,
		version.Commit,
		version.Version,
	}, nil
}

func (s *service) load() {
	s.logger.Info("service", "health", "status", "init")
	s.ready = true
}
