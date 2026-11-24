package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("config: DATABASE_URL wajib diisi")
	}
	return cfg, nil
}
