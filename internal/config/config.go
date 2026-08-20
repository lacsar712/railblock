package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds runtime settings for the railblock service.
type Config struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ForceClearSecret string
	MaxBodyBytes    int64
	EnableCORS      bool
	LogRequests     bool
}

// Load reads configuration from environment variables with defaults.
func Load() (*Config, error) {
	cfg := Defaults()
	if v := os.Getenv("RAILBLOCK_ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("RAILBLOCK_FORCE_SECRET"); v != "" {
		cfg.ForceClearSecret = v
	}
	if v := os.Getenv("RAILBLOCK_MAX_BODY"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RAILBLOCK_MAX_BODY: %w", err)
		}
		cfg.MaxBodyBytes = n
	}
	if v := os.Getenv("RAILBLOCK_CORS"); strings.EqualFold(v, "true") || v == "1" {
		cfg.EnableCORS = true
	}
	if v := os.Getenv("RAILBLOCK_LOG"); strings.EqualFold(v, "true") || v == "1" {
		cfg.LogRequests = true
	}
	return cfg, nil
}

// MustLoad loads configuration or panics (for main).
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	return cfg
}
