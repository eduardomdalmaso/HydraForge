import React from 'react';

export default function CyberAlertModal({
  isOpen,
  title = 'SYSTEM ALERT',
  message,
  type = 'danger', // 'danger' | 'cyan' | 'yellow'
  confirmText = 'CONFIRM',
  cancelText = 'ABORT',
  checkboxLabel,
  checkboxChecked,
  onToggleCheckbox,
  onConfirm,
  onClose,
  isProcessing = false
}) {
  if (!isOpen) return null;

  const boxClass = type === 'cyan' ? 'cyber-alert-box cyan' : (type === 'yellow' ? 'cyber-alert-box yellow' : 'cyber-alert-box');
  const confirmBtnClass = type === 'cyan' ? 'cyber-action-btn' : (type === 'yellow' ? 'cyber-action-btn' : 'cyber-action-btn');
  const confirmStyle = type === 'danger' ? { background: 'var(--cb-magenta)', borderColor: 'var(--cb-magenta)', color: '#fff' } : {};

  return (
    <div className="cyber-alert-overlay" onClick={onClose}>
      <div className={boxClass} onClick={(e) => e.stopPropagation()}>
        <div className="hud-corner-tl" />
        <div className="hud-corner-br" />

        <div className="cyber-alert-header">
          <div className="cyber-alert-title">
            <span>{type === 'danger' ? '⚠️' : (type === 'yellow' ? '⚡' : 'ℹ️')}</span>
            <span>{title}</span>
          </div>
          <button className="cyber-pill" onClick={onClose} style={{ padding: '0.15rem 0.45rem', fontSize: '0.7rem' }}>✕</button>
        </div>

        <div className="cyber-alert-body">
          {message}
        </div>

        {checkboxLabel && (
          <div style={{ margin: '0.75rem 0 1.25rem', padding: '0.5rem 0.65rem', background: 'rgba(255, 0, 85, 0.08)', border: '1px dashed rgba(255, 0, 85, 0.3)', borderRadius: '3px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer', fontSize: '0.78rem', color: '#fff', fontFamily: 'var(--font-mono)' }}>
              <input
                type="checkbox"
                checked={checkboxChecked}
                onChange={(e) => onToggleCheckbox && onToggleCheckbox(e.target.checked)}
              />
              <span>{checkboxLabel}</span>
            </label>
          </div>
        )}

        <div className="cyber-alert-actions">
          <button className="cyber-action-btn secondary" onClick={onClose} disabled={isProcessing}>
            {cancelText}
          </button>
          <button className={confirmBtnClass} style={confirmStyle} onClick={onConfirm} disabled={isProcessing}>
            {isProcessing ? 'PROCESSING...' : confirmText}
          </button>
        </div>
      </div>
    </div>
  );
}
