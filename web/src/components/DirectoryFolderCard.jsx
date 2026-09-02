import React from 'react';
import CyberNavIcon from './CyberNavIcons';

export default function DirectoryFolderCard({ folder, onOpen }) {
  const { id, label, code, tag, desc, stats, status } = folder;

  return (
    <div className="cyber-folder-card" onClick={() => onOpen(id)}>
      <div className="folder-tab-ear">
        // {code} [{tag}]
      </div>
      
      <div className="folder-header">
        <div className="folder-title-group">
          <div className="folder-icon cyber-icon-glow">
            <CyberNavIcon type={id} />
          </div>
          <div>
            <div className="folder-title">{label}</div>
            <div className="text-mono" style={{ fontSize: '0.7rem', color: 'var(--cb-yellow)' }}>
              PATH: /{id}
            </div>
          </div>
        </div>
        <span className="card-badge">{status || 'READY'}</span>
      </div>

      <div className="folder-desc">
        {desc}
      </div>

      <div className="folder-footer">
        <span style={{ color: '#64748b' }}>
          METRIC: <strong style={{ color: 'var(--cb-green)' }}>{stats}</strong>
        </span>
        <span className="folder-action-btn">
          MOUNT DIRECTORY ❯❯
        </span>
      </div>
    </div>
  );
}
