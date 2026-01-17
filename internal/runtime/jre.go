package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// JRE version to bundle - Eclipse Temurin 21 LTS
const (
	JREVersion      = "21.0.5"
	JREBuildNumber  = "11"
	JREMajorVersion = 21
)

// JREArtifact represents a downloadable JRE.
type JREArtifact struct {
	Name      string // Display name
	URL       string // Download URL
	Target    string // Download target path (.tar.gz or .zip)
	ExtractTo string // Extraction destination directory
}

// GetBundledJREPath returns the path to the bundled java binary if it exists.
// Returns empty string if bundled JRE is not found.
func GetBundledJREPath(installDir string) string {
	var candidates []string

	switch runtime.GOOS {
	case "darwin":
		// macOS uses app bundle structure: jre/Contents/Home/bin/java
		candidates = []string{
			filepath.Join(installDir, "jre", "Contents", "Home", "bin", "java"),
			filepath.Join(installDir, "jre", "bin", "java"), // fallback for non-bundle JREs
		}
	case "windows":
		candidates = []string{
			filepath.Join(installDir, "jre", "bin", "java.exe"),
		}
	default: // linux
		candidates = []string{
			filepath.Join(installDir, "jre", "bin", "java"),
		}
	}

	for _, javaBin := range candidates {
		if _, err := os.Stat(javaBin); err == nil {
			return javaBin
		}
	}
	return ""
}

// HasBundledJRE checks if a bundled JRE exists in the install directory.
func HasBundledJRE(installDir string) bool {
	return GetBundledJREPath(installDir) != ""
}

// HasSystemJava checks if Java is available on the system PATH.
func HasSystemJava() bool {
	_, err := GetJavaInfo()
	return err == nil
}

// ResolveJavaPath determines which Java to use based on strategy.
// Returns the path to the java binary.
func ResolveJavaPath(installDir string, strategy string) (string, error) {
	switch strategy {
	case "system":
		// Only use system Java
		info, err := GetJavaInfo()
		if err != nil {
			return "", fmt.Errorf("system java required but not found: %w", err)
		}
		return info.Path, nil

	case "bundled":
		// Only use bundled JRE
		path := GetBundledJREPath(installDir)
		if path == "" {
			return "", fmt.Errorf("bundled JRE required but not found in %s", installDir)
		}
		return path, nil

	default: // "auto"
		// Prefer bundled, fallback to system
		if path := GetBundledJREPath(installDir); path != "" {
			return path, nil
		}
		info, err := GetJavaInfo()
		if err != nil {
			return "", fmt.Errorf("no java found (bundled or system)")
		}
		return info.Path, nil
	}
}

// NeedsJRE determines if JRE needs to be downloaded.
// Returns true if no Java is available (bundled or system).
func NeedsJRE(installDir string, strategy string) bool {
	switch strategy {
	case "system":
		// Never download JRE if system-only strategy
		return false
	case "bundled":
		// Need JRE if bundled doesn't exist
		return !HasBundledJRE(installDir)
	default: // "auto"
		// Need JRE if neither bundled nor system available
		return !HasBundledJRE(installDir) && !HasSystemJava()
	}
}

// BuildJREArtifact creates a JREArtifact for the current platform.
func BuildJREArtifact(installDir string) (JREArtifact, error) {
	url, ext, err := jreDownloadURL()
	if err != nil {
		return JREArtifact{}, err
	}

	archiveName := fmt.Sprintf("jre-%s%s", JREVersion, ext)

	return JREArtifact{
		Name:      fmt.Sprintf("Eclipse Temurin JRE %s", JREVersion),
		URL:       url,
		Target:    filepath.Join(installDir, archiveName),
		ExtractTo: installDir,
	}, nil
}

// jreDownloadURL returns the Adoptium download URL for the current platform.
func jreDownloadURL() (url string, ext string, err error) {
	var osName, archName string

	switch runtime.GOOS {
	case "linux":
		osName = "linux"
		ext = ".tar.gz"
	case "darwin":
		osName = "mac"
		ext = ".tar.gz"
	case "windows":
		osName = "windows"
		ext = ".zip"
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "amd64":
		archName = "x64"
	case "arm64":
		archName = "aarch64"
	default:
		return "", "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	// Eclipse Adoptium Temurin URL pattern
	// Example: https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.5+11/OpenJDK21U-jre_x64_mac_hotspot_21.0.5_11.tar.gz
	baseURL := "https://github.com/adoptium/temurin21-binaries/releases/download"
	tag := fmt.Sprintf("jdk-%s+%s", JREVersion, JREBuildNumber)
	fileName := fmt.Sprintf("OpenJDK21U-jre_%s_%s_hotspot_%s_%s%s",
		archName, osName, JREVersion, JREBuildNumber, ext)

	url = fmt.Sprintf("%s/%s/%s", baseURL, tag, fileName)
	return url, ext, nil
}

// JREExtractedDir returns the expected directory name after extraction.
// Adoptium archives extract to a directory like "jdk-21.0.5+11-jre"
func JREExtractedDir() string {
	return fmt.Sprintf("jdk-%s+%s-jre", JREVersion, JREBuildNumber)
}
