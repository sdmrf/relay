package main

import (
	"fmt"
	"os"

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

	installPlan, err := burp.ResolveInstall()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("burp install plan:")
	fmt.Printf("  product: %s\n", installPlan.Product)
	fmt.Printf("  edition: %s\n", installPlan.Edition)
	fmt.Printf("  version: %s\n", installPlan.Version)
	fmt.Printf("  layout:  %s\n", installPlan.Layout)
	fmt.Printf("  java_min: %d\n", installPlan.JavaMin)
	fmt.Printf("  jvm_args: %v\n", installPlan.JVMArgs)
	fmt.Println("  paths:")
	fmt.Printf("    install: %s\n", installPlan.Paths.InstallDir)
	fmt.Printf("    data:    %s\n", installPlan.Paths.DataDir)
	fmt.Printf("    bin:     %s\n", installPlan.Paths.BinDir)
	fmt.Printf("    config:  %s\n", installPlan.Paths.ConfigDir)
	fmt.Printf("    cache:   %s\n", installPlan.Paths.CacheDir)
}
