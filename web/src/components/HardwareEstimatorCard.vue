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
      <span class="card-title">3. HARDWARE & ESTIMADOR</span>
      <span class="badge-green">NVIDIA RTX 5090</span>
    </div>

    <div class="selector-group">
      <div class="selector-label">DATASET ALVO</div>
      <select
        class="cyber-select"
        :value="selectedDataset"
        @change="(e) => emit('update:selectedDataset', (e.target as HTMLSelectElement).value)"
      >
        <option v-if="datasets.length === 0" value="">
          [VAZIO] Nenhum dataset importado (importe via DATASETS)
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

    <div class="telemetry-box" style="background: var(--vms-bg-elevated); padding: 0.85rem; border-radius: var(--vms-radius-sm); border: 1px solid var(--vms-border);">
      <div style="display: flex; justify-content: space-between; margin-bottom: 0.45rem;">
        <span style="font-size: 0.775rem; color: var(--vms-text-muted);">PICO ESTIMADO DE VRAM:</span>
        <span style="font-family: var(--font-mono); color: var(--vms-primary); font-weight: 600; font-size: 0.8125rem;">
          ~{{ estimatedVRAM }} GB / 32.0 GB
        </span>
      </div>
      <div style="display: flex; justify-content: space-between; margin-bottom: 0.45rem;">
        <span style="font-size: 0.775rem; color: var(--vms-text-muted);">DISPOSITIVO:</span>
        <span style="font-family: var(--font-mono); color: var(--vms-success); font-weight: 600; font-size: 0.8125rem;">
          cuda:0 (NVIDIA RTX 5090)
        </span>
      </div>
      <div style="display: flex; justify-content: space-between;">
        <span style="font-size: 0.775rem; color: var(--vms-text-muted);">EXPORTACAO ZERO-COPY:</span>
        <span style="font-family: var(--font-mono); color: var(--vms-info); font-weight: 600; font-size: 0.8125rem;">
          TensorRT 10.x Engine
        </span>
      </div>
    </div>
  </div>
</template>
