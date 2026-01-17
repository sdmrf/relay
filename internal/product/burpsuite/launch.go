package burpsuite

import "github.com/sdmrf/relay/internal/plan"

// ResolveLaunch creates an immutable LaunchPlan for Burp Suite.
// Pure function - no filesystem or network access.
func (b *BurpSuite) ResolveLaunch() (plan.LaunchPlan, error) {
	return plan.LaunchPlan{
		Product:      b.Name(),
		Version:      b.cfg.Product.Version,
		Paths:        plan.FromResolved(b.paths),
		JVMArgs:      b.cfg.Runtime.Java.JVMArgs,
		JavaMin:      b.cfg.Runtime.Java.MinVersion,
		JavaStrategy: b.cfg.Runtime.Java.Strategy,
	}, nil
}
