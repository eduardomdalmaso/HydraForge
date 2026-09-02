package application_test

import (
	"context"
	"testing"
	"time"

	"hydraforge/internal/adapters/secondary/gpu"
	"hydraforge/internal/adapters/secondary/memory"
	"hydraforge/internal/adapters/secondary/worker"
	"hydraforge/internal/application"
	"hydraforge/internal/domain"
)

func TestBenchmarkServiceLifecycle(t *testing.T) {
	memStore := memory.NewMemoryStore()
	detector := gpu.NewDetector()
	pyWorker := worker.NewPythonWorker()

	service := application.NewTrainingService(memStore, memStore, memStore, memStore, pyWorker, detector)

	// 1. Create Benchmark Job
	job := &domain.BenchmarkJob{
		JobID:         "bmk_test_yolo26n",
		Model:         "yolo26n.pt",
		Data:          "coco8.yaml",
		ImageSize:     640,
		Quantize:      16,
		Device:        "0",
		TargetFormats: []string{"PyTorch", "ONNX", "TensorRT"},
	}

	created, err := service.CreateBenchmarkJob(context.Background(), job)
	if err != nil {
		t.Fatalf("failed to create benchmark job: %v", err)
	}
	if created.JobID != "bmk_test_yolo26n" {
		t.Fatalf("expected bmk_test_yolo26n, got %s", created.JobID)
	}

	// 2. Fetch Benchmark Job
	fetched, err := service.GetBenchmarkJob(context.Background(), "bmk_test_yolo26n")
	if err != nil {
		t.Fatalf("failed to get benchmark job: %v", err)
	}
	if fetched.Model != "yolo26n.pt" {
		t.Fatalf("expected yolo26n.pt, got %s", fetched.Model)
	}

	// 3. List Benchmark Jobs
	list, err := service.ListBenchmarkJobs(context.Background(), "")
	if err != nil {
		t.Fatalf("failed to list benchmark jobs: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 benchmark job, got %d", len(list))
	}

	// 4. Supported formats
	formats, err := service.GetSupportedBenchmarkFormats(context.Background())
	if err != nil {
		t.Fatalf("failed to get supported formats: %v", err)
	}
	if len(formats) < 5 {
		t.Fatalf("expected at least 5 supported export formats, got %d", len(formats))
	}

	// 5. Wait briefly for worker simulation to populate results
	time.Sleep(900 * time.Millisecond)

	updated, err := service.GetBenchmarkJob(context.Background(), "bmk_test_yolo26n")
	if err != nil {
		t.Fatalf("failed to get updated job: %v", err)
	}
	if len(updated.Results) == 0 {
		t.Fatalf("expected at least 1 benchmark format result, got 0")
	}
}
