package plan

import "github.com/sdmrf/relay/internal/paths"

// Kind represents the lifecycle intent of a plan.
type Kind string

const (
	Install Kind = "install"
	Update  Kind = "update"
	Launch  Kind = "launch"
	Remove  Kind = "remove"
)

// Plan is implemented by all plan types.
type Plan interface {
	Kind() Kind
}

// Paths contains resolved filesystem paths for execution.
type Paths struct {
	InstallDir string
	DataDir    string
	BinDir     string
	ConfigDir  string
	CacheDir   string
}

// FromResolved converts paths.Paths to plan.Paths.
func FromResolved(p paths.Paths) Paths {
	return Paths{
		InstallDir: p.InstallDir,
		DataDir:    p.DataDir,
		BinDir:     p.BinDir,
		ConfigDir:  p.ConfigDir,
		CacheDir:   p.CacheDir,
	}
}
