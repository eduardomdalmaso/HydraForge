import React from 'react';

export default function LiveMetricsChartsCard({ job, recentMetrics }) {
  const latest = (recentMetrics && recentMetrics.length > 0) ? recentMetrics[recentMetrics.length - 1] : null;
  const map50 = latest?.map50 || job?.best_map50 || 0;
  const map5095 = latest?.map50_95 || 0;
  const boxLoss = latest?.box_loss || 0;
  const clsLoss = latest?.cls_loss || 0;
  const hasData = (recentMetrics && recentMetrics.length > 0) || (job && job.current_epoch > 0);

  const getPoints = (type) => {
    if (!recentMetrics || recentMetrics.length < 2) return null;
    const total = recentMetrics[0].total_epochs || 50;
    return recentMetrics.map((m) => {
      const x = Math.round((m.epoch / total) * 380 + 10);
      let y = 90;
      if (type === 'loss') {
        const lossVal = Math.min(3.0, (m.box_loss || 0) + (m.cls_loss || 0));
        y = Math.round(95 - (lossVal / 3.0) * 75);
      } else {
        const mapVal = Math.min(1.0, m.map50_95 || m.map50 || 0);
        y = Math.round(95 - mapVal * 80);
      }
      return `${x},${y}`;
    }).join(' ');
  };

  const lossPoints = getPoints('loss');
  const mapPoints = getPoints('map');

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">2. ACCURACY & LOSS DYNAMICS</span>
        <span className="badge-cyan">{hasData ? 'ACTIVE METRICS' : 'STANDBY'}</span>
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
          {hasData && lossPoints && mapPoints ? (
            <>
              <polyline fill="none" stroke="var(--cb-magenta)" strokeWidth="2" points={lossPoints} />
              <polyline fill="none" stroke="var(--cb-green)" strokeWidth="2" points={mapPoints} />
            </>
          ) : hasData ? (
            <text x="200" y="60" textAnchor="middle" fill="var(--cb-cyan)" fontFamily="var(--font-mono)" fontSize="11">
              Processando Épocas no PyTorch... acumulando pontos da curva
            </text>
          ) : (
            <text x="200" y="60" textAnchor="middle" fill="#64748b" fontFamily="var(--font-mono)" fontSize="11">
              {job ? 'PyTorch inicializado na GPU... aguardando Época 1' : 'Aguardando início do treinamento para desenhar curvas...'}
            </text>
          )}
        </svg>
      </div>
    </div>
  );
}
