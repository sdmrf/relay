package main

import (
	"fmt"
	"os"

	"github.com/sdmrf/relay/internal/app"
	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/internal/product/burpsuite"
	"github.com/sdmrf/relay/pkg/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	p, err := paths.Resolve(paths.Options{
		Layout:      paths.Layout(cfg.Layout.Mode),
		InstallHint: cfg.Paths.Install,
		DataHint:    cfg.Paths.Data,
		BinHint:     cfg.Paths.Bin,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	burp, err := burpsuite.New(cfg, p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exec := app.FSExecutor{DryRun: true}

	// Test launch plan (includes Java validation)
	launchPlan, err := burp.ResolveLaunch()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := exec.Execute(launchPlan); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("execution completed")
}
