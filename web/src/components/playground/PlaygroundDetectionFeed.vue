<script setup lang="ts">
import { ref, computed } from 'vue'
import type { DetectionEvent } from './useDetectionHistory'
import { sendToHydraVaultInbox } from '../../api/media_client'

const props = defineProps<{ events: DetectionEvent[]; selectedEntity?: any }>()
const emit = defineEmits<{ (e: 'selectEntity', det: any): void; (e: 'clearHistory'): void }>()

const exportedMsg = ref<string | null>(null)
const isSyncingAll = ref(false)
const selectedClasses = ref<Record<string, boolean>>({})
const maxPerClass = ref<number>(0)
const onlyHardCases = ref<boolean>(false)
const showFilters = ref(true)

const availableClasses = computed(() => {
  const counts: Record<string, number> = {}
  for (const e of props.events) counts[e.label] = (counts[e.label] || 0) + 1
  return counts
})

const isClassActive = (label: string) => selectedClasses.value[label] !== false
const toggleClass = (label: string) => { selectedClasses.value[label] = selectedClasses.value[label] === false }

const filteredEventsToSend = computed(() => {
  const counts: Record<string, number> = {}
  const result: DetectionEvent[] = []
  for (const evt of props.events) {
    if (selectedClasses.value[evt.label] === false) continue
    if (onlyHardCases.value && evt.conf > 0.65) continue
    const current = counts[evt.label] || 0
    if (maxPerClass.value > 0 && current >= maxPerClass.value) continue
    counts[evt.label] = current + 1
    result.push(evt)
  }
  return result
})

const handleSyncSelected = async () => {
  const toSend = filteredEventsToSend.value
  if (toSend.length === 0) return
  isSyncingAll.value = true
  let successCount = 0
  let lastError = ''

  for (const evt of toSend) {
    if (evt.vaultStatus === 'synced') continue
    evt.vaultStatus = 'sending'
    const bx = evt.box || [10, 10, 30, 30]
    const bboxPayload = [{
      class_id: 0,
      class_name: evt.label.toLowerCase(),
      x_center: Math.min(0.99, Math.max(0.01, (bx[0] + bx[2] / 2) / 100)),
      y_center: Math.min(0.99, Math.max(0.01, (bx[1] + bx[3] / 2) / 100)),
      width: Math.min(0.99, Math.max(0.02, bx[2] / 100)),
      height: Math.min(0.99, Math.max(0.02, bx[3] / 100)),
      confidence: evt.conf
    }]
    const datasetTag = evt.sourceCategory && evt.sourceCategory !== 'root' ? evt.sourceCategory : 'playground_curation'
    const cameraTag = evt.sourceFile ? `video_${evt.sourceCategory}_${evt.sourceFile}` : 'video_loop'
    const res = await sendToHydraVaultInbox(datasetTag, cameraTag, evt.thumbnailUrl || '', bboxPayload)
    if (res.success) {
      evt.vaultStatus = 'synced'
      successCount++
    } else {
      evt.vaultStatus = 'error'
      lastError = res.message
    }
  }

  isSyncingAll.value = false
  if (successCount > 0) {
    exportedMsg.value = `[HYDRAVAULT] ${successCount} FRAMES ENVIADOS COM SUCESSO AO INBOX (:8082)!`
  } else if (lastError) {
    exportedMsg.value = `[HYDRAVAULT ERRO] ${lastError}`
  } else {
    exportedMsg.value = `[HYDRAVAULT] NENHUM NOVO FRAME PENDENTE`
  }
  setTimeout(() => exportedMsg.value = null, 5000)
}
</script>

