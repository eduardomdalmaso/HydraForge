<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  config: Record<string, any>
  modelsList?: any[]
  hydraStreams?: any[]
  isHydraOnline?: boolean
  isRunning?: boolean
  isContinuous?: boolean
}>(), {
  modelsList: () => [],
  hydraStreams: () => [],
  isHydraOnline: false,
  isRunning: false,
  isContinuous: false
})

const emit = defineEmits<{
  (e: 'update:config', cfg: Record<string, any>): void
  (e: 'update:isContinuous', val: boolean): void
  (e: 'runInference'): void
}>()

const customModels = computed(() => props.modelsList.filter(m => m.isCustom))
const baseModels = computed(() => props.modelsList.filter(m => !m.isCustom))
const updateField = (field: string, val: any) => emit('update:config', { ...props.config, [field]: val })
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">1. MODEL & INPUT SOURCE</span>
      <span :class="isHydraOnline ? 'badge-green' : 'badge-yellow'">{{ isHydraOnline ? 'HYDRASTREAM LINKED' : 'STANDALONE' }}</span>
    </div>

    <div class="selector-group">
      <div class="selector-label"><span>YOLO MODEL ARCHITECTURE</span><span style="font-size: 0.65rem; color: var(--cb-yellow);">{{ customModels.length }} CUSTOM TRAINED</span></div>
      <select class="cyber-select" :value="config.model" @change="(e) => updateField('model', (e.target as HTMLSelectElement).value)">
        <optgroup v-if="customModels.length > 0" label="[CUSTOM TRAINED // RTX 5090]">
          <option v-for="m in customModels" :key="m.id" :value="m.id">[TRAINED] {{ m.name.toUpperCase() }} // mAP {{ m.map5095?.toFixed(1) }}%</option>
        </optgroup>
        <optgroup label="[BASE ARCHITECTURES // ULTRALYTICS]">
          <option v-for="m in baseModels" :key="m.id" :value="m.id">[BASE] {{ m.name.toUpperCase() }} // {{ m.task }}</option>
        </optgroup>
      </select>
    </div>

    <div class="selector-group">
      <div class="selector-label"><span>INPUT MEDIA SOURCE</span><span style="font-size: 0.65rem; color: var(--cb-cyan);">HARDWARE // CAMERAS</span></div>
      <select class="cyber-select" :value="config.source" @change="(e) => updateField('source', (e.target as HTMLSelectElement).value)">
        <optgroup label="[LOCAL HARDWARE DEVICES]"><option value="webcam">[WEBCAM] LOCAL WEBCAM // /dev/video0</option></optgroup>
        <optgroup label="[HYDRASTREAM ZERO-COPY // /dev/shm]">
          <template v-if="hydraStreams.length > 0">
            <option v-for="s in hydraStreams" :key="s.stream_id" :value="s.stream_id">[STREAM] {{ s.stream_id.toUpperCase() }} // {{ s.resolution || '1080P' }} @ {{ s.ingest_fps || 30 }} FPS</option>
          </template>
          <option v-else value="cam_entrance_01">[STREAM] CAM_ENTRANCE_01 // 1080P @ 30 FPS</option>
        </optgroup>
      </select>
    </div>

    <div class="selector-group">
      <div class="selector-label">RUNTIME ENGINE</div>
      <select class="cyber-select" :value="config.runtime" @change="(e) => updateField('runtime', (e.target as HTMLSelectElement).value)">
        <option value="pytorch">PYTORCH CUDA 13.3 // RTX 5090 DIRECT</option>
        <option value="tensorrt">TENSORRT 10.X // ZERO-LATENCY ENGINE</option>
      </select>
    </div>

    <div style="display: flex; gap: 0.5rem; margin-top: 0.5rem;">
      <button class="cyber-action-btn" style="flex: 1; padding: 0.75rem;" :disabled="isRunning" @click="emit('runInference')">
        {{ isRunning ? 'SCANNING...' : 'SCAN FRAME // TRIGGER' }}
      </button>
      <button class="cyber-action-btn" :style="{ padding: '0.75rem', background: isContinuous ? 'var(--cb-green)' : 'rgba(0,240,255,0.1)', color: isContinuous ? '#07080c' : 'var(--cb-cyan)', border: '1px solid var(--cb-cyan)' }" title="Toggle Continuous Realtime HUD Scanner" @click="emit('update:isContinuous', !isContinuous)">
        {{ isContinuous ? 'LIVE SCAN [ACTIVE]' : 'LIVE SCAN [IDLE]' }}
      </button>
    </div>
  </div>
</template>
