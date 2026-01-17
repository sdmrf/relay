package config

import (
	"fmt"
)

func (c Config) Validate() error {
	if c.Product.Name == "" {
		return fmt.Errorf("product.name is required")
	}

	if c.Product.Version == "" {
		return fmt.Errorf("product.version is required")
	}

	switch c.Layout.Mode {
	case "system", "portable":
	default:
		return fmt.Errorf("invalid layout.mode: %s", c.Layout.Mode)
	}

	switch c.Runtime.Java.Strategy {
	case "auto", "system":
	default:
		return fmt.Errorf("invalid runtime.java.strategy: %s", c.Runtime.Java.Strategy)
	}

	if c.Runtime.Java.MinVersion < 8 {
		return fmt.Errorf("runtime.java.min_version must be >= 8")
	}

	switch c.Logging.Level {
	case "info", "debug", "trace":
	default:
		return fmt.Errorf("invalid logging.level: %s", c.Logging.Level)
	}

	return nil
}
