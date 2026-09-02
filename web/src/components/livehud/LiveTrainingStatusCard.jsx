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
          <button className="cyber-pill active" style={{ marginTop: '1rem', padding: '0.4rem 1rem' }} onClick={() => window.location.hash = 'cockpit'}>
            IR PARA O COCKPIT →
          </button>
        </div>
      </div>
    );
  }

  const currentEpoch = job.current_epoch || 0;
  const totalEpochs = job.epochs || 50;
  const progressPct = totalEpochs > 0 ? Math.min(100, Math.round((currentEpoch / totalEpochs) * 100)) : 0;
  const isStage2 = currentEpoch > (job.stage1_epochs || 20);

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. ACTIVE TRAINING STATUS</span>
        <span className={job.status === 'RUNNING' ? 'badge-green' : 'badge-cyan'}>
          ● {job.status}
        </span>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
        <div>
          <span style={{ fontFamily: 'var(--font-oxanium)', fontSize: '1.05rem', fontWeight: '700', color: '#fff' }}>
            {job.model_name || job.model || 'YOLO26'}
          </span>
          <span style={{ marginLeft: '0.5rem', fontSize: '0.72rem', color: '#94a3b8' }}>
            // DATASET: {job.dataset_id}
          </span>
        </div>
        {job.two_stage && (
          <div className={`stage-pill ${isStage2 ? 'stage2' : 'stage1'}`}>
            {isStage2 ? '⚡ STAGE 2: ALL LAYERS' : `🔒 STAGE 1: FROZEN (0-${job.stage1_freeze || 10})`}
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
          <div className="epoch-progress-fill" style={{ width: `${progressPct}%` }}></div>
        </div>
      </div>

      <div className="telemetry-row">
        <span className="k">OPTIMIZER & PRECISION</span>
        <span className="v" style={{ color: 'var(--cb-cyan)' }}>{job.optimizer || 'AdamW'} (AMP FP16: {job.amp ? 'ON' : 'OFF'})</span>
      </div>
      <div className="telemetry-row">
        <span className="k">JOB IDENTIFIER</span>
        <span className="v text-mono">{job.job_id}</span>
      </div>
      <div className="telemetry-row" style={{ borderBottom: 'none' }}>
        <span className="k">CREATED AT</span>
        <span className="v text-mono">{job.created_at ? new Date(job.created_at).toLocaleTimeString() : '-'}</span>
      </div>

      {job.status === 'RUNNING' && (
        <button className="cyber-pill" style={{ width: '100%', marginTop: '0.75rem', borderColor: 'var(--cb-magenta)', color: 'var(--cb-magenta)', textAlign: 'center', padding: '0.5rem' }} onClick={onAbortJob}>
          ⏹ ABORT RUNNING JOB
        </button>
      )}
    </div>
  );
}
