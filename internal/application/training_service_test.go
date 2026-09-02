package application_test

import (
	"context"
	"testing"

	"hydraforge/internal/adapters/secondary/gpu"
	"hydraforge/internal/adapters/secondary/memory"
	"hydraforge/internal/adapters/secondary/worker"
	"hydraforge/internal/application"
	"hydraforge/internal/domain"
)

func TestTrainingServiceLifecycle(t *testing.T) {
	memStore := memory.NewMemoryStore()
	detector := gpu.NewDetector()
	pyWorker := worker.NewPythonWorker()

	service := application.NewTrainingService(memStore, memStore, memStore, memStore, pyWorker, detector)


	// 1. Create Job
	job := &domain.TrainingJob{
		JobID:             "test_training_run_01",
		ModelArchitecture: "yolo11n",
		Task:              domain.TaskDetect,
		Hyperparameters: domain.Hyperparameters{
			Epochs:    2,
			BatchSize: 8,
			ImageSize: 640,
		},
	}

	created, err := service.CreateTrainingJob(context.Background(), job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}
	if created.JobID != "test_training_run_01" {
		t.Fatalf("expected job ID test_training_run_01, got %s", created.JobID)
	}

	// 2. Fetch Job
	fetched, err := service.GetTrainingJob(context.Background(), "test_training_run_01")
	if err != nil {
		t.Fatalf("failed to fetch job: %v", err)
	}
	if fetched.ModelArchitecture != "yolo11n" {
		t.Fatalf("expected yolo11n, got %s", fetched.ModelArchitecture)
	}

	// 3. Telemetry
	telemetry, err := service.GetCockpitTelemetry(context.Background())
	if err != nil {
		t.Fatalf("failed to get cockpit telemetry: %v", err)
	}
	if telemetry.TotalModelsTrained < 1 {
		t.Fatalf("expected at least 1 job tracked, got %d", telemetry.TotalModelsTrained)
	}
}
