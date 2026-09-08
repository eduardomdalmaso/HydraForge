<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import ModelZooFilterBar from '../components/modelzoo/ModelZooFilterBar.vue'
import ModelCardGrid from '../components/modelzoo/ModelCardGrid.vue'
import ModelDetailsDrawer from '../components/modelzoo/ModelDetailsDrawer.vue'
import ModelDistillationCard from '../components/modelzoo/ModelDistillationCard.vue'
import { fetchModelsAPI } from '../api/client'

const models = ref<any[]>([])
const selectedFamily = ref('ALL')
const selectedTask = ref('ALL')
const selectedModel = ref<any>(null)
let timer: any = null

const loadModels = async () => {
  const list = await fetchModelsAPI()
  if (list && list.length > 0) {
    models.value = list
    if (!selectedModel.value) {
      selectedModel.value = list[0]
    } else {
      selectedModel.value = list.find((m: any) => m.id === selectedModel.value.id) || list[0]
    }
  }
}

onMounted(() => {
  loadModels()
  timer = setInterval(loadModels, 3000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const filteredModels = computed(() => {
  return models.value.filter(m => {
    const famMatch = selectedFamily.value === 'ALL' || m.family === selectedFamily.value
    const taskMatch = selectedTask.value === 'ALL' || m.task === selectedTask.value
    return famMatch && taskMatch
  })
})

const navigateTo = (tab: string) => {
  window.location.hash = tab
}
</script>

<template>
  <div class="view-container zoo-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">MODEL ZOO REPOSITORY</h1>
      <p class="cockpit-main-subtitle">
        OFFICIAL CHECKPOINTS // CUSTOM TRAINED MODELS ON RTX 5090 // DISTILLATION & RUNTIMES
      </p>
    </div>

    <ModelZooFilterBar
      :selectedFamily="selectedFamily"
      :selectedTask="selectedTask"
      :totalCount="filteredModels.length"
      @update:selectedFamily="(f) => selectedFamily = f"
      @update:selectedTask="(t) => selectedTask = t"
    />

    <div class="zoo-layout">
      <div>
        <ModelCardGrid
          :models="filteredModels"
          :selectedModel="selectedModel"
          @selectModel="(m) => selectedModel = m"
        />
        <ModelDistillationCard @launchDistill="() => navigateTo('cockpit')" />
      </div>

      <div>
        <ModelDetailsDrawer
          :model="selectedModel"
          @sendToCockpit="() => navigateTo('cockpit')"
          @sendToBenchmark="() => navigateTo('benchmarks')"
          @sendToPlayground="() => navigateTo('playground')"
        />
      </div>
    </div>
  </div>
</template>
