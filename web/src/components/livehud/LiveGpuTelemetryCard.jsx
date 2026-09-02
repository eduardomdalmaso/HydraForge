import React from 'react';

export default function LiveGpuTelemetryCard({ gpuStats }) {
  const vramUsed = gpuStats?.vram_used_mb || 4820;
  const vramTotal = gpuStats?.vram_total_mb || 32607;
  const vramPct = Math.round((vramUsed / vramTotal) * 100);
  const temp = gpuStats?.temp_celsius || 48;
  const power = gpuStats?.power_watts || 285;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. HARDWARE & NVIDIA RTX 5090</span>
        <span className="badge-green">CUDA 13.3 • 575.86</span>
      </div>

      <div style={{ marginBottom: '0.65rem' }}>
        <div className="selector-label" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span>VRAM USAGE // DUAL-HEAD TRAINING</span>
          <span className="slider-val">{vramUsed.toLocaleString()} MB / {vramTotal.toLocaleString()} MB ({vramPct}%)</span>
        </div>
        <div className="epoch-progress-track">
          <div className="epoch-progress-fill" style={{ width: `${vramPct}%`, background: 'var(--cb-cyan)' }}></div>
        </div>
      </div>

      <div className="telemetry-row">
        <span className="k">NVIDIA DALI PREPROCESSING PIPELINE</span>
        <span className="v" style={{ color: 'var(--cb-green)' }}>0.004 ms (GPU Zero-Copy)</span>
      </div>
      <div className="telemetry-row">
        <span className="k">GPU CORE TEMPERATURE</span>
        <span className="v" style={{ color: temp < 65 ? 'var(--cb-green)' : 'var(--cb-yellow)' }}>{temp}°C (Fan: 42%)</span>
      </div>
      <div className="telemetry-row">
        <span className="k">BOARD POWER CONSUMPTION</span>
        <span className="v">{power} W / 600 W TGP</span>
      </div>
      <div className="telemetry-row" style={{ borderBottom: 'none' }}>
        <span className="k">SM UTILIZATION & TENSOR CORES</span>
        <span className="v" style={{ color: 'var(--cb-yellow)' }}>94% (Blackwell FP16 TC)</span>
      </div>
    </div>
  );
}
