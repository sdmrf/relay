package main

import (
	"fmt"
	"os"

	"github.com/sdmrf/relay/internal/app"
	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/internal/product/burpsuite"
	"github.com/sdmrf/relay/pkg/config"
	"github.com/spf13/cobra"
)

var (
	updateForce bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Burp Suite to the latest version",
	Long:  `Check for updates and download the latest version of Burp Suite.`,
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVarP(&updateForce, "force", "f", false, "force update even if already at latest version")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	p, err := paths.Resolve(paths.Options{
		Layout:      paths.Layout(cfg.Layout.Mode),
		InstallHint: cfg.Paths.Install,
		DataHint:    cfg.Paths.Data,
		BinHint:     cfg.Paths.Bin,
	})
	if err != nil {
		return fmt.Errorf("resolve paths: %w", err)
	}

	burp, err := burpsuite.New(cfg, p)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	updatePlan, err := burp.ResolveUpdate()
	if err != nil {
		// No installed version - suggest running install
		fmt.Fprintln(os.Stderr, "No installation found. Run 'relay install' first.")
		return fmt.Errorf("resolve update: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Current version: %s\n", updatePlan.CurrentVersion)
		fmt.Fprintf(os.Stderr, "Target version:  %s\n", updatePlan.TargetVersion)
	}

	// For now, always proceed with update (downloads latest)
	// Future: check actual version from server
	if !updateForce && updatePlan.TargetVersion != "latest" {
		cmp := burpsuite.CompareVersions(updatePlan.CurrentVersion, updatePlan.TargetVersion)
		if cmp >= 0 {
			fmt.Println("Already at latest version:", updatePlan.CurrentVersion)
			return nil
		}
	}

	fmt.Printf("Updating %s from %s to %s...\n",
		burp.Name(), updatePlan.CurrentVersion, updatePlan.TargetVersion)

	// Use install plan for the actual download
	installPlan, err := burp.ResolveInstall()
	if err != nil {
		return fmt.Errorf("resolve install: %w", err)
	}

	exec := app.FSExecutor{DryRun: dryRun}
	if err := exec.Execute(cmd.Context(), installPlan); err != nil {
		return fmt.Errorf("execute update: %w", err)
	}

	// Write version marker
	if !dryRun {
		version := updatePlan.TargetVersion
		if version == "latest" {
			version = "unknown" // Will be determined when we fetch actual version
		}
		if err := burpsuite.WriteVersionMarker(p.InstallDir, version); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write version marker: %v\n", err)
		}
		fmt.Println("Update complete")
	}

	return nil
}
