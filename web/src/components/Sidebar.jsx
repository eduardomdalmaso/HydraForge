import React from 'react';
import CyberNavIcon from './CyberNavIcons';

const NAV_ITEMS = [
  { id: 'cockpit', label: 'TRAINING COCKPIT', code: '01', tag: 'LAUNCH' },
  { id: 'live-hud', label: 'LIVE TELEMETRY HUD', code: '02', tag: 'STREAM' },
  { id: 'benchmarks', label: 'BENCHMARK STUDIO', code: '03', tag: 'SPEED' },
  { id: 'datasets', label: 'DATASET STUDIO', code: '04', tag: 'YAML' },
  { id: 'model-zoo', label: 'MODEL ZOO & EXPORT', code: '05', tag: 'TRT' },
  { id: 'playground', label: 'KIROSHI PLAYGROUND', code: '06', tag: 'OPTIC' }
];

export default function Sidebar({ activeTab, onSelectTab }) {
  return (
    <aside className="sidebar-nav">
      <div className="brand-header">
        <div className="brand-logo-icon">
          <svg width="28" height="28" viewBox="0 0 32 32" fill="none">
            <polygon points="16,2 30,10 30,22 16,30 2,22 2,10" stroke="var(--cb-magenta)" strokeWidth="2.5" fill="rgba(255,0,85,0.15)" />
            <polygon points="16,7 25,12 25,20 16,25 7,20 7,12" stroke="var(--cb-cyan)" strokeWidth="1.5" />
            <circle cx="16" cy="16" r="3" fill="var(--cb-yellow)" />
          </svg>
        </div>
        <div className="brand-text">
          <div className="brand-title">HYDRAFORGE</div>
          <div className="brand-subtitle">// YOLO.AI.STUDIO.v26</div>
        </div>
      </div>

      <nav className="nav-menu">
        {NAV_ITEMS.map((item) => {
          const isActive = activeTab === item.id;
          return (
            <button
              key={item.id}
              className={`nav-btn ${isActive ? 'active' : ''}`}
              onClick={() => onSelectTab(item.id)}
            >
              <span className="nav-icon cyber-icon-glow">
                <CyberNavIcon type={item.id} />
              </span>
              <div className="nav-label-group">
                <span className="nav-code">{item.code}</span>
                <span className="nav-label">{item.label}</span>
              </div>
              <span className="nav-tag">[{item.tag}]</span>
            </button>
          );
        })}
      </nav>

      <div className="sidebar-footer">
        <div className="gpu-badge">
          <span className="status-dot green"></span>
          <span>RTX 5090 // CUDA 13.3</span>
        </div>
        <div className="version-label text-mono">NET:ONLINE</div>
      </div>
    </aside>
  );
}
