package ports

import (
	"context"
	"hydraforge/internal/domain"
)

// GPUProvider queries physical GPU metrics.
type GPUProvider interface {
	QueryGPU(ctx context.Context) (*domain.GPUTelemetry, error)
}
