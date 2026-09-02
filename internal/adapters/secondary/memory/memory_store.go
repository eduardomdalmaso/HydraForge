package memory

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	benchmarks  map[string]*domain.BenchmarkJob
}

// NewMemoryStore initializes a new empty memory store and discovers on-disk datasets.
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		jobs:        make(map[string]*domain.TrainingJob),
		datasets:    make(map[string]*domain.Dataset),
		checkpoints: make(map[string]*domain.ModelCheckpoint),
		benchmarks:  make(map[string]*domain.BenchmarkJob),
	}

	// Auto-discover existing datasets in /home/hades/datasets
	datasetsDir := "/home/hades/datasets"
	if entries, err := os.ReadDir(datasetsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				id := entry.Name()
				yamlPath := filepath.Join(datasetsDir, id, "data.yaml")
				if _, err := os.Stat(yamlPath); err == nil {
					// Count images
					trainCount := countDirFiles(filepath.Join(datasetsDir, id, "train", "images"))
					valCount := countDirFiles(filepath.Join(datasetsDir, id, "valid", "images"))
					if valCount == 0 {
						valCount = countDirFiles(filepath.Join(datasetsDir, id, "val", "images"))
					}

					classes := parseYamlClasses(yamlPath)
					name := strings.ReplaceAll(id, "_", " ")
					name = strings.ToUpper(name[:1]) + name[1:]

					store.datasets[id] = &domain.Dataset{
						DatasetID:   id,
						Name:        name,
						Task:        domain.TaskDetect,
						YAMLPath:    yamlPath,
						Classes:     classes,
						NumClasses:  len(classes),
						TrainImages: trainCount,
						ValImages:   valCount,
						CreatedAt:   time.Now(),
					}
				}
			}
		}
	}

	return store
}

func countDirFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count
}

func parseYamlClasses(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{"object"}
	}
	content := string(data)
	lines := strings.Split(content, "\n")
	classes := make([]string, 0)
	inNames := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "names:") {
			if strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]") {
				// Inline array: names: ['car', 'truck']
				start := strings.Index(trimmed, "[")
				end := strings.Index(trimmed, "]")
				if start != -1 && end != -1 && end > start {
					items := strings.Split(trimmed[start+1:end], ",")
					for _, item := range items {
						clean := strings.Trim(strings.TrimSpace(item), "'\"")
						if clean != "" {
							classes = append(classes, clean)
						}
					}
					return classes
				}
			}
			inNames = true
			continue
		}
		if inNames {
			if strings.HasPrefix(trimmed, "-") {
				cls := strings.Trim(strings.TrimPrefix(trimmed, "-"), " '\"")
				if cls != "" {
					classes = append(classes, cls)
				}
			} else if strings.Contains(trimmed, ":") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					cls := strings.Trim(parts[1], " '\"")
					if cls != "" {
						classes = append(classes, cls)
					}
				}
			} else if trimmed == "" || (!strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")) {
				inNames = false
			}
		}
	}
	if len(classes) == 0 {
		classes = []string{"object"}
	}
	return classes
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

// ListDatasets returns all datasets sorted deterministically by ID.
func (s *MemoryStore) ListDatasets(ctx context.Context) ([]*domain.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*domain.Dataset, 0, len(s.datasets))
	for _, d := range s.datasets {
		res = append(res, d)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].DatasetID < res[j].DatasetID
	})
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

// Ensure interface compliances
var (
	_ ports.JobRepository        = (*MemoryStore)(nil)
	_ ports.DatasetRepository    = (*MemoryStore)(nil)
	_ ports.CheckpointRepository = (*MemoryStore)(nil)
	_ ports.BenchmarkRepository  = (*MemoryStore)(nil)
)

