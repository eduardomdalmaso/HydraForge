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
    if (isWebcam.value && videoRef.value && videoRef.value.videoWidth) {
      const canvas = document.createElement('canvas')
      canvas.width = videoRef.value.videoWidth
      canvas.height = videoRef.value.videoHeight
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.drawImage(videoRef.value, 0, 0)
        return canvas.toDataURL('image/jpeg', 0.85)
      }
    }
    return null
  }

  return { videoRef, captureFrame }
}
