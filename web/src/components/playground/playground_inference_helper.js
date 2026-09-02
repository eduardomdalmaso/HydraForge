/* ==========================================================================
   PLAYGROUND INFERENCE & TARGET GENERATOR HELPER (< 100 LINES)
   ========================================================================== */

export const SAMPLE_CLASSES = [
  'Person', 'Vehicle', 'Cyber-Bicycle', 'Traffic Light',
  'Drone', 'License Plate', 'Optical Target', 'Security Robot'
];

export function generatePlaygroundDetections(config, streamId) {
  const seed = (streamId || 'stream').length;
  const isSeg = config.model.includes('seg');
  const isPose = config.model.includes('pose');

  const baseItems = [
    { id: 1, label: isSeg ? 'Mask: Vehicle' : (isPose ? 'Pose: Security Agent' : 'Vehicle'), conf: 0.94, box: [18, 28, 36, 42], color: '#00f0ff' },
    { id: 2, label: isPose ? 'Pose: Pedestrian' : 'Pedestrian (Person)', conf: 0.89, box: [62, 34, 16, 45], color: '#fcee0a' },
    { id: 3, label: 'License Plate (OCR)', conf: 0.96, box: [32, 58, 12, 10], color: '#00ff66' },
    { id: 4, label: 'Optical Sensor / Drone', conf: 0.82, box: [75, 15, 14, 18], color: '#ff0055' }
  ];

  if (config.isolateBg) {
    return [{ id: 1, label: 'Foreground Segment Mask (SAM 2)', conf: 0.98, box: [15, 20, 70, 65], color: '#00f0ff' }];
  }

  return baseItems.filter((item, idx) => {
    const jitter = ((seed + idx * 7) % 15) / 100;
    item.conf = Math.min(0.99, Math.max(0.3, item.conf - jitter));
    return item.conf >= config.conf;
  });
}

export function calculateInferenceTelemetry(config) {
  const isTRT = config.runtime === 'tensorrt';
  const baseInf = isTRT ? 0.36 : 1.15;
  const jitter = (Math.random() * 0.06).toFixed(2);
  const infMs = (baseInf + parseFloat(jitter)).toFixed(2);
  const fps = isTRT ? '2,631' : '869';

  return {
    preprocess_ms: isTRT ? '0.01' : '0.12',
    inference_ms: infMs,
    postprocess_ms: config.nmsFree ? '0.00 (NMS-Free)' : '0.03',
    fps: fps,
    vram_mb: isTRT ? '1,420' : '2,180'
  };
}
