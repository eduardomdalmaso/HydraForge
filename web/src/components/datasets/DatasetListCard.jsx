import React from 'react';
import DatasetCyberIcon from './DatasetCyberIcons';

export default function DatasetListCard({
  datasets = [],
  selectedId,
  onSelectDataset,
  onOpenImportModal,
  onOpenAnnotateModal,
  onOpenMergeModal,
  onRescan,
  onDeleteDataset
}) {
  const isMerged = (ds) => ds.dataset_id.includes('merged') || ds.dataset_id.includes('frota') || ds.dataset_id.includes('fusion');
  const mergedCount = datasets.filter(isMerged).length;
  const rawCount = datasets.length - mergedCount;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. REGISTERED DATASETS</span>
        <div style={{ display: 'flex', gap: '0.35rem', alignItems: 'center' }}>
          {onRescan && (
            <button className="cyber-pill" onClick={onRescan} title="Rescan local folders" style={{ fontSize: '0.62rem', padding: '0.15rem 0.4rem' }}>
              🔄 RESCAN
            </button>
          )}
          <span className="badge-cyan">{rawCount} REPOS</span>
        </div>
      </div>

      <div className="ds-list-scroll">
        {datasets.length === 0 ? (
          <div style={{ fontSize: '0.8rem', color: '#94a3b8', textAlign: 'center', padding: '1.2rem 0' }}>
            No registered datasets found. Click <strong>IMPORT</strong> or <strong>RESCAN</strong>.
          </div>
        ) : (
          datasets.map((ds) => {
            const isSelected = selectedId === ds.dataset_id;
            const merged = isMerged(ds);
            const totalImgs = (ds.train_images || 0) + (ds.val_images || 0);
            return (
              <div key={ds.dataset_id} className={`dataset-card-item ${isSelected ? 'active' : ''}`} onClick={() => onSelectDataset(ds)}>
                <div style={{ minWidth: 0, flex: 1, paddingRight: '0.5rem' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
                    <span style={{ color: isSelected ? 'var(--cb-yellow)' : (merged ? 'var(--cb-green)' : 'var(--cb-cyan)'), fontSize: '0.8rem' }}>
                      {merged ? '⚡' : (isSelected ? '▶' : '◈')}
                    </span>
                    <span className="ds-name">{ds.name}</span>
                    <span className={merged ? "badge-green" : "badge-cyan"} style={{ fontSize: '0.55rem', padding: '0.05rem 0.25rem' }}>
                      {merged ? 'MERGED' : 'RAW'}
                    </span>
                  </div>
                  <div className="ds-meta">
                    <span className="ds-tag">📁 {totalImgs.toLocaleString()} imgs</span>
                    <span className="ds-tag">🏷️ {ds.num_classes} cls</span>
                    <span className="ds-tag">⚡ train: {ds.train_images || 0} / val: {ds.val_images || 0}</span>
                  </div>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
                  <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '0.2rem' }}>
                    <span className="badge-yellow" style={{ fontSize: '0.62rem', padding: '0.1rem 0.35rem' }}>{ds.task || 'DETECT'}</span>
                    <span style={{ fontSize: '0.62rem', fontFamily: 'var(--font-mono)', color: isSelected ? 'var(--cb-yellow)' : '#64748b' }}>{isSelected ? 'ACTIVE' : 'SELECT'}</span>
                  </div>
                  {onDeleteDataset && (
                    <button type="button" className="cyber-pill" title="Delete dataset" style={{ padding: '0.2rem 0.4rem', fontSize: '0.7rem', color: 'var(--cb-magenta)', borderColor: 'rgba(255,0,60,0.3)' }} onClick={(e) => { e.stopPropagation(); onDeleteDataset(ds); }}>
                      🗑️
                    </button>
                  )}
                </div>
              </div>
            );
          })
        )}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '0.4rem', marginTop: '0.5rem' }}>
        <button className="cyber-action-btn secondary" onClick={onOpenAnnotateModal} disabled={datasets.length === 0}>
          <DatasetCyberIcon name="tag" size={12} /> <span>ANNOTATE</span>
        </button>
        <button className="cyber-action-btn secondary" onClick={onOpenMergeModal}>
          <DatasetCyberIcon name="fusion" size={12} /> <span>MERGE</span>
        </button>
        <button className="cyber-action-btn" onClick={onOpenImportModal}>
          <DatasetCyberIcon name="zip" size={12} color="#07080c" /> <span>IMPORT</span>
        </button>
      </div>
    </div>
  );
}
