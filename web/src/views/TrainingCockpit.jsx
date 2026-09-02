import React, { useState } from 'react';
import ArchitectureCard from '../components/ArchitectureCard';
import HyperparameterCard from '../components/HyperparameterCard';
import TwoStageControlCard from '../components/TwoStageControlCard';
import HardwareEstimatorCard from '../components/HardwareEstimatorCard';
import { launchTrainingJobAPI } from '../api/client';

export default function TrainingCockpit({ datasets, onJobLaunched }) {
  const [family, setFamily] = useState('yolo26');
  const [scale, setScale] = useState('s');
  const [task, setTask] = useState('detect');
  const [selectedDataset, setSelectedDataset] = useState(datasets?.[0]?.dataset_id || 'coco8');
  const [isLaunching, setIsLaunching] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');
  const [params, setParams] = useState({
    epochs: 50,
    batch_size: 16,
    imgsz: 640,
    optimizer: 'AdamW',
    lr0: 0.001,
    amp: true,
    two_stage: false,
    stage1_epochs: 20,
    stage1_freeze: 10,
    stage2_epochs: 30
  });

  const handleLaunch = async () => {
    setIsLaunching(true);
    setErrorMsg('');
    try {
      const modelArch = `${family}${scale}`;
      const jobId = `run_${modelArch}_${Date.now()}`;
      const config = {
        job_id: jobId,
        model_architecture: modelArch,
        task: task,
        dataset_id: selectedDataset,
        hyperparameters: params
      };

      const result = await launchTrainingJobAPI(config);
      if (onJobLaunched) onJobLaunched(result);
    } catch (err) {
      setErrorMsg(err.message || 'Launch error');
    } finally {
      setIsLaunching(false);
    }
  };

  return (
    <div className="view-container cockpit-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">TRAINING COCKPIT</h1>
        <p className="cockpit-main-subtitle">
          CONFIGURE ARCHITECTURE, HYPERPARAMETERS & INITIALIZE GPU TRAINING RUNS
        </p>
      </div>

      {errorMsg && (
        <div className="cyber-alert alert-critical" style={{ marginBottom: '1rem' }}>
          <span>⚠️ {errorMsg}</span>
        </div>
      )}

      <div className="cockpit-quad-grid">
        <ArchitectureCard family={family} setFamily={setFamily} scale={scale} setScale={setScale} task={task} setTask={setTask} />
        <HyperparameterCard params={params} setParams={setParams} />
        <TwoStageControlCard params={params} setParams={setParams} />
        <HardwareEstimatorCard
          datasets={datasets}
          selectedDataset={selectedDataset}
          setSelectedDataset={setSelectedDataset}
          scale={scale}
          batchSize={params.batch_size}
          imgsz={params.imgsz}
        />
      </div>

      <div className="cockpit-bottom-bar">
        <button className="cockpit-launch-action-btn" onClick={handleLaunch} disabled={isLaunching}>
          {isLaunching ? 'INITIALIZING WORKER...' : 'INITIATE TRAINING RUN'}
        </button>
      </div>
    </div>
  );
}
