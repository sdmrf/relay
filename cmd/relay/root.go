package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	dryRun  bool
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "relay",
	Short: "A modern CLI for managing Burp Suite installations",
	Long:  `relay is a command-line tool for installing, launching, and managing Burp Suite.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview actions without executing")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
