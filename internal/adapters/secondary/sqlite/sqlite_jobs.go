package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"hydraforge/internal/domain"
)

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
