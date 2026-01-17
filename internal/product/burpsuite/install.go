package burpsuite

import (
	"path/filepath"

	"github.com/sdmrf/relay/internal/plan"
)

// ResolveInstall creates an immutable InstallPlan for Burp Suite.
// Pure function - no filesystem or network access.
func (b *BurpSuite) ResolveInstall() (plan.InstallPlan, error) {
	return plan.InstallPlan{
		Product: b.Name(),
		Edition: b.cfg.Product.Edition,
		Version: b.cfg.Product.Version,
		Paths:   plan.FromResolved(b.paths),
		JavaMin: b.cfg.Runtime.Java.MinVersion,
		JVMArgs: b.cfg.Runtime.Java.JVMArgs,
		Layout:  b.cfg.Layout.Mode,
		Artifact: plan.Artifact{
			Name:   JarName,
			URL:    burpDownloadURL(b.cfg.Product.Edition),
			Target: filepath.Join(b.paths.InstallDir, JarName),
		},
	}, nil
}
