import React from 'react';

export default function ModelDetailsDrawer({ model, onSendToCockpit, onSendToBenchmark, onSendToPlayground }) {
  if (!model) return null;

  const yamlConfig = `# Ultralytics Model Architecture (${model.id}.yaml)
scales:
  ${model.id[6] || 'n'}: [${model.depth}, ${model.width}, 1024]

# Dual-Head End-to-End NMS-Free Backbone
backbone:
  - [-1, 1, Conv, [64, 3, 2]]       # 0-P1/2
  - [-1, 1, Conv, [128, 3, 2]]      # 1-P2/4
  - [-1, 2, C3k2, [256, False, 0.25]] # 2
  - [-1, 1, C2PSA, [512]]           # 9 Attention
head:
  - [-1, 1, Detect, [nc, [256, 512, 1024]]] # Dual One-to-One / One-to-Many`;

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">SPECIFICATION // {model.name}</span>
        <span className="badge-cyan">{model.task}</span>
      </div>

      <div style={{ marginBottom: '0.85rem' }}>
        <div className="telemetry-row">
          <span className="k">ARCHITECTURE HEAD</span>
          <span className="v" style={{ color: 'var(--cb-yellow)' }}>{model.nmsFree ? '⚡ NMS-Free End-to-End' : 'Standard Decoupled Head'}</span>
        </div>
        <div className="telemetry-row">
          <span className="k">FLOPs / COMPLEXITY</span>
          <span className="v">{model.flops} GFLOPs</span>
        </div>
        <div className="telemetry-row">
          <span className="k">TENSORRT LATENCY (RTX 5090)</span>
          <span className="v" style={{ color: 'var(--cb-cyan)' }}>{model.trtLatency} ms</span>
        </div>
        <div className="telemetry-row" style={{ borderBottom: 'none' }}>
          <span className="k">WEIGHT FILE</span>
          <span className="v" style={{ color: 'var(--cb-green)' }}>{model.id}.pt (PyTorch Checkpoint)</span>
        </div>
      </div>

      <pre className="yaml-code-box" style={{ maxHeight: '140px', fontSize: '0.72rem', marginBottom: '0.85rem' }}>
        <code>{yamlConfig}</code>
      </pre>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem', marginBottom: '0.5rem' }}>
        <button className="cyber-pill active" style={{ padding: '0.55rem', textAlign: 'center' }} onClick={() => onSendToCockpit(model)}>
          🚀 TRAIN IN COCKPIT
        </button>
        <button className="cyber-pill" style={{ padding: '0.55rem', textAlign: 'center' }} onClick={() => onSendToBenchmark(model)}>
          ⚡ BENCHMARK TRT
        </button>
      </div>

      <button className="cyber-pill" style={{ width: '100%', padding: '0.55rem', textAlign: 'center', borderColor: 'var(--cb-yellow)', color: 'var(--cb-yellow)' }} onClick={() => onSendToPlayground(model)}>
        👁️ TEST LIVE IN PLAYGROUND →
      </button>
    </div>
  );
}
