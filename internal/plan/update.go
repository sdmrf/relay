package plan

// UpdatePlan is an immutable plan for updating a product.
type UpdatePlan struct {
	Product        string
	Edition        string
	CurrentVersion string
	TargetVersion  string
	Paths          Paths
	Artifact       Artifact
}

func (p UpdatePlan) Kind() Kind {
	return Update
}
