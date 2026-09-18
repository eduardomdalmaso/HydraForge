<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import PlaygroundViewport from '../components/playground/PlaygroundViewport.vue'
import PlaygroundUnifiedPanel from '../components/playground/PlaygroundUnifiedPanel.vue'
import PlaygroundDetectionInspector from '../components/playground/PlaygroundDetectionInspector.vue'
import PlaygroundDetectionFeed from '../components/playground/PlaygroundDetectionFeed.vue'
import PlaygroundMediaModal from '../components/playground/PlaygroundMediaModal.vue'
import { useWebcamStream } from '../components/playground/useWebcamStream'
import { usePlaygroundStream } from '../components/playground/usePlaygroundStream'
import { usePlaygroundInit } from '../components/playground/usePlaygroundInit'
import { useVaultCollector } from '../components/playground/useVaultCollector'
import { fetchMediaSources, type MediaFolder } from '../api/media_client'
import { runRealInferenceAPI } from '../api/inference_client'

const config = ref({
  model: 'yolo26n',
  source: 'folder:carros',
  runtime: 'pytorch',
  conf: 0.25,
  iou: 0.45,
  nmsFree: true,
  sahi: false,
  isolateBg: false
})

const { modelsList, hydraStreams, gpuStats } = usePlaygroundInit(config)
const detections = ref<any[]>([])
const selectedEntity = ref<any>(null)
const isScanning = ref(false)
const isContinuous = ref(false)
const isPaused = ref(false)
const telemetry = ref<any>(null)
const imageSrc = ref('/api/v1/hydrastream/api/v1/streams/cam_entrance_01/snapshot.jpg')
const mediaFolders = ref<MediaFolder[]>([])
const showMediaModal = ref(false)

const {
  isVaultOnline, vaultDatasets, selectedDatasetId, activeClasses,
  maxPerClass, onlyHardCases, hardCaseThreshold, autoStreamToVault, collectedStats,
  totalCollected, sentLogs, isSending, checkHealth, resetStats,
  processFrameDetections
} = useVaultCollector()

const loadMediaFolders = async () => {
  const data = await fetchMediaSources()
  mediaFolders.value = data.folders || []
  if (config.value.source === 'folder:carros' && !mediaFolders.value.some(f => f.name === 'carros')) {
    const firstWithFiles = mediaFolders.value.find(f => f.file_count > 0)
    if (firstWithFiles) config.value.source = `folder:${firstWithFiles.name}`
  }
}

onMounted(() => {
  loadMediaFolders()
  checkHealth()
})

const isWebcam = computed(() => config.value.source === 'webcam')
const isVideoMedia = computed(() => config.value.source?.startsWith('folder:') || config.value.source?.startsWith('video:'))
const activeStream = computed(() => hydraStreams.value.find(s => s.stream_id === config.value.source) || (hydraStreams.value[0] || null))

const { videoRef, captureFrame } = useWebcamStream(isWebcam)
usePlaygroundStream(config, isContinuous, isWebcam, captureFrame, detections, telemetry, imageSrc, isPaused)

const onGoalReached = () => {
  isContinuous.value = false
  isPaused.value = true
}

watch(detections, (newDets) => {
  if (newDets && newDets.length > 0 && !isPaused.value) {
    const el = (isWebcam.value || isVideoMedia.value)
      ? document.getElementById('hud-viewport-video') as HTMLVideoElement
      : document.getElementById('hud-viewport-image') as any
    const src = config.value.source || ''
    const cat = src.startsWith('folder:') ? src.replace('folder:', '') : 'playground'
    const fn = src.includes('/') ? src.split('/').pop() || 'stream' : src
    processFrameDetections(newDets, el, cat, fn, onGoalReached)
  }
}, { deep: true })

