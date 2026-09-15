package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

var (
	validFolderNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	allowedVideoExts     = map[string]bool{
		".mp4":  true,
		".mkv":  true,
		".avi":  true,
		".mov":  true,
		".webm": true,
	}
)

// MediaService orchestrates media folder and video asset lifecycle operations.
type MediaService struct {
	store ports.IMediaStore
}

// NewMediaService creates a new MediaService instance.
func NewMediaService(store ports.IMediaStore) *MediaService {
	return &MediaService{store: store}
}

// ListSources lists all media folders and video files.
func (s *MediaService) ListSources(ctx context.Context) (*domain.MediaSourcesResponse, error) {
	return s.store.ListSources(ctx)
}

// CreateFolder sanitizes and creates a new media folder.
func (s *MediaService) CreateFolder(ctx context.Context, folderName string) error {
	folderName = strings.TrimSpace(folderName)
	if folderName == "" {
		return errors.New("folder name cannot be empty")
	}
	if !validFolderNameRegex.MatchString(folderName) {
		return errors.New("folder name contains invalid characters (alphanumeric, underscore, hyphen only)")
	}
	if folderName == "root" || folderName == "general" {
		return errors.New("reserved folder name")
	}
	return s.store.CreateFolder(ctx, folderName)
}

// DeleteFolder removes a media folder and its contents.
func (s *MediaService) DeleteFolder(ctx context.Context, folderName string) error {
	folderName = strings.TrimSpace(folderName)
	if folderName == "" || folderName == "root" || folderName == "general" {
		return errors.New("cannot delete root/general media folder")
	}
	return s.store.DeleteFolder(ctx, folderName)
}

// SaveFile validates extension and saves a video file.
func (s *MediaService) SaveFile(ctx context.Context, folderName string, filename string, r io.Reader) (*domain.MediaFile, error) {
	folderName = strings.TrimSpace(folderName)
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedVideoExts[ext] {
		return nil, fmt.Errorf("unsupported video format '%s'; allowed: .mp4, .mkv, .avi, .mov, .webm", ext)
	}

	return s.store.SaveFile(ctx, folderName, filename, r)
}

// DeleteFile removes a specific video file.
func (s *MediaService) DeleteFile(ctx context.Context, folderName string, filename string) error {
	folderName = strings.TrimSpace(folderName)
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return errors.New("filename cannot be empty")
	}
	return s.store.DeleteFile(ctx, folderName, filename)
}

// OpenFile opens a media file for streaming.
func (s *MediaService) OpenFile(ctx context.Context, folderName string, filename string) (*os.File, os.FileInfo, error) {
	folderName = strings.TrimSpace(folderName)
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, nil, errors.New("filename cannot be empty")
	}
	return s.store.OpenFile(ctx, folderName, filename)
}
