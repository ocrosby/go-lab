// Package http provides HTTP adapter implementation for the User Management API.
package http

import (
	"context"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/config"
	"github.com/ocrosby/go-lab/projects/api/pkg/health"
)

type Server struct {
	httpServer   *http.Server
	healthServer *http.Server
	logger       *zap.Logger
	config       *config.Config
}

func NewServer(
	userHandler *UserHandler,
	healthChecker health.HealthChecker,
	logger *zap.Logger,
	cfg *config.Config,
) *Server {
	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", health.LivenessHandler(healthChecker))
	healthMux.HandleFunc("/readyz", health.ReadinessHandler(healthChecker))
	healthMux.HandleFunc("/startupz", health.StartupHandler(healthChecker))

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Handler:      mux,
			ReadTimeout:  cfg.GetReadTimeout(),
			WriteTimeout: cfg.GetWriteTimeout(),
			IdleTimeout:  cfg.GetIdleTimeout(),
		},
		healthServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Health.Port),
			Handler:      healthMux,
			ReadTimeout:  cfg.GetHealthReadTimeout(),
			WriteTimeout: cfg.GetHealthWriteTimeout(),
			IdleTimeout:  cfg.GetHealthIdleTimeout(),
		},
		logger: logger,
		config: cfg,
	}
}

func (s *Server) Start() error {
	go func() {
		s.logger.Info("Starting health server",
			zap.String("addr", s.healthServer.Addr))
		if err := s.healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Health server error", zap.Error(err))
		}
	}()

	s.logger.Info("Starting HTTP server",
		zap.String("addr", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping servers")

	if err := s.healthServer.Shutdown(ctx); err != nil {
		s.logger.Error("Error stopping health server", zap.Error(err))
	}

	return s.httpServer.Shutdown(ctx)
}
