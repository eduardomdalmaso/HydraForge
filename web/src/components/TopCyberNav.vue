<script setup lang="ts">
import CyberNavIcons from './CyberNavIcons.vue'
import HydraLogo from './HydraLogo.vue'
import type { GpuStats } from '../types/telemetry'

const QUICK_TABS = [
  { id: 'cockpit', label: 'COCKPIT', tag: 'R-01' },
  { id: 'live-hud', label: 'LIVE HUD', tag: 'R-02' },
  { id: 'benchmarks', label: 'BENCHMARKS', tag: 'R-03' },
  { id: 'datasets', label: 'DATASETS', tag: 'R-04' },
  { id: 'model-zoo', label: 'MODEL ZOO', tag: 'R-05' },
  { id: 'playground', label: 'PLAYGROUND', tag: 'R-06' }
]

defineProps<{
  activeTab: string
  gpuStats?: GpuStats
}>()

const emit = defineEmits<{
  (e: 'selectTab', tabId: string): void
}>()
</script>

<template>
  <header class="top-cyber-nav">
    <div class="nav-brand-btn" @click="emit('selectTab', 'cockpit')">
      <HydraLogo :size="34" />
      <div class="nav-brand-title-box">
        <span class="brand-hydra">HYDRA</span>
        <span class="brand-forge">FORGE</span>
      </div>
    </div>

    <div class="nav-folder-tabs">
      <button
        v-for="tab in QUICK_TABS"
        :key="tab.id"
        class="cybr-btn"
        :class="{ 'cybr-btn--active': activeTab === tab.id }"
        @click="emit('selectTab', tab.id)"
      >
        <span class="cybr-btn__content">
          <CyberNavIcons :type="tab.id" />
          <span>{{ tab.label }}</span>
          <span aria-hidden="true">_</span>
        </span>
        <span aria-hidden="true" class="cybr-btn__glitch">
          <CyberNavIcons :type="tab.id" />
          <span>{{ tab.label }}_</span>
        </span>
        <span aria-hidden="true" class="cybr-btn__tag">
          {{ tab.tag }}
        </span>
      </button>
    </div>

    <div class="top-nav-right">
      <div class="header-metric-pill">
        <span class="status-dot green"></span>
        <span class="pill-val cyan">{{ gpuStats?.name || 'RTX 5090' }}</span>
        <span class="pill-sub">{{ (gpuStats?.temp_c || 36).toFixed(0) }}°C</span>
      </div>
      <a href="/swagger/" target="_blank" rel="noreferrer" class="swagger-link-btn">
        [ API DOCS ]
      </a>
    </div>
  </header>
</template>
