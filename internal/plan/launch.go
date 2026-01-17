package plan

// LaunchPlan is an immutable plan for launching a product.
// Contains only what's needed to launch - knows nothing about install logic.
type LaunchPlan struct {
	Product string
	Version string
	Paths   Paths
	JVMArgs []string
}

func (p LaunchPlan) Kind() Kind {
	return Launch
}
