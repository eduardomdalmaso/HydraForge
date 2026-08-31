import React, { useState, useEffect } from 'react';
import Sidebar from './components/Sidebar';
import Header from './components/Header';
import TrainingCockpit from './views/TrainingCockpit';
import { fetchTelemetryAPI, fetchDatasetsAPI } from './api/client';

export default function App() {
  const [activeTab, setActiveTab] = useState('cockpit');
  const [telemetry, setTelemetry] = useState(null);
  const [datasets, setDatasets] = useState([
    { dataset_id: 'coco8', name: 'COCO8 Tiny Sample', num_classes: 8, train_images: 4, val_images: 4 }
  ]);

  useEffect(() => {
    const loadInitialData = async () => {
      const tel = await fetchTelemetryAPI();
      if (tel) setTelemetry(tel);
      const ds = await fetchDatasetsAPI();
      if (ds && ds.length > 0) setDatasets(ds);
    };
    loadInitialData();

    const interval = setInterval(loadInitialData, 3000);
    return () => clearInterval(interval);
  }, []);

  const handleJobLaunched = (job) => {
    setActiveTab('live-hud');
  };

  return (
    <div className="app-layout">
      <Sidebar activeTab={activeTab} onSelectTab={setActiveTab} />
      <div className="main-content">
        <Header gpuTelemetry={telemetry?.gpu_stats} activeJobsCount={telemetry?.active_jobs_count} />
        <main className="content-body">
          {activeTab === 'cockpit' && (
            <TrainingCockpit datasets={datasets} onJobLaunched={handleJobLaunched} />
          )}
          {activeTab !== 'cockpit' && (
            <div className="view-container">
              <div className="cyber-card" style={{ textAlign: 'center', padding: '3rem' }}>
                <h2 style={{ color: 'var(--cb-cyan)', fontFamily: 'var(--font-oxanium)' }}>
                  VIEW UNDER CONSTRUCTION: {activeTab.toUpperCase()}
                </h2>
                <p style={{ color: '#94a3b8', marginTop: '0.5rem' }}>
                  Next module to be built in sequence. Select <strong>Training Cockpit</strong> to configure models.
                </p>
                <button className="cyber-pill active" style={{ marginTop: '1rem' }} onClick={() => setActiveTab('cockpit')}>
                  ← Return to Training Cockpit
                </button>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  );
}
