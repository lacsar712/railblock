package config

import (
	"fmt"
	"strings"
)

// Validate checks that configuration values are usable.
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("listen address is required")
	}
	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}
	if c.IdleTimeout <= 0 {
		return fmt.Errorf("idle timeout must be positive")
	}
	if len(c.ForceClearSecret) < 8 {
		return fmt.Errorf("force clear secret must be at least 8 bytes")
	}
	if c.MaxBodyBytes <= 0 {
		return fmt.Errorf("max body bytes must be positive")
	}
	return nil
}
