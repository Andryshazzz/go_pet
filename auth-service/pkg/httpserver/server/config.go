package httpserver

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration for the HTTP server.
//
// Required environment variables:
//
//	HTTP_ADDR              - server address to listen on (e.g., ":5050")
//	HTTP_SHUTDOWN_TIMEOUT  - maximum duration to wait for graceful shutdown
type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}

// NewConfig reads HTTP server configuration from environment variables
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("{Process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust is like NewConfig but panics if the configuration is invalid.
// Use only during application startup for critical configuration
// that must be valid for the app to run.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("Get HTTP server config: %w", err)

		panic(err)
	}

	return config
}
