export interface DatasetInfo {
  id?: string
  dataset_id?: string
  name: string
  yaml_path?: string
  classes?: string[]
  num_classes?: number
  train_images?: number
  val_images?: number
  test_images?: number
  train_count?: number
  val_count?: number
  test_count?: number
  format?: string
  status?: string
}
