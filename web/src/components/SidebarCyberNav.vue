<script setup lang="ts">
import HydraLogo from './HydraLogo.vue'
import CyberNavIcons from './CyberNavIcons.vue'
import type { GpuStats } from '../types/telemetry'

const PIPELINE_STEPS = [
  { id: 'datasets', step: '01', title: 'DATASETS', desc: 'Curadoria & Especificação YAML' },
  { id: 'cockpit', step: '02', title: 'CONFIGURAÇÃO', desc: 'Modelos YOLO & Hiperparâmetros' },
  { id: 'live-hud', step: '03', title: 'LIVE TRAINING', desc: 'Monitoramento & Curvas mAP' },
  { id: 'model-zoo', step: '04', title: 'MODEL ZOO', desc: 'Checkpoints & TensorRT Export' },
  { id: 'playground', step: '05', title: 'INFERÊNCIA & FEED', desc: 'Streams & Feed HydraVault' }
]

const AUX_TABS = [
  { id: 'benchmarks', title: 'BENCHMARKS', desc: 'Latência & Throughput FP16/INT8' }
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
  <aside class="sidebar-cyber-nav">
    <div class="sidebar-brand-box" @click="emit('selectTab', 'datasets')">
      <HydraLogo :size="32" />
      <div class="sidebar-brand-text">
        <div class="sidebar-brand-title">
          <span class="brand-hydra">HYDRA</span><span class="brand-forge">FORGE</span>
        </div>
        <span class="sidebar-brand-subtitle">AI TRAINING FOUNDRY // SOTA</span>
      </div>
    </div>

    <nav class="sidebar-nav-section">
      <div class="sidebar-nav-label">PIPELINE DE TREINAMENTO</div>
      
      <button
        v-for="step in PIPELINE_STEPS"
        :key="step.id"
        class="sidebar-step-btn"
        :class="{ active: activeTab === step.id }"
        @click="emit('selectTab', step.id)"
      >
        <span class="sidebar-step-num">{{ step.step }}</span>
        <CyberNavIcons :type="step.id" />
        <div class="sidebar-step-content">
          <span class="sidebar-step-title">{{ step.title }}</span>
          <span class="sidebar-step-desc">{{ step.desc }}</span>
        </div>
      </button>

      <div class="sidebar-nav-label" style="margin-top: 0.75rem;">ANALISE & BENCHMARKS</div>
      <button
        v-for="aux in AUX_TABS"
        :key="aux.id"
        class="sidebar-step-btn"
        :class="{ active: activeTab === aux.id }"
        @click="emit('selectTab', aux.id)"
      >
        <CyberNavIcons :type="aux.id" />
        <div class="sidebar-step-content">
          <span class="sidebar-step-title">{{ aux.title }}</span>
          <span class="sidebar-step-desc">{{ aux.desc }}</span>
        </div>
      </button>
    </nav>

    <footer class="sidebar-footer">
      <div class="sidebar-gpu-badge">
        <div style="display: flex; align-items: center; gap: 0.5rem;">
          <span class="status-dot green"></span>
          <span class="text-mono" style="font-size: 0.72rem; color: #ffffff; font-weight: 700;">
            {{ gpuStats?.name || 'NVIDIA RTX 5090' }}
          </span>
        </div>
        <span class="text-mono" style="font-size: 0.72rem; color: var(--vms-primary, #ff5e3a);">
          {{ (gpuStats?.temp_c || 36).toFixed(0) }}°C
        </span>
      </div>

      <div class="sidebar-eco-links">
        <a href="http://localhost:5173" target="_blank" rel="noreferrer" class="sidebar-eco-link">
          <span>HYDRA VMS</span>
          <span class="sidebar-eco-port">:5173 ↗</span>
        </a>
        <a href="http://localhost:8082" target="_blank" rel="noreferrer" class="sidebar-eco-link">
          <span>HYDRA VAULT</span>
          <span class="sidebar-eco-port">:8082 ↗</span>
        </a>
        <a href="http://localhost:8080" target="_blank" rel="noreferrer" class="sidebar-eco-link">
          <span>HYDRA STREAM</span>
          <span class="sidebar-eco-port">:8080 ↗</span>
        </a>
        <a href="/swagger/" target="_blank" rel="noreferrer" class="sidebar-eco-link">
          <span>SWAGGER DOCS</span>
          <span class="sidebar-eco-port">API ↗</span>
        </a>
      </div>
    </footer>
  </aside>
</template>
