//go:build linux

package launcher

import (
	"fmt"
	"os"
	"path/filepath"
)

// LinuxShortcut creates .desktop files for Linux.
type LinuxShortcut struct {
	config  ShortcutConfig
	homeDir string
}

// NewLinuxShortcut creates a new Linux shortcut generator.
func NewLinuxShortcut(cfg ShortcutConfig) (*LinuxShortcut, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}
	return &LinuxShortcut{config: cfg, homeDir: home}, nil
}

func (l *LinuxShortcut) Path() string {
	return filepath.Join(l.homeDir, ".local", "share", "applications", l.config.Name+".desktop")
}

func (l *LinuxShortcut) Create() error {
	// Ensure directory exists
	dir := filepath.Dir(l.Path())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create applications dir: %w", err)
	}

	categories := l.config.Categories
	if categories == "" {
		categories = "Development;Security;"
	}

	icon := l.config.IconPath
	if icon == "" {
		icon = "applications-security"
	}

	desktop := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Comment=%s
Exec=%s
Icon=%s
Terminal=false
Categories=%s
`, l.config.Name, l.config.Description, l.config.LauncherPath, icon, categories)

	if err := os.WriteFile(l.Path(), []byte(desktop), 0o644); err != nil {
		return fmt.Errorf("write desktop file: %w", err)
	}

	return nil
}

func (l *LinuxShortcut) Remove() error {
	return os.Remove(l.Path())
}
