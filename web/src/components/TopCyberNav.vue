<script setup lang="ts">
import CyberNavIcons from './CyberNavIcons.vue'
import HydraLogo from './HydraLogo.vue'
import type { GpuStats } from '../types/telemetry'

const QUICK_TABS = [
  { id: 'cockpit', label: 'Cockpit' },
  { id: 'live-hud', label: 'Live HUD' },
  { id: 'benchmarks', label: 'Benchmarks' },
  { id: 'datasets', label: 'Datasets' },
  { id: 'model-zoo', label: 'Model Zoo' },
  { id: 'playground', label: 'Playground' }
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
      <HydraLogo :size="30" />
      <div class="nav-brand-title-box">
        <span class="brand-hydra">HYDRA</span>
        <span class="brand-forge">FORGE</span>
      </div>
    </div>

    <nav class="nav-folder-tabs">
      <button
        v-for="tab in QUICK_TABS"
        :key="tab.id"
        class="tab-btn"
        :class="{ 'tab-btn--active': activeTab === tab.id }"
        @click="emit('selectTab', tab.id)"
      >
        <CyberNavIcons :type="tab.id" />
        <span>{{ tab.label }}</span>
      </button>
    </nav>

    <div class="top-nav-right">
      <div class="header-metric-pill">
        <span class="status-dot green"></span>
        <span class="pill-val">{{ gpuStats?.name || 'RTX 5090' }}</span>
        <span class="pill-sub">{{ (gpuStats?.temp_c || 36).toFixed(0) }}°C</span>
      </div>
      <a href="/swagger/" target="_blank" rel="noreferrer" class="swagger-link-btn">
        API DOCS
      </a>
    </div>
  </header>
</template>
