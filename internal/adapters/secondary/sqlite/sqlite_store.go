package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// SQLiteStore provides persistent storage for HydraForge entities with primary keys.
type SQLiteStore struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteStore initializes a SQLite database connection and runs auto-migrations.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}

	store.autoDiscoverDatasets("/home/hades/Documents/HydraForge/datasets")
	store.autoDiscoverDatasets("datasets")
	if _, err := os.Stat("/home/hades/datasets"); err == nil {
		store.autoDiscoverDatasets("/home/hades/datasets")
	}
	return store, nil
}

func (s *SQLiteStore) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS registered_datasets (
		dataset_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		task TEXT NOT NULL,
		yaml_path TEXT NOT NULL,
		classes_json TEXT NOT NULL,
		num_classes INTEGER NOT NULL,
		train_images INTEGER NOT NULL,
		val_images INTEGER NOT NULL,
		test_images INTEGER DEFAULT 0,
		size_bytes INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS training_jobs (
		job_id TEXT PRIMARY KEY,
		model_architecture TEXT NOT NULL,
		task TEXT NOT NULL,
		dataset_id TEXT NOT NULL,
		dataset_path TEXT,
		hyperparameters_json TEXT NOT NULL,
		status TEXT NOT NULL,
		current_epoch INTEGER DEFAULT 0,
		total_epochs INTEGER DEFAULT 50,
		best_map50 REAL DEFAULT 0.0,
		total_energy_kwh REAL DEFAULT 0.0,
		avg_power_watts REAL DEFAULT 0.0,
		peak_vram_mb REAL DEFAULT 0.0,
		avg_fps REAL DEFAULT 0.0,
		duration_sec REAL DEFAULT 0.0,
		output_weights TEXT,
		error_message TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY(dataset_id) REFERENCES registered_datasets(dataset_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS epoch_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		epoch INTEGER NOT NULL,
		total_epochs INTEGER NOT NULL,
		box_loss REAL,
		cls_loss REAL,
		dfl_loss REAL,
		val_box_loss REAL,
		val_cls_loss REAL,
		map50 REAL,
		map50_95 REAL,
		precision REAL,
		recall REAL,
		learning_rate REAL,
		gpu_vram_mb REAL,
		power_watts REAL,
		temp_celsius REAL,
		gpu_util_pct REAL,
		fps REAL,
		epoch_duration_sec REAL,
		created_at DATETIME NOT NULL,
		UNIQUE(job_id, epoch),
		FOREIGN KEY(job_id) REFERENCES training_jobs(job_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS model_checkpoints (
		checkpoint_id TEXT PRIMARY KEY,
		job_id TEXT NOT NULL,
		architecture TEXT NOT NULL,
		task TEXT NOT NULL,
		epoch INTEGER NOT NULL,
		map50 REAL,
		map50_95 REAL,
		weights_path TEXT NOT NULL,
		export_format TEXT,
		export_path TEXT,
		size_bytes INTEGER NOT NULL,
		precision TEXT,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(job_id) REFERENCES training_jobs(job_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS benchmark_jobs (
		job_id TEXT PRIMARY KEY,
		model TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '',
		model_architecture TEXT DEFAULT '',
		weights_path TEXT DEFAULT '',
		imgsz INTEGER NOT NULL,
		quantize INTEGER NOT NULL,
		device TEXT NOT NULL,
		target_formats_json TEXT NOT NULL,
		status TEXT NOT NULL,
		error_message TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS benchmark_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		format TEXT NOT NULL,
		status TEXT NOT NULL,
		latency_ms REAL NOT NULL,
		fps REAL NOT NULL,
		size_mb REAL NOT NULL,
		map50_95 REAL NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(job_id) REFERENCES benchmark_jobs(job_id) ON DELETE CASCADE
	);
	`
	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN current_batch INTEGER DEFAULT 0;`)
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN total_batches INTEGER DEFAULT 0;`)
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN best_map50_95 REAL DEFAULT 0.0;`)
	_, _ = s.db.Exec(`ALTER TABLE benchmark_jobs ADD COLUMN model TEXT DEFAULT '';`)
	_, _ = s.db.Exec(`ALTER TABLE benchmark_jobs ADD COLUMN data TEXT DEFAULT '';`)
	return nil
}

func (s *SQLiteStore) RescanDatasets() {
	s.autoDiscoverDatasets("/home/hades/Documents/HydraForge/datasets")
	s.autoDiscoverDatasets("datasets")
	if _, err := os.Stat("/home/hades/datasets"); err == nil {
		s.autoDiscoverDatasets("/home/hades/datasets")
	}
}

func countImages(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".bmp" {
				count++
			}
		}
	}
	return count
}

func sanitizeYamlPaths(yamlPath, datasetDir string) {
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		return
	}
	s := string(content)
	modified := false
	if strings.Contains(s, "../train/images") {
		s = strings.ReplaceAll(s, "../train/images", "train/images")
		modified = true
	}
	if strings.Contains(s, "../valid/images") {
		s = strings.ReplaceAll(s, "../valid/images", "valid/images")
		modified = true
	}
	if strings.Contains(s, "../test/images") {
		s = strings.ReplaceAll(s, "../test/images", "test/images")
		modified = true
	}
	if !strings.Contains(s, "path:") {
		s = fmt.Sprintf("path: %s\n%s", datasetDir, s)
		modified = true
	}
	if modified {
		_ = os.WriteFile(yamlPath, []byte(s), 0644)
	}
}

