/* ==========================================================================
   HYDRASTREAM INTEGRATION API CLIENT MODULE (TypeScript)
   ========================================================================== */

const HYDRASTREAM_BASE = '/api/v1/hydrastream'

export async function fetchHydraStreams(): Promise<any[]> {
  try {
    const res = await fetch(`${HYDRASTREAM_BASE}/api/v1/streams`)
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    const data = await res.json()
    return (data && Array.isArray(data.streams)) ? data.streams : []
  } catch {
    try {
      const fallback = await fetch('http://localhost:8080/api/v1/streams')
      if (!fallback.ok) return []
      const data = await fallback.json()
      return (data && Array.isArray(data.streams)) ? data.streams : []
    } catch {
      return []
    }
  }
}

export async function fetchHydraTelemetry(): Promise<any | null> {
  try {
    const res = await fetch(`${HYDRASTREAM_BASE}/api/v1/telemetry/stats`)
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    try {
      const fallback = await fetch('http://localhost:8080/api/v1/telemetry/stats')
      if (fallback.ok) return await fallback.json()
    } catch {}
    return null
  }
}

export function getStreamSnapshotUrl(streamId?: string): string {
  if (!streamId) return '/hydra-logo.jpg'
  return `${HYDRASTREAM_BASE}/api/v1/streams/${encodeURIComponent(streamId)}/snapshot.jpg?t=${Date.now()}`
}

export function getStreamMJPEGUrl(streamId?: string): string {
  if (!streamId) return ''
  return `${HYDRASTREAM_BASE}/api/v1/streams/${encodeURIComponent(streamId)}/mjpeg`
}
