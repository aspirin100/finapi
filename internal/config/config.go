package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Hostname    string `env:"FINAPI_HOSTNAME" env-default:"localhost"`
	Port        string `env:"FINAPI_PORT" env-default:"8080"`
	PostgresDSN string `env:"FINAPI_POSTGRES_DSN" env-default:"postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read env for config: %w", err)
	}

	return &cfg, nil
}
