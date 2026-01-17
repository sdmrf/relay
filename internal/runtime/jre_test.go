package runtime

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestJREDownloadURL(t *testing.T) {
	url, ext, err := jreDownloadURL()
	if err != nil {
		// Skip on unsupported platforms
		if runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
			t.Skip("unsupported platform")
		}
		if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
			t.Skip("unsupported architecture")
		}
		t.Fatalf("jreDownloadURL() error = %v", err)
	}

	// Check URL format
	if url == "" {
		t.Error("jreDownloadURL() returned empty URL")
	}

	// Check extension
	if runtime.GOOS == "windows" {
		if ext != ".zip" {
			t.Errorf("jreDownloadURL() ext = %v, want .zip on Windows", ext)
		}
	} else {
		if ext != ".tar.gz" {
			t.Errorf("jreDownloadURL() ext = %v, want .tar.gz on Unix", ext)
		}
	}

	// Verify URL contains expected components
	expectedComponents := []string{
		"adoptium",
		"temurin21",
		"OpenJDK21U-jre",
		JREVersion,
	}
	for _, comp := range expectedComponents {
		if !contains(url, comp) {
			t.Errorf("jreDownloadURL() URL missing %q: %s", comp, url)
		}
	}
}

func TestGetBundledJREPath(t *testing.T) {
	// Create a temp directory structure
	tmpDir := t.TempDir()
	installDir := filepath.Join(tmpDir, "relay")

	// Initially should not find JRE
	if path := GetBundledJREPath(installDir); path != "" {
		t.Errorf("GetBundledJREPath() = %v, want empty when JRE not present", path)
	}

	// Create bundled JRE structure
	var javaBin string
	if runtime.GOOS == "windows" {
		javaBin = filepath.Join(installDir, "jre", "bin", "java.exe")
	} else {
		javaBin = filepath.Join(installDir, "jre", "bin", "java")
	}

	if err := os.MkdirAll(filepath.Dir(javaBin), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(javaBin, []byte("fake java"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Now should find JRE
	if path := GetBundledJREPath(installDir); path != javaBin {
		t.Errorf("GetBundledJREPath() = %v, want %v", path, javaBin)
	}
}

func TestHasBundledJRE(t *testing.T) {
	tmpDir := t.TempDir()
	installDir := filepath.Join(tmpDir, "relay")

	// Initially should not have JRE
	if HasBundledJRE(installDir) {
		t.Error("HasBundledJRE() = true, want false when JRE not present")
	}

	// Create bundled JRE structure
	var javaBin string
	if runtime.GOOS == "windows" {
		javaBin = filepath.Join(installDir, "jre", "bin", "java.exe")
	} else {
		javaBin = filepath.Join(installDir, "jre", "bin", "java")
	}

	if err := os.MkdirAll(filepath.Dir(javaBin), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(javaBin, []byte("fake java"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Now should have JRE
	if !HasBundledJRE(installDir) {
		t.Error("HasBundledJRE() = false, want true when JRE present")
	}
}

func TestResolveJavaPath(t *testing.T) {
	tmpDir := t.TempDir()
	installDir := filepath.Join(tmpDir, "relay")

	// Create bundled JRE
	var javaBin string
	if runtime.GOOS == "windows" {
		javaBin = filepath.Join(installDir, "jre", "bin", "java.exe")
	} else {
		javaBin = filepath.Join(installDir, "jre", "bin", "java")
	}

	if err := os.MkdirAll(filepath.Dir(javaBin), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(javaBin, []byte("fake java"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name       string
		strategy   string
		wantPath   string
		wantErr    bool
		skipReason string
	}{
		{
			name:     "bundled strategy with bundled JRE",
			strategy: "bundled",
			wantPath: javaBin,
			wantErr:  false,
		},
		{
			name:     "auto strategy with bundled JRE",
			strategy: "auto",
			wantPath: javaBin,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipReason != "" {
				t.Skip(tt.skipReason)
			}

			got, err := ResolveJavaPath(installDir, tt.strategy)

			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveJavaPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantPath {
				t.Errorf("ResolveJavaPath() = %v, want %v", got, tt.wantPath)
			}
		})
	}
}

func TestNeedsJRE(t *testing.T) {
	tmpDir := t.TempDir()
	installDir := filepath.Join(tmpDir, "relay")

	// Create bundled JRE
	var javaBin string
	if runtime.GOOS == "windows" {
		javaBin = filepath.Join(installDir, "jre", "bin", "java.exe")
	} else {
		javaBin = filepath.Join(installDir, "jre", "bin", "java")
	}

	tests := []struct {
		name       string
		strategy   string
		hasJRE     bool
		wantNeeds  bool
	}{
		{
			name:      "system strategy never needs bundled",
			strategy:  "system",
			hasJRE:    false,
			wantNeeds: false,
		},
		{
			name:      "bundled strategy needs JRE when not present",
			strategy:  "bundled",
			hasJRE:    false,
			wantNeeds: true,
		},
		{
			name:      "bundled strategy doesn't need JRE when present",
			strategy:  "bundled",
			hasJRE:    true,
			wantNeeds: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(tmpDir, tt.name)
			if tt.hasJRE {
				testJavaBin := javaBin
				if runtime.GOOS == "windows" {
					testJavaBin = filepath.Join(testDir, "jre", "bin", "java.exe")
				} else {
					testJavaBin = filepath.Join(testDir, "jre", "bin", "java")
				}
				if err := os.MkdirAll(filepath.Dir(testJavaBin), 0o755); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
				if err := os.WriteFile(testJavaBin, []byte("fake java"), 0o755); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			}

			got := NeedsJRE(testDir, tt.strategy)
			if got != tt.wantNeeds {
				t.Errorf("NeedsJRE() = %v, want %v", got, tt.wantNeeds)
			}
		})
	}
}

func TestBuildJREArtifact(t *testing.T) {
	tmpDir := t.TempDir()
	installDir := filepath.Join(tmpDir, "relay")

	artifact, err := BuildJREArtifact(installDir)
	if err != nil {
		// Skip on unsupported platforms
		if runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
			t.Skip("unsupported platform")
		}
		if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
			t.Skip("unsupported architecture")
		}
		t.Fatalf("BuildJREArtifact() error = %v", err)
	}

	if artifact.Name == "" {
		t.Error("BuildJREArtifact() Name is empty")
	}
	if artifact.URL == "" {
		t.Error("BuildJREArtifact() URL is empty")
	}
	if artifact.Target == "" {
		t.Error("BuildJREArtifact() Target is empty")
	}
	if artifact.ExtractTo != installDir {
		t.Errorf("BuildJREArtifact() ExtractTo = %v, want %v", artifact.ExtractTo, installDir)
	}
}

func TestJREExtractedDir(t *testing.T) {
	dir := JREExtractedDir()
	if dir == "" {
		t.Error("JREExtractedDir() returned empty string")
	}

	// Should contain version
	if !contains(dir, JREVersion) {
		t.Errorf("JREExtractedDir() = %v, should contain version %v", dir, JREVersion)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
