import { ref, watch, onUnmounted, type Ref } from 'vue'

export function useWebcamStream(isWebcam: Ref<boolean>) {
  const videoRef = ref<HTMLVideoElement | null>(null)
  let mediaStream: MediaStream | null = null

  const startStream = () => {
    if (navigator.mediaDevices?.getUserMedia) {
      navigator.mediaDevices.getUserMedia({ video: { width: 640, height: 480 } })
        .then(s => {
          mediaStream = s
          if (videoRef.value) videoRef.value.srcObject = s
        })
        .catch(err => console.warn('Webcam permission:', err))
    }
  }

  const stopStream = () => {
    if (mediaStream) {
      mediaStream.getTracks().forEach(t => t.stop())
      mediaStream = null
    }
  }

  watch(isWebcam, (active) => {
    if (active) startStream()
    else stopStream()
  }, { immediate: true })

  onUnmounted(() => {
    stopStream()
  })

  const captureFrame = (): string | null => {
    const el = videoRef.value || (document.getElementById('hud-viewport-video') as HTMLVideoElement)
    if (el && el.videoWidth > 0 && el.videoHeight > 0) {
      const canvas = document.createElement('canvas')
      const targetW = Math.min(el.videoWidth, 640)
      const targetH = Math.round(targetW * (el.videoHeight / el.videoWidth))
      canvas.width = targetW
      canvas.height = targetH
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.drawImage(el, 0, 0, targetW, targetH)
        return canvas.toDataURL('image/jpeg', 0.8)
      }
    }
    return null
  }

  return { videoRef, captureFrame }
}
