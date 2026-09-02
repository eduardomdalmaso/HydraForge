import React, { useState } from 'react';
import PlaygroundControlBar from '../components/playground/PlaygroundControlBar';
import PlaygroundTuningCard from '../components/playground/PlaygroundTuningCard';
import PlaygroundViewport from '../components/playground/PlaygroundViewport';
import PlaygroundTelemetryCard from '../components/playground/PlaygroundTelemetryCard';

export default function PlaygroundStudio() {
  const [config, setConfig] = useState({
    model: 'yolo26n',
    source: 'sample_bus',
    runtime: 'tensorrt',
    conf: 0.25,
    iou: 0.45,
    nmsFree: true,
    sahi: false,
    isolateBg: false
  });

  const [detections, setDetections] = useState([]);
  const [selectedEntity, setSelectedEntity] = useState(null);
  const [isScanning, setIsScanning] = useState(false);
  const [telemetry, setTelemetry] = useState(null);

  const handleRunInference = () => {
    setIsScanning(true);
    setTimeout(() => {
      const generated = [
        { id: 1, label: config.isolateBg ? 'Segment Mask' : 'Target Object', conf: 0.94, box: [20, 25, 38, 46], color: '#00f0ff' },
        { id: 2, label: 'Detected Entity', conf: 0.88, box: [65, 35, 18, 40], color: '#fcee0a' }
      ].filter(d => d.conf >= config.conf);

      setDetections(generated);
      setIsScanning(false);
      setTelemetry({
        preprocess_ms: '0.01',
        inference_ms: (0.35 + Math.random() * 0.05).toFixed(2),
        postprocess_ms: '0.02',
        fps: '2,439',
        vram_mb: '1,420'
      });
    }, 450);
  };

  return (
    <div className="view-container playground-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">KIROSHI PLAYGROUND</h1>
        <p className="cockpit-main-subtitle">
          REAL-TIME OPTICAL SCANNER • SAHI 4K SLICING • BACKGROUND ISOLATOR & TENSORRT INFERENCE
        </p>
      </div>

      <div className="playground-layout">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <PlaygroundControlBar config={config} setConfig={setConfig} onRunInference={handleRunInference} isRunning={isScanning} />
          <PlaygroundTuningCard config={config} setConfig={setConfig} />
        </div>

        <PlaygroundViewport
          imageSrc="/hydra-logo.jpg"
          detections={detections}
          selectedEntity={selectedEntity}
          setSelectedEntity={setSelectedEntity}
          isScanning={isScanning}
        />

        <PlaygroundTelemetryCard
          telemetry={telemetry}
          detections={detections}
          selectedEntity={selectedEntity}
          onSelectEntity={setSelectedEntity}
        />
      </div>
    </div>
  );
}
