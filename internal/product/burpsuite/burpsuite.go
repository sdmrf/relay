package burpsuite

import (
	"fmt"

	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/pkg/config"
)

// BurpSuite encapsulates all Burp-specific resolution logic.
// No filesystem access, no network access - pure data transformation.
type BurpSuite struct {
	cfg   config.Config
	paths paths.Paths
}

// New creates a BurpSuite product instance.
// Returns error if config is not for burpsuite.
func New(cfg config.Config, p paths.Paths) (*BurpSuite, error) {
	if cfg.Product.Name != "burpsuite" {
		return nil, fmt.Errorf("invalid product: expected burpsuite, got %s", cfg.Product.Name)
	}

	return &BurpSuite{
		cfg:   cfg,
		paths: p,
	}, nil
}

func (b *BurpSuite) Name() string {
	return "burpsuite"
}
