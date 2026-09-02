import React from 'react';

export default function PlaygroundViewport({ imageSrc, detections, selectedEntity, setSelectedEntity, isScanning }) {
  return (
    <div className="viewport-hud">
      <div className="viewport-top-bar">
        <span>⚡ KIROSHI OPTICAL HUD // LIVE SCANNER</span>
        <span style={{ color: 'var(--cb-yellow)' }}>
          {isScanning ? 'SCANNING FRAME...' : `ACTIVE DETECTIONS: ${detections.length}`}
        </span>
      </div>

      <div className="viewport-canvas-area">
        <div className="scanline-overlay"></div>
        <div className="hud-corner-tl"></div>
        <div className="hud-corner-br"></div>

        <div style={{ position: 'relative', display: 'inline-block' }}>
          <img
            src={imageSrc || '/hydra-logo.jpg'}
            alt="Inference Target"
            className="viewport-img"
          />

          {!isScanning && detections.map((det, idx) => {
            const isSelected = selectedEntity?.id === det.id;
            return (
              <div
                key={det.id || idx}
                className="detection-box"
                style={{
                  left: `${det.box[0]}%`,
                  top: `${det.box[1]}%`,
                  width: `${det.box[2]}%`,
                  height: `${det.box[3]}%`,
                  borderColor: isSelected ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
                  boxShadow: isSelected ? '0 0 15px var(--cb-yellow)' : 'none'
                }}
                onClick={() => setSelectedEntity(det)}
              >
                <span
                  className="detection-tag"
                  style={{
                    background: isSelected ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
                    color: '#07080c'
                  }}
                >
                  {det.label.toUpperCase()} {(det.conf * 100).toFixed(0)}%
                </span>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
