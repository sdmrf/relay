package main

import (
	"fmt"
	"os"

	"github.com/sdmrf/relay/internal/diagnostics"
	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/pkg/config"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system readiness",
	Long:  `Run diagnostic checks to verify the system is ready to run Burp Suite.`,
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("relay doctor")
	fmt.Println()

	report := diagnostics.NewReport()

	// Load config first to get java version requirement
	cfg, err := config.Load(cfgFile)
	if err != nil {
		cfg = config.Default()
	}

	// Check Java
	report.Add(diagnostics.CheckJava(cfg.Runtime.Java.MinVersion))

	// Check Config
	report.Add(diagnostics.CheckConfig(cfgFile))

	// Resolve and check paths
	p, err := paths.Resolve(paths.Options{
		Layout:      paths.Layout(cfg.Layout.Mode),
		InstallHint: cfg.Paths.Install,
		DataHint:    cfg.Paths.Data,
		BinHint:     cfg.Paths.Bin,
	})
	if err == nil {
		report.AddAll(diagnostics.CheckPaths(p))
		report.Add(diagnostics.CheckProduct(p.InstallDir))
	} else {
		report.Add(diagnostics.Check{
			Name:    "Paths",
			Status:  diagnostics.StatusFail,
			Message: fmt.Sprintf("Failed to resolve: %v", err),
		})
	}

	// Check network
	report.Add(diagnostics.CheckNetwork())

	// Print report
	if verbose {
		report.PrintVerbose(os.Stdout)
	} else {
		report.Print(os.Stdout)
	}

	fmt.Println()

	if report.HasFailures() {
		fmt.Println("Some checks failed. Please fix the issues above.")
		return fmt.Errorf("diagnostics failed")
	}

	if report.HasWarnings() {
		fmt.Println("All critical checks passed with warnings.")
	} else {
		fmt.Println("All checks passed.")
	}

	return nil
}
