<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { MediaFolder } from '../../api/media_client'

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
  mediaFolders?: MediaFolder[]
}>()

const emit = defineEmits<{ (e: 'update:selectedEntity', det: any): void }>()
const imgError = ref(false)
watch(() => props.imageSrc, () => { imgError.value = false })

const isFolderLoop = computed(() => props.config?.source?.startsWith('folder:'))
const activeFolderName = computed(() => isFolderLoop.value ? props.config?.source?.replace('folder:', '') : '')
const activeFolder = computed(() => (props.mediaFolders || []).find(f => f.name === activeFolderName.value) || null)
const folderFiles = computed(() => activeFolder.value?.files || [])

const currentFileIndex = ref(0)
watch(() => props.config?.source, () => { currentFileIndex.value = 0 })
const currentFile = computed(() => folderFiles.value[currentFileIndex.value] || null)
const isVideoSource = computed(() => props.isWebcam || isFolderLoop.value || props.config?.source?.startsWith('video:'))

const videoSrc = computed(() => {
  if (isFolderLoop.value && currentFile.value) return currentFile.value.stream_url
  if (props.config?.source?.startsWith('video:')) return `/api/v1/media/stream/${props.config.source.replace(/^video:/, '')}`
  return undefined
})

const handleVideoEnded = () => {
  if (isFolderLoop.value && folderFiles.value.length > 0) {
    currentFileIndex.value = (currentFileIndex.value + 1) % folderFiles.value.length
  }
}
const prevVideo = () => {
  if (folderFiles.value.length > 0) {
    currentFileIndex.value = (currentFileIndex.value - 1 + folderFiles.value.length) % folderFiles.value.length
  }
}
const nextVideo = () => {
  if (folderFiles.value.length > 0) {
    currentFileIndex.value = (currentFileIndex.value + 1) % folderFiles.value.length
  }
}

const mediaAspectRatio = ref<string>('16 / 9')
const onVideoMetadata = (e: Event) => {
  const el = e.target as HTMLVideoElement
  if (el.videoWidth && el.videoHeight) mediaAspectRatio.value = `${el.videoWidth} / ${el.videoHeight}`
}
const onImageLoaded = (e: Event) => {
  const el = e.target as HTMLImageElement
  if (el.naturalWidth && el.naturalHeight) mediaAspectRatio.value = `${el.naturalWidth} / ${el.naturalHeight}`
}
</script>

