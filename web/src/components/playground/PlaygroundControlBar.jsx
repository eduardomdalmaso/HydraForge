import React from 'react';

export default function PlaygroundControlBar({
  config, setConfig, modelsList = [], hydraStreams = [],
  isHydraOnline, onRunInference, isRunning, isContinuous, setIsContinuous
}) {
  const customModels = modelsList.filter(m => m.isCustom);
  const baseModels = modelsList.filter(m => !m.isCustom);

  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. MODEL & INPUT SOURCE</span>
        <span className={isHydraOnline ? 'badge-green' : 'badge-yellow'}>
          {isHydraOnline ? 'HYDRASTREAM LINKED' : 'STANDALONE'}
        </span>
      </div>

      <div className="selector-group">
        <div className="selector-label">
          <span>YOLO MODEL ARCHITECTURE</span>
          <span style={{ fontSize: '0.65rem', color: 'var(--cb-yellow)' }}>{customModels.length} CUSTOM TRAINED</span>
        </div>
        <select className="cyber-select" value={config.model} onChange={(e) => setConfig({ ...config, model: e.target.value })}>
          {customModels.length > 0 && (
            <optgroup label="[CUSTOM TRAINED // RTX 5090]">
              {customModels.map((m) => (
                <option key={m.id} value={m.id}>[TRAINED] {m.name.toUpperCase()} // mAP {m.map5095?.toFixed(1)}%</option>
              ))}
            </optgroup>
          )}
          <optgroup label="[BASE ARCHITECTURES // ULTRALYTICS]">
            {baseModels.map((m) => (
              <option key={m.id} value={m.id}>[BASE] {m.name.toUpperCase()} // {m.task}</option>
            ))}
          </optgroup>
        </select>
      </div>

      <div className="selector-group">
        <div className="selector-label">
          <span>INPUT MEDIA SOURCE</span>
          <span style={{ fontSize: '0.65rem', color: 'var(--cb-cyan)' }}>HARDWARE // CAMERAS</span>
        </div>
        <select className="cyber-select" value={config.source} onChange={(e) => setConfig({ ...config, source: e.target.value })}>
          <optgroup label="[LOCAL HARDWARE DEVICES]">
            <option value="webcam">[WEBCAM] LOCAL WEBCAM // /dev/video0</option>
          </optgroup>
          <optgroup label="[HYDRASTREAM ZERO-COPY // /dev/shm]">
            {hydraStreams.length > 0 ? (
              hydraStreams.map((s) => (
                <option key={s.stream_id} value={s.stream_id}>[STREAM] {s.stream_id.toUpperCase()} // {s.resolution || '1080P'} @ {s.ingest_fps || 30} FPS</option>
              ))
            ) : (
              <option value="cam_entrance_01">[STREAM] CAM_ENTRANCE_01 // 1080P @ 30 FPS</option>
            )}
          </optgroup>
          <optgroup label="[HYDRAVAULT ARCHIVE CAMERAS]">
            <option value="vault_cam_active" disabled>[HYDRAVAULT] NO VAULT CAMERAS REGISTERED (STANDBY)</option>
          </optgroup>
        </select>
      </div>

      <div className="selector-group">
        <div className="selector-label">RUNTIME ENGINE</div>
        <select className="cyber-select" value={config.runtime} onChange={(e) => setConfig({ ...config, runtime: e.target.value })}>
          <option value="pytorch">PYTORCH CUDA 13.3 // RTX 5090 DIRECT</option>
          <option value="tensorrt">TENSORRT 10.X // ZERO-LATENCY ENGINE</option>
        </select>
      </div>

      <div style={{ display: 'flex', gap: '0.5rem', marginTop: '0.5rem' }}>
        <button className="cyber-action-btn" style={{ flex: 1, padding: '0.75rem' }} onClick={onRunInference} disabled={isRunning}>
          {isRunning ? 'SCANNING...' : 'SCAN FRAME // TRIGGER'}
        </button>
        <button
          className="cyber-action-btn"
          style={{
            padding: '0.75rem',
            background: isContinuous ? 'var(--cb-green)' : 'rgba(0,240,255,0.1)',
            color: isContinuous ? '#07080c' : 'var(--cb-cyan)',
            border: '1px solid var(--cb-cyan)'
          }}
          onClick={() => setIsContinuous(!isContinuous)}
          title="Toggle Continuous Realtime HUD Scanner"
        >
          {isContinuous ? 'LIVE SCAN [ACTIVE]' : 'LIVE SCAN [IDLE]'}
        </button>
      </div>
    </div>
  );
}
