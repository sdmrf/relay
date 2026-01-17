package app

import (
	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/internal/plan"
	"github.com/sdmrf/relay/pkg/config"
)

// ResolveInstall creates an InstallPlan from config and resolved paths.
// Pure function with no side effects - no filesystem or environment access.
func ResolveInstall(cfg config.Config, p paths.Paths) plan.InstallPlan {
	return plan.InstallPlan{
		Product: cfg.Product.Name,
		Edition: cfg.Product.Edition,
		Version: cfg.Product.Version,
		Paths:   plan.FromResolved(p),
		JavaMin: cfg.Runtime.Java.MinVersion,
		JVMArgs: cfg.Runtime.Java.JVMArgs,
		Layout:  cfg.Layout.Mode,
	}
}

// ResolveLaunch creates a LaunchPlan from config and resolved paths.
// Pure function with no side effects - no filesystem or environment access.
func ResolveLaunch(cfg config.Config, p paths.Paths) plan.LaunchPlan {
	return plan.LaunchPlan{
		Product: cfg.Product.Name,
		Version: cfg.Product.Version,
		Paths:   plan.FromResolved(p),
		JVMArgs: cfg.Runtime.Java.JVMArgs,
	}
}

// ResolveRemove creates a RemovePlan from config and resolved paths.
// Pure function with no side effects - no filesystem or environment access.
func ResolveRemove(cfg config.Config, p paths.Paths) plan.RemovePlan {
	return plan.RemovePlan{
		Product: cfg.Product.Name,
		Paths:   plan.FromResolved(p),
	}
}
