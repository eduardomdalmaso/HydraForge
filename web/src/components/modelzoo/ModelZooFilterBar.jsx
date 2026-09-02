import React from 'react';

const FAMILIES = ['ALL', 'YOLO26', 'YOLO11', 'YOLOv8', 'SAM/DETR'];
const TASKS = ['ALL', 'DETECT', 'SEGMENT', 'POSE', 'OBB'];

export default function ModelZooFilterBar({ selectedFamily, setSelectedFamily, selectedTask, setSelectedTask, totalCount }) {
  return (
    <div className="cyber-card" style={{ padding: '0.85rem 1rem', marginBottom: '0.5rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '0.75rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', flexWrap: 'wrap' }}>
          <span style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.75rem', color: 'var(--cb-cyan)', fontWeight: '700' }}>
            FAMILY:
          </span>
          <div className="pills-container" style={{ margin: 0 }}>
            {FAMILIES.map(fam => (
              <button
                key={fam}
                type="button"
                className={`cyber-pill ${selectedFamily === fam ? 'active' : ''}`}
                style={{ padding: '0.2rem 0.55rem', fontSize: '0.72rem' }}
                onClick={() => setSelectedFamily(fam)}
              >
                {fam}
              </button>
            ))}
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', flexWrap: 'wrap' }}>
          <span style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.75rem', color: 'var(--cb-yellow)', fontWeight: '700' }}>
            TASK:
          </span>
          <div className="pills-container" style={{ margin: 0 }}>
            {TASKS.map(tsk => (
              <button
                key={tsk}
                type="button"
                className={`cyber-pill ${selectedTask === tsk ? 'active' : ''}`}
                style={{ padding: '0.2rem 0.55rem', fontSize: '0.72rem' }}
                onClick={() => setSelectedTask(tsk)}
              >
                {tsk}
              </button>
            ))}
          </div>
        </div>

        <span className="badge-cyan" style={{ fontSize: '0.75rem' }}>
          {totalCount} CHECKPOINTS
        </span>
      </div>
    </div>
  );
}
