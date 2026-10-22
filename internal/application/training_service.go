package application

import (
	"context"
	"fmt"
	"sync"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// TrainingService coordinates training workflows between ports.
type TrainingService struct {
	mu                sync.RWMutex
	jobRepo           ports.JobRepository
	datasetRepo       ports.DatasetRepository
	ckptRepo          ports.CheckpointRepository
	benchmarkRepo     ports.BenchmarkRepository
	worker            ports.TrainingWorker
	gpuProvider       ports.GPUProvider
	metricSubscribers []chan domain.TrainingMetrics
	recentMetrics     []domain.TrainingMetrics
	rawLogs           []string
	activeJob         *domain.TrainingJob
}

// NewTrainingService initializes a new TrainingService instance.
func NewTrainingService(
	jobRepo ports.JobRepository,
	datasetRepo ports.DatasetRepository,
	ckptRepo ports.CheckpointRepository,
	benchmarkRepo ports.BenchmarkRepository,
	worker ports.TrainingWorker,
	gpuProvider ports.GPUProvider,
) *TrainingService {
	if jobs, err := jobRepo.ListJobs(context.Background(), ""); err == nil {
		for _, j := range jobs {
			if j.Status == domain.StatusTraining || j.Status == domain.StatusQueued {
				_ = jobRepo.UpdateJobStatus(context.Background(), j.JobID, domain.StatusStopped, "Interrupted by server restart")
			}
		}
	}

	return &TrainingService{
		jobRepo:       jobRepo,
		datasetRepo:   datasetRepo,
		ckptRepo:      ckptRepo,
		benchmarkRepo: benchmarkRepo,
		worker:        worker,
		gpuProvider:   gpuProvider,
		recentMetrics: make([]domain.TrainingMetrics, 0, 100),
		rawLogs:       make([]string, 0, 100),
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

	go func(targetJob domain.TrainingJob) {
		targetJob.Status = domain.StatusTraining
		_ = s.jobRepo.UpdateJobStatus(context.Background(), targetJob.JobID, domain.StatusTraining, "")

		initLog1 := fmt.Sprintf("[HydraForge] Training Job %s (%s) started on %s", targetJob.JobID, targetJob.ModelArchitecture, targetJob.DatasetID)
		initLog2 := fmt.Sprintf("[Configuration] Epochs: %d | Batch: %d | ImgSz: %d | Optimizer: %s (lr0=%v)", targetJob.Hyperparameters.Epochs, targetJob.Hyperparameters.BatchSize, targetJob.Hyperparameters.ImageSize, targetJob.Hyperparameters.Optimizer, targetJob.Hyperparameters.LearningRate)
		s.mu.Lock()
		s.rawLogs = []string{initLog1, initLog2}
		s.activeJob = &targetJob
		s.mu.Unlock()

		var peakVRAM float64

		err := s.worker.StartTraining(context.Background(), &targetJob, func(m domain.TrainingMetrics) {
			s.mu.Lock()
			if m.CurrentBatch > 0 {
				targetJob.CurrentBatch = m.CurrentBatch
			}
			if m.TotalBatches > 0 {
				targetJob.TotalBatches = m.TotalBatches
			}
			if m.Epoch > 0 {
				targetJob.CurrentEpoch = m.Epoch
			}
			s.activeJob = &targetJob
			if m.EpochDurationSec > 0 && m.MAP50 > 0 {
				s.recentMetrics = append(s.recentMetrics, m)
				if len(s.recentMetrics) > 100 {
					s.recentMetrics = s.recentMetrics[1:]
				}
				if m.MAP50 > targetJob.BestMAP50 {
					targetJob.BestMAP50 = m.MAP50
				}
				if m.MAP50_95 > targetJob.BestMAP50_95 {
					targetJob.BestMAP50_95 = m.MAP50_95
				}
				logLine := fmt.Sprintf("[Epoch %d/%d] Box Loss: %.4f | Cls: %.4f | mAP50: %.1f%% | %.0f FPS (%.1fs)",
					m.Epoch, targetJob.Hyperparameters.Epochs, m.BoxLoss, m.ClsLoss, m.MAP50*100.0, m.FPS, m.EpochDurationSec)
				s.rawLogs = append(s.rawLogs, logLine)
				if len(s.rawLogs) > 100 {
					s.rawLogs = s.rawLogs[1:]
				}
			} else if targetJob.CurrentBatch > 0 && targetJob.CurrentBatch%50 == 0 {
				totalLoss := m.BoxLoss + m.ClsLoss
				batchLog := fmt.Sprintf("[Batch %d/%d] Loss: %.4f (Box: %.4f, Cls: %.4f) | Power: %.0fW | %.0f FPS (%.1fs)",
					targetJob.CurrentBatch, targetJob.TotalBatches, totalLoss, m.BoxLoss, m.ClsLoss, m.PowerWatts, m.FPS, m.EpochDurationSec)
				s.rawLogs = append(s.rawLogs, batchLog)
				if len(s.rawLogs) > 100 {
					s.rawLogs = s.rawLogs[1:]
				}
			}

			if m.GPUVRAMMB > peakVRAM {
				peakVRAM = m.GPUVRAMMB
				targetJob.PeakVRAMMB = peakVRAM
			}

			_ = s.jobRepo.SaveJob(context.Background(), &targetJob)
			s.mu.Unlock()

			s.broadcastMetric(m)
		})

		if err != nil {
			_ = s.jobRepo.UpdateJobStatus(context.Background(), targetJob.JobID, domain.StatusFailed, err.Error())
			s.mu.Lock()
			s.rawLogs = append(s.rawLogs, fmt.Sprintf("[FAILED] Training aborted: %s", err.Error()))
			s.mu.Unlock()
		} else {
			_ = s.jobRepo.UpdateJobStatus(context.Background(), targetJob.JobID, domain.StatusCompleted, "")
			s.mu.Lock()
			s.rawLogs = append(s.rawLogs, fmt.Sprintf("[SUCCESS] Training run %s completed. Weights best.pt saved.", targetJob.JobID))
			s.mu.Unlock()
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

var _ ports.TrainingUseCase = (*TrainingService)(nil)

