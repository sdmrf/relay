package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// HTTPDownloader fetches artifacts over HTTP/HTTPS.
type HTTPDownloader struct {
	Timeout time.Duration
	Retries int
}

// Fetch downloads an artifact with retry support.
// Uses atomic writes (.tmp → rename) for safety.
func (d HTTPDownloader) Fetch(ctx context.Context, a Artifact) error {
	client := &http.Client{
		Timeout: d.Timeout,
	}

	var lastErr error

	for i := 0; i <= d.Retries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.URL, nil)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		err = d.writeResponse(resp, a.Target)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return lastErr
}

// writeResponse writes the response body to target atomically.
func (d HTTPDownloader) writeResponse(resp *http.Response, target string) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	tmp := target + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	_, err = io.Copy(out, resp.Body)
	if closeErr := out.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	if err != nil {
		os.Remove(tmp)
		return fmt.Errorf("write file: %w", err)
	}

	if err := os.Rename(tmp, target); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}
