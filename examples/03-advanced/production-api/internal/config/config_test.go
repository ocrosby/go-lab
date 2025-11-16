package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestNewConfig_Defaults(t *testing.T) {
	// Clear any existing environment variables
	viper.Reset()

	config := NewConfig()

	if config.Server.Port != 8080 {
		t.Errorf("Expected default server port 8080, got %d", config.Server.Port)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default server host '0.0.0.0', got %s", config.Server.Host)
	}

	if config.Health.Port != 8081 {
		t.Errorf("Expected default health port 8081, got %d", config.Health.Port)
	}
}

func TestNewConfig_EnvironmentVariables(t *testing.T) {
	// Clear any existing configuration
	viper.Reset()

	// Set environment variables - viper expects uppercase with underscores
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("HEALTH_PORT", "9001")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("HEALTH_PORT")
		viper.Reset()
	}()

	// Set viper to use env vars and configure key binding
	viper.SetEnvKeyReplacer(nil)
	viper.AutomaticEnv()
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.host", "SERVER_HOST")
	_ = viper.BindEnv("health.port", "HEALTH_PORT")

	config := NewConfig()

	if config.Server.Port != 9000 {
		t.Errorf("Expected server port 9000, got %d", config.Server.Port)
	}

	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host 'localhost', got %s", config.Server.Host)
	}

	if config.Health.Port != 9001 {
		t.Errorf("Expected health port 9001, got %d", config.Health.Port)
	}
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectValid bool
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "0.0.0.0",
				},
				Health: HealthConfig{
					Port: 8081,
				},
			},
			expectValid: true,
		},
		{
			name: "invalid server port (too low)",
			config: Config{
				Server: ServerConfig{
					Port: 0,
					Host: "0.0.0.0",
				},
				Health: HealthConfig{
					Port: 8081,
				},
			},
			expectValid: false,
		},
		{
			name: "invalid server port (too high)",
			config: Config{
				Server: ServerConfig{
					Port: 65536,
					Host: "0.0.0.0",
				},
				Health: HealthConfig{
					Port: 8081,
				},
			},
			expectValid: false,
		},
		{
			name: "empty host",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "",
				},
				Health: HealthConfig{
					Port: 8081,
				},
			},
			expectValid: false,
		},
		{
			name: "invalid health port",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "0.0.0.0",
				},
				Health: HealthConfig{
					Port: 0,
				},
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			isValid := err == nil

			if isValid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v, error=%v", tt.expectValid, isValid, err)
			}
		})
	}
}

func TestConfig_GetServerAddress(t *testing.T) {
	config := Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	expected := "localhost:8080"
	actual := config.GetServerAddress()

	if actual != expected {
		t.Errorf("Expected server address %s, got %s", expected, actual)
	}
}

func TestConfig_GetHealthAddress(t *testing.T) {
	config := Config{
		Server: ServerConfig{
			Host: "localhost",
		},
		Health: HealthConfig{
			Port: 8081,
		},
	}

	expected := "localhost:8081"
	actual := config.GetHealthAddress()

	if actual != expected {
		t.Errorf("Expected health address %s, got %s", expected, actual)
	}
}
