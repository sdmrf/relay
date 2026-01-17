package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sdmrf/relay/internal/downloader"
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

// execInstall creates required directories and downloads artifacts.
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

	// Download artifact
	artifact := downloader.Artifact{
		Name:   "burpsuite.jar",
		URL:    burpDownloadURL(p.Edition),
		Target: filepath.Join(p.Paths.InstallDir, "burpsuite.jar"),
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

	return dl.Fetch(context.Background(), artifact)
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

// burpDownloadURL constructs the download URL for a Burp Suite release.
// Temporary duplication - will be refactored to use product module.
func burpDownloadURL(edition string) string {
	product := "pro"
	if edition == "community" {
		product = "community"
	}
	return "https://portswigger-cdn.net/burp/releases/download?product=" + product + "&type=Jar"
}
