package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// JavaInfo contains information about a Java installation.
type JavaInfo struct {
	Version int
	Path    string
	Output  string
}

// ValidateJava checks that Java is installed and meets minimum version.
// Returns the version output on success.
func ValidateJava(minVersion int) (string, error) {
	info, err := GetJavaInfo()
	if err != nil {
		return "", err
	}

	if info.Version < minVersion {
		return "", fmt.Errorf("java %d+ required, found %d", minVersion, info.Version)
	}

	return info.Output, nil
}

// GetJavaInfo retrieves information about the installed Java.
func GetJavaInfo() (JavaInfo, error) {
	cmd := exec.Command("java", "-version")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return JavaInfo{}, fmt.Errorf("java not found in PATH")
	}

	output := stderr.String()

	version, err := ParseJavaVersion(output)
	if err != nil {
		return JavaInfo{}, err
	}

	path := getJavaPath()

	return JavaInfo{
		Version: version,
		Path:    path,
		Output:  output,
	}, nil
}

// ParseJavaVersion extracts the major version number from java -version output.
func ParseJavaVersion(output string) (int, error) {
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

// getJavaPath attempts to find the java executable path.
func getJavaPath() string {
	cmd := exec.Command("which", "java")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
