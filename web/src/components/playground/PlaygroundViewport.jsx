import React from 'react';

export default function PlaygroundViewport({
  imageSrc, detections = [], selectedEntity, setSelectedEntity,
  isScanning, isContinuous, isHydraLinked, activeStream, isWebcam, videoRef
}) {
  return (
    <div className="viewport-hud">
      <div className="viewport-top-bar">
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
          <span>// KIROSHI OPTICAL HUD // REALTIME TRACKING</span>
          <span style={{
            fontSize: '0.65rem', padding: '2px 6px', borderRadius: '2px',
            background: isWebcam ? 'rgba(0,240,255,0.15)' : (isHydraLinked ? 'rgba(0,255,157,0.15)' : 'rgba(252,238,10,0.15)'),
            color: isWebcam ? 'var(--cb-cyan)' : (isHydraLinked ? 'var(--cb-green)' : 'var(--cb-yellow)'),
            border: `1px solid ${isWebcam ? 'var(--cb-cyan)' : (isHydraLinked ? 'var(--cb-green)' : 'var(--cb-yellow)')}`
          }}>
            {isWebcam ? '[HARDWARE] WEBCAM /dev/video0' : (isHydraLinked ? `[STREAM] ${activeStream?.stream_id?.toUpperCase() || 'SHM'}` : '[STANDBY]')}
          </span>
        </div>
        <span style={{ color: 'var(--cb-yellow)', fontFamily: 'var(--font-mono)' }}>
          {isContinuous ? '● REALTIME STREAM TRACKING' : (isScanning ? 'INFERENCE EXECUTING...' : `ACTIVE TARGETS: ${detections.length}`)}
        </span>
      </div>

      <div className="viewport-canvas-area" style={{ background: '#000', overflow: 'hidden' }}>
        <div style={{ position: 'relative', display: 'inline-block', maxWidth: '100%', width: '100%', textAlign: 'center' }}>
          {isWebcam ? (
            <video
              ref={videoRef}
              autoPlay
              playsInline
              muted
              className="viewport-img"
              style={{ display: 'block', maxWidth: '100%', width: '100%', height: 'auto', maxHeight: '520px', objectFit: 'contain', margin: '0 auto' }}
            />
          ) : (
            <img
              src={imageSrc || '/hydra-logo.jpg'}
              alt="Tracking Feed"
              className="viewport-img"
              style={{ display: 'block', maxWidth: '100%', maxHeight: '520px', objectFit: 'contain', margin: '0 auto' }}
              onError={(e) => { e.target.onerror = null; e.target.src = '/hydra-logo.jpg'; }}
            />
          )}

          {detections.map((det, idx) => {
            const isSelected = selectedEntity?.id === det.id;
            return (
              <div
                key={det.id || idx}
                className="detection-box"
                style={{
                  left: `${det.box[0]}%`, top: `${det.box[1]}%`, width: `${det.box[2]}%`, height: `${det.box[3]}%`,
                  borderColor: isSelected ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
                  boxShadow: isSelected ? '0 0 10px var(--cb-yellow)' : '0 0 6px rgba(0,240,255,0.4)',
                  transition: 'left 0.06s linear, top 0.06s linear, width 0.06s linear, height 0.06s linear'
                }}
                onClick={() => setSelectedEntity(det)}
              >
                <span className="detection-tag" style={{ background: isSelected ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'), color: '#07080c' }}>
                  {det.label.toUpperCase()} {(det.conf * 100).toFixed(0)}%
                </span>
              </div>
            );
          })}
        </div>
      </div>

      <div className="viewport-footer-hud">
        <span>FEED: {isWebcam ? 'USB V4L2 WEBCAM' : (activeStream?.stream_id?.toUpperCase() || 'CAM_ENTRANCE_01')}</span>
        <span>LATENCY: CUDA 13.3 (~3.8ms)</span>
        <span>RTX 5090 REALTIME TRACKER</span>
        <span style={{ color: 'var(--cb-green)' }}>8.46 GB/s ZERO-COPY</span>
      </div>
    </div>
  );
}
