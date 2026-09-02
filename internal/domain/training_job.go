package domain

import (
	"fmt"
	"time"
)

// TaskType defines the vision task for the YOLO model.
type TaskType string

const (
	TaskDetect   TaskType = "detect"
	TaskSegment  TaskType = "segment"
	TaskPose     TaskType = "pose"
	TaskClassify TaskType = "classify"
	TaskOBB      TaskType = "obb"
)

// JobStatus defines the state lifecycle of a training job.
type JobStatus string

const (
	StatusQueued    JobStatus = "QUEUED"
	StatusTraining  JobStatus = "TRAINING"
	StatusPaused    JobStatus = "PAUSED"
	StatusCompleted JobStatus = "COMPLETED"
	StatusFailed    JobStatus = "FAILED"
	StatusStopped   JobStatus = "STOPPED"
)

// Hyperparameters represents YOLO training hyperparameter configuration.
type Hyperparameters struct {
	Epochs        int     `json:"epochs"`
	BatchSize     int     `json:"batch_size"` // -1 for auto-batch
	ImageSize     int     `json:"imgsz"`      // 640, 1280, etc.
	Optimizer     string  `json:"optimizer"`  // "AdamW", "SGD", "RMSprop"
	LearningRate  float64 `json:"lr0"`
	FinalLR       float64 `json:"lrf"`
	Momentum      float64 `json:"momentum"`
	WeightDecay   float64 `json:"weight_decay"`
	WarmupEpochs  float64 `json:"warmup_epochs"`
	UseAMP        bool    `json:"amp"`           // Automatic Mixed Precision FP16
	CosineLR      bool    `json:"cos_lr"`        // Cosine learning rate scheduler
	Patience      int     `json:"patience"`      // Early stopping epochs
	DeviceID      string  `json:"device"`        // "0", "cuda:0", "cpu"
	Workers       int     `json:"workers"`       // Dataloader workers
	Mosaic        float64 `json:"mosaic"`        // Mosaic augmentation probability
	Mixup         float64 `json:"mixup"`         // Mixup augmentation probability
	Freeze        int     `json:"freeze"`        // Layer freeze depth (0, 10, 11, 23)
	TwoStage      bool    `json:"two_stage"`     // Enable Two-Stage Fine-Tuning Pipeline
	Stage1Epochs  int     `json:"stage1_epochs"` // Stage 1 Head Adaptation Epochs
	Stage1Freeze  int     `json:"stage1_freeze"` // Stage 1 Freeze depth (10=Backbone, 23=Head-Only)
	Stage2Epochs  int     `json:"stage2_epochs"` // Stage 2 Full Model Refinement Epochs
	CloseMosaic   int     `json:"close_mosaic"`  // Disable mosaic for final N epochs
	RecipePreset  string  `json:"recipe_preset"` // "default", "yolo26_recipe", "small_dataset", "two_stage"
	Pretrained    bool    `json:"pretrained"`    // Start from pretrained weights (.pt)
}

// TrainingMetrics represents epoch-level metrics emitted by PyTorch.
type TrainingMetrics struct {
	Epoch            int     `json:"epoch"`
	TotalEpochs      int     `json:"total_epochs"`
	CurrentBatch     int     `json:"current_batch,omitempty"`
	TotalBatches     int     `json:"total_batches,omitempty"`
	BoxLoss          float64 `json:"box_loss"`
	ClsLoss          float64 `json:"cls_loss"`
	DFLLoss          float64 `json:"dfl_loss"`
	ValBoxLoss       float64 `json:"val_box_loss"`
	ValClsLoss       float64 `json:"val_cls_loss"`
	MAP50            float64 `json:"map50"`
	MAP50_95         float64 `json:"map50_95"`
	Precision        float64 `json:"precision"`
	Recall           float64 `json:"recall"`
	LearningRate     float64 `json:"lr"`
	GPUVRAMMB        float64 `json:"gpu_vram_mb"`
	PowerWatts       float64 `json:"power_watts"`
	TempCelsius      float64 `json:"temp_celsius"`
	GPUUtilPct       float64 `json:"gpu_util_pct"`
	FPS              float64 `json:"fps"`
	EpochDurationSec float64 `json:"epoch_duration_sec"`
}

