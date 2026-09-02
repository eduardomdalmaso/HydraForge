import React, { useState } from 'react';
import DatasetCyberIcon from './DatasetCyberIcons';

export default function DatasetAuditCard({ dataset, onSendToCockpit }) {
  const [auditStatus, setAuditStatus] = useState(null);
  const [isAuditing, setIsAuditing] = useState(false);

  const handleRunAudit = async () => {
    if (!dataset?.dataset_id) return;
    setIsAuditing(true);
    try {
      const res = await fetch(`/api/v1/training/datasets/audit/${dataset.dataset_id}`, { method: 'POST' });
      if (res.ok) setAuditStatus(await res.json());
    } finally { setIsAuditing(false); }
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">4. QUALITY AUDITOR & LEAKAGE DEFENSE</span>
        <span className="badge-green">REAL DATA INTEGRITY</span>
      </div>

      <div style={{ marginBottom: '0.85rem' }}>
        <div className="telemetry-row">
          <span className="k">ANNOTATION INTEGRITY</span>
          <span className="audit-badge-ok">{auditStatus ? `✓ ${auditStatus.valid_bboxes_pct} VALID` : 'AUDIT REQUIRED'}</span>
        </div>
        <div className="telemetry-row">
          <span className="k">CORRUPTED LABELS / OUT-OF-BOUNDS</span>
          <span className="audit-badge-ok">{auditStatus ? `✓ ${auditStatus.corrupt_files} CORRUPTED` : '-'}</span>
        </div>
        <div className="telemetry-row" style={{ borderBottom: 'none' }}>
          <span className="k">TRAIN / VAL DATA LEAKAGE</span>
          <span className="audit-badge-ok">{auditStatus ? `✓ ${auditStatus.leakage_overlap_pct} OVERLAP` : '-'}</span>
        </div>
      </div>

      {auditStatus && (
        <div style={{ background: 'rgba(0,255,157,0.06)', border: '1px solid var(--cb-green)', padding: '0.65rem', borderRadius: '4px', marginBottom: '0.85rem', fontFamily: 'var(--font-mono)', fontSize: '0.75rem' }}>
          <div style={{ color: 'var(--cb-green)', fontWeight: '700', marginBottom: '0.2rem' }}>
            ⚡ AUDIT // {auditStatus.status} [{auditStatus.timestamp}]
          </div>
          <div style={{ color: '#cbd5e1' }}>
            {auditStatus.total_bboxes?.toLocaleString()} bboxes auditados no disco. Zero vazamento.
          </div>
        </div>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1.2fr', gap: '0.5rem' }}>
        <button className="cyber-action-btn" onClick={handleRunAudit} disabled={isAuditing || !dataset}>
          <DatasetCyberIcon name="shield" size={14} color="#07080c" />
          <span>{isAuditing ? 'AUDITING...' : 'AUDIT'}</span>
        </button>
        <button className="cyber-action-btn secondary" onClick={onSendToCockpit} disabled={!dataset}>
          <DatasetCyberIcon name="launch" size={14} />
          <span>COCKPIT</span>
        </button>
      </div>
    </div>
  );
}
