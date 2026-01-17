package plan

// RemovePlan is an immutable plan for uninstalling a product.
// Minimal struct - can only act on owned paths (enforced during execution).
type RemovePlan struct {
	Product string
	Paths   Paths
}

func (p RemovePlan) Kind() Kind {
	return Remove
}
