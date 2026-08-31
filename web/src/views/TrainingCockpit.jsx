import React, { useState } from 'react';
import ArchitectureCard from '../components/ArchitectureCard';
import HyperparameterCard from '../components/HyperparameterCard';
import HardwareEstimatorCard from '../components/HardwareEstimatorCard';
import { launchTrainingJobAPI } from '../api/client';

export default function TrainingCockpit({ datasets, onJobLaunched }) {
  const [family, setFamily] = useState('yolo11');
  const [scale, setScale] = useState('s');
  const [task, setTask] = useState('detect');
  const [selectedDataset, setSelectedDataset] = useState('coco8');
  const [isLaunching, setIsLaunching] = useState(false);
  const [params, setParams] = useState({
    epochs: 100,
    batch_size: 16,
    imgsz: 640,
    optimizer: 'AdamW',
    lr0: 0.001,
    amp: true
  });

  const handleLaunch = async () => {
    setIsLaunching(true);
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
      alert(`Launch error: ${err.message}`);
    } finally {
      setIsLaunching(false);
    }
  };

  return (
    <div className="view-container">
      <div className="view-header">
        <h1 className="view-title">TRAINING COCKPIT</h1>
        <p className="view-subtitle">CONFIGURE ARCHITECTURE, HYPERPARAMETERS & INITIALIZE GPU TRAINING RUNS</p>
      </div>

      <div className="cockpit-grid">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <ArchitectureCard family={family} setFamily={setFamily} scale={scale} setScale={setScale} task={task} setTask={setTask} />
          <HyperparameterCard params={params} setParams={setParams} />
        </div>

        <div>
          <HardwareEstimatorCard
            datasets={datasets}
            selectedDataset={selectedDataset}
            setSelectedDataset={setSelectedDataset}
            scale={scale}
            batchSize={params.batch_size}
            imgsz={params.imgsz}
            onLaunch={handleLaunch}
            isLaunching={isLaunching}
          />
        </div>
      </div>
    </div>
  );
}
