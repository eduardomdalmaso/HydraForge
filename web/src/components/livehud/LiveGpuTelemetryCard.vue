<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  gpuStats: any
  job: any
}>()

const vramUsed = computed(() => Math.round(props.gpuStats?.used_vram_mb || props.gpuStats?.vram_used_mb || 0))
const vramTotal = computed(() => Math.round(props.gpuStats?.total_vram_mb || props.gpuStats?.vram_total_mb || 32607))
const vramPct = computed(() => Math.round(props.gpuStats?.vram_usage_pct || (vramTotal.value > 0 ? (vramUsed.value / vramTotal.value) * 100 : 0)))
const temp = computed(() => Math.round(props.gpuStats?.temp_celsius || props.gpuStats?.temp_c || 0))
const power = computed(() => Math.round(props.gpuStats?.power_watts || props.gpuStats?.power_w || 0))
const util = computed(() => Math.round(props.gpuStats?.gpu_util_pct || props.gpuStats?.utilization_pct || 0))
const energyKWh = computed(() => props.job?.total_energy_kwh || 0)
const avgFPS = computed(() => Math.round(props.job?.avg_fps || props.job?.fps || 0))
const durationSec = computed(() => Math.round(props.job?.duration_sec || 0))
const clusterName = computed(() => props.gpuStats?.model || props.gpuStats?.name || 'NVIDIA RTX 5090')
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">3. HARDWARE & {{ clusterName.toUpperCase() }}</span>
      <span class="badge-green">NVML LIVE TELEMETRY</span>
    </div>

    <div style="margin-bottom: 0.65rem;">
      <div class="selector-label" style="display: flex; justify-content: space-between;">
        <span>VRAM USAGE // TOTAL VRAM</span>
        <span class="slider-val">
          {{ (vramUsed / 1024).toFixed(1) }} GB / {{ (vramTotal / 1024).toFixed(1) }} GB ({{ vramPct }}%)
        </span>
      </div>
      <div class="epoch-progress-track">
        <div class="epoch-progress-fill" :style="{ width: `${Math.max(vramPct, 2)}%`, background: 'var(--cb-cyan)' }" />
      </div>
    </div>

    <div class="telemetry-row">
      <span class="k">ACTIVE ACCELERATOR</span>
      <span class="v" style="color: var(--cb-cyan);">
        {{ clusterName }} (CUDA 13.3)
      </span>
    </div>
    <div class="telemetry-row">
      <span class="k">GPU TEMP & POWER DRAW</span>
      <span class="v" :style="{ color: temp < 65 ? 'var(--cb-green)' : 'var(--cb-yellow)' }">
        {{ temp }}°C • {{ power }} W
      </span>
    </div>
    <div class="telemetry-row">
      <span class="k">TOTAL ENERGY BUDGET (SESSION)</span>
      <span class="v" style="color: var(--cb-yellow);">
        {{ energyKWh > 0 ? `${energyKWh.toFixed(4)} kWh` : '0.0000 kWh' }}
      </span>
    </div>
    <div class="telemetry-row">
      <span class="k">TRAIN THROUGHPUT & TIME</span>
      <span class="v" style="color: var(--cb-green);">
        {{ avgFPS > 0 ? `${avgFPS.toLocaleString()} FPS` : '-' }} • {{ durationSec > 0 ? `${Math.floor(durationSec / 60)}m ${durationSec % 60}s` : '0s' }}
      </span>
    </div>
    <div class="telemetry-row" style="border-bottom: none;">
      <span class="k">SM & TENSOR CORE UTILIZATION</span>
      <span class="v" :style="{ color: util > 30 ? 'var(--cb-yellow)' : 'var(--cb-green)' }">
        {{ util }}% (Mixed-Precision TC)
      </span>
    </div>
  </div>
</template>
