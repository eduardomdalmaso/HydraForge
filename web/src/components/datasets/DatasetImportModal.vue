<script setup lang="ts">
import { ref } from 'vue'
import DatasetCyberIcons from './DatasetCyberIcons.vue'
import { registerDatasetPathAPI } from '../../api/client'

defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'datasetImported', ds: any): void }>()

const activeTab = ref<'zip' | 'path'>('zip')
const zipFiles = ref<File[]>([])
const folderPath = ref('')
const isProcessing = ref(false)
const error = ref<string | null>(null)

const handleExecuteImport = async () => {
  isProcessing.value = true
  error.value = null
  try {
    if (activeTab.value === 'zip' && zipFiles.value.length > 0) {
      for (const file of zipFiles.value) {
        const formData = new FormData()
        formData.append('file', file)
        const cleanId = file.name.replace(/\.zip$/i, '').toLowerCase().replace(/[^a-z0-9_]/g, '_')
        formData.append('dataset_id', cleanId)
        formData.append('task', 'detect')
        const res = await fetch('/api/v1/training/datasets/import', { method: 'POST', body: formData })
        if (res.ok) emit('datasetImported', await res.json())
      }
    } else if (activeTab.value === 'path' && folderPath.value.trim()) {
      const registered = await registerDatasetPathAPI({ name: folderPath.value.trim().split('/').pop() || 'dataset', yaml_path: folderPath.value.trim() })
      emit('datasetImported', registered)
    }
    emit('close')
  } catch (err: any) {
    error.value = err.message || 'Import failed'
  } finally {
    isProcessing.value = false
  }
}
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click="emit('close')">
    <div class="import-modal-box" style="max-width: 600px;" @click.stop>
      <div class="card-header" style="margin-bottom: 0.6rem;">
        <span class="card-title">IMPORT DATASET (MULTI-ZIP / FOLDER)</span>
        <button class="cyber-pill" @click="emit('close')">✕</button>
      </div>

      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.4rem; margin-bottom: 0.6rem;">
        <button class="cyber-pill" :class="{ active: activeTab === 'zip' }" style="width: 100%; justify-content: center;" @click="activeTab = 'zip'">[ZIP] MULTIPLE .ZIP FILES</button>
        <button class="cyber-pill" :class="{ active: activeTab === 'path' }" style="width: 100%; justify-content: center;" @click="activeTab = 'path'">[PATH] LOCAL FOLDER PATH</button>
      </div>

      <label v-if="activeTab === 'zip'" class="dropzone-box" style="cursor: pointer; display: block; margin: 0.5rem 0;" @dragover.prevent @drop.prevent="(e) => zipFiles = Array.from(e.dataTransfer?.files || []).filter(f => f.name.endsWith('.zip'))">
        <input type="file" accept=".zip" multiple style="display: none;" @change="(e) => zipFiles = Array.from((e.target as HTMLInputElement).files || [])" />
        <div style="text-align: center;">
          <DatasetCyberIcons name="zip" :size="28" :color="zipFiles.length ? 'var(--cb-yellow)' : 'var(--cb-cyan)'" />
          <div style="font-family: var(--font-oxanium); font-weight: 700; margin-top: 0.2rem;" :style="{ color: zipFiles.length ? 'var(--cb-yellow)' : 'var(--cb-cyan)' }">
            {{ zipFiles.length ? `[OK] ${zipFiles.length} ZIP FILE(S) SELECTED` : 'CLICK OR DRAG MULTIPLE .ZIP ARCHIVES' }}
          </div>
        </div>
      </label>

      <div v-else style="margin: 0.6rem 0;">
        <div style="font-size: 0.72rem; color: #94a3b8; margin-bottom: 0.3rem;">FILESYSTEM PATH (e.g. /home/hades/Downloads/cell-phone):</div>
        <input class="cyber-input" type="text" placeholder="/home/hades/Downloads/cell-phone" :value="folderPath" style="width: 100%; padding: 0.5rem; font-family: var(--font-mono); font-size: 0.8rem;" @input="(e) => folderPath = (e.target as HTMLInputElement).value" />
      </div>

      <div v-if="error" style="color: var(--cb-magenta); font-size: 0.75rem; margin-top: 0.3rem;">[ALERT] {{ error }}</div>

      <div style="display: flex; gap: 0.75rem; justify-content: flex-end; margin-top: 0.75rem;">
        <button class="cyber-action-btn secondary" :disabled="isProcessing" @click="emit('close')">CANCEL</button>
        <button class="cyber-action-btn" :disabled="isProcessing || (activeTab === 'zip' ? zipFiles.length === 0 : !folderPath.trim())" @click="handleExecuteImport">
          <span>{{ isProcessing ? 'IMPORTING...' : 'CONFIRM IMPORT' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>
