<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import TopCyberNav from './components/TopCyberNav.vue'
import TrainingCockpit from './views/TrainingCockpit.vue'
import LiveHudStudio from './views/LiveHudStudio.vue'
import BenchmarkStudio from './views/BenchmarkStudio.vue'
import PlaygroundStudio from './views/PlaygroundStudio.vue'
import DatasetStudio from './views/DatasetStudio.vue'
import ModelZooStudio from './views/ModelZooStudio.vue'
import { fetchTelemetryAPI, fetchDatasetsAPI } from './api/client'
import type { TelemetryData } from './types/telemetry'
import type { DatasetInfo } from './types/dataset'

const getInitialTab = () => {
  const hash = window.location.hash.replace('#', '')
  return hash || 'cockpit'
}

const activeTab = ref(getInitialTab())
const telemetry = ref<TelemetryData | null>(null)
const datasets = ref<DatasetInfo[]>([])
let pollTimer: any = null

const handleHashChange = () => {
  const hash = window.location.hash.replace('#', '')
  activeTab.value = hash || 'cockpit'
}

const navigateTo = (tabId: string) => {
  window.location.hash = tabId
  activeTab.value = tabId
}

const loadInitialData = async () => {
  const tel = await fetchTelemetryAPI()
  if (tel) telemetry.value = tel
  const ds = await fetchDatasetsAPI()
  if (ds) datasets.value = ds
}

onMounted(() => {
  window.addEventListener('hashchange', handleHashChange)
  loadInitialData()
  pollTimer = setInterval(loadInitialData, 3000)
})

onUnmounted(() => {
  window.removeEventListener('hashchange', handleHashChange)
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<template>
  <div class="app-layout" style="flex-direction: column;">
    <TopCyberNav
      :activeTab="activeTab"
      :gpuStats="telemetry?.gpu_stats"
      @selectTab="navigateTo"
    />
    <main class="content-body">
      <TrainingCockpit
        v-if="activeTab === 'cockpit'"
        :datasets="datasets"
        @jobLaunched="() => navigateTo('live-hud')"
      />
      <LiveHudStudio v-else-if="activeTab === 'live-hud'" />
      <BenchmarkStudio v-else-if="activeTab === 'benchmarks'" />
      <DatasetStudio v-else-if="activeTab === 'datasets'" :datasets="datasets" />
      <PlaygroundStudio v-else-if="activeTab === 'playground'" />
      <ModelZooStudio v-else-if="activeTab === 'model-zoo'" />
    </main>
  </div>
</template>
