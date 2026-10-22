package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"hydraforge/internal/domain"
)

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
