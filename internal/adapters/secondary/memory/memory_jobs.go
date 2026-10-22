package memory

import (
	"context"
	"sort"
	"time"

	"hydraforge/internal/domain"
)

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

// ListJobs returns jobs sorted by creation time (newest first).
func (s *MemoryStore) ListJobs(ctx context.Context, status string) ([]*domain.TrainingJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*domain.TrainingJob, 0, len(s.jobs))
	for _, j := range s.jobs {
		if status == "" || string(j.Status) == status {
			res = append(res, j)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})
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

// SaveBenchmark stores a benchmark job.
func (s *MemoryStore) SaveBenchmark(ctx context.Context, job *domain.BenchmarkJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.benchmarks[job.JobID] = job
	return nil
}

// GetBenchmark fetches a benchmark job by ID.
func (s *MemoryStore) GetBenchmark(ctx context.Context, jobID string) (*domain.BenchmarkJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.benchmarks[jobID]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	return b, nil
}

// ListBenchmarks returns all benchmark runs, optionally filtered by status.
func (s *MemoryStore) ListBenchmarks(ctx context.Context, status string) ([]*domain.BenchmarkJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*domain.BenchmarkJob, 0, len(s.benchmarks))
	for _, b := range s.benchmarks {
		if status == "" || string(b.Status) == status {
			res = append(res, b)
		}
	}
	return res, nil
}

// UpdateBenchmarkStatus updates the execution status and error of a benchmark job.
func (s *MemoryStore) UpdateBenchmarkStatus(ctx context.Context, jobID string, status domain.BenchmarkStatus, errMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.benchmarks[jobID]
	if !ok {
		return domain.ErrJobNotFound
	}
	b.Status = status
	b.ErrorMessage = errMsg
	b.UpdatedAt = time.Now().UTC()
	return nil
}

// AddBenchmarkResult appends a format result to a benchmark job.
func (s *MemoryStore) AddBenchmarkResult(ctx context.Context, jobID string, result domain.FormatBenchmarkResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.benchmarks[jobID]
	if !ok {
		return domain.ErrJobNotFound
	}
	b.Results = append(b.Results, result)
	b.UpdatedAt = time.Now().UTC()
	return nil
}

// DeleteBenchmark removes a benchmark job.
func (s *MemoryStore) DeleteBenchmark(ctx context.Context, jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.benchmarks, jobID)
	return nil
}
