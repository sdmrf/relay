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
	Timeout    time.Duration
	Retries    int
	OnProgress ProgressFunc // Optional progress callback
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

// FetchWithProgress downloads an artifact with a progress bar.
func (d HTTPDownloader) FetchWithProgress(ctx context.Context, a Artifact) error {
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

		err = d.writeResponseWithProgress(resp, a.Target, a.Name)
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

	var reader io.Reader = resp.Body
	if d.OnProgress != nil {
		reader = &progressReader{
			reader:     resp.Body,
			total:      resp.ContentLength,
			onProgress: d.OnProgress,
		}
	}

	_, err = io.Copy(out, reader)
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

// writeResponseWithProgress writes the response body with a progress bar.
func (d HTTPDownloader) writeResponseWithProgress(resp *http.Response, target, name string) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	tmp := target + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	// Create progress bar
	bar := NewProgressBar(name, resp.ContentLength)
	reader := &progressReader{
		reader: resp.Body,
		total:  resp.ContentLength,
		onProgress: func(downloaded, total int64) {
			bar.Update(downloaded)
		},
	}

	_, err = io.Copy(out, reader)
	bar.Finish()

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
