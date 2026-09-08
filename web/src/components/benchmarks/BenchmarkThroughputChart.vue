<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  results?: any[]
}>(), {
  results: () => []
})

const valid = computed(() => props.results?.filter(r => r.status === 'SUCCESS') || [])
const maxFPS = computed(() => valid.value.length > 0 ? Math.max(...valid.value.map(r => r.fps), 100) : 100)

const getFillColor = (format: string) => {
  const f = format.toLowerCase()
  if (f.includes('tensorrt') || f.includes('engine')) return 'var(--cb-green)'
  if (f.includes('onnx')) return 'var(--cb-cyan)'
  return 'var(--cb-yellow)'
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">THROUGHPUT COMPARISON (FPS)</span>
      <span class="badge-cyan">{{ valid.length > 0 ? 'REAL BENCHMARK DATA' : 'STANDBY' }}</span>
    </div>
    <div style="margin-top: 0.5rem;">
      <div v-if="valid.length === 0" style="text-align: center; padding: 2rem 1rem; color: #94a3b8; font-size: 0.8rem;">
        No benchmarks executed yet. Select runtime formats and execute the speed test.
      </div>
      <div
        v-for="(r, idx) in valid"
        v-else
        :key="idx"
        class="bar-chart-row"
      >
        <div class="bar-chart-label">{{ r.format }}</div>
        <div class="bar-chart-track">
          <div
            class="bar-chart-fill"
            :style="{
              width: `${Math.min(100, Math.max(6, (r.fps / maxFPS) * 100))}%`,
              background: getFillColor(r.format)
            }"
          />
        </div>
        <div class="bar-chart-val" :style="{ color: getFillColor(r.format) }">
          {{ r.fps?.toFixed(0) || 0 }} <span style="font-size: 0.65rem;">FPS</span>
        </div>
      </div>
    </div>
  </div>
</template>
