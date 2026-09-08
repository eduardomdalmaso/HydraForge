import { watch, onUnmounted, type Ref } from 'vue'

export function usePlaygroundStream(
  config: Ref<any>,
  isContinuous: Ref<boolean>,
  isWebcam: Ref<boolean>,
  captureFrame: () => string | null,
  detections: Ref<any[]>,
  telemetry: Ref<any>,
  imageSrc: Ref<string>
) {
  let eventSource: EventSource | null = null
  let webcamInterval: any = null

  const startLiveStream = () => {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
    if (webcamInterval) {
      clearInterval(webcamInterval)
      webcamInterval = null
    }

    if (!isContinuous.value) {
      if (!isWebcam.value) {
        imageSrc.value = `/api/v1/hydrastream/api/v1/streams/${config.value.source}/snapshot.jpg?t=${Date.now()}`
      }
      return
    }

    if (!isWebcam.value) {
      imageSrc.value = `/api/v1/hydrastream/api/v1/streams/${config.value.source}/mjpeg?t=${Date.now()}`
    }

    if (isWebcam.value) {
      webcamInterval = setInterval(() => {
        const b64 = captureFrame()
        if (b64) {
          fetch('/api/v1/inference/frame', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ image_base64: b64 })
          }).catch(() => {})
        }
      }, 40)
    }

    eventSource = new EventSource(`/api/v1/inference/live?model=${encodeURIComponent(config.value.model)}&source=${encodeURIComponent(config.value.source)}&conf=${config.value.conf}`)
    eventSource.onmessage = (e) => {
      try {
        const d = JSON.parse(e.data)
        if (Array.isArray(d.detections)) detections.value = d.detections
        if (d.telemetry) telemetry.value = d.telemetry
      } catch {}
    }
  }

  watch(() => [isContinuous.value, config.value.source, config.value.model, config.value.conf], () => {
    startLiveStream()
  })

  onUnmounted(() => {
    if (eventSource) eventSource.close()
    if (webcamInterval) clearInterval(webcamInterval)
  })

  return { startLiveStream }
}
