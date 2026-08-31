package domain

import "errors"

var (
	ErrJobNotFound         = errors.New("training job not found")
	ErrInvalidJobConfig    = errors.New("invalid training job configuration")
	ErrDatasetNotFound     = errors.New("dataset not found")
	ErrInvalidDataset      = errors.New("invalid dataset configuration")
	ErrJobAlreadyRunning   = errors.New("a training job is already active on target GPU")
	ErrCheckpointNotFound  = errors.New("model checkpoint not found")
	ErrExportFailed        = errors.New("failed to export model format")
)
