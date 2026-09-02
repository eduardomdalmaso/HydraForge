/* ==========================================================================
   HYDRAFORGE BENCHMARK API CLIENT MODULE
   ========================================================================== */

export async function fetchBenchmarksAPI(status = '') {
  try {
    const res = await fetch(`/api/v1/benchmarks?status=${status}`);
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return [];
  }
}

export async function fetchBenchmarkByIDAPI(jobId) {
  try {
    const res = await fetch(`/api/v1/benchmarks/${jobId}`);
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return null;
  }
}

export async function launchBenchmarkAPI(config) {
  const res = await fetch('/api/v1/benchmarks', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config)
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || `HTTP error ${res.status}`);
  }
  return await res.json();
}

export async function stopBenchmarkAPI(jobId) {
  try {
    const res = await fetch(`/api/v1/benchmarks/${jobId}`, { method: 'DELETE' });
    return res.ok;
  } catch (err) {
    return false;
  }
}

export async function fetchBenchmarkFormatsAPI() {
  try {
    const res = await fetch('/api/v1/benchmarks/formats');
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return [];
  }
}
