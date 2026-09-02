import React, { useState } from 'react';
import ModelZooFilterBar from '../components/modelzoo/ModelZooFilterBar';
import ModelCardGrid from '../components/modelzoo/ModelCardGrid';
import ModelDetailsDrawer from '../components/modelzoo/ModelDetailsDrawer';
import ModelDistillationCard from '../components/modelzoo/ModelDistillationCard';

const ZOO_MODELS = [
  { id: 'yolo26n', name: 'YOLO26 Nano', family: 'YOLO26', task: 'DETECT', desc: 'End-to-End NMS-Free ultra-fast edge detector.', map5095: 39.8, params: 2.4, flops: 6.2, trtLatency: 0.38, nmsFree: true, depth: 0.33, width: 0.25 },
  { id: 'yolo26s', name: 'YOLO26 Small', family: 'YOLO26', task: 'DETECT', desc: 'Balanced accuracy and latency for smart city cameras.', map5095: 46.2, params: 9.8, flops: 21.5, trtLatency: 0.65, nmsFree: true, depth: 0.33, width: 0.50 },
  { id: 'yolo26m', name: 'YOLO26 Medium', family: 'YOLO26', task: 'DETECT', desc: 'High-precision vehicle and pedestrian specialist.', map5095: 51.4, params: 20.4, flops: 68.0, trtLatency: 1.12, nmsFree: true, depth: 0.67, width: 0.75 },
  { id: 'yolo26x', name: 'YOLO26 XLarge', family: 'YOLO26', task: 'DETECT', desc: 'Oracle teacher checkpoint for maximum mAP benchmark.', map5095: 56.1, params: 58.2, flops: 195.0, trtLatency: 2.30, nmsFree: true, depth: 1.00, width: 1.25 },
  { id: 'yolo26n-seg', name: 'YOLO26n Segment', family: 'YOLO26', task: 'SEGMENT', desc: 'Real-time instance segmentation with pixel masks.', map5095: 35.6, params: 3.1, flops: 9.8, trtLatency: 0.55, nmsFree: true, depth: 0.33, width: 0.25 },
  { id: 'yolo26n-pose', name: 'YOLO26n Pose', family: 'YOLO26', task: 'POSE', desc: '17 Human keypoint posture and ergonomics tracker.', map5095: 52.4, params: 3.3, flops: 9.2, trtLatency: 0.52, nmsFree: true, depth: 0.33, width: 0.25 },
  { id: 'yolo26n-obb', name: 'YOLO26n OBB', family: 'YOLO26', task: 'OBB', desc: 'Oriented bounding boxes for aerial and drone footage.', map5095: 41.2, params: 2.8, flops: 7.4, trtLatency: 0.48, nmsFree: true, depth: 0.33, width: 0.25 },
  { id: 'yolo11n', name: 'YOLO11 Nano', family: 'YOLO11', task: 'DETECT', desc: 'Standard decoupled head with NMS postprocess.', map5095: 39.5, params: 2.6, flops: 6.5, trtLatency: 0.45, nmsFree: false, depth: 0.33, width: 0.25 }
];

export default function ModelZooStudio() {
  const [selectedFamily, setSelectedFamily] = useState('ALL');
  const [selectedTask, setSelectedTask] = useState('ALL');
  const [selectedModel, setSelectedModel] = useState(ZOO_MODELS[0]);

  const filteredModels = ZOO_MODELS.filter(m => {
    const famMatch = selectedFamily === 'ALL' || m.family === selectedFamily;
    const taskMatch = selectedTask === 'ALL' || m.task === selectedTask;
    return famMatch && taskMatch;
  });

  const navigateTo = (tab) => { window.location.hash = tab; };

  return (
    <div className="view-container zoo-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">MODEL ZOO REPOSITORY</h1>
        <p className="cockpit-main-subtitle">
          OFFICIAL CHECKPOINTS • DUAL-HEAD YOLO26 ARCHITECTURES • KNOWLEDGE DISTILLATION & BENCHMARKS
        </p>
      </div>

      <ModelZooFilterBar
        selectedFamily={selectedFamily}
        setSelectedFamily={setSelectedFamily}
        selectedTask={selectedTask}
        setSelectedTask={setSelectedTask}
        totalCount={filteredModels.length}
      />

      <div className="zoo-layout">
        <div>
          <ModelCardGrid
            models={filteredModels}
            selectedModel={selectedModel}
            onSelectModel={setSelectedModel}
          />
          <ModelDistillationCard onLaunchDistill={() => navigateTo('cockpit')} />
        </div>

        <div>
          <ModelDetailsDrawer
            model={selectedModel}
            onSendToCockpit={() => navigateTo('cockpit')}
            onSendToBenchmark={() => navigateTo('benchmarks')}
            onSendToPlayground={() => navigateTo('playground')}
          />
        </div>
      </div>
    </div>
  );
}
