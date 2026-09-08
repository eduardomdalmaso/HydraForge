<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import { formatLoss } from '../../utils/formatters'
import type { TrainingMetricPayload } from '../../types/training'

const props = defineProps<{
  job: any
  rawLogs?: string[]
  recentMetrics?: TrainingMetricPayload[]
  gpuStats?: any
}>()

const terminalRef = ref<HTMLElement | null>(null)
const gpuName = computed(() => props.gpuStats?.model || 'NVIDIA RTX 5090')
const gpuVramGB = computed(() => Math.round((props.gpuStats?.total_vram_mb || 32607) / 1024))

const logs = computed(() => {
  if (props.rawLogs && props.rawLogs.length > 0) {
    return props.rawLogs.map(text => {
      let type = 'normal'
      if (text.startsWith('[HydraForge]') || text.includes('TRAINING')) type = 'active'
      else if (text.startsWith('[Epoch') || text.includes('Loss:')) type = 'metric'
      else if (text.includes('SUCCESS') || text.includes('saved')) type = 'success'
      return { text, type }
    })
  }

  if (!props.job) {
    return [
      { text: '[HydraForge Control Plane] Standby listener initialized on :8081', type: 'active' },
      { text: `[Hardware Telemetry] ${gpuName.value} (${gpuVramGB.value}GB VRAM) detected & ready for CUDA workloads`, type: 'metric' },
      { text: '[Status] Awaiting job launch from Training Cockpit...', type: 'normal' }
    ]
  }

  const modelName = props.job.model_architecture || props.job.model_name || 'yolo26m'
  const epochs = props.job.hyperparameters?.epochs || props.job.total_epochs || 50
  const lines = [
    { text: `[HydraForge] Training Job ${props.job.job_id} (${modelName}) started on ${props.job.dataset_id}`, type: 'active' },
    { text: `[Configuration] Epochs: ${epochs} // Batch: ${props.job.hyperparameters?.batch || 32} // ImgSz: ${props.job.hyperparameters?.imgsz || 640}`, type: 'normal' }
  ]

  if (props.recentMetrics && props.recentMetrics.length > 0) {
    props.recentMetrics.slice(-6).forEach(m => {
      lines.push({
        text: `[Epoch ${m.epoch}/${epochs}] Box Loss: ${formatLoss(m.box_loss)} // Cls: ${formatLoss(m.cls_loss)} // mAP50: ${((m.map50 || 0) * 100).toFixed(1)}% // ${Math.round(m.fps || 0)} FPS`,
        type: 'metric'
      })
    })
  }
  return lines
})

watch(() => logs.value.length, async () => {
  await nextTick()
  if (terminalRef.value) {
    terminalRef.value.scrollTop = terminalRef.value.scrollHeight
  }
})
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">4. STREAMING PYTORCH LOGS</span>
      <span class="badge-cyan">{{ job ? job.status : 'IDLE' }}</span>
    </div>

    <div ref="terminalRef" class="terminal-box" style="margin-top: 0.5rem; max-height: 180px; overflow-y: auto;">
      <div
        v-for="(log, i) in logs"
        :key="i"
        :class="{
          'terminal-row-active': log.type === 'active',
          'terminal-row-metric': log.type === 'metric',
          'terminal-row-success': log.type === 'success'
        }"
        style="margin-bottom: 0.2rem; font-family: var(--font-mono); font-size: 0.75rem;"
      >
        {{ log.text }}
      </div>
    </div>
  </div>
</template>