func (s *SQLiteStore) autoDiscoverDatasets(datasetsDir string) {
	entries, err := os.ReadDir(datasetsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		id := entry.Name()
		fullPath := filepath.Join(datasetsDir, id)
		fi, err := os.Stat(fullPath)
		if err != nil || !fi.IsDir() {
			continue
		}

		yamlPath := filepath.Join(fullPath, "data.yaml")
		if _, err := os.Stat(yamlPath); err != nil {
			imgCount := countImages(fullPath)
			lblDir := filepath.Join(fullPath, "labels")
			if imgCount > 0 {
				if _, err := os.Stat(lblDir); err == nil {
					clsName := strings.ReplaceAll(id, "_", " ")
					clsName = strings.ReplaceAll(clsName, "-", " ")
					content := fmt.Sprintf("path: %s\ntrain: .\nval: .\nnc: 1\nnames:\n  0: %s\n", fullPath, clsName)
					_ = os.WriteFile(yamlPath, []byte(content), 0644)
				}
			}
		}

		if _, err := os.Stat(yamlPath); err == nil {
			sanitizeYamlPaths(yamlPath, fullPath)
			trainCount := countDir(filepath.Join(fullPath, "train", "images"))
			valCount := countDir(filepath.Join(fullPath, "valid", "images"))
			if valCount == 0 {
				valCount = countDir(filepath.Join(fullPath, "val", "images"))
			}
			testCount := countDir(filepath.Join(fullPath, "test", "images"))
			if trainCount == 0 && valCount == 0 {
				trainCount = countImages(fullPath)
			}
			name := strings.ReplaceAll(id, "_", " ")
			name = strings.ReplaceAll(name, "-", " ")
			name = strings.ToUpper(name[:1]) + name[1:]
			classes := parseYamlClasses(yamlPath)
			classesJSON, _ := json.Marshal(classes)

			query := `INSERT INTO registered_datasets 
				(dataset_id, name, task, yaml_path, classes_json, num_classes, train_images, val_images, test_images, size_bytes, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(dataset_id) DO UPDATE SET 
					name=excluded.name, yaml_path=excluded.yaml_path,
					classes_json=excluded.classes_json, num_classes=excluded.num_classes,
					train_images=excluded.train_images, val_images=excluded.val_images, test_images=excluded.test_images`
			_, _ = s.db.Exec(query, id, name, "detect", yamlPath, string(classesJSON), len(classes), trainCount, valCount, testCount, 0, time.Now())
		}
	}
}

