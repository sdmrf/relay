//go:build darwin

package launcher

import (
	"fmt"
	"os"
	"path/filepath"
)

// DarwinShortcut creates macOS .app bundles or aliases.
type DarwinShortcut struct {
	config  ShortcutConfig
	homeDir string
}

// NewDarwinShortcut creates a new macOS shortcut generator.
func NewDarwinShortcut(cfg ShortcutConfig) (*DarwinShortcut, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}
	return &DarwinShortcut{config: cfg, homeDir: home}, nil
}

func (d *DarwinShortcut) Path() string {
	return filepath.Join(d.homeDir, "Applications", d.config.Name+".app")
}

func (d *DarwinShortcut) Create() error {
	appBundle := d.Path()
	contentsDir := filepath.Join(appBundle, "Contents")
	macosDir := filepath.Join(contentsDir, "MacOS")

	// Create directory structure
	if err := os.MkdirAll(macosDir, 0o755); err != nil {
		return fmt.Errorf("create app bundle: %w", err)
	}

	// Create Info.plist
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>launcher</string>
	<key>CFBundleIdentifier</key>
	<string>com.relay.%s</string>
	<key>CFBundleName</key>
	<string>%s</string>
	<key>CFBundleDisplayName</key>
	<string>%s</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
</dict>
</plist>
`, d.config.Name, d.config.Name, d.config.Name)

	plistPath := filepath.Join(contentsDir, "Info.plist")
	if err := os.WriteFile(plistPath, []byte(plist), 0o644); err != nil {
		return fmt.Errorf("write Info.plist: %w", err)
	}

	// Create launcher script
	launcherScript := fmt.Sprintf(`#!/bin/sh
exec "%s" "$@"
`, d.config.LauncherPath)

	launcherPath := filepath.Join(macosDir, "launcher")
	if err := os.WriteFile(launcherPath, []byte(launcherScript), 0o755); err != nil {
		return fmt.Errorf("write launcher: %w", err)
	}

	return nil
}

func (d *DarwinShortcut) Remove() error {
	return os.RemoveAll(d.Path())
}
