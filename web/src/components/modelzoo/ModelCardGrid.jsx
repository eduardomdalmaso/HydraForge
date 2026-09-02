import React from 'react';

export default function ModelCardGrid({ models, selectedModel, onSelectModel }) {
  return (
    <div className="model-cards-grid">
      {models.map((m) => {
        const isSelected = selectedModel?.id === m.id;
        return (
          <div
            key={m.id}
            className={`model-item-card ${isSelected ? 'selected' : ''}`}
            onClick={() => onSelectModel(m)}
          >
            <div className="model-header-row">
              <span className="model-name">{m.name}</span>
              <span className="model-task-badge">{m.task}</span>
            </div>

            <div style={{ fontSize: '0.72rem', color: '#94a3b8', minHeight: '32px' }}>
              {m.desc}
            </div>

            <div className="model-stat-grid">
              <div className="model-stat-box">
                <div className="model-stat-val" style={{ color: 'var(--cb-green)' }}>{m.map5095}%</div>
                <div className="model-stat-lbl">mAP 50-95</div>
              </div>
              <div className="model-stat-box">
                <div className="model-stat-val">{m.params}M</div>
                <div className="model-stat-lbl">PARAMS</div>
              </div>
              <div className="model-stat-box">
                <div className="model-stat-val" style={{ color: 'var(--cb-yellow)' }}>{m.trtLatency}ms</div>
                <div className="model-stat-lbl">RTX 5090</div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
