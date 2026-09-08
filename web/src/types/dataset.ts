export interface DatasetInfo {
  id: string
  dataset_id?: string
  name: string
  yaml_path: string
  classes: string[]
  train_count: number
  val_count: number
  test_count?: number
  format?: string
  status?: string
}
