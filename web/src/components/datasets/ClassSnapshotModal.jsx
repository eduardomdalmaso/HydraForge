import React, { useState, useEffect } from 'react';
import { YOLO_CLASSES, decodeYoloClass } from './yolo_coco_classes';
import DatasetCyberIcon from './DatasetCyberIcons';

export default function ClassSnapshotModal({ isOpen, onClose, datasetId, className, currentTarget, onSaveTarget }) {
  const [target, setTarget] = useState(currentTarget || 'carro');
  const [search, setSearch] = useState('');
  const [sample, setSample] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const filtered = YOLO_CLASSES.filter(c => c.label.toLowerCase().includes(search.toLowerCase()) || c.en.toLowerCase().includes(search.toLowerCase()));

  const fetchSample = async () => {
    if (!datasetId || !className) return;
    setIsLoading(true);
    try {
      const res = await fetch(`/api/v1/training/datasets/sample?id=${datasetId}&class=${encodeURIComponent(className)}&t=${Date.now()}`);
      if (res.ok) setSample(await res.json());
    } finally { setIsLoading(false); }
  };

  useEffect(() => {
    if (isOpen && className) {
      setTarget(currentTarget || 'carro');
      setSearch('');
      fetchSample();
    }
  }, [isOpen, className, currentTarget]);

  if (!isOpen) return null;

  const bx = sample?.bbox ? (sample.bbox[0] - sample.bbox[2] / 2) * 1000 : 0;
  const by = sample?.bbox ? (sample.bbox[1] - sample.bbox[3] / 2) * 1000 : 0;
  const bw = sample?.bbox ? sample.bbox[2] * 1000 : 0;
  const bh = sample?.bbox ? sample.bbox[3] * 1000 : 0;

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="import-modal-box" style={{ maxWidth: '680px' }} onClick={(e) => e.stopPropagation()}>
        <div className="card-header" style={{ marginBottom: '0.65rem' }}>
          <span className="card-title">SNAPSHOT: "{className}" ({decodeYoloClass(className)})</span>
          <button className="cyber-pill" style={{ padding: '0.2rem 0.5rem' }} onClick={onClose}>✕</button>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1.2fr 1fr', gap: '0.85rem' }}>
          <div>
            <div style={{ width: '100%', height: '215px', background: '#05070a', border: '1px solid rgba(0,240,255,0.3)', borderRadius: '4px', display: 'flex', alignItems: 'center', justifyContent: 'center', overflow: 'hidden' }}>
              {isLoading ? (
                <div style={{ color: 'var(--cb-yellow)', fontFamily: 'var(--font-mono)', fontSize: '0.75rem' }}>CARREGANDO...</div>
              ) : sample?.image_url ? (
                <svg viewBox="0 0 1000 1000" preserveAspectRatio="xMidYMid meet" style={{ width: '100%', height: '100%', display: 'block' }}>
                  <image href={sample.image_url} x="0" y="0" width="1000" height="1000" preserveAspectRatio="none" />
                  {sample.bbox && (
                    <>
                      <rect x={bx} y={by} width={bw} height={bh} fill="rgba(0, 240, 255, 0.2)" stroke="#00f0ff" strokeWidth="8" />
                      <text x={bx + 10} y={Math.max(by - 15, 35)} fill="#00f0ff" fontSize="42" fontFamily="monospace" fontWeight="bold">{className}</text>
                    </>
                  )}
                </svg>
              ) : <div style={{ color: '#64748b', fontSize: '0.75rem' }}>Sem preview.</div>}
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.4rem', marginTop: '0.4rem' }}>
              <button className="cyber-action-btn secondary" style={{ padding: '0.45rem' }} onClick={fetchSample} disabled={isLoading}>NEXT</button>
              <button className="cyber-action-btn danger" style={{ padding: '0.45rem' }} onClick={() => { onSaveTarget(className, 'ignorar'); onClose(); }}>IGNORE</button>
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
            <div>
              <div className="selector-label" style={{ marginBottom: '0.3rem' }}>ESCOLHA A CLASSE DESTINO:</div>
              <input type="text" placeholder="🔍 Filtrar..." value={search} onChange={e => setSearch(e.target.value)} style={{ width: '100%', background: '#05070a', border: '1px solid rgba(0,240,255,0.25)', color: '#fff', padding: '0.25rem 0.5rem', fontSize: '0.7rem', borderRadius: '3px', marginBottom: '0.35rem', fontFamily: 'var(--font-mono)' }} />
              <div style={{ height: '135px', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '0.25rem', paddingRight: '3px' }}>
                {filtered.map(c => {
                  const isSel = target === c.id;
                  return (
                    <div key={c.id} className="class-mapper-row" style={{ cursor: 'pointer', height: '28px', padding: '0.2rem 0.5rem', borderColor: isSel ? 'var(--cb-yellow)' : 'rgba(0,240,255,0.15)', background: isSel ? 'rgba(252,238,10,0.15)' : 'rgba(0,0,0,0.4)' }} onClick={() => setTarget(c.id)}>
                      <span style={{ color: isSel ? 'var(--cb-yellow)' : '#fff', fontSize: '0.7rem', fontFamily: 'var(--font-mono)', fontWeight: isSel ? '700' : '400' }}>
                        {isSel ? '✓ ' : '◇ '}{c.label} <span style={{ color: isSel ? '#fef08a' : '#64748b', fontSize: '0.62rem' }}>[{c.en}]</span>
                      </span>
                      {isSel && <span className="badge-yellow" style={{ fontSize: '0.55rem', padding: '0.05rem 0.3rem' }}>SELECTED</span>}
                    </div>
                  );
                })}
              </div>
            </div>

            <button className="cyber-action-btn" style={{ width: '100%', padding: '0.6rem', marginTop: '0.5rem' }} onClick={() => { onSaveTarget(className, target); onClose(); }}>
              <DatasetCyberIcon name="launch" size={14} color="#07080c" />
              <span>SAVE</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
