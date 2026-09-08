<script setup lang="ts">
import { computed } from 'vue'
import type { DatasetInfo } from '../types/dataset'

const props = withDefaults(defineProps<{
  datasets?: DatasetInfo[]
  selectedDataset?: string
  scale?: string
  batchSize?: number
  imgsz?: number
}>(), {
  datasets: () => [],
  selectedDataset: '',
  scale: 'n',
  batchSize: 32,
  imgsz: 640
})

const emit = defineEmits<{
  (e: 'update:selectedDataset', val: string): void
}>()

const estimatedVRAM = computed(() => {
  const scaleMultiplier = props.scale === 'n' ? 1.0 : props.scale === 's' ? 1.8 : props.scale === 'm' ? 3.5 : props.scale === 'l' ? 5.5 : 8.0
  const batchMultiplier = props.batchSize === -1 ? 16 : props.batchSize
  const resMultiplier = (props.imgsz / 640) * (props.imgsz / 640)
  return Math.min(31.5, 0.8 + (scaleMultiplier * batchMultiplier * 0.04 * resMultiplier)).toFixed(1)
})
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">3. HARDWARE & ESTIMATOR</span>
      <span class="badge-green">NVIDIA RTX 5090</span>
    </div>

    <div class="selector-group">
      <div class="selector-label">TARGET DATASET</div>
      <select
        class="cyber-select"
        :value="selectedDataset"
        @change="(e) => emit('update:selectedDataset', (e.target as HTMLSelectElement).value)"
      >
        <option v-if="datasets.length === 0" value="">
          [EMPTY] No dataset imported (import via DATASET STUDIO)
        </option>
        <option
          v-for="d in datasets"
          :key="d.id"
          :value="d.id"
        >
          {{ d.name }} ({{ d.classes?.length || 0 }} classes // {{ d.train_count || 0 }} train / {{ d.val_count || 0 }} val)
        </option>
      </select>
    </div>

    <div class="telemetry-box" style="background: rgba(0,0,0,0.4); padding: 1rem; border-radius: 4px; border: 1px solid rgba(0, 240, 255, 0.15);">
      <div style="display: flex; justify-content: space-between; margin-bottom: 0.5rem;">
        <span style="font-size: 0.8rem; color: #94a3b8;">ESTIMATED VRAM PEAK:</span>
        <span style="font-family: var(--font-mono); color: var(--cb-yellow); font-weight: 700;">
          ~{{ estimatedVRAM }} GB / 32.0 GB
        </span>
      </div>
      <div style="display: flex; justify-content: space-between; margin-bottom: 0.5rem;">
        <span style="font-size: 0.8rem; color: #94a3b8;">HARDWARE DEVICE:</span>
        <span style="font-family: var(--font-mono); color: var(--cb-green); font-weight: 700;">
          cuda:0 (NVIDIA RTX 5090)
        </span>
      </div>
      <div style="display: flex; justify-content: space-between;">
        <span style="font-size: 0.8rem; color: #94a3b8;">ZERO-COPY EXPORT:</span>
        <span style="font-family: var(--font-mono); color: var(--cb-cyan); font-weight: 700;">
          TensorRT 10.x Engine
        </span>
      </div>
    </div>
  </div>
</template>
