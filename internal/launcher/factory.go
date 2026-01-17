package launcher

import (
	"fmt"
	"runtime"

	"github.com/sdmrf/relay/internal/plan"
)

// New returns the appropriate launcher generator for the current OS.
func New(p plan.LaunchPlan, javaPath string) (Generator, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		return ShellLauncher{BinDir: p.Paths.BinDir, JavaPath: javaPath}, nil
	case "windows":
		return PowerShellLauncher{BinDir: p.Paths.BinDir, JavaPath: javaPath}, nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
