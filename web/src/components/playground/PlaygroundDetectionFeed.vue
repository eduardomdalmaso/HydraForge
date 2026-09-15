<script setup lang="ts">
import { ref } from 'vue'
import type { DetectionEvent } from './useDetectionHistory'

defineProps<{
  events: DetectionEvent[]
  selectedEntity?: any
}>()

const emit = defineEmits<{
  (e: 'selectEntity', det: any): void
  (e: 'clearHistory'): void
}>()

const exportedMsg = ref<string | null>(null)

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
</script>

<template>
  <div class="cyber-card playground-feed-card">
    <div class="card-header">
      <div style="display: flex; align-items: center; gap: 0.5rem;">
        <span class="card-title">HISTORICO DE DETECCOES // REALTIME</span>
        <span class="badge-cyan">{{ events.length }}</span>
      </div>
      <button v-if="events.length > 0" class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.65rem;" @click="emit('clearHistory')">
        [CLEAR]
      </button>
    </div>

    <div v-if="exportedMsg" style="background: rgba(0,255,157,0.08); border: 1px solid var(--cb-green); padding: 0.35rem; border-radius: 3px; font-size: 0.7rem; color: var(--cb-green); font-family: var(--font-mono); margin-bottom: 0.5rem;">
      [OK] {{ exportedMsg }}
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
          <img
            v-if="evt.thumbnailUrl"
            :src="evt.thumbnailUrl"
            :alt="evt.label"
            class="feed-thumbnail-img"
          />
          <div v-else class="feed-thumbnail-fallback">
            {{ evt.label.slice(0, 3) }}
          </div>
        </div>

        <div class="feed-info-col">
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <span class="feed-entity-label">{{ evt.label }}</span>
            <span class="feed-conf-tag">{{ (evt.conf * 100).toFixed(0) }}%</span>
          </div>
          <div style="display: flex; align-items: center; justify-content: space-between; margin-top: 2px;">
            <span class="feed-timestamp">{{ evt.timestamp }}</span>
            <div style="display: flex; gap: 0.25rem;">
              <button class="cyber-pill feed-btn" title="Exportar snapshot" @click="(e) => handleExportPng(e, evt)">[PNG]</button>
              <a href="http://localhost:8082" target="_blank" rel="noreferrer" class="cyber-pill feed-btn feed-vault-btn" title="Curadoria no HydraVault" @click.stop>[VAULT]</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
