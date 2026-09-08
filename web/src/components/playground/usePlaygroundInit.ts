import { ref, onMounted, onUnmounted, type Ref } from 'vue'
import { fetchHydraStreams, fetchHydraTelemetry } from '../../api/hydrastream_client'
import { fetchModelsAPI, fetchTelemetryAPI } from '../../api/client'

export function usePlaygroundInit(config: Ref<any>) {
  const modelsList = ref<any[]>([])
  const hydraStreams = ref<any[]>([])
  const hydraTelemetry = ref<any>(null)
  const gpuStats = ref<any>(null)
  let pollTimer: any = null

  onMounted(async () => {
    const [models, streams, stats, hw] = await Promise.all([
      fetchModelsAPI(), fetchHydraStreams(), fetchHydraTelemetry(), fetchTelemetryAPI()
    ])
    if (models?.length) {
      modelsList.value = models
      const b = models.find((m: any) => m.isCustom) || models[0]
      config.value.model = b.id
    }
    if (streams?.length) hydraStreams.value = streams
    if (stats) hydraTelemetry.value = stats
    if (hw?.gpu_stats) gpuStats.value = hw.gpu_stats

    pollTimer = setInterval(async () => {
      const [s, h] = await Promise.all([fetchHydraTelemetry(), fetchTelemetryAPI()])
      if (s) hydraTelemetry.value = s
      if (h?.gpu_stats) gpuStats.value = h.gpu_stats
    }, 3000)
  })

  onUnmounted(() => {
    if (pollTimer) clearInterval(pollTimer)
  })

  return { modelsList, hydraStreams, hydraTelemetry, gpuStats }
}
