package application

import (
	"context"
	"fmt"
	"time"

	"hydraforge/internal/domain"
)

// ListCheckpoints returns all checkpoints for a job.
func (s *TrainingService) ListCheckpoints(ctx context.Context, jobID string) ([]*domain.ModelCheckpoint, error) {
	return s.ckptRepo.ListCheckpoints(ctx, jobID)
}

// ExportModel converts a checkpoint into TensorRT or ONNX.
func (s *TrainingService) ExportModel(ctx context.Context, checkpointID string, format domain.ExportFormat, precision string) (*domain.ModelCheckpoint, error) {
	ckpt, err := s.ckptRepo.GetCheckpoint(ctx, checkpointID)
	if err != nil {
		return nil, err
	}

	exportPath, err := s.worker.ExportWeights(ctx, ckpt.WeightsPath, format, precision)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrExportFailed, err)
	}

	exportedCkpt := &domain.ModelCheckpoint{
		CheckpointID: fmt.Sprintf("%s_%s", checkpointID, format),
		JobID:        ckpt.JobID,
		Architecture: ckpt.Architecture,
		Task:         ckpt.Task,
		Epoch:        ckpt.Epoch,
		MAP50:        ckpt.MAP50,
		MAP50_95:     ckpt.MAP50_95,
		WeightsPath:  ckpt.WeightsPath,
		ExportFormat: format,
		ExportPath:   exportPath,
		Precision:    precision,
		CreatedAt:    time.Now(),
	}

	if err := s.ckptRepo.SaveCheckpoint(ctx, exportedCkpt); err != nil {
		return nil, err
	}

	return exportedCkpt, nil
}

// GetCockpitTelemetry aggregates real-time telemetry.
func (s *TrainingService) GetCockpitTelemetry(ctx context.Context) (*domain.CockpitTelemetry, error) {
	gpu, err := s.gpuProvider.QueryGPU(ctx)
	if err != nil {
		gpu = &domain.GPUTelemetry{Model: "NVIDIA GeForce RTX 5090", TotalVRAMMB: 32768, Timestamp: time.Now()}
	}

	jobs, _ := s.jobRepo.ListJobs(ctx, "")
	activeCount := 0
	s.mu.RLock()
	var activeJob *domain.TrainingJob
	if s.activeJob != nil && s.activeJob.Status == domain.StatusTraining {
		jobCopy := *s.activeJob
		activeJob = &jobCopy
		activeCount = 1
	} else {
		for _, j := range jobs {
			if j.Status == domain.StatusTraining {
				activeCount++
				activeJob = j
			}
		}
	}

	metricsCopy := make([]domain.TrainingMetrics, len(s.recentMetrics))
	copy(metricsCopy, s.recentMetrics)
	logsCopy := make([]string, len(s.rawLogs))
	copy(logsCopy, s.rawLogs)
	s.mu.RUnlock()

	return &domain.CockpitTelemetry{
		ActiveJobsCount:    activeCount,
		TotalModelsTrained: len(jobs),
		GPUStats:           *gpu,
		ActiveJob:          activeJob,
		RecentMetrics:      metricsCopy,
		RawLogs:            logsCopy,
	}, nil
}

// GetGPUTelemetry queries live GPU readings.
func (s *TrainingService) GetGPUTelemetry(ctx context.Context) (*domain.GPUTelemetry, error) {
	return s.gpuProvider.QueryGPU(ctx)
}

func (s *TrainingService) broadcastMetric(m domain.TrainingMetrics) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range s.metricSubscribers {
		select {
		case ch <- m:
		default:
		}
	}
}
