package downloader

import "context"

// Artifact describes a downloadable resource.
type Artifact struct {
	Name     string
	URL      string
	Checksum string
	Target   string
}

// Downloader fetches artifacts from remote sources.
type Downloader interface {
	Fetch(ctx context.Context, a Artifact) error
}
