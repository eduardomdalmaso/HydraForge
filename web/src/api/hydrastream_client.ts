/* ==========================================================================
   HYDRASTREAM INTEGRATION API CLIENT MODULE (TypeScript)
   ========================================================================== */

const HYDRASTREAM_BASE = '/api/v1/hydrastream'

export async function fetchHydraStreams(): Promise<any[]> {
  try {
    const res = await fetch(`${HYDRASTREAM_BASE}/api/v1/streams`)
    if (!res.ok) return []
    const data = await res.json()
    return (data && Array.isArray(data.streams)) ? data.streams : []
  } catch {
    return []
  }
}

export async function fetchHydraTelemetry(): Promise<any | null> {
  try {
    const res = await fetch(`${HYDRASTREAM_BASE}/api/v1/telemetry/stats`)
    if (!res.ok) return null
    return await res.json()
  } catch {
    return null
  }
}

export function getStreamSnapshotUrl(streamId?: string): string {
  if (!streamId) return ''
  return `${HYDRASTREAM_BASE}/api/v1/streams/${encodeURIComponent(streamId)}/snapshot.jpg?t=${Date.now()}`
}

export function getStreamMJPEGUrl(streamId?: string): string {
  if (!streamId) return ''
  return `${HYDRASTREAM_BASE}/api/v1/streams/${encodeURIComponent(streamId)}/mjpeg`
}
