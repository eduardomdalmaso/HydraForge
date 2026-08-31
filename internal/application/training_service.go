package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// TrainingService coordinates training workflows between ports.
type TrainingService struct {
	mu           sync.RWMutex
	jobRepo      ports.JobRepository
	datasetRepo  ports.DatasetRepository
	ckptRepo     ports.CheckpointRepository
	worker       ports.TrainingWorker
	gpuProvider  ports.GPUProvider
	metricSubscribers []chan domain.TrainingMetrics
	recentMetrics     []domain.TrainingMetrics
}

// NewTrainingService initializes a new TrainingService instance.
func NewTrainingService(
	jobRepo ports.JobRepository,
	datasetRepo ports.DatasetRepository,
	ckptRepo ports.CheckpointRepository,
	worker ports.TrainingWorker,
	gpuProvider ports.GPUProvider,
) *TrainingService {
	return &TrainingService{
		jobRepo:       jobRepo,
		datasetRepo:   datasetRepo,
		ckptRepo:      ckptRepo,
		worker:        worker,
		gpuProvider:   gpuProvider,
		recentMetrics: make([]domain.TrainingMetrics, 0, 100),
	}
}

// CreateTrainingJob validates, registers, and triggers a new training execution.
func (s *TrainingService) CreateTrainingJob(ctx context.Context, job *domain.TrainingJob) (*domain.TrainingJob, error) {
	if job == nil {
		return nil, domain.ErrInvalidJobConfig
	}
	job.SetDefaults()
	if err := job.Validate(); err != nil {
		return nil, err
	}

	if err := s.jobRepo.SaveJob(ctx, job); err != nil {
		return nil, err
	}

	// Trigger asynchronous training loop via worker port
	go func(targetJob domain.TrainingJob) {
		_ = s.jobRepo.UpdateJobStatus(context.Background(), targetJob.JobID, domain.StatusTraining, "")
		
		err := s.worker.StartTraining(context.Background(), &targetJob, func(m domain.TrainingMetrics) {
			s.mu.Lock()
			s.recentMetrics = append(s.recentMetrics, m)
			if len(s.recentMetrics) > 100 {
				s.recentMetrics = s.recentMetrics[1:]
			}
			s.mu.Unlock()
			s.broadcastMetric(m)
		})

		if err != nil {
			_ = s.jobRepo.UpdateJobStatus(context.Background(), targetJob.JobID, domain.StatusFailed, err.Error())
		} else {
			_ = s.jobRepo.UpdateJobStatus(context.Background(), targetJob.JobID, domain.StatusCompleted, "")
		}
	}(*job)

	return job, nil
}

// GetTrainingJob retrieves a specific job by ID.
func (s *TrainingService) GetTrainingJob(ctx context.Context, jobID string) (*domain.TrainingJob, error) {
	return s.jobRepo.GetJob(ctx, jobID)
}

// ListTrainingJobs lists jobs filtered by status.
func (s *TrainingService) ListTrainingJobs(ctx context.Context, status string) ([]*domain.TrainingJob, error) {
	return s.jobRepo.ListJobs(ctx, status)
}

// StopTrainingJob cancels an active training job.
func (s *TrainingService) StopTrainingJob(ctx context.Context, jobID string) error {
	if err := s.worker.StopTraining(ctx, jobID); err != nil {
		return err
	}
	return s.jobRepo.UpdateJobStatus(ctx, jobID, domain.StatusStopped, "Manually stopped by operator")
}

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
	var activeJob *domain.TrainingJob
	for _, j := range jobs {
		if j.Status == domain.StatusTraining {
			activeCount++
			activeJob = j
		}
	}

	s.mu.RLock()
	metricsCopy := make([]domain.TrainingMetrics, len(s.recentMetrics))
	copy(metricsCopy, s.recentMetrics)
	s.mu.RUnlock()

	return &domain.CockpitTelemetry{
		ActiveJobsCount:    activeCount,
		TotalModelsTrained: len(jobs),
		GPUStats:           *gpu,
		ActiveJob:          activeJob,
		RecentMetrics:      metricsCopy,
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

// Ensure interface compliance
var _ ports.TrainingUseCase = (*TrainingService)(nil)
