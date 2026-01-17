package burpsuite

import "github.com/sdmrf/relay/internal/plan"

// ResolveRemove creates an immutable RemovePlan for Burp Suite.
// Pure function - no filesystem or network access.
func (b *BurpSuite) ResolveRemove() (plan.RemovePlan, error) {
	return plan.RemovePlan{
		Product: b.Name(),
		Paths:   plan.FromResolved(b.paths),
	}, nil
}
