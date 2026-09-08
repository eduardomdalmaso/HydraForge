<script setup lang="ts">
import { ref, computed } from 'vue'

const props = withDefaults(defineProps<{ jobs?: any[] }>(), { jobs: () => [] })
const emit = defineEmits<{ (e: 'testInPlayground', run: any): void }>()

const selectedRunId = ref<string | null>(props.jobs[0]?.job_id || null)
const activeRun = computed(() => props.jobs.find(r => r.job_id === selectedRunId.value) || props.jobs[0])

const navigate = (hash: string) => { if (typeof window !== 'undefined') window.location.hash = hash }
const fmtMap = (val: any) => (!val || val <= 0) ? '-' : `${(Number(val) * 100).toFixed(2)}%`
</script>

<template>
  <div class="cyber-card" style="margin-top: 1.25rem;">
    <div class="card-header">
      <span class="card-title">5. COMPARACAO DE EXPERIMENTOS</span>
      <span class="badge-yellow">{{ jobs.length }} RUNS GRAVADAS</span>
    </div>

    <div v-if="jobs.length === 0" style="text-align: center; padding: 1.5rem; color: var(--vms-text-muted); font-size: 0.8rem;">
      Nenhum historico de treinamento gravado no banco de dados.
    </div>

    <template v-else>
      <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 0.65rem; margin-bottom: 0.85rem;">
        <div
          v-for="run in jobs"
          :key="run.job_id"
          class="dataset-card-item"
          :class="{ active: (selectedRunId || jobs[0]?.job_id) === run.job_id }"
          style="margin-bottom: 0; padding: 0.65rem; cursor: pointer;"
          @click="selectedRunId = run.job_id"
        >
          <div>
            <div style="font-family: var(--font-inter); font-size: 0.85rem; font-weight: 600;" :style="{ color: (selectedRunId || jobs[0]?.job_id) === run.job_id ? 'var(--vms-primary)' : '#fff' }">
              {{ run.model_architecture || run.model }} ({{ run.dataset_id }})
            </div>
            <div style="font-family: var(--font-mono); font-size: 0.7rem; color: var(--vms-text-muted); margin-top: 0.2rem;">
              mAP50: <strong style="color: #ffffff;">{{ fmtMap(run.best_map50 || run.metrics?.map50) }}</strong>
              <span v-if="(run.best_map50_95 || run.metrics?.map50_95 || 0) > 0"> • 50-95: <strong style="color: var(--vms-primary);">{{ fmtMap(run.best_map50_95 || run.metrics?.map50_95) }}</strong></span>
              <span> • {{ run.hyperparameters?.epochs || run.total_epochs || 50 }}e</span>
            </div>
          </div>
          <span :class="run.status === 'COMPLETED' ? 'badge-green' : 'badge-cyan'" style="font-size: 0.65rem;">{{ run.status }}</span>
        </div>
      </div>

      <div v-if="activeRun" style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.5rem; margin-bottom: 0.75rem;">
        <div class="matrix-cell highlight">
          <div style="font-size: 0.65rem; color: var(--vms-text-muted);">mAP@50 {{ activeRun.best_map50_95 > 0 ? '/ @50-95' : '' }}</div>
          <div style="font-family: var(--font-mono); font-size: 0.95rem; font-weight: 700; color: #ffffff;">
            {{ activeRun.best_map50 > 0 ? fmtMap(activeRun.best_map50) : 'PENDING' }}
            {{ activeRun.best_map50_95 > 0 ? ` / ${fmtMap(activeRun.best_map50_95)}` : '' }}
          </div>
        </div>
        <div class="matrix-cell">
          <div style="font-size: 0.65rem; color: var(--vms-text-muted);">OTIMIZADOR / LR0</div>
          <div style="font-family: var(--font-mono); font-size: 0.85rem; font-weight: 700; color: #ffffff;">
            {{ activeRun.hyperparameters?.optimizer || 'AdamW' }} ({{ activeRun.hyperparameters?.lr0 || 0.001 }})
          </div>
        </div>
        <div class="matrix-cell">
          <div style="font-size: 0.65rem; color: var(--vms-text-muted);">RESOLUCAO</div>
          <div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: var(--vms-primary);">{{ activeRun.hyperparameters?.imgsz || 640 }} px</div>
        </div>
        <div class="matrix-cell">
          <div style="font-size: 0.65rem; color: var(--vms-text-muted);">BATCH SIZE</div>
          <div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: #fff;">{{ activeRun.hyperparameters?.batch || 32 }}</div>
        </div>
      </div>

      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; margin-top: 0.75rem;">
        <button class="cyber-action-btn" style="width: 100%; padding: 0.75rem;" @click="emit('testInPlayground', activeRun)">PLAYGROUND</button>
        <button class="cyber-action-btn secondary" style="width: 100%; padding: 0.75rem;" @click="navigate('benchmarks')">EXPORTAR</button>
      </div>
    </template>
  </div>
</template>
