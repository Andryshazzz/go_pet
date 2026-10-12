package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
	Logger   LoggerConfig   `yaml:"logger"`
	JWT      JWTConfig      `yaml:"jwt"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

// HTTPConfig holds HTTP server settings.
type HTTPConfig struct {
	Addr            string        `yaml:"addr" env:"HTTP_ADDR"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
}

// PostgresConfig holds PostgreSQL connection settings.
type PostgresConfig struct {
	Host     string        `yaml:"host" env:"POSTGRES_HOST"`
	Port     string        `yaml:"port" env:"POSTGRES_PORT"`
	User     string        `yaml:"user" env:"POSTGRES_USER"`
	Password string        `yaml:"password" env:"POSTGRES_PASSWORD"`
	DB       string        `yaml:"db" env:"POSTGRES_DB"`
	SSLMode  string        `yaml:"sslmode" env:"POSTGRES_SSLMODE"`
	Timeout  time.Duration `yaml:"timeout" env:"POSTGRES_TIMEOUT"`
}

// LoggerConfig holds logger settings.
type LoggerConfig struct {
	Level  string `yaml:"level" env:"LOGGER_LEVEL"`
	Folder string `yaml:"folder" env:"LOGGER_FOLDER"`
}

// JWTConfig holds JWT settings for validation.
type JWTConfig struct {
	Secret string `yaml:"secret" env:"JWT_SECRET"`
}

// KafkaConfig holds Kafka connection settings.
type KafkaConfig struct {
	Brokers            []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
	OrdersCreatedTopic string   `yaml:"orders_created_topic" env:"KAFKA_ORDERS_CREATED_TOPIC"`
	OrderStatusTopic   string   `yaml:"order_status_topic" env:"KAFKA_ORDER_STATUS_TOPIC"`
	ClientID           string   `yaml:"client_id" env:"KAFKA_CLIENT_ID"`
	ConsumerGroup      string   `yaml:"consumer_group" env:"KAFKA_CONSUMER_GROUP"`
}

// NewConfig loads configuration from YAML file and environment variables.
func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("load .env: %w", err)
		}
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return cfg, nil
}
