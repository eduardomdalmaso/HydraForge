<script setup lang="ts">
import { ref } from 'vue'
import { sendToHydraVaultInbox } from '../../api/media_client'

const props = defineProps<{
  entity: any | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const statusMsg = ref<string | null>(null)
const isSending = ref(false)

const handleExportPng = () => {
  if (!props.entity?.thumbnailUrl) return
  const a = document.createElement('a')
  a.href = props.entity.thumbnailUrl
  a.download = `detection_${(props.entity.label || 'target').toLowerCase()}_${Date.now()}.jpg`
  a.click()
  statusMsg.value = '[OK] FOTO EXPORTADA'
  setTimeout(() => statusMsg.value = null, 2500)
}

const handleSendToVault = async () => {
  if (!props.entity?.thumbnailUrl) return
  isSending.value = true
  const datasetTag = props.entity.sourceCategory && props.entity.sourceCategory !== 'root' ? props.entity.sourceCategory : 'playground_curation'
  const cameraTag = props.entity.sourceFile ? `video_${props.entity.sourceCategory}_${props.entity.sourceFile}` : 'video_loop'

  const res = await sendToHydraVaultInbox(datasetTag, cameraTag, props.entity.thumbnailUrl, [
    { label: props.entity.label, confidence: props.entity.conf, box: props.entity.box }
  ])

  isSending.value = false
  if (res.success) {
    props.entity.vaultStatus = 'synced'
    statusMsg.value = `[HYDRAVAULT :8082] SINCRONIZADO // ${datasetTag.toUpperCase()}`
  } else {
    props.entity.vaultStatus = 'error'
    statusMsg.value = `[HYDRAVAULT] ${res.message}`
  }
  setTimeout(() => statusMsg.value = null, 3500)
}
</script>

<template>
  <div v-if="entity" class="cyber-card detection-inspector-card">
    <div class="card-header" style="padding-bottom: 0.5rem; border-bottom: 1px solid var(--vms-border);">
      <div style="display: flex; align-items: center; gap: 0.5rem;">
        <span class="badge-cyan">[TARGET INSPECTOR]</span>
        <span class="card-title" style="font-size: 0.82rem;">{{ entity.label || 'OBJETO' }}</span>
      </div>
      <button class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.65rem;" @click="emit('close')">[✕]</button>
    </div>

    <!-- LARGE SNAPSHOT PREVIEW -->
    <div class="inspector-media-container">
      <img v-if="entity.thumbnailUrl" :src="entity.thumbnailUrl" :alt="entity.label" class="inspector-media-img" />
      <div v-else class="inspector-media-fallback">SEM FOTO // {{ entity.label }}</div>
      <div class="inspector-reticle-overlay"></div>
      <div class="inspector-conf-badge">{{ (entity.conf * 100).toFixed(1) }}%</div>
    </div>

    <!-- METADATA GRID -->
    <div class="inspector-meta-grid">
      <div class="meta-item">
        <span class="meta-k">TRACK ID:</span>
        <span class="meta-v text-mono">{{ entity.track_id ? `#${entity.track_id}` : (entity.id || 'N/A') }}</span>
      </div>
      <div class="meta-item">
        <span class="meta-k">TIMESTAMP:</span>
        <span class="meta-v text-mono">{{ entity.timestamp || 'REALTIME' }}</span>
      </div>
      <div class="meta-item" style="grid-column: 1 / -1;">
        <span class="meta-k">BOUNDING BOX:</span>
        <span class="meta-v text-mono" style="font-size: 0.68rem; color: var(--cb-cyan);">
          {{ entity.box ? `X:${entity.box[0]}% Y:${entity.box[1]}% W:${entity.box[2]}% H:${entity.box[3]}%` : 'N/A' }}
        </span>
      </div>
      <div v-if="entity.sourceCategory" class="meta-item" style="grid-column: 1 / -1;">
        <span class="meta-k">FONTE / ORIGEM:</span>
        <span class="meta-v text-mono" style="font-size: 0.68rem; color: #cbd5e1;">{{ entity.sourceCategory }}</span>
      </div>
    </div>

    <div v-if="statusMsg" style="background: rgba(0,255,157,0.08); border: 1px solid var(--cb-green); padding: 0.35rem 0.5rem; border-radius: 3px; font-size: 0.7rem; color: var(--cb-green); font-family: var(--font-mono); margin-top: 0.5rem;">
      {{ statusMsg }}
    </div>

    <!-- ACTIONS -->
    <div style="display: flex; gap: 0.4rem; margin-top: 0.6rem;">
      <button class="cyber-action-btn secondary" style="flex: 1; padding: 0.4rem; font-size: 0.72rem;" @click="handleExportPng">
        <span>BAIXAR PNG</span>
      </button>
      <button
        class="cyber-action-btn"
        style="flex: 1.4; padding: 0.4rem; font-size: 0.72rem;"
        :disabled="isSending || entity.vaultStatus === 'synced'"
        @click="handleSendToVault"
      >
        <span>{{ entity.vaultStatus === 'synced' ? '[VAULT: OK]' : (isSending ? 'ENVIANDO...' : '↑ HYDRAVAULT') }}</span>
      </button>
    </div>
  </div>
</template>
