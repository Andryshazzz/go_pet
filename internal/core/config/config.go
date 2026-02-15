package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type AppConfig struct {
	PublicHTTPPort   string `env:"PUBLIC_HTTP_PORT" env-default:"8080"`
	InternalHTTPPort string `env:"INTERNAL_HTTP_PORT" env-default:"8081"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DBName   string `env:"POSTGRES_DB" env-default:"go_pet"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	return &cfg, err
}
