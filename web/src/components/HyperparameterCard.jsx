import React from 'react';

export default function HyperparameterCard({ params, setParams }) {
  const handleChange = (field, val) => {
    setParams((prev) => ({ ...prev, [field]: val }));
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">2. HYPERPARAMETER COCKPIT</span>
        <span className="badge-yellow">PYTORCH CUDA</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">
          <span>TRAINING EPOCHS</span>
          <span className="slider-val">{params.epochs}</span>
        </div>
        <div className="slider-row">
          <input type="range" min="1" max="500" value={params.epochs} onChange={(e) => handleChange('epochs', parseInt(e.target.value))} />
        </div>
      </div>

      <div className="selector-group">
        <div className="selector-label">BATCH SIZE & RESOLUTION</div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
          <div>
            <label style={{ fontSize: '0.75rem', color: '#94a3b8' }}>BATCH SIZE</label>
            <select className="cyber-select" value={params.batch_size} onChange={(e) => handleChange('batch_size', parseInt(e.target.value))}>
              <option value="-1">Auto-Batch (-1)</option>
              <option value="8">8 Images</option>
              <option value="16">16 Images</option>
              <option value="32">32 Images</option>
              <option value="64">64 Images</option>
            </select>
          </div>
          <div>
            <label style={{ fontSize: '0.75rem', color: '#94a3b8' }}>IMAGE SIZE (IMGSZ)</label>
            <select className="cyber-select" value={params.imgsz} onChange={(e) => handleChange('imgsz', parseInt(e.target.value))}>
              <option value="640">640 x 640 (Fast)</option>
              <option value="1280">1280 x 1280 (HD)</option>
              <option value="1920">1920 x 1920 (4K)</option>
            </select>
          </div>
        </div>
      </div>

      <div className="selector-group" style={{ marginBottom: 0 }}>
        <div className="selector-label">OPTIMIZER & PRECISION</div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
          <select className="cyber-select" value={params.optimizer} onChange={(e) => handleChange('optimizer', e.target.value)}>
            <option value="AdamW">AdamW (Recommended)</option>
            <option value="SGD">SGD with Momentum</option>
            <option value="RMSprop">RMSprop</option>
          </select>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.2rem' }}>
            <input type="checkbox" id="ampCheck" checked={params.amp} onChange={(e) => handleChange('amp', e.target.checked)} />
            <label htmlFor="ampCheck" style={{ fontSize: '0.85rem', color: '#cbd5e1', cursor: 'pointer' }}>
              AMP FP16 (RTX 5090)
            </label>
          </div>
        </div>
      </div>
    </div>
  );
}
