import React from 'react';

export default function Header({ gpuTelemetry, activeJobsCount }) {
  const vramUsed = gpuTelemetry ? (gpuTelemetry.used_vram_mb / 1024).toFixed(1) : '2.0';
  const vramTotal = gpuTelemetry ? (gpuTelemetry.total_vram_mb / 1024).toFixed(0) : '32';
  const temp = gpuTelemetry ? gpuTelemetry.temp_celsius.toFixed(0) : '37';

  return (
    <header className="top-header">
      <div className="header-breadcrumbs">
        <span className="crumb-root">CONTROL PLANE</span>
        <span className="crumb-sep">/</span>
        <span className="crumb-active">TRAINING COCKPIT</span>
      </div>

      <div className="header-status-group">
        <div className="status-badge online">
          <span className="pulse-dot"></span>
          <span>STUDIO READY</span>
        </div>

        <div className="header-metric-pill">
          <span className="pill-label">ACTIVE RUNS:</span>
          <span className="pill-val cyan">{activeJobsCount || 0}</span>
        </div>

        <div className="header-metric-pill">
          <span className="pill-label">VRAM:</span>
          <span className="pill-val yellow">{vramUsed} / {vramTotal} GB</span>
          <span className="pill-sub">({temp}°C)</span>
        </div>

        <a href="/swagger/" target="_blank" rel="noreferrer" className="swagger-link-btn">
          DOCS
        </a>
      </div>
    </header>
  );
}
