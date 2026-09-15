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
  let frameInterval: any = null

  const startLiveStream = () => {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
    if (frameInterval) {
      clearInterval(frameInterval)
      frameInterval = null
    }

    const isVideoLoop = config.value.source?.startsWith('video:')
    const isFrameUploader = isWebcam.value || isVideoLoop

    if (!isContinuous.value) {
      if (!isFrameUploader) {
        imageSrc.value = `/api/v1/hydrastream/api/v1/streams/${config.value.source}/snapshot.jpg?t=${Date.now()}`
      }
      return
    }

    if (!isFrameUploader) {
      imageSrc.value = `/api/v1/hydrastream/api/v1/streams/${config.value.source}/mjpeg?t=${Date.now()}`
    } else {
      frameInterval = setInterval(() => {
        const b64 = captureFrame()
        if (b64) {
          fetch('/api/v1/inference/frame', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ image_base64: b64 })
          }).catch(() => {})
        }
      }, 50)
    }

    const liveSrc = isFrameUploader ? 'webcam' : config.value.source
    eventSource = new EventSource(`/api/v1/inference/live?model=${encodeURIComponent(config.value.model)}&source=${encodeURIComponent(liveSrc)}&conf=${config.value.conf}`)
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
    if (frameInterval) clearInterval(frameInterval)
  })

  return { startLiveStream }
}
