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

	rawLines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var validLines []string
	for _, l := range rawLines {
		if strings.TrimSpace(l) != "" {
			validLines = append(validLines, strings.TrimSpace(l))
		}
	}

	if len(validLines) == 0 {
		return nil, domain.ErrInvalidJobConfig
	}

	var primaryModel string
	var sumTotalVRAM, sumUsedVRAM, sumPower, sumGPUUtil, maxTemp float64
	deviceCount := len(validLines)
	var devices []domain.GPUTelemetry

	for i, line := range validLines {
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}
		mod := strings.TrimSpace(parts[0])
		if i == 0 {
			primaryModel = mod
		}
		tVRAM, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		uVRAM, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		gUtil, _ := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		tTemp, _ := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
		pPwr, _ := strconv.ParseFloat(strings.TrimSpace(parts[5]), 64)

		sumTotalVRAM += tVRAM
		sumUsedVRAM += uVRAM
		sumPower += pPwr
		sumGPUUtil += gUtil
		if tTemp > maxTemp {
			maxTemp = tTemp
		}

		devPct := 0.0
		if tVRAM > 0 {
			devPct = (uVRAM / tVRAM) * 100.0
		}
		devices = append(devices, domain.GPUTelemetry{
			DeviceID:     strconv.Itoa(i),
			DeviceCount:  1,
			Model:        mod,
			TotalVRAMMB:  tVRAM,
			UsedVRAMMB:   uVRAM,
			VRAMUsagePct: devPct,
			GPUUtilPct:   gUtil,
			TempCelsius:  tTemp,
			PowerWatts:   pPwr,
			Timestamp:    time.Now(),
		})
	}

	displayName := primaryModel
	if deviceCount > 1 {
		displayName = strconv.Itoa(deviceCount) + "x " + primaryModel
	}

	var pct float64
	if sumTotalVRAM > 0 {
		pct = (sumUsedVRAM / sumTotalVRAM) * 100.0
	}

	avgUtil := sumGPUUtil / float64(deviceCount)

	return &domain.GPUTelemetry{
		DeviceID:     "0",
		DeviceCount:  deviceCount,
		Model:        displayName,
		TotalVRAMMB:  sumTotalVRAM,
		UsedVRAMMB:   sumUsedVRAM,
		VRAMUsagePct: pct,
		GPUUtilPct:   avgUtil,
		TempCelsius:  maxTemp,
		PowerWatts:   sumPower,
		Devices:      devices,
		Timestamp:    time.Now(),
	}, nil
}

var _ ports.GPUProvider = (*Detector)(nil)
