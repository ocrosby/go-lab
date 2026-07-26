// Package config provides configuration management for the User Management API.
package config

import "time"

// Default configuration values
const (
	// Default server configuration
	DefaultServerPort = 8080
	DefaultServerHost = "0.0.0.0"
	DefaultHealthPort = 8081

	// Default timeout values
	DefaultReadTimeout  = 30 * time.Second
	DefaultWriteTimeout = 30 * time.Second
	DefaultIdleTimeout  = 60 * time.Second

	// Health check timeout values
	DefaultHealthReadTimeout  = 10 * time.Second
	DefaultHealthWriteTimeout = 10 * time.Second
	DefaultHealthIdleTimeout  = 30 * time.Second

	// Health check context timeout
	DefaultHealthCheckTimeout = 30 * time.Second

	// Default pagination values
	DefaultPaginationLimit  = 10
	DefaultPaginationOffset = 0

	// Application shutdown timeout
	DefaultShutdownTimeout = 30 * time.Second
)
