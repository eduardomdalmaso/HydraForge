import React, { useState, useEffect } from 'react';
import ModelZooFilterBar from '../components/modelzoo/ModelZooFilterBar';
import ModelCardGrid from '../components/modelzoo/ModelCardGrid';
import ModelDetailsDrawer from '../components/modelzoo/ModelDetailsDrawer';
import ModelDistillationCard from '../components/modelzoo/ModelDistillationCard';
import { fetchModelsAPI } from '../api/client';

export default function ModelZooStudio() {
  const [models, setModels] = useState([]);
  const [selectedFamily, setSelectedFamily] = useState('ALL');
  const [selectedTask, setSelectedTask] = useState('ALL');
  const [selectedModel, setSelectedModel] = useState(null);

  const loadModels = async () => {
    const list = await fetchModelsAPI();
    if (list && list.length > 0) {
      setModels(list);
      setSelectedModel(prev => prev ? (list.find(m => m.id === prev.id) || list[0]) : list[0]);
    }
  };

  useEffect(() => {
    loadModels();
    const interval = setInterval(loadModels, 3000);
    return () => clearInterval(interval);
  }, []);

  const filteredModels = models.filter(m => {
    const famMatch = selectedFamily === 'ALL' || m.family === selectedFamily;
    const taskMatch = selectedTask === 'ALL' || m.task === selectedTask;
    return famMatch && taskMatch;
  });

  const navigateTo = (tab) => { window.location.hash = tab; };

  return (
    <div className="view-container zoo-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">MODEL ZOO REPOSITORY</h1>
        <p className="cockpit-main-subtitle">
          OFFICIAL CHECKPOINTS • CUSTOM TRAINED MODELS ON RTX 5090 • DISTILLATION & RUNTIMES
        </p>
      </div>

      <ModelZooFilterBar
        selectedFamily={selectedFamily}
        setSelectedFamily={setSelectedFamily}
        selectedTask={selectedTask}
        setSelectedTask={setSelectedTask}
        totalCount={filteredModels.length}
      />

      <div className="zoo-layout">
        <div>
          <ModelCardGrid
            models={filteredModels}
            selectedModel={selectedModel}
            onSelectModel={setSelectedModel}
          />
          <ModelDistillationCard onLaunchDistill={() => navigateTo('cockpit')} />
        </div>

        <div>
          <ModelDetailsDrawer
            model={selectedModel}
            onSendToCockpit={() => navigateTo('cockpit')}
            onSendToBenchmark={() => navigateTo('benchmarks')}
            onSendToPlayground={() => navigateTo('playground')}
          />
        </div>
      </div>
    </div>
  );
}
