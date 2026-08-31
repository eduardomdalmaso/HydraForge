import React from 'react';

const FAMILIES = ['yolov8', 'yolo11', 'yolo26'];
const SCALES = [
  { id: 'n', label: 'Nano (n)' },
  { id: 's', label: 'Small (s)' },
  { id: 'm', label: 'Medium (m)' },
  { id: 'l', label: 'Large (l)' },
  { id: 'x', label: 'XLarge (x)' }
];
const TASKS = [
  { id: 'detect', label: 'Detection' },
  { id: 'segment', label: 'Segmentation' },
  { id: 'pose', label: 'Pose Estimation' },
  { id: 'classify', label: 'Classification' },
  { id: 'obb', label: 'Oriented Box (OBB)' }
];

export default function ArchitectureCard({ family, setFamily, scale, setScale, task, setTask }) {
  const modelTag = `${family}${scale}${task !== 'detect' ? '-' + task : ''}.pt`;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. MODEL ARCHITECTURE MATRIX</span>
        <span className="badge-cyan">{modelTag}</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">MODEL FAMILY</div>
        <div className="pills-container">
          {FAMILIES.map((f) => (
            <button key={f} className={`cyber-pill ${family === f ? 'active' : ''}`} onClick={() => setFamily(f)}>
              {f.toUpperCase()}
            </button>
          ))}
        </div>
      </div>

      <div className="selector-group">
        <div className="selector-label">SCALE & COMPLEXITY</div>
        <div className="pills-container">
          {SCALES.map((s) => (
            <button key={s.id} className={`cyber-pill ${scale === s.id ? 'active' : ''}`} onClick={() => setScale(s.id)}>
              {s.label}
            </button>
          ))}
        </div>
      </div>

      <div className="selector-group" style={{ marginBottom: 0 }}>
        <div className="selector-label">VISION TASK</div>
        <div className="pills-container">
          {TASKS.map((t) => (
            <button key={t.id} className={`cyber-pill ${task === t.id ? 'active' : ''}`} onClick={() => setTask(t.id)}>
              {t.label}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
