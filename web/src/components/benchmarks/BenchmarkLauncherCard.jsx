import React, { useState } from 'react';

const FORMATS = [
  { id: 'TensorRT', label: 'TensorRT (.engine)' },
  { id: 'PyTorch', label: 'PyTorch CUDA' },
  { id: 'ONNX', label: 'ONNX DirectML' },
  { id: 'OpenVINO', label: 'OpenVINO' },
  { id: 'TorchScript', label: 'TorchScript' },
  { id: 'LiteRT', label: 'LiteRT (TFLite)' }
];

export default function BenchmarkLauncherCard({ onLaunch, isRunning }) {
  const [model, setModel] = useState('yolo26n.pt');
  const [data, setData] = useState('coco8.yaml');
  const [imgsz, setImgsz] = useState(640);
  const [quantize, setQuantize] = useState(16);
  const [device, setDevice] = useState('0');
  const [formats, setFormats] = useState(['TensorRT', 'PyTorch', 'ONNX', 'OpenVINO']);

  const toggleFormat = (fmt) => {
    setFormats(prev => prev.includes(fmt) ? prev.filter(f => f !== fmt) : [...prev, fmt]);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onLaunch({ model, data, imgsz: Number(imgsz), quantize: Number(quantize), device, target_formats: formats });
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. BENCHMARK LAUNCHER</span>
        <span className="badge-cyan">RTX 5090 • CUDA 13.3</span>
      </div>
      <form onSubmit={handleSubmit} style={{ marginTop: '0.65rem' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.65rem' }}>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">MODEL CHECKPOINT</div>
            <input className="cyber-input" value={model} onChange={e => setModel(e.target.value)} />
          </div>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">VALIDATION DATASET</div>
            <input className="cyber-input" value={data} onChange={e => setData(e.target.value)} />
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '0.65rem', marginTop: '0.65rem' }}>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">RESOLUTION</div>
            <select className="cyber-select" value={imgsz} onChange={e => setImgsz(e.target.value)}>
              <option value={320}>320x320</option>
              <option value={640}>640x640</option>
              <option value={1280}>1280x1280</option>
            </select>
          </div>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">PRECISION</div>
            <select className="cyber-select" value={quantize} onChange={e => setQuantize(e.target.value)}>
              <option value={16}>FP16 (TensorRT)</option>
              <option value={8}>INT8 (Quantized)</option>
              <option value={32}>FP32 (Float)</option>
            </select>
          </div>
          <div className="selector-group" style={{ marginBottom: 0 }}>
            <div className="selector-label">DEVICE</div>
            <input className="cyber-input" value={device} onChange={e => setDevice(e.target.value)} />
          </div>
        </div>

        <div style={{ marginTop: '0.75rem' }}>
          <div className="selector-label">TARGET RUNTIMES ({formats.length} SELECTED)</div>
          <div className="format-checkbox-grid">
            {FORMATS.map(f => (
              <div key={f.id} className={`format-chip ${formats.includes(f.id) ? 'selected' : ''}`} onClick={() => toggleFormat(f.id)}>
                <span>{formats.includes(f.id) ? '◈' : '◇'}</span>
                <span>{f.label}</span>
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
