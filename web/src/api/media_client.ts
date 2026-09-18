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
    let blob: Blob
    if (imageDataUrl && imageDataUrl.startsWith('data:')) {
      const parts = imageDataUrl.split(',')
      const mimeMatch = parts[0].match(/:(.*?);/)
      const mime = mimeMatch ? mimeMatch[1] : 'image/jpeg'
      const byteStr = atob(parts[1])
      const u8Arr = new Uint8Array(byteStr.length)
      for (let i = 0; i < byteStr.length; i++) {
        u8Arr[i] = byteStr.charCodeAt(i)
      }
      blob = new Blob([u8Arr], { type: mime })
    } else if (imageDataUrl) {
      const res = await fetch(imageDataUrl)
      blob = await res.blob()
    } else {
      const canvas = document.createElement('canvas')
      canvas.width = 320
      canvas.height = 240
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.fillStyle = '#07080c'
        ctx.fillRect(0, 0, 320, 240)
      }
      blob = await new Promise((resolve) => canvas.toBlob((b) => resolve(b || new Blob()), 'image/jpeg', 0.8))
    }

    const form = new FormData()
    form.append('image', blob, `snap_${Date.now()}.jpg`)
    form.append('dataset_id', datasetId || 'inbox_playground')
    form.append('camera_id', cameraId || 'video_loop_sample')
    if (bboxes && bboxes.length > 0) {
      form.append('bboxes', JSON.stringify(bboxes))
    }

    const host = typeof window !== 'undefined' && window.location.hostname ? window.location.hostname : 'localhost'
    const vaultUrl = `http://${host}:8082/api/v1/inbox/upload`

    const vaultRes = await fetch(vaultUrl, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer hydravault_vault_token'
      },
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