// SaveJob persists a training job into SQLite.
func (s *SQLiteStore) SaveJob(ctx context.Context, job *domain.TrainingJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	paramsJSON, err := json.Marshal(job.Hyperparameters)
	if err != nil {
		return err
	}

	query := `INSERT INTO training_jobs (
		job_id, model_architecture, task, dataset_id, dataset_path,
		hyperparameters_json, status, current_epoch, total_epochs,
		current_batch, total_batches,
		best_map50, best_map50_95, total_energy_kwh, avg_power_watts, peak_vram_mb,
		avg_fps, duration_sec, output_weights, error_message, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(job_id) DO UPDATE SET
		status=excluded.status,
		current_epoch=excluded.current_epoch,
		current_batch=excluded.current_batch,
		total_batches=excluded.total_batches,
		best_map50=excluded.best_map50,
		best_map50_95=excluded.best_map50_95,
		total_energy_kwh=excluded.total_energy_kwh,
		avg_power_watts=excluded.avg_power_watts,
		peak_vram_mb=excluded.peak_vram_mb,
		avg_fps=excluded.avg_fps,
		duration_sec=excluded.duration_sec,
		output_weights=excluded.output_weights,
		error_message=excluded.error_message,
		updated_at=excluded.updated_at`

	now := time.Now()
	if job.CreatedAt.IsZero() {
		job.CreatedAt = now
	}
	job.UpdatedAt = now

	_, err = s.db.ExecContext(ctx, query,
		job.JobID, job.ModelArchitecture, string(job.Task), job.DatasetID, job.DatasetPath,
		string(paramsJSON), string(job.Status), job.CurrentEpoch, job.TotalEpochs,
		job.CurrentBatch, job.TotalBatches,
		job.BestMAP50, job.BestMAP50_95, job.TotalEnergyKWh, job.AvgPowerWatts, job.PeakVRAMMB,
		job.AvgFPS, job.DurationSec, job.OutputWeights, job.ErrorMessage, job.CreatedAt, job.UpdatedAt,
	)
	return err
}

