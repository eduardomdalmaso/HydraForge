<script setup lang="ts">
import { toRef } from 'vue'
import DatasetListCard from '../components/datasets/DatasetListCard.vue'
import DatasetYamlViewer from '../components/datasets/DatasetYamlViewer.vue'
import DatasetKFoldCard from '../components/datasets/DatasetKFoldCard.vue'
import DatasetAuditCard from '../components/datasets/DatasetAuditCard.vue'
import DatasetModalsContainer from '../components/datasets/DatasetModalsContainer.vue'
import { useDatasetStudio } from '../components/datasets/useDatasetStudio'
import type { DatasetInfo } from '../types/dataset'

const props = defineProps<{ datasets?: DatasetInfo[] }>()
const {
  datasets, selectedDataset, isModalOpen, isAnnotateOpen, isMergeOpen,
  deleteTarget, deleteDisk, isDeleting, globalMappings,
  saveMapping, handleRescan, confirmDelete
} = useDatasetStudio(toRef(props, 'datasets'))

const navigate = (hash: string) => {
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
      <h1 class="cockpit-main-title">DATASET STUDIO</h1>
      <p class="cockpit-main-subtitle">DATA REPOSITORY // YAML SPECIFICATION // K-FOLD CROSS-VALIDATION & LEAKAGE AUDIT</p>
    </div>

    <div class="datasets-grid">
      <div style="display: flex; flex-direction: column; gap: 1.25rem;">
        <DatasetListCard
          :datasets="datasets"
          :selectedId="selectedDataset?.id || selectedDataset?.dataset_id"
          @selectDataset="(ds) => selectedDataset = ds"
          @openImportModal="isModalOpen = true"
          @openAnnotateModal="isAnnotateOpen = true"
          @openMergeModal="isMergeOpen = true"
          @rescan="handleRescan"
          @deleteDataset="(ds) => deleteTarget = ds"
        />
        <DatasetYamlViewer :dataset="selectedDataset" @saveMappings="saveMapping" />
      </div>

      <div style="display: flex; flex-direction: column; gap: 1.25rem;">
        <DatasetKFoldCard :dataset="selectedDataset" />
        <DatasetAuditCard :dataset="selectedDataset" @sendToCockpit="navigate('cockpit')" />
      </div>
    </div>

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
