package ports

import (
	"context"
	"hydraforge/internal/domain"
)

// JobRepository handles persistence of training jobs.
type JobRepository interface {
	SaveJob(ctx context.Context, job *domain.TrainingJob) error
	GetJob(ctx context.Context, jobID string) (*domain.TrainingJob, error)
	ListJobs(ctx context.Context, status string) ([]*domain.TrainingJob, error)
	UpdateJobStatus(ctx context.Context, jobID string, status domain.JobStatus, errMsg string) error
	DeleteJob(ctx context.Context, jobID string) error
}

// DatasetRepository handles persistence of datasets.
type DatasetRepository interface {
	SaveDataset(ctx context.Context, dataset *domain.Dataset) error
	GetDataset(ctx context.Context, datasetID string) (*domain.Dataset, error)
	ListDatasets(ctx context.Context) ([]*domain.Dataset, error)
}

// CheckpointRepository handles persistence of model checkpoints and export artifacts.
type CheckpointRepository interface {
	SaveCheckpoint(ctx context.Context, checkpoint *domain.ModelCheckpoint) error
	GetCheckpoint(ctx context.Context, checkpointID string) (*domain.ModelCheckpoint, error)
	ListCheckpoints(ctx context.Context, jobID string) ([]*domain.ModelCheckpoint, error)
}
