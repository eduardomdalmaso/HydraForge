import React, { useState } from 'react';

export default function ExperimentsComparisonCard({ jobs = [], onTestInPlayground }) {
  const [selectedRunId, setSelectedRunId] = useState(jobs[0]?.job_id || null);
  const activeRun = jobs.find(r => r.job_id === selectedRunId) || jobs[0];

  return (
    <div className="cyber-card" style={{ marginTop: '1.25rem' }}>
      <div className="card-header">
        <span className="card-title">5. EXPERIMENT COMPARISON & METRICS VALIDATION</span>
        <span className="badge-yellow">{jobs.length} RUNS RECORDED</span>
      </div>

      {jobs.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '1.5rem', color: '#94a3b8', fontSize: '0.8rem' }}>
          No training history recorded in database. Runs launched from the Cockpit will appear here for cross-validation.
        </div>
      ) : (
        <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '0.65rem', marginBottom: '0.85rem' }}>
            {jobs.map(run => {
              const isSelected = (selectedRunId || jobs[0]?.job_id) === run.job_id;
              const map50 = run.best_map50 || run.metrics?.map50 || 0;
              const map5095 = run.best_map50_95 || run.metrics?.map50_95 || 0;
              return (
                <div
                  key={run.job_id}
                  className={`dataset-card-item ${isSelected ? 'active' : ''}`}
                  style={{ marginBottom: 0, padding: '0.65rem', cursor: 'pointer' }}
                  onClick={() => setSelectedRunId(run.job_id)}
                >
                  <div>
                    <div style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.85rem', fontWeight: '700', color: isSelected ? 'var(--cb-yellow)' : '#fff' }}>
                      {run.model_architecture || run.model} ({run.dataset_id})
                    </div>
                    <div style={{ fontFamily: 'var(--font-mono)', fontSize: '0.7rem', color: '#94a3b8', marginTop: '0.2rem' }}>
                      mAP50: <strong style={{ color: 'var(--cb-green)' }}>{map50 > 0 ? `${(map50 * 100).toFixed(1)}%` : '-'}</strong>
                      {map5095 > 0 && <span> • 50-95: <strong style={{ color: 'var(--cb-cyan)' }}>{(map5095 * 100).toFixed(1)}%</strong></span>}
                      <span> • {run.hyperparameters?.epochs || run.total_epochs || 50}e</span>
                    </div>
                  </div>
                  <span className={run.status === 'COMPLETED' ? 'badge-online' : 'badge-cyan'} style={{ fontSize: '0.65rem' }}>
                    {run.status}
                  </span>
                </div>
              );
            })}
          </div>

          {activeRun && (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '0.5rem', marginBottom: '0.75rem' }}>
              <div className="matrix-cell highlight">
                <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>mAP@50 {activeRun.best_map50_95 > 0 ? '/ @50-95' : ''}</div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '0.95rem', fontWeight: '700', color: 'var(--cb-green)' }}>
                  {activeRun.best_map50 > 0 ? `${(activeRun.best_map50 * 100).toFixed(1)}%` : 'PENDING'}
                  {activeRun.best_map50_95 > 0 ? ` / ${(activeRun.best_map50_95 * 100).toFixed(1)}%` : ''}
                </div>
              </div>
              <div className="matrix-cell">
                <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>OPTIMIZER / LR0</div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '0.85rem', fontWeight: '700', color: 'var(--cb-cyan)' }}>
                  {activeRun.hyperparameters?.optimizer || 'AdamW'} ({activeRun.hyperparameters?.lr0 || 0.001})
                </div>
              </div>
              <div className="matrix-cell">
                <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>RESOLUTION</div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1rem', fontWeight: '700', color: 'var(--cb-yellow)' }}>
                  {activeRun.hyperparameters?.imgsz || 640} px
                </div>
              </div>
              <div className="matrix-cell">
                <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>BATCH SIZE</div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1rem', fontWeight: '700', color: '#fff' }}>
                  {activeRun.hyperparameters?.batch_size || 32}
                </div>
              </div>
            </div>
          )}

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem', marginTop: '0.75rem' }}>
            <button className="cyber-action-btn" style={{ width: '100%', padding: '0.75rem' }} onClick={() => onTestInPlayground?.(activeRun)}>PLAYGROUND</button>
            <button className="cyber-action-btn secondary" style={{ width: '100%', padding: '0.75rem' }} onClick={() => window.location.hash = 'benchmarks'}>EXPORT</button>
          </div>
        </>
      )}
    </div>
  );
}
