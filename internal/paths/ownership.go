package paths

// Owned returns paths that relay owns and can safely modify/delete.
// Excludes ConfigDir to preserve user configuration.
// Excludes BinDir as it may be a shared system directory (e.g., /usr/local/bin).
func (p Paths) Owned() []string {
	return []string{
		p.InstallDir,
		p.DataDir,
		p.CacheDir,
	}
}
