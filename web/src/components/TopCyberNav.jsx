import React from 'react';
import CyberNavIcon from './CyberNavIcons';
import HydraLogo from './HydraLogo';

const QUICK_TABS = [
  { id: 'cockpit', label: 'COCKPIT', tag: 'R-01' },
  { id: 'live-hud', label: 'LIVE HUD', tag: 'R-02' },
  { id: 'benchmarks', label: 'BENCHMARKS', tag: 'R-03' },
  { id: 'datasets', label: 'DATASETS', tag: 'R-04' },
  { id: 'model-zoo', label: 'MODEL ZOO', tag: 'R-05' },
  { id: 'playground', label: 'PLAYGROUND', tag: 'R-06' }
];

export default function TopCyberNav({ activeTab, onSelectTab, gpuStats }) {
  return (
    <header className="top-cyber-nav">
      <div className="nav-brand-btn" onClick={() => onSelectTab('cockpit')}>
        <HydraLogo size={34} />
        <div className="nav-brand-title-box">
          <span className="brand-hydra">HYDRA</span>
          <span className="brand-forge">FORGE</span>
        </div>
      </div>

      <div className="nav-folder-tabs">
        {QUICK_TABS.map(tab => {
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              className={`cybr-btn ${isActive ? 'cybr-btn--active' : ''}`}
              onClick={() => onSelectTab(tab.id)}
            >
              <span className="cybr-btn__content">
                <CyberNavIcon type={tab.id} />
                <span>{tab.label}</span>
                <span aria-hidden="true">_</span>
              </span>
              <span aria-hidden="true" className="cybr-btn__glitch">
                <CyberNavIcon type={tab.id} />
                <span>{tab.label}_</span>
              </span>
              <span aria-hidden="true" className="cybr-btn__tag">
                {tab.tag}
              </span>
            </button>
          );
        })}
      </div>

      <div className="top-nav-right">
        <div className="header-metric-pill">
          <span className="status-dot green"></span>
          <span className="pill-val cyan">{gpuStats?.model || 'RTX 5090'}</span>
          <span className="pill-sub">{(gpuStats?.temp_celsius || 36).toFixed(0)}°C</span>
        </div>
        <a href="/swagger/" target="_blank" rel="noreferrer" className="swagger-link-btn">
          [ API DOCS ]
        </a>
      </div>
    </header>
  );
}
