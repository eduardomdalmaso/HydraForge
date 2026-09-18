import { ref, computed } from 'vue'
import { checkVaultHealth, fetchVaultDatasets, sendDirectFrameToVault, type VaultHealth, type VaultDataset } from '../../api/vault_client'

export interface SentFrameLog {
  id: string
  label: string
  conf: number
  timestamp: string
  status: 'synced' | 'error'
  thumbnailUrl?: string
}

export function useVaultCollector() {
  const isVaultOnline = ref(false)
  const vaultHealth = ref<VaultHealth | null>(null)
  const vaultDatasets = ref<VaultDataset[]>([])
  const selectedDatasetId = ref('frota_urbana_fusion')
  const activeClasses = ref<Record<string, boolean>>({})
  const maxPerClass = ref<number>(20)
  const onlyHardCases = ref<boolean>(false)
  const autoStreamToVault = ref<boolean>(true)
  const collectedStats = ref<Record<string, number>>({})
  const sentLogs = ref<SentFrameLog[]>([])
  const isSending = ref(false)
  let lastSentTime = 0

  const totalCollected = computed(() => Object.values(collectedStats.value).reduce((a, b) => a + b, 0))

  const checkHealth = async () => {
    const res = await checkVaultHealth()
    isVaultOnline.value = res.online
    vaultHealth.value = res
    if (res.online) {
      const dsets = await fetchVaultDatasets()
      if (dsets.length > 0) {
        vaultDatasets.value = dsets
        if (!vaultDatasets.value.some(d => d.dataset_id === selectedDatasetId.value)) {
          selectedDatasetId.value = vaultDatasets.value[0].dataset_id
        }
      }
    }
  }

  const resetStats = () => { collectedStats.value = {}; sentLogs.value = [] }

  const processFrameDetections = async (
    dets: any[], sourceEl: HTMLImageElement | HTMLVideoElement | null,
    sourceCategory: string, sourceFile: string, onGoalReached?: () => void
  ) => {
    if (!autoStreamToVault.value || !isVaultOnline.value || isSending.value) return
    const now = Date.now()
    if (now - lastSentTime < 350) return // Throttling

    const validDets = (dets || []).filter(d => {
      const lbl = String(d.label || d.class_name || '').toLowerCase()
      if (activeClasses.value[lbl] === false) return false
      const conf = d.conf ?? d.confidence ?? 0
      if (onlyHardCases.value && conf > 0.65) return false
      const current = collectedStats.value[lbl] || 0
      return !(maxPerClass.value > 0 && current >= maxPerClass.value)
    })

    if (validDets.length === 0) return
    const targetDet = validDets[0]
    const lbl = String(targetDet.label || targetDet.class_name || '').toLowerCase()

    isSending.value = true
    lastSentTime = now
    try {
      let blob: Blob | null = null
      let thumbUrl = ''
      if (sourceEl) {
        const canvas = document.createElement('canvas')
        const nw = (sourceEl as HTMLVideoElement).videoWidth || (sourceEl as HTMLImageElement).naturalWidth || 640
        const nh = (sourceEl as HTMLVideoElement).videoHeight || (sourceEl as HTMLImageElement).naturalHeight || 480
        canvas.width = Math.min(nw, 1280); canvas.height = Math.min(nh, 720)
        const ctx = canvas.getContext('2d')
        if (ctx) {
          ctx.drawImage(sourceEl, 0, 0, canvas.width, canvas.height)
          thumbUrl = canvas.toDataURL('image/jpeg', 0.6)
          blob = await new Promise<Blob | null>(r => canvas.toBlob(r, 'image/jpeg', 0.82))
        }
      }

      if (blob) {
        const bx = targetDet.box || [10, 10, 30, 30]
        const bboxPayload = [{
          class_id: 0, class_name: lbl,
          x_center: Math.min(0.99, Math.max(0.01, (bx[0] + bx[2] / 2) / 100)),
          y_center: Math.min(0.99, Math.max(0.01, (bx[1] + bx[3] / 2) / 100)),
          width: Math.min(0.99, Math.max(0.02, bx[2] / 100)),
          height: Math.min(0.99, Math.max(0.02, bx[3] / 100)),
          confidence: targetDet.conf ?? targetDet.confidence ?? 0.85
        }]

        const res = await sendDirectFrameToVault(selectedDatasetId.value, `cam_${sourceCategory}_${sourceFile}`, blob, bboxPayload)
        if (res.success) {
          collectedStats.value[lbl] = (collectedStats.value[lbl] || 0) + 1
          sentLogs.value = [{
            id: res.frameId || `frm_${Date.now()}`, label: lbl.toUpperCase(),
            conf: targetDet.conf ?? targetDet.confidence ?? 0.85,
            timestamp: new Date().toTimeString().split(' ')[0],
            status: 'synced' as const, thumbnailUrl: thumbUrl
          }, ...sentLogs.value].slice(0, 15)

          const activeKeys = Object.keys(activeClasses.value).filter(k => activeClasses.value[k] !== false)
          const allDone = activeKeys.length > 0 && activeKeys.every(k => (collectedStats.value[k] || 0) >= maxPerClass.value)
          if (allDone && maxPerClass.value > 0 && onGoalReached) onGoalReached()
        }
      }
    } finally {
      isSending.value = false
    }
  }

  return {
    isVaultOnline, vaultHealth, vaultDatasets, selectedDatasetId,
    activeClasses, maxPerClass, onlyHardCases, autoStreamToVault,
    collectedStats, totalCollected, sentLogs, isSending,
    checkHealth, resetStats, processFrameDetections
  }
}
