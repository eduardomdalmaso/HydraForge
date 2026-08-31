package domain

import "time"

// ExportFormat defines target deployment artifact formats.
type ExportFormat string

const (
	FormatTensorRT ExportFormat = "engine"
	FormatONNX     ExportFormat = "onnx"
	FormatTorchScript ExportFormat = "torchscript"
	FormatOpenVINO ExportFormat = "openvino"
)

// ModelCheckpoint represents a saved model weight file.
type ModelCheckpoint struct {
	CheckpointID string       `json:"checkpoint_id"`
	JobID        string       `json:"job_id"`
	Architecture string       `json:"architecture"`
	Task         TaskType     `json:"task"`
	Epoch        int          `json:"epoch"`
	MAP50        float64      `json:"map50"`
	MAP50_95     float64      `json:"map50_95"`
	WeightsPath  string       `json:"weights_path"`
	ExportFormat ExportFormat `json:"export_format,omitempty"`
	ExportPath   string       `json:"export_path,omitempty"`
	SizeBytes    int64        `json:"size_bytes"`
	Precision    string       `json:"precision"` // "FP32", "FP16", "INT8"
	CreatedAt    time.Time    `json:"created_at"`
}
