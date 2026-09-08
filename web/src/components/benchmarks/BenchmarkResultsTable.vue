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
      <span class="card-title">MATRIZ DE RUNTIMES & EXPORTACAO</span>
      <span class="badge-cyan">{{ hasResults ? `${results.length} RUNTIMES` : '0 RUNTIMES' }}</span>
    </div>
    <div style="padding: 0; overflow-x: auto; margin-top: 0.5rem;">
      <div v-if="!hasResults" style="text-align: center; padding: 2rem; color: var(--vms-text-muted); font-size: 0.8rem;">
        Nenhum resultado de benchmark gravado.
      </div>
      <table v-else class="cyber-table">
        <thead>
          <tr>
            <th>FORMATO</th>
            <th>STATUS</th>
            <th>TAMANHO</th>
            <th>LATENCIA (GPU)</th>
            <th>THROUGHPUT</th>
            <th>mAP 50-95</th>
            <th>ARGUMENTOS</th>
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
              <span :class="r.status === 'SUCCESS' ? 'badge-green' : 'badge-cyan'">
                {{ r.status === 'SUCCESS' ? '[OK] SUCESSO' : r.status }}
              </span>
            </td>
            <td class="text-mono">{{ r.size_mb?.toFixed(1) || '0.0' }} MB</td>
            <td class="text-mono" style="color: var(--vms-primary); font-weight: 700;">
              {{ r.inference_time_ms?.toFixed(2) || '0.00' }} ms
            </td>
            <td class="text-mono" style="color: #ffffff; font-weight: 700;">
              {{ r.fps?.toFixed(0) || '0' }} FPS
            </td>
            <td class="text-mono">
              {{ r.map50_95 > 0 ? `${(r.map50_95 * 100).toFixed(2)}%` : '-' }}
            </td>
            <td class="text-mono" style="font-size: 0.72rem; color: var(--vms-text-muted);">
              {{ r.export_args || '-' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
