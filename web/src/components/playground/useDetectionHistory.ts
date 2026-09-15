import { ref } from 'vue'

export interface DetectionEvent {
  id: string
  label: string
  conf: number
  box: [number, number, number, number]
  timestamp: string
  color: string
  thumbnailUrl?: string
}

export function useDetectionHistory() {
  const history = ref<DetectionEvent[]>([])
  const maxEvents = 30

  const cropThumbnail = (
    sourceEl: HTMLImageElement | HTMLVideoElement | null,
    box: [number, number, number, number],
    label: string
  ): string => {
    try {
      const canvas = document.createElement('canvas')
      canvas.width = 72
      canvas.height = 72
      const ctx = canvas.getContext('2d')
      if (!ctx) return ''

      if (sourceEl) {
        const nw = (sourceEl as HTMLVideoElement).videoWidth || (sourceEl as HTMLImageElement).naturalWidth || sourceEl.clientWidth || 640
        const nh = (sourceEl as HTMLVideoElement).videoHeight || (sourceEl as HTMLImageElement).naturalHeight || sourceEl.clientHeight || 480
        const sx = Math.max(0, (box[0] / 100) * nw)
        const sy = Math.max(0, (box[1] / 100) * nh)
        const sw = Math.max(10, (box[2] / 100) * nw)
        const sh = Math.max(10, (box[3] / 100) * nh)
        ctx.drawImage(sourceEl, sx, sy, sw, sh, 0, 0, 72, 72)
      } else {
        ctx.fillStyle = '#0e1117'
        ctx.fillRect(0, 0, 72, 72)
        ctx.fillStyle = '#ff5e3a'
        ctx.font = 'bold 10px JetBrains Mono, monospace'
        ctx.textAlign = 'center'
        ctx.fillText(label.slice(0, 8).toUpperCase(), 36, 40)
      }
      return canvas.toDataURL('image/jpeg', 0.82)
    } catch {
      return ''
    }
  }

  const pushDetections = (
    newDets: any[],
    sourceEl: HTMLImageElement | HTMLVideoElement | null = null
  ) => {
    if (!newDets || newDets.length === 0) return

    const now = new Date()
    const timeStr = now.toTimeString().split(' ')[0] + '.' + Math.floor(now.getMilliseconds() / 100)

    const events: DetectionEvent[] = newDets.map((d, i) => {
      const evtId = d.id ? `${d.id}` : `det_${Date.now()}_${i}`
      const thumb = cropThumbnail(sourceEl, d.box || [10, 10, 30, 30], d.label || 'target')
      return {
        id: evtId,
        label: (d.label || 'object').toUpperCase(),
        conf: d.conf || 0.85,
        box: d.box || [10, 10, 30, 30],
        timestamp: timeStr,
        color: d.color || '#00f0ff',
        thumbnailUrl: thumb
      }
    })

    history.value = [...events, ...history.value].slice(0, maxEvents)
  }

  const clearHistory = () => {
    history.value = []
  }

  return {
    history,
    pushDetections,
    clearHistory
  }
}
