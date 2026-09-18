<script setup lang="ts">
import { ref, computed } from 'vue'
import type { MediaFolder } from '../../api/media_client'
import type { VaultDataset } from '../../api/vault_client'

const props = withDefaults(defineProps<{
  config: Record<string, any>; modelsList?: any[]; hydraStreams?: any[]; mediaFolders?: MediaFolder[]
  isRunning?: boolean; isContinuous?: boolean; isPaused?: boolean; isVaultOnline?: boolean
  vaultDatasets?: VaultDataset[]; selectedDatasetId?: string; activeClasses?: Record<string, boolean>
  maxPerClass?: number; onlyHardCases?: boolean; hardCaseThreshold?: number; autoStreamToVault?: boolean
  collectedStats?: Record<string, number>; totalCollected?: number
}>(), {
  modelsList: () => [], hydraStreams: () => [], mediaFolders: () => [],
  isRunning: false, isContinuous: false, isPaused: false, isVaultOnline: false,
  vaultDatasets: () => [], selectedDatasetId: 'frota_urbana_fusion',
  activeClasses: () => ({}), maxPerClass: 20, onlyHardCases: false, hardCaseThreshold: 0.65,
  autoStreamToVault: true, collectedStats: () => ({}), totalCollected: 0
})

const emit = defineEmits<{
  (e: 'update:config', cfg: Record<string, any>): void; (e: 'update:isContinuous', val: boolean): void
  (e: 'update:isPaused', val: boolean): void; (e: 'update:selectedDatasetId', id: string): void
  (e: 'toggleTargetClass', label: string): void; (e: 'update:maxPerClass', limit: number): void
  (e: 'update:onlyHardCases', val: boolean): void; (e: 'update:hardCaseThreshold', val: number): void
  (e: 'update:autoStreamToVault', val: boolean): void; (e: 'resetCollectorStats'): void
  (e: 'checkVaultHealth'): void; (e: 'runInference'): void; (e: 'openMediaModal'): void
}>()

const activeTab = ref<'source' | 'classes' | 'vault_stream'>('source')
const customModels = computed(() => props.modelsList.filter(m => m.isCustom))
const baseModels = computed(() => props.modelsList.filter(m => !m.isCustom))
const updateField = (field: string, val: any) => emit('update:config', { ...props.config, [field]: val })

const currentDatasetClasses = computed(() => {
  const d = props.vaultDatasets.find(ds => ds.dataset_id === props.selectedDatasetId)
  if (d && d.classes?.length && d.classes[0].name !== 'object') return d.classes.map(c => c.name.toLowerCase())
  return ['motorcycle', 'truck', 'bus', 'car', 'person']
})
</script>

