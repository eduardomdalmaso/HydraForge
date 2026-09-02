import React from 'react';

export default function LiveMetricsChartsCard({ job }) {
  const map50 = job?.metrics?.map50 || 0;
  const map5095 = job?.metrics?.map50_95 || 0;
  const boxLoss = job?.metrics?.box_loss || 0;
  const clsLoss = job?.metrics?.cls_loss || 0;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">2. ACCURACY & LOSS DYNAMICS</span>
        <span className="badge-cyan">{job ? 'ACTIVE METRICS' : 'STANDBY'}</span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '0.5rem', marginBottom: '0.75rem' }}>
        <div className="matrix-cell highlight">
          <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>mAP@50</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1rem', fontWeight: '700', color: 'var(--cb-green)' }}>
            {map50 > 0 ? `${(map50 * 100).toFixed(1)}%` : '-'}
          </div>
        </div>
        <div className="matrix-cell highlight">
          <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>mAP@50-95</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1rem', fontWeight: '700', color: 'var(--cb-cyan)' }}>
            {map5095 > 0 ? `${(map5095 * 100).toFixed(1)}%` : '-'}
          </div>
        </div>
        <div className="matrix-cell">
          <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>BOX LOSS</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1rem', fontWeight: '700', color: 'var(--cb-yellow)' }}>
            {boxLoss > 0 ? boxLoss.toFixed(3) : '-'}
          </div>
        </div>
        <div className="matrix-cell">
          <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>CLS LOSS</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1rem', fontWeight: '700', color: 'var(--cb-magenta)' }}>
            {clsLoss > 0 ? clsLoss.toFixed(3) : '-'}
          </div>
        </div>
      </div>

      <div style={{ position: 'relative', background: 'rgba(0,0,0,0.4)', borderRadius: '4px', border: '1px solid rgba(0,240,255,0.15)', padding: '0.5rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.65rem', fontFamily: 'var(--font-mono)', color: '#94a3b8', marginBottom: '0.25rem' }}>
          <span><span style={{ color: 'var(--cb-green)' }}>●</span> mAP 50-95 Trajectory</span>
          <span><span style={{ color: 'var(--cb-magenta)' }}>●</span> Loss Convergence</span>
        </div>
        <svg className="hud-svg-chart" viewBox="0 0 400 110">
          <line x1="0" y1="30" x2="400" y2="30" className="hud-grid-line" />
          <line x1="0" y1="60" x2="400" y2="60" className="hud-grid-line" />
          <line x1="0" y1="90" x2="400" y2="90" className="hud-grid-line" />
          {job ? (
            <>
              <path d="M 10,95 Q 80,70 160,45 T 320,25 T 390,18" className="hud-chart-line-loss" />
              <path d="M 10,100 Q 80,80 160,50 T 280,32 T 390,20" className="hud-chart-line-map" />
            </>
          ) : (
            <text x="200" y="60" textAnchor="middle" fill="#64748b" fontFamily="var(--font-mono)" fontSize="11">
              Aguardando início do treinamento para desenhar curvas...
            </text>
          )}
        </svg>
      </div>
    </div>
  );
}
