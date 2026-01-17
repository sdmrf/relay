package downloader

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract auto-detects the archive format and extracts to dest.
// Supports .tar.gz, .tgz, and .zip formats.
func Extract(src, dest string) error {
	switch {
	case strings.HasSuffix(src, ".tar.gz"), strings.HasSuffix(src, ".tgz"):
		return ExtractTarGz(src, dest)
	case strings.HasSuffix(src, ".zip"):
		return ExtractZip(src, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", src)
	}
}

// ExtractTarGz extracts a .tar.gz archive to dest.
// Creates dest if it doesn't exist.
func ExtractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read: %w", err)
		}

		target, err := sanitizePath(dest, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}

		case tar.TypeReg:
			if err := extractFile(tr, target, header.FileInfo().Mode()); err != nil {
				return err
			}

		case tar.TypeSymlink:
			// Validate symlink target doesn't escape
			linkTarget := header.Linkname
			if filepath.IsAbs(linkTarget) {
				return fmt.Errorf("absolute symlink not allowed: %s", linkTarget)
			}
			targetDir := filepath.Dir(target)
			resolved := filepath.Join(targetDir, linkTarget)
			if !strings.HasPrefix(filepath.Clean(resolved), filepath.Clean(dest)) {
				return fmt.Errorf("symlink escape attempt: %s -> %s", header.Name, linkTarget)
			}
			if err := os.Symlink(linkTarget, target); err != nil {
				return fmt.Errorf("create symlink: %w", err)
			}
		}
	}

	return nil
}

// ExtractZip extracts a .zip archive to dest.
// Creates dest if it doesn't exist.
func ExtractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		target, err := sanitizePath(dest, f.Name)
		if err != nil {
			return err
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create parent directory: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open file in zip: %w", err)
		}

		err = extractFile(rc, target, f.Mode())
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// sanitizePath ensures the file path doesn't escape the destination directory.
// Prevents zip-slip vulnerability.
func sanitizePath(dest, name string) (string, error) {
	// Reject absolute paths immediately
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("illegal file path: %s", name)
	}

	// Clean the name to handle any .. sequences
	cleanName := filepath.Clean(name)

	// Reject paths that start with .. after cleaning
	if strings.HasPrefix(cleanName, ".."+string(os.PathSeparator)) || cleanName == ".." {
		return "", fmt.Errorf("illegal file path: %s", name)
	}

	target := filepath.Join(dest, cleanName)

	// Final check: ensure the target is within dest (belt and suspenders)
	cleanDest := filepath.Clean(dest)
	cleanTarget := filepath.Clean(target)
	if !strings.HasPrefix(cleanTarget+string(os.PathSeparator), cleanDest+string(os.PathSeparator)) &&
		cleanTarget != cleanDest {
		return "", fmt.Errorf("illegal file path: %s", name)
	}

	return target, nil
}

// extractFile writes content from r to target with the given mode.
func extractFile(r io.Reader, target string, mode os.FileMode) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}

	// Preserve executable bit, ensure at least 0644
	if mode == 0 {
		mode = 0o644
	}
	if mode&0o111 != 0 {
		mode = mode | 0o755
	}

	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	_, err = io.Copy(f, r)
	if closeErr := f.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
