import React from 'react';

export default function BenchmarkThroughputChart({ results = [] }) {
  const valid = results?.filter(r => r.status === 'SUCCESS') || [];
  const maxFPS = valid.length > 0 ? Math.max(...valid.map(r => r.fps), 100) : 100;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">THROUGHPUT COMPARISON (FPS)</span>
        <span className="badge-cyan">{valid.length > 0 ? 'REAL BENCHMARK DATA' : 'STANDBY'}</span>
      </div>
      <div style={{ marginTop: '0.5rem' }}>
        {valid.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '2rem 1rem', color: '#94a3b8', fontSize: '0.8rem' }}>
            Nenhum benchmark executado ainda. Selecione os runtimes e execute o teste.
          </div>
        ) : (
          valid.map((r, idx) => {
            const pct = Math.min(100, Math.max(6, (r.fps / maxFPS) * 100));
            const isTensorRT = r.format?.toLowerCase().includes('tensorrt') || r.format?.toLowerCase().includes('engine');
            const isONNX = r.format?.toLowerCase().includes('onnx');

            let fillColor = 'var(--cb-yellow)';
            if (isTensorRT) fillColor = 'var(--cb-green)';
            else if (isONNX) fillColor = 'var(--cb-cyan)';

            return (
              <div key={idx} className="bar-chart-row">
                <div className="bar-chart-label">{r.format}</div>
                <div className="bar-chart-track">
                  <div className="bar-chart-fill" style={{ width: `${pct}%`, background: fillColor }} />
                </div>
                <div className="bar-chart-val" style={{ color: fillColor }}>
                  {r.fps?.toFixed(0) || 0} <span style={{ fontSize: '0.65rem' }}>FPS</span>
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}
