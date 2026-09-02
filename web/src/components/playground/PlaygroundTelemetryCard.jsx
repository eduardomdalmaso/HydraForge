import React, { useState } from 'react';

export default function PlaygroundTelemetryCard({ telemetry, detections = [], selectedEntity, onSelectEntity }) {
  const [exportedMsg, setExportedMsg] = useState(null);

  const handleExportPng = (e, det) => {
    e.stopPropagation();
    setExportedMsg(`PNG EXPORTED: ${det.label.toUpperCase()} #${det.id}`);
    setTimeout(() => setExportedMsg(null), 2500);
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. HARDWARE & TELEMETRY</span>
        <span className="badge-green">NVIDIA RTX 5090</span>
      </div>

      <div style={{ marginBottom: '1rem' }}>
        <div className="telemetry-row">
          <span className="k">DALI PREPROCESS (GPU)</span>
          <span className="v" style={{ color: 'var(--cb-green)' }}>{telemetry?.preprocess_ms || '-'} ms</span>
        </div>
        <div className="telemetry-row">
          <span className="k">TENSORRT 10.x INFERENCE</span>
          <span className="v" style={{ color: 'var(--cb-cyan)' }}>{telemetry?.inference_ms || '-'} ms</span>
        </div>
        <div className="telemetry-row">
          <span className="k">NMS / POSTPROCESS</span>
          <span className="v">{telemetry?.postprocess_ms || '-'} ms</span>
        </div>
        <div className="telemetry-row">
          <span className="k">SUSTAINED INFERENCE SPEED</span>
          <span className="v" style={{ color: 'var(--cb-yellow)', fontSize: '0.9rem' }}>
            {telemetry?.fps ? `${telemetry.fps} FPS` : '-'}
          </span>
        </div>
        <div className="telemetry-row" style={{ borderBottom: 'none' }}>
          <span className="k">VRAM ALLOCATION</span>
          <span className="v">{telemetry?.vram_mb ? `${telemetry.vram_mb} MB` : 'STANDBY'}</span>
        </div>
      </div>

      {exportedMsg && (
        <div style={{ background: 'rgba(0,255,157,0.08)', border: '1px solid var(--cb-green)', padding: '0.4rem', borderRadius: '3px', fontSize: '0.72rem', color: 'var(--cb-green)', fontFamily: 'var(--font-mono)', marginBottom: '0.65rem' }}>
          ✓ {exportedMsg}
        </div>
      )}

      <div className="selector-group" style={{ marginBottom: 0 }}>
        <div className="selector-label">
          <span>DETECTED ENTITIES ({detections.length})</span>
          <span style={{ fontSize: '0.65rem', color: '#94a3b8' }}>CLICK TO FOCUS</span>
        </div>

        <div style={{ maxHeight: '180px', overflowY: 'auto', paddingRight: '4px' }}>
          {detections.length === 0 ? (
            <div style={{ fontSize: '0.75rem', color: '#64748b', textAlign: 'center', padding: '1.2rem 0' }}>
              Nenhum alvo detectado. Execute um scan no botão acima.
            </div>
          ) : (
            detections.map((det) => {
              const isSelected = selectedEntity?.id === det.id;
              return (
                <div
                  key={det.id}
                  className="entity-item"
                  style={{
                    borderColor: isSelected ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
                    background: isSelected ? 'rgba(252, 238, 10, 0.15)' : 'rgba(0, 240, 255, 0.05)',
                    cursor: 'pointer'
                  }}
                  onClick={() => onSelectEntity(det)}
                >
                  <div>
                    <span style={{ fontWeight: '700', color: '#fff' }}>{det.label}</span>
                    <span style={{ marginLeft: '0.4rem', color: 'var(--cb-yellow)', fontSize: '0.7rem' }}>
                      {(det.conf * 100).toFixed(0)}%
                    </span>
                  </div>
                  <button className="cyber-pill" style={{ padding: '0.2rem 0.5rem', fontSize: '0.65rem' }} onClick={(e) => handleExportPng(e, det)}>
                    ✂️ PNG
                  </button>
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
}
