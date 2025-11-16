// Package config provides configuration management for the User Management API.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Health HealthConfig `mapstructure:"health"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type HealthConfig struct {
	Port int `mapstructure:"port"`
}

func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if c.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	if c.Health.Port < 1 || c.Health.Port > 65535 {
		return fmt.Errorf("health port must be between 1 and 65535")
	}
	return nil
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetHealthAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Health.Port)
}

// Timeout configuration methods

func (c *Config) GetReadTimeout() time.Duration {
	return DefaultReadTimeout
}

func (c *Config) GetWriteTimeout() time.Duration {
	return DefaultWriteTimeout
}

func (c *Config) GetIdleTimeout() time.Duration {
	return DefaultIdleTimeout
}

func (c *Config) GetHealthReadTimeout() time.Duration {
	return DefaultHealthReadTimeout
}

func (c *Config) GetHealthWriteTimeout() time.Duration {
	return DefaultHealthWriteTimeout
}

func (c *Config) GetHealthIdleTimeout() time.Duration {
	return DefaultHealthIdleTimeout
}

func (c *Config) GetHealthCheckTimeout() time.Duration {
	return DefaultHealthCheckTimeout
}

func (c *Config) GetShutdownTimeout() time.Duration {
	return DefaultShutdownTimeout
}

func NewConfig() *Config {
	viper.SetDefault("server.port", DefaultServerPort)
	viper.SetDefault("server.host", DefaultServerHost)
	viper.SetDefault("health.port", DefaultHealthPort)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
