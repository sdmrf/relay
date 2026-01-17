package launcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sdmrf/relay/internal/plan"
)

func TestShellLauncherPath(t *testing.T) {
	s := ShellLauncher{BinDir: "/usr/local/bin"}

	want := "/usr/local/bin/burpsuite"
	if got := s.Path(); got != want {
		t.Errorf("ShellLauncher.Path() = %v, want %v", got, want)
	}
}

func TestShellLauncherGenerate(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "launcher-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := ShellLauncher{BinDir: tmpDir}

	p := plan.LaunchPlan{
		Product: "burpsuite",
		Version: "2024.1",
		JVMArgs: []string{"-Xmx4g", "-noverify"},
		Paths: plan.Paths{
			InstallDir: "/opt/relay",
		},
	}

	if err := s.Generate(p); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check file exists
	content, err := os.ReadFile(s.Path())
	if err != nil {
		t.Fatalf("failed to read generated launcher: %v", err)
	}

	script := string(content)

	// Verify shebang
	if !strings.HasPrefix(script, "#!/bin/sh") {
		t.Error("launcher script should start with #!/bin/sh")
	}

	// Verify exec java
	if !strings.Contains(script, "exec java") {
		t.Error("launcher script should contain 'exec java'")
	}

	// Verify JVM args
	if !strings.Contains(script, "-Xmx4g") {
		t.Error("launcher script should contain JVM arg -Xmx4g")
	}
	if !strings.Contains(script, "-noverify") {
		t.Error("launcher script should contain JVM arg -noverify")
	}

	// Verify jar path
	if !strings.Contains(script, "/opt/relay/burpsuite.jar") {
		t.Error("launcher script should contain jar path")
	}

	// Verify background execution
	if !strings.Contains(script, "&") {
		t.Error("launcher script should run in background (&)")
	}

	// Check file is executable
	info, err := os.Stat(s.Path())
	if err != nil {
		t.Fatalf("failed to stat launcher: %v", err)
	}

	mode := info.Mode()
	if mode&0o111 == 0 {
		t.Error("launcher script should be executable")
	}
}

func TestShellLauncherGenerateEmptyJVMArgs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "launcher-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := ShellLauncher{BinDir: tmpDir}

	p := plan.LaunchPlan{
		Product: "burpsuite",
		Version: "2024.1",
		JVMArgs: []string{},
		Paths: plan.Paths{
			InstallDir: "/opt/relay",
		},
	}

	if err := s.Generate(p); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	content, err := os.ReadFile(s.Path())
	if err != nil {
		t.Fatalf("failed to read generated launcher: %v", err)
	}

	script := string(content)

	// Should still have valid script structure
	if !strings.HasPrefix(script, "#!/bin/sh") {
		t.Error("launcher script should start with #!/bin/sh")
	}
	if !strings.Contains(script, "exec java") {
		t.Error("launcher script should contain 'exec java'")
	}
	if !strings.Contains(script, "-jar") {
		t.Error("launcher script should contain '-jar'")
	}
}

func TestNewLauncher(t *testing.T) {
	p := plan.LaunchPlan{
		Product: "burpsuite",
		Paths: plan.Paths{
			BinDir: "/usr/local/bin",
		},
	}

	gen, err := New(p)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if gen == nil {
		t.Error("New() returned nil generator")
	}

	// Path should be set correctly
	path := gen.Path()
	if !strings.Contains(path, "burpsuite") {
		t.Errorf("generator path should contain 'burpsuite', got %v", path)
	}
	if filepath.Dir(path) != "/usr/local/bin" {
		t.Errorf("generator path should be in bin dir, got %v", filepath.Dir(path))
	}
}