// GetJob retrieves a job and its hyperparameters from SQLite.
func (s *SQLiteStore) GetJob(ctx context.Context, jobID string) (*domain.TrainingJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT job_id, model_architecture, task, dataset_id, dataset_path,
		hyperparameters_json, status, current_epoch, total_epochs,
		current_batch, total_batches, best_map50, best_map50_95,
		total_energy_kwh, avg_power_watts, peak_vram_mb, avg_fps, duration_sec,
		output_weights, error_message, created_at, updated_at
		FROM training_jobs WHERE job_id = ?`

	var job domain.TrainingJob
	var taskStr, statusStr, paramsJSON string
	row := s.db.QueryRowContext(ctx, query, jobID)
	err := row.Scan(
		&job.JobID, &job.ModelArchitecture, &taskStr, &job.DatasetID, &job.DatasetPath,
		&paramsJSON, &statusStr, &job.CurrentEpoch, &job.TotalEpochs,
		&job.CurrentBatch, &job.TotalBatches, &job.BestMAP50, &job.BestMAP50_95,
		&job.TotalEnergyKWh, &job.AvgPowerWatts, &job.PeakVRAMMB, &job.AvgFPS, &job.DurationSec,
		&job.OutputWeights, &job.ErrorMessage, &job.CreatedAt, &job.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrJobNotFound
	} else if err != nil {
		return nil, err
	}

	job.Task = domain.TaskType(taskStr)
	job.Status = domain.JobStatus(statusStr)
	_ = json.Unmarshal([]byte(paramsJSON), &job.Hyperparameters)
	return &job, nil
}

// ListJobs retrieves all training jobs ordered by creation date.
func (s *SQLiteStore) ListJobs(ctx context.Context, status string) ([]*domain.TrainingJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var query string
	var rows *sql.Rows
	var err error

	if status != "" {
		query = `SELECT job_id, model_architecture, task, dataset_id, dataset_path,
			hyperparameters_json, status, current_epoch, total_epochs,
			current_batch, total_batches, best_map50, best_map50_95,
			total_energy_kwh, avg_power_watts, peak_vram_mb, avg_fps, duration_sec,
			output_weights, error_message, created_at, updated_at
			FROM training_jobs WHERE status = ? ORDER BY created_at DESC`
		rows, err = s.db.QueryContext(ctx, query, status)
	} else {
		query = `SELECT job_id, model_architecture, task, dataset_id, dataset_path,
			hyperparameters_json, status, current_epoch, total_epochs,
			current_batch, total_batches, best_map50, best_map50_95,
			total_energy_kwh, avg_power_watts, peak_vram_mb, avg_fps, duration_sec,
			output_weights, error_message, created_at, updated_at
			FROM training_jobs ORDER BY created_at DESC`
		rows, err = s.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.TrainingJob
	for rows.Next() {
		var job domain.TrainingJob
		var taskStr, statusStr, paramsJSON string
		if err := rows.Scan(
			&job.JobID, &job.ModelArchitecture, &taskStr, &job.DatasetID, &job.DatasetPath,
			&paramsJSON, &statusStr, &job.CurrentEpoch, &job.TotalEpochs,
			&job.CurrentBatch, &job.TotalBatches, &job.BestMAP50, &job.BestMAP50_95,
			&job.TotalEnergyKWh, &job.AvgPowerWatts, &job.PeakVRAMMB, &job.AvgFPS, &job.DurationSec,
			&job.OutputWeights, &job.ErrorMessage, &job.CreatedAt, &job.UpdatedAt,
		); err != nil {
			return nil, err
		}
		job.Task = domain.TaskType(taskStr)
		job.Status = domain.JobStatus(statusStr)
		_ = json.Unmarshal([]byte(paramsJSON), &job.Hyperparameters)
		jobs = append(jobs, &job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

// UpdateJobStatus updates the status and error message for a training job.
func (s *SQLiteStore) UpdateJobStatus(ctx context.Context, jobID string, status domain.JobStatus, errMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE training_jobs SET status = ?, error_message = ?, updated_at = ? WHERE job_id = ?`
	_, err := s.db.ExecContext(ctx, query, string(status), errMsg, time.Now(), jobID)
	return err
}

