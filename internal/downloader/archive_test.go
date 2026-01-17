package downloader

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTarGz(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	extractDir := filepath.Join(tmpDir, "extracted")

	// Create a test tar.gz archive
	if err := createTestTarGz(archivePath); err != nil {
		t.Fatalf("createTestTarGz() error = %v", err)
	}

	// Extract it
	if err := ExtractTarGz(archivePath, extractDir); err != nil {
		t.Fatalf("ExtractTarGz() error = %v", err)
	}

	// Verify extracted files
	verifyExtractedFiles(t, extractDir)
}

func TestExtractZip(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.zip")
	extractDir := filepath.Join(tmpDir, "extracted")

	// Create a test zip archive
	if err := createTestZip(archivePath); err != nil {
		t.Fatalf("createTestZip() error = %v", err)
	}

	// Extract it
	if err := ExtractZip(archivePath, extractDir); err != nil {
		t.Fatalf("ExtractZip() error = %v", err)
	}

	// Verify extracted files
	verifyExtractedFiles(t, extractDir)
}

func TestExtract(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		archiveName string
		createFunc  func(string) error
	}{
		{
			name:        "tar.gz",
			archiveName: "test.tar.gz",
			createFunc:  createTestTarGz,
		},
		{
			name:        "tgz",
			archiveName: "test.tgz",
			createFunc:  createTestTarGz,
		},
		{
			name:        "zip",
			archiveName: "test.zip",
			createFunc:  createTestZip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archivePath := filepath.Join(tmpDir, tt.archiveName)
			extractDir := filepath.Join(tmpDir, "extracted_"+tt.name)

			if err := tt.createFunc(archivePath); err != nil {
				t.Fatalf("create archive error = %v", err)
			}

			if err := Extract(archivePath, extractDir); err != nil {
				t.Fatalf("Extract() error = %v", err)
			}

			verifyExtractedFiles(t, extractDir)
		})
	}
}

func TestExtractUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.rar")
	extractDir := filepath.Join(tmpDir, "extracted")

	// Create a dummy file
	if err := os.WriteFile(archivePath, []byte("dummy"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := Extract(archivePath, extractDir)
	if err == nil {
		t.Error("Extract() expected error for unsupported format, got nil")
	}
}

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name    string
		dest    string
		path    string
		wantErr bool
	}{
		{
			name:    "normal path",
			dest:    "/tmp/extract",
			path:    "file.txt",
			wantErr: false,
		},
		{
			name:    "nested path",
			dest:    "/tmp/extract",
			path:    "dir/subdir/file.txt",
			wantErr: false,
		},
		{
			name:    "zip slip attack",
			dest:    "/tmp/extract",
			path:    "../../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "absolute path",
			dest:    "/tmp/extract",
			path:    "/etc/passwd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sanitizePath(tt.dest, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("sanitizePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractPreservesExecutableBit(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	extractDir := filepath.Join(tmpDir, "extracted")

	// Create archive with executable file
	if err := createTarGzWithExecutable(archivePath); err != nil {
		t.Fatalf("createTarGzWithExecutable() error = %v", err)
	}

	// Extract
	if err := ExtractTarGz(archivePath, extractDir); err != nil {
		t.Fatalf("ExtractTarGz() error = %v", err)
	}

	// Check executable bit is preserved
	execPath := filepath.Join(extractDir, "bin", "java")
	info, err := os.Stat(execPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	// Check if executable bit is set
	if info.Mode()&0o111 == 0 {
		t.Errorf("executable bit not preserved: mode = %o", info.Mode())
	}
}

// Helper functions to create test archives

func createTestTarGz(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Add a directory
	if err := tw.WriteHeader(&tar.Header{
		Name:     "testdir/",
		Mode:     0o755,
		Typeflag: tar.TypeDir,
	}); err != nil {
		return err
	}

	// Add a file
	content := []byte("test content")
	if err := tw.WriteHeader(&tar.Header{
		Name: "testdir/file.txt",
		Mode: 0o644,
		Size: int64(len(content)),
	}); err != nil {
		return err
	}
	if _, err := tw.Write(content); err != nil {
		return err
	}

	return nil
}

func createTestZip(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	// Add a directory
	if _, err := zw.Create("testdir/"); err != nil {
		return err
	}

	// Add a file
	w, err := zw.Create("testdir/file.txt")
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte("test content")); err != nil {
		return err
	}

	return nil
}

func createTarGzWithExecutable(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Add bin directory
	if err := tw.WriteHeader(&tar.Header{
		Name:     "bin/",
		Mode:     0o755,
		Typeflag: tar.TypeDir,
	}); err != nil {
		return err
	}

	// Add executable file
	content := []byte("#!/bin/sh\necho hello")
	if err := tw.WriteHeader(&tar.Header{
		Name: "bin/java",
		Mode: 0o755,
		Size: int64(len(content)),
	}); err != nil {
		return err
	}
	if _, err := tw.Write(content); err != nil {
		return err
	}

	return nil
}

func verifyExtractedFiles(t *testing.T, extractDir string) {
	t.Helper()

	// Check directory exists
	dirPath := filepath.Join(extractDir, "testdir")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("directory not extracted: %s", dirPath)
	}

	// Check file exists and has correct content
	filePath := filepath.Join(extractDir, "testdir", "file.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	if string(content) != "test content" {
		t.Errorf("file content = %q, want %q", string(content), "test content")
	}
}
