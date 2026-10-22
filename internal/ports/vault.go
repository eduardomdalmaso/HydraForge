package ports

import (
	"context"
	"image"
	"io"

	"hydraforge/internal/domain"
)

// IVaultDatasetRepository defines persistence operations for curated datasets.
type IVaultDatasetRepository interface {
	Save(ctx context.Context, dataset *domain.VaultDataset) error
	GetByID(ctx context.Context, datasetID string) (*domain.VaultDataset, error)
	List(ctx context.Context) ([]*domain.VaultDataset, error)
	Delete(ctx context.Context, datasetID string) error
}

// FrameFilter provides criteria for listing frames.
type FrameFilter struct {
	DatasetID string
	CameraID  string
	Status    domain.FrameStatus
	Split     domain.SplitType
	Limit     int
	Offset    int
}

// IFrameRepository defines persistence operations for frames.
type IFrameRepository interface {
	Save(ctx context.Context, frame *domain.Frame) error
	GetByID(ctx context.Context, frameID string) (*domain.Frame, error)
	List(ctx context.Context, filter FrameFilter) ([]*domain.Frame, error)
	Count(ctx context.Context, filter FrameFilter) (int, error)
	FindByPHash(ctx context.Context, phash uint64, maxHammingDistance int) (*domain.Frame, error)
	UpdateStatus(ctx context.Context, frameID string, status domain.FrameStatus, split domain.SplitType) error
	Delete(ctx context.Context, frameID string) error
}

// IDeduplicator computes perceptual hashes and measures similarity distances.
type IDeduplicator interface {
	ComputePHash(img image.Image) (uint64, error)
	HammingDistance(hashA, hashB uint64) int
}

// IFileStore provides sandboxed, path-traversal safe storage for image blobs.
type IFileStore interface {
	Save(ctx context.Context, relativePath string, r io.Reader) (string, int64, error)
	Open(ctx context.Context, relativePath string) (io.ReadCloser, error)
	Exists(ctx context.Context, relativePath string) bool
	Delete(ctx context.Context, relativePath string) error
	GetAbsolutePath(relativePath string) (string, error)
}

// ExportResult contains paths and hashes of the exported dataset.
type ExportResult struct {
	DatasetID   string
	YAMLPath    string
	OutputDir   string
	TrainCount  int
	ValCount    int
	TestCount   int
	ChecksumSHA string
}

// IYOLOExporter generates normalized .txt labels and data.yaml for Ultralytics/HydraForge.
type IYOLOExporter interface {
	Export(ctx context.Context, dataset *domain.VaultDataset, frames []*domain.Frame) (*ExportResult, error)
}
