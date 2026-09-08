<script setup lang="ts">
import { computed } from 'vue'
import { formatLoss } from '../../utils/formatters'
import type { TrainingMetricPayload } from '../../types/training'

const props = defineProps<{ job: any; recentMetrics?: TrainingMetricPayload[] }>()
const latest = computed(() => props.recentMetrics?.length ? props.recentMetrics[props.recentMetrics.length - 1] : null)
const map50 = computed(() => latest.value?.map50 || props.job?.best_map50 || 0)
const map5095 = computed(() => latest.value?.map50_95 || 0)
const hasData = computed(() => (props.recentMetrics && props.recentMetrics.length > 0) || (props.job && props.job.current_epoch > 0))

const calcPoints = (key: 'box_loss' | 'cls_loss' | 'map', strokeScale: number, isMap?: boolean) => {
  if (!props.recentMetrics || props.recentMetrics.length < 2) return null
  const total = props.recentMetrics[0].total_epochs || 50
  const maxVal = isMap ? 1.0 : props.recentMetrics.reduce((m, item) => Math.max(m, (item as any)[key] || 0), 0.005)
  return props.recentMetrics.map((m) => {
    const x = Math.round((m.epoch / total) * 380 + 10)
    const val = isMap ? Math.min(1.0, m.map50_95 || m.map50 || 0) : Math.min(1.0, ((m as any)[key] || 0) / maxVal)
    const y = Math.round(95 - val * strokeScale)
    return `${x},${y}`
  }).join(' ')
}

const boxPoints = computed(() => calcPoints('box_loss', 75))
const clsPoints = computed(() => calcPoints('cls_loss', 75))
const mapPoints = computed(() => calcPoints('map', 80, true))
</script>

<template>
  <div class="cyber-card">
    <div class="card-header"><span class="card-title">2. ACCURACY & LOSS DYNAMICS</span><span class="badge-cyan">{{ hasData ? 'ACTIVE METRICS' : 'STANDBY' }}</span></div>

    <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.5rem; margin-bottom: 0.75rem;">
      <div class="matrix-cell highlight"><div style="font-size: 0.65rem; color: #94a3b8;">mAP@50</div><div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: var(--cb-green);">{{ map50 > 0 ? `${(map50 * 100).toFixed(1)}%` : '-' }}</div></div>
      <div class="matrix-cell highlight"><div style="font-size: 0.65rem; color: #94a3b8;">mAP@50-95</div><div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: var(--cb-cyan);">{{ map5095 > 0 ? `${(map5095 * 100).toFixed(1)}%` : '-' }}</div></div>
      <div class="matrix-cell"><div style="font-size: 0.65rem; color: #94a3b8;">BOX LOSS</div><div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: var(--cb-yellow);">{{ latest ? formatLoss(latest.box_loss) : '-' }}</div></div>
      <div class="matrix-cell"><div style="font-size: 0.65rem; color: #94a3b8;">CLS LOSS</div><div style="font-family: var(--font-mono); font-size: 1rem; font-weight: 700; color: var(--cb-magenta);">{{ latest ? formatLoss(latest.cls_loss) : '-' }}</div></div>
    </div>

    <div style="position: relative; background: rgba(0,0,0,0.4); border-radius: 4px; border: 1px solid rgba(0,240,255,0.15); padding: 0.5rem;">
      <div style="display: flex; justify-content: space-between; font-size: 0.65rem; font-family: var(--font-mono); color: #94a3b8; margin-bottom: 0.25rem;">
        <span><span style="color: var(--cb-green);">●</span> mAP 50-95</span><span><span style="color: var(--cb-yellow);">●</span> Box Loss</span><span><span style="color: var(--cb-magenta);">●</span> Cls Loss</span>
      </div>
      <svg class="hud-svg-chart" viewBox="0 0 400 110">
        <line x1="0" y1="30" x2="400" y2="30" class="hud-grid-line" /><line x1="0" y1="60" x2="400" y2="60" class="hud-grid-line" /><line x1="0" y1="90" x2="400" y2="90" class="hud-grid-line" />
        <template v-if="hasData && (boxPoints || clsPoints || mapPoints)">
          <polyline v-if="boxPoints" fill="none" stroke="var(--cb-yellow)" stroke-width="1.8" :points="boxPoints" />
          <polyline v-if="clsPoints" fill="none" stroke="var(--cb-magenta)" stroke-width="1.8" :points="clsPoints" />
          <polyline v-if="mapPoints" fill="none" stroke="var(--cb-green)" stroke-width="2" :points="mapPoints" />
        </template>
        <template v-else-if="hasData">
          <text x="200" y="60" text-anchor="middle" fill="var(--cb-cyan)" font-family="var(--font-mono)" font-size="11">Processing PyTorch Epochs... accumulating loss & mAP trajectory</text>
        </template>
        <template v-else>
          <text x="200" y="60" text-anchor="middle" fill="#64748b" font-family="var(--font-mono)" font-size="11">{{ job ? 'PyTorch initialized on GPU... awaiting Epoch 1 metrics' : 'Awaiting training session launch to render dynamics...' }}</text>
        </template>
      </svg>
    </div>
  </div>
</template>
