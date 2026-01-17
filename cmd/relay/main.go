package main

import (
	"fmt"
	"os"
	"github.com/sdmrf/relay/pkg/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}

	fmt.Println("relay config loaded")
	_ = cfg
}
