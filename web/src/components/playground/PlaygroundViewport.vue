<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  imageSrc?: string
  detections?: any[]
  selectedEntity?: any
  isScanning?: boolean
  isContinuous?: boolean
  isHydraLinked?: boolean
  activeStream?: any
  isWebcam?: boolean
  videoRef?: any
  telemetry?: any
  gpuStats?: any
  config?: Record<string, any>
}>()

const emit = defineEmits<{
  (e: 'update:selectedEntity', det: any): void
}>()

const imgError = ref(false)
watch(() => props.imageSrc, () => {
  imgError.value = false
})

const isVideoLoop = computed(() => props.config?.source?.startsWith('video:'))
const videoLoopSrc = computed(() => {
  if (isVideoLoop.value) {
    const raw = props.config?.source?.replace(/^video:/, '') || ''
    return `/api/v1/media/stream/${raw}`
  }
  return undefined
})
const videoLoopLabel = computed(() => {
  if (isVideoLoop.value) {
    const raw = props.config?.source?.replace(/^video:/, '') || ''
    return `[LOOP] ${raw.toUpperCase()}`
  }
  return ''
})
</script>

<template>
  <div class="viewport-hud">
    <div class="viewport-top-bar">
      <div style="display: flex; align-items: center; gap: 0.6rem;">
        <span>// KIROSHI OPTICAL HUD // REALTIME TRACKING</span>
        <span
          :style="{
            fontSize: '0.65rem',
            padding: '2px 6px',
            borderRadius: '2px',
            background: isWebcam ? 'rgba(0,240,255,0.15)' : (isVideoLoop ? 'rgba(252,238,10,0.15)' : (isHydraLinked ? 'rgba(0,255,157,0.15)' : 'rgba(252,238,10,0.15)')),
            color: isWebcam ? 'var(--cb-cyan)' : (isVideoLoop ? 'var(--cb-yellow)' : (isHydraLinked ? 'var(--cb-green)' : 'var(--cb-yellow)')),
            border: `1px solid ${isWebcam ? 'var(--cb-cyan)' : (isVideoLoop ? 'var(--cb-yellow)' : (isHydraLinked ? 'var(--cb-green)' : 'var(--cb-yellow)'))}`
          }"
        >
          {{ isWebcam ? '[HARDWARE] WEBCAM /dev/video0' : (isVideoLoop ? videoLoopLabel : (isHydraLinked ? `[STREAM] ${activeStream?.stream_id?.toUpperCase() || 'SHM'}` : '[STANDBY]')) }}
        </span>
      </div>
      <span style="color: var(--cb-yellow); font-family: var(--font-mono);">
        {{ isContinuous ? '● REALTIME SCANNING' : (isScanning ? 'INFERENCE EXECUTING...' : `ACTIVE TARGETS: ${detections?.length || 0}`) }}
      </span>
    </div>

    <div class="viewport-stage">
      <!-- WEBCAM OR VIDEO LOOP PLAYER -->
      <video
        v-if="isWebcam || isVideoLoop"
        id="hud-viewport-video"
        :ref="videoRef as any"
        :src="videoLoopSrc"
        autoplay
        loop
        playsinline
        muted
        class="viewport-img"
      />
      <!-- HYDRASTREAM STATIC SNAPSHOT / MJPEG -->
      <img
        v-else-if="imageSrc && !imgError && imageSrc !== '/hydra-logo.jpg'"
        id="hud-viewport-image"
        :src="imageSrc"
        alt="Tracking Feed"
        class="viewport-img"
        @error="imgError = true"
      />
      <!-- STANDBY PLACEHOLDER -->
      <div v-else class="camera-placeholder-hud">
        <div style="position: relative; margin-bottom: 1.25rem;">
          <svg width="68" height="68" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" style="filter: drop-shadow(0 0 14px rgba(255, 94, 58, 0.55));">
            <path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z" />
            <circle cx="12" cy="13" r="4" />
            <circle cx="12" cy="13" r="1.5" fill="#ff5e3a" />
          </svg>
          <div style="position: absolute; top: -6px; right: -6px; width: 8px; height: 8px; border-radius: 50%; background: #ff5e3a; box-shadow: 0 0 8px #ff5e3a;" />
        </div>
        <div style="color: #ff5e3a; font-family: var(--font-mono); font-size: 0.85rem; font-weight: 600; letter-spacing: 0.08em; margin-bottom: 0.4rem;">
          [HYDRASTREAM SENSOR FEED // 16:9]
        </div>
        <div style="color: #8b94a0; font-family: var(--font-mono); font-size: 0.72rem; letter-spacing: 0.04em;">
          {{ isHydraLinked ? '// FLUXO DE VIDEO ZERO-COPY CALIBRADO' : '// STANDBY // SELECIONE UMA FONTE OU VIDEO EM LOOP' }}
        </div>
      </div>

      <!-- DETECTIONS HUD OVERLAY -->
      <div
        v-for="(det, idx) in (detections || [])"
        :key="det.id || idx"
        class="detection-box"
        :style="{
          left: `${det.box[0]}%`,
          top: `${det.box[1]}%`,
          width: `${det.box[2]}%`,
          height: `${det.box[3]}%`,
          borderColor: selectedEntity?.id === det.id ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
          boxShadow: selectedEntity?.id === det.id ? '0 0 10px var(--cb-yellow)' : '0 0 6px rgba(0,240,255,0.4)'
        }"
        @click.stop="emit('update:selectedEntity', det)"
      >
        <span class="detection-tag">
          {{ det.class_name?.toUpperCase() }} // {{ (det.confidence * 100).toFixed(0) }}%
        </span>
      </div>
    </div>

    <div class="viewport-footer-hud">
      <div style="display: flex; gap: 0.6rem; align-items: center;">
        <span style="color: var(--cb-green); font-weight: 700;">[STATUS: NATIVE 16:9]</span>
        <span>LATENCY: {{ telemetry?.inference_ms ? `${telemetry.inference_ms}ms` : '0.92ms' }}</span>
      </div>
      <div style="color: var(--cb-cyan);">
        ENGINE: {{ (config?.runtime || 'PYTORCH').toUpperCase() }} // {{ (config?.model || 'YOLO26M').toUpperCase() }}
      </div>
    </div>
  </div>
</template>
