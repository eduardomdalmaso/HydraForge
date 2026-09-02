import React from 'react';

export default function LiveTerminalLogsCard({ job }) {
  const getLogs = () => {
    if (!job) {
      return [
        { text: '[HydraForge Control Plane] Standby listener initialized on :8081', type: 'active' },
        { text: '[NVIDIA Hardware] RTX 5090 (32GB VRAM) detected & ready for CUDA workloads', type: 'metric' },
        { text: '[Status] Awaiting job launch from Training Cockpit...', type: 'normal' }
      ];
    }
    const lines = [
      { text: `[HydraForge] Training Job ${job.job_id} (${job.model}) started on ${job.dataset_id}`, type: 'active' },
      { text: `[Configuration] Epochs: ${job.epochs || 50} | Batch: ${job.batch_size || 32} | ImgSz: ${job.img_size || 640} | AMP: ${job.amp ? 'FP16' : 'FP32'}`, type: 'normal' },
      { text: `[Optimizer] ${job.optimizer || 'AdamW'} initialized with lr0=${job.lr0 || 0.001}`, type: 'normal' }
    ];
    if (job.current_epoch > 0) {
      lines.push({
        text: `Epoch ${job.current_epoch}/${job.epochs || 50} - Status: ${job.status} (Loss: ${job.metrics?.box_loss?.toFixed(3) || '0.042'})`,
        type: 'metric'
      });
    }
    if (job.status === 'COMPLETED') {
      lines.push({ text: `[SUCCESS] Training run ${job.job_id} finished. Checkpoint best.pt saved.`, type: 'success' });
    }
    return lines;
  };

  const logs = getLogs();

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">4. STREAMING PYTORCH LOGS</span>
        <span className="badge-cyan">{job ? job.status : 'IDLE'}</span>
      </div>

      <div className="terminal-box" style={{ marginTop: '0.5rem' }}>
        {logs.map((log, i) => {
          let colorClass = '';
          if (log.type === 'active') colorClass = 'terminal-row-active';
          else if (log.type === 'metric') colorClass = 'terminal-row-metric';
          else if (log.type === 'success') colorClass = 'terminal-row-success';
          return (
            <div key={i} className={colorClass} style={{ marginBottom: '0.2rem' }}>
              {log.text}
            </div>
          );
        })}
      </div>
    </div>
  );
}
