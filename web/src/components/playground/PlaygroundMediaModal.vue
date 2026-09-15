<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  fetchMediaSources,
  createMediaFolder,
  deleteMediaFolder,
  uploadMediaVideo,
  deleteMediaFile,
  type MediaFolder
} from '../../api/media_client'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'mediaUpdated'): void
}>()

const folders = ref<MediaFolder[]>([])
const newFolderName = ref('')
const selectedUploadFolder = ref('root')
const isUploading = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)
const statusMsg = ref('')

const loadSources = async () => {
  const data = await fetchMediaSources()
  folders.value = data.folders || []
}

onMounted(loadSources)

const handleCreateFolder = async () => {
  if (!newFolderName.value.trim()) return
  const ok = await createMediaFolder(newFolderName.value.trim())
  if (ok) {
    statusMsg.value = `PASTA [${newFolderName.value.toUpperCase()}] CRIADA COM SUCESSO`
    newFolderName.value = ''
    await loadSources()
    emit('mediaUpdated')
  } else {
    statusMsg.value = 'ERRO AO CRIAR PASTA (USE ALFANUMERICO/UNDERLINE)'
  }
}

const handleDeleteFolder = async (folder: string) => {
  if (confirm(`Deseja excluir a pasta [${folder.toUpperCase()}] e todos os seus vídeos?`)) {
    await deleteMediaFolder(folder)
    await loadSources()
    emit('mediaUpdated')
  }
}

const handleFileUpload = async (e: Event) => {
  const target = e.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return
  isUploading.value = true
  statusMsg.value = 'UPLOADING VIDEO TO STORAGE...'

  for (let i = 0; i < target.files.length; i++) {
    await uploadMediaVideo(selectedUploadFolder.value, target.files[i])
  }

  isUploading.value = false
  statusMsg.value = 'UPLOAD CONCLUIDO // ARQUIVOS PRONTOS PARA USO'
  if (fileInputRef.value) fileInputRef.value.value = ''
  await loadSources()
  emit('mediaUpdated')
}

const handleDeleteFile = async (folder: string, filename: string) => {
  if (confirm(`Deseja remover ${filename}?`)) {
    await deleteMediaFile(folder, filename)
    await loadSources()
    emit('mediaUpdated')
  }
}
</script>

<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="cyber-card media-manager-modal">
      <div class="modal-header">
        <div style="display: flex; align-items: center; gap: 0.6rem;">
          <span class="hud-tag">[STORAGE]</span>
          <h2 style="font-size: 1rem; color: var(--cb-cyan); font-family: var(--font-mono); margin: 0;">
            GERENCIADOR DE PASTAS & VIDEOS DE ENTRADA (.MP4)
          </h2>
        </div>
        <button class="modal-close-btn" @click="emit('close')">✕</button>
      </div>

      <div class="modal-body">
        <!-- TOP ACTIONS: NEW FOLDER & UPLOAD -->
        <div class="media-actions-row">
          <div class="media-action-box">
            <div class="action-box-title">// CRIAR NOVA PASTA DE CATEGORIA</div>
            <div style="display: flex; gap: 0.4rem;">
              <input
                v-model="newFolderName"
                type="text"
                placeholder="Ex: carros, noite, calcada..."
                class="cyber-input"
                style="flex: 1; font-size: 0.75rem;"
                @keyup.enter="handleCreateFolder"
              />
              <button class="cyber-action-btn" style="padding: 0 0.8rem; font-size: 0.75rem;" @click="handleCreateFolder">
                + CRIAR
              </button>
            </div>
          </div>

          <div class="media-action-box">
            <div class="action-box-title">// UPLOAD DE VIDEO (.MP4 / .MKV)</div>
            <div style="display: flex; gap: 0.4rem;">
              <select v-model="selectedUploadFolder" class="cyber-select" style="width: 140px; font-size: 0.75rem;">
                <option value="root">[GERAL // RAIZ]</option>
                <option v-for="f in folders.filter(x => x.name !== 'root')" :key="f.name" :value="f.name">
                  [{{ f.name.toUpperCase() }}]
                </option>
              </select>
              <input ref="fileInputRef" type="file" accept="video/mp4,video/mkv,video/webm" multiple style="display: none;" @change="handleFileUpload" />
              <button class="cyber-action-btn" style="flex: 1; font-size: 0.75rem;" :disabled="isUploading" @click="fileInputRef?.click()">
                {{ isUploading ? 'ENVIANDO...' : 'SELECIONAR .MP4' }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="statusMsg" class="media-status-banner">{{ statusMsg }}</div>

        <!-- FOLDERS & FILES TREE LIST -->
        <div class="media-tree-scroll">
          <div v-if="folders.length === 0" style="color: var(--cb-muted); font-size: 0.8rem; text-align: center; padding: 2rem;">
            Nenhuma pasta ou vídeo cadastrado. Adicione vídeos na pasta storage/media_sources ou faça upload acima.
          </div>

          <div v-for="folder in folders" :key="folder.name" class="media-folder-card">
            <div class="folder-header">
              <div style="display: flex; align-items: center; gap: 0.5rem;">
                <span class="folder-icon">📁</span>
                <span style="font-weight: 700; color: #fff; font-family: var(--font-mono); font-size: 0.82rem;">
                  [{{ folder.name.toUpperCase() }}]
                </span>
                <span style="font-size: 0.7rem; color: var(--cb-cyan);">({{ folder.file_count }} vídeos • {{ (folder.total_bytes / 1024 / 1024).toFixed(1) }} MB)</span>
              </div>
              <button v-if="folder.name !== 'root'" class="del-btn" title="Excluir pasta" @click="handleDeleteFolder(folder.name)">
                DELETAR PASTA
              </button>
            </div>

            <div class="folder-files-list">
              <div v-if="folder.files.length === 0" style="font-size: 0.7rem; color: var(--cb-muted); padding: 0.3rem 0.5rem;">
                (Pasta vazia)
              </div>
              <div v-for="file in folder.files" :key="file.id" class="media-file-item">
                <div style="display: flex; align-items: center; gap: 0.5rem; overflow: hidden;">
                  <span style="color: var(--cb-green); font-size: 0.75rem;">▶</span>
                  <span class="file-name-txt" :title="file.name">{{ file.name }}</span>
                  <span style="font-size: 0.65rem; color: var(--cb-muted);">({{ (file.size_bytes / 1024 / 1024).toFixed(1) }} MB)</span>
                </div>
                <div style="display: flex; gap: 0.4rem; align-items: center;">
                  <span class="file-tag">// 16:9 LOOP</span>
                  <button class="del-file-btn" title="Remover vídeo" @click="handleDeleteFile(folder.name, file.name)">✕</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
