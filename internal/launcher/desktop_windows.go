//go:build windows

package launcher

import (
	"fmt"
	"os"
	"path/filepath"
)

// WindowsShortcut creates .lnk shortcuts for Windows.
// Note: Creating proper .lnk files requires COM interfaces or external tools.
// This implementation creates a batch file as a simple alternative.
type WindowsShortcut struct {
	config  ShortcutConfig
	homeDir string
}

// NewWindowsShortcut creates a new Windows shortcut generator.
func NewWindowsShortcut(cfg ShortcutConfig) (*WindowsShortcut, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}
	return &WindowsShortcut{config: cfg, homeDir: home}, nil
}

func (w *WindowsShortcut) Path() string {
	return filepath.Join(w.homeDir, "Desktop", w.config.Name+".bat")
}

func (w *WindowsShortcut) Create() error {
	// Create a batch file launcher as a simple shortcut
	// For proper .lnk files, would need to use COM interfaces or PowerShell
	batch := fmt.Sprintf(`@echo off
start "" "%s"
`, w.config.LauncherPath)

	if err := os.WriteFile(w.Path(), []byte(batch), 0o644); err != nil {
		return fmt.Errorf("write batch file: %w", err)
	}

	return nil
}

func (w *WindowsShortcut) Remove() error {
	return os.Remove(w.Path())
}
