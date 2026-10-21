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

export async function registerHydraStream(streamId: string, sourceUrl: string, ingestFps = 25): Promise<any> {
  const res = await fetch(`${HYDRASTREAM_BASE}/api/v1/streams`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      stream_id: streamId,
      source_url: sourceUrl,
      ingest_fps: ingestFps
    })
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  return await res.json()
}

