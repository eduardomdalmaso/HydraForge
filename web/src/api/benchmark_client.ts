/* ==========================================================================
   HYDRAFORGE BENCHMARK API CLIENT MODULE (TypeScript)
   ========================================================================== */

export async function fetchBenchmarksAPI(status = ''): Promise<any[]> {
  try {
    const res = await fetch(`/api/v1/benchmarks?status=${status}`)
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return []
  }
}

export async function fetchBenchmarkByIDAPI(jobId: string): Promise<any | null> {
  try {
    const res = await fetch(`/api/v1/benchmarks/${jobId}`)
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return null
  }
}

export async function launchBenchmarkAPI(config: Record<string, any>): Promise<any> {
  const res = await fetch('/api/v1/benchmarks', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config)
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error || `HTTP error ${res.status}`)
  }
  return await res.json()
}

export async function stopBenchmarkAPI(jobId: string): Promise<boolean> {
  try {
    const res = await fetch(`/api/v1/benchmarks/${jobId}`, { method: 'DELETE' })
    return res.ok
  } catch {
    return false
  }
}

export async function fetchBenchmarkFormatsAPI(): Promise<any[]> {
  try {
    const res = await fetch('/api/v1/benchmarks/formats')
    if (!res.ok) throw new Error(`HTTP error ${res.status}`)
    return await res.json()
  } catch {
    return []
  }
}
