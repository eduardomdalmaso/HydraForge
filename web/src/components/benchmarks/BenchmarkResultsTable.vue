<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  results?: any[]
}>(), {
  results: () => []
})

const hasResults = computed(() => props.results && props.results.length > 0)
</script>

<template>
  <div class="cyber-card" style="margin-top: 0.5rem;">
    <div class="card-header">
      <span class="card-title">EXPORT RUNTIME EVALUATION MATRIX</span>
      <span class="badge-cyan">{{ hasResults ? `${results.length} RUNTIMES RECORDED` : '0 RUNTIMES' }}</span>
    </div>
    <div style="padding: 0; overflow-x: auto; margin-top: 0.5rem;">
      <div v-if="!hasResults" style="text-align: center; padding: 2rem; color: #94a3b8; font-size: 0.8rem;">
        No export benchmark results recorded. Run a benchmark compilation above to evaluate latency and binary engine size metrics.
      </div>
      <table v-else class="cyber-table">
        <thead>
          <tr>
            <th>RUNTIME FORMAT</th>
            <th>STATUS</th>
            <th>ENGINE SIZE</th>
            <th>LATENCY (GPU)</th>
            <th>THROUGHPUT</th>
            <th>mAP 50-95</th>
            <th>COMPILATION ARGS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(r, idx) in results" :key="idx">
            <td>
              <span
                class="speedup-badge"
                :class="{ 'top-tier': r.format?.toLowerCase().includes('tensorrt') || r.format?.toLowerCase().includes('engine') }"
              >
                {{ r.format }}
              </span>
            </td>
            <td>
              <span :class="r.status === 'SUCCESS' ? 'badge-online' : 'badge-offline'">
                {{ r.status === 'SUCCESS' ? '[OK] SUCCESS' : r.status }}
              </span>
            </td>
            <td class="text-mono">{{ r.size_mb?.toFixed(1) || '0.0' }} MB</td>
            <td class="text-mono" style="color: var(--cb-green); font-weight: 700;">
              {{ r.inference_time_ms?.toFixed(2) || '0.00' }} ms
            </td>
            <td class="text-mono" style="color: var(--cb-cyan); font-weight: 700;">
              {{ r.fps?.toFixed(0) || '0' }} FPS
            </td>
            <td class="text-mono">
              {{ r.map50_95 > 0 ? `${(r.map50_95 * 100).toFixed(1)}%` : '-' }}
            </td>
            <td class="text-mono" style="font-size: 0.72rem; color: #94a3b8;">
              {{ r.export_args || '-' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
