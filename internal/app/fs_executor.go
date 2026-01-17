package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

	dl := downloader.HTTPDownloader{
		Timeout: 5 * time.Minute,
		Retries: 3,
	}

	// Download and extract JRE if needed
	if p.JREArtifact != nil {
		if err := e.downloadAndExtractJRE(ctx, dl, p.JREArtifact, p.Paths.InstallDir); err != nil {
			return fmt.Errorf("install JRE: %w", err)
		}
	}

	// Download product artifact
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

	fmt.Println("Downloading", artifact.Name)
	return dl.FetchWithProgress(ctx, artifact)
}

// downloadAndExtractJRE downloads and extracts the JRE archive.
func (e FSExecutor) downloadAndExtractJRE(ctx context.Context, dl downloader.HTTPDownloader, jre *plan.JREArtifact, installDir string) error {
	if e.DryRun {
		fmt.Println("[dry-run] download JRE:", jre.Name)
		fmt.Println("[dry-run]   url:", jre.URL)
		fmt.Println("[dry-run]   target:", jre.Target)
		fmt.Println("[dry-run] extract JRE to:", jre.ExtractTo)
		return nil
	}

	// Download JRE archive
	artifact := downloader.Artifact{
		Name:   jre.Name,
		URL:    jre.URL,
		Target: jre.Target,
	}

	if err := dl.FetchWithProgress(ctx, artifact); err != nil {
		return fmt.Errorf("download JRE: %w", err)
	}

	fmt.Println("Extracting JRE...")

	// Extract to temp directory first (atomic extraction)
	tmpDir := filepath.Join(installDir, ".jre-extract-tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("cleanup temp dir: %w", err)
	}
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	if err := downloader.Extract(jre.Target, tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("extract JRE: %w", err)
	}

	// Find the extracted directory (Adoptium extracts to jdk-21.0.5+11-jre)
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("read extracted dir: %w", err)
	}

	if len(entries) == 0 {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("JRE archive was empty")
	}

	// Move to final location
	jreDir := filepath.Join(installDir, "jre")
	if err := os.RemoveAll(jreDir); err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("remove existing jre: %w", err)
	}

	// The first (and usually only) directory in the archive is the JRE root
	extractedDir := filepath.Join(tmpDir, entries[0].Name())
	if err := os.Rename(extractedDir, jreDir); err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("move JRE to final location: %w", err)
	}

	// Cleanup
	os.RemoveAll(tmpDir)
	os.Remove(jre.Target) // Remove downloaded archive

	fmt.Println("JRE installed to:", jreDir)
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

// execLaunch validates Java, generates the launcher, and runs it.
func (e FSExecutor) execLaunch(ctx context.Context, p plan.LaunchPlan) error {
	// Resolve Java path based on strategy
	javaPath, err := runtime.ResolveJavaPath(p.Paths.InstallDir, string(p.JavaStrategy))
	if err != nil {
		return fmt.Errorf("resolve java: %w", err)
	}

	gen, err := launcher.New(p, javaPath)
	if err != nil {
		return fmt.Errorf("create launcher: %w", err)
	}

	if e.DryRun {
		fmt.Printf("[dry-run] using java: %s\n", javaPath)
		fmt.Println("[dry-run] generate launcher:", gen.Path())
		fmt.Println("[dry-run] run launcher:", gen.Path())
		return nil
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

	fmt.Printf("Updating %s -> %s\n", p.CurrentVersion, p.TargetVersion)

	dl := downloader.HTTPDownloader{
		Timeout: 5 * time.Minute,
		Retries: 3,
	}

	return dl.FetchWithProgress(ctx, artifact)
}
