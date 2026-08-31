/* ==========================================================================
   HYDRAFORGE API CLIENT MODULE
   ========================================================================== */

export async function fetchTelemetryAPI() {
  try {
    const res = await fetch('/api/v1/training/telemetry');
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return null;
  }
}

export async function fetchDatasetsAPI() {
  try {
    const res = await fetch('/api/v1/training/datasets');
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return [];
  }
}

export async function fetchJobsAPI(status = '') {
  try {
    const res = await fetch(`/api/v1/training/jobs?status=${status}`);
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return [];
  }
}

export async function launchTrainingJobAPI(jobConfig) {
  try {
    const res = await fetch('/api/v1/training/jobs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(jobConfig)
    });
    if (!res.ok) {
      const err = await res.json();
      throw new Error(err.error || `HTTP error ${res.status}`);
    }
    return await res.json();
  } catch (err) {
    throw err;
  }
}

export async function stopTrainingJobAPI(jobId) {
  try {
    const res = await fetch(`/api/v1/training/jobs/${jobId}`, { method: 'DELETE' });
    return res.ok;
  } catch (err) {
    return false;
  }
}
