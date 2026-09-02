import React, { useState } from 'react';

const FORMATS = [
  { id: 'TensorRT', label: 'TensorRT (.engine)', gpuOnly: true }, { id: 'PyTorch', label: 'PyTorch' },
  { id: 'ONNX', label: 'ONNX DirectML' }, { id: 'OpenVINO', label: 'OpenVINO' },
  { id: 'TorchScript', label: 'TorchScript' }, { id: 'LiteRT', label: 'LiteRT (TFLite)' }
];
const BASE_MODELS = ['yolo26n.pt', 'yolo26s.pt', 'yolo26m.pt', 'yolo26l.pt', 'yolov8n.pt', 'yolo11n.pt'];

export default function BenchmarkLauncherCard({ onLaunch, isRunning, jobs = [], datasets = [], gpuStats }) {
  const [model, setModel] = useState('yolo26n.pt');
  const [data, setData] = useState('coco8.yaml');
  const [imgsz, setImgsz] = useState(640);
  const [quantize, setQuantize] = useState(16);
  const [device, setDevice] = useState('0');
  const [formats, setFormats] = useState(['TensorRT', 'PyTorch', 'ONNX', 'OpenVINO']);

  const onDeviceChange = (val) => {
    setDevice(val);
    setFormats(p => val === 'cpu' ? p.filter(f => f !== 'TensorRT') : (!p.includes('TensorRT') ? ['TensorRT', ...p] : p));
  };
  const toggle = (id, gpuOnly) => !(device === 'cpu' && gpuOnly) && setFormats(p => p.includes(id) ? p.filter(f => f !== id) : [...p, id]);
  const handleSubmit = (e) => {
    e.preventDefault();
    onLaunch({ model, data, imgsz: Number(imgsz), quantize: Number(quantize), device, target_formats: formats });
  };

  const gpuName = gpuStats?.model || 'NVIDIA RTX 5090';
  const vramGB = Math.round((gpuStats?.total_vram_mb || 32607) / 1024);

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. BENCHMARK LAUNCHER</span>
        <span className="badge-cyan">{device === 'cpu' ? 'CPU HOST' : `${gpuName} • CUDA 13.3`}</span>
      </div>
      <form onSubmit={handleSubmit} style={{ marginTop: '0.65rem' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.65rem' }}>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">MODEL CHECKPOINT</div>
            <select className="cyber-select" value={model} onChange={e => setModel(e.target.value)}>
              <optgroup label="Trained Models">
                {jobs.map(j => <option key={j.job_id} value={j.output_weights || j.model_architecture}>🎯 [Trained] {j.model_architecture} ({j.dataset_id}){j.best_map50 ? ` - ${(j.best_map50 * 100).toFixed(1)}%` : ''}</option>)}
              </optgroup>
              <optgroup label="Base Pretrained Models">{BASE_MODELS.map(m => <option key={m} value={m}>{m}</option>)}</optgroup>
            </select>
          </div>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">VALIDATION DATASET</div>
            <select className="cyber-select" value={data} onChange={e => setData(e.target.value)}>
              <optgroup label="Workspace Datasets">
                {datasets.map(d => <option key={d.dataset_id} value={d.yaml_path || d.dataset_id}>📁 {d.name || d.dataset_id} ({d.num_classes || d.classes?.length || 1} cls)</option>)}
              </optgroup>
              <optgroup label="Reference Datasets"><option value="coco8.yaml">coco8.yaml (YOLO Demo)</option></optgroup>
            </select>
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '0.65rem', marginTop: '0.65rem' }}>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">RESOLUTION</div>
            <select className="cyber-select" value={imgsz} onChange={e => setImgsz(e.target.value)}>
              <option value={320}>320x320</option><option value={640}>640x640</option><option value={1280}>1280x1280</option>
            </select>
          </div>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">PRECISION</div>
            <select className="cyber-select" value={quantize} onChange={e => setQuantize(e.target.value)}>
              <option value={16}>FP16 (TensorRT)</option><option value={8}>INT8 (Quantized)</option><option value={32}>FP32 (Float)</option>
            </select>
          </div>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">DEVICE</div>
            <select className="cyber-select" value={device} onChange={e => onDeviceChange(e.target.value)}>
              <option value="0">🎮 GPU 0: {gpuName} ({vramGB} GB)</option>
              <option value="cpu">💻 CPU: Host Multi-Core</option>
            </select>
          </div>
        </div>

        <div style={{ marginTop: '0.75rem' }}>
          <div className="selector-label">TARGET RUNTIMES ({formats.length} SELECTED)</div>
          <div className="format-checkbox-grid">
            {FORMATS.map(f => (
              <div key={f.id} className={`format-chip ${formats.includes(f.id) ? 'selected' : ''}`} style={{ opacity: device === 'cpu' && f.gpuOnly ? 0.35 : 1, cursor: device === 'cpu' && f.gpuOnly ? 'not-allowed' : 'pointer' }} onClick={() => toggle(f.id, f.gpuOnly)}>
                <span>{formats.includes(f.id) ? '◈' : '◇'}</span> <span>{f.label} {device === 'cpu' && f.gpuOnly ? '(CUDA)' : ''}</span>
              </div>
            ))}
          </div>
        </div>

        <button type="submit" disabled={isRunning || formats.length === 0} className="cyber-action-btn" style={{ width: '100%', marginTop: '1rem', padding: '0.75rem' }}>
          {isRunning ? 'BENCHMARKING...' : 'BENCHMARK'}
        </button>
      </form>
    </div>
  );
}
