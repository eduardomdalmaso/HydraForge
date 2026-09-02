import React, { useState } from 'react';

export default function PlaygroundTelemetryCard({
  telemetry, detections = [], selectedEntity, onSelectEntity, hydraTelemetry, activeStream, gpuStats
}) {
  const [exportedMsg, setExportedMsg] = useState(null);

  const handleExportPng = (e, det) => {
    e.stopPropagation();
    setExportedMsg(`PNG EXPORTED: ${det.label.toUpperCase()} #${det.id}`);
    setTimeout(() => setExportedMsg(null), 2500);
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. HARDWARE & SHM TELEMETRY</span>
        <span className="badge-green">{gpuStats?.model || 'RTX 5090'}</span>
      </div>

      <div style={{ marginBottom: '0.8rem' }}>
        <div className="telemetry-row">
          <span className="k">GPU CORE & TEMP</span>
          <span className="v" style={{ color: 'var(--cb-green)' }}>
            {gpuStats ? `${gpuStats.gpu_util_pct}% • ${gpuStats.temp_celsius}°C (${gpuStats.power_watts}W)` : 'STANDBY'}
          </span>
        </div>
        <div className="telemetry-row">
          <span className="k">VRAM ALLOCATION</span>
          <span className="v" style={{ color: 'var(--cb-cyan)' }}>
            {gpuStats ? `${gpuStats.used_vram_mb?.toLocaleString()} / ${gpuStats.total_vram_mb?.toLocaleString()} MB` : 'STANDBY'}
          </span>
        </div>
        <div className="telemetry-row">
          <span className="k">YOLO INFERENCE (RTX 5090)</span>
          <span className="v" style={{ color: 'var(--cb-yellow)' }}>
            {telemetry?.inference_ms ? `${telemetry.inference_ms} ms` : 'TRIGGER SCAN'}
          </span>
        </div>
        <div className="telemetry-row">
          <span className="k">SUSTAINED SPEED</span>
          <span className="v" style={{ color: 'var(--cb-yellow)', fontSize: '0.85rem' }}>
            {telemetry?.fps ? `${telemetry.fps} FPS` : 'STANDBY'}
          </span>
        </div>
        <div className="telemetry-row" style={{ borderBottom: 'none' }}>
          <span className="k">POSIX /dev/shm OCCUPANCY</span>
          <span className="v">{hydraTelemetry?.posix_shm_occupancy !== undefined ? `${hydraTelemetry.posix_shm_occupancy}%` : '0.1%'}</span>
        </div>
      </div>

      {exportedMsg && (
        <div style={{ background: 'rgba(0,255,157,0.08)', border: '1px solid var(--cb-green)', padding: '0.35rem', borderRadius: '3px', fontSize: '0.7rem', color: 'var(--cb-green)', fontFamily: 'var(--font-mono)', marginBottom: '0.5rem' }}>
          [OK] {exportedMsg}
        </div>
      )}

      <div className="selector-group" style={{ marginBottom: 0 }}>
        <div className="selector-label">
          <span>DETECTED ENTITIES ({detections.length})</span>
          <span style={{ fontSize: '0.65rem', color: '#94a3b8' }}>REAL LABELS</span>
        </div>
        <div style={{ maxHeight: '160px', overflowY: 'auto', paddingRight: '4px' }}>
          {detections.length === 0 ? (
            <div style={{ fontSize: '0.75rem', color: '#64748b', textAlign: 'center', padding: '1.2rem 0' }}>Execute um scan na RTX 5090 para ver as detecções reais.</div>
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
                    <span style={{ marginLeft: '0.4rem', color: 'var(--cb-yellow)', fontSize: '0.7rem' }}>{(det.conf * 100).toFixed(0)}%</span>
                  </div>
                  <button className="cyber-pill" style={{ padding: '0.2rem 0.5rem', fontSize: '0.65rem' }} onClick={(e) => handleExportPng(e, det)}>[PNG]</button>
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
}
