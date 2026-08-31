import React from 'react';

const NAV_ITEMS = [
  { id: 'cockpit', label: 'Training Cockpit', icon: '⚡', tag: 'LAUNCHER' },
  { id: 'live-hud', label: 'Live Training HUD', icon: '📈', tag: 'WEBSOCKETS' },
  { id: 'datasets', label: 'Dataset Studio', icon: '🗂️', tag: 'DATA.YAML' },
  { id: 'model-zoo', label: 'Model Zoo & Export', icon: '📦', tag: 'TENSORRT' },
  { id: 'playground', label: 'Live Playground', icon: '🎯', tag: 'INFERENCE' }
];

export default function Sidebar({ activeTab, onSelectTab }) {
  return (
    <aside className="sidebar-nav">
      <div className="brand-header">
        <div className="brand-logo-icon">🔥</div>
        <div className="brand-text">
          <div className="brand-title">HYDRAFORGE</div>
          <div className="brand-subtitle">YOLO AI STUDIO</div>
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
              <span className="nav-icon">{item.icon}</span>
              <span className="nav-label">{item.label}</span>
              <span className="nav-tag">{item.tag}</span>
            </button>
          );
        })}
      </nav>

      <div className="sidebar-footer">
        <div className="gpu-badge">
          <span className="status-dot green"></span>
          <span>RTX 5090 • CUDA 13.3</span>
        </div>
        <div className="version-label">v1.0.0-PROD</div>
      </div>
    </aside>
  );
}
