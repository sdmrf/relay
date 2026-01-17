package launcher

import "github.com/sdmrf/relay/internal/plan"

// Generator creates platform-specific launchers.
type Generator interface {
	Generate(plan.LaunchPlan) error
	Path() string
}
