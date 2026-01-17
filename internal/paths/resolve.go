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
	var install, data, bin string

	switch os {
	case Linux:
		install = "/usr/lib/relay"
		data = "/var/lib/relay"
		bin = "/usr/bin"
	case Darwin:
		install = "/Applications/relay"
		data = "/Library/Application Support/relay"
		bin = "/usr/local/bin"
	case Windows:
		install = `C:\Program Files\relay`
		data = `C:\ProgramData\relay`
		bin = `C:\Program Files\relay\bin`
	default:
		return Paths{}, fmt.Errorf("unsupported OS: %s", os)
	}

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
