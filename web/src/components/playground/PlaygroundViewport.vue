<script setup lang="ts">
defineProps<{
  imageSrc?: string
  detections?: any[]
  selectedEntity?: any
  isScanning?: boolean
  isContinuous?: boolean
  isHydraLinked?: boolean
  activeStream?: any
  isWebcam?: boolean
  videoRef?: any
}>()

const emit = defineEmits<{
  (e: 'update:selectedEntity', det: any): void
}>()
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
            background: isWebcam ? 'rgba(0,240,255,0.15)' : (isHydraLinked ? 'rgba(0,255,157,0.15)' : 'rgba(252,238,10,0.15)'),
            color: isWebcam ? 'var(--cb-cyan)' : (isHydraLinked ? 'var(--cb-green)' : 'var(--cb-yellow)'),
            border: `1px solid ${isWebcam ? 'var(--cb-cyan)' : (isHydraLinked ? 'var(--cb-green)' : 'var(--cb-yellow)')}`
          }"
        >
          {{ isWebcam ? '[HARDWARE] WEBCAM /dev/video0' : (isHydraLinked ? `[STREAM] ${activeStream?.stream_id?.toUpperCase() || 'SHM'}` : '[STANDBY]') }}
        </span>
      </div>
      <span style="color: var(--cb-yellow); font-family: var(--font-mono);">
        {{ isContinuous ? '● REALTIME STREAM TRACKING' : (isScanning ? 'INFERENCE EXECUTING...' : `ACTIVE TARGETS: ${detections?.length || 0}`) }}
      </span>
    </div>

    <div class="viewport-canvas-area" style="background: #000; overflow: hidden;">
      <div style="position: relative; display: inline-block; max-width: 100%; width: 100%; text-align: center;">
        <video
          v-if="isWebcam"
          :ref="videoRef as any"
          autoplay
          playsinline
          muted
          class="viewport-img"
          style="display: block; max-width: 100%; width: 100%; height: auto; max-height: 520px; object-fit: contain; margin: 0 auto;"
        />
        <img
          v-else
          :src="imageSrc || '/hydra-logo.jpg'"
          alt="Tracking Feed"
          class="viewport-img"
          style="display: block; max-width: 100%; max-height: 520px; object-fit: contain; margin: 0 auto;"
          @error="(e) => (e.target as HTMLImageElement).src = '/hydra-logo.jpg'"
        />

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
            boxShadow: selectedEntity?.id === det.id ? '0 0 10px var(--cb-yellow)' : '0 0 6px rgba(0,240,255,0.4)',
            transition: 'left 0.06s linear, top 0.06s linear, width 0.06s linear, height 0.06s linear'
          }"
          @click="emit('update:selectedEntity', det)"
        >
          <span
            class="detection-tag"
            :style="{ background: selectedEntity?.id === det.id ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'), color: '#07080c' }"
          >
            {{ det.label.toUpperCase() }} {{ (det.conf * 100).toFixed(0) }}%
          </span>
        </div>
      </div>
    </div>

    <div class="viewport-footer-hud">
      <span>FEED: {{ isWebcam ? 'USB V4L2 WEBCAM' : (activeStream?.stream_id?.toUpperCase() || 'CAM_ENTRANCE_01') }}</span>
      <span>LATENCY: CUDA 13.3 (~3.8ms)</span>
      <span>RTX 5090 REALTIME TRACKER</span>
      <span style="color: var(--cb-green);">8.46 GB/s ZERO-COPY</span>
    </div>
  </div>
</template>
