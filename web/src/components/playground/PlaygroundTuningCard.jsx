import React from 'react';

export default function PlaygroundTuningCard({ config, setConfig }) {
  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">2. HYPER-TUNING & FILTERS</span>
        <span className="badge-yellow">NMS & SAHI</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">
          <span>CONFIDENCE THRESHOLD</span>
          <span className="slider-val">{(config.conf * 100).toFixed(0)}%</span>
        </div>
        <div className="slider-row">
          <input
            type="range"
            min="0.05"
            max="0.95"
            step="0.05"
            value={config.conf}
            onChange={(e) => setConfig({ ...config, conf: parseFloat(e.target.value) })}
          />
        </div>
      </div>

      <div className="selector-group">
        <div className="selector-label">
          <span>IOU NMS THRESHOLD</span>
          <span className="slider-val">{(config.iou * 100).toFixed(0)}%</span>
        </div>
        <div className="slider-row">
          <input
            type="range"
            min="0.10"
            max="0.90"
            step="0.05"
            value={config.iou}
            onChange={(e) => setConfig({ ...config, iou: parseFloat(e.target.value) })}
            disabled={config.nmsFree}
          />
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.4rem' }}>
          <input
            type="checkbox"
            id="nmsFreeCheck"
            checked={config.nmsFree}
            onChange={(e) => setConfig({ ...config, nmsFree: e.target.checked })}
          />
          <label htmlFor="nmsFreeCheck" style={{ fontSize: '0.75rem', color: '#cbd5e1', cursor: 'pointer' }}>
            NMS-FREE END-TO-END OUTPUT (YOLO26)
          </label>
        </div>
      </div>

      <div className="selector-group" style={{ borderTop: '1px solid rgba(0,240,255,0.15)', paddingTop: '0.75rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.85rem', color: config.sahi ? 'var(--cb-cyan)' : '#fff', fontWeight: '700' }}>
              SAHI 4K SLICED INFERENCE
            </div>
            <div style={{ fontSize: '0.7rem', color: '#94a3b8' }}>Detect tiny objects in ultra-HD</div>
          </div>
          <button
            type="button"
            className={`cyber-pill ${config.sahi ? 'active' : ''}`}
            onClick={() => setConfig({ ...config, sahi: !config.sahi })}
          >
            {config.sahi ? 'ON' : 'OFF'}
          </button>
        </div>
      </div>

      <div className="selector-group" style={{ borderTop: '1px solid rgba(0,240,255,0.15)', paddingTop: '0.75rem', marginBottom: 0 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.85rem', color: config.isolateBg ? 'var(--cb-green)' : '#fff', fontWeight: '700' }}>
              ALPHA BACKGROUND REMOVER (PNG)
            </div>
            <div style={{ fontSize: '0.7rem', color: '#94a3b8' }}>Extract transparent alpha masks</div>
          </div>
          <button
            type="button"
            className={`cyber-pill ${config.isolateBg ? 'active' : ''}`}
            onClick={() => setConfig({ ...config, isolateBg: !config.isolateBg })}
          >
            {config.isolateBg ? 'ON' : 'OFF'}
          </button>
        </div>
      </div>
    </div>
  );
}