// DeleteJob removes a training job from SQLite.
func (s *SQLiteStore) DeleteJob(ctx context.Context, jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM training_jobs WHERE job_id = ?`
	_, err := s.db.ExecContext(ctx, query, jobID)
	return err
}

// SaveDataset inserts or updates a dataset record.
func (s *SQLiteStore) SaveDataset(ctx context.Context, dataset *domain.Dataset) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	classesJSON, _ := json.Marshal(dataset.Classes)
	query := `INSERT INTO registered_datasets (
		dataset_id, name, task, yaml_path, classes_json, num_classes,
		train_images, val_images, test_images, size_bytes, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(dataset_id) DO UPDATE SET
		name=excluded.name,
		yaml_path=excluded.yaml_path,
		classes_json=excluded.classes_json,
		num_classes=excluded.num_classes,
		train_images=excluded.train_images,
		val_images=excluded.val_images,
		test_images=excluded.test_images,
		size_bytes=excluded.size_bytes`

	now := dataset.CreatedAt
	if now.IsZero() {
		now = time.Now()
	}

	_, err := s.db.ExecContext(ctx, query,
		dataset.DatasetID, dataset.Name, string(dataset.Task), dataset.YAMLPath,
		string(classesJSON), dataset.NumClasses, dataset.TrainImages, dataset.ValImages,
		dataset.TestImages, dataset.SizeBytes, now,
	)
	return err
}

// GetDataset retrieves a dataset by its unique primary key ID.
func (s *SQLiteStore) GetDataset(ctx context.Context, datasetID string) (*domain.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT dataset_id, name, task, yaml_path, classes_json, num_classes,
		train_images, val_images, test_images, size_bytes, created_at FROM registered_datasets WHERE dataset_id = ?`

	var ds domain.Dataset
	var taskStr, classesJSON string
	row := s.db.QueryRowContext(ctx, query, datasetID)
	err := row.Scan(
		&ds.DatasetID, &ds.Name, &taskStr, &ds.YAMLPath, &classesJSON,
		&ds.NumClasses, &ds.TrainImages, &ds.ValImages, &ds.TestImages, &ds.SizeBytes, &ds.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrDatasetNotFound
	} else if err != nil {
		return nil, err
	}

	ds.Task = domain.TaskType(taskStr)
	_ = json.Unmarshal([]byte(classesJSON), &ds.Classes)
	return &ds, nil
}

// ListDatasets returns all registered datasets.
func (s *SQLiteStore) ListDatasets(ctx context.Context) ([]*domain.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT dataset_id, name, task, yaml_path, classes_json, num_classes,
		train_images, val_images, test_images, size_bytes, created_at FROM registered_datasets ORDER BY dataset_id ASC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Dataset
	for rows.Next() {
		var ds domain.Dataset
		var taskStr, classesJSON string
		if err := rows.Scan(
			&ds.DatasetID, &ds.Name, &taskStr, &ds.YAMLPath, &classesJSON,
			&ds.NumClasses, &ds.TrainImages, &ds.ValImages, &ds.TestImages, &ds.SizeBytes, &ds.CreatedAt,
		); err != nil {
			return nil, err
		}
		ds.Task = domain.TaskType(taskStr)
		_ = json.Unmarshal([]byte(classesJSON), &ds.Classes)
		list = append(list, &ds)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteDataset removes a dataset record from the database.
func (s *SQLiteStore) DeleteDataset(ctx context.Context, datasetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx, "DELETE FROM registered_datasets WHERE dataset_id = ?", datasetID)
	return err
}

// SaveCheckpoint saves a model checkpoint with primary key ID.
func (s *SQLiteStore) SaveCheckpoint(ctx context.Context, ckpt *domain.ModelCheckpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `INSERT INTO model_checkpoints (
		checkpoint_id, job_id, architecture, task, epoch, map50, map50_95,
		weights_path, export_format, export_path, size_bytes, precision, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(checkpoint_id) DO UPDATE SET
		map50=excluded.map50, map50_95=excluded.map50_95, export_format=excluded.export_format,
		export_path=excluded.export_path, size_bytes=excluded.size_bytes, precision=excluded.precision`

	now := ckpt.CreatedAt
	if now.IsZero() {
		now = time.Now()
	}

	_, err := s.db.ExecContext(ctx, query,
		ckpt.CheckpointID, ckpt.JobID, ckpt.Architecture, string(ckpt.Task),
		ckpt.Epoch, ckpt.MAP50, ckpt.MAP50_95, ckpt.WeightsPath,
		string(ckpt.ExportFormat), ckpt.ExportPath, ckpt.SizeBytes, ckpt.Precision, now,
	)
	return err
}

// GetCheckpoint retrieves a checkpoint by ID.
func (s *SQLiteStore) GetCheckpoint(ctx context.Context, checkpointID string) (*domain.ModelCheckpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT checkpoint_id, job_id, architecture, task, epoch, map50, map50_95,
		weights_path, export_format, export_path, size_bytes, precision, created_at FROM model_checkpoints WHERE checkpoint_id = ?`

	var c domain.ModelCheckpoint
	var taskStr, expFmtStr string
	row := s.db.QueryRowContext(ctx, query, checkpointID)
	err := row.Scan(&c.CheckpointID, &c.JobID, &c.Architecture, &taskStr, &c.Epoch, &c.MAP50, &c.MAP50_95, &c.WeightsPath, &expFmtStr, &c.ExportPath, &c.SizeBytes, &c.Precision, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrJobNotFound
	} else if err != nil {
		return nil, err
	}
	c.Task = domain.TaskType(taskStr)
	c.ExportFormat = domain.ExportFormat(expFmtStr)
	return &c, nil
}

// ListCheckpoints returns all checkpoints for a job.
func (s *SQLiteStore) ListCheckpoints(ctx context.Context, jobID string) ([]*domain.ModelCheckpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT checkpoint_id, job_id, architecture, task, epoch, map50, map50_95,
		weights_path, export_format, export_path, size_bytes, precision, created_at FROM model_checkpoints WHERE job_id = ? ORDER BY epoch DESC`

	rows, err := s.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ModelCheckpoint
	for rows.Next() {
		var c domain.ModelCheckpoint
		var taskStr, expFmtStr string
		if err := rows.Scan(&c.CheckpointID, &c.JobID, &c.Architecture, &taskStr, &c.Epoch, &c.MAP50, &c.MAP50_95, &c.WeightsPath, &expFmtStr, &c.ExportPath, &c.SizeBytes, &c.Precision, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.Task = domain.TaskType(taskStr)
		c.ExportFormat = domain.ExportFormat(expFmtStr)
		list = append(list, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// SaveBenchmark saves a benchmark job.
func (s *SQLiteStore) SaveBenchmark(ctx context.Context, job *domain.BenchmarkJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmtsJSON, _ := json.Marshal(job.TargetFormats)
	query := `INSERT INTO benchmark_jobs (
		job_id, model, data, model_architecture, weights_path, imgsz, quantize, device,
		target_formats_json, status, error_message, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(job_id) DO UPDATE SET status=excluded.status, error_message=excluded.error_message, updated_at=excluded.updated_at`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query,
		job.JobID, job.Model, job.Data, job.Model, job.Data, job.ImageSize, job.Quantize,
		job.Device, string(fmtsJSON), string(job.Status), job.ErrorMessage, now, now,
	)
	return err
}

// GetBenchmark retrieves a benchmark job by ID.
func (s *SQLiteStore) GetBenchmark(ctx context.Context, jobID string) (*domain.BenchmarkJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT job_id, model, data, imgsz, quantize, device,
		target_formats_json, status, error_message, created_at, updated_at FROM benchmark_jobs WHERE job_id = ?`

	var b domain.BenchmarkJob
	var fmtsJSON, statusStr string
	row := s.db.QueryRowContext(ctx, query, jobID)
	err := row.Scan(&b.JobID, &b.Model, &b.Data, &b.ImageSize, &b.Quantize, &b.Device, &fmtsJSON, &statusStr, &b.ErrorMessage, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrJobNotFound
	} else if err != nil {
		return nil, err
	}
	b.Status = domain.BenchmarkStatus(statusStr)
	_ = json.Unmarshal([]byte(fmtsJSON), &b.TargetFormats)
	b.Results = s.loadBenchmarkResults(ctx, jobID)
	return &b, nil
}

// ListBenchmarks returns all benchmark runs.
func (s *SQLiteStore) ListBenchmarks(ctx context.Context, status string) ([]*domain.BenchmarkJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT job_id, model, data, imgsz, quantize, device,
		target_formats_json, status, error_message, created_at, updated_at FROM benchmark_jobs ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.BenchmarkJob
	for rows.Next() {
		var b domain.BenchmarkJob
		var fmtsJSON, statusStr string
		if err := rows.Scan(&b.JobID, &b.Model, &b.Data, &b.ImageSize, &b.Quantize, &b.Device, &fmtsJSON, &statusStr, &b.ErrorMessage, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.Status = domain.BenchmarkStatus(statusStr)
		_ = json.Unmarshal([]byte(fmtsJSON), &b.TargetFormats)
		b.Results = s.loadBenchmarkResults(ctx, b.JobID)
		list = append(list, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateBenchmarkStatus updates the status of a benchmark job.
func (s *SQLiteStore) UpdateBenchmarkStatus(ctx context.Context, jobID string, status domain.BenchmarkStatus, errMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE benchmark_jobs SET status = ?, error_message = ?, updated_at = ? WHERE job_id = ?`
	_, err := s.db.ExecContext(ctx, query, string(status), errMsg, time.Now(), jobID)
	return err
}

// AddBenchmarkResult persists a benchmark result for a format.
func (s *SQLiteStore) AddBenchmarkResult(ctx context.Context, jobID string, result domain.FormatBenchmarkResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `INSERT INTO benchmark_results (job_id, format, status, latency_ms, fps, size_mb, map50_95, export_args, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, jobID, result.Format, string(result.Status), result.InferenceTimeMS, result.FPS, result.SizeMB, result.MAP50_95, result.ExportArgs, result.ErrorMessage, time.Now())
	return err
}

func (s *SQLiteStore) loadBenchmarkResults(ctx context.Context, jobID string) []domain.FormatBenchmarkResult {
	rows, err := s.db.QueryContext(ctx, `SELECT format, status, latency_ms, fps, size_mb, map50_95, export_args, error_message FROM benchmark_results WHERE job_id = ? ORDER BY id ASC`, jobID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []domain.FormatBenchmarkResult
	for rows.Next() {
		var r domain.FormatBenchmarkResult
		var status, exportArgs, errMsg sql.NullString
		if err := rows.Scan(&r.Format, &status, &r.InferenceTimeMS, &r.FPS, &r.SizeMB, &r.MAP50_95, &exportArgs, &errMsg); err == nil {
			r.Status = status.String
			r.ExportArgs = exportArgs.String
			r.ErrorMessage = errMsg.String
			results = append(results, r)
		}
	}
	return results
}

// DeleteBenchmark removes a benchmark job.
func (s *SQLiteStore) DeleteBenchmark(ctx context.Context, jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM benchmark_jobs WHERE job_id = ?`
	_, err := s.db.ExecContext(ctx, query, jobID)
	return err
}

func countDir(dirPath string) int {
	entries, err := os.ReadDir(dirPath)
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

func parseYamlClasses(yamlPath string) []string {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return []string{"object"}
	}
	lines := strings.Split(string(data), "\n")
	var classes []string
	inNames := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "names:") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "names:"))
			if strings.HasPrefix(rest, "[") && strings.HasSuffix(rest, "]") {
				inner := strings.Trim(rest, "[]")
				parts := strings.Split(inner, ",")
				for _, p := range parts {
					c := strings.TrimSpace(p)
					c = strings.Trim(c, "'\"")
					if c != "" {
						classes = append(classes, c)
					}
				}
				if len(classes) > 0 {
					return classes
				}
			}
			inNames = true
			continue
		}
		if inNames {
			if strings.HasPrefix(trimmed, "-") {
				classes = append(classes, strings.TrimSpace(strings.TrimPrefix(trimmed, "-")))
			} else if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "#") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					val := strings.TrimSpace(parts[1])
					val = strings.Trim(val, "'\"")
					if val != "" {
						classes = append(classes, val)
					}
				}
			} else if len(trimmed) > 0 && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
				break
			}
		}
	}
	if len(classes) == 0 {
		return []string{"object"}
	}
	return classes
}

var _ ports.JobRepository = (*SQLiteStore)(nil)
var _ ports.DatasetRepository = (*SQLiteStore)(nil)
var _ ports.CheckpointRepository = (*SQLiteStore)(nil)
var _ ports.BenchmarkRepository = (*SQLiteStore)(nil)
