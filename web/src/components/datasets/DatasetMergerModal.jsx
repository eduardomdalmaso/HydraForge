import React, { useState, useEffect } from 'react';
import DatasetCyberIcon from './DatasetCyberIcons';
import { autoSuggestCategory } from './yolo_coco_classes';

export default function DatasetMergerModal({ isOpen, onClose, datasets = [], onDatasetsMerged, savedMappings = {} }) {
  const [mergedName, setMergedName] = useState('fusion_custom_dataset');
  const [isMerging, setIsMerging] = useState(false);
  const [selectedDsIds, setSelectedDsIds] = useState([]);
  const [mappings, setMappings] = useState({});

  const getSuggestedName = (ids) => {
    const bases = Array.from(new Set(ids.map(id => id.replace(/\.v\d+.*$/, '').replace(/_v\d+.*$/, '').replace(/[-_]raw$/, ''))));
    return bases.length === 1 ? `${bases[0]}_merged_${ids.length}ds` : (bases.length ? `${bases.slice(0, 2).join('_')}_fusion` : 'fusion_custom_dataset');
  };

  useEffect(() => {
    if (!isOpen) return;
    const initial = {};
    datasets.forEach(ds => {
      initial[ds.dataset_id] = {};
      (ds.classes || []).forEach(cls => { initial[ds.dataset_id][cls] = savedMappings[ds.dataset_id]?.[cls] || autoSuggestCategory(cls); });
    });
    setMappings(initial); setSelectedDsIds([]); setMergedName('fusion_custom_dataset');
  }, [isOpen]);

  if (!isOpen) return null;

  const toggle = (id) => setSelectedDsIds(p => { const next = p.includes(id) ? p.filter(x => x !== id) : [...p, id]; setMergedName(getSuggestedName(next)); return next; });
  const setMap = (dsId, cls, val) => setMappings(p => ({ ...p, [dsId]: { ...(p[dsId] || {}), [cls]: val } }));
  const targetOptions = Array.from(new Set([...datasets.flatMap(d => d.classes || []), ...Object.values(mappings).flatMap(m => Object.values(m)), 'car', 'motorcycle', 'truck', 'bus', 'cell-phone', 'ignore'])).filter(Boolean);
  const totalImgs = datasets.filter(d => selectedDsIds.includes(d.dataset_id)).reduce((acc, d) => acc + (d.train_images || 0) + (d.val_images || 0) + (d.test_images || 0), 0);

  const handleExecute = async () => {
    setIsMerging(true);
    try {
      const activeTargets = Array.from(new Set(selectedDsIds.flatMap(id => Object.values(mappings[id] || {})))).filter(c => c && c !== 'ignore');
      const payload = { target_name: mergedName, dataset_ids: selectedDsIds, mappings, classes: activeTargets.length ? activeTargets : ['cell-phone'] };
      const res = await fetch('/api/v1/training/datasets/merge', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
      if (res.ok && onDatasetsMerged) { onDatasetsMerged(await res.json()); window.location.reload(); }
    } finally { setIsMerging(false); onClose(); }
  };

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="import-modal-box" style={{ maxWidth: '720px' }} onClick={e => e.stopPropagation()}>
        <div className="card-header" style={{ marginBottom: '0.5rem' }}>
          <span className="card-title">DATASET HARMONIZER & FUSION ENGINE (MERGE)</span>
          <button className="cyber-pill" style={{ padding: '0.2rem 0.5rem' }} onClick={onClose}>✕</button>
        </div>

        <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', marginBottom: '0.6rem' }}>
          <span style={{ fontSize: '0.72rem', color: '#94a3b8', fontFamily: 'var(--font-mono)' }}>NAME:</span>
          <input type="text" value={mergedName} onChange={e => setMergedName(e.target.value)} style={{ flex: 1, background: '#05070a', border: '1px solid rgba(0,240,255,0.3)', color: '#fff', padding: '0.25rem 0.5rem', fontSize: '0.75rem', borderRadius: '3px', fontFamily: 'var(--font-mono)' }} />
          <button type="button" className="cyber-pill" style={{ padding: '0.15rem 0.4rem', fontSize: '0.62rem' }} onClick={() => { const all = datasets.map(d => d.dataset_id); setSelectedDsIds(all); setMergedName(getSuggestedName(all)); }}>ALL</button>
          <button type="button" className="cyber-pill" style={{ padding: '0.15rem 0.4rem', fontSize: '0.62rem' }} onClick={() => { setSelectedDsIds([]); setMergedName('fusion_custom_dataset'); }}>NONE</button>
        </div>

        <div style={{ maxHeight: '240px', overflowY: 'auto', paddingRight: '4px', marginBottom: '0.75rem' }}>
          {datasets.map(ds => {
            const inc = selectedDsIds.includes(ds.dataset_id);
            const imgs = (ds.train_images || 0) + (ds.val_images || 0) + (ds.test_images || 0);
            return (
              <div key={ds.dataset_id} style={{ background: inc ? 'rgba(0,240,255,0.05)' : 'rgba(0,0,0,0.2)', border: inc ? '1px solid rgba(0,240,255,0.3)' : '1px solid rgba(255,255,255,0.08)', padding: '0.5rem', borderRadius: '4px', marginBottom: '0.5rem', opacity: inc ? 1 : 0.5 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.35rem' }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', cursor: 'pointer' }}>
                    <input type="checkbox" checked={inc} onChange={() => toggle(ds.dataset_id)} />
                    <span style={{ fontSize: '0.78rem', fontFamily: 'var(--font-oxanium)', color: inc ? 'var(--cb-yellow)' : '#94a3b8', fontWeight: '700' }}>📁 {ds.name} ({imgs.toLocaleString()} imgs)</span>
                  </label>
                  <span style={{ fontSize: '0.65rem', color: '#64748b' }}>{inc ? 'INCLUDED' : 'EXCLUDED'}</span>
                </div>
                {inc && (
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: '0.35rem' }}>
                    {(ds.classes || []).map(cls => (
                      <div key={cls} className="class-mapper-row" style={{ height: '28px', padding: '0.15rem 0.4rem' }}>
                        <span className="class-mapper-label" style={{ fontSize: '0.65rem', maxWidth: '110px' }}>"{cls}" ➔</span>
                        <select className="cyber-select" style={{ width: '100px', padding: '0.1rem 0.2rem', fontSize: '0.65rem' }} value={mappings[ds.dataset_id]?.[cls] || ''} onChange={e => setMap(ds.dataset_id, cls, e.target.value)}>
                          {targetOptions.map(t => <option key={t} value={t}>{t}</option>)}
                        </select>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            );
          })}
        </div>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span style={{ fontSize: '0.72rem', color: 'var(--cb-cyan)', fontFamily: 'var(--font-mono)' }}>{selectedDsIds.length} repos selected • {totalImgs.toLocaleString()} images</span>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <button className="cyber-action-btn secondary" onClick={onClose} disabled={isMerging}>CANCEL</button>
            <button className="cyber-action-btn" onClick={handleExecute} disabled={isMerging || selectedDsIds.length === 0}><DatasetCyberIcon name="fusion" size={14} color="#07080c" /> <span>{isMerging ? 'MERGING...' : 'FUSE DATASETS'}</span></button>
          </div>
        </div>
      </div>
    </div>
  );
}
