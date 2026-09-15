<script setup lang="ts">
import { ref, computed } from 'vue'

const props = withDefaults(defineProps<{
  config: Record<string, any>
  modelsList?: any[]
  hydraStreams?: any[]
  isHydraOnline?: boolean
  isRunning?: boolean
  isContinuous?: boolean
  telemetry?: any
  gpuStats?: any
  hydraTelemetry?: any
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

const activeTab = ref<'source' | 'tuning' | 'telemetry'>('source')
const customModels = computed(() => props.modelsList.filter(m => m.isCustom))
const baseModels = computed(() => props.modelsList.filter(m => !m.isCustom))
const updateField = (field: string, val: any) => emit('update:config', { ...props.config, [field]: val })
</script>

<template>
  <div class="cyber-card playground-unified-panel">
    <div class="panel-tab-nav">
      <button class="panel-tab-btn" :class="{ active: activeTab === 'source' }" @click="activeTab = 'source'">
        1. MODEL & SOURCE
      </button>
      <button class="panel-tab-btn" :class="{ active: activeTab === 'tuning' }" @click="activeTab = 'tuning'">
        2. HYPER-TUNING
      </button>
      <button class="panel-tab-btn" :class="{ active: activeTab === 'telemetry' }" @click="activeTab = 'telemetry'">
        3. HARDWARE & SHM
      </button>
    </div>

    <!-- TAB 1: SOURCE & MODEL -->
    <div v-if="activeTab === 'source'" class="tab-pane">
      <div class="selector-group">
        <div class="selector-label">
          <span>YOLO MODEL ARCHITECTURE</span>
          <span style="font-size: 0.65rem; color: var(--cb-yellow);">{{ customModels.length }} CUSTOM TRAINED</span>
        </div>
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
        <div class="selector-label">
          <span>INPUT MEDIA SOURCE</span>
          <span style="font-size: 0.65rem; color: var(--cb-cyan);">HARDWARE // CAMERAS</span>
        </div>
        <select class="cyber-select" :value="config.source" @change="(e) => updateField('source', (e.target as HTMLSelectElement).value)">
          <optgroup label="[LOCAL HARDWARE DEVICES]">
            <option value="webcam">[WEBCAM] LOCAL WEBCAM // /dev/video0</option>
          </optgroup>
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

    <!-- TAB 2: TUNING & FILTERS -->
    <div v-else-if="activeTab === 'tuning'" class="tab-pane">
      <div class="selector-group">
        <div class="selector-label">
          <span>CONFIDENCE THRESHOLD</span>
          <span class="slider-val">{{ (config.conf * 100).toFixed(0) }}%</span>
        </div>
        <div class="slider-row">
          <input type="range" min="0.05" max="0.95" step="0.05" :value="config.conf" @input="(e) => updateField('conf', parseFloat((e.target as HTMLInputElement).value))" />
        </div>
      </div>

      <div class="selector-group">
        <div class="selector-label">
          <span>IOU NMS THRESHOLD</span>
          <span class="slider-val">{{ (config.iou * 100).toFixed(0) }}%</span>
        </div>
        <div class="slider-row">
          <input type="range" min="0.10" max="0.90" step="0.05" :value="config.iou" :disabled="config.nmsFree" @input="(e) => updateField('iou', parseFloat((e.target as HTMLInputElement).value))" />
        </div>
        <div style="display: flex; align-items: center; gap: 0.5rem; margin-top: 0.4rem;">
          <input id="nmsFreeCheckUnified" type="checkbox" :checked="config.nmsFree" @change="(e) => updateField('nmsFree', (e.target as HTMLInputElement).checked)" />
          <label for="nmsFreeCheckUnified" style="font-size: 0.75rem; color: #cbd5e1; cursor: pointer;">NMS-FREE END-TO-END OUTPUT (YOLO26)</label>
        </div>
      </div>

      <div class="selector-group" style="border-top: 1px solid rgba(0,240,255,0.15); padding-top: 0.65rem;">
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <div style="font-family: var(--font-oxanium); font-size: 0.82rem; font-weight: 700;" :style="{ color: config.sahi ? 'var(--cb-cyan)' : '#fff' }">SAHI 4K SLICED INFERENCE</div>
            <div style="font-size: 0.68rem; color: #94a3b8;">Fatiamento dinâmico para alvos pequenos</div>
          </div>
          <button type="button" class="cyber-pill" :class="{ active: config.sahi }" @click="updateField('sahi', !config.sahi)">{{ config.sahi ? 'ON' : 'OFF' }}</button>
        </div>
      </div>

      <div class="selector-group" style="border-top: 1px solid rgba(0,240,255,0.15); padding-top: 0.65rem; margin-bottom: 0;">
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <div style="font-family: var(--font-oxanium); font-size: 0.82rem; font-weight: 700;" :style="{ color: config.isolateBg ? 'var(--cb-green)' : '#fff' }">ALPHA BACKGROUND REMOVER (PNG)</div>
            <div style="font-size: 0.68rem; color: #94a3b8;">Máscaras transparentes por instância</div>
          </div>
          <button type="button" class="cyber-pill" :class="{ active: config.isolateBg }" @click="updateField('isolateBg', !config.isolateBg)">{{ config.isolateBg ? 'ON' : 'OFF' }}</button>
        </div>
      </div>
    </div>

    <!-- TAB 3: TELEMETRY & HARDWARE -->
    <div v-else class="tab-pane">
      <div class="telemetry-row">
        <span class="k">GPU CORE & TEMP</span>
        <span class="v" style="color: var(--cb-green);">{{ gpuStats ? `${gpuStats.gpu_util_pct || gpuStats.utilization_pct || 0}% • ${gpuStats.temp_celsius || gpuStats.temp_c || 0}°C (${gpuStats.power_watts || gpuStats.power_w || 0}W)` : 'STANDBY' }}</span>
      </div>
      <div class="telemetry-row">
        <span class="k">VRAM ALLOCATION</span>
        <span class="v" style="color: var(--cb-cyan);">{{ gpuStats ? `${(gpuStats.used_vram_mb || gpuStats.vram_used_mb || 0).toLocaleString()} / ${(gpuStats.total_vram_mb || gpuStats.vram_total_mb || 32607).toLocaleString()} MB` : 'STANDBY' }}</span>
      </div>
      <div class="telemetry-row">
        <span class="k">YOLO INFERENCE (RTX 5090)</span>
        <span class="v" style="color: var(--cb-yellow);">{{ telemetry?.inference_ms ? `${telemetry.inference_ms} ms` : 'TRIGGER SCAN' }}</span>
      </div>
      <div class="telemetry-row">
        <span class="k">SUSTAINED SPEED</span>
        <span class="v" style="color: var(--cb-yellow); font-size: 0.85rem;">{{ telemetry?.fps ? `${telemetry.fps} FPS` : 'STANDBY' }}</span>
      </div>
      <div class="telemetry-row" style="border-bottom: none;">
        <span class="k">POSIX /dev/shm OCCUPANCY</span>
        <span class="v">{{ hydraTelemetry?.posix_shm_occupancy !== undefined ? `${hydraTelemetry.posix_shm_occupancy}%` : '0.1%' }}</span>
      </div>
    </div>
  </div>
</template>
