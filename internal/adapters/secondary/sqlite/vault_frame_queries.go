package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"math/bits"
	"strings"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

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
