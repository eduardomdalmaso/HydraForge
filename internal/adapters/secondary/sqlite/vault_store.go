package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"strings"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// VaultDatasetRepository implements ports.IVaultDatasetRepository using SQLite WAL.
type VaultDatasetRepository struct {
	db *sql.DB
}

// VaultFrameRepository implements ports.IFrameRepository using SQLite WAL.
type VaultFrameRepository struct {
	db *sql.DB
}

// VaultStore bundles Vault dataset and frame repositories.
type VaultStore struct {
	DatasetRepo *VaultDatasetRepository
	FrameRepo   *VaultFrameRepository
}

// NewVaultStore creates and initializes Vault dataset and frame repositories.
func NewVaultStore(db *sql.DB) (*VaultStore, error) {
	if err := initVaultSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize vault schema: %w", err)
	}
	return &VaultStore{
		DatasetRepo: &VaultDatasetRepository{db: db},
		FrameRepo:   &VaultFrameRepository{db: db},
	}, nil
}

var _ ports.IVaultDatasetRepository = (*VaultDatasetRepository)(nil)
var _ ports.IFrameRepository = (*VaultFrameRepository)(nil)

func initVaultSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS vault_datasets (
		dataset_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		task TEXT NOT NULL DEFAULT 'detect',
		classes_json TEXT NOT NULL DEFAULT '[]',
		train_images INTEGER NOT NULL DEFAULT 0,
		val_images INTEGER NOT NULL DEFAULT 0,
		test_images INTEGER NOT NULL DEFAULT 0,
		total_frames INTEGER NOT NULL DEFAULT 0,
		size_bytes INTEGER NOT NULL DEFAULT 0,
		yaml_path TEXT NOT NULL DEFAULT '',
		checksum_sha TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS frames (
		frame_id TEXT PRIMARY KEY,
		dataset_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_name TEXT NOT NULL,
		phash INTEGER NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		size_bytes INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'RAW',
		split TEXT NOT NULL DEFAULT 'none',
		bboxes_json TEXT NOT NULL DEFAULT '[]',
		polygons_json TEXT NOT NULL DEFAULT '[]',
		source_event TEXT NOT NULL DEFAULT '',
		captured_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_frames_dataset ON frames(dataset_id);
	CREATE INDEX IF NOT EXISTS idx_frames_status ON frames(status);
	CREATE INDEX IF NOT EXISTS idx_frames_split ON frames(split);
	CREATE INDEX IF NOT EXISTS idx_frames_camera ON frames(camera_id);
	CREATE INDEX IF NOT EXISTS idx_frames_phash ON frames(phash);
	`
	_, err := db.Exec(schema)
	return err
}

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

// Save creates or updates a frame entity.
func (r *VaultFrameRepository) Save(ctx context.Context, f *domain.Frame) error {
	if err := f.Validate(); err != nil {
		return err
	}

	bboxesJSON, err := json.Marshal(f.BBoxes)
	if err != nil {
		return fmt.Errorf("failed to marshal bboxes: %w", err)
	}

	polygonsJSON, err := json.Marshal(f.Polygons)
	if err != nil {
		return fmt.Errorf("failed to marshal polygons: %w", err)
	}

	now := time.Now().UTC()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	f.UpdatedAt = now

	query := `
	INSERT INTO frames (
		frame_id, dataset_id, camera_id, file_path, file_name,
		phash, width, height, size_bytes, status, split,
		bboxes_json, polygons_json, source_event, captured_at, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(frame_id) DO UPDATE SET
		dataset_id = excluded.dataset_id,
		camera_id = excluded.camera_id,
		file_path = excluded.file_path,
		file_name = excluded.file_name,
		phash = excluded.phash,
		width = excluded.width,
		height = excluded.height,
		size_bytes = excluded.size_bytes,
		status = excluded.status,
		split = excluded.split,
		bboxes_json = excluded.bboxes_json,
		polygons_json = excluded.polygons_json,
		source_event = excluded.source_event,
		captured_at = excluded.captured_at,
		updated_at = excluded.updated_at;
	`

	_, err = r.db.ExecContext(ctx, query,
		f.FrameID, f.DatasetID, f.CameraID, f.FilePath, f.FileName,
		int64(f.PHash), f.Width, f.Height, f.SizeBytes, string(f.Status), string(f.Split),
		string(bboxesJSON), string(polygonsJSON), f.SourceEvent, f.CapturedAt, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save frame: %w", err)
	}
	return nil
}

// GetByID retrieves a frame by ID.
func (r *VaultFrameRepository) GetByID(ctx context.Context, frameID string) (*domain.Frame, error) {
	query := `
	SELECT frame_id, dataset_id, camera_id, file_path, file_name,
	       phash, width, height, size_bytes, status, split,
	       bboxes_json, polygons_json, source_event, captured_at, created_at, updated_at
	FROM frames WHERE frame_id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, frameID)

	var f domain.Frame
	var phashInt int64
	var statusStr, splitStr, bboxesJSON, polygonsJSON string

	err := row.Scan(
		&f.FrameID, &f.DatasetID, &f.CameraID, &f.FilePath, &f.FileName,
		&phashInt, &f.Width, &f.Height, &f.SizeBytes, &statusStr, &splitStr,
		&bboxesJSON, &polygonsJSON, &f.SourceEvent, &f.CapturedAt, &f.CreatedAt, &f.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrFrameNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query frame: %w", err)
	}

	f.PHash = uint64(phashInt)
	f.Status = domain.FrameStatus(statusStr)
	f.Split = domain.SplitType(splitStr)
	_ = json.Unmarshal([]byte(bboxesJSON), &f.BBoxes)
	_ = json.Unmarshal([]byte(polygonsJSON), &f.Polygons)

	return &f, nil
}

// List lists frames matching filter criteria.
func (r *VaultFrameRepository) List(ctx context.Context, filter ports.FrameFilter) ([]*domain.Frame, error) {
	var conditions []string
	var args []any

	if filter.DatasetID != "" {
		conditions = append(conditions, "dataset_id = ?")
		args = append(args, filter.DatasetID)
	}
	if filter.CameraID != "" {
		conditions = append(conditions, "camera_id = ?")
		args = append(args, filter.CameraID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, string(filter.Status))
	}
	if filter.Split != "" {
		conditions = append(conditions, "split = ?")
		args = append(args, string(filter.Split))
	}

	query := `
	SELECT frame_id, dataset_id, camera_id, file_path, file_name,
	       phash, width, height, size_bytes, status, split,
	       bboxes_json, polygons_json, source_event, captured_at, created_at, updated_at
	FROM frames
	`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY captured_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", filter.Offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list frames: %w", err)
	}
	defer rows.Close()

	var result []*domain.Frame
	for rows.Next() {
		var f domain.Frame
		var phashInt int64
		var statusStr, splitStr, bboxesJSON, polygonsJSON string

		err := rows.Scan(
			&f.FrameID, &f.DatasetID, &f.CameraID, &f.FilePath, &f.FileName,
			&phashInt, &f.Width, &f.Height, &f.SizeBytes, &statusStr, &splitStr,
			&bboxesJSON, &polygonsJSON, &f.SourceEvent, &f.CapturedAt, &f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan frame: %w", err)
		}
		f.PHash = uint64(phashInt)
		f.Status = domain.FrameStatus(statusStr)
		f.Split = domain.SplitType(splitStr)
		_ = json.Unmarshal([]byte(bboxesJSON), &f.BBoxes)
		_ = json.Unmarshal([]byte(polygonsJSON), &f.Polygons)
		result = append(result, &f)
	}

	return result, nil
}

