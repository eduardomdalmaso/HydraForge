import React from 'react';

export default function LiveGpuTelemetryCard({ gpuStats, job }) {
  const vramUsed = Math.round(gpuStats?.used_vram_mb || 0);
  const vramTotal = Math.round(gpuStats?.total_vram_mb || 32607);
  const vramPct = Math.round(gpuStats?.vram_usage_pct || (vramTotal > 0 ? (vramUsed / vramTotal) * 100 : 0));
  const temp = Math.round(gpuStats?.temp_celsius || 0);
  const power = Math.round(gpuStats?.power_watts || 0);
  const util = Math.round(gpuStats?.gpu_util_pct || 0);
  const energyKWh = job?.total_energy_kwh || 0;
  const avgFPS = Math.round(job?.avg_fps || 0);
  const durationSec = Math.round(job?.duration_sec || 0);

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. HARDWARE & NVIDIA RTX 5090</span>
        <span className="badge-green">NVML LIVE TELEMETRY</span>
      </div>

      <div style={{ marginBottom: '0.65rem' }}>
        <div className="selector-label" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span>VRAM USAGE // REAL HARDWARE</span>
          <span className="slider-val">{vramUsed.toLocaleString()} MB / {vramTotal.toLocaleString()} MB ({vramPct}%)</span>
        </div>
        <div className="epoch-progress-track">
          <div className="epoch-progress-fill" style={{ width: `${Math.max(vramPct, 2)}%`, background: 'var(--cb-cyan)' }} />
        </div>
      </div>

      <div className="telemetry-row">
        <span className="k">ACTIVE ACCELERATOR</span>
        <span className="v" style={{ color: 'var(--cb-cyan)' }}>{gpuStats?.model || 'NVIDIA GeForce RTX 5090'}</span>
      </div>
      <div className="telemetry-row">
        <span className="k">GPU TEMP & POWER DRAW</span>
        <span className="v" style={{ color: temp < 65 ? 'var(--cb-green)' : 'var(--cb-yellow)' }}>{temp}°C • {power} W / 600 W</span>
      </div>
      <div className="telemetry-row">
        <span className="k">TOTAL ENERGY BUDGET (SESSION)</span>
        <span className="v" style={{ color: 'var(--cb-yellow)' }}>
          {energyKWh > 0 ? `${energyKWh.toFixed(4)} kWh` : '0.0000 kWh'} {job?.avg_power_watts ? `(Avg ${Math.round(job.avg_power_watts)}W)` : ''}
        </span>
      </div>
      <div className="telemetry-row">
        <span className="k">TRAIN THROUGHPUT & TIME</span>
        <span className="v" style={{ color: 'var(--cb-green)' }}>
          {avgFPS > 0 ? `${avgFPS.toLocaleString()} FPS` : '-'} • {durationSec > 0 ? `${Math.floor(durationSec / 60)}m ${durationSec % 60}s` : '0s'}
        </span>
      </div>
      <div className="telemetry-row" style={{ borderBottom: 'none' }}>
        <span className="k">SM & TENSOR CORE UTILIZATION</span>
        <span className="v" style={{ color: util > 30 ? 'var(--cb-yellow)' : 'var(--cb-green)' }}>{util}% (Blackwell FP16 TC)</span>
      </div>
    </div>
  );
}
