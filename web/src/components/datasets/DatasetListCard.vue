<script setup lang="ts">
import DatasetCyberIcons from './DatasetCyberIcons.vue'
import type { DatasetInfo } from '../../types/dataset'

withDefaults(defineProps<{
  datasets?: DatasetInfo[]
  selectedId?: string
}>(), {
  datasets: () => [],
  selectedId: ''
})

const emit = defineEmits<{
  (e: 'selectDataset', ds: DatasetInfo): void
  (e: 'openImportModal'): void
  (e: 'openAnnotateModal'): void
  (e: 'openMergeModal'): void
  (e: 'rescan'): void
  (e: 'deleteDataset', ds: DatasetInfo): void
}>()

const isMerged = (ds: DatasetInfo) => ds.id.includes('merged') || ds.id.includes('frota') || ds.id.includes('fusion')
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">1. REGISTERED DATASETS</span>
      <div style="display: flex; gap: 0.35rem; align-items: center;">
        <button class="cyber-pill" title="Rescan local folders" style="font-size: 0.62rem; padding: 0.15rem 0.4rem;" @click="emit('rescan')">
          [RESCAN]
        </button>
        <span class="badge-cyan">{{ datasets.length }} REPOS</span>
      </div>
    </div>

    <div class="ds-list-scroll">
      <div v-if="datasets.length === 0" style="font-size: 0.8rem; color: #94a3b8; textAlign: center; padding: 1.2rem 0;">
        No registered datasets found. Click <strong>IMPORT</strong> or <strong>RESCAN</strong>.
      </div>
      <div
        v-for="ds in datasets"
        v-else
        :key="ds.id"
        class="dataset-card-item"
        :class="{ active: selectedId === ds.id }"
        @click="emit('selectDataset', ds)"
      >
        <div style="min-width: 0; flex: 1; padding-right: 0.5rem;">
          <div style="display: flex; align-items: center; gap: 0.4rem;">
            <span :style="{ color: selectedId === ds.id ? 'var(--cb-yellow)' : (isMerged(ds) ? 'var(--cb-green)' : 'var(--cb-cyan)'), fontSize: '0.8rem' }">
              {{ selectedId === ds.id ? '▶' : '◈' }}
            </span>
            <span class="ds-name">{{ ds.name }}</span>
            <span :class="isMerged(ds) ? 'badge-green' : 'badge-cyan'" style="font-size: 0.55rem; padding: 0.05rem 0.25rem;">
              {{ isMerged(ds) ? 'MERGED' : 'RAW' }}
            </span>
          </div>
          <div class="ds-meta">
            <span class="ds-tag">[IMGS] {{ ((ds.train_count || 0) + (ds.val_count || 0)).toLocaleString() }}</span>
            <span class="ds-tag">[CLS] {{ ds.classes?.length || 0 }}</span>
            <span class="ds-tag">[SPLIT] {{ ds.train_count || 0 }} train / {{ ds.val_count || 0 }} val</span>
          </div>
        </div>
        <div style="display: flex; align-items: center; gap: 0.4rem;">
          <div style="display: flex; flex-direction: column; align-items: flex-end; gap: 0.2rem;">
            <span class="badge-yellow" style="font-size: 0.62rem; padding: 0.1rem 0.35rem;">{{ ds.format || 'YOLO' }}</span>
            <span :style="{ fontSize: '0.62rem', fontFamily: 'var(--font-mono)', color: selectedId === ds.id ? 'var(--cb-yellow)' : '#64748b' }">{{ selectedId === ds.id ? 'ACTIVE' : 'SELECT' }}</span>
          </div>
          <button
            type="button"
            class="cyber-pill"
            title="Delete dataset"
            style="padding: 0.2rem 0.4rem; font-size: 0.7rem; color: var(--cb-magenta); border-color: rgba(255,0,60,0.3);"
            @click.stop="emit('deleteDataset', ds)"
          >
            [DEL]
          </button>
        </div>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.4rem; margin-top: 0.5rem;">
      <button class="cyber-action-btn secondary" :disabled="datasets.length === 0" @click="emit('openAnnotateModal')">
        <DatasetCyberIcons name="tag" :size="12" /> <span>ANNOTATE</span>
      </button>
      <button class="cyber-action-btn secondary" @click="emit('openMergeModal')">
        <DatasetCyberIcons name="fusion" :size="12" /> <span>MERGE</span>
      </button>
      <button class="cyber-action-btn" @click="emit('openImportModal')">
        <DatasetCyberIcons name="zip" :size="12" color="#07080c" /> <span>IMPORT</span>
      </button>
    </div>
  </div>
</template>
