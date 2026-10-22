package http

// BoxAnnotation represents a single detected or labeled bounding box in sample inspector.
type BoxAnnotation struct {
	ClassID   string    `json:"class_id"`
	ClassName string    `json:"class_name"`
	BBox      []float64 `json:"bbox"`
	IsTarget  bool      `json:"is_target"`
}

// SampleMatch represents a matched sample image with its bounding boxes.
type SampleMatch struct {
	ImageURL    string          `json:"image_url"`
	Filename    string          `json:"filename"`
	BBox        []float64       `json:"bbox"`
	ClassID     string          `json:"class_id"`
	Annotations []BoxAnnotation `json:"annotations"`
}

type candidatePair struct {
	imagesDir string
	labelsDir string
	subRel    string
}
