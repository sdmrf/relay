package burpsuite

import (
	"path/filepath"

	"github.com/sdmrf/relay/internal/downloader"
)

// ArtifactJar returns the download artifact for the Burp Suite JAR.
func (b *BurpSuite) ArtifactJar() downloader.Artifact {
	return downloader.Artifact{
		Name:   "burpsuite.jar",
		URL:    burpDownloadURL(b.cfg.Product.Edition),
		Target: filepath.Join(b.paths.InstallDir, "burpsuite.jar"),
	}
}
