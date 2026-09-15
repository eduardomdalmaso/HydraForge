<script setup lang="ts">
import { ref } from 'vue'
import type { DetectionEvent } from './useDetectionHistory'
import { sendToHydraVaultInbox } from '../../api/media_client'

const props = defineProps<{
  events: DetectionEvent[]
  selectedEntity?: any
}>()

const emit = defineEmits<{
  (e: 'selectEntity', det: any): void
  (e: 'clearHistory'): void
}>()

const exportedMsg = ref<string | null>(null)
const isSyncingAll = ref(false)

const handleExportPng = (e: MouseEvent, evt: DetectionEvent) => {
  e.stopPropagation()
  if (evt.thumbnailUrl) {
    const a = document.createElement('a')
    a.href = evt.thumbnailUrl
    a.download = `detection_${evt.label.toLowerCase()}_${Date.now()}.jpg`
    a.click()
  }
  exportedMsg.value = `EXPORTED // ${evt.label} #${evt.id}`
  setTimeout(() => exportedMsg.value = null, 2500)
}

const handleSendToVault = async (e: MouseEvent, evt: DetectionEvent) => {
  e.stopPropagation()
  if (!evt.thumbnailUrl) return
  evt.vaultStatus = 'sending'
  const datasetTag = evt.sourceCategory && evt.sourceCategory !== 'root' ? evt.sourceCategory : 'playground_curation'
  const cameraTag = evt.sourceFile ? `video_${evt.sourceCategory}_${evt.sourceFile}` : 'video_loop'

  const res = await sendToHydraVaultInbox(datasetTag, cameraTag, evt.thumbnailUrl, [
    { label: evt.label, confidence: evt.conf, box: evt.box }
  ])

  if (res.success) {
    evt.vaultStatus = 'synced'
    exportedMsg.value = `[HYDRAVAULT :8082] ENVIADO // DATASET: ${datasetTag.toUpperCase()}`
  } else {
    evt.vaultStatus = 'error'
    exportedMsg.value = `[HYDRAVAULT] ${res.message}`
  }
  setTimeout(() => exportedMsg.value = null, 3500)
}

const handleSyncAll = async () => {
  if (props.events.length === 0) return
  isSyncingAll.value = true
  let successCount = 0

  for (const evt of props.events) {
    if (evt.thumbnailUrl && evt.vaultStatus !== 'synced') {
      const datasetTag = evt.sourceCategory && evt.sourceCategory !== 'root' ? evt.sourceCategory : 'playground_curation'
      const cameraTag = evt.sourceFile ? `video_${evt.sourceCategory}_${evt.sourceFile}` : 'video_loop'
      const res = await sendToHydraVaultInbox(datasetTag, cameraTag, evt.thumbnailUrl, [{ label: evt.label, confidence: evt.conf, box: evt.box }])
      if (res.success) {
        evt.vaultStatus = 'synced'
        successCount++
      }
    }
  }

  isSyncingAll.value = false
  exportedMsg.value = `[HYDRAVAULT] SINCRONIZADOS ${successCount} FRAMES COM SUCESSO!`
  setTimeout(() => exportedMsg.value = null, 4000)
}
</script>

<template>
  <div class="cyber-card playground-feed-card">
    <div class="card-header">
      <div style="display: flex; align-items: center; gap: 0.5rem;">
        <span class="card-title">HISTORICO DE DETECCOES // REALTIME</span>
        <span class="badge-cyan">{{ events.length }}</span>
      </div>
      <div style="display: flex; gap: 0.4rem; align-items: center;">
        <button
          v-if="events.length > 0"
          class="cyber-pill"
          style="padding: 0.15rem 0.5rem; font-size: 0.65rem; color: var(--cb-cyan); border-color: rgba(0,240,255,0.4);"
          :disabled="isSyncingAll"
          @click="handleSyncAll"
        >
          {{ isSyncingAll ? 'ENVIANDO...' : '↑ ENVIAR TODOS AO HYDRAVAULT (:8082)' }}
        </button>
        <button v-if="events.length > 0" class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.65rem;" @click="emit('clearHistory')">
          [CLEAR]
        </button>
      </div>
    </div>

    <div v-if="exportedMsg" style="background: rgba(0,255,157,0.08); border: 1px solid var(--cb-green); padding: 0.35rem 0.6rem; border-radius: 3px; font-size: 0.72rem; color: var(--cb-green); font-family: var(--font-mono); margin-bottom: 0.5rem;">
      {{ exportedMsg }}
    </div>

    <div class="feed-list-scroll">
      <div v-if="events.length === 0" class="feed-empty-state">
        <div style="font-family: var(--font-mono); font-size: 0.75rem; color: #64748b; margin-bottom: 0.25rem;">
          // NENHUMA DETECCAO NO BUFFER
        </div>
        <div style="font-size: 0.68rem; color: #475569;">
          Dispare um Scan ou ative o Live Scan para registrar alvos em tempo real.
        </div>
      </div>

      <div
        v-for="evt in events"
        v-else
        :key="evt.id"
        class="feed-item"
        :class="{ active: selectedEntity?.id === evt.id }"
        @click="emit('selectEntity', evt)"
      >
        <div class="feed-thumbnail-wrapper">
          <img v-if="evt.thumbnailUrl" :src="evt.thumbnailUrl" :alt="evt.label" class="feed-thumbnail-img" />
          <div v-else class="feed-thumbnail-fallback">{{ evt.label.slice(0, 3) }}</div>
        </div>

        <div class="feed-info-col">
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <div style="display: flex; align-items: center; gap: 0.4rem;">
              <span class="feed-entity-label">{{ evt.label }}</span>
              <span v-if="evt.sourceCategory" class="file-tag" style="font-size: 0.58rem;">
                [{{ evt.sourceCategory.toUpperCase() }}]
              </span>
            </div>
            <span class="feed-conf-tag">{{ (evt.conf * 100).toFixed(0) }}%</span>
          </div>

          <div style="display: flex; align-items: center; justify-content: space-between; margin-top: 2px;">
            <span class="feed-timestamp">{{ evt.timestamp }}</span>
            <div style="display: flex; gap: 0.25rem; align-items: center;">
              <button class="cyber-pill feed-btn" title="Baixar PNG" @click="(e) => handleExportPng(e, evt)">[PNG]</button>
              <button
                class="cyber-pill feed-btn"
                :style="{
                  color: evt.vaultStatus === 'synced' ? 'var(--cb-green)' : 'var(--cb-yellow)',
                  borderColor: evt.vaultStatus === 'synced' ? 'rgba(0,255,157,0.5)' : 'rgba(252,238,10,0.4)'
                }"
                :title="evt.vaultStatus === 'synced' ? 'Já enviado ao HydraVault' : 'Enviar snapshot para curadoria no HydraVault (:8082)'"
                :disabled="evt.vaultStatus === 'sending' || evt.vaultStatus === 'synced'"
                @click="(e) => handleSendToVault(e, evt)"
              >
                {{ evt.vaultStatus === 'synced' ? '[VAULT: OK]' : (evt.vaultStatus === 'sending' ? '[...]' : '[+ VAULT]') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
