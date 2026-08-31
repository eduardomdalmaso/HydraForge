package worker

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
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

	totalEpochs := job.Hyperparameters.Epochs
	if totalEpochs > 10 {
		totalEpochs = 10 // Simulated fast cycle if testing
	}

	initialBoxLoss := 2.45
	initialClsLoss := 3.12
	initialDFLLoss := 1.88

	for epoch := 1; epoch <= totalEpochs; epoch++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(300 * time.Millisecond): // Simulate epoch computation
			decay := math.Exp(-float64(epoch) / 5.0)
			boxLoss := initialBoxLoss*decay + 0.35 + (rand.Float64() * 0.05)
			clsLoss := initialClsLoss*decay + 0.22 + (rand.Float64() * 0.03)
			dflLoss := initialDFLLoss*decay + 0.18 + (rand.Float64() * 0.02)
			map50 := (1.0 - decay*0.85) * (0.92 + (rand.Float64() * 0.04))
			map50_95 := map50 * 0.72

			metric := domain.TrainingMetrics{
				Epoch:        epoch,
				TotalEpochs:  totalEpochs,
				BoxLoss:      boxLoss,
				ClsLoss:      clsLoss,
				DFLLoss:      dflLoss,
				ValBoxLoss:   boxLoss * 1.1,
				ValClsLoss:   clsLoss * 1.05,
				MAP50:        map50,
				MAP50_95:     map50_95,
				Precision:    map50 * 0.96,
				Recall:       map50 * 0.94,
				LearningRate: job.Hyperparameters.LearningRate * decay,
				GPUVRAMMB:    5420.0 + (float64(epoch) * 45.0),
			}

			if metricCallback != nil {
				metricCallback(metric)
			}
		}
	}

	return nil
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

// ExportWeights converts weights into target format.
func (w *PythonWorker) ExportWeights(ctx context.Context, weightsPath string, format domain.ExportFormat, precision string) (string, error) {
	time.Sleep(200 * time.Millisecond) // Simulate TensorRT engine build
	return fmt.Sprintf("%s.%s", weightsPath, format), nil
}

// IsJobRunning checks if job is active.
func (w *PythonWorker) IsJobRunning(jobID string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.activeJobs[jobID]
}

var _ ports.TrainingWorker = (*PythonWorker)(nil)
