package ports

import (
	"context"
	"io"
	"os"

	"hydraforge/internal/domain"
)

// IMediaStore defines the outbound port for managing media directories and video assets.
type IMediaStore interface {
	ListSources(ctx context.Context) (*domain.MediaSourcesResponse, error)
	CreateFolder(ctx context.Context, folderName string) error
	DeleteFolder(ctx context.Context, folderName string) error
	SaveFile(ctx context.Context, folderName string, filename string, r io.Reader) (*domain.MediaFile, error)
	DeleteFile(ctx context.Context, folderName string, filename string) error
	OpenFile(ctx context.Context, folderName string, filename string) (*os.File, os.FileInfo, error)
}
