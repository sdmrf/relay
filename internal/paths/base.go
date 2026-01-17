package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

func baseDirs(currentOS OS) (config, cache string, err error) {
	home, err := osUserHome()
	if err != nil {
		return "", "", err
	}

	switch currentOS {
	case Linux:
		config = filepath.Join(home, ".config", "relay")
		cache = filepath.Join(home, ".cache", "relay")
	case Darwin:
		config = filepath.Join(home, "Library", "Application Support", "relay")
		cache = filepath.Join(home, "Library", "Caches", "relay")
	case Windows:
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", "", fmt.Errorf("APPDATA not set")
		}
		config = filepath.Join(appData, "relay")
		cache = filepath.Join(config, "cache")
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", currentOS)
	}

	return config, cache, nil
}

func osUserHome() (string, error) {
	return os.UserHomeDir()
}
