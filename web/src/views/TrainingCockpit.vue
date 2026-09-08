<script setup lang="ts">
import { ref, watch } from 'vue'
import ArchitectureCard from '../components/ArchitectureCard.vue'
import HyperparameterCard from '../components/HyperparameterCard.vue'
import TwoStageControlCard from '../components/TwoStageControlCard.vue'
import HardwareEstimatorCard from '../components/HardwareEstimatorCard.vue'
import { launchTrainingJobAPI } from '../api/client'
import type { DatasetInfo } from '../types/dataset'

const props = defineProps<{ datasets?: DatasetInfo[] }>()
const emit = defineEmits<{ (e: 'jobLaunched', result: any): void }>()

const family = ref('yolo26')
const scale = ref('s')
const task = ref('detect')
const selectedDataset = ref(props.datasets?.[0]?.id || '')
const isLaunching = ref(false)
const errorMsg = ref('')
const params = ref({
  epochs: 50, batch: 16, imgsz: 640, optimizer: 'AdamW', lr0: 0.001, amp: true,
  workers: 8, patience: 20, close_mosaic: 10, two_stage: false,
  stage1_epochs: 20, stage1_freeze: 10, stage2_epochs: 30
})

watch(() => props.datasets, (newDs) => {
  if (newDs && newDs.length > 0 && !selectedDataset.value) selectedDataset.value = newDs[0].id
}, { immediate: true })

const handleLaunch = async () => {
  isLaunching.value = true
  errorMsg.value = ''
  try {
    const modelArch = `${family.value}${scale.value}`
    const result = await launchTrainingJobAPI({
      job_id: `run_${modelArch}_${Date.now()}`,
      model_architecture: modelArch,
      task: task.value,
      dataset_id: selectedDataset.value,
      hyperparameters: params.value
    })
    emit('jobLaunched', result)
    if (typeof window !== 'undefined') window.location.hash = 'live-hud'
  } catch (err: any) {
    errorMsg.value = err.message || 'Erro ao iniciar treinamento'
  } finally {
    isLaunching.value = false
  }
}
</script>

<template>
  <div class="view-container cockpit-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">ESTUDIO DE TREINAMENTO YOLO</h1>
      <p class="cockpit-main-subtitle">CONFIGURE ARQUITETURA, HIPERPARAMETROS E EXECUTE TREINOS NA GPU NVIDIA RTX 5090</p>
    </div>

    <div v-if="errorMsg" class="cyber-alerts-banner" style="margin-bottom: 1rem;">
      <span class="cyber-alert-badge">ERRO</span>
      <span>{{ errorMsg }}</span>
    </div>

    <div class="cockpit-quad-grid">
      <ArchitectureCard v-model:family="family" v-model:scale="scale" v-model:task="task" />
      <HyperparameterCard :params="params as any" @update:params="(p) => params = p as any" />
      <TwoStageControlCard :params="params" @update:params="(p) => params = p as any" />
      <HardwareEstimatorCard :datasets="datasets" v-model:selectedDataset="selectedDataset" :scale="scale" :batchSize="params.batch" :imgsz="params.imgsz" />
    </div>

    <div class="cockpit-bottom-bar">
      <button class="cockpit-launch-action-btn" :disabled="isLaunching" @click="handleLaunch">
        {{ isLaunching ? 'INICIANDO TREINO...' : 'INICIAR TREINAMENTO' }}
      </button>
    </div>
  </div>
</template>
