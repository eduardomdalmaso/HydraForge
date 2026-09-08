<script setup lang="ts">
import { ref } from 'vue'
import BenchmarkFormatGrid from './BenchmarkFormatGrid.vue'

const FORMATS = [
  { id: 'TensorRT', label: 'TensorRT (.engine)', gpuOnly: true },
  { id: 'PyTorch', label: 'PyTorch' },
  { id: 'ONNX', label: 'ONNX DirectML' },
  { id: 'OpenVINO', label: 'OpenVINO' },
  { id: 'TorchScript', label: 'TorchScript' },
  { id: 'LiteRT', label: 'LiteRT (TFLite)' }
]
const BASE_MODELS = ['yolo26n.pt', 'yolo26s.pt', 'yolo26m.pt', 'yolo26l.pt', 'yolov8n.pt', 'yolo11n.pt']

defineProps<{ isRunning?: boolean; jobs?: any[]; datasets?: any[]; gpuStats?: any }>()
const emit = defineEmits<{ (e: 'launch', payload: any): void }>()

const model = ref('yolo26n.pt')
const data = ref('coco8.yaml')
const imgsz = ref(640)
const quantize = ref(16)
const device = ref('0')
const formats = ref(['TensorRT', 'PyTorch', 'ONNX', 'OpenVINO'])

const onDeviceChange = (val: string) => {
  device.value = val
  if (val === 'cpu') formats.value = formats.value.filter(f => f !== 'TensorRT')
  else if (!formats.value.includes('TensorRT')) formats.value = ['TensorRT', ...formats.value]
}

const toggle = (id: string, gpuOnly?: boolean) => {
  if (device.value === 'cpu' && gpuOnly) return
  formats.value = formats.value.includes(id) ? formats.value.filter(f => f !== id) : [...formats.value, id]
}

const handleSubmit = () => {
  emit('launch', { model: model.value, data: data.value, imgsz: Number(imgsz.value), quantize: Number(quantize.value), device: device.value, target_formats: formats.value })
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header"><span class="card-title">1. BENCHMARK LAUNCHER</span><span class="badge-cyan">{{ device === 'cpu' ? 'CPU HOST' : `${gpuStats?.model || 'RTX 5090'} // CUDA 13.3` }}</span></div>
    <form style="margin-top: 0.65rem;" @submit.prevent="handleSubmit">
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.65rem;">
        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label">MODEL CHECKPOINT</div>
          <select class="cyber-select" :value="model" @change="(e) => model = (e.target as HTMLSelectElement).value">
            <optgroup label="Trained Models">
              <option v-for="j in jobs" :key="j.job_id" :value="j.output_weights || j.model_architecture">[TRAINED] {{ j.model_architecture }} ({{ j.dataset_id }}){{ j.best_map50 ? ` - ${(j.best_map50 * 100).toFixed(1)}%` : '' }}</option>
            </optgroup>
            <optgroup label="Base Pretrained Models"><option v-for="m in BASE_MODELS" :key="m" :value="m">{{ m }}</option></optgroup>
          </select>
        </div>
        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label">VALIDATION DATASET</div>
          <select class="cyber-select" :value="data" @change="(e) => data = (e.target as HTMLSelectElement).value">
            <optgroup label="Workspace Datasets">
              <option v-for="d in datasets" :key="d.id || d.dataset_id" :value="d.yaml_path || d.id || d.dataset_id">[DATASET] {{ d.name || d.dataset_id }} ({{ d.classes?.length || 1 }} CLS)</option>
            </optgroup>
            <optgroup label="Reference Datasets"><option value="coco8.yaml">[DEMO] coco8.yaml</option></optgroup>
          </select>
        </div>
      </div>

      <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 0.65rem; margin-top: 0.65rem;">
        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label">RESOLUTION</div>
          <select class="cyber-select" :value="imgsz" @change="(e) => imgsz = Number((e.target as HTMLSelectElement).value)">
            <option :value="320">320x320</option><option :value="640">640x640</option><option :value="1280">1280x1280</option>
          </select>
        </div>
        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label">PRECISION</div>
          <select class="cyber-select" :value="quantize" @change="(e) => quantize = Number((e.target as HTMLSelectElement).value)">
            <option :value="16">FP16 (TensorRT)</option><option :value="8">INT8 (Quantized)</option><option :value="32">FP32 (Float)</option>
          </select>
        </div>
        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label">DEVICE</div>
          <select class="cyber-select" :value="device" @change="(e) => onDeviceChange((e.target as HTMLSelectElement).value)">
            <option value="0">GPU 0: {{ gpuStats?.model || 'RTX 5090' }}</option><option value="cpu">CPU: HOST MULTI-CORE</option>
          </select>
        </div>
      </div>

      <BenchmarkFormatGrid :formats="formats" :formatList="FORMATS" :device="device" @toggle="toggle" />

      <button type="submit" :disabled="isRunning || formats.length === 0" class="cyber-action-btn" style="width: 100%; margin-top: 1rem; padding: 0.75rem;">
        {{ isRunning ? 'BENCHMARKING...' : 'BENCHMARK' }}
      </button>
    </form>
  </div>
</template>
