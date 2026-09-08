export interface GpuStats {
  name: string
  driver_version?: string
  cuda_version?: string
  utilization_pct: number
  vram_used_mb: number
  vram_total_mb: number
  temp_c: number
  power_w: number
}

export interface TelemetryData {
  status: string
  uptime_sec: number
  gpu_stats: GpuStats
  active_jobs: number
  active_job?: any
  cpu_util_pct?: number
  ram_used_mb?: number
  ram_total_mb?: number
}
