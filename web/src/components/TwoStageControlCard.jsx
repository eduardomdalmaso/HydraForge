import React from 'react';

export default function TwoStageControlCard({ params, setParams }) {
  const isTwoStage = !!params.two_stage;

  const handleToggle = () => {
    setParams(prev => ({
      ...prev,
      two_stage: !prev.two_stage,
      stage1_epochs: prev.stage1_epochs || 20,
      stage1_freeze: prev.stage1_freeze || 10,
      stage2_epochs: prev.stage2_epochs || 30,
      close_mosaic: prev.close_mosaic || 10
    }));
  };

  const applyPreset = (presetKey) => {
    if (presetKey === 'yolo26') {
      setParams(prev => ({ ...prev, recipe_preset: 'yolo26_recipe', close_mosaic: 10, optimizer: 'AdamW', lr0: 0.00038, two_stage: false }));
    } else if (presetKey === 'small_data') {
      setParams(prev => ({ ...prev, recipe_preset: 'small_dataset', freeze: 10, patience: 20, epochs: 50, lr0: 0.001, two_stage: false }));
    } else if (presetKey === 'two_stage') {
      setParams(prev => ({ ...prev, recipe_preset: 'two_stage', two_stage: true, stage1_epochs: 20, stage1_freeze: 10, stage2_epochs: 30, close_mosaic: 10 }));
    } else {
      setParams(prev => ({ ...prev, recipe_preset: 'default', two_stage: false, freeze: 0 }));
    }
  };

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">3. RECIPES & TWO-STAGE TUNING</span>
        <span className="badge-green">TRANSFER LEARNING</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">RECIPE PRESETS</div>
        <div className="pills-container">
          <button type="button" className={`cyber-pill ${!isTwoStage && !params.freeze ? 'active' : ''}`} onClick={() => applyPreset('default')}>
            Standard
          </button>
          <button type="button" className={`cyber-pill ${params.recipe_preset === 'yolo26_recipe' ? 'active' : ''}`} onClick={() => applyPreset('yolo26')}>
            ⚡ YOLO26 Recipe
          </button>
          <button type="button" className={`cyber-pill ${params.freeze === 10 ? 'active' : ''}`} onClick={() => applyPreset('small_data')}>
            🎯 Small Data (Freeze 10)
          </button>
          <button type="button" className={`cyber-pill ${isTwoStage ? 'active' : ''}`} onClick={() => applyPreset('two_stage')}>
            🚀 Two-Stage
          </button>
        </div>
      </div>

      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', background: 'rgba(0,0,0,0.3)', padding: '0.65rem 0.8rem', borderRadius: '4px', border: isTwoStage ? '1px solid var(--cb-cyan)' : '1px solid rgba(255,255,255,0.08)' }}>
        <div>
          <div style={{ fontFamily: 'var(--font-oxanium)', fontSize: '0.85rem', fontWeight: '700', color: isTwoStage ? 'var(--cb-cyan)' : '#fff' }}>
            ⚡ TWO-STAGE FINE-TUNING PIPELINE
          </div>
          <div style={{ fontSize: '0.72rem', color: '#94a3b8' }}>
            Stage 1 (Freeze Backbone) ➔ Stage 2 (Global Unfreeze)
          </div>
        </div>
        <button type="button" className={`cyber-pill ${isTwoStage ? 'active' : ''}`} onClick={handleToggle} style={{ minWidth: '80px', textAlign: 'center' }}>
          {isTwoStage ? 'ACTIVE' : 'ENABLE'}
        </button>
      </div>

      {isTwoStage && (
        <div style={{ marginTop: '0.75rem', display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem', background: 'rgba(0,240,255,0.03)', padding: '0.75rem', border: '1px dashed rgba(0,240,255,0.3)', borderRadius: '4px' }}>
          <div>
            <div className="selector-label" style={{ fontSize: '0.75rem' }}>
              <span>STAGE 1: HEAD ADAPT</span>
              <span className="slider-val">{params.stage1_epochs || 20} ep</span>
            </div>
            <input type="range" min="5" max="50" value={params.stage1_epochs || 20} onChange={(e) => setParams(p => ({ ...p, stage1_epochs: parseInt(e.target.value) }))} style={{ width: '100%' }} />
            <div style={{ marginTop: '0.2rem', fontSize: '0.7rem', color: 'var(--cb-yellow)' }}>
              Freeze: <strong>Backbone (0-10)</strong>
            </div>
          </div>

          <div>
            <div className="selector-label" style={{ fontSize: '0.75rem' }}>
              <span>STAGE 2: FULL REFINE</span>
              <span className="slider-val">{params.stage2_epochs || 30} ep</span>
            </div>
            <input type="range" min="10" max="100" value={params.stage2_epochs || 30} onChange={(e) => setParams(p => ({ ...p, stage2_epochs: parseInt(e.target.value) }))} style={{ width: '100%' }} />
            <div style={{ marginTop: '0.2rem', fontSize: '0.7rem', color: 'var(--cb-green)' }}>
              Unfreeze: <strong>All • lr0=0.001</strong>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
