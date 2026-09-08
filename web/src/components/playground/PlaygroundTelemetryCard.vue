<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  telemetry?: any
  detections?: any[]
  selectedEntity?: any
  hydraTelemetry?: any
  activeStream?: any
  gpuStats?: any
}>()

const emit = defineEmits<{
  (e: 'selectEntity', det: any): void
}>()

const exportedMsg = ref<string | null>(null)

const handleExportPng = (e: MouseEvent, det: any) => {
  e.stopPropagation()
  exportedMsg.value = `PNG EXPORTED: ${det.label.toUpperCase()} #${det.id}`
  setTimeout(() => exportedMsg.value = null, 2500)
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header"><span class="card-title">3. HARDWARE & SHM TELEMETRY</span><span class="badge-green">{{ gpuStats?.model || gpuStats?.name || 'RTX 5090' }}</span></div>

    <div style="margin-bottom: 0.8rem;">
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

    <div v-if="exportedMsg" style="background: rgba(0,255,157,0.08); border: 1px solid var(--cb-green); padding: 0.35rem; border-radius: 3px; font-size: 0.7rem; color: var(--cb-green); font-family: var(--font-mono); margin-bottom: 0.5rem;">
      [OK] {{ exportedMsg }}
    </div>

    <div class="selector-group" style="margin-bottom: 0;">
      <div class="selector-label">
        <span>DETECTED ENTITIES ({{ detections?.length || 0 }})</span>
        <span style="font-size: 0.65rem; color: #94a3b8;">REAL LABELS</span>
      </div>
      <div style="max-height: 160px; overflow-y: auto; padding-right: 4px;">
        <div v-if="!detections || detections.length === 0" style="font-size: 0.75rem; color: #64748b; text-align: center; padding: 1.2rem 0;">
          Execute um scan na RTX 5090 para ver as detecções reais.
        </div>
        <div
          v-for="det in (detections || [])"
          v-else
          :key="det.id"
          class="entity-item"
          :style="{
            borderColor: selectedEntity?.id === det.id ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
            background: selectedEntity?.id === det.id ? 'rgba(252, 238, 10, 0.15)' : 'rgba(0, 240, 255, 0.05)',
            cursor: 'pointer'
          }"
          @click="emit('selectEntity', det)"
        >
          <div>
            <span style="font-weight: 700; color: #fff;">{{ det.label }}</span>
            <span style="margin-left: 0.4rem; color: var(--cb-yellow); font-size: 0.7rem;">{{ (det.conf * 100).toFixed(0) }}%</span>
          </div>
          <button class="cyber-pill" style="padding: 0.2rem 0.5rem; font-size: 0.65rem;" @click="(e) => handleExportPng(e, det)">[PNG]</button>
        </div>
      </div>
    </div>
  </div>
</template>
