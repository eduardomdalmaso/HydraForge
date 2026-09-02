import React, { useState, useEffect } from 'react';
import DatasetListCard from '../components/datasets/DatasetListCard';
import DatasetYamlViewer from '../components/datasets/DatasetYamlViewer';
import DatasetKFoldCard from '../components/datasets/DatasetKFoldCard';
import DatasetAuditCard from '../components/datasets/DatasetAuditCard';
import DatasetImportModal from '../components/datasets/DatasetImportModal';
import DatasetAnnotateModal from '../components/datasets/DatasetAnnotateModal';
import DatasetMergerModal from '../components/datasets/DatasetMergerModal';
import CyberAlertModal from '../components/CyberAlertModal';
import { fetchDatasetsAPI, rescanDatasetsAPI, deleteDatasetAPI } from '../api/client';

export default function DatasetStudio({ datasets: initialDatasets }) {
  const [datasets, setDatasets] = useState(initialDatasets || []);
  const [selectedDataset, setSelectedDataset] = useState(initialDatasets?.[0] || null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isAnnotateOpen, setIsAnnotateOpen] = useState(false);
  const [isMergeOpen, setIsMergeOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [deleteDisk, setDeleteDisk] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [globalMappings, setGlobalMappings] = useState(() => {
    try { return JSON.parse(localStorage.getItem('hydraforge_mappings') || '{}'); } catch { return {}; }
  });

  const saveMapping = (dsId, mapping) => {
    setGlobalMappings(prev => {
      const next = { ...prev, [dsId]: { ...(prev[dsId] || {}), ...mapping } };
      try { localStorage.setItem('hydraforge_mappings', JSON.stringify(next)); } catch {}
      return next;
    });
  };

  const handleRescan = async () => {
    const data = await rescanDatasetsAPI();
    if (data?.length) { setDatasets(data); if (!selectedDataset || !data.some(d => d.dataset_id === selectedDataset.dataset_id)) setSelectedDataset(data[0]); }
  };

  const confirmDelete = async () => {
    if (!deleteTarget) return;
    setIsDeleting(true);
    if (await deleteDatasetAPI(deleteTarget.dataset_id, deleteDisk)) {
      const rem = datasets.filter(d => d.dataset_id !== deleteTarget.dataset_id);
      setDatasets(rem);
      if (selectedDataset?.dataset_id === deleteTarget.dataset_id) setSelectedDataset(rem[0] || null);
    }
    setIsDeleting(false); setDeleteTarget(null);
  };

  useEffect(() => {
    if (!initialDatasets?.length) {
      fetchDatasetsAPI().then(data => { if (data?.length) { setDatasets(data); setSelectedDataset(data[0]); } });
    } else {
      setDatasets(initialDatasets);
      if (!selectedDataset) setSelectedDataset(initialDatasets[0]);
    }
  }, [initialDatasets]);

  const handleDatasetImported = (newDs) => { setDatasets(prev => [newDs, ...prev]); setSelectedDataset(newDs); };

  return (
    <div className="view-container datasets-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">DATASET STUDIO</h1>
        <p className="cockpit-main-subtitle">DATA REPOSITORY • YAML SPECIFICATION • K-FOLD CROSS-VALIDATION & LEAKAGE AUDIT</p>
      </div>

      <div className="datasets-grid">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <DatasetListCard
            datasets={datasets} selectedId={selectedDataset?.dataset_id} onSelectDataset={setSelectedDataset}
            onOpenImportModal={() => setIsModalOpen(true)} onOpenAnnotateModal={() => setIsAnnotateOpen(true)}
            onOpenMergeModal={() => setIsMergeOpen(true)} onRescan={handleRescan} onDeleteDataset={setDeleteTarget}
          />
          <DatasetYamlViewer dataset={selectedDataset} initialMappings={globalMappings[selectedDataset?.dataset_id] || {}} onSaveMappings={saveMapping} />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <DatasetKFoldCard dataset={selectedDataset} />
          <DatasetAuditCard dataset={selectedDataset} onSendToCockpit={() => { window.location.hash = 'cockpit'; }} />
        </div>
      </div>

      <DatasetImportModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} onDatasetImported={handleDatasetImported} />
      <DatasetAnnotateModal isOpen={isAnnotateOpen} onClose={() => setIsAnnotateOpen(false)} dataset={selectedDataset} />
      <DatasetMergerModal isOpen={isMergeOpen} onClose={() => setIsMergeOpen(false)} datasets={datasets} onDatasetsMerged={handleDatasetImported} savedMappings={globalMappings} />
      <CyberAlertModal
        isOpen={Boolean(deleteTarget)} title="DELETE DATASET" type="danger"
        message={`Delete dataset "${deleteTarget?.name}" (${deleteTarget?.dataset_id})?`}
        confirmText="CONFIRM DELETE" cancelText="ABORT" checkboxLabel="Also purge physical files from disk"
        checkboxChecked={deleteDisk} onToggleCheckbox={setDeleteDisk}
        onConfirm={confirmDelete} onClose={() => setDeleteTarget(null)} isProcessing={isDeleting}
      />
    </div>
  );
}
