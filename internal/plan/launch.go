package plan

import "github.com/sdmrf/relay/pkg/config"

// LaunchPlan is an immutable plan for launching a product.
// Contains only what's needed to launch - knows nothing about install logic.
type LaunchPlan struct {
	Product      string
	Version      string
	Paths        Paths
	JVMArgs      []string
	JavaMin      int
	JavaStrategy config.JavaStrategy // How to resolve Java (auto/system/bundled)
}

func (p LaunchPlan) Kind() Kind {
	return Launch
}
