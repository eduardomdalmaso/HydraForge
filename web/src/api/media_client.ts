export interface MediaFile {
  id: string
  name: string
  folder: string
  path: string
  size_bytes: number
  mod_time: string
  stream_url: string
}

export interface MediaFolder {
  name: string
  label: string
  file_count: number
  total_bytes: number
  files: MediaFile[]
}

export interface MediaSourcesResponse {
  total_folders: number
  total_files: number
  folders: MediaFolder[]
}

export async function fetchMediaSources(): Promise<MediaSourcesResponse> {
  try {
    const res = await fetch('/api/v1/media/sources')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return await res.json()
  } catch (err) {
    console.error('Failed to fetch media sources:', err)
    return { total_folders: 0, total_files: 0, folders: [] }
  }
}

export async function createMediaFolder(name: string): Promise<boolean> {
  try {
    const res = await fetch('/api/v1/media/folders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name })
    })
    return res.ok
  } catch {
    return false
  }
}

export async function deleteMediaFolder(name: string): Promise<boolean> {
  try {
    const res = await fetch(`/api/v1/media/folders/${encodeURIComponent(name)}`, {
      method: 'DELETE'
    })
    return res.ok
  } catch {
    return false
  }
}

export async function uploadMediaVideo(folder: string, file: File): Promise<boolean> {
  try {
    const form = new FormData()
    form.append('folder', folder)
    form.append('file', file)
    const res = await fetch('/api/v1/media/upload', {
      method: 'POST',
      body: form
    })
    return res.ok
  } catch {
    return false
  }
}

export async function deleteMediaFile(folder: string, filename: string): Promise<boolean> {
  try {
    const res = await fetch(`/api/v1/media/files/${encodeURIComponent(folder)}/${encodeURIComponent(filename)}`, {
      method: 'DELETE'
    })
    return res.ok
  } catch {
    return false
  }
}

export async function sendToHydraVaultInbox(
  datasetId: string,
  cameraId: string,
  imageDataUrl: string,
  bboxes: any[] = []
): Promise<{ success: boolean; message: string }> {
  try {
    // Convert dataURL to Blob
    const res = await fetch(imageDataUrl)
    const blob = await res.blob()

    const form = new FormData()
    form.append('image', blob, `snap_${Date.now()}.jpg`)
    form.append('dataset_id', datasetId || 'inbox_playground')
    form.append('camera_id', cameraId || 'video_loop_sample')
    if (bboxes && bboxes.length > 0) {
      form.append('bboxes', JSON.stringify(bboxes))
    }

    const vaultRes = await fetch('http://localhost:8082/api/v1/inbox/upload', {
      method: 'POST',
      body: form
    })

    if (!vaultRes.ok) {
      const errTxt = await vaultRes.text()
      return { success: false, message: `HydraVault HTTP ${vaultRes.status}: ${errTxt}` }
    }
    return { success: true, message: 'Frame enviado com sucesso ao Inbox do HydraVault (:8082)!' }
  } catch (err: any) {
    return { success: false, message: `Falha ao conectar no HydraVault (:8082): ${err.message}` }
  }
}
