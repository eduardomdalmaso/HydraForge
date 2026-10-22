package memory

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"hydraforge/internal/config"
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

	datasetsDirs := config.GetDatasetsSearchDirs()
	for _, datasetsDir := range datasetsDirs {
		if entries, err := os.ReadDir(datasetsDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					id := entry.Name()
					yamlPath := filepath.Join(datasetsDir, id, "data.yaml")
					if _, err := os.Stat(yamlPath); err == nil {
						if _, exists := store.datasets[id]; exists {
							continue
						}
						trainCount := countDirFiles(filepath.Join(datasetsDir, id, "images", "train"))
						if trainCount == 0 {
							trainCount = countDirFiles(filepath.Join(datasetsDir, id, "train", "images"))
						}
						valCount := countDirFiles(filepath.Join(datasetsDir, id, "images", "val"))
						if valCount == 0 {
							valCount = countDirFiles(filepath.Join(datasetsDir, id, "val", "images"))
						}
						if valCount == 0 {
							valCount = countDirFiles(filepath.Join(datasetsDir, id, "images", "valid"))
						}
						if valCount == 0 {
							valCount = countDirFiles(filepath.Join(datasetsDir, id, "valid", "images"))
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
	lines := strings.Split(string(data), "\n")
	classes := make([]string, 0)
	inNames := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "names:") {
			if strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]") {
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

// DeleteDataset removes a dataset from memory.
func (s *MemoryStore) DeleteDataset(ctx context.Context, datasetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.datasets, datasetID)
	return nil
}

var (
	_ ports.JobRepository        = (*MemoryStore)(nil)
	_ ports.DatasetRepository    = (*MemoryStore)(nil)
	_ ports.CheckpointRepository = (*MemoryStore)(nil)
	_ ports.BenchmarkRepository  = (*MemoryStore)(nil)
)
