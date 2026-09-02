import React from 'react';

export default function PlaygroundControlBar({ config, setConfig, onRunInference, isRunning }) {
  return (
    <div className="cyber-card">
      <div className="card-header">
        <span className="card-title">1. MODEL & INPUT SOURCE</span>
        <span className="badge-cyan">INFERENCE ENGINE</span>
      </div>

      <div className="selector-group">
        <div className="selector-label">YOLO MODEL ARCHITECTURE</div>
        <select
          className="cyber-select"
          value={config.model}
          onChange={(e) => setConfig({ ...config, model: e.target.value })}
        >
          <option value="yolo26n">⚡ YOLO26 Nano (End-to-End NMS-Free)</option>
          <option value="yolo26s">⚡ YOLO26 Small (Balanced)</option>
          <option value="yolo26m">⚡ YOLO26 Medium (High Precision)</option>
          <option value="yolo26x">⚡ YOLO26 XLarge (Max mAP)</option>
          <option value="yolo26n-seg">✂️ YOLO26n Segment (Instance Masks)</option>
          <option value="yolo26n-pose">🤸 YOLO26n Pose (17 Keypoints)</option>
          <option value="yolo26n-obb">📐 YOLO26n OBB (Oriented Bounding Boxes)</option>
          <option value="yolo11n">YOLO11 Nano</option>
          <option value="yolov8n">YOLOv8 Nano</option>
        </select>
      </div>

      <div className="selector-group">
        <div className="selector-label">INPUT MEDIA SOURCE</div>
        <select
          className="cyber-select"
          value={config.source}
          onChange={(e) => setConfig({ ...config, source: e.target.value })}
        >
          <option value="sample_bus">📷 Sample: Urban Traffic (Bus & Pedestrians)</option>
          <option value="sample_zidane">📷 Sample: Sports Motion (Zidane & Players)</option>
          <option value="rtsp_cam01">🎥 RTSP Live Stream (HydraStream /dev/shm)</option>
          <option value="webcam">📹 USB WebCam (/dev/video0)</option>
        </select>
      </div>

      <div className="selector-group">
        <div className="selector-label">RUNTIME PROVIDER & PREPROCESSING</div>
        <select
          className="cyber-select"
          value={config.runtime}
          onChange={(e) => setConfig({ ...config, runtime: e.target.value })}
        >
          <option value="tensorrt">🚀 TensorRT 10.x + NVIDIA DALI (RTX 5090)</option>
          <option value="rust">🦀 Rust ultralytics-inference (Zero-Copy SHM)</option>
          <option value="pytorch">PyTorch CUDA 13.3 (FP16)</option>
          <option value="openvino">Intel OpenVINO (Latency Mode)</option>
          <option value="onnx">ONNX Runtime (CPU/DirectML)</option>
        </select>
      </div>

      <button
        className="cyber-action-btn"
        style={{ width: '100%', padding: '0.75rem', marginTop: '0.5rem' }}
        onClick={onRunInference}
        disabled={isRunning}
      >
        {isRunning ? 'SCANNING...' : 'PREDICT'}
      </button>
    </div>
  );
}
