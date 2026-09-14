<script setup lang="ts">
import { ref, toRef } from 'vue'
import DatasetFolderExplorer from '../components/datasets/DatasetFolderExplorer.vue'
import DatasetPipelineModal from '../components/datasets/DatasetPipelineModal.vue'
import DatasetClassInspectorModal from '../components/datasets/DatasetClassInspectorModal.vue'
import DatasetModalsContainer from '../components/datasets/DatasetModalsContainer.vue'
import { useDatasetStudio } from '../components/datasets/useDatasetStudio'
import type { DatasetInfo } from '../types/dataset'

const props = defineProps<{ datasets?: DatasetInfo[] }>()
const {
  datasets, selectedDataset, isModalOpen, isAnnotateOpen, isMergeOpen,
  deleteTarget, deleteDisk, isDeleting, isRescanning, globalMappings,
  handleRescan, confirmDelete
} = useDatasetStudio(toRef(props, 'datasets'))

const isPipelineOpen = ref(false)
const isInspectorOpen = ref(false)
const inspectingDataset = ref<any>(null)

const handleInspect = (ds: any) => {
  inspectingDataset.value = ds
  selectedDataset.value = ds
  isInspectorOpen.value = true
}

const navigateTo = (hash: string) => {
  if (typeof window !== 'undefined') window.location.hash = hash
}

const handleDatasetImported = (newDs: any) => {
  datasets.value = [newDs, ...datasets.value]
  selectedDataset.value = newDs
}
</script>

<template>
  <div class="view-container datasets-container">
    <div class="cockpit-full-header">
      <h1 class="cockpit-main-title">REPOSITORIO DE DATASETS YOLO</h1>
      <p class="cockpit-main-subtitle">EXPLORADOR DE PASTAS // INSPECAO VISUAL DE CLASSES COM FOTOS // PIPELINE DE FUSAO EM 4 ETAPAS</p>
    </div>

    <!-- WINDOWS EXPLORER FOLDER GRID -->
    <DatasetFolderExplorer
      :datasets="datasets"
      :selectedId="selectedDataset?.id || selectedDataset?.dataset_id"
      :isRescanning="isRescanning"
      @selectDataset="handleInspect"
      @inspectClasses="handleInspect"
      @openPipelineModal="isPipelineOpen = true"
      @openImportModal="isModalOpen = true"
      @openAnnotateModal="isAnnotateOpen = true"
      @rescan="handleRescan"
      @deleteDataset="(ds) => deleteTarget = ds"
    />

    <!-- CLASS INSPECTOR MODAL WITH PHOTOS & BBOXES -->
    <DatasetClassInspectorModal
      :isOpen="isInspectorOpen"
      :dataset="inspectingDataset"
      @close="isInspectorOpen = false"
      @train="() => navigateTo('cockpit')"
    />

    <!-- 4-STEP PIPELINE CONFIGURATION MODAL -->
    <DatasetPipelineModal
      :isOpen="isPipelineOpen"
      :datasets="datasets"
      :savedMappings="globalMappings"
      @close="isPipelineOpen = false"
      @pipelineComplete="handleRescan"
    />

    <!-- OTHER MODALS (IMPORT / DELETE / ANNOTATE) -->
    <DatasetModalsContainer
      v-model:isModalOpen="isModalOpen"
      v-model:isAnnotateOpen="isAnnotateOpen"
      v-model:isMergeOpen="isMergeOpen"
      v-model:deleteDisk="deleteDisk"
      :deleteTarget="deleteTarget"
      :isDeleting="isDeleting"
      :selectedDataset="selectedDataset"
      :datasets="datasets"
      :globalMappings="globalMappings"
      @confirmDelete="confirmDelete"
      @cancelDelete="deleteTarget = null"
      @datasetImported="handleDatasetImported"
    />
  </div>
</template>
