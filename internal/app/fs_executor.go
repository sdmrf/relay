package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sdmrf/relay/internal/downloader"
	"github.com/sdmrf/relay/internal/launcher"
	"github.com/sdmrf/relay/internal/plan"
	"github.com/sdmrf/relay/internal/runtime"
)

// FSExecutor executes plans by performing filesystem operations.
type FSExecutor struct {
	DryRun bool
}

// Execute dispatches to the appropriate handler based on plan type.
func (e FSExecutor) Execute(ctx context.Context, p plan.Plan) error {
	switch p := p.(type) {
	case plan.InstallPlan:
		return e.execInstall(ctx, p)
	case plan.RemovePlan:
		return e.execRemove(p)
	case plan.LaunchPlan:
		return e.execLaunch(ctx, p)
	case plan.UpdatePlan:
		return e.execUpdate(ctx, p)
	default:
		return fmt.Errorf("unsupported plan kind: %s", p.Kind())
	}
}

// execInstall creates required directories and downloads artifacts.
// Uses MkdirAll for idempotency - safe to run multiple times.
func (e FSExecutor) execInstall(ctx context.Context, p plan.InstallPlan) error {
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

	artifact := downloader.Artifact{
		Name:   p.Artifact.Name,
		URL:    p.Artifact.URL,
		Target: p.Artifact.Target,
	}

	if e.DryRun {
		fmt.Println("[dry-run] download:", artifact.Name)
		fmt.Println("[dry-run]   url:", artifact.URL)
		fmt.Println("[dry-run]   target:", artifact.Target)
		return nil
	}

	dl := downloader.HTTPDownloader{
		Timeout: 5 * time.Minute,
		Retries: 3,
	}

	return dl.Fetch(ctx, artifact)
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

// execLaunch validates Java, generates the launcher, and runs it.
func (e FSExecutor) execLaunch(ctx context.Context, p plan.LaunchPlan) error {
	gen, err := launcher.New(p)
	if err != nil {
		return fmt.Errorf("create launcher: %w", err)
	}

	if e.DryRun {
		fmt.Printf("[dry-run] validate java %d+\n", p.JavaMin)
		fmt.Println("[dry-run] generate launcher:", gen.Path())
		fmt.Println("[dry-run] run launcher:", gen.Path())
		return nil
	}

	// Validate Java version
	if _, err := runtime.ValidateJava(p.JavaMin); err != nil {
		return fmt.Errorf("java validation: %w", err)
	}

	// Generate the launcher script
	if err := gen.Generate(p); err != nil {
		return fmt.Errorf("generate launcher: %w", err)
	}

	// Run the launcher
	runner := runtime.ExecRunner{}
	return runner.Run(ctx, gen.Path())
}

// execUpdate downloads the new version, replacing the existing JAR.
func (e FSExecutor) execUpdate(ctx context.Context, p plan.UpdatePlan) error {
	artifact := downloader.Artifact{
		Name:   p.Artifact.Name,
		URL:    p.Artifact.URL,
		Target: p.Artifact.Target,
	}

	if e.DryRun {
		fmt.Printf("[dry-run] update %s -> %s\n", p.CurrentVersion, p.TargetVersion)
		fmt.Println("[dry-run] download:", artifact.Name)
		fmt.Println("[dry-run]   url:", artifact.URL)
		fmt.Println("[dry-run]   target:", artifact.Target)
		return nil
	}

	dl := downloader.HTTPDownloader{
		Timeout: 5 * time.Minute,
		Retries: 3,
	}

	return dl.Fetch(ctx, artifact)
}
