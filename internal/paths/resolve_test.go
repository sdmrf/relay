package paths

import (
	"path/filepath"
	"testing"
)

func TestResolvePortable(t *testing.T) {
	root := "/home/user/burpsuite"

	p, err := Resolve(Options{
		Layout:      PortableLayout,
		InstallHint: root,
	})

	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if p.InstallDir != root {
		t.Errorf("InstallDir = %v, want %v", p.InstallDir, root)
	}
	if p.DataDir != filepath.Join(root, "data") {
		t.Errorf("DataDir = %v, want %v", p.DataDir, filepath.Join(root, "data"))
	}
	if p.BinDir != filepath.Join(root, "bin") {
		t.Errorf("BinDir = %v, want %v", p.BinDir, filepath.Join(root, "bin"))
	}
	if p.ConfigDir != filepath.Join(root, "config") {
		t.Errorf("ConfigDir = %v, want %v", p.ConfigDir, filepath.Join(root, "config"))
	}
	if p.CacheDir != filepath.Join(root, "cache") {
		t.Errorf("CacheDir = %v, want %v", p.CacheDir, filepath.Join(root, "cache"))
	}
}

func TestResolvePortableRequiresInstallHint(t *testing.T) {
	_, err := Resolve(Options{
		Layout:      PortableLayout,
		InstallHint: "",
	})

	if err == nil {
		t.Error("Resolve() expected error for portable layout without install hint")
	}
}

func TestResolvePortableRejectsAuto(t *testing.T) {
	_, err := Resolve(Options{
		Layout:      PortableLayout,
		InstallHint: "auto",
	})

	if err == nil {
		t.Error("Resolve() expected error for portable layout with 'auto' install hint")
	}
}

func TestResolveInvalidLayout(t *testing.T) {
	_, err := Resolve(Options{
		Layout: Layout("invalid"),
	})

	if err == nil {
		t.Error("Resolve() expected error for invalid layout")
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

	// Owned should include: InstallDir, DataDir, CacheDir
	// Owned should NOT include: ConfigDir (user config) or BinDir (shared system dir)
	if len(owned) != 3 {
		t.Fatalf("Owned() returned %d paths, want 3", len(owned))
	}

	contains := func(slice []string, item string) bool {
		for _, s := range slice {
			if s == item {
				return true
			}
		}
		return false
	}

	if !contains(owned, p.InstallDir) {
		t.Error("Owned() should include InstallDir")
	}
	if !contains(owned, p.DataDir) {
		t.Error("Owned() should include DataDir")
	}
	if !contains(owned, p.CacheDir) {
		t.Error("Owned() should include CacheDir")
	}
	if contains(owned, p.ConfigDir) {
		t.Error("Owned() should not include ConfigDir")
	}
	if contains(owned, p.BinDir) {
		t.Error("Owned() should not include BinDir (shared system directory)")
	}
}

func TestCurrentOS(t *testing.T) {
	os := CurrentOS()

	// Just verify it returns a valid value
	switch os {
	case Linux, Darwin, Windows:
		// valid
	default:
		t.Errorf("CurrentOS() = %v, want linux, darwin, or windows", os)
	}
}
