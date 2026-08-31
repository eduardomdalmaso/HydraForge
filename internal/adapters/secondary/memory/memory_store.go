package memory

import (
	"context"
	"sync"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// MemoryStore provides thread-safe in-memory repositories.
type MemoryStore struct {
	mu          sync.RWMutex
	jobs        map[string]*domain.TrainingJob
	datasets    map[string]*domain.Dataset
	checkpoints map[string]*domain.ModelCheckpoint
}

// NewMemoryStore initializes a new populated memory store.
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		jobs:        make(map[string]*domain.TrainingJob),
		datasets:    make(map[string]*domain.Dataset),
		checkpoints: make(map[string]*domain.ModelCheckpoint),
	}

	// Seed with sample dataset
	store.datasets["coco8"] = &domain.Dataset{
		DatasetID:   "coco8",
		Name:        "COCO8 Tiny Sample",
		Task:        domain.TaskDetect,
		YAMLPath:    "datasets/coco8/data.yaml",
		Classes:     []string{"person", "bicycle", "car", "motorcycle", "airplane", "bus", "train", "truck"},
		NumClasses:  8,
		TrainImages: 4,
		ValImages:   4,
		SizeBytes:   2048576,
		CreatedAt:   time.Now(),
	}

	return store
}

// SaveJob stores a job.
func (s *MemoryStore) SaveJob(ctx context.Context, job *domain.TrainingJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.JobID] = job
	return nil
}

// GetJob fetches a job.
func (s *MemoryStore) GetJob(ctx context.Context, jobID string) (*domain.TrainingJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[jobID]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	return j, nil
}

// ListJobs returns jobs.
func (s *MemoryStore) ListJobs(ctx context.Context, status string) ([]*domain.TrainingJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*domain.TrainingJob, 0, len(s.jobs))
	for _, j := range s.jobs {
		if status == "" || string(j.Status) == status {
			res = append(res, j)
		}
	}
	return res, nil
}

// UpdateJobStatus modifies job status and error.
func (s *MemoryStore) UpdateJobStatus(ctx context.Context, jobID string, status domain.JobStatus, errMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[jobID]
	if !ok {
		return domain.ErrJobNotFound
	}
	j.Status = status
	j.ErrorMessage = errMsg
	j.UpdatedAt = time.Now()
	return nil
}

// DeleteJob removes a job.
func (s *MemoryStore) DeleteJob(ctx context.Context, jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, jobID)
	return nil
}

// SaveDataset stores a dataset.
func (s *MemoryStore) SaveDataset(ctx context.Context, dataset *domain.Dataset) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.datasets[dataset.DatasetID] = dataset
	return nil
}

// GetDataset fetches a dataset.
func (s *MemoryStore) GetDataset(ctx context.Context, datasetID string) (*domain.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.datasets[datasetID]
	if !ok {
		return nil, domain.ErrDatasetNotFound
	}
	return d, nil
}

// ListDatasets returns all datasets.
func (s *MemoryStore) ListDatasets(ctx context.Context) ([]*domain.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*domain.Dataset, 0, len(s.datasets))
	for _, d := range s.datasets {
		res = append(res, d)
	}
	return res, nil
}

// SaveCheckpoint stores a checkpoint.
func (s *MemoryStore) SaveCheckpoint(ctx context.Context, checkpoint *domain.ModelCheckpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkpoints[checkpoint.CheckpointID] = checkpoint
	return nil
}

// GetCheckpoint retrieves a checkpoint.
func (s *MemoryStore) GetCheckpoint(ctx context.Context, checkpointID string) (*domain.ModelCheckpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.checkpoints[checkpointID]
	if !ok {
		return nil, domain.ErrCheckpointNotFound
	}
	return c, nil
}

// ListCheckpoints returns checkpoints for a job.
func (s *MemoryStore) ListCheckpoints(ctx context.Context, jobID string) ([]*domain.ModelCheckpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*domain.ModelCheckpoint, 0, len(s.checkpoints))
	for _, c := range s.checkpoints {
		if jobID == "" || c.JobID == jobID {
			res = append(res, c)
		}
	}
	return res, nil
}

// Ensure interface compliances
var (
	_ ports.JobRepository        = (*MemoryStore)(nil)
	_ ports.DatasetRepository    = (*MemoryStore)(nil)
	_ ports.CheckpointRepository = (*MemoryStore)(nil)
)
