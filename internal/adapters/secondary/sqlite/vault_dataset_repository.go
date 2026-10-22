package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// VaultDatasetRepository implements ports.IVaultDatasetRepository using SQLite WAL.
type VaultDatasetRepository struct {
	db *sql.DB
}

var _ ports.IVaultDatasetRepository = (*VaultDatasetRepository)(nil)

// Save creates or updates a curated dataset.
func (r *VaultDatasetRepository) Save(ctx context.Context, d *domain.VaultDataset) error {
	if err := d.Validate(); err != nil {
		return err
	}

	classesJSON, err := json.Marshal(d.Classes)
	if err != nil {
		return fmt.Errorf("failed to marshal classes: %w", err)
	}

	now := time.Now().UTC()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	query := `
	INSERT INTO vault_datasets (
		dataset_id, name, description, task, classes_json,
		train_images, val_images, test_images, total_frames, size_bytes,
		yaml_path, checksum_sha, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(dataset_id) DO UPDATE SET
		name = excluded.name,
		description = excluded.description,
		task = excluded.task,
		classes_json = excluded.classes_json,
		train_images = excluded.train_images,
		val_images = excluded.val_images,
		test_images = excluded.test_images,
		total_frames = excluded.total_frames,
		size_bytes = excluded.size_bytes,
		yaml_path = excluded.yaml_path,
		checksum_sha = excluded.checksum_sha,
		updated_at = excluded.updated_at;
	`

	_, err = r.db.ExecContext(ctx, query,
		d.DatasetID, d.Name, d.Description, string(d.Task), string(classesJSON),
		d.TrainImages, d.ValImages, d.TestImages, d.TotalFrames, d.SizeBytes,
		d.YAMLPath, d.ChecksumSHA, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save vault dataset: %w", err)
	}
	return nil
}

// GetByID retrieves a curated dataset by ID.
func (r *VaultDatasetRepository) GetByID(ctx context.Context, datasetID string) (*domain.VaultDataset, error) {
	query := `
	SELECT dataset_id, name, description, task, classes_json,
	       train_images, val_images, test_images, total_frames, size_bytes,
	       yaml_path, checksum_sha, created_at, updated_at
	FROM vault_datasets WHERE dataset_id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, datasetID)

	var d domain.VaultDataset
	var taskStr, classesJSON string

	err := row.Scan(
		&d.DatasetID, &d.Name, &d.Description, &taskStr, &classesJSON,
		&d.TrainImages, &d.ValImages, &d.TestImages, &d.TotalFrames, &d.SizeBytes,
		&d.YAMLPath, &d.ChecksumSHA, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrDatasetNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query vault dataset: %w", err)
	}

	d.Task = domain.TaskType(taskStr)
	if err := json.Unmarshal([]byte(classesJSON), &d.Classes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal classes: %w", err)
	}

	return &d, nil
}

// List retrieves all curated datasets.
func (r *VaultDatasetRepository) List(ctx context.Context) ([]*domain.VaultDataset, error) {
	query := `
	SELECT dataset_id, name, description, task, classes_json,
	       train_images, val_images, test_images, total_frames, size_bytes,
	       yaml_path, checksum_sha, created_at, updated_at
	FROM vault_datasets ORDER BY updated_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list vault datasets: %w", err)
	}
	defer rows.Close()

	var result []*domain.VaultDataset
	for rows.Next() {
		var d domain.VaultDataset
		var taskStr, classesJSON string
		err := rows.Scan(
			&d.DatasetID, &d.Name, &d.Description, &taskStr, &classesJSON,
			&d.TrainImages, &d.ValImages, &d.TestImages, &d.TotalFrames, &d.SizeBytes,
			&d.YAMLPath, &d.ChecksumSHA, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vault dataset: %w", err)
		}
		d.Task = domain.TaskType(taskStr)
		_ = json.Unmarshal([]byte(classesJSON), &d.Classes)
		result = append(result, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating vault datasets: %w", err)
	}
	return result, nil
}

// Delete removes a curated dataset by ID.
func (r *VaultDatasetRepository) Delete(ctx context.Context, datasetID string) error {
	query := `DELETE FROM vault_datasets WHERE dataset_id = ?;`
	res, err := r.db.ExecContext(ctx, query, datasetID)
	if err != nil {
		return fmt.Errorf("failed to delete vault dataset: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrDatasetNotFound
	}
	return nil
}
