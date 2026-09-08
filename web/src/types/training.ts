export interface TrainingParams {
  model: string
  dataset_id: string
  epochs: number
  batch: number
  imgsz: number
  device: string
  optimizer: string
  amp: boolean
  compile?: boolean
  workers: number
  lr0: number
  patience: number
  close_mosaic: number
}

export interface TrainingMetricPayload {
  epoch: number
  total_epochs: number
  box_loss: number
  cls_loss: number
  dfl_loss: number
  val_box_loss: number
  val_cls_loss: number
  map50: number
  map50_95: number
  precision: number
  recall: number
  lr: number
  gpu_vram_mb: number
  power_watts: number
  temp_celsius: number
  gpu_util_pct: number
  fps: number
  epoch_duration_sec: number
}
