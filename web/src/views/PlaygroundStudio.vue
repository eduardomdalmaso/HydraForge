<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import PlaygroundViewport from '../components/playground/PlaygroundViewport.vue'
import PlaygroundUnifiedPanel from '../components/playground/PlaygroundUnifiedPanel.vue'
import PlaygroundDetectionFeed from '../components/playground/PlaygroundDetectionFeed.vue'
import { useWebcamStream } from '../components/playground/useWebcamStream'
import { usePlaygroundStream } from '../components/playground/usePlaygroundStream'
import { usePlaygroundInit } from '../components/playground/usePlaygroundInit'
import { useDetectionHistory } from '../components/playground/useDetectionHistory'
import { runRealInferenceAPI } from '../api/inference_client'

const config = ref({
  model: 'yolo26n',
  source: 'cam_entrance_01',
  runtime: 'pytorch',
  conf: 0.25,
  iou: 0.45,
  nmsFree: true,
  sahi: false,
  isolateBg: false
})

const { modelsList, hydraStreams, hydraTelemetry, gpuStats } = usePlaygroundInit(config)
const detections = ref<any[]>([])
const selectedEntity = ref<any>(null)
const isScanning = ref(false)
const isContinuous = ref(false)
const telemetry = ref<any>(null)
const imageSrc = ref('/api/v1/hydrastream/api/v1/streams/cam_entrance_01/snapshot.jpg')

const isWebcam = computed(() => config.value.source === 'webcam')
const activeStream = computed(() => hydraStreams.value.find(s => s.stream_id === config.value.source) || (hydraStreams.value[0] || null))

const { videoRef, captureFrame } = useWebcamStream(isWebcam)
usePlaygroundStream(config, isContinuous, isWebcam, captureFrame, detections, telemetry, imageSrc)

const { history, pushDetections, clearHistory } = useDetectionHistory()

watch(detections, (newDets) => {
  if (newDets && newDets.length > 0) {
    const el = (isWebcam.value ? document.getElementById('hud-viewport-video') : document.getElementById('hud-viewport-image')) as any
    pushDetections(newDets, el)
  }
}, { deep: true })

const handleInference = async () => {
  const b64 = captureFrame()
  const res = await runRealInferenceAPI(config.value, b64)
  if (res) {
    if (res.displayImageUrl && !isWebcam.value && !isContinuous.value) imageSrc.value = res.displayImageUrl
    if (Array.isArray(res.detections)) {
      detections.value = res.detections
      const el = (isWebcam.value ? document.getElementById('hud-viewport-video') : document.getElementById('hud-viewport-image')) as any
      pushDetections(res.detections, el)
    }
    if (res.telemetry) telemetry.value = res.telemetry
  }
}
</script>

<template>
  <div class="view-container playground-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">PLAYGROUND DE INFERENCIA & TRACKING</h1>
      <p class="cockpit-main-subtitle">HARDWARE NATIVO RTX 5090 // FLUXO CONTINUO HYDRASTREAM // INFERENCIA AO VIVO</p>
    </div>

    <!-- TOP ROW: LEFT VIEWPORT (720P) + RIGHT CONTROLS PANEL -->
    <div class="playground-top-row">
      <div class="playground-viewport-col">
        <PlaygroundViewport
          :imageSrc="imageSrc"
          :detections="detections"
          :selectedEntity="selectedEntity"
          :isScanning="isScanning"
          :isContinuous="isContinuous"
          :isHydraLinked="!isWebcam && hydraStreams.length > 0"
          :activeStream="activeStream"
          :isWebcam="isWebcam"
          :videoRef="videoRef"
          :telemetry="telemetry"
          :gpuStats="gpuStats"
          :config="config"
          @update:selectedEntity="(e) => selectedEntity = e"
        />
      </div>

      <div class="playground-controls-col">
        <PlaygroundUnifiedPanel
          :config="config"
          :modelsList="modelsList"
          :hydraStreams="hydraStreams"
          :isHydraOnline="hydraStreams.length > 0"
          :isRunning="isScanning"
          :isContinuous="isContinuous"
          :telemetry="telemetry"
          :gpuStats="gpuStats"
          :hydraTelemetry="hydraTelemetry"
          @update:config="(c) => config = c as any"
          @update:isContinuous="(v) => isContinuous = v"
          @runInference="async () => { isScanning = true; await handleInference(); isScanning = false; }"
        />
      </div>
    </div>

    <!-- BOTTOM ROW: FULL-WIDTH REALTIME DETECTION FEED -->
    <div class="playground-bottom-row">
      <PlaygroundDetectionFeed
        :events="history"
        :selectedEntity="selectedEntity"
        @selectEntity="(e) => selectedEntity = e"
        @clearHistory="clearHistory"
      />
    </div>
  </div>
</template>
