package app

import (
	"fmt"
	"os"

	"github.com/sdmrf/relay/internal/plan"
)

// FSExecutor executes plans by performing filesystem operations.
type FSExecutor struct {
	DryRun bool
}

// Execute dispatches to the appropriate handler based on plan type.
func (e FSExecutor) Execute(p plan.Plan) error {
	switch p := p.(type) {
	case plan.InstallPlan:
		return e.execInstall(p)
	case plan.RemovePlan:
		return e.execRemove(p)
	case plan.LaunchPlan:
		return fmt.Errorf("launch execution not yet implemented")
	default:
		return fmt.Errorf("unsupported plan kind: %s", p.Kind())
	}
}

// execInstall creates required directories for installation.
// Uses MkdirAll for idempotency - safe to run multiple times.
func (e FSExecutor) execInstall(p plan.InstallPlan) error {
	dirs := []string{
		p.Paths.InstallDir,
		p.Paths.DataDir,
		p.Paths.BinDir,
		p.Paths.CacheDir,
	}

	for _, dir := range dirs {
		if e.DryRun {
			fmt.Println("[dry-run] mkdir:", dir)
			continue
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	return nil
}

// execRemove deletes only owned paths.
// Preserves ConfigDir to retain user configuration.
func (e FSExecutor) execRemove(p plan.RemovePlan) error {
	for _, dir := range p.Paths.Owned() {
		if e.DryRun {
			fmt.Println("[dry-run] rm -rf:", dir)
			continue
		}

		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("remove directory %s: %w", dir, err)
		}
	}

	return nil
}
