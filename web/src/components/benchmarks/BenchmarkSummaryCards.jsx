import React from 'react';

export default function BenchmarkSummaryCards({ results = [] }) {
  const valid = results?.filter(r => r.status === 'SUCCESS') || [];
  const hasResults = valid.length > 0;

  const bestFPS = hasResults ? valid.reduce((max, r) => (r.fps > max ? r.fps : max), 0) : null;
  const bestFormat = hasResults ? (valid.find(r => r.fps === bestFPS)?.format || '-') : 'AWAITING BENCHMARK RUN';
  const minLatency = hasResults ? valid.reduce((min, r) => (r.inference_time_ms < min ? r.inference_time_ms : min), 9999) : null;

  const ptResult = valid.find(r => r.format?.toLowerCase().includes('pytorch'));
  const ptFPS = ptResult?.fps || (valid[0]?.fps || null);
  const speedup = (hasResults && ptFPS > 0) ? (bestFPS / ptFPS).toFixed(1) : null;
  const bestMAP = hasResults ? valid.reduce((max, r) => (r.map50_95 > max ? r.map50_95 : max), 0) : null;

  return (
    <div className="telemetry-grid" style={{ marginBottom: '0.5rem' }}>
      <div className="cyber-card metric-card">
        <div className="metric-title">MAX THROUGHPUT</div>
        <div className="metric-value" style={{ color: 'var(--cb-cyan)' }}>
          {bestFPS ? `${bestFPS.toFixed(0)} FPS` : '-'}
        </div>
        <div className="metric-subtitle">{bestFormat}</div>
      </div>

      <div className="cyber-card metric-card">
        <div className="metric-title">MIN LATENCY (GPU)</div>
        <div className="metric-value" style={{ color: 'var(--cb-green)' }}>
          {minLatency ? `${minLatency.toFixed(2)} ms` : '-'}
        </div>
        <div className="metric-subtitle">Per Frame Inference (RTX 5090)</div>
      </div>

      <div className="cyber-card metric-card">
        <div className="metric-title">SPEEDUP MULTIPLIER</div>
        <div className="metric-value" style={{ color: 'var(--cb-yellow)' }}>
          {speedup ? `${speedup}x` : '-'}
        </div>
        <div className="metric-subtitle">vs PyTorch Native</div>
      </div>

      <div className="cyber-card metric-card">
        <div className="metric-title">PRECISION RETENTION</div>
        <div className="metric-value" style={{ color: 'var(--cb-magenta)' }}>
          {bestMAP ? `${(bestMAP * 100).toFixed(1)}%` : '-'}
        </div>
        <div className="metric-subtitle">mAP@50-95 Retention</div>
      </div>
    </div>
  );
}