<template>
  <div class="viewport-hud">
    <div class="viewport-top-bar">
      <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
        <span>// KIROSHI OPTICAL HUD</span>
        <span
          :style="{
            fontSize: '0.65rem', padding: '2px 6px', borderRadius: '2px',
            background: isWebcam ? 'rgba(0,240,255,0.15)' : (isFolderLoop ? 'rgba(252,238,10,0.15)' : 'rgba(0,255,157,0.15)'),
            color: isWebcam ? 'var(--cb-cyan)' : (isFolderLoop ? 'var(--cb-yellow)' : 'var(--cb-green)'),
            border: `1px solid ${isWebcam ? 'var(--cb-cyan)' : (isFolderLoop ? 'var(--cb-yellow)' : 'var(--cb-green)')}`
          }"
        >
          {{ isWebcam ? '[HARDWARE] WEBCAM' : (isFolderLoop ? `[PASTA // ${activeFolderName.toUpperCase()}]` : (isHydraLinked ? `[STREAM] ${activeStream?.stream_id?.toUpperCase() || 'SHM'}` : '[STANDBY]')) }}
        </span>

        <div v-if="isFolderLoop && folderFiles.length > 1" style="display: flex; align-items: center; gap: 0.3rem;">
          <button class="cyber-pill" style="padding: 1px 5px; font-size: 0.62rem;" title="Vídeo anterior" @click="prevVideo">◀</button>
          <span style="font-size: 0.65rem; color: #fff; font-family: var(--font-mono);">{{ currentFileIndex + 1 }}/{{ folderFiles.length }}: {{ currentFile?.name }}</span>
          <button class="cyber-pill" style="padding: 1px 5px; font-size: 0.62rem;" title="Próximo vídeo" @click="nextVideo">▶</button>
        </div>
        <span v-else-if="isFolderLoop && currentFile" style="font-size: 0.65rem; color: var(--cb-muted);">{{ currentFile.name }}</span>
      </div>

      <span style="color: var(--cb-yellow); font-family: var(--font-mono); font-size: 0.72rem;">
        {{ isContinuous ? '● REALTIME STREAM TRACKING' : (isScanning ? 'INFERENCE...' : `TARGETS: ${detections?.length || 0}`) }}
      </span>
    </div>

    <div class="viewport-stage">
      <div v-if="isVideoSource || (imageSrc && !imgError && imageSrc !== '/hydra-logo.jpg')" class="viewport-media-wrapper" :style="{ aspectRatio: mediaAspectRatio }">
        <video
          v-if="isVideoSource"
          id="hud-viewport-video"
          :ref="videoRef as any"
          :src="videoSrc"
          autoplay
          :loop="!isFolderLoop || folderFiles.length <= 1"
          playsinline
          muted
          class="viewport-media"
          @loadedmetadata="onVideoMetadata"
          @ended="handleVideoEnded"
        />
        <img v-else id="hud-viewport-image" :src="imageSrc" alt="Tracking Feed" class="viewport-media" @load="onImageLoaded" @error="imgError = true" />

        <div
          v-for="(det, idx) in (detections || [])"
          :key="det.id || idx"
          class="detection-box"
          :style="{
            left: `${det.box[0]}%`, top: `${det.box[1]}%`, width: `${det.box[2]}%`, height: `${det.box[3]}%`,
            borderColor: selectedEntity?.id === det.id ? 'var(--cb-yellow)' : (det.color || 'var(--cb-cyan)'),
            boxShadow: selectedEntity?.id === det.id ? '0 0 10px var(--cb-yellow)' : '0 0 6px rgba(0,240,255,0.4)'
          }"
          @click.stop="emit('update:selectedEntity', det)"
        >
          <span class="detection-tag">
            <template v-if="det.track_id">#{{ det.track_id }} </template>
            {{ (det.label || det.class_name || 'TARGET').toUpperCase() }} // {{ (((det.conf ?? det.confidence ?? 0.85)) * 100).toFixed(0) }}%
          </span>
        </div>
      </div>

      <div v-else class="camera-placeholder-hud">
        <div style="position: relative; margin-bottom: 1.25rem;">
          <svg width="68" height="68" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" style="filter: drop-shadow(0 0 14px rgba(255, 94, 58, 0.55));">
            <path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z" />
            <circle cx="12" cy="13" r="4" />
            <circle cx="12" cy="13" r="1.5" fill="#ff5e3a" />
          </svg>
          <div style="position: absolute; top: -6px; right: -6px; width: 8px; height: 8px; border-radius: 50%; background: #ff5e3a; box-shadow: 0 0 8px #ff5e3a;" />
        </div>
        <div style="color: #ff5e3a; font-family: var(--font-mono); font-size: 0.85rem; font-weight: 600; letter-spacing: 0.08em; margin-bottom: 0.4rem;">[HYDRASTREAM SENSOR FEED // 16:9]</div>
        <div style="color: #8b94a0; font-family: var(--font-mono); font-size: 0.72rem; letter-spacing: 0.04em;">{{ isHydraLinked ? '// FLUXO DE VIDEO ZERO-COPY CALIBRADO' : '// SELECIONE UMA PASTA DE VIDEOS OU ATIVE A WEBCAM' }}</div>
      </div>
    </div>

    <div class="viewport-footer-hud">
      <div style="display: flex; gap: 0.6rem; align-items: center;">
        <span style="color: var(--cb-green); font-weight: 700;">[16:9 NATIVE HUD]</span>
        <span>LATENCY: {{ telemetry?.inference_ms ? `${telemetry.inference_ms}ms` : '0.92ms' }}</span>
      </div>
      <div style="color: var(--cb-cyan);">
        ENGINE: {{ (config?.runtime || 'PYTORCH').toUpperCase() }} // {{ (config?.model || 'YOLO26M').toUpperCase() }}
      </div>
    </div>
  </div>
</template>
