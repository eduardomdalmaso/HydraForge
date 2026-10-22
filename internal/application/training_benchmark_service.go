package application

import (
	"context"

	"hydraforge/internal/domain"
)

// CreateBenchmarkJob validates and dispatches a multi-format benchmark run.
func (s *TrainingService) CreateBenchmarkJob(ctx context.Context, job *domain.BenchmarkJob) (*domain.BenchmarkJob, error) {
	if job == nil {
		return nil, domain.ErrInvalidJobConfig
	}
	job.SetDefaults()
	if err := job.Validate(); err != nil {
		return nil, err
	}

	if err := s.benchmarkRepo.SaveBenchmark(ctx, job); err != nil {
		return nil, err
	}

	go func(targetJob domain.BenchmarkJob) {
		_ = s.benchmarkRepo.UpdateBenchmarkStatus(context.Background(), targetJob.JobID, domain.BenchmarkStatusRunning, "")

		err := s.worker.RunBenchmark(context.Background(), &targetJob, func(res domain.FormatBenchmarkResult) {
			_ = s.benchmarkRepo.AddBenchmarkResult(context.Background(), targetJob.JobID, res)
		})

		if err != nil {
			_ = s.benchmarkRepo.UpdateBenchmarkStatus(context.Background(), targetJob.JobID, domain.BenchmarkStatusFailed, err.Error())
		} else {
			_ = s.benchmarkRepo.UpdateBenchmarkStatus(context.Background(), targetJob.JobID, domain.BenchmarkStatusCompleted, "")
		}
	}(*job)

	return job, nil
}

// GetBenchmarkJob retrieves a specific benchmark execution with its format results.
func (s *TrainingService) GetBenchmarkJob(ctx context.Context, jobID string) (*domain.BenchmarkJob, error) {
	return s.benchmarkRepo.GetBenchmark(ctx, jobID)
}

// ListBenchmarkJobs returns all benchmark jobs filtered optionally by status.
func (s *TrainingService) ListBenchmarkJobs(ctx context.Context, status string) ([]*domain.BenchmarkJob, error) {
	return s.benchmarkRepo.ListBenchmarks(ctx, status)
}

// StopBenchmarkJob cancels a running benchmark execution.
func (s *TrainingService) StopBenchmarkJob(ctx context.Context, jobID string) error {
	if err := s.worker.StopBenchmark(ctx, jobID); err != nil {
		return err
	}
	return s.benchmarkRepo.UpdateBenchmarkStatus(ctx, jobID, domain.BenchmarkStatusStopped, "Manually stopped by operator")
}

// GetSupportedBenchmarkFormats returns the complete catalog of Ultralytics export formats with metadata.
func (s *TrainingService) GetSupportedBenchmarkFormats(ctx context.Context) ([]map[string]interface{}, error) {
	formats := []map[string]interface{}{
		{"format": "PyTorch", "arg": "-", "extension": ".pt", "metadata": true, "target_hardware": "NVIDIA GPU / CPU Baseline", "speedup": "1.0x (Baseline)"},
		{"format": "TensorRT", "arg": "engine", "extension": ".engine", "metadata": true, "target_hardware": "NVIDIA RTX 5090 (Max Throughput)", "speedup": "Up to 5.5x GPU Speedup"},
		{"format": "ONNX", "arg": "onnx", "extension": ".onnx", "metadata": true, "target_hardware": "Universal / CPU / Edge", "speedup": "Up to 2.5x CPU Speedup"},
		{"format": "OpenVINO", "arg": "openvino", "extension": "_openvino_model", "metadata": true, "target_hardware": "Intel CPU / iGPU / NPU", "speedup": "Up to 3.0x Intel Speedup"},
		{"format": "TorchScript", "arg": "torchscript", "extension": ".torchscript", "metadata": true, "target_hardware": "C++ libtorch Runtime", "speedup": "1.2x C++ Speedup"},
		{"format": "CoreML", "arg": "coreml", "extension": ".mlpackage", "metadata": true, "target_hardware": "Apple Neural Engine (ANE)", "speedup": "Up to 4.0x Apple Silicon"},
		{"format": "LiteRT", "arg": "litert", "extension": ".tflite", "metadata": true, "target_hardware": "Android / Edge Devices", "speedup": "Optimized Lightweight"},
		{"format": "NCNN", "arg": "ncnn", "extension": "_ncnn_model", "metadata": true, "target_hardware": "Mobile ARM Vulkan", "speedup": "Vulkan Acceleration"},
		{"format": "RKNN", "arg": "rknn", "extension": "_rknn_model", "metadata": true, "target_hardware": "Rockchip NPU (RK3588)", "speedup": "Dedicated NPU"},
		{"format": "Hailo", "arg": "hailo", "extension": "_hailo_model", "metadata": true, "target_hardware": "Hailo-8 / Hailo-15 NPU", "speedup": "Ultra Low Power Edge"},
		{"format": "Qualcomm QNN", "arg": "qnn", "extension": "_qnn.onnx", "metadata": true, "target_hardware": "Snapdragon NPU / Hexagon", "speedup": "Qualcomm NPU Speedup"},
	}
	return formats, nil
}
