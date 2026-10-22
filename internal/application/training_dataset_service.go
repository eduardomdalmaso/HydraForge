package application

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"hydraforge/internal/domain"
)

// RegisterDataset saves dataset metadata.
func (s *TrainingService) RegisterDataset(ctx context.Context, dataset *domain.Dataset) (*domain.Dataset, error) {
	if dataset == nil {
		return nil, domain.ErrInvalidDataset
	}
	if dataset.CreatedAt.IsZero() {
		dataset.CreatedAt = time.Now()
	}
	if err := dataset.Validate(); err != nil {
		return nil, err
	}
	if err := s.datasetRepo.SaveDataset(ctx, dataset); err != nil {
		return nil, err
	}
	return dataset, nil
}

// GetDataset retrieves a dataset by ID.
func (s *TrainingService) GetDataset(ctx context.Context, datasetID string) (*domain.Dataset, error) {
	return s.datasetRepo.GetDataset(ctx, datasetID)
}

// ListDatasets retrieves all registered datasets.
func (s *TrainingService) ListDatasets(ctx context.Context) ([]*domain.Dataset, error) {
	return s.datasetRepo.ListDatasets(ctx)
}

// RescanDatasets triggers an on-disk discovery rescan and returns updated datasets.
func (s *TrainingService) RescanDatasets(ctx context.Context) ([]*domain.Dataset, error) {
	if rescanable, ok := s.datasetRepo.(interface{ RescanDatasets() }); ok {
		rescanable.RescanDatasets()
	}
	return s.datasetRepo.ListDatasets(ctx)
}

// DeleteDataset removes a dataset from the repository and optionally removes its directory from disk.
func (s *TrainingService) DeleteDataset(ctx context.Context, datasetID string, deleteFiles bool) error {
	ds, err := s.datasetRepo.GetDataset(ctx, datasetID)
	if err == nil {
		_ = s.datasetRepo.DeleteDataset(ctx, datasetID)
	}
	if deleteFiles && ds != nil && ds.YAMLPath != "" {
		dir := filepath.Dir(ds.YAMLPath)
		if dir != "" && dir != "/" && dir != "." {
			_ = os.RemoveAll(dir)
		}
	}
	return nil
}
