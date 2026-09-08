<script setup lang="ts">
import { ref, computed } from 'vue'
import PlaygroundControlBar from '../components/playground/PlaygroundControlBar.vue'
import PlaygroundTuningCard from '../components/playground/PlaygroundTuningCard.vue'
import PlaygroundViewport from '../components/playground/PlaygroundViewport.vue'
import PlaygroundTelemetryCard from '../components/playground/PlaygroundTelemetryCard.vue'
import { useWebcamStream } from '../components/playground/useWebcamStream'
import { usePlaygroundStream } from '../components/playground/usePlaygroundStream'
import { usePlaygroundInit } from '../components/playground/usePlaygroundInit'
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

const handleInference = async () => {
  const b64 = captureFrame()
  const res = await runRealInferenceAPI(config.value, b64)
  if (res) {
    if (res.displayImageUrl && !isWebcam.value && !isContinuous.value) imageSrc.value = res.displayImageUrl
    if (Array.isArray(res.detections)) detections.value = res.detections
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

    <div class="playground-layout">
      <div style="display: flex; flex-direction: column; gap: 1.25rem;">
        <PlaygroundControlBar
          :config="config"
          :modelsList="modelsList"
          :hydraStreams="hydraStreams"
          :isHydraOnline="hydraStreams.length > 0"
          :isRunning="isScanning"
          :isContinuous="isContinuous"
          @update:config="(c) => config = c as any"
          @update:isContinuous="(v) => isContinuous = v"
          @runInference="async () => { isScanning = true; await handleInference(); isScanning = false; }"
        />
        <PlaygroundTuningCard :config="config" @update:config="(c) => config = c as any" />
      </div>

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
        @update:selectedEntity="(e) => selectedEntity = e"
      />

      <PlaygroundTelemetryCard
        :telemetry="telemetry"
        :detections="detections"
        :selectedEntity="selectedEntity"
        :hydraTelemetry="hydraTelemetry"
        :activeStream="activeStream"
        :gpuStats="gpuStats"
        @selectEntity="(e) => selectedEntity = e"
      />
    </div>
  </div>
</template>
