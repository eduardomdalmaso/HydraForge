export interface BenchmarkRun {
  id: string
  name: string
  device: string
  batch_size: number
  imgsz: number
  fps_throughput: number
  avg_latency_ms: number
  p99_latency_ms: number
  gpu_power_watts: number
  timestamp: string
}
