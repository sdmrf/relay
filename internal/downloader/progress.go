package downloader

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// ProgressFunc is called with download progress updates.
type ProgressFunc func(downloaded, total int64)

// progressReader wraps an io.Reader and reports progress.
type progressReader struct {
	reader     io.Reader
	total      int64
	downloaded int64
	onProgress ProgressFunc
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.downloaded += int64(n)
	if pr.onProgress != nil {
		pr.onProgress(pr.downloaded, pr.total)
	}
	return n, err
}

// ProgressBar displays a terminal progress bar.
type ProgressBar struct {
	name       string
	total      int64
	current    int64
	width      int
	mu         sync.Mutex
	lastUpdate time.Time
	started    time.Time
}

// NewProgressBar creates a new progress bar.
func NewProgressBar(name string, total int64) *ProgressBar {
	return &ProgressBar{
		name:    name,
		total:   total,
		width:   40,
		started: time.Now(),
	}
}

// Update updates the progress bar with the current value.
func (p *ProgressBar) Update(current int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = current

	// Throttle updates to avoid flickering (max 10 updates/sec)
	if time.Since(p.lastUpdate) < 100*time.Millisecond && current < p.total {
		return
	}
	p.lastUpdate = time.Now()

	p.render()
}

// Finish completes the progress bar.
func (p *ProgressBar) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = p.total
	p.render()
	fmt.Println() // New line after progress bar
}

func (p *ProgressBar) render() {
	var percent float64
	if p.total > 0 {
		percent = float64(p.current) / float64(p.total) * 100
	}

	// Calculate filled width
	filled := int(float64(p.width) * percent / 100)
	if filled > p.width {
		filled = p.width
	}

	// Build progress bar
	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)

	// Format sizes
	currentStr := formatBytes(p.current)
	totalStr := formatBytes(p.total)

	// Calculate speed
	elapsed := time.Since(p.started).Seconds()
	var speedStr string
	if elapsed > 0 {
		speed := float64(p.current) / elapsed
		speedStr = formatBytes(int64(speed)) + "/s"
	}

	// Print progress bar (carriage return to overwrite)
	fmt.Printf("\r%s [%s] %5.1f%% %s/%s %s",
		truncate(p.name, 20),
		bar,
		percent,
		currentStr,
		totalStr,
		speedStr,
	)
}

// formatBytes formats bytes as human-readable string.
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// truncate shortens a string to maxLen, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}
