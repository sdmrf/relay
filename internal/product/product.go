package product

import "github.com/sdmrf/relay/internal/plan"

// Product defines the contract for product-specific behavior.
// Implementations provide resolution logic but no filesystem or network access.
type Product interface {
	Name() string
	ResolveInstall() (plan.InstallPlan, error)
	ResolveLaunch() (plan.LaunchPlan, error)
}
