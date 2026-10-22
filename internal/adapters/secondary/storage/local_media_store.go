package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

var _ ports.IMediaStore = (*LocalMediaStore)(nil)

// LocalMediaStore implements IMediaStore on the local filesystem.
type LocalMediaStore struct {
	baseDir string
}

// NewLocalMediaStore initializes a LocalMediaStore at baseDir.
func NewLocalMediaStore(baseDir string) (*LocalMediaStore, error) {
	abs, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("invalid media base directory: %w", err)
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return nil, fmt.Errorf("failed to create media base directory: %w", err)
	}
	return &LocalMediaStore{baseDir: abs}, nil
}

func (s *LocalMediaStore) resolvePath(folderName, filename string) (string, error) {
	folderName = strings.TrimSpace(folderName)
	if folderName == "root" || folderName == "general" || folderName == "geral" {
		folderName = ""
	}

	var targetPath string
	if folderName == "" {
		targetPath = filepath.Join(s.baseDir, filename)
	} else {
		targetPath = filepath.Join(s.baseDir, folderName, filename)
	}

	clean := filepath.Clean(targetPath)
	if !strings.HasPrefix(clean, s.baseDir) {
		return "", errors.New("directory traversal detected")
	}
	return clean, nil
}

// CreateFolder creates a new category directory under baseDir.
func (s *LocalMediaStore) CreateFolder(ctx context.Context, folderName string) error {
	folderPath, err := s.resolvePath(folderName, "")
	if err != nil {
		return err
	}
	return os.MkdirAll(folderPath, 0755)
}

// DeleteFolder removes a directory and all its files.
func (s *LocalMediaStore) DeleteFolder(ctx context.Context, folderName string) error {
	folderPath, err := s.resolvePath(folderName, "")
	if err != nil {
		return err
	}
	return os.RemoveAll(folderPath)
}

// SaveFile writes incoming stream to disk.
func (s *LocalMediaStore) SaveFile(ctx context.Context, folderName string, filename string, r io.Reader) (*domain.MediaFile, error) {
	filePath, err := s.resolvePath(folderName, filename)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer out.Close()

	n, err := io.Copy(out, r)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	info, err := out.Stat()
	if err != nil {
		return nil, err
	}

	folderTag := folderName
	if folderTag == "" {
		folderTag = "root"
	}

	return &domain.MediaFile{
		ID:        folderTag + "/" + filename,
		Name:      filename,
		Folder:    folderTag,
		Path:      filePath,
		SizeBytes: n,
		ModTime:   info.ModTime(),
		StreamURL: fmt.Sprintf("/api/v1/media/stream/%s/%s", folderTag, filename),
	}, nil
}

// DeleteFile deletes a file from disk.
func (s *LocalMediaStore) DeleteFile(ctx context.Context, folderName string, filename string) error {
	filePath, err := s.resolvePath(folderName, filename)
	if err != nil {
		return err
	}
	return os.Remove(filePath)
}

// OpenFile opens a file for streaming with http.ServeContent.
func (s *LocalMediaStore) OpenFile(ctx context.Context, folderName string, filename string) (*os.File, os.FileInfo, error) {
	filePath, err := s.resolvePath(folderName, filename)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return f, info, nil
}

func isVideoExtension(ext string) bool {
	switch ext {
	case ".mp4", ".mkv", ".avi", ".mov", ".webm", ".m4v":
		return true
	default:
		return false
	}
}
