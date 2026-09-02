import React, { useState, useEffect } from 'react';
import DatasetListCard from '../components/datasets/DatasetListCard';
import DatasetYamlViewer from '../components/datasets/DatasetYamlViewer';
import DatasetKFoldCard from '../components/datasets/DatasetKFoldCard';
import DatasetAuditCard from '../components/datasets/DatasetAuditCard';
import DatasetImportModal from '../components/datasets/DatasetImportModal';
import DatasetAnnotateModal from '../components/datasets/DatasetAnnotateModal';
import DatasetMergerModal from '../components/datasets/DatasetMergerModal';
import { fetchDatasetsAPI } from '../api/client';

export default function DatasetStudio({ datasets: initialDatasets }) {
  const [datasets, setDatasets] = useState(initialDatasets || []);
  const [selectedDataset, setSelectedDataset] = useState(initialDatasets?.[0] || null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isAnnotateOpen, setIsAnnotateOpen] = useState(false);
  const [isMergeOpen, setIsMergeOpen] = useState(false);

  useEffect(() => {
    if (!initialDatasets || initialDatasets.length === 0) {
      fetchDatasetsAPI().then(data => {
        if (data && data.length > 0) {
          setDatasets(data);
          setSelectedDataset(data[0]);
        }
      });
    } else {
      setDatasets(initialDatasets);
      if (!selectedDataset) setSelectedDataset(initialDatasets[0]);
    }
  }, [initialDatasets]);

  const handleDatasetImported = (newDs) => {
    setDatasets(prev => [newDs, ...prev]);
    setSelectedDataset(newDs);
  };

  return (
    <div className="view-container datasets-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">DATASET STUDIO</h1>
        <p className="cockpit-main-subtitle">
          DATA REPOSITORY • YAML SPECIFICATION • K-FOLD CROSS-VALIDATION & LEAKAGE AUDIT
        </p>
      </div>

      <div className="datasets-grid">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <DatasetListCard
            datasets={datasets}
            selectedId={selectedDataset?.dataset_id}
            onSelectDataset={setSelectedDataset}
            onOpenImportModal={() => setIsModalOpen(true)}
            onOpenAnnotateModal={() => setIsAnnotateOpen(true)}
            onOpenMergeModal={() => setIsMergeOpen(true)}
          />
          <DatasetYamlViewer dataset={selectedDataset} />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <DatasetKFoldCard dataset={selectedDataset} />
          <DatasetAuditCard dataset={selectedDataset} onSendToCockpit={() => { window.location.hash = 'cockpit'; }} />
        </div>
      </div>

      <DatasetImportModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} onDatasetImported={handleDatasetImported} />
      <DatasetAnnotateModal isOpen={isAnnotateOpen} onClose={() => setIsAnnotateOpen(false)} dataset={selectedDataset} />
      <DatasetMergerModal isOpen={isMergeOpen} onClose={() => setIsMergeOpen(false)} datasets={datasets} onDatasetsMerged={handleDatasetImported} />
    </div>
  );
}
