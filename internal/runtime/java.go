package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ValidateJava checks that Java is installed and meets minimum version.
// Returns the version output on success.
func ValidateJava(minVersion int) (string, error) {
	cmd := exec.Command("java", "-version")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("java not found in PATH")
	}

	out := stderr.String()

	version, err := parseJavaVersion(out)
	if err != nil {
		return "", err
	}

	if version < minVersion {
		return "", fmt.Errorf("java %d+ required, found %d", minVersion, version)
	}

	return out, nil
}

// parseJavaVersion extracts the major version number from java -version output.
func parseJavaVersion(output string) (int, error) {
	// Handles formats:
	// - java version "17.0.8"
	// - openjdk version "21.0.1"
	// - java version "1.8.0_392" (Java 8)
	start := strings.Index(output, `"`)
	if start == -1 {
		return 0, fmt.Errorf("unable to parse java version")
	}

	end := strings.Index(output[start+1:], `"`)
	if end == -1 {
		return 0, fmt.Errorf("unable to parse java version")
	}

	raw := output[start+1 : start+1+end]
	parts := strings.Split(raw, ".")

	// Handle 1.8 style versions (Java 8)
	if parts[0] == "1" && len(parts) > 1 {
		return strconv.Atoi(parts[1])
	}

	return strconv.Atoi(parts[0])
}
