import { ref, watch, onMounted, type Ref } from 'vue'
import { fetchDatasetsAPI, rescanDatasetsAPI, deleteDatasetAPI } from '../../api/client'
import type { DatasetInfo } from '../../types/dataset'

export function useDatasetStudio(propsDatasets: Ref<DatasetInfo[] | undefined>) {
  const datasets = ref<DatasetInfo[]>(propsDatasets.value || [])
  const selectedDataset = ref<DatasetInfo | null>(propsDatasets.value?.[0] || null)
  const isModalOpen = ref(false)
  const isAnnotateOpen = ref(false)
  const isMergeOpen = ref(false)
  const deleteTarget = ref<any>(null)
  const deleteDisk = ref(false)
  const isDeleting = ref(false)

  const globalMappings = ref<Record<string, any>>(() => {
    try { return JSON.parse(localStorage.getItem('hydraforge_mappings') || '{}') } catch { return {} }
  })

  const saveMapping = (dsId: string, mapping: Record<string, string>) => {
    const next = { ...globalMappings.value, [dsId]: { ...(globalMappings.value[dsId] || {}), ...mapping } }
    globalMappings.value = next
    try { localStorage.setItem('hydraforge_mappings', JSON.stringify(next)) } catch {}
  }

  const handleRescan = async () => {
    const data = await rescanDatasetsAPI()
    if (data?.length) {
      datasets.value = data
      if (!selectedDataset.value || !data.some(d => d.id === selectedDataset.value?.id)) {
        selectedDataset.value = data[0]
      }
    }
  }

  const confirmDelete = async () => {
    if (!deleteTarget.value) return
    isDeleting.value = true
    const id = deleteTarget.value.id || deleteTarget.value.dataset_id
    if (await deleteDatasetAPI(id, deleteDisk.value)) {
      const rem = datasets.value.filter(d => (d.id || (d as any).dataset_id) !== id)
      datasets.value = rem
      if ((selectedDataset.value?.id || (selectedDataset.value as any)?.dataset_id) === id) {
        selectedDataset.value = rem[0] || null
      }
    }
    isDeleting.value = false
    deleteTarget.value = null
  }

  watch(() => propsDatasets.value, (newDs) => {
    if (newDs && newDs.length > 0) {
      datasets.value = newDs
      if (!selectedDataset.value) selectedDataset.value = newDs[0]
    }
  }, { immediate: true })

  onMounted(() => {
    if (!propsDatasets.value || propsDatasets.value.length === 0) {
      fetchDatasetsAPI().then(data => {
        if (data && data.length > 0) {
          datasets.value = data
          selectedDataset.value = data[0]
        }
      })
    }
  })

  return {
    datasets, selectedDataset, isModalOpen, isAnnotateOpen, isMergeOpen,
    deleteTarget, deleteDisk, isDeleting, globalMappings,
    saveMapping, handleRescan, confirmDelete
  }
}
