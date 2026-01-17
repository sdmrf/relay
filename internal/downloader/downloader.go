package downloader

import "context"

// Artifact describes a downloadable resource.
type Artifact struct {
	Name     string // Display name (e.g., "burpsuite.jar")
	URL      string // Download URL
	Checksum string // Expected checksum for verification (optional, future use)
	Target   string // Target file path
}

// Downloader fetches artifacts from remote sources.
type Downloader interface {
	Fetch(ctx context.Context, a Artifact) error
}
