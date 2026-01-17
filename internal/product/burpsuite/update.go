package burpsuite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdmrf/relay/internal/plan"
)

const versionMarkerFile = ".relay-version"

// ResolveUpdate creates an immutable UpdatePlan for Burp Suite.
// Returns an error if update is not needed or cannot be determined.
func (b *BurpSuite) ResolveUpdate() (plan.UpdatePlan, error) {
	currentVersion, err := b.readInstalledVersion()
	if err != nil {
		// No version marker - treat as fresh install needed
		return plan.UpdatePlan{}, fmt.Errorf("no installed version found: %w", err)
	}

	// For now, we always suggest updating to "latest"
	// Future: could fetch actual latest version from PortSwigger
	targetVersion := b.cfg.Product.Version
	if targetVersion == "" || targetVersion == "latest" {
		targetVersion = "latest"
	}

	return plan.UpdatePlan{
		Product:        b.Name(),
		Edition:        b.cfg.Product.Edition,
		CurrentVersion: currentVersion,
		TargetVersion:  targetVersion,
		Paths:          plan.FromResolved(b.paths),
		Artifact: plan.Artifact{
			Name:   JarName,
			URL:    burpDownloadURL(b.cfg.Product.Edition),
			Target: filepath.Join(b.paths.InstallDir, JarName),
		},
	}, nil
}

// readInstalledVersion reads the version from the marker file.
func (b *BurpSuite) readInstalledVersion() (string, error) {
	markerPath := filepath.Join(b.paths.InstallDir, versionMarkerFile)

	data, err := os.ReadFile(markerPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// WriteVersionMarker writes the version to the marker file.
func WriteVersionMarker(installDir, version string) error {
	markerPath := filepath.Join(installDir, versionMarkerFile)
	return os.WriteFile(markerPath, []byte(version+"\n"), 0o644)
}

// CompareVersions compares two version strings.
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
//
// Handles version formats like "2024.1.1", "2023.12.1.4"
func CompareVersions(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var numA, numB int
		if i < len(partsA) {
			fmt.Sscanf(partsA[i], "%d", &numA)
		}
		if i < len(partsB) {
			fmt.Sscanf(partsB[i], "%d", &numB)
		}

		if numA < numB {
			return -1
		}
		if numA > numB {
			return 1
		}
	}

	return 0
}
