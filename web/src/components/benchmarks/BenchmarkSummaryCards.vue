<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  results?: any[]
}>(), {
  results: () => []
})

const valid = computed(() => props.results?.filter(r => r.status === 'SUCCESS') || [])
const hasResults = computed(() => valid.value.length > 0)

const bestFPS = computed(() => hasResults.value ? valid.value.reduce((max, r) => (r.fps > max ? r.fps : max), 0) : null)
const bestFormat = computed(() => hasResults.value ? (valid.value.find(r => r.fps === bestFPS.value)?.format || '-') : 'AWAITING BENCHMARK RUN')
const minLatency = computed(() => hasResults.value ? valid.value.reduce((min, r) => (r.inference_time_ms < min ? r.inference_time_ms : min), 9999) : null)

const ptResult = computed(() => valid.value.find(r => r.format?.toLowerCase().includes('pytorch')))
const ptFPS = computed(() => ptResult.value?.fps || (valid.value[0]?.fps || null))
const speedup = computed(() => (hasResults.value && ptFPS.value && ptFPS.value > 0 && bestFPS.value) ? (bestFPS.value / ptFPS.value).toFixed(1) : null)
const bestMAP = computed(() => hasResults.value ? valid.value.reduce((max, r) => (r.map50_95 > max ? r.map50_95 : max), 0) : null)
</script>

<template>
  <div class="telemetry-grid" style="margin-bottom: 0.5rem;">
    <div class="cyber-card metric-card">
      <div class="metric-title">MAX THROUGHPUT</div>
      <div class="metric-value" style="color: var(--cb-cyan);">
        {{ bestFPS ? `${bestFPS.toFixed(0)} FPS` : '-' }}
      </div>
      <div class="metric-subtitle">{{ bestFormat }}</div>
    </div>

    <div class="cyber-card metric-card">
      <div class="metric-title">MIN LATENCY (GPU)</div>
      <div class="metric-value" style="color: var(--cb-green);">
        {{ minLatency ? `${minLatency.toFixed(2)} ms` : '-' }}
      </div>
      <div class="metric-subtitle">Per Frame Inference (RTX 5090)</div>
    </div>

    <div class="cyber-card metric-card">
      <div class="metric-title">SPEEDUP MULTIPLIER</div>
      <div class="metric-value" style="color: var(--cb-yellow);">
        {{ speedup ? `${speedup}x` : '-' }}
      </div>
      <div class="metric-subtitle">vs PyTorch Native</div>
    </div>

    <div class="cyber-card metric-card">
      <div class="metric-title">PRECISION RETENTION</div>
      <div class="metric-value" style="color: var(--cb-magenta);">
        {{ bestMAP ? `${(bestMAP * 100).toFixed(1)}%` : '-' }}
      </div>
      <div class="metric-subtitle">mAP@50-95 Retention</div>
    </div>
  </div>
</template>
