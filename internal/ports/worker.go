package ports

import (
	"context"
	"hydraforge/internal/domain"
)

// TrainingWorker defines the interface for running Python/PyTorch training processes.
type TrainingWorker interface {
	StartTraining(ctx context.Context, job *domain.TrainingJob, metricCallback func(domain.TrainingMetrics)) error
	StopTraining(ctx context.Context, jobID string) error
	ExportWeights(ctx context.Context, weightsPath string, format domain.ExportFormat, precision string) (string, error)
	IsJobRunning(jobID string) bool
}
