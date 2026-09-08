<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  job: any
}>()

const emit = defineEmits<{
  (e: 'abortJob'): void
  (e: 'restartJob', job: any): void
  (e: 'resumeJob', job: any): void
}>()

const isTraining = computed(() => props.job?.status === 'TRAINING' || props.job?.status === 'RUNNING')
const currentEpoch = computed(() => props.job?.current_epoch || 0)
const totalEpochs = computed(() => props.job?.hyperparameters?.epochs || props.job?.total_epochs || props.job?.epochs || 50)
const currentBatch = computed(() => props.job?.current_batch || 0)
const totalBatches = computed(() => props.job?.total_batches || 0)

const epochBatchPct = computed(() => {
  if (totalBatches.value > 0) return Math.min(100, Math.round((currentBatch.value / totalBatches.value) * 100))
  return props.job?.status === 'COMPLETED' ? 100 : 0
})

const globalPct = computed(() => {
  if (totalEpochs.value <= 0) return 0
  if (props.job?.status === 'COMPLETED') return 100
  const base = currentEpoch.value > 0 ? (currentEpoch.value - 1) : 0
  const frac = totalBatches.value > 0 ? (currentBatch.value / totalBatches.value) : 0
  return Math.min(100, ((base + frac) / totalEpochs.value) * 100)
})

const navigate = (hash: string) => {
  if (typeof window !== 'undefined') window.location.hash = hash
}
</script>

<template>
  <div v-if="!job" class="cyber-card">
    <div class="card-header"><span class="card-title">1. ACTIVE TRAINING STATUS</span><span class="badge-cyan">STANDBY / IDLE</span></div>
    <div style="text-align: center; padding: 2.5rem 1rem; color: #94a3b8;">
      <div style="font-family: var(--font-oxanium); font-size: 0.95rem; color: #fff;">NO ACTIVE TRAINING IN PROGRESS</div>
      <div style="font-size: 0.75rem; margin-top: 0.35rem;">Navigate to <strong>TRAINING COCKPIT</strong> to launch a session on PyTorch / RTX 5090.</div>
      <button class="cyber-action-btn" style="margin-top: 1rem;" @click="navigate('cockpit')">COCKPIT</button>
    </div>
  </div>

  <div v-else class="cyber-card">
    <div class="card-header"><span class="card-title">1. ACTIVE TRAINING STATUS</span><span :class="isTraining ? 'badge-green' : 'badge-cyan'">● {{ job.status }}</span></div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
      <div>
        <span style="font-family: var(--font-oxanium); font-size: 1.05rem; font-weight: 700; color: #fff;">{{ job.model_architecture || job.model_name || 'YOLO26' }}</span>
        <span style="margin-left: 0.5rem; font-size: 0.72rem; color: #94a3b8;">// DATASET: {{ job.dataset_id }}</span>
      </div>
      <div v-if="job.hyperparameters?.two_stage" class="stage-pill stage1">STAGE 1 (FROZEN)</div>
    </div>

    <div class="selector-group" style="margin-bottom: 0.65rem;">
      <div class="selector-label" style="display: flex; justify-content: space-between;">
        <span>CURRENT EPOCH PROGRESS</span>
        <span class="slider-val" style="color: var(--cb-cyan);">{{ totalBatches > 0 ? `BATCH ${currentBatch} / ${totalBatches} (${epochBatchPct}%)` : `EPOCH ${currentEpoch > 0 ? currentEpoch : 1}` }}</span>
      </div>
      <div class="epoch-progress-track"><div class="epoch-progress-fill" :style="{ width: `${epochBatchPct}%`, transition: 'width 0.3s ease' }" /></div>
    </div>

    <div class="selector-group" style="margin-bottom: 0.65rem;">
      <div class="selector-label" style="display: flex; justify-content: space-between;">
        <span>TOTAL SESSION COMPLETION</span>
        <span class="slider-val" style="color: var(--cb-yellow);">EPOCH {{ currentEpoch > 0 ? currentEpoch : 1 }} / {{ totalEpochs }} ({{ globalPct.toFixed(1) }}%)</span>
      </div>
      <div class="epoch-progress-track">
        <div class="epoch-progress-fill" :style="{ width: `${globalPct}%`, background: 'linear-gradient(90deg, var(--cb-yellow) 0%, var(--cb-magenta) 100%)', boxShadow: '0 0 12px rgba(252, 238, 10, 0.4)', transition: 'width 0.3s ease' }" />
      </div>
    </div>

    <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 0.65rem; flex-wrap: wrap; gap: 0.4rem;">
      <span style="font-size: 0.72rem; color: #94a3b8; font-family: var(--font-mono);">JOB ID: {{ job.job_id }}</span>
      <div style="display: flex; gap: 0.4rem;">
        <template v-if="isTraining">
          <button class="cyber-action-btn secondary" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('restartJob', job)">RESTART</button>
          <button class="cyber-action-btn danger" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('abortJob')">ABORT</button>
        </template>
        <template v-else>
          <button class="cyber-action-btn" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('restartJob', job)">RESTART</button>
          <button v-if="job.current_epoch > 0 && job.status !== 'COMPLETED'" class="cyber-action-btn secondary" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('resumeJob', job)">RESUME</button>
        </template>
      </div>
    </div>
  </div>
</template>
