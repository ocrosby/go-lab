// Package di provides dependency injection container setup for the application.
package di

import (
	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/application"
	"github.com/ocrosby/go-lab/projects/api/internal/config"
	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	httpAdapter "github.com/ocrosby/go-lab/projects/api/internal/infrastructure/adapters/http"
	"github.com/ocrosby/go-lab/projects/api/internal/infrastructure/adapters/repository"
	"github.com/ocrosby/go-lab/projects/api/pkg/health"
)

type Container struct {
	*dig.Container
}

func NewContainer() *Container {
	container := dig.New()
	return &Container{Container: container}
}

func (c *Container) BuildContainer() error {
	providers := []interface{}{
		c.provideLogger,
		c.provideConfig,
		c.provideHealthChecker,
		c.provideUserRepository,
		c.provideUserService,
		c.provideUserHandler,
		c.provideServer,
	}

	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) provideLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func (c *Container) provideConfig() *config.Config {
	return config.NewConfig()
}

func (c *Container) provideHealthChecker() health.HealthChecker {
	return health.NewHealthChecker()
}

func (c *Container) provideUserRepository() domain.UserRepository {
	return repository.NewMemoryUserRepository()
}

func (c *Container) provideUserService(
	userRepo domain.UserRepository,
	logger *zap.Logger,
) domain.UserService {
	return application.NewUserService(userRepo, logger)
}

func (c *Container) provideUserHandler(
	userService domain.UserService,
	logger *zap.Logger,
) *httpAdapter.UserHandler {
	return httpAdapter.NewUserHandler(userService, logger)
}

func (c *Container) provideServer(
	userHandler *httpAdapter.UserHandler,
	healthChecker health.HealthChecker,
	logger *zap.Logger,
	cfg *config.Config,
) *httpAdapter.Server {
	return httpAdapter.NewServer(userHandler, healthChecker, logger, cfg)
}
