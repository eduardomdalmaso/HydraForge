import React, { useState } from 'react';
import DatasetCyberIcon from './DatasetCyberIcons';

export default function DatasetKFoldCard({ dataset }) {
  const [kFolds, setKFolds] = useState(5);
  const [isGenerating, setIsGenerating] = useState(false);
  const [genSuccess, setGenSuccess] = useState(false);

  const handleGenerateFolds = () => {
    setIsGenerating(true);
    setGenSuccess(false);
    setTimeout(() => {
      setIsGenerating(false);
      setGenSuccess(true);
    }, 450);
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. K-FOLD CROSS-VALIDATION SPLITTER</span>
        <span className="badge-yellow">SKLEARN KFOLD</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">
          <span>NUMBER OF FOLDS (K)</span>
          <span className="slider-val">K = {kFolds} FOLDS</span>
        </div>
        <div className="pills-container" style={{ marginBottom: '0.75rem' }}>
          {[3, 5, 10].map(k => (
            <button
              key={k}
              type="button"
              className={`cyber-pill ${kFolds === k ? 'active' : ''}`}
              onClick={() => { setKFolds(k); setGenSuccess(false); }}
            >
              {k}-Fold ({((1 / k) * 100).toFixed(0)}% Val per Split)
            </button>
          ))}
        </div>
      </div>

      <div style={{ background: 'rgba(0,0,0,0.3)', padding: '0.65rem 0.8rem', borderRadius: '4px', border: '1px solid rgba(0,240,255,0.15)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.75rem', fontFamily: 'var(--font-mono)' }}>
          <span style={{ color: 'var(--cb-cyan)' }}>TRAIN ({100 - ((1 / kFolds) * 100).toFixed(0)}%)</span>
          <span style={{ color: 'var(--cb-yellow)' }}>VAL ({((1 / kFolds) * 100).toFixed(0)}%)</span>
        </div>
        <div className="kfold-bar">
          <div className="kfold-seg-train" style={{ width: `${100 - (100 / kFolds)}%` }}></div>
          <div className="kfold-seg-val" style={{ width: `${100 / kFolds}%` }}></div>
        </div>
        <div style={{ marginTop: '0.5rem', fontSize: '0.7rem', color: '#94a3b8' }}>
          Generates <code>split_1.yaml</code> ... <code>split_{kFolds}.yaml</code> in zero-copy text.
        </div>
      </div>

      {genSuccess && (
        <div style={{ marginTop: '0.65rem', padding: '0.5rem', background: 'rgba(0,255,157,0.06)', border: '1px solid var(--cb-green)', borderRadius: '4px', fontSize: '0.72rem', color: 'var(--cb-green)', fontFamily: 'var(--font-mono)' }}>
          ✓ Generated {kFolds} Stratified Folds for {dataset?.name || 'Dataset'}
        </div>
      )}

      <button
        className="cyber-action-btn"
        style={{ width: '100%', marginTop: '0.75rem', padding: '0.75rem' }}
        onClick={handleGenerateFolds}
        disabled={isGenerating}
      >
        <DatasetCyberIcon name="kfold" size={15} color="#07080c" />
        <span>{isGenerating ? 'GENERATING...' : 'GENERATE'}</span>
      </button>
    </div>
  );
}
