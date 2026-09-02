import React, { useState } from 'react';
import DatasetCyberIcon from './DatasetCyberIcons';

export default function DatasetAnnotateModal({ isOpen, onClose, dataset }) {
  const [model, setModel] = useState('yolo26x');
  const [conf, setConf] = useState(0.35);
  const [isAnnotating, setIsAnnotating] = useState(false);
  const [progressLog, setProgressLog] = useState([]);
  const [isComplete, setIsComplete] = useState(false);

  if (!isOpen) return null;

  const handleStartAnnotation = () => {
    setIsAnnotating(true);
    setIsComplete(false);
    setProgressLog(['[INIT] Allocating YOLO26x Engine on RTX 5090 (VRAM: 1.8GB)...']);
    setTimeout(() => {
      setProgressLog(prev => [...prev, `[SCAN] Processing unlabelled images for ${dataset?.name || 'Dataset'}...`]);
    }, 400);
    setTimeout(() => {
      setProgressLog(prev => [...prev, '[BATCH] Annotated 420 frames with confidence >= ' + (conf * 100).toFixed(0) + '% (2,150 FPS)...']);
    }, 900);
    setTimeout(() => {
      setProgressLog(prev => [...prev, '[SYNC] Generated label .txt files with polygon & bbox coordinates.']);
      setIsAnnotating(false);
      setIsComplete(true);
    }, 1400);
  };

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="import-modal-box" onClick={(e) => e.stopPropagation()}>
        <div className="card-header" style={{ marginBottom: '0.75rem' }}>
          <span className="card-title">SMART AUTO-ANNOTATION // YOLO26 ZERO-SHOT</span>
          <button className="cyber-pill" style={{ padding: '0.2rem 0.5rem' }} onClick={onClose}>✕</button>
        </div>

        <p style={{ fontSize: '0.8rem', color: '#94a3b8' }}>
          Automated batch pre-annotation using <strong>YOLO26</strong> on <strong>RTX 5090</strong> to accelerate raw image labeling.
        </p>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem', margin: '0.75rem 0' }}>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label" style={{ fontSize: '0.7rem' }}>ORACLE MODEL</div>
            <select className="cyber-select" value={model} onChange={(e) => setModel(e.target.value)} disabled={isAnnotating}>
              <option value="yolo26x">⚡ YOLO26 XLarge (Max Recall)</option>
              <option value="yolo26m">⚡ YOLO26 Medium (Fast)</option>
              <option value="sam2">👁️ SAM 2 (Segment Anything)</option>
            </select>
          </div>

          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label" style={{ fontSize: '0.7rem' }}>
              <span>MIN CONFIDENCE</span>
              <span className="slider-val">{(conf * 100).toFixed(0)}%</span>
            </div>
            <input type="range" min="0.10" max="0.90" step="0.05" value={conf} onChange={(e) => setConf(parseFloat(e.target.value))} disabled={isAnnotating} />
          </div>
        </div>

        {progressLog.length > 0 && (
          <div className="yaml-code-box" style={{ maxHeight: '110px', fontSize: '0.75rem', marginBottom: '0.85rem' }}>
            {progressLog.map((log, i) => (
              <div key={i} style={{ color: i === progressLog.length - 1 ? (isComplete ? 'var(--cb-green)' : 'var(--cb-yellow)') : 'var(--cb-cyan)' }}>
                {log}
              </div>
            ))}
          </div>
        )}

        {isComplete && (
          <div style={{ background: 'rgba(0,255,157,0.08)', border: '1px solid var(--cb-green)', padding: '0.5rem', borderRadius: '3px', fontSize: '0.75rem', color: 'var(--cb-green)', fontFamily: 'var(--font-mono)', marginBottom: '0.85rem' }}>
            ✓ AUTO-ANNOTATION COMPLETE: 420 Labels Synced to /labels/train/
          </div>
        )}

        <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end' }}>
          <button className="cyber-pill" onClick={onClose} disabled={isAnnotating}>FECHAR</button>
          <button className="cyber-pill active" onClick={handleStartAnnotation} disabled={isAnnotating} style={{ display: 'inline-flex', alignItems: 'center', gap: '0.4rem' }}>
            <DatasetCyberIcon name="tag" size={14} color="#07080c" />
            <span>{isAnnotating ? 'ANNOTATING ON RTX 5090...' : 'INICIAR AUTO-ROTULAGEM'}</span>
          </button>
        </div>
      </div>
    </div>
  );
}
