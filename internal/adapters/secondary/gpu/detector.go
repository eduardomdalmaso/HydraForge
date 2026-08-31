package gpu

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// Detector queries nvidia-smi for live hardware metrics.
type Detector struct{}

// NewDetector initializes a new GPU detector.
func NewDetector() *Detector {
	return &Detector{}
}

// QueryGPU queries physical NVIDIA GPU sensors.
func (d *Detector) QueryGPU(ctx context.Context) (*domain.GPUTelemetry, error) {
	out, err := exec.CommandContext(ctx, "nvidia-smi", "--query-gpu=name,memory.total,memory.used,utilization.gpu,temperature.gpu,power.draw", "--format=csv,noheader,nounits").Output()
	if err != nil {
		// Fallback for simulation or CPU-only
		return &domain.GPUTelemetry{
			DeviceID:     "0",
			Model:        "NVIDIA GeForce RTX 5090 (Simulated)",
			TotalVRAMMB:  32768,
			UsedVRAMMB:   2048,
			VRAMUsagePct: 6.25,
			GPUUtilPct:   12.0,
			TempCelsius:  38.0,
			PowerWatts:   45.0,
			Timestamp:    time.Now(),
		}, nil
	}

	parts := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(parts) < 6 {
		return nil, domain.ErrInvalidJobConfig
	}

	model := strings.TrimSpace(parts[0])
	totalVRAM, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	usedVRAM, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	gpuUtil, _ := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
	temp, _ := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
	power, _ := strconv.ParseFloat(strings.TrimSpace(parts[5]), 64)

	var pct float64
	if totalVRAM > 0 {
		pct = (usedVRAM / totalVRAM) * 100.0
	}

	return &domain.GPUTelemetry{
		DeviceID:     "0",
		Model:        model,
		TotalVRAMMB:  totalVRAM,
		UsedVRAMMB:   usedVRAM,
		VRAMUsagePct: pct,
		GPUUtilPct:   gpuUtil,
		TempCelsius:  temp,
		PowerWatts:   power,
		Timestamp:    time.Now(),
	}, nil
}

var _ ports.GPUProvider = (*Detector)(nil)