<template>
  <div class="cyber-card playground-feed-card">
    <div class="card-header">
      <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
        <span class="card-title">HISTORICO DE DETECCOES // REALTIME</span>
        <span class="badge-cyan">{{ events.length }}</span>
        <button v-if="events.length > 0" class="cyber-pill" style="padding: 0.15rem 0.45rem; font-size: 0.65rem;" :class="{ active: showFilters }" @click="showFilters = !showFilters">[FILTROS & COTAS]</button>
      </div>
      <div style="display: flex; gap: 0.4rem; align-items: center; flex-wrap: wrap;">
        <button v-if="events.length > 0" class="cyber-pill" style="padding: 0.2rem 0.55rem; font-size: 0.68rem; color: var(--cb-cyan); border-color: rgba(0,240,255,0.4);" :disabled="isSyncingAll || filteredEventsToSend.length === 0" @click="handleSyncSelected">
          {{ isSyncingAll ? 'ENVIANDO...' : `↑ ENVIAR AO VAULT (${filteredEventsToSend.length})` }}
        </button>
        <button v-if="events.length > 0" class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.65rem;" @click="emit('clearHistory')">[CLEAR]</button>
      </div>
    </div>

    <!-- BALANCE AND ACTIVE LEARNING CONTROL BAR -->
    <div v-if="showFilters && events.length > 0" style="background: var(--vms-bg-elevated); border: 1px solid var(--vms-border); border-radius: 4px; padding: 0.5rem 0.65rem; margin-bottom: 0.6rem; display: flex; flex-direction: column; gap: 0.45rem;">
      <div style="display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.4rem;">
        <div style="font-size: 0.68rem; font-family: var(--font-mono); color: var(--vms-text-muted);">CLASSES ATIVAS:</div>
        <div style="display: flex; gap: 0.3rem; flex-wrap: wrap;">
          <button v-for="(cnt, lbl) in availableClasses" :key="lbl" class="cyber-pill" :class="{ active: isClassActive(lbl) }" style="font-size: 0.63rem; padding: 0.12rem 0.4rem;" @click="toggleClass(lbl)">
            {{ isClassActive(lbl) ? '[X]' : '[ ]' }} {{ lbl }} ({{ cnt }})
          </button>
        </div>
      </div>
      <div style="display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.4rem; border-top: 1px solid rgba(255,255,255,0.06); padding-top: 0.4rem;">
        <div style="display: flex; align-items: center; gap: 0.3rem; flex-wrap: wrap;">
          <span style="font-size: 0.68rem; font-family: var(--font-mono); color: var(--vms-text-muted);">COTA/CLASSE:</span>
          <button class="cyber-pill" :class="{ active: maxPerClass === 0 }" style="font-size: 0.62rem; padding: 0.1rem 0.35rem;" @click="maxPerClass = 0">[TODOS]</button>
          <button class="cyber-pill" :class="{ active: maxPerClass === 3 }" style="font-size: 0.62rem; padding: 0.1rem 0.35rem;" @click="maxPerClass = 3">[MAX 3]</button>
          <button class="cyber-pill" :class="{ active: maxPerClass === 5 }" style="font-size: 0.62rem; padding: 0.1rem 0.35rem;" @click="maxPerClass = 5">[MAX 5]</button>
          <button class="cyber-pill" :class="{ active: maxPerClass === 10 }" style="font-size: 0.62rem; padding: 0.1rem 0.35rem;" @click="maxPerClass = 10">[MAX 10]</button>
        </div>
        <div>
          <button class="cyber-pill" :class="{ active: onlyHardCases }" style="font-size: 0.62rem; padding: 0.1rem 0.4rem;" @click="onlyHardCases = !onlyHardCases">
            {{ onlyHardCases ? '★ HARD CASES (<65%) ATIVO' : '☆ CASOS DIFICEIS (<65%)' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="exportedMsg" style="background: rgba(0,255,157,0.08); border: 1px solid var(--cb-green); padding: 0.35rem 0.6rem; border-radius: 3px; font-size: 0.72rem; color: var(--cb-green); font-family: var(--font-mono); margin-bottom: 0.5rem;">
      {{ exportedMsg }}
    </div>

    <div class="feed-photo-grid">
      <div v-if="events.length === 0" class="feed-empty-state">
        <div style="font-family: var(--font-mono); font-size: 0.75rem; color: #64748b; margin-bottom: 0.25rem;">// NENHUMA DETECCAO NO BUFFER</div>
        <div style="font-size: 0.68rem; color: #475569;">Dispare um Scan ou ative o Live Scan para registrar fotos de alvos em tempo real.</div>
      </div>
      <div v-for="evt in events" v-else :key="evt.id" class="feed-photo-card" :class="{ active: selectedEntity?.id === evt.id }" :title="`Clique para inspecionar ${evt.label} (${(evt.conf * 100).toFixed(0)}%)`" @click="emit('selectEntity', evt)">
        <div class="photo-card-wrapper">
          <img v-if="evt.thumbnailUrl" :src="evt.thumbnailUrl" :alt="evt.label" class="photo-card-img" />
          <div v-else class="photo-card-fallback">{{ evt.label.slice(0, 4) }}</div>
          <div class="photo-card-badge">{{ (evt.conf * 100).toFixed(0) }}%</div>
          <div class="photo-card-footer">
            <span class="photo-card-label">{{ evt.label }}</span>
            <span v-if="evt.track_id" class="photo-card-track">#{{ evt.track_id }}</span>
          </div>
          <div v-if="selectedEntity?.id === evt.id" class="photo-card-reticle"></div>
        </div>
      </div>
    </div>
  </div>
</template>