const handleInference = async () => {
  if (!isVaultOnline.value) return
  const b64 = captureFrame()
  const res = await runRealInferenceAPI(config.value, b64)
  if (res) {
    if (res.displayImageUrl && !isWebcam.value && !isVideoMedia.value && !isContinuous.value) {
      imageSrc.value = res.displayImageUrl
    }
    if (Array.isArray(res.detections)) {
      detections.value = res.detections
      const el = (isWebcam.value || isVideoMedia.value)
        ? document.getElementById('hud-viewport-video') as HTMLVideoElement
        : document.getElementById('hud-viewport-image') as any
      const src = config.value.source || ''
      const cat = src.startsWith('folder:') ? src.replace('folder:', '') : 'playground'
      const fn = src.includes('/') ? src.split('/').pop() || 'stream' : src
      processFrameDetections(res.detections, el, cat, fn, onGoalReached)
    }
    if (res.telemetry) telemetry.value = res.telemetry
  }
}
</script>

<template>
  <div class="view-container playground-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">PLAYGROUND DE INFERENCIA & TRACKING</h1>
      <p class="cockpit-main-subtitle">HARDWARE NATIVO RTX 5090 // STREAM DIRETO HYDRAVAULT // CURADORIA ATIVA</p>
    </div>

    <!-- TOP ROW: LEFT VIEWPORT + RIGHT CONTROLS PANEL -->
    <div class="playground-top-row">
      <div class="playground-viewport-col">
        <PlaygroundViewport
          :imageSrc="imageSrc"
          :detections="detections"
          :selectedEntity="selectedEntity"
          :isScanning="isScanning"
          :isContinuous="isContinuous"
          :isPaused="isPaused"
          :isHydraLinked="!isWebcam && !isVideoMedia && hydraStreams.length > 0"
          :activeStream="activeStream"
          :isWebcam="isWebcam"
          :videoRef="videoRef"
          :telemetry="telemetry"
          :gpuStats="gpuStats"
          :config="config"
          :mediaFolders="mediaFolders"
          @update:selectedEntity="(e) => selectedEntity = e"
          @update:isPaused="(p) => isPaused = p"
        />
      </div>

      <div class="playground-controls-col">
        <PlaygroundUnifiedPanel
          :config="config"
          :modelsList="modelsList"
          :hydraStreams="hydraStreams"
          :mediaFolders="mediaFolders"
          :isRunning="isScanning"
          :isContinuous="isContinuous"
          :isPaused="isPaused"
          :isVaultOnline="isVaultOnline"
          :vaultDatasets="vaultDatasets"
          :selectedDatasetId="selectedDatasetId"
          :activeClasses="activeClasses"
          :maxPerClass="maxPerClass"
          :onlyHardCases="onlyHardCases"
          :hardCaseThreshold="hardCaseThreshold"
          :autoStreamToVault="autoStreamToVault"
          :collectedStats="collectedStats"
          :totalCollected="totalCollected"
          @update:config="(c) => config = c as any"
          @update:isContinuous="(v) => isContinuous = v"
          @update:isPaused="(p) => isPaused = p"
          @update:selectedDatasetId="(id) => selectedDatasetId = id"
          @toggleTargetClass="(cls) => activeClasses[cls] = activeClasses[cls] === false"
          @update:maxPerClass="(m) => maxPerClass = m"
          @update:onlyHardCases="(v) => onlyHardCases = v"
          @update:hardCaseThreshold="(v) => hardCaseThreshold = v"
          @update:autoStreamToVault="(v) => autoStreamToVault = v"
          @resetCollectorStats="resetStats"
          @checkVaultHealth="checkHealth"
          @runInference="async () => { isScanning = true; await handleInference(); isScanning = false; }"
          @openMediaModal="showMediaModal = true"
        />
      </div>

      <!-- RIGHT COLUMN: DETECTION INSPECTOR -->
      <div v-if="selectedEntity" class="playground-inspector-col">
        <PlaygroundDetectionInspector :entity="selectedEntity" @close="selectedEntity = null" />
      </div>
    </div>

    <!-- BOTTOM ROW: VAULT TRANSMISSION STREAM FEED -->
    <div class="playground-bottom-row">
      <PlaygroundDetectionFeed
        :sentLogs="sentLogs"
        :totalCollected="totalCollected"
        :isSending="isSending"
        :selectedEntity="selectedEntity"
        @selectEntity="(e) => selectedEntity = e"
        @clearHistory="resetStats"
      />
    </div>

    <!-- MEDIA MANAGER MODAL -->
    <PlaygroundMediaModal v-if="showMediaModal" @close="showMediaModal = false" @mediaUpdated="loadMediaFolders" />
  </div>
</template>
