package launcher

// DesktopShortcut represents a platform-specific desktop shortcut.
type DesktopShortcut interface {
	// Create generates the desktop shortcut.
	Create() error
	// Remove deletes the desktop shortcut.
	Remove() error
	// Path returns the shortcut file path.
	Path() string
}

// ShortcutConfig contains configuration for creating desktop shortcuts.
type ShortcutConfig struct {
	Name        string // Display name (e.g., "Burp Suite Professional")
	Description string // Short description
	LauncherPath string // Path to the launcher script
	IconPath    string // Path to the icon (optional)
	Categories  string // Desktop categories (Linux only)
}
