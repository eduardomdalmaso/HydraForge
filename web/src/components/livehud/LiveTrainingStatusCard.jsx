import React from 'react';

export default function LiveTrainingStatusCard({ job, onAbortJob, onRestartJob, onResumeJob }) {
  if (!job) {
    return (
      <div className="cyber-card">
        <div className="card-header">
          <span className="card-title">1. ACTIVE TRAINING STATUS</span>
          <span className="badge-cyan">STANDBY / IDLE</span>
        </div>
        <div style={{ textAlign: 'center', padding: '2.5rem 1rem', color: '#94a3b8' }}>
          <div style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.95rem', color: '#fff' }}>NO ACTIVE TRAINING IN PROGRESS</div>
          <div style={{ fontSize: '0.75rem', marginTop: '0.35rem' }}>
            Navigate to <strong>TRAINING COCKPIT</strong> to launch a session on PyTorch / RTX 5090.
          </div>
          <button className="cyber-action-btn" style={{ marginTop: '1rem' }} onClick={() => window.location.hash = 'cockpit'}>
            COCKPIT
          </button>
        </div>
      </div>
    );
  }

  const currentEpoch = job.current_epoch || 0;
  const totalEpochs = job.hyperparameters?.epochs || job.total_epochs || job.epochs || 50;
  const currentBatch = job.current_batch || 0;
  const totalBatches = job.total_batches || 0;
  const isTraining = job.status === 'TRAINING' || job.status === 'RUNNING';

  const epochBatchPct = totalBatches > 0 ? Math.min(100, Math.round((currentBatch / totalBatches) * 100)) : (job.status === 'COMPLETED' ? 100 : 0);
  const baseEpoch = currentEpoch > 0 ? (currentEpoch - 1) : 0;
  const batchFraction = totalBatches > 0 ? (currentBatch / totalBatches) : 0;
  const globalPct = totalEpochs > 0 ? (job.status === 'COMPLETED' ? 100 : Math.min(100, ((baseEpoch + batchFraction) / totalEpochs) * 100)) : 0;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. ACTIVE TRAINING STATUS</span>
        <span className={isTraining ? 'badge-green' : 'badge-cyan'}>● {job.status}</span>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
        <div>
          <span style={{ fontFamily: 'var(--font-oxanium)', fontSize: '1.05rem', fontWeight: '700', color: '#fff' }}>
            {job.model_architecture || job.model_name || 'YOLO26'}
          </span>
          <span style={{ marginLeft: '0.5rem', fontSize: '0.72rem', color: '#94a3b8' }}>// DATASET: {job.dataset_id}</span>
        </div>
        {job.hyperparameters?.two_stage && <div className="stage-pill stage1">STAGE 1 (FROZEN)</div>}
      </div>

      <div className="selector-group" style={{ marginBottom: '0.65rem' }}>
        <div className="selector-label" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span>CURRENT EPOCH PROGRESS</span>
          <span className="slider-val" style={{ color: 'var(--cb-cyan)' }}>
            {totalBatches > 0 ? `BATCH ${currentBatch} / ${totalBatches} (${epochBatchPct}%)` : `EPOCH ${currentEpoch > 0 ? currentEpoch : 1}`}
          </span>
        </div>
        <div className="epoch-progress-track">
          <div className="epoch-progress-fill" style={{ width: `${epochBatchPct}%`, transition: 'width 0.3s ease' }} />
        </div>
      </div>

      <div className="selector-group" style={{ marginBottom: '0.65rem' }}>
        <div className="selector-label" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span>TOTAL SESSION COMPLETION</span>
          <span className="slider-val" style={{ color: 'var(--cb-yellow)' }}>
            EPOCH {currentEpoch > 0 ? currentEpoch : 1} / {totalEpochs} ({globalPct.toFixed(1)}%)
          </span>
        </div>
        <div className="epoch-progress-track">
          <div className="epoch-progress-fill" style={{ width: `${globalPct}%`, background: 'linear-gradient(90deg, var(--cb-yellow) 0%, var(--cb-magenta) 100%)', boxShadow: '0 0 12px rgba(252, 238, 10, 0.4)', transition: 'width 0.3s ease' }} />
        </div>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '0.65rem', flexWrap: 'wrap', gap: '0.4rem' }}>
        <span style={{ fontSize: '0.72rem', color: '#94a3b8', fontFamily: 'var(--font-mono)' }}>JOB ID: {job.job_id}</span>
        <div style={{ display: 'flex', gap: '0.4rem' }}>
          {isTraining ? (
            <>
              <button className="cyber-action-btn secondary" style={{ padding: '0.35rem 0.75rem', fontSize: '0.7rem' }} onClick={() => onRestartJob?.(job)}>RESTART</button>
              <button className="cyber-action-btn danger" style={{ padding: '0.35rem 0.75rem', fontSize: '0.7rem' }} onClick={onAbortJob}>ABORT</button>
            </>
          ) : (
            <>
              <button className="cyber-action-btn" style={{ padding: '0.35rem 0.75rem', fontSize: '0.7rem' }} onClick={() => onRestartJob?.(job)}>RESTART</button>
              {job.current_epoch > 0 && job.status !== 'COMPLETED' && (
                <button className="cyber-action-btn secondary" style={{ padding: '0.35rem 0.75rem', fontSize: '0.7rem' }} onClick={() => onResumeJob?.(job)}>RESUME</button>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
