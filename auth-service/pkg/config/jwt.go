package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// JWTConfig holds JWT configuration from environment variables.
type JWTConfig struct {
	Secret            string        `envconfig:"SECRET" required:"true"`
	AccessExpiration  time.Duration `envconfig:"ACCESS_EXPIRATION" default:"15m"`
	RefreshExpiration time.Duration `envconfig:"REFRESH_EXPIRATION" default:"168h"`
}

// NewJWTConfig reads JWT configuration from environment variables
// with the JWT_ prefix.
func NewJWTConfig() (JWTConfig, error) {
	var config JWTConfig

	if err := envconfig.Process("JWT", &config); err != nil {
		return JWTConfig{}, fmt.Errorf("Process envconfig: %w", err)
	}

	return config, nil
}

func NewJWTConfigMust() JWTConfig {
	config, err := NewJWTConfig()

	if err != nil {
		panic(fmt.Errorf("Get JWT config: %w", err))
	}
	
	return config
}