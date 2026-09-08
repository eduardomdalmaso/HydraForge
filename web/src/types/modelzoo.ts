export interface ModelZooItem {
  id: string
  name: string
  family: string
  size: string
  task: 'detect' | 'segment' | 'pose' | 'obb' | 'classify'
  params_m: number
  flops_g: number
  map50_95: number
  inference_speed_ms: number
  status: 'downloaded' | 'remote' | 'trained'
  description: string
}
