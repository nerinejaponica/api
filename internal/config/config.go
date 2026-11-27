package config

import (
	"errors"
	"os"
	"strings"

	"github.com/nerinejaponica/api/internal/constant"
)

type Config struct {
	DatabaseURL    string
	AllowedOrigins []string
}

func Load() (Config, error) {
	// GET CONFIGURATION FROM ENVIRONMENT
	cfg := Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AllowedOrigins: splitList(os.Getenv("CORS_ALLOWED_ORIGINS")),
	}
	// WHEN REQUIRED CONFIGURATION IS EMPTY
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New(constant.MsgDatabaseURLRequired)
	}
	return cfg, nil
}

func splitList(raw string) []string {
	// SPLIT COMMA SEPARATED VALUE AND DROP EMPTY ITEMS
	var out []string
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}
