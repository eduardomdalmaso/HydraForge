package domain

import (
	"fmt"
	"time"
)

// BenchmarkStatus defines the lifecycle status of a benchmark job.
type BenchmarkStatus string

const (
	BenchmarkStatusQueued    BenchmarkStatus = "QUEUED"
	BenchmarkStatusRunning   BenchmarkStatus = "RUNNING"
	BenchmarkStatusCompleted BenchmarkStatus = "COMPLETED"
	BenchmarkStatusFailed    BenchmarkStatus = "FAILED"
	BenchmarkStatusStopped   BenchmarkStatus = "STOPPED"
)

// FormatBenchmarkResult stores benchmark metrics for a single exported runtime format.
type FormatBenchmarkResult struct {
	Format          string  `json:"format"`            // "pytorch", "onnx", "engine", "openvino", etc.
	Status          string  `json:"status"`            // "SUCCESS", "FAILED", "SKIPPED"
	SizeMB          float64 `json:"size_mb"`           // File or directory size in megabytes
	InferenceTimeMS float64 `json:"inference_time_ms"` // Average inference time per image in milliseconds
	FPS             float64 `json:"fps"`               // Throughput (Frames Per Second)
	MAP50_95        float64 `json:"map50_95"`          // Accuracy on detection/segmentation/pose
	AccuracyTop1    float64 `json:"accuracy_top1"`     // Accuracy on classification tasks
	ExportArgs      string  `json:"export_args"`       // Arguments used during export/quantization
	ErrorMessage    string  `json:"error_message,omitempty"`
}

// BenchmarkJob represents a multi-format benchmark run across runtime formats.
type BenchmarkJob struct {
	JobID         string                  `json:"job_id"`
	Model         string                  `json:"model"`          // e.g., "yolo26n.pt", "best.pt"
	Data          string                  `json:"data"`           // e.g., "coco8.yaml"
	ImageSize     int                     `json:"imgsz"`          // e.g., 640
	Quantize      int                     `json:"quantize"`       // 16 (FP16), 8 (INT8), 32/0 (FP32)
	Device        string                  `json:"device"`         // "0", "cuda:0", "cpu"
	TargetFormats []string                `json:"target_formats"` // List of formats or empty for all
	Status        BenchmarkStatus         `json:"status"`
	Results       []FormatBenchmarkResult `json:"results"`
	Summary       string                  `json:"summary,omitempty"`
	ErrorMessage  string                  `json:"error_message,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// SetDefaults applies sensible defaults to the benchmark job.
func (b *BenchmarkJob) SetDefaults() {
	if b.JobID == "" {
		b.JobID = fmt.Sprintf("bmk_%d", time.Now().UnixNano())
	}
	if b.Model == "" {
		b.Model = "yolo26n.pt"
	}
	if b.Data == "" {
		b.Data = "coco8.yaml"
	}
	if b.ImageSize <= 0 {
		b.ImageSize = 640
	}
	if b.Device == "" {
		b.Device = "0"
	}
	if b.Status == "" {
		b.Status = BenchmarkStatusQueued
	}
	if b.Results == nil {
		b.Results = make([]FormatBenchmarkResult, 0)
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	b.UpdatedAt = time.Now().UTC()
}

// Validate checks business invariants according to Ultralytics Benchmark rules.
func (b *BenchmarkJob) Validate() error {
	if b.JobID == "" {
		return fmt.Errorf("%w: job_id required", ErrInvalidJobConfig)
	}
	if b.Model == "" {
		return fmt.Errorf("%w: model required", ErrInvalidJobConfig)
	}
	if b.ImageSize%32 != 0 || b.ImageSize <= 0 {
		return fmt.Errorf("%w: imgsz must be a positive multiple of 32", ErrInvalidJobConfig)
	}
	if b.Quantize != 0 && b.Quantize != 8 && b.Quantize != 16 && b.Quantize != 32 {
		return fmt.Errorf("%w: quantize must be 8, 16, 32 or unset (0)", ErrInvalidJobConfig)
	}
	return nil
}