// TrainingJob represents a complete YOLO model training execution.
type TrainingJob struct {
	JobID             string          `json:"job_id"`
	ModelArchitecture string          `json:"model_architecture"` // "yolov8n", "yolo11s", "yolov26m"
	Task              TaskType        `json:"task"`
	DatasetID         string          `json:"dataset_id"`
	DatasetPath       string          `json:"dataset_path"`
	Hyperparameters   Hyperparameters `json:"hyperparameters"`
	Status            JobStatus       `json:"status"`
	CurrentEpoch      int             `json:"current_epoch"`
	TotalEpochs       int             `json:"total_epochs"`
	CurrentBatch      int             `json:"current_batch,omitempty"`
	TotalBatches      int             `json:"total_batches,omitempty"`
	BestMAP50         float64         `json:"best_map50"`
	BestMAP50_95      float64         `json:"best_map50_95"`
	Checkpoints       []string        `json:"checkpoints"`
	OutputWeights     string          `json:"output_weights"`
	TotalEnergyKWh    float64         `json:"total_energy_kwh"`
	AvgPowerWatts     float64         `json:"avg_power_watts"`
	PeakVRAMMB        float64         `json:"peak_vram_mb"`
	AvgFPS            float64         `json:"avg_fps"`
	DurationSec       float64         `json:"duration_sec"`
	ErrorMessage      string          `json:"error_message,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// Validate verifies domain business invariants.
func (j *TrainingJob) Validate() error {
	if j.JobID == "" {
		return fmt.Errorf("%w: job_id required", ErrInvalidJobConfig)
	}
	if j.ModelArchitecture == "" {
		return fmt.Errorf("%w: model_architecture required", ErrInvalidJobConfig)
	}
	if j.Hyperparameters.Epochs <= 0 || j.Hyperparameters.Epochs > 1000 {
		return fmt.Errorf("%w: epochs must be between 1 and 1000", ErrInvalidJobConfig)
	}
	if j.Hyperparameters.ImageSize%32 != 0 || j.Hyperparameters.ImageSize <= 0 {
		return fmt.Errorf("%w: imgsz must be a positive multiple of 32", ErrInvalidJobConfig)
	}
	return nil
}

// SetDefaults applies standard defaults.
func (j *TrainingJob) SetDefaults() {
	if j.JobID == "" {
		arch := j.ModelArchitecture
		if arch == "" {
			arch = "yolo26n"
		}
		j.JobID = fmt.Sprintf("run_%s_%d", arch, time.Now().UnixMilli())
	}
	if j.CreatedAt.IsZero() {
		j.CreatedAt = time.Now()
	}
	j.UpdatedAt = time.Now()
	if j.Status == "" {
		j.Status = StatusQueued
	}
	if j.Task == "" {
		j.Task = TaskDetect
	}
	if j.Hyperparameters.Epochs == 0 {
		j.Hyperparameters.Epochs = 100
	}
	if j.Hyperparameters.BatchSize == 0 {
		j.Hyperparameters.BatchSize = 16
	}
	if j.Hyperparameters.ImageSize == 0 {
		j.Hyperparameters.ImageSize = 640
	}
	if j.Hyperparameters.Optimizer == "" {
		j.Hyperparameters.Optimizer = "AdamW"
	}
	if j.Hyperparameters.LearningRate == 0 {
		j.Hyperparameters.LearningRate = 0.001
	}
	if j.Hyperparameters.DeviceID == "" {
		j.Hyperparameters.DeviceID = "0"
	}
	if j.Hyperparameters.TwoStage {
		if j.Hyperparameters.Stage1Epochs == 0 {
			j.Hyperparameters.Stage1Epochs = 20
		}
		if j.Hyperparameters.Stage1Freeze == 0 {
			j.Hyperparameters.Stage1Freeze = 10
		}
		if j.Hyperparameters.Stage2Epochs == 0 {
			j.Hyperparameters.Stage2Epochs = 30
		}
		j.TotalEpochs = j.Hyperparameters.Stage1Epochs + j.Hyperparameters.Stage2Epochs
		j.Hyperparameters.Epochs = j.TotalEpochs
	} else if j.TotalEpochs == 0 {
		j.TotalEpochs = j.Hyperparameters.Epochs
	}
}
