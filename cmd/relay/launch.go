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

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch Burp Suite",
	Long:  `Start Burp Suite with the configured settings.`,
	RunE:  runLaunch,
}

func init() {
	rootCmd.AddCommand(launchCmd)
}

func runLaunch(cmd *cobra.Command, args []string) error {
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

	launchPlan, err := burp.ResolveLaunch()
	if err != nil {
		return fmt.Errorf("resolve launch: %w", err)
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Launching", burp.Name(), cfg.Product.Version)
	}

	exec := app.FSExecutor{DryRun: dryRun}
	if err := exec.Execute(cmd.Context(), launchPlan); err != nil {
		return fmt.Errorf("execute launch: %w", err)
	}

	return nil
}
