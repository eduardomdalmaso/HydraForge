import React, { useState } from 'react';
import DatasetCyberIcon from './DatasetCyberIcons';
import { registerDatasetPathAPI } from '../../api/client';

export default function DatasetImportModal({ isOpen, onClose, onDatasetImported }) {
  const [activeTab, setActiveTab] = useState('zip');
  const [zipFiles, setZipFiles] = useState([]);
  const [folderPath, setFolderPath] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState(null);

  if (!isOpen) return null;

  const handleExecuteImport = async () => {
    setIsProcessing(true);
    setError(null);
    try {
      if (activeTab === 'zip' && zipFiles.length > 0) {
        for (let i = 0; i < zipFiles.length; i++) {
          const file = zipFiles[i];
          const formData = new FormData();
          formData.append('file', file);
          const cleanId = file.name.replace(/\.zip$/i, '').toLowerCase().replace(/[^a-z0-9_]/g, '_');
          formData.append('dataset_id', cleanId);
          formData.append('task', 'detect');
          const res = await fetch('/api/v1/training/datasets/import', { method: 'POST', body: formData });
          if (res.ok) {
            const registered = await res.json();
            if (onDatasetImported) onDatasetImported(registered);
          }
        }
      } else if (activeTab === 'path' && folderPath.trim()) {
        const registered = await registerDatasetPathAPI({ path: folderPath.trim() });
        if (onDatasetImported) onDatasetImported(registered);
      }
      onClose();
    } catch (err) {
      setError(err.message || 'Import failed');
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="import-modal-box" style={{ maxWidth: '600px' }} onClick={(e) => e.stopPropagation()}>
        <div className="card-header" style={{ marginBottom: '0.6rem' }}>
          <span className="card-title">IMPORT DATASET (MULTI-ZIP / FOLDER)</span>
          <button className="cyber-pill" onClick={onClose}>✕</button>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.4rem', marginBottom: '0.6rem' }}>
          <button className={`cyber-pill ${activeTab === 'zip' ? 'active' : ''}`} style={{ width: '100%', justifyContent: 'center' }} onClick={() => setActiveTab('zip')}>
            📦 MULTIPLE .ZIP FILES
          </button>
          <button className={`cyber-pill ${activeTab === 'path' ? 'active' : ''}`} style={{ width: '100%', justifyContent: 'center' }} onClick={() => setActiveTab('path')}>
            📁 LOCAL FOLDER PATH
          </button>
        </div>

        {activeTab === 'zip' ? (
          <label className="dropzone-box" style={{ cursor: 'pointer', display: 'block', margin: '0.5rem 0' }}
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => { e.preventDefault(); setZipFiles(Array.from(e.dataTransfer.files).filter(f => f.name.endsWith('.zip'))); }}>
            <input type="file" accept=".zip" multiple style={{ display: 'none' }} onChange={(e) => setZipFiles(Array.from(e.target.files || []))} />
            <div style={{ textAlign: 'center' }}>
              <DatasetCyberIcon name="zip" size={28} color={zipFiles.length ? 'var(--cb-yellow)' : 'var(--cb-cyan)'} />
              <div style={{ fontFamily: 'var(--font-oxanium)', color: zipFiles.length ? 'var(--cb-yellow)' : 'var(--cb-cyan)', fontWeight: '700', marginTop: '0.2rem' }}>
                {zipFiles.length ? `✓ ${zipFiles.length} ZIP FILE(S) SELECTED` : 'CLICK OR DRAG MULTIPLE .ZIP ARCHIVES'}
              </div>
            </div>
          </label>
        ) : (
          <div style={{ margin: '0.6rem 0' }}>
            <div style={{ fontSize: '0.72rem', color: '#94a3b8', marginBottom: '0.3rem' }}>FILESYSTEM PATH (e.g. /home/hades/Downloads/cell-phone):</div>
            <input className="cyber-input" type="text" placeholder="/home/hades/Downloads/cell-phone" value={folderPath} onChange={(e) => setFolderPath(e.target.value)} style={{ width: '100%', padding: '0.5rem', fontFamily: 'var(--font-mono)', fontSize: '0.8rem' }} />
          </div>
        )}

        {error && <div style={{ color: 'var(--cb-magenta)', fontSize: '0.75rem', marginTop: '0.3rem' }}>⚠️ {error}</div>}

        <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end', marginTop: '0.75rem' }}>
          <button className="cyber-action-btn secondary" onClick={onClose} disabled={isProcessing}>CANCEL</button>
          <button className="cyber-action-btn" onClick={handleExecuteImport} disabled={isProcessing || (activeTab === 'zip' ? zipFiles.length === 0 : !folderPath.trim())}>
            <span>{isProcessing ? 'IMPORTING...' : 'CONFIRM IMPORT'}</span>
          </button>
        </div>
      </div>
    </div>
  );
}
