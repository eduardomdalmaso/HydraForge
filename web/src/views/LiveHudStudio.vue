<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import LiveTrainingStatusCard from '../components/livehud/LiveTrainingStatusCard.vue'
import LiveMetricsChartsCard from '../components/livehud/LiveMetricsChartsCard.vue'
import LiveGpuTelemetryCard from '../components/livehud/LiveGpuTelemetryCard.vue'
import LiveTerminalLogsCard from '../components/livehud/LiveTerminalLogsCard.vue'
import ExperimentsComparisonCard from '../components/livehud/ExperimentsComparisonCard.vue'
import { fetchJobsAPI, fetchTelemetryAPI, stopTrainingJobAPI, launchTrainingJobAPI } from '../api/client'

const jobs = ref<any[]>([])
const activeJob = ref<any>(null)
const telemetry = ref<any>(null)
let timer: any = null

const loadLiveData = async () => {
  const [list, tel] = await Promise.all([fetchJobsAPI(), fetchTelemetryAPI()])
  if (list) jobs.value = list
  if (tel) {
    telemetry.value = tel
    if (tel.active_job) {
      activeJob.value = tel.active_job
    } else if (list) {
      const running = list.find((j: any) => j.status === 'TRAINING' || j.status === 'RUNNING' || j.status === 'QUEUED')
      activeJob.value = running || list[0] || null
    }
  } else if (list) {
    const running = list.find((j: any) => j.status === 'TRAINING' || j.status === 'RUNNING' || j.status === 'QUEUED')
    activeJob.value = running || list[0] || null
  }
}

onMounted(() => {
  loadLiveData()
  timer = setInterval(loadLiveData, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const handleAbort = async () => {
  if (activeJob.value?.job_id) {
    await stopTrainingJobAPI(activeJob.value.job_id)
    loadLiveData()
  }
}

const handleRestart = async (job: any) => {
  if (!job) return
  if (job.status === 'TRAINING' || job.status === 'RUNNING') {
    await stopTrainingJobAPI(job.job_id)
  }
  await launchTrainingJobAPI({
    model_architecture: job.model_architecture || 'yolo26m',
    task: job.task || 'detect',
    dataset_id: job.dataset_id,
    hyperparameters: { ...job.hyperparameters, pretrained: true }
  })
  loadLiveData()
}

const navigate = (hash: string) => {
  if (typeof window !== 'undefined') window.location.hash = hash
}
</script>

<template>
  <div class="view-container livehud-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">LIVE HUD TELEMETRY & VALIDATION STUDIO</h1>
      <p class="cockpit-main-subtitle">
        REAL-TIME PYTORCH RUNTIMES // LOSS & mAP DYNAMICS // GPU HARDWARE TELEMETRY & EXPERIMENT MATRIX
      </p>
    </div>

    <div class="livehud-grid">
      <div style="display: flex; flexDirection: column; gap: 1.25rem;">
        <LiveTrainingStatusCard
          :job="activeJob"
          @abortJob="handleAbort"
          @restartJob="handleRestart"
          @resumeJob="handleRestart"
        />
        <LiveGpuTelemetryCard :gpuStats="telemetry?.gpu_stats" :job="activeJob" />
      </div>

      <div style="display: flex; flexDirection: column; gap: 1.25rem;">
        <LiveMetricsChartsCard :job="activeJob" :recentMetrics="telemetry?.recent_metrics" />
        <LiveTerminalLogsCard :job="activeJob" :rawLogs="telemetry?.raw_logs" :recentMetrics="telemetry?.recent_metrics" :gpuStats="telemetry?.gpu_stats" />
      </div>
    </div>

    <ExperimentsComparisonCard :jobs="jobs" @testInPlayground="navigate('playground')" />
  </div>
</template>
