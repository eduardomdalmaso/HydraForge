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

// VaultFrameRepository implements ports.IFrameRepository using SQLite WAL.
type VaultFrameRepository struct {
	db *sql.DB
}

var _ ports.IFrameRepository = (*VaultFrameRepository)(nil)

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
		status = excluded.status,
		split = excluded.split,
		bboxes_json = excluded.bboxes_json,
		polygons_json = excluded.polygons_json,
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

// GetByID retrieves a frame by its ID.
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
