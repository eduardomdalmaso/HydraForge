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
          <span>TRAINING EPOCHS & EARLY STOPPING</span>
          <span className="slider-val">{params.epochs} EP (Patience: {params.patience || 20})</span>
        </div>
        <div className="slider-row">
          <input type="range" min="5" max="300" value={params.epochs} onChange={(e) => handleChange('epochs', parseInt(e.target.value))} />
        </div>
      </div>

      <div className="selector-group">
        <div className="selector-label">BATCH SIZE & RESOLUTION</div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.65rem' }}>
          <div>
            <label style={{ fontSize: '0.72rem', color: '#94a3b8' }}>BATCH SIZE</label>
            <select className="cyber-select" value={params.batch_size} onChange={(e) => handleChange('batch_size', parseInt(e.target.value))}>
              <option value="-1">Auto-Batch (-1)</option>
              <option value="16">16 Images</option>
              <option value="32">32 Images</option>
              <option value="64">64 Images (RTX 5090)</option>
            </select>
          </div>
          <div>
            <label style={{ fontSize: '0.72rem', color: '#94a3b8' }}>IMAGE SIZE (IMGSZ)</label>
            <select className="cyber-select" value={params.imgsz} onChange={(e) => handleChange('imgsz', parseInt(e.target.value))}>
              <option value="640">640 x 640 (Standard)</option>
              <option value="1280">1280 x 1280 (HD)</option>
            </select>
          </div>
        </div>
      </div>

      <div className="selector-group" style={{ marginBottom: 0 }}>
        <div className="selector-label">OPTIMIZER & EARLY STOPPING (PATIENCE)</div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.65rem' }}>
          <div>
            <label style={{ fontSize: '0.72rem', color: '#94a3b8' }}>OPTIMIZER</label>
            <select className="cyber-select" value={params.optimizer} onChange={(e) => handleChange('optimizer', e.target.value)}>
              <option value="AdamW">AdamW (Recommended)</option>
              <option value="SGD">SGD with Momentum</option>
            </select>
          </div>
          <div>
            <label style={{ fontSize: '0.72rem', color: '#94a3b8' }}>EARLY STOPPING</label>
            <select className="cyber-select" value={params.patience || 20} onChange={(e) => handleChange('patience', parseInt(e.target.value))}>
              <option value="15">15 Epochs</option>
              <option value="20">20 Epochs</option>
              <option value="30">30 Epochs</option>
              <option value="50">50 Epochs</option>
              <option value="0">Desativado (0)</option>
            </select>
          </div>
        </div>
      </div>
    </div>
  );
}
