import React, { useState } from 'react';

export default function ModelDistillationCard({ onLaunchDistill }) {
  const [teacher, setTeacher] = useState('yolo26x');
  const [student, setStudent] = useState('yolo26n');
  const [disLoss, setDisLoss] = useState(1.5);

  return (
    <div className="cyber-card" style={{ marginTop: '1.25rem' }}>
      <div className="card-header">
        <span className="card-title">KNOWLEDGE DISTILLATION STUDIO</span>
        <span className="badge-yellow">TEACHER ➔ STUDENT</span>
      </div>

      <div style={{ fontSize: '0.75rem', color: '#94a3b8', marginBottom: '0.65rem' }}>
        Transfira o poder preditivo do modelo <strong>Teacher ({teacher})</strong> para o <strong>Student ({student})</strong> aumentando mAP com custo zero de inferência.
      </div>

      <div className="distill-connector-box">
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '0.65rem', color: 'var(--cb-yellow)', fontFamily: 'var(--font-mono)' }}>TEACHER (ORACLE)</div>
          <div style={{ fontFamily: 'var(--font-oxanium)', fontWeight: '700', color: '#fff', fontSize: '0.9rem' }}>{teacher}</div>
          <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>56.1% mAP</div>
        </div>

        <div style={{ color: 'var(--cb-cyan)', fontSize: '1.2rem', fontWeight: '700' }}>➔ ⚡ ➔</div>

        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '0.65rem', color: 'var(--cb-cyan)', fontFamily: 'var(--font-mono)' }}>STUDENT (TARGET)</div>
          <div style={{ fontFamily: 'var(--font-oxanium)', fontWeight: '700', color: '#fff', fontSize: '0.9rem' }}>{student}</div>
          <div style={{ fontSize: '0.65rem', color: '#94a3b8' }}>0.38ms Latency</div>
        </div>
      </div>

      <div className="selector-group" style={{ marginBottom: '0.75rem' }}>
        <div className="selector-label">
          <span>DISTILLATION LOSS WEIGHT (dis)</span>
          <span className="slider-val">{disLoss.toFixed(1)}</span>
        </div>
        <input type="range" min="0.5" max="3.0" step="0.1" value={disLoss} onChange={(e) => setDisLoss(parseFloat(e.target.value))} />
      </div>

      <button
        className="cyber-pill active"
        style={{ width: '100%', padding: '0.6rem', textAlign: 'center' }}
        onClick={() => onLaunchDistill({ teacher, student, disLoss })}
      >
        ⚡ START DISTILLATION ON RTX 5090
      </button>
    </div>
  );
}
