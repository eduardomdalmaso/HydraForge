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

	// Vault Curation & Active Learning Errors
	ErrInvalidBBox         = errors.New("invalid bounding box coordinates")
	ErrInvalidPolygon      = errors.New("invalid polygon points")
	ErrInvalidFrame        = errors.New("invalid frame metadata")
	ErrFrameNotFound       = errors.New("frame not found")
	ErrDuplicateFrame      = errors.New("frame already exists or is a perceptual duplicate")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrInvalidPath         = errors.New("invalid file path or path traversal detected")
)
