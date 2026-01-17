package plan

import "github.com/sdmrf/relay/pkg/config"

// InstallPlan is an immutable plan for installing a product.
// All fields are required and fully resolved before execution.
type InstallPlan struct {
	Product string
	Edition string
	Version string
	Paths   Paths
	JavaMin int
	JVMArgs []string
	Layout  config.LayoutMode
}

func (p InstallPlan) Kind() Kind {
	return Install
}
