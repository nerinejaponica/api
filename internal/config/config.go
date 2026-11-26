package config

import (
	"errors"
	"os"

	"github.com/nerinejaponica/api/internal/constant"
)

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	// GET CONFIGURATION FROM ENVIRONMENT
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	// WHEN REQUIRED CONFIGURATION IS EMPTY
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New(constant.MsgDatabaseURLRequired)
	}
	return cfg, nil
}