// Count counts frames matching filter criteria.
func (r *VaultFrameRepository) Count(ctx context.Context, filter ports.FrameFilter) (int, error) {
	var conditions []string
	var args []any

	if filter.DatasetID != "" {
		conditions = append(conditions, "dataset_id = ?")
		args = append(args, filter.DatasetID)
	}
	if filter.CameraID != "" {
		conditions = append(conditions, "camera_id = ?")
		args = append(args, filter.CameraID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, string(filter.Status))
	}
	if filter.Split != "" {
		conditions = append(conditions, "split = ?")
		args = append(args, string(filter.Split))
	}

	query := `SELECT COUNT(*) FROM frames`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count frames: %w", err)
	}
	return count, nil
}

// FindByPHash searches for an existing frame whose pHash is within maxHammingDistance.
func (r *VaultFrameRepository) FindByPHash(ctx context.Context, phash uint64, maxHammingDistance int) (*domain.Frame, error) {
	query := `
	SELECT frame_id, dataset_id, camera_id, file_path, file_name,
	       phash, width, height, size_bytes, status, split,
	       bboxes_json, polygons_json, source_event, captured_at, created_at, updated_at
	FROM frames ORDER BY captured_at DESC LIMIT 500;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query frames for phash: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f domain.Frame
		var phashInt int64
		var statusStr, splitStr, bboxesJSON, polygonsJSON string

		err := rows.Scan(
			&f.FrameID, &f.DatasetID, &f.CameraID, &f.FilePath, &f.FileName,
			&phashInt, &f.Width, &f.Height, &f.SizeBytes, &statusStr, &splitStr,
			&bboxesJSON, &polygonsJSON, &f.SourceEvent, &f.CapturedAt, &f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan candidate frame: %w", err)
		}
		candidateHash := uint64(phashInt)
		dist := bits.OnesCount64(phash ^ candidateHash)
		if dist <= maxHammingDistance {
			f.PHash = candidateHash
			f.Status = domain.FrameStatus(statusStr)
			f.Split = domain.SplitType(splitStr)
			_ = json.Unmarshal([]byte(bboxesJSON), &f.BBoxes)
			_ = json.Unmarshal([]byte(polygonsJSON), &f.Polygons)
			return &f, nil
		}
	}

	return nil, domain.ErrFrameNotFound
}

// UpdateStatus updates the lifecycle status and split of a frame.
func (r *VaultFrameRepository) UpdateStatus(ctx context.Context, frameID string, status domain.FrameStatus, split domain.SplitType) error {
	query := `UPDATE frames SET status = ?, split = ?, updated_at = ? WHERE frame_id = ?;`
	res, err := r.db.ExecContext(ctx, query, string(status), string(split), time.Now().UTC(), frameID)
	if err != nil {
		return fmt.Errorf("failed to update frame status: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrFrameNotFound
	}
	return nil
}

// Delete removes a frame by ID.
func (r *VaultFrameRepository) Delete(ctx context.Context, frameID string) error {
	query := `DELETE FROM frames WHERE frame_id = ?;`
	res, err := r.db.ExecContext(ctx, query, frameID)
	if err != nil {
		return fmt.Errorf("failed to delete frame: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrFrameNotFound
	}
	return nil
}
