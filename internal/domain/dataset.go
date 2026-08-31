package domain

import (
	"fmt"
	"time"
)

// Dataset represents an annotated computer vision dataset.
type Dataset struct {
	DatasetID    string    `json:"dataset_id"`
	Name         string    `json:"name"`
	Task         TaskType  `json:"task"`
	YAMLPath     string    `json:"yaml_path"`
	Classes      []string  `json:"classes"`
	NumClasses   int       `json:"num_classes"`
	TrainImages  int       `json:"train_images"`
	ValImages    int       `json:"val_images"`
	TestImages   int       `json:"test_images"`
	SizeBytes    int64     `json:"size_bytes"`
	CreatedAt    time.Time `json:"created_at"`
}

// Validate verifies dataset configuration invariants.
func (d *Dataset) Validate() error {
	if d.DatasetID == "" {
		return fmt.Errorf("%w: dataset_id required", ErrInvalidDataset)
	}
	if len(d.Classes) == 0 {
		return fmt.Errorf("%w: at least 1 class name required", ErrInvalidDataset)
	}
	return nil
}
