/* ==========================================================================
   REAL-TIME INFERENCE API CLIENT (< 100 LINES)
   Communicates directly with NVIDIA RTX 5090 and Ultralytics YOLO26
   ========================================================================== */

export async function runRealInferenceAPI(config, imageBase64 = null) {
  try {
    const isWebcam = config.source === 'webcam';
    let sourcePath = config.source || 'cam_entrance_01';
    let displayImageUrl = null;

    if (isWebcam) {
      sourcePath = '/dev/shm/hydra_webcam_frame.jpg';
    } else {
      displayImageUrl = `/api/v1/hydrastream/api/v1/streams/${encodeURIComponent(sourcePath)}/snapshot.jpg?t=${Date.now()}`;
    }

    const payload = {
      model: config.model || 'yolo26n',
      source: sourcePath,
      image_base64: imageBase64,
      conf: config.conf || 0.25,
      iou: config.iou || 0.45,
      device: '0'
    };

    const res = await fetch('/api/v1/inference/predict', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (!res.ok) throw new Error(`HTTP error ${res.status}`);
    const data = await res.json();
    return { ...data, displayImageUrl };
  } catch (err) {
    console.error('Real inference API error:', err);
    return null;
  }
}
