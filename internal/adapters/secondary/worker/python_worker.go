package worker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"hydraforge/internal/config"
	"hydraforge/internal/domain"
)

// PythonWorker handles execution of YOLO training runs.
type PythonWorker struct {
	mu          sync.Mutex
	activeJobs  map[string]bool
	cancelFuncs map[string]context.CancelFunc
}

// NewPythonWorker creates a new worker adapter.
func NewPythonWorker() *PythonWorker {
	return &PythonWorker{
		activeJobs:  make(map[string]bool),
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

// StartTraining executes the training loop with metric streaming.
func (w *PythonWorker) StartTraining(ctx context.Context, job *domain.TrainingJob, metricCallback func(domain.TrainingMetrics)) error {
	w.mu.Lock()
	ctx, cancel := context.WithCancel(ctx)
	w.activeJobs[job.JobID] = true
	w.cancelFuncs[job.JobID] = cancel
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.activeJobs, job.JobID)
		delete(w.cancelFuncs, job.JobID)
		w.mu.Unlock()
	}()

	pythonBin := config.GetPythonBin()

	yamlPath := job.DatasetPath
	if yamlPath == "" {
		for _, base := range config.GetDatasetsSearchDirs() {
			cand := filepath.Join(base, job.DatasetID, "data.yaml")
			if _, err := os.Stat(cand); err == nil {
				yamlPath = cand
				break
			}
		}
		if yamlPath == "" {
			yamlPath = filepath.Join("datasets", job.DatasetID, "data.yaml")
		}
	}

	scriptPath := "worker_python/train.py"
	modelArg := job.ModelArchitecture
	if !strings.HasSuffix(modelArg, ".pt") && !strings.HasSuffix(modelArg, ".yaml") {
		modelArg += ".pt"
	}
	if _, err := os.Stat(filepath.Join("weights", filepath.Base(modelArg))); err == nil {
		modelArg = filepath.Join("weights", filepath.Base(modelArg))
	}

	args := []string{
		scriptPath,
		"--model", modelArg,
		"--data", yamlPath,
		"--epochs", strconv.Itoa(job.Hyperparameters.Epochs),
		"--batch", strconv.Itoa(job.Hyperparameters.BatchSize),
		"--imgsz", strconv.Itoa(job.Hyperparameters.ImageSize),
		"--optimizer", job.Hyperparameters.Optimizer,
		"--lr0", fmt.Sprintf("%f", job.Hyperparameters.LearningRate),
		"--close-mosaic", strconv.Itoa(job.Hyperparameters.CloseMosaic),
		"--job-id", job.JobID,
	}
	if job.Hyperparameters.Patience > 0 {
		args = append(args, "--patience", strconv.Itoa(job.Hyperparameters.Patience))
	}
	if job.Hyperparameters.FLGamma > 0 {
		args = append(args, "--fl-gamma", fmt.Sprintf("%.2f", job.Hyperparameters.FLGamma))
	}
	if job.Hyperparameters.Freeze > 0 {
		args = append(args, "--freeze", strconv.Itoa(job.Hyperparameters.Freeze))
	}
	if job.Hyperparameters.UseAMP {
		args = append(args, "--amp")
	}

	cmd := exec.CommandContext(ctx, pythonBin, args...)
	cmd.Dir = "."
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "HYDRA_METRIC:"); idx != -1 {
			jsonStr := line[idx+len("HYDRA_METRIC:"):]
			var m domain.TrainingMetrics
			if err := json.Unmarshal([]byte(jsonStr), &m); err == nil {
				if metricCallback != nil {
					metricCallback(m)
				}
			}
		} else if idx := strings.Index(line, "HYDRA_BATCH:"); idx != -1 {
			jsonStr := line[idx+len("HYDRA_BATCH:"):]
			var b struct {
				Epoch          int     `json:"epoch"`
				Batch          int     `json:"batch"`
				TotalBatches   int     `json:"total_batches"`
				BoxLoss        float64 `json:"box_loss"`
				ClsLoss        float64 `json:"cls_loss"`
				PowerWatts     float64 `json:"power_watts"`
				TotalEnergyKWh float64 `json:"total_energy_kwh"`
				FPS            float64 `json:"fps"`
				DurationSec    float64 `json:"duration_sec"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &b); err == nil {
				job.CurrentEpoch = b.Epoch
				job.CurrentBatch = b.Batch
				job.TotalBatches = b.TotalBatches
				job.TotalEnergyKWh = b.TotalEnergyKWh
				job.AvgFPS = b.FPS
				job.DurationSec = b.DurationSec
				job.AvgPowerWatts = b.PowerWatts
				if metricCallback != nil {
					metricCallback(domain.TrainingMetrics{
						Epoch:            b.Epoch,
						TotalEpochs:      job.Hyperparameters.Epochs,
						CurrentBatch:     b.Batch,
						TotalBatches:     b.TotalBatches,
						BoxLoss:          b.BoxLoss,
						ClsLoss:          b.ClsLoss,
						PowerWatts:       b.PowerWatts,
						FPS:              b.FPS,
						EpochDurationSec: b.DurationSec,
					})
				}
			}
		}
	}
	if err := scanner.Err(); err != nil && ctx.Err() == nil {
		_ = err
	}

	return cmd.Wait()
}

// StopTraining terminates an active training process.
func (w *PythonWorker) StopTraining(ctx context.Context, jobID string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if cancel, ok := w.cancelFuncs[jobID]; ok {
		cancel()
		delete(w.activeJobs, jobID)
		delete(w.cancelFuncs, jobID)
		return nil
	}
	return nil
}
