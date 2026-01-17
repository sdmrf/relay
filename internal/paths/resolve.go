package paths

import (
	"fmt"
	"path/filepath"
)

type Options struct {
	Layout      Layout
	InstallHint string
	DataHint    string
	BinHint     string
}

func Resolve(opts Options) (Paths, error) {
	os := CurrentOS()

	configDir, cacheDir, err := baseDirs(os)
	if err != nil {
		return Paths{}, err
	}

	switch opts.Layout {
	case SystemLayout:
		return resolveSystem(os, configDir, cacheDir, opts)
	case PortableLayout:
		return resolvePortable(configDir, cacheDir, opts)
	default:
		return Paths{}, fmt.Errorf("invalid layout: %s", opts.Layout)
	}
}

func resolveSystem(os OS, configDir, cacheDir string, opts Options) (Paths, error) {
	home, err := osUserHome()
	if err != nil {
		return Paths{}, err
	}

	var install, data string

	switch os {
	case Linux:
		// User-writable locations
		install = filepath.Join(home, ".local", "share", "relay")
		data = filepath.Join(home, ".local", "share", "relay", "data")
	case Darwin:
		// User-writable locations in ~/Library
		install = filepath.Join(home, "Library", "Application Support", "relay")
		data = filepath.Join(home, "Library", "Application Support", "relay", "data")
	case Windows:
		// User's local app data
		localAppData := filepath.Join(home, "AppData", "Local")
		install = filepath.Join(localAppData, "relay")
		data = filepath.Join(localAppData, "relay", "data")
	default:
		return Paths{}, fmt.Errorf("unsupported OS: %s", os)
	}

	// Bin directory inside install dir - no sudo needed for launcher
	bin := filepath.Join(install, "bin")

	return Paths{
		InstallDir: install,
		DataDir:    data,
		BinDir:     bin,
		ConfigDir:  configDir,
		CacheDir:   cacheDir,
	}, nil
}

func resolvePortable(configDir, cacheDir string, opts Options) (Paths, error) {
	root := opts.InstallHint
	if root == "" || root == "auto" {
		return Paths{}, fmt.Errorf("portable layout requires explicit install path")
	}

	root = filepath.Clean(root)

	return Paths{
		InstallDir: root,
		DataDir:    filepath.Join(root, "data"),
		BinDir:     filepath.Join(root, "bin"),
		ConfigDir:  filepath.Join(root, "config"),
		CacheDir:   filepath.Join(root, "cache"),
	}, nil
}
