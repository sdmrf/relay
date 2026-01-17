package main

import (
	"fmt"

	"github.com/sdmrf/relay/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, commit, and build date of relay.`,
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("relay %s\n", version.Version)
	fmt.Printf("  commit:  %s\n", version.Commit)
	fmt.Printf("  built:   %s\n", version.BuildDate)
}
