<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import BenchmarkLauncherCard from '../components/benchmarks/BenchmarkLauncherCard.vue'
import BenchmarkSummaryCards from '../components/benchmarks/BenchmarkSummaryCards.vue'
import BenchmarkResultsTable from '../components/benchmarks/BenchmarkResultsTable.vue'
import BenchmarkThroughputChart from '../components/benchmarks/BenchmarkThroughputChart.vue'
import { fetchBenchmarksAPI, launchBenchmarkAPI, fetchBenchmarkByIDAPI } from '../api/benchmark_client'
import { fetchJobsAPI, fetchDatasetsAPI, fetchTelemetryAPI } from '../api/client'

const activeJob = ref<any>(null)
const isRunning = ref(false)
const errorMsg = ref('')
const jobs = ref<any[]>([])
const datasets = ref<any[]>([])
const gpuStats = ref<any>(null)
let pollTimer: any = null

const loadLatestBenchmark = async () => {
  const list = await fetchBenchmarksAPI()
  if (list && list.length > 0) {
    const latest = list[list.length - 1]
    activeJob.value = latest
    isRunning.value = latest.status === 'RUNNING' || latest.status === 'QUEUED'
  }
}

onMounted(() => {
  fetchJobsAPI().then(j => { if (j) jobs.value = j })
  fetchDatasetsAPI().then(d => { if (d) datasets.value = d })
  fetchTelemetryAPI().then(t => { if (t?.gpu_stats) gpuStats.value = t.gpu_stats })
  loadLatestBenchmark()

  pollTimer = setInterval(async () => {
    if (activeJob.value?.job_id && isRunning.value) {
      const updated = await fetchBenchmarkByIDAPI(activeJob.value.job_id)
      if (updated) {
        activeJob.value = updated
        if (updated.status !== 'RUNNING' && updated.status !== 'QUEUED') {
          isRunning.value = false
        }
      }
    }
  }, 1000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

const handleLaunch = async (config: any) => {
  errorMsg.value = ''
  try {
    isRunning.value = true
    const created = await launchBenchmarkAPI(config)
    activeJob.value = created
  } catch (err: any) {
    errorMsg.value = err.message || 'Failed to start benchmark'
    isRunning.value = false
  }
}
</script>

<template>
  <div class="view-container benchmark-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">BENCHMARK & EXPORT STUDIO</h1>
      <p class="cockpit-main-subtitle">
        HARDWARE ACCELERATION // TENSORRT 10.x // ONNX // OPENVINO & THROUGHPUT PROFILING
      </p>
    </div>

    <div v-if="errorMsg" class="cyber-alert alert-critical">
      <span>[ALERT] {{ errorMsg }}</span>
    </div>

    <BenchmarkSummaryCards :results="activeJob?.results" />

    <div class="benchmark-grid-top">
      <BenchmarkLauncherCard
        :isRunning="isRunning"
        :jobs="jobs"
        :datasets="datasets"
        :gpuStats="gpuStats"
        @launch="handleLaunch"
      />
      <BenchmarkThroughputChart :results="activeJob?.results" />
    </div>

    <BenchmarkResultsTable :results="activeJob?.results" />
  </div>
</template>
