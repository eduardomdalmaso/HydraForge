package domain

import (
	"testing"
)

func TestBenchmarkJob_Validate(t *testing.T) {
	tests := []struct {
		name    string
		job     BenchmarkJob
		wantErr bool
	}{
		{
			name: "valid default job",
			job: BenchmarkJob{
				JobID:     "bmk_test_1",
				Model:     "yolo26n.pt",
				ImageSize: 640,
				Quantize:  16,
			},
			wantErr: false,
		},
		{
			name: "missing job_id",
			job: BenchmarkJob{
				Model:     "yolo26n.pt",
				ImageSize: 640,
			},
			wantErr: true,
		},
		{
			name: "missing model",
			job: BenchmarkJob{
				JobID:     "bmk_test_2",
				ImageSize: 640,
			},
			wantErr: true,
		},
		{
			name: "invalid imgsz non-multiple of 32",
			job: BenchmarkJob{
				JobID:     "bmk_test_3",
				Model:     "yolo26n.pt",
				ImageSize: 500,
			},
			wantErr: true,
		},
		{
			name: "invalid quantize value",
			job: BenchmarkJob{
				JobID:     "bmk_test_4",
				Model:     "yolo26n.pt",
				ImageSize: 640,
				Quantize:  64,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.job.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BenchmarkJob.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBenchmarkJob_SetDefaults(t *testing.T) {
	job := BenchmarkJob{}
	job.SetDefaults()

	if job.JobID == "" {
		t.Errorf("expected JobID to be generated")
	}
	if job.Model != "yolo26n.pt" {
		t.Errorf("expected default model yolo26n.pt, got %s", job.Model)
	}
	if job.ImageSize != 640 {
		t.Errorf("expected default imgsz 640, got %d", job.ImageSize)
	}
	if job.Status != BenchmarkStatusQueued {
		t.Errorf("expected default status QUEUED, got %s", job.Status)
	}
}
