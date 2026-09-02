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

export async function fetchModelsAPI() {
  try {
    const res = await fetch('/api/v1/training/models');
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return [];
  }
}

export async function rescanDatasetsAPI() {
  try {
    const res = await fetch('/api/v1/training/datasets/rescan', { method: 'POST' });
    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    return await res.json();
  } catch (err) {
    return [];
  }
}

export async function registerDatasetPathAPI(payload) {
  const res = await fetch('/api/v1/training/datasets/register-path', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || `HTTP ${res.status}`);
  }
  return await res.json();
}

export async function deleteDatasetAPI(datasetId, deleteFiles = false) {
  const res = await fetch(`/api/v1/training/datasets/${datasetId}?delete_files=${deleteFiles}`, { method: 'DELETE' });
  return res.ok;
}
