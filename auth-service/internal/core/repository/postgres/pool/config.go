package postgrespool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration for PostgreSQL connection pool.
// Values are populated from environment variables with the POSTGRES_ prefix.
//
// Required environment variables:
//
//	POSTGRES_HOST     - database host (e.g., "localhost" or "postgres")
//	POSTGRES_PORT     - database port (default: "5432")
//	POSTGRES_USER     - database user
//	POSTGRES_PASSWORD - database password
//	POSTGRES_DB       - database name
//	POSTGRES_TIMEOUT  - operation timeout (e.g., "5s", "30s")
type Config struct {
	// Host is the PostgreSQL server address.
	// In Docker: "postgres" (service name)
	// Local: "localhost"
	Host string `envconfig:"HOST" required:"true"`

	// Port is the PostgreSQL server port.
	// Default: "5432" (standard PostgreSQL port)
	Port string `envconfig:"PORT" default:"5432"`

	// User is the database user for authentication.
	User string `envconfig:"USER" required:"true"`

	// Password is the database password for authentication.
	Password string `envconfig:"PASSWORD" required:"true"`

	// Database is the database name to connect to.
	Database string `envconfig:"DB" required:"true"`

	// Timeout is the maximum duration for database operations.
	// Format: duration string (e.g., "5s", "30s", "1m")
	Timeout time.Duration `envconfig:"TIMEOUT" required:"true"`
}

// NewConfig reads PostgreSQL configuration from environment variables
// with the POSTGRES_ prefix. Returns an error if required variables
// are missing.
//
// Usage:
//
//	cfg, err := postgrespool.NewConfig()
//	if err != nil {
//	    // handle configuration error
//	}
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust is like NewConfig but panics if the configuration is invalid.
// Use only during application startup for critical configuration.
//
// Usage:
//
//	cfg := postgrespool.NewConfigMust()
//	pool, _ := postgrespool.NewConnectionPool(ctx, cfg)
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get Postgres connection pool config: %w", err))
	}

	return config
}
