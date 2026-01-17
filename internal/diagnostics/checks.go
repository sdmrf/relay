package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sdmrf/relay/internal/paths"
	"github.com/sdmrf/relay/internal/runtime"
	"github.com/sdmrf/relay/pkg/config"
)

// Check represents a single diagnostic check.
type Check struct {
	Name    string
	Status  Status
	Message string
	Details string
}

// Status represents the result status of a check.
type Status int

const (
	StatusOK Status = iota
	StatusWarn
	StatusFail
)

// CheckJava verifies Java installation and version.
func CheckJava(minVersion int) Check {
	check := Check{Name: "Java"}

	info, err := runtime.GetJavaInfo()
	if err != nil {
		check.Status = StatusFail
		check.Message = "Java not found in PATH"
		return check
	}

	if info.Version < minVersion {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Java %d+ required, found %d", minVersion, info.Version)
		check.Details = info.Path
		return check
	}

	check.Status = StatusOK
	check.Message = fmt.Sprintf("Java %d found", info.Version)
	check.Details = info.Path

	return check
}

// CheckConfig verifies the configuration file.
func CheckConfig(cfgPath string) Check {
	check := Check{Name: "Config"}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Failed to load: %v", err)
		return check
	}

	if err := cfg.Validate(); err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Validation failed: %v", err)
		return check
	}

	check.Status = StatusOK
	check.Message = fmt.Sprintf("Loaded from %s", cfgPath)

	return check
}

// CheckPaths verifies required directories exist and are writable.
func CheckPaths(p paths.Paths) []Check {
	checks := []Check{}

	// Check install directory
	checks = append(checks, checkDirectory("Install directory", p.InstallDir))

	// Check data directory
	checks = append(checks, checkDirectory("Data directory", p.DataDir))

	// Check bin directory
	checks = append(checks, checkDirectory("Bin directory", p.BinDir))

	// Check cache directory
	checks = append(checks, checkDirectory("Cache directory", p.CacheDir))

	return checks
}

// CheckProduct verifies the product installation.
func CheckProduct(installDir string) Check {
	check := Check{Name: "Product"}

	jarPath := filepath.Join(installDir, "burpsuite.jar")
	info, err := os.Stat(jarPath)
	if os.IsNotExist(err) {
		check.Status = StatusWarn
		check.Message = "Burp Suite JAR not found"
		check.Details = jarPath
		return check
	}
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Error checking JAR: %v", err)
		return check
	}

	check.Status = StatusOK
	check.Message = fmt.Sprintf("Burp Suite JAR present (%d MB)", info.Size()/1024/1024)
	check.Details = jarPath

	return check
}

// CheckNetwork verifies network connectivity to PortSwigger.
func CheckNetwork() Check {
	check := Check{Name: "Network"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://portswigger.net", nil)
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Failed to create request: %v", err)
		return check
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.Status = StatusWarn
		check.Message = "portswigger.net unreachable"
		check.Details = err.Error()
		return check
	}
	defer resp.Body.Close()

	check.Status = StatusOK
	check.Message = "portswigger.net reachable"

	return check
}

func checkDirectory(name, path string) Check {
	check := Check{Name: name}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		check.Status = StatusWarn
		check.Message = "Directory does not exist"
		check.Details = path
		return check
	}
	if err != nil {
		check.Status = StatusFail
		check.Message = fmt.Sprintf("Error: %v", err)
		check.Details = path
		return check
	}
	if !info.IsDir() {
		check.Status = StatusFail
		check.Message = "Path is not a directory"
		check.Details = path
		return check
	}

	// Check if writable by attempting to create a temp file
	testFile := filepath.Join(path, ".relay-test")
	f, err := os.Create(testFile)
	if err != nil {
		check.Status = StatusWarn
		check.Message = "Directory not writable"
		check.Details = path
		return check
	}
	f.Close()
	os.Remove(testFile)

	check.Status = StatusOK
	check.Message = "Directory exists and writable"
	check.Details = path

	return check
}
