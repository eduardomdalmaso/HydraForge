<script setup lang="ts">
import DatasetImportModal from './DatasetImportModal.vue'
import DatasetAnnotateModal from './DatasetAnnotateModal.vue'
import DatasetMergerModal from './DatasetMergerModal.vue'
import CyberAlertModal from '../CyberAlertModal.vue'
import type { DatasetInfo } from '../../types/dataset'

defineProps<{
  isModalOpen: boolean
  isAnnotateOpen: boolean
  isMergeOpen: boolean
  deleteTarget: any
  deleteDisk: boolean
  isDeleting: boolean
  selectedDataset: DatasetInfo | null
  datasets: DatasetInfo[]
  globalMappings: Record<string, any>
}>()

const emit = defineEmits<{
  (e: 'update:isModalOpen', val: boolean): void
  (e: 'update:isAnnotateOpen', val: boolean): void
  (e: 'update:isMergeOpen', val: boolean): void
  (e: 'update:deleteDisk', val: boolean): void
  (e: 'cancelDelete'): void
  (e: 'confirmDelete'): void
  (e: 'datasetImported', newDs: any): void
}>()
</script>

<template>
  <DatasetImportModal :isOpen="isModalOpen" @close="emit('update:isModalOpen', false)" @datasetImported="(ds) => emit('datasetImported', ds)" />
  <DatasetAnnotateModal :isOpen="isAnnotateOpen" :dataset="selectedDataset" @close="emit('update:isAnnotateOpen', false)" />
  <DatasetMergerModal :isOpen="isMergeOpen" :datasets="datasets" :savedMappings="globalMappings" @close="emit('update:isMergeOpen', false)" @datasetsMerged="(ds) => emit('datasetImported', ds)" />
  <CyberAlertModal
    :isOpen="Boolean(deleteTarget)"
    title="DELETE DATASET"
    type="danger"
    :message="`Delete dataset '${deleteTarget?.name}' (${deleteTarget?.id || deleteTarget?.dataset_id})?`"
    confirmText="CONFIRM DELETE"
    cancelText="ABORT"
    checkboxLabel="Also purge physical files from disk"
    :checkboxChecked="deleteDisk"
    :isProcessing="isDeleting"
    @update:checkboxChecked="(val: boolean) => emit('update:deleteDisk', val)"
    @confirm="emit('confirmDelete')"
    @cancel="emit('cancelDelete')"
  />
</template>
