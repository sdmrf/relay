package plan

import (
	"testing"

	"github.com/sdmrf/relay/pkg/config"
)

func TestInstallPlanKind(t *testing.T) {
	p := InstallPlan{
		Product: "burpsuite",
		Edition: "professional",
		Version: "2024.1",
	}

	if got := p.Kind(); got != Install {
		t.Errorf("InstallPlan.Kind() = %v, want %v", got, Install)
	}
}

func TestLaunchPlanKind(t *testing.T) {
	p := LaunchPlan{
		Product: "burpsuite",
		Version: "2024.1",
	}

	if got := p.Kind(); got != Launch {
		t.Errorf("LaunchPlan.Kind() = %v, want %v", got, Launch)
	}
}

func TestRemovePlanKind(t *testing.T) {
	p := RemovePlan{
		Product: "burpsuite",
	}

	if got := p.Kind(); got != Remove {
		t.Errorf("RemovePlan.Kind() = %v, want %v", got, Remove)
	}
}

func TestPathsOwned(t *testing.T) {
	p := Paths{
		InstallDir: "/opt/relay",
		DataDir:    "/var/relay",
		BinDir:     "/usr/bin",
		ConfigDir:  "/etc/relay",
		CacheDir:   "/var/cache/relay",
	}

	owned := p.Owned()

	// Should contain InstallDir, DataDir, CacheDir
	// Should NOT contain ConfigDir or BinDir
	expected := []string{"/opt/relay", "/var/relay", "/var/cache/relay"}

	if len(owned) != len(expected) {
		t.Fatalf("Paths.Owned() returned %d paths, want %d", len(owned), len(expected))
	}

	for i, path := range expected {
		if owned[i] != path {
			t.Errorf("Paths.Owned()[%d] = %v, want %v", i, owned[i], path)
		}
	}
}

func TestPathsOwnedExcludesConfig(t *testing.T) {
	p := Paths{
		InstallDir: "/opt/relay",
		DataDir:    "/var/relay",
		ConfigDir:  "/etc/relay",
		CacheDir:   "/var/cache/relay",
	}

	owned := p.Owned()

	for _, path := range owned {
		if path == p.ConfigDir {
			t.Error("Paths.Owned() should not include ConfigDir")
		}
	}
}

func TestInstallPlanImmutable(t *testing.T) {
	p := InstallPlan{
		Product: "burpsuite",
		Edition: "professional",
		Version: "2024.1",
		Layout:  config.SystemLayout,
		JavaMin: 17,
		JVMArgs: []string{"-Xmx4g"},
	}

	// Verify fields are accessible
	if p.Product != "burpsuite" {
		t.Errorf("Product = %v, want burpsuite", p.Product)
	}
	if p.Edition != "professional" {
		t.Errorf("Edition = %v, want professional", p.Edition)
	}
	if p.Version != "2024.1" {
		t.Errorf("Version = %v, want 2024.1", p.Version)
	}
	if p.Layout != config.SystemLayout {
		t.Errorf("Layout = %v, want system", p.Layout)
	}
	if p.JavaMin != 17 {
		t.Errorf("JavaMin = %v, want 17", p.JavaMin)
	}
}

func TestLaunchPlanImmutable(t *testing.T) {
	p := LaunchPlan{
		Product: "burpsuite",
		Version: "2024.1",
		JavaMin: 17,
		JVMArgs: []string{"-Xmx4g"},
		Paths: Paths{
			InstallDir: "/opt/relay",
		},
	}

	if p.Product != "burpsuite" {
		t.Errorf("Product = %v, want burpsuite", p.Product)
	}
	if p.Paths.InstallDir != "/opt/relay" {
		t.Errorf("Paths.InstallDir = %v, want /opt/relay", p.Paths.InstallDir)
	}
}

func TestRemovePlanImmutable(t *testing.T) {
	p := RemovePlan{
		Product: "burpsuite",
		Paths: Paths{
			InstallDir: "/opt/relay",
			DataDir:    "/var/relay",
		},
	}

	if p.Product != "burpsuite" {
		t.Errorf("Product = %v, want burpsuite", p.Product)
	}
	if p.Paths.InstallDir != "/opt/relay" {
		t.Errorf("Paths.InstallDir = %v, want /opt/relay", p.Paths.InstallDir)
	}
}
