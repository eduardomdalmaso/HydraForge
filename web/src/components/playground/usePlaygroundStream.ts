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
  let isLoopActive = false
  let currentRunId = 0

  const runVideoFrameLoop = async (runId: number) => {
    if (!isLoopActive || runId !== currentRunId) return

    const b64 = captureFrame()
    if (b64) {
      try {
        const res = await fetch('/api/v1/inference/predict', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            model: config.value.model || 'yolo26n',
            image_base64: b64,
            conf: config.value.conf || 0.25,
            iou: config.value.iou || 0.45,
            sahi: !!config.value.sahi,
            nms_free: !!config.value.nmsFree,
            track: true,
            device: '0'
          })
        })
        if (res.ok && isLoopActive && runId === currentRunId) {
          const d = await res.json()
          if (Array.isArray(d.detections)) detections.value = d.detections
          if (d.telemetry) telemetry.value = d.telemetry
        }
      } catch {}
    }

    if (isLoopActive && runId === currentRunId) {
      setTimeout(() => {
        requestAnimationFrame(() => runVideoFrameLoop(runId))
      }, 30)
    }
  }

  const startLiveStream = () => {
    isLoopActive = false
    currentRunId++

    if (eventSource) {
      eventSource.close()
      eventSource = null
    }

    const isVideoLoop = config.value.source?.startsWith('video:') || config.value.source?.startsWith('folder:')
    const isFrameUploader = isWebcam.value || isVideoLoop

    if (!isContinuous.value) {
      if (!isFrameUploader) {
        imageSrc.value = `/api/v1/hydrastream/api/v1/streams/${config.value.source}/snapshot.jpg?t=${Date.now()}`
      }
      return
    }

    if (isFrameUploader) {
      isLoopActive = true
      runVideoFrameLoop(currentRunId)
    } else {
      imageSrc.value = `/api/v1/hydrastream/api/v1/streams/${config.value.source}/mjpeg?t=${Date.now()}`
      const sParams = new URLSearchParams({
        model: config.value.model || 'yolo26n',
        source: config.value.source || 'cam_entrance_01',
        conf: String(config.value.conf || 0.25),
        iou: String(config.value.iou || 0.45),
        sahi: config.value.sahi ? 'true' : 'false',
        nms_free: config.value.nmsFree ? 'true' : 'false'
      })
      eventSource = new EventSource(`/api/v1/inference/live?${sParams.toString()}`)
      eventSource.onmessage = (e) => {
        try {
          const d = JSON.parse(e.data)
          if (Array.isArray(d.detections)) detections.value = d.detections
          if (d.telemetry) telemetry.value = d.telemetry
        } catch {}
      }
    }
  }

  watch(() => [
    isContinuous.value,
    config.value.source,
    config.value.model,
    config.value.conf,
    config.value.iou,
    config.value.sahi,
    config.value.nmsFree
  ], () => {
    startLiveStream()
  })

  onUnmounted(() => {
    isLoopActive = false
    currentRunId++
    if (eventSource) eventSource.close()
  })

  return { startLiveStream }
}
