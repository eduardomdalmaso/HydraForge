import React, { useState, useEffect } from 'react';
import TopCyberNav from './components/TopCyberNav';
import TrainingCockpit from './views/TrainingCockpit';
import LiveHudStudio from './views/LiveHudStudio';
import BenchmarkStudio from './views/BenchmarkStudio';
import PlaygroundStudio from './views/PlaygroundStudio';
import DatasetStudio from './views/DatasetStudio';
import ModelZooStudio from './views/ModelZooStudio';
import { fetchTelemetryAPI, fetchDatasetsAPI } from './api/client';

export default function App() {
  const getInitialTab = () => {
    const hash = window.location.hash.replace('#', '');
    return hash || 'cockpit';
  };

  const [activeTab, setActiveTab] = useState(getInitialTab);
  const [telemetry, setTelemetry] = useState(null);
  const [datasets, setDatasets] = useState([]);

  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.replace('#', '');
      setActiveTab(hash || 'cockpit');
    };
    window.addEventListener('hashchange', handleHashChange);
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, []);

  const navigateTo = (tabId) => {
    window.location.hash = tabId;
    setActiveTab(tabId);
  };

  useEffect(() => {
    const loadInitialData = async () => {
      const tel = await fetchTelemetryAPI();
      if (tel) setTelemetry(tel);
      const ds = await fetchDatasetsAPI();
      if (ds) setDatasets(ds);
    };
    loadInitialData();
    const interval = setInterval(loadInitialData, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="app-layout" style={{ flexDirection: 'column' }}>
      <TopCyberNav activeTab={activeTab} onSelectTab={navigateTo} gpuStats={telemetry?.gpu_stats} />
      <main className="content-body">
        {activeTab === 'cockpit' && (
          <TrainingCockpit datasets={datasets} onJobLaunched={() => navigateTo('live-hud')} />
        )}
        {activeTab === 'live-hud' && <LiveHudStudio />}
        {activeTab === 'benchmarks' && <BenchmarkStudio />}
        {activeTab === 'datasets' && <DatasetStudio datasets={datasets} />}
        {activeTab === 'playground' && <PlaygroundStudio />}
        {activeTab === 'model-zoo' && <ModelZooStudio />}
      </main>
    </div>
  );
}
