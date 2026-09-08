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
    <div class="card-header"><span class="card-title">1. STATUS DO TREINAMENTO</span><span class="badge-cyan">AGUARDANDO</span></div>
    <div style="text-align: center; padding: 2.5rem 1rem; color: var(--vms-text-muted);">
      <div style="font-family: var(--font-inter); font-size: 0.95rem; font-weight: 600; color: #fff;">NENHUM TREINAMENTO EM ANDAMENTO</div>
      <div style="font-size: 0.75rem; margin-top: 0.35rem;">Acesse o <strong>ESTUDIO DE TREINO</strong> para iniciar uma sessao na GPU RTX 5090.</div>
      <button class="cyber-action-btn" style="margin-top: 1rem;" @click="navigate('cockpit')">IR PARA TREINO</button>
    </div>
  </div>

  <div v-else class="cyber-card">
    <div class="card-header"><span class="card-title">1. STATUS DO TREINAMENTO</span><span :class="isTraining ? 'badge-green' : 'badge-cyan'">● {{ job.status }}</span></div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
      <div>
        <span style="font-family: var(--font-inter); font-size: 1rem; font-weight: 700; color: #fff;">{{ job.model_architecture || job.model_name || 'YOLO26' }}</span>
        <span style="margin-left: 0.5rem; font-size: 0.72rem; color: var(--vms-text-muted);">// DATASET: {{ job.dataset_id }}</span>
      </div>
      <div v-if="job.hyperparameters?.two_stage" class="stage-pill stage1">ESTAGIO 1 (FROZEN)</div>
    </div>

    <div class="selector-group" style="margin-bottom: 0.65rem;">
      <div class="selector-label" style="display: flex; justify-content: space-between;">
        <span>PROGRESSO DA EPOCA ATUAL</span>
        <span class="slider-val" style="color: var(--vms-info);">{{ totalBatches > 0 ? `BATCH ${currentBatch} / ${totalBatches} (${epochBatchPct}%)` : `EPOCA ${currentEpoch > 0 ? currentEpoch : 1}` }}</span>
      </div>
      <div class="epoch-progress-track"><div class="epoch-progress-fill" :style="{ width: `${epochBatchPct}%`, transition: 'width 0.3s ease' }" /></div>
    </div>

    <div class="selector-group" style="margin-bottom: 0.65rem;">
      <div class="selector-label" style="display: flex; justify-content: space-between;">
        <span>CONCLUSAO TOTAL DA SESSAO</span>
        <span class="slider-val" style="color: var(--vms-primary);">EPOCA {{ currentEpoch > 0 ? currentEpoch : 1 }} / {{ totalEpochs }} ({{ globalPct.toFixed(1) }}%)</span>
      </div>
      <div class="epoch-progress-track">
        <div class="epoch-progress-fill" :style="{ width: `${globalPct}%`, background: 'var(--vms-primary)', transition: 'width 0.3s ease' }" />
      </div>
    </div>

    <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 0.65rem; flex-wrap: wrap; gap: 0.4rem;">
      <span style="font-size: 0.72rem; color: var(--vms-text-muted); font-family: var(--font-mono);">JOB ID: {{ job.job_id }}</span>
      <div style="display: flex; gap: 0.4rem;">
        <template v-if="isTraining">
          <button class="cyber-action-btn secondary" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('restartJob', job)">REINICIAR</button>
          <button class="cyber-action-btn danger" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('abortJob')">INTERROMPER</button>
        </template>
        <template v-else>
          <button class="cyber-action-btn" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('restartJob', job)">REINICIAR</button>
          <button v-if="job.current_epoch > 0 && job.status !== 'COMPLETED'" class="cyber-action-btn secondary" style="padding: 0.35rem 0.75rem; font-size: 0.7rem;" @click="emit('resumeJob', job)">RETOMAR</button>
        </template>
      </div>
    </div>
  </div>
</template>
