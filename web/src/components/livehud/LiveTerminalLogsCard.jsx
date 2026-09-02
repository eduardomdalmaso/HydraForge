import React, { useRef, useEffect } from 'react';
import { formatLoss } from '../../utils/formatters';

export default function LiveTerminalLogsCard({ job, rawLogs, recentMetrics, gpuStats }) {
  const terminalRef = useRef(null);
  const gpuName = gpuStats?.model || 'NVIDIA GPU';
  const gpuVramGB = Math.round((gpuStats?.total_vram_mb || 32607) / 1024);

  const getLogs = () => {
    if (rawLogs && rawLogs.length > 0) {
      return rawLogs.map(text => {
        let type = 'normal';
        if (text.startsWith('[HydraForge]') || text.includes('TRAINING')) type = 'active';
        else if (text.startsWith('[Epoch') || text.includes('Loss:')) type = 'metric';
        else if (text.includes('SUCCESS') || text.includes('saved')) type = 'success';
        return { text, type };
      });
    }

    if (!job) {
      return [
        { text: '[HydraForge Control Plane] Standby listener initialized on :8081', type: 'active' },
        { text: `[Hardware Telemetry] ${gpuName} (${gpuVramGB}GB VRAM) detected & ready for CUDA workloads`, type: 'metric' },
        { text: '[Status] Awaiting job launch from Training Cockpit...', type: 'normal' }
      ];
    }

    const modelName = job.model_architecture || job.model_name || 'yolo26m';
    const epochs = job.hyperparameters?.epochs || job.total_epochs || 50;
    const lines = [
      { text: `[HydraForge] Training Job ${job.job_id} (${modelName}) started on ${job.dataset_id}`, type: 'active' },
      { text: `[Configuration] Epochs: ${epochs} | Batch: ${job.hyperparameters?.batch_size || 32} | ImgSz: ${job.hyperparameters?.imgsz || 640}`, type: 'normal' }
    ];

    if (recentMetrics && recentMetrics.length > 0) {
      recentMetrics.slice(-6).forEach(m => {
        lines.push({
          text: `[Epoch ${m.epoch}/${epochs}] Box Loss: ${formatLoss(m.box_loss)} | Cls: ${formatLoss(m.cls_loss)} | mAP50: ${((m.map50 || 0) * 100).toFixed(1)}% | ${Math.round(m.fps || 0)} FPS`,
          type: 'metric'
        });
      });
    }
    return lines;
  };

  const logs = getLogs();

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [logs.length]);

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">4. STREAMING PYTORCH LOGS</span>
        <span className="badge-cyan">{job ? job.status : 'IDLE'}</span>
      </div>

      <div ref={terminalRef} className="terminal-box" style={{ marginTop: '0.5rem', maxHeight: '180px', overflowY: 'auto' }}>
        {logs.map((log, i) => {
          let colorClass = '';
          if (log.type === 'active') colorClass = 'terminal-row-active';
          else if (log.type === 'metric') colorClass = 'terminal-row-metric';
          else if (log.type === 'success') colorClass = 'terminal-row-success';
          return (
            <div key={i} className={colorClass} style={{ marginBottom: '0.2rem', fontFamily: 'var(--font-mono)', fontSize: '0.75rem' }}>
              {log.text}
            </div>
          );
        })}
      </div>
    </div>
  );
}
