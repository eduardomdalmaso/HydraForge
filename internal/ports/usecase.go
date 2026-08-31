package ports

import (
	"context"
	"hydraforge/internal/domain"
)

// TrainingUseCase defines primary business methods for the training studio.
type TrainingUseCase interface {
	CreateTrainingJob(ctx context.Context, job *domain.TrainingJob) (*domain.TrainingJob, error)
	GetTrainingJob(ctx context.Context, jobID string) (*domain.TrainingJob, error)
	ListTrainingJobs(ctx context.Context, status string) ([]*domain.TrainingJob, error)
	StopTrainingJob(ctx context.Context, jobID string) error
	
	// Dataset operations
	RegisterDataset(ctx context.Context, dataset *domain.Dataset) (*domain.Dataset, error)
	GetDataset(ctx context.Context, datasetID string) (*domain.Dataset, error)
	ListDatasets(ctx context.Context) ([]*domain.Dataset, error)

	// Checkpoints & Model Zoo
	ListCheckpoints(ctx context.Context, jobID string) ([]*domain.ModelCheckpoint, error)
	ExportModel(ctx context.Context, checkpointID string, format domain.ExportFormat, precision string) (*domain.ModelCheckpoint, error)

	// Telemetry & Hardware
	GetCockpitTelemetry(ctx context.Context) (*domain.CockpitTelemetry, error)
	GetGPUTelemetry(ctx context.Context) (*domain.GPUTelemetry, error)
}
