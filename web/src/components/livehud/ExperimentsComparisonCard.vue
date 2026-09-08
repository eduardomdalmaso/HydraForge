<script setup lang="ts">
import { ref, computed } from 'vue'

const props = withDefaults(defineProps<{
  jobs?: any[]
}>(), {
  jobs: () => []
})

const emit = defineEmits<{
  (e: 'testInPlayground', run: any): void
}>()

const selectedRunId = ref<string | null>(props.jobs[0]?.job_id || null)
const activeRun = computed(() => props.jobs.find(r => r.job_id === selectedRunId.value) || props.jobs[0])

const navigate = (hash: string) => {
  if (typeof window !== 'undefined') window.location.hash = hash
}
</script>

<template>
  <div class="cyber-card" style="margin-top: 1.25rem;">
    <div class="card-header">
      <span class="card-title">5. EXPERIMENT COMPARISON & METRICS VALIDATION</span>
      <span class="badge-yellow">{{ jobs.length }} RUNS RECORDED</span>
    </div>

    <div v-if="jobs.length === 0" style="text-align: center; padding: 1.5rem; color: #94a3b8; font-size: 0.8rem;">
      No training history recorded in database. Runs launched from the Cockpit will appear here for cross-validation.
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
            <div
              style="font-family: var(--font-oxanium); font-size: 0.85rem; font-weight: 700;"
              :style="{ color: (selectedRunId || jobs[0]?.job_id) === run.job_id ? 'var(--cb-yellow)' : '#fff' }"
            >
              {{ run.model_architecture || run.model }} ({{ run.dataset_id }})
            </div>
            <div style="font-family: var(--font-mono); font-size: 0.7rem; color: #94a3b8; margin-top: 0.2rem;">
              mAP50: <strong style="color: var(--cb-green);">{{ (run.best_map50 || run.metrics?.map50 || 0) > 0 ? `${((run.best_map50 || run.metrics?.map50) * 100).toFixed(1)}%` : '-' }}</strong>
              <span v-if="(run.best_map50_95 || run.metrics?.map50_95 || 0) > 0"> • 50-95: <strong style="color: var(--cb-cyan);">{{ (((run.best_map50_95 || run.metrics?.map50_95 || 0) * 100)).toFixed(1) }}%</strong></span>
              <span> • {{ run.hyperparameters?.epochs || run.total_epochs || 50 }}e</span>
            </div>
          </div>
          <span :class="run.status === 'COMPLETED' ? 'badge-online' : 'badge-cyan'" style="font-size: 0.65rem;">
            {{ run.status }}
          </span>
        </div>
      </div>

      <div v-if="activeRun" style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.5rem; margin-bottom: 0.75rem;">
        <div class="matrix-cell highlight">
          <div style="font-size: 0.65rem; color: #94a3b8;">mAP@50 {{ activeRun.best_map50_95 > 0 ? '/ @50-95' : '' }}</div>
          <div style="font-family: var(--font-mono); font-size: 0.95rem; font-weight: 700; color: var(--cb-green);">
            {{ activeRun.best_map50 > 0 ? `${(activeRun.best_map50 * 100).toFixed(1)}%` : 'PENDING' }}
            {{ activeRun.best_map50_95 > 0 ? ` / ${(activeRun.best_map50_95 * 100).toFixed(1)}%` : '' }}
          </div>
        </div>
        <div class="matrix-cell">
          <div style="font-size: 0.65rem; color: #94a3b8;">OPTIMIZER / LR0</div>
          <div style="font-family: var(--font-mono); font-size: 0.85rem; font-weight: 700; color: var(--cb-cyan);">
            {{ activeRun.hyperparameters?.optimizer || 'AdamW' }} ({{ activeRun.hyperparameters?.lr0 || 0.001 }})
          </div>
        </div>
        <div class="matrix-cell">
          <div style="font-size: 0.65rem; color: #94a3b8;">RESOLUTION</div>
          <div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: var(--cb-yellow);">
            {{ activeRun.hyperparameters?.imgsz || 640 }} px
          </div>
        </div>
        <div class="matrix-cell">
          <div style="font-size: 0.65rem; color: #94a3b8;">BATCH SIZE</div>
          <div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: #fff;">
            {{ activeRun.hyperparameters?.batch || activeRun.hyperparameters?.batch_size || 32 }}
          </div>
        </div>
      </div>

      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; margin-top: 0.75rem;">
        <button class="cyber-action-btn" style="width: 100%; padding: 0.75rem;" @click="emit('testInPlayground', activeRun)">PLAYGROUND</button>
        <button class="cyber-action-btn secondary" style="width: 100%; padding: 0.75rem;" @click="navigate('benchmarks')">EXPORT</button>
      </div>
    </template>
  </div>
</template>
