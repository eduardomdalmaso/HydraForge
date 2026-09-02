import React from 'react';

export default function HardwareEstimatorCard({
  datasets = [],
  selectedDataset,
  setSelectedDataset,
  scale,
  batchSize,
  imgsz
}) {
  const scaleMultiplier = scale === 'n' ? 1.0 : scale === 's' ? 1.8 : scale === 'm' ? 3.5 : scale === 'l' ? 5.5 : 8.0;
  const batchMultiplier = batchSize === -1 ? 16 : batchSize;
  const resMultiplier = (imgsz / 640) * (imgsz / 640);
  const estimatedVRAM = Math.min(31.5, (0.8 + (scaleMultiplier * batchMultiplier * 0.04 * resMultiplier))).toFixed(1);

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. HARDWARE & ESTIMATOR</span>
        <span className="badge-green">NVIDIA RTX 5090</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">TARGET DATASET</div>
        <select className="cyber-select" value={selectedDataset || ''} onChange={(e) => setSelectedDataset(e.target.value)}>
          {datasets.length === 0 ? (
            <option value="">No dataset imported (import via DATASET STUDIO)</option>
          ) : (
            datasets.map((d) => (
              <option key={d.dataset_id} value={d.dataset_id}>
                {d.name} ({d.num_classes} classes • {d.train_images} train / {d.val_images} val)
              </option>
            ))
          )}
        </select>
      </div>

      <div className="telemetry-box" style={{ background: 'rgba(0,0,0,0.4)', padding: '1rem', borderRadius: '4px', border: '1px solid rgba(0, 240, 255, 0.15)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
          <span style={{ fontSize: '0.8rem', color: '#94a3b8' }}>ESTIMATED VRAM PEAK:</span>
          <span style={{ fontFamily: 'var(--font-mono)', color: 'var(--cb-yellow)', fontWeight: '700' }}>
            ~{estimatedVRAM} GB / 32.0 GB
          </span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
          <span style={{ fontSize: '0.8rem', color: '#94a3b8' }}>HARDWARE DEVICE:</span>
          <span style={{ fontFamily: 'var(--font-mono)', color: 'var(--cb-green)', fontWeight: '700' }}>
            cuda:0 (NVIDIA RTX 5090)
          </span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span style={{ fontSize: '0.8rem', color: '#94a3b8' }}>ZERO-COPY EXPORT:</span>
          <span style={{ fontFamily: 'var(--font-mono)', color: 'var(--cb-cyan)', fontWeight: '700' }}>
            TensorRT 10.x Engine
          </span>
        </div>
      </div>
    </div>
  );
}
