package application

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// VaultDatasetService handles curated dataset lifecycle, class curation, and YOLO export.
type VaultDatasetService struct {
	datasetRepo ports.IVaultDatasetRepository
	frameRepo   ports.IFrameRepository
	exporter    ports.IYOLOExporter
}

// NewVaultDatasetService creates a new VaultDatasetService instance.
func NewVaultDatasetService(
	datasetRepo ports.IVaultDatasetRepository,
	frameRepo ports.IFrameRepository,
	exporter ports.IYOLOExporter,
) *VaultDatasetService {
	return &VaultDatasetService{
		datasetRepo: datasetRepo,
		frameRepo:   frameRepo,
		exporter:    exporter,
	}
}

// CreateDataset initializes a new dataset template.
func (s *VaultDatasetService) CreateDataset(ctx context.Context, d *domain.VaultDataset) error {
	if err := d.Validate(); err != nil {
		return err
	}
	return s.datasetRepo.Save(ctx, d)
}

// GetDataset retrieves a curated dataset by ID.
func (s *VaultDatasetService) GetDataset(ctx context.Context, datasetID string) (*domain.VaultDataset, error) {
	return s.datasetRepo.GetByID(ctx, datasetID)
}

// ListDatasets returns all registered curated datasets.
func (s *VaultDatasetService) ListDatasets(ctx context.Context) ([]*domain.VaultDataset, error) {
	return s.datasetRepo.List(ctx)
}

// ListFrames lists frames for a dataset with pagination.
func (s *VaultDatasetService) ListFrames(ctx context.Context, filter ports.FrameFilter) ([]*domain.Frame, int, error) {
	frames, err := s.frameRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.frameRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return frames, count, nil
}

// ExportDataset splits frames into train/val (e.g. 80/20) and generates data.yaml.
func (s *VaultDatasetService) ExportDataset(ctx context.Context, datasetID string, valRatio float64) (*ports.ExportResult, error) {
	d, err := s.datasetRepo.GetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	frames, err := s.frameRepo.List(ctx, ports.FrameFilter{
		DatasetID: datasetID,
		Limit:     100000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch frames for export: %w", err)
	}
	if len(frames) == 0 {
		return nil, fmt.Errorf("%w: dataset has no frames", domain.ErrExportFailed)
	}

	if valRatio <= 0.0 || valRatio >= 1.0 {
		valRatio = 0.2 // Default 20% validation split
	}

	// Deterministic split assignment
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(frames), func(i, j int) { frames[i], frames[j] = frames[j], frames[i] })

	valThreshold := int(float64(len(frames)) * valRatio)
	for i, f := range frames {
		if i < valThreshold {
			f.Split = domain.SplitVal
		} else {
			f.Split = domain.SplitTrain
		}
		_ = s.frameRepo.UpdateStatus(ctx, f.FrameID, domain.FrameStatusExported, f.Split)
	}

	exportRes, err := s.exporter.Export(ctx, d, frames)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrExportFailed, err)
	}

	// Update dataset aggregates
	d.TrainImages = exportRes.TrainCount
	d.ValImages = exportRes.ValCount
	d.TestImages = exportRes.TestCount
	d.TotalFrames = len(frames)
	d.YAMLPath = exportRes.YAMLPath
	d.ChecksumSHA = exportRes.ChecksumSHA
	_ = s.datasetRepo.Save(ctx, d)

	return exportRes, nil
}
