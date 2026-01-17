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
	installEdition string
	installVersion string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Burp Suite",
	Long:  `Download and install Burp Suite with the specified edition and version.`,
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().StringVar(&installEdition, "edition", "", "edition to install (professional, community)")
	installCmd.Flags().StringVar(&installVersion, "version", "", "version to install (default: latest)")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Override config with flags if provided
	if installEdition != "" {
		cfg.Product.Edition = installEdition
	}
	if installVersion != "" {
		cfg.Product.Version = installVersion
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

	installPlan, err := burp.ResolveInstall()
	if err != nil {
		return fmt.Errorf("resolve install: %w", err)
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Installing", burp.Name(), cfg.Product.Edition, cfg.Product.Version)
	}

	exec := app.FSExecutor{DryRun: dryRun}
	if err := exec.Execute(cmd.Context(), installPlan); err != nil {
		return fmt.Errorf("execute install: %w", err)
	}

	if !dryRun {
		fmt.Println("Installation complete")
	}

	return nil
}
