package main

import (
	"fmt"
	"os"

	"github.com/sdmrf/relay/internal/app"
	"github.com/sdmrf/relay/internal/paths"
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

	plan := app.ResolveInstall(cfg, p)

	fmt.Println("install plan:")
	fmt.Printf("  product: %s\n", plan.Product)
	fmt.Printf("  edition: %s\n", plan.Edition)
	fmt.Printf("  version: %s\n", plan.Version)
	fmt.Printf("  layout:  %s\n", plan.Layout)
	fmt.Printf("  java_min: %d\n", plan.JavaMin)
	fmt.Printf("  jvm_args: %v\n", plan.JVMArgs)
	fmt.Println("  paths:")
	fmt.Printf("    install: %s\n", plan.Paths.InstallDir)
	fmt.Printf("    data:    %s\n", plan.Paths.DataDir)
	fmt.Printf("    bin:     %s\n", plan.Paths.BinDir)
	fmt.Printf("    config:  %s\n", plan.Paths.ConfigDir)
	fmt.Printf("    cache:   %s\n", plan.Paths.CacheDir)
}
