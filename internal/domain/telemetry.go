package domain

import "time"

// GPUTelemetry represents physical GPU sensor readings.
type GPUTelemetry struct {
	DeviceID      string    `json:"device_id"`
	Model         string    `json:"model"`
	TotalVRAMMB   float64   `json:"total_vram_mb"`
	UsedVRAMMB    float64   `json:"used_vram_mb"`
	VRAMUsagePct  float64   `json:"vram_usage_pct"`
	GPUUtilPct    float64   `json:"gpu_util_pct"`
	TempCelsius   float64   `json:"temp_celsius"`
	PowerWatts    float64   `json:"power_watts"`
	FanSpeedPct   float64   `json:"fan_speed_pct"`
	Timestamp     time.Time `json:"timestamp"`
}

// CockpitTelemetry aggregates system and training telemetry for the React HUD.
type CockpitTelemetry struct {
	ActiveJobsCount int              `json:"active_jobs_count"`
	TotalModelsTrained int           `json:"total_models_trained"`
	GPUStats        GPUTelemetry     `json:"gpu_stats"`
	ActiveJob       *TrainingJob     `json:"active_job,omitempty"`
	RecentMetrics   []TrainingMetrics `json:"recent_metrics,omitempty"`
}
