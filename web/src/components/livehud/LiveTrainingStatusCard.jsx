import React from 'react';

export default function LiveTrainingStatusCard({ job, onAbortJob }) {
  if (!job) {
    return (
      <div className="cyber-card">
        <div className="card-header">
          <span className="card-title">1. ACTIVE TRAINING STATUS</span>
          <span className="badge-cyan">STANDBY / IDLE</span>
        </div>
        <div style={{ textAlign: 'center', padding: '2.5rem 1rem', color: '#94a3b8' }}>
          <div style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>⏸</div>
          <div style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.95rem', color: '#fff' }}>NENHUM TREINO EM EXECUÇÃO</div>
          <div style={{ fontSize: '0.75rem', marginTop: '0.35rem' }}>
            Acesse o <strong>TRAINING COCKPIT</strong> para iniciar uma sessão no PyTorch / RTX 5090.
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
  const progressPct = totalEpochs > 0 ? Math.min(100, Math.round((currentEpoch / totalEpochs) * 100)) : 0;
  const isTraining = job.status === 'TRAINING' || job.status === 'RUNNING';

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. ACTIVE TRAINING STATUS</span>
        <span className={isTraining ? 'badge-green' : 'badge-cyan'}>
          ● {job.status}
        </span>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
        <div>
          <span style={{ fontFamily: 'var(--font-oxanium)', fontSize: '1.05rem', fontWeight: '700', color: '#fff' }}>
            {job.model_architecture || job.model_name || 'YOLO26'}
          </span>
          <span style={{ marginLeft: '0.5rem', fontSize: '0.72rem', color: '#94a3b8' }}>
            // DATASET: {job.dataset_id}
          </span>
        </div>
        {job.hyperparameters?.two_stage && (
          <div className="stage-pill stage1">
            STAGE 1 (FROZEN)
          </div>
        )}
      </div>

      <div className="selector-group" style={{ marginBottom: '0.65rem' }}>
        <div className="selector-label" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span>EPOCH PROGRESSION</span>
          <span className="slider-val" style={{ color: 'var(--cb-yellow)' }}>
            EPOCH {currentEpoch} / {totalEpochs} ({progressPct}%)
          </span>
        </div>
        <div className="epoch-progress-track">
          <div className="epoch-progress-fill" style={{ width: `${progressPct}%` }} />
        </div>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '0.75rem' }}>
        <span style={{ fontSize: '0.72rem', color: '#94a3b8', fontFamily: 'var(--font-mono)' }}>
          JOB ID: {job.job_id}
        </span>
        {isTraining && (
          <button className="cyber-action-btn danger" style={{ padding: '0.35rem 0.8rem', fontSize: '0.7rem' }} onClick={onAbortJob}>
            ABORT
          </button>
        )}
      </div>
    </div>
  );
}
