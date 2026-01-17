package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sdmrf/relay/internal/app"
	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/internal/plan"
	"github.com/sdmrf/relay/internal/product/burpsuite"
	"github.com/sdmrf/relay/internal/runtime"
	"github.com/sdmrf/relay/pkg/config"
	"github.com/spf13/cobra"
)

var (
	installEdition string
	installVersion string
	installYes     bool // Skip confirmation prompts
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
	installCmd.Flags().BoolVarP(&installYes, "yes", "y", false, "skip confirmation prompts")
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

	// Check if JRE needs to be downloaded
	needsJRE := runtime.NeedsJRE(p.InstallDir, string(cfg.Runtime.Java.Strategy))
	if needsJRE {
		// Prompt user for confirmation
		if !installYes && !dryRun {
			fmt.Println("Java runtime not found on your system.")
			fmt.Print("Download bundled JRE (~50MB)? [y/N]: ")

			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read response: %w", err)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("\nJava 17+ is required to run Burp Suite.")
				fmt.Println("Install Java manually, then run 'relay install' again.")
				fmt.Println("\nInstallation options:")
				fmt.Println("  - macOS:   brew install openjdk@21")
				fmt.Println("  - Ubuntu:  apt install openjdk-21-jre")
				fmt.Println("  - Windows: Download from https://adoptium.net")
				return nil
			}
		}

		// Build JRE artifact
		jreArtifact, err := runtime.BuildJREArtifact(p.InstallDir)
		if err != nil {
			return fmt.Errorf("build JRE artifact: %w", err)
		}

		installPlan.JREArtifact = &plan.JREArtifact{
			Name:      jreArtifact.Name,
			URL:       jreArtifact.URL,
			Target:    jreArtifact.Target,
			ExtractTo: jreArtifact.ExtractTo,
		}
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Installing", burp.Name(), cfg.Product.Edition, cfg.Product.Version)
		if installPlan.JREArtifact != nil {
			fmt.Fprintln(os.Stderr, "Including bundled JRE:", installPlan.JREArtifact.Name)
		}
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