<template>
  <div class="cyber-card playground-unified-panel">
    <!-- VAULT HEALTH GATE STATUS -->
    <div style="display: flex; justify-content: space-between; align-items: center; padding: 0.45rem 0.6rem; border-bottom: 1px solid var(--vms-border); background: var(--vms-bg-elevated);">
      <div style="display: flex; align-items: center; gap: 0.4rem;">
        <span style="font-size: 0.68rem; font-family: var(--font-mono); color: var(--vms-text-muted);">VAULT GATE:</span>
        <span class="cyber-pill" :style="{ background: isVaultOnline ? 'rgba(0,255,157,0.15)' : 'rgba(255,0,60,0.2)', color: isVaultOnline ? 'var(--cb-green)' : 'var(--cb-magenta)', borderColor: isVaultOnline ? 'var(--cb-green)' : 'var(--cb-magenta)', fontSize: '0.65rem', padding: '1px 6px' }">
          {{ isVaultOnline ? '[ONLINE // :8082]' : '[OFFLINE // BLOQUEADO]' }}
        </span>
      </div>
      <button class="cyber-pill" style="font-size: 0.62rem; padding: 1px 5px;" @click="emit('checkVaultHealth')">TESTAR</button>
    </div>

    <div class="panel-tab-nav">
      <button class="panel-tab-btn" :class="{ active: activeTab === 'source' }" @click="activeTab = 'source'">1. MODELO & FONTE</button>
      <button class="panel-tab-btn" :class="{ active: activeTab === 'classes' }" @click="activeTab = 'classes'">2. DATASET & CLASSES</button>
      <button class="panel-tab-btn" :class="{ active: activeTab === 'vault_stream' }" @click="activeTab = 'vault_stream'">3. META & STREAM</button>
    </div>

    <!-- TAB 1: MODEL & SOURCE -->
    <div v-if="activeTab === 'source'" class="tab-pane">
      <div class="selector-group">
        <div class="selector-label"><span>ARQUITETURA YOLO</span><span style="font-size: 0.65rem; color: var(--cb-yellow);">{{ customModels.length }} CUSTOM</span></div>
        <select class="cyber-select" :value="config.model" @change="(e) => updateField('model', (e.target as HTMLSelectElement).value)">
          <optgroup v-if="customModels.length > 0" label="[CUSTOM TRAINED // RTX 5090]">
            <option v-for="m in customModels" :key="m.id" :value="m.id">[TREINADO] {{ m.name.toUpperCase() }} // mAP {{ m.map5095?.toFixed(1) }}%</option>
          </optgroup>
          <optgroup label="[BASE ARCHITECTURES // ULTRALYTICS]">
            <option v-for="m in baseModels" :key="m.id" :value="m.id">[BASE] {{ m.name.toUpperCase() }}</option>
          </optgroup>
        </select>
      </div>

      <div class="selector-group">
        <div class="selector-label" style="display: flex; justify-content: space-between; align-items: center;">
          <span>FONTE DE MÍDIA</span>
          <button type="button" class="del-btn" style="color: var(--cb-cyan); border-color: rgba(0,240,255,0.4);" @click="emit('openMediaModal')">+ VÍDEOS</button>
        </div>
        <select class="cyber-select" :value="config.source" @change="(e) => updateField('source', (e.target as HTMLSelectElement).value)">
          <optgroup label="[PASTAS DE VIDEOS]">
            <option v-for="folder in mediaFolders.filter(f => (f.files && f.files.length > 0) || f.file_count > 0)" :key="folder.name" :value="`folder:${folder.name}`">[PASTA] {{ folder.label }} ({{ folder.file_count }})</option>
          </optgroup>
          <optgroup label="[HARDWARE]">
            <option value="webcam">[WEBCAM] WEBCAM LOCAL</option>
          </optgroup>
          <optgroup label="[HYDRASTREAM ZERO-COPY]">
            <option v-for="s in hydraStreams" :key="s.stream_id" :value="s.stream_id">[STREAM] {{ s.stream_id.toUpperCase() }}</option>
          </optgroup>
        </select>
      </div>

      <div class="selector-group">
        <div class="selector-label"><span>CONFIANÇA MÍNIMA (CONF)</span><span class="slider-val">{{ (config.conf * 100).toFixed(0) }}%</span></div>
        <input type="range" min="0.05" max="0.95" step="0.05" :value="config.conf" @input="(e) => updateField('conf', parseFloat((e.target as HTMLInputElement).value))" />
      </div>
    </div>

    <!-- TAB 2: DATASET & CLASSES -->
    <div v-else-if="activeTab === 'classes'" class="tab-pane">
      <div class="selector-group">
        <div class="selector-label">DATASET DESTINO (HYDRAVAULT)</div>
        <select class="cyber-select" :value="selectedDatasetId" @change="(e) => emit('update:selectedDatasetId', (e.target as HTMLSelectElement).value)">
          <option v-for="ds in vaultDatasets" :key="ds.dataset_id" :value="ds.dataset_id">[DATASET] {{ ds.name }} ({{ ds.dataset_id }})</option>
          <option v-if="vaultDatasets.length === 0" value="frota_urbana_fusion">[DATASET] Frota Urbana Fusion (frota_urbana_fusion)</option>
        </select>
      </div>

      <div class="selector-group">
        <div class="selector-label">CLASSES ATIVAS PARA BUSCA & MINERAÇÃO</div>
        <div style="display: flex; gap: 0.35rem; flex-wrap: wrap; margin-top: 0.3rem;">
          <button v-for="cls in currentDatasetClasses" :key="cls" class="cyber-pill" :class="{ active: activeClasses[cls] !== false }" style="font-size: 0.65rem; padding: 0.2rem 0.5rem;" @click="emit('toggleTargetClass', cls)">
            {{ activeClasses[cls] !== false ? '[X]' : '[ ]' }} {{ cls }} ({{ collectedStats[cls] || 0 }})
          </button>
        </div>
      </div>

      <div class="selector-group" style="border-top: 1px solid rgba(255,255,255,0.06); padding-top: 0.5rem;">
        <button class="cyber-pill" :class="{ active: onlyHardCases }" style="font-size: 0.68rem; padding: 0.25rem 0.5rem; width: 100%; text-align: center; margin-bottom: 0.3rem;" @click="emit('update:onlyHardCases', !onlyHardCases)">
          {{ onlyHardCases ? '★ MINERAR APENAS CASOS DIFÍCEIS' : '☆ MINERAR TODAS AS DETECÇÕES' }}
        </button>
        <div v-if="onlyHardCases" style="margin-top: 0.3rem;">
          <div class="selector-label">
            <span>TETO MÁXIMO DE CONFIANÇA (ABAIXO DE %)</span>
            <span class="slider-val" style="color: var(--cb-yellow);">{{ ((hardCaseThreshold || 0.65) * 100).toFixed(0) }}%</span>
          </div>
          <input type="range" min="0.30" max="0.90" step="0.05" :value="hardCaseThreshold || 0.65" @input="(e) => emit('update:hardCaseThreshold', parseFloat((e.target as HTMLInputElement).value))" />
        </div>
      </div>
    </div>

    <!-- TAB 3: META & VAULT DIRECT STREAM -->
    <div v-else class="tab-pane">
      <div class="selector-group">
        <div class="selector-label">COTA MÁXIMA DE FOTOS POR CLASSE</div>
        <div style="display: flex; gap: 0.3rem; margin-top: 0.3rem;">
          <button v-for="lim in [5, 10, 25, 50, 100]" :key="lim" class="cyber-pill" :class="{ active: maxPerClass === lim }" style="flex: 1; font-size: 0.65rem; padding: 0.2rem;" @click="emit('update:maxPerClass', lim)">{{ lim }}</button>
        </div>
      </div>

      <div class="telemetry-row" style="margin-top: 0.6rem;">
        <span class="k">TOTAL ENVIADO AO VAULT</span>
        <span class="v" style="color: var(--cb-green); font-size: 0.9rem; font-weight: 700;">{{ totalCollected }} FRAMES</span>
      </div>

      <div style="background: var(--vms-bg-elevated); border: 1px solid var(--vms-border); border-radius: 4px; padding: 0.45rem 0.6rem; margin: 0.5rem 0; font-size: 0.68rem; font-family: var(--font-mono);">
        <div v-for="cls in currentDatasetClasses.filter(c => activeClasses[c] !== false)" :key="cls" style="display: flex; justify-content: space-between; margin-bottom: 0.2rem;">
          <span style="color: var(--cb-cyan);">{{ cls.toUpperCase() }}:</span>
          <span :style="{ color: (collectedStats[cls] || 0) >= maxPerClass ? 'var(--cb-green)' : 'var(--cb-yellow)' }">
            {{ collectedStats[cls] || 0 }} / {{ maxPerClass }} {{ (collectedStats[cls] || 0) >= maxPerClass ? '[COMPLETO]' : '' }}
          </span>
        </div>
      </div>
      <button class="cyber-pill" style="width: 100%; font-size: 0.65rem; padding: 0.2rem;" @click="emit('resetCollectorStats')">[RESETAR CONTADORES]</button>
    </div>

    <!-- ACTION CONTROLS (LOCKED IF VAULT OFFLINE) -->
    <div style="padding: 0.6rem; border-top: 1px solid var(--vms-border); background: rgba(0,0,0,0.2);">
      <div v-if="!isVaultOnline" style="color: var(--cb-magenta); font-size: 0.68rem; font-family: var(--font-mono); margin-bottom: 0.4rem; text-align: center;">
        // CONEXAO COM HYDRAVAULT (:8082) REQUERIDA PARA INFERENCIA
      </div>

      <div style="display: flex; gap: 0.4rem;">
        <button class="cyber-action-btn" style="flex: 1; padding: 0.65rem; font-size: 0.72rem;" :disabled="!isVaultOnline || isRunning" @click="emit('runInference')">
          {{ isRunning ? 'SCANNING...' : 'SCAN 1 FRAME' }}
        </button>
        <button
          class="cyber-action-btn"
          style="flex: 1.2; padding: 0.65rem; font-size: 0.72rem;"
          :style="{ background: isContinuous ? (isPaused ? 'var(--cb-yellow)' : 'var(--cb-green)') : 'rgba(0,240,255,0.1)', color: isContinuous ? '#07080c' : 'var(--cb-cyan)', border: '1px solid var(--cb-cyan)' }"
          :disabled="!isVaultOnline"
          @click="emit('update:isContinuous', !isContinuous)"
        >
          {{ isContinuous ? (isPaused ? 'STREAM [PAUSADO]' : 'STREAM [COLETANDO]') : '▶ INICIAR STREAM & COLETAR' }}
        </button>
      </div>
    </div>
  </div>
</template>
