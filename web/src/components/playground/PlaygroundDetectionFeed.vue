<script setup lang="ts">
import type { SentFrameLog } from './useVaultCollector'

defineProps<{
  sentLogs?: SentFrameLog[]
  totalCollected?: number
  isSending?: boolean
  selectedEntity?: any
}>()

const emit = defineEmits<{
  (e: 'selectEntity', det: any): void
  (e: 'clearHistory'): void
}>()
</script>

<template>
  <div class="cyber-card playground-feed-card">
    <div class="card-header">
      <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
        <span class="card-title">STREAM DIRETO HYDRAVAULT // CURADORIA ATIVA (:8082)</span>
        <span class="badge-cyan">{{ totalCollected || 0 }} COLETADOS</span>
        <span v-if="isSending" style="font-size: 0.65rem; color: var(--cb-green); font-family: var(--font-mono);">● ENVIANDO FRAME...</span>
      </div>
      <div style="display: flex; gap: 0.4rem; align-items: center;">
        <button v-if="(sentLogs?.length || 0) > 0" class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.65rem;" @click="emit('clearHistory')">[LIMPAR FEED]</button>
      </div>
    </div>

    <div class="feed-photo-grid">
      <div v-if="!sentLogs || sentLogs.length === 0" class="feed-empty-state">
        <div style="font-family: var(--font-mono); font-size: 0.75rem; color: #64748b; margin-bottom: 0.25rem;">
          // AGUARDANDO DETECÇÕES DAS CLASSES SELECIONADAS
        </div>
        <div style="font-size: 0.68rem; color: #475569;">
          Inicie o Stream no painel lateral. As fotos das classes selecionadas serão gravadas diretamente no HydraVault sem sobrecarregar a memória.
        </div>
      </div>

      <div
        v-for="log in (sentLogs || [])"
        :key="log.id"
        class="feed-photo-card"
        :class="{ active: selectedEntity?.id === log.id }"
        :title="`${log.label} enviado às ${log.timestamp}`"
        @click="emit('selectEntity', log)"
      >
        <div class="photo-card-wrapper">
          <img v-if="log.thumbnailUrl" :src="log.thumbnailUrl" :alt="log.label" class="photo-card-img" />
          <div v-else class="photo-card-fallback">{{ log.label.slice(0, 4) }}</div>
          <div class="photo-card-badge">{{ (log.conf * 100).toFixed(0) }}%</div>
          <div class="photo-card-footer">
            <span class="photo-card-label">{{ log.label }}</span>
            <span style="font-size: 0.58rem; color: var(--cb-green); font-family: var(--font-mono);">[VAULT]</span>
          </div>
          <div v-if="selectedEntity?.id === log.id" class="photo-card-reticle"></div>
        </div>
      </div>
    </div>
  </div>
</template>
