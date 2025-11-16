package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "github.com/ocrosby/go-lab/projects/api/docs"
	"github.com/ocrosby/go-lab/projects/api/internal/config"
	"github.com/ocrosby/go-lab/projects/api/internal/di"
	httpAdapter "github.com/ocrosby/go-lab/projects/api/internal/infrastructure/adapters/http"
	"github.com/ocrosby/go-lab/projects/api/pkg/health"
)

// @title User Management API
// @version 1.0
// @description A comprehensive Go API example demonstrating hexagonal architecture and dependency injection
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "A sample Go API demonstrating hexagonal architecture",
	Long: `A comprehensive Go API example that demonstrates:
- Hexagonal architecture (ports and adapters)
- Dependency injection with uber-go/dig
- Kubernetes health probes
- ServeMux for HTTP routing
- Gang of Four design patterns
- Low cyclomatic complexity`,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API server",
	RunE:  runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	container := di.NewContainer()
	if err := container.BuildContainer(); err != nil {
		return fmt.Errorf("failed to build container: %w", err)
	}

	var server *httpAdapter.Server
	var logger *zap.Logger
	var healthChecker health.HealthChecker
	var cfg *config.Config

	if err := container.Invoke(func(
		s *httpAdapter.Server,
		l *zap.Logger,
		h health.HealthChecker,
		c *config.Config,
	) {
		server = s
		logger = l
		healthChecker = h
		cfg = c
	}); err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	healthChecker.AddCheck("basic", func(ctx context.Context) error {
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutdown signal received")
		cancel()
	}()

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("Server error", zap.Error(err))
			cancel()
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.GetShutdownTimeout())
	defer shutdownCancel()

	return server.Stop(shutdownCtx)
}
