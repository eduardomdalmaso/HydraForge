<script setup lang="ts">
import { ref, computed } from 'vue'
import DatasetCyberIcons from './DatasetCyberIcons.vue'

const props = withDefaults(defineProps<{
  datasets?: any[]
  selectedId?: string
  isRescanning?: boolean
}>(), {
  datasets: () => [],
  selectedId: '',
  isRescanning: false
})

const emit = defineEmits<{
  (e: 'selectDataset', ds: any): void
  (e: 'inspectClasses', ds: any): void
  (e: 'openPipelineModal'): void
  (e: 'openImportModal'): void
  (e: 'openAnnotateModal'): void
  (e: 'rescan'): void
  (e: 'deleteDataset', ds: any): void
}>()

const searchQuery = ref('')

const filteredDatasets = computed(() => {
  if (!searchQuery.value.trim()) return props.datasets
  const q = searchQuery.value.toLowerCase()
  return props.datasets.filter(d => 
    (d.name || d.id || d.dataset_id || '').toLowerCase().includes(q) ||
    (d.classes || []).some((c: string) => c.toLowerCase().includes(q))
  )
})
</script>

<template>
  <div class="cyber-card" style="padding: 1.25rem;">
    <!-- WINDOWS EXPLORER PATH & ACTION BAR -->
    <div style="display: flex; justify-content: space-between; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; flex-wrap: wrap;">
      <div style="display: flex; align-items: center; gap: 0.5rem; flex: 1; min-width: 280px; background: #040609; border: 1px solid var(--vms-border); border-radius: 4px; padding: 0.4rem 0.75rem;">
        <span style="color: var(--vms-primary); font-family: var(--font-mono); font-size: 0.82rem; font-weight: 700;">DIR:</span>
        <span class="text-mono" style="color: #cbd5e1; font-size: 0.8rem; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
          datasets/
        </span>
        <span class="badge-cyan" style="font-size: 0.7rem; padding: 2px 7px;">{{ filteredDatasets.length }} PASTAS</span>
      </div>

      <div style="display: flex; align-items: center; gap: 0.5rem;">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Buscar dataset ou classe..."
          style="background: #040609; border: 1px solid var(--vms-border); color: #fff; padding: 0.4rem 0.75rem; border-radius: 4px; font-family: var(--font-mono); font-size: 0.8rem; width: 220px;"
        />
        <button class="cyber-action-btn" style="padding: 0.45rem 1rem; font-size: 0.8rem;" @click="emit('openPipelineModal')">
          <DatasetCyberIcons name="pipeline" :size="16" color="#ffffff" />
          <span>NOVO PIPELINE</span>
        </button>
        <button class="cyber-action-btn secondary" style="padding: 0.45rem 0.8rem; font-size: 0.8rem;" @click="emit('openImportModal')">
          <span>+ IMPORTAR ZIP</span>
        </button>
        <button class="cyber-pill" :disabled="isRescanning" style="padding: 0.45rem 0.75rem; font-size: 0.8rem;" @click="emit('rescan')">
          <span>{{ isRescanning ? 'REESCANEANDO...' : 'REESCANEAR' }}</span>
        </button>
      </div>
    </div>

    <!-- WINDOWS EXPLORER FOLDER GRID -->
    <div v-if="filteredDatasets.length === 0" style="padding: 3rem 1rem; text-align: center; color: #8b94a0; font-family: var(--font-mono); font-size: 0.85rem;">
      Nenhum dataset encontrado no diretório. Clique em "+ NOVO PIPELINE (4 ETAPAS)" ou "+ IMPORTAR ZIP".
    </div>

    <div v-else style="display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem;">
      <div
        v-for="ds in filteredDatasets"
        :key="ds.id || ds.dataset_id"
        class="cyber-card"
        :style="{
          padding: '1.1rem',
          cursor: 'pointer',
          border: (selectedId === (ds.id || ds.dataset_id)) ? '1px solid var(--vms-primary)' : '1px solid var(--vms-border)',
          background: (selectedId === (ds.id || ds.dataset_id)) ? 'rgba(255, 94, 58, 0.08)' : 'var(--vms-bg-elevated)',
          display: 'flex',
          flexDirection: 'column',
          gap: '0.65rem',
          transition: 'all 0.15s ease'
        }"
        @click="emit('selectDataset', ds)"
      >
        <div style="display: flex; justify-content: space-between; align-items: flex-start; gap: 0.5rem;">
          <div style="display: flex; align-items: center; gap: 0.65rem; min-width: 0;">
            <!-- FOLDER SVG ICON -->
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="var(--vms-primary)" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0;">
              <path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"></path>
              <polygon points="12 10 12 14 16 14"></polygon>
            </svg>
            <div style="min-width: 0;">
              <div class="text-mono" style="font-size: 0.9rem; font-weight: 700; color: #ffffff; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
                {{ ds.name || ds.id || ds.dataset_id }}
              </div>
              <div class="text-mono" style="font-size: 0.68rem; color: #8b94a0;">
                {{ ds.task ? `[${ds.task.toUpperCase()}]` : '[DETECT]' }} // YAML ATIVO
              </div>
            </div>
          </div>
          <button
            class="cyber-pill"
            style="padding: 2px 6px; font-size: 0.7rem; color: #8b94a0;"
            title="Excluir Dataset"
            @click.stop="emit('deleteDataset', ds)"
          >✕</button>
        </div>

        <div style="display: flex; gap: 0.4rem; flex-wrap: wrap; margin-top: 0.2rem;">
          <span class="ds-tag">{{ (ds.classes || []).length || ds.num_classes || 0 }} classes</span>
          <span class="ds-tag">{{ ((ds.train_count || ds.train_images || 0) + (ds.val_count || ds.val_images || 0)).toLocaleString() }} imgs</span>
          <span class="ds-tag" v-if="ds.size_bytes">{{ (ds.size_bytes / (1024 * 1024)).toFixed(1) }} MB</span>
        </div>

        <div style="display: flex; justify-content: space-between; align-items: center; border-top: 1px solid rgba(255, 255, 255, 0.05); padding-top: 0.5rem; margin-top: 0.35rem;">
          <span class="text-mono" style="font-size: 0.68rem; color: var(--vms-success);">[PRONTO]</span>
          <button
            class="cyber-pill"
            style="padding: 0.2rem 0.6rem; font-size: 0.72rem; color: var(--vms-primary); font-weight: 700;"
            @click.stop="emit('inspectClasses', ds)"
          >
            VER CLASSES & FOTOS ➔
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
