<script setup lang="ts">
import DirectoryFolderCard from '../components/DirectoryFolderCard.vue'

const DIRECTORIES = [
  { id: 'cockpit', code: 'DIR_01', label: 'TRAINING COCKPIT', tag: 'LAUNCHER', desc: 'Configure YOLOv8, YOLO11, YOLO26 architectures, hyperparameters, AMP FP16/BF16, and launch training loops.', stats: '5 ARCHITECTURES', status: 'ACTIVE' },
  { id: 'live-hud', code: 'DIR_02', label: 'LIVE TELEMETRY HUD', tag: 'METRICS', desc: 'Realtime WebSocket HUD with live Bézier loss curves (box/cls/dfl), mAP50 precision, and RTX 5090 sensors.', stats: 'WS STREAMING', status: 'STREAM' },
  { id: 'benchmarks', code: 'DIR_03', label: 'BENCHMARK STUDIO', tag: 'SPEED_MAP', desc: 'Ultralytics multi-format speed evaluation. Compare PyTorch vs TensorRT (5.5x) vs ONNX across FPS, latency, and mAP.', stats: '11 RUNTIMES', status: 'OPTIMIZED' },
  { id: 'datasets', code: 'DIR_04', label: 'DATASET STUDIO', tag: 'DATA_YAML', desc: 'Validate data.yaml configs, inspect class distribution bar charts for bias detection, and manage train/val splits.', stats: 'COCO / YOLO', status: 'READY' },
  { id: 'model-zoo', code: 'DIR_05', label: 'MODEL ZOO & EXPORT', tag: 'TENSORRT', desc: 'Manage checkpoints (best.pt vs last.pt), compare accuracy retention, and compile 1-click TensorRT .engine models.', stats: '1-CLICK TRT', status: 'EXPORT' },
  { id: 'playground', code: 'DIR_06', label: 'KIROSHI PLAYGROUND', tag: 'INFERENCE', desc: 'Drag & drop image/video inference testing with real-time confidence and IoU threshold sliders powered by PyTorch.', stats: 'LIVE TEST', status: 'OPTIC' }
]

defineProps<{ gpuStats?: any }>()
const emit = defineEmits<{ (e: 'navigate', tabId: string): void }>()
</script>

<template>
  <div class="hub-container">
    <div class="hub-hero-banner">
      <div>
        <div class="hub-title">HYDRAFORGE CYBERDECK OS</div>
        <div class="hub-subtitle">// NEURAL STUDIO SYSTEM ARCHITECTURE • SELECT A DIRECTORY MODULE TO MOUNT</div>
      </div>
      <div class="header-status-group">
        <div class="header-metric-pill">
          <span class="pill-label">GPU:</span>
          <span class="pill-val cyan">{{ gpuStats?.model || gpuStats?.name || 'RTX 5090' }}</span>
        </div>
        <div class="header-metric-pill">
          <span class="pill-label">VRAM:</span>
          <span class="pill-val yellow">{{ (gpuStats?.total_vram_mb || gpuStats?.vram_total_mb || 32607).toFixed(0) }} MB</span>
        </div>
      </div>
    </div>

    <div class="folder-grid">
      <DirectoryFolderCard v-for="dir in DIRECTORIES" :key="dir.id" :folder="dir" @open="(id) => emit('navigate', id)" />
    </div>
  </div>
</template>
