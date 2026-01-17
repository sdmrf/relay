package plan

import "github.com/sdmrf/relay/pkg/config"

// Artifact represents a downloadable artifact.
type Artifact struct {
	Name   string // Display name (e.g., "burpsuite.jar")
	URL    string // Download URL
	Target string // Target file path
}

// InstallPlan is an immutable plan for installing a product.
// All fields are required and fully resolved before execution.
type InstallPlan struct {
	Product  string
	Edition  string
	Version  string
	Paths    Paths
	JavaMin  int
	JVMArgs  []string
	Layout   config.LayoutMode
	Artifact Artifact
}

func (p InstallPlan) Kind() Kind {
	return Install
}
