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

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove Burp Suite",
	Long:  `Uninstall Burp Suite and remove related files. Configuration is preserved.`,
	RunE:  runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
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

	removePlan, err := burp.ResolveRemove()
	if err != nil {
		return fmt.Errorf("resolve remove: %w", err)
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Removing", burp.Name())
	}

	exec := app.FSExecutor{DryRun: dryRun}
	if err := exec.Execute(cmd.Context(), removePlan); err != nil {
		return fmt.Errorf("execute remove: %w", err)
	}

	if !dryRun {
		fmt.Println("Removal complete")
	}

	return nil
}
