package domain_test

import (
	"testing"

	"hydraforge/internal/domain"
)

func TestTrainingJobValidation(t *testing.T) {
	job := &domain.TrainingJob{
		JobID:             "job_yolo11s_custom",
		ModelArchitecture: "yolo11s",
		Hyperparameters: domain.Hyperparameters{
			Epochs:    100,
			ImageSize: 640,
		},
	}

	if err := job.Validate(); err != nil {
		t.Fatalf("expected valid job, got error: %v", err)
	}

	invalidJob := &domain.TrainingJob{
		JobID:             "",
		ModelArchitecture: "yolov8n",
	}
	if err := invalidJob.Validate(); err == nil {
		t.Fatalf("expected validation error for empty job_id, got nil")
	}
}

func TestDatasetValidation(t *testing.T) {
	ds := &domain.Dataset{
		DatasetID: "dataset_coco8",
		Classes:   []string{"person", "car"},
	}
	if err := ds.Validate(); err != nil {
		t.Fatalf("expected valid dataset, got: %v", err)
	}
}
