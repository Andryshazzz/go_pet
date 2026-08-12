package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration for the Logger.
// Values are populated from environment variables with the LOGGER_ prefix.
//
// Required environment variables:
//
//	LOGGER_LEVEL  - log level: debug, info, warn, error
//	LOGGER_FOLDER - directory path for log files
//
// Example .env:
//
//	LOGGER_LEVEL=debug
//	LOGGER_FOLDER=/app/out/logs
type Config struct {
	Level  string `envconfig:"LEVEL"  default:"info"`
	Folder string `envconfig:"FOLDER" default:"./out/logs"`
}

// NewConfig reads Logger configuration from environment variables
// with the LOGGER_ prefix. Returns an error if required variables
// are missing or invalid.
//
// Usage:
//
//	cfg, err := logger.NewConfig()
//	if err != nil {
//	    // handle configuration error
//	}
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("Process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust is like NewConfig but panics if the configuration is invalid.
// Use this only during application startup for critical configuration
// that must be valid for the app to run.
//
// Usage:
//
//	cfg := logger.NewConfigMust()
//	log, _ := logger.NewLogger(cfg)
//
// Typical error panic message:
//
//	panic: get Logger config: process envconfig: required key LOGGER_LEVEL missing value
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("Get Logger config: %w", err)
		panic(err)
	}

	return config
}
