package plan

import "github.com/sdmrf/relay/pkg/config"

// Artifact represents a downloadable artifact.
type Artifact struct {
	Name   string // Display name (e.g., "burpsuite.jar")
	URL    string // Download URL
	Target string // Target file path
}

// JREArtifact represents a JRE to download and extract.
type JREArtifact struct {
	Name      string // Display name (e.g., "Eclipse Temurin JRE 21.0.5")
	URL       string // Download URL
	Target    string // Download target path (.tar.gz or .zip)
	ExtractTo string // Extraction destination directory
}

// InstallPlan is an immutable plan for installing a product.
// All fields are required and fully resolved before execution.
type InstallPlan struct {
	Product     string
	Edition     string
	Version     string
	Paths       Paths
	JavaMin     int
	JVMArgs     []string
	Layout      config.LayoutMode
	Artifact    Artifact
	JREArtifact *JREArtifact // Optional: nil if JRE not needed
}

func (p InstallPlan) Kind() Kind {
	return Install
}
