package config

import "time"

// Defaults returns the baseline service configuration.
func Defaults() *Config {
	return &Config{
		Addr:             ":8080",
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		IdleTimeout:      60 * time.Second,
		ForceClearSecret: "railblock-dev-secret",
		MaxBodyBytes:     1 << 20,
		EnableCORS:       true,
		LogRequests:      false,
	}
}
