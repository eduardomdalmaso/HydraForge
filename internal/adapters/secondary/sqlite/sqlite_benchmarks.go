package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"hydraforge/internal/domain"
)

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
	if err := rows.Err(); err != nil {
		return nil
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
