/* ==========================================================================
   HYDRAVAULT CLIENT MODULE (TypeScript) — Curation & Active Learning Plane
   ========================================================================== */

export interface VaultHealth {
  online: boolean
  service?: string
  version?: string
  message?: string
}

export interface VaultClassMetadata {
  id: number
  name: string
  count?: number
}

export interface VaultDataset {
  dataset_id: string
  name: string
  description: string
  task: string
  classes: VaultClassMetadata[]
  total_frames: number
}

function getVaultBaseUrl(): string {
  const host = typeof window !== 'undefined' && window.location.hostname ? window.location.hostname : 'localhost'
  return `http://${host}:8082`
}

export async function checkVaultHealth(): Promise<VaultHealth> {
  try {
    const res = await fetch(`${getVaultBaseUrl()}/health`, { method: 'GET' })
    if (res.ok) {
      const data = await res.json()
      return { online: true, service: data.service, version: data.version }
    }
    return { online: false, message: `HTTP ${res.status}` }
  } catch (err: any) {
    return { online: false, message: err.message || 'Connection refused' }
  }
}

export async function fetchVaultDatasets(): Promise<VaultDataset[]> {
  try {
    const res = await fetch(`${getVaultBaseUrl()}/api/v1/datasets`, {
      method: 'GET',
      headers: { 'Authorization': 'Bearer hydravault_vault_token' }
    })
    if (!res.ok) return []
    return await res.json()
  } catch {
    return []
  }
}

export async function sendDirectFrameToVault(
  datasetId: string,
  cameraId: string,
  imageBlob: Blob,
  bboxes: any[] = []
): Promise<{ success: boolean; frameId?: string; message: string }> {
  try {
    const form = new FormData()
    form.append('image', imageBlob, `frame_${Date.now()}.jpg`)
    form.append('dataset_id', datasetId || 'frota_urbana_fusion')
    form.append('camera_id', cameraId || 'playground_stream')
    if (bboxes && bboxes.length > 0) {
      form.append('bboxes', JSON.stringify(bboxes))
    }

    const res = await fetch(`${getVaultBaseUrl()}/api/v1/inbox/upload`, {
      method: 'POST',
      headers: { 'Authorization': 'Bearer hydravault_vault_token' },
      body: form
    })

    if (!res.ok) {
      const txt = await res.text()
      return { success: false, message: `Vault HTTP ${res.status}: ${txt}` }
    }
    const data = await res.json()
    return { success: true, frameId: data.frame_id, message: 'Frame enviado ao Vault com sucesso' }
  } catch (err: any) {
    return { success: false, message: `Falha de rede Vault: ${err.message}` }
  }
}
