/* ==========================================================================
   HYDRAFORGE API CLIENT MODULE (TypeScript)
   ========================================================================== */

import type { TelemetryData } from '../types/telemetry'
import type { DatasetInfo } from '../types/dataset'

export async function fetchTelemetryAPI(): Promise<TelemetryData | null> {
  try {
    const res = await fetch('/api/v1/training/telemetry')
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return null
  }
}

export async function fetchDatasetsAPI(): Promise<DatasetInfo[]> {
  try {
    const res = await fetch('/api/v1/training/datasets')
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return []
  }
}

export async function fetchJobsAPI(status = ''): Promise<any[]> {
  try {
    const res = await fetch(`/api/v1/training/jobs?status=${status}`)
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return []
  }
}

export async function launchTrainingJobAPI(jobConfig: Record<string, any>): Promise<any> {
  const res = await fetch('/api/v1/training/jobs', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(jobConfig)
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error || `HTTP error ${res.status}`)
  }
  return await res.json()
}

export async function stopTrainingJobAPI(jobId: string): Promise<boolean> {
  try {
    const res = await fetch(`/api/v1/training/jobs/${jobId}`, { method: 'DELETE' })
    return res.ok
  } catch {
    return false
  }
}

export async function fetchModelsAPI(): Promise<any[]> {
  try {
    const res = await fetch('/api/v1/training/models')
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return []
  }
}

export async function rescanDatasetsAPI(): Promise<DatasetInfo[]> {
  try {
    const res = await fetch('/api/v1/training/datasets/rescan', { method: 'POST' })
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return []
  }
}

export async function registerDatasetPathAPI(payload: { name: string; yaml_path: string }): Promise<DatasetInfo> {
  const res = await fetch('/api/v1/training/datasets/register-path', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  return await res.json()
}

export async function deleteDatasetAPI(datasetId: string, deleteFiles = false): Promise<boolean> {
  const res = await fetch(`/api/v1/training/datasets/${datasetId}?delete_files=${deleteFiles}`, { method: 'DELETE' })
  return res.ok
}
