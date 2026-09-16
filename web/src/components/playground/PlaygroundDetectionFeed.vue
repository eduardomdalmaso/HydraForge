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

    <div class="feed-photo-grid">
      <div v-if="events.length === 0" class="feed-empty-state">
        <div style="font-family: var(--font-mono); font-size: 0.75rem; color: #64748b; margin-bottom: 0.25rem;">
          // NENHUMA DETECCAO NO BUFFER
        </div>
        <div style="font-size: 0.68rem; color: #475569;">
          Dispare um Scan ou ative o Live Scan para registrar fotos de alvos em tempo real.
        </div>
      </div>

      <div
        v-for="evt in events"
        v-else
        :key="evt.id"
        class="feed-photo-card"
        :class="{ active: selectedEntity?.id === evt.id }"
        :title="`Clique para inspecionar ${evt.label} (${(evt.conf * 100).toFixed(0)}%)`"
        @click="emit('selectEntity', evt)"
      >
        <div class="photo-card-wrapper">
          <img v-if="evt.thumbnailUrl" :src="evt.thumbnailUrl" :alt="evt.label" class="photo-card-img" />
          <div v-else class="photo-card-fallback">{{ evt.label.slice(0, 4) }}</div>
          
          <!-- TOP CONFIDENCE BADGE -->
          <div class="photo-card-badge">{{ (evt.conf * 100).toFixed(0) }}%</div>

          <!-- BOTTOM LABEL OVERLAY -->
          <div class="photo-card-footer">
            <span class="photo-card-label">{{ evt.label }}</span>
            <span v-if="evt.track_id" class="photo-card-track">#{{ evt.track_id }}</span>
          </div>

          <!-- SELECTION RETICLE -->
          <div v-if="selectedEntity?.id === evt.id" class="photo-card-reticle"></div>
        </div>
      </div>
    </div>
  </div>
</template>
