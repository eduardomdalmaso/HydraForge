<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import DatasetCyberIcons from './DatasetCyberIcons.vue'
import DatasetMergerRow from './DatasetMergerRow.vue'
import { autoSuggestCategory } from './yolo_coco_classes'

const props = withDefaults(defineProps<{ isOpen: boolean; datasets?: any[]; savedMappings?: Record<string, any> }>(), {
  datasets: () => [], savedMappings: () => ({})
})
const emit = defineEmits<{ (e: 'close'): void; (e: 'datasetsMerged', result: any): void }>()

const mergedName = ref('fusion_custom_dataset')
const isMerging = ref(false)
const selectedDsIds = ref<string[]>([])
const mappings = ref<Record<string, any>>({})

const getSuggestedName = (ids: string[]) => {
  const bases = Array.from(new Set(ids.map(id => id.replace(/\.v\d+.*$/, '').replace(/_v\d+.*$/, '').replace(/[-_]raw$/, ''))))
  return bases.length === 1 ? `${bases[0]}_merged_${ids.length}ds` : (bases.length ? `${bases.slice(0, 2).join('_')}_fusion` : 'fusion_custom_dataset')
}

watch(() => props.isOpen, (open) => {
  if (!open) return
  const initial: Record<string, any> = {}
  props.datasets.forEach(ds => {
    const id = ds.id || ds.dataset_id
    initial[id] = {}
    ;(ds.classes || []).forEach((cls: string) => {
      initial[id][cls] = props.savedMappings[id]?.[cls] || autoSuggestCategory(cls)
    })
  })
  mappings.value = initial
  selectedDsIds.value = []
  mergedName.value = 'fusion_custom_dataset'
})

const toggle = (id: string) => {
  selectedDsIds.value = selectedDsIds.value.includes(id) ? selectedDsIds.value.filter(x => x !== id) : [...selectedDsIds.value, id]
  mergedName.value = getSuggestedName(selectedDsIds.value)
}

const setMap = (dsId: string, cls: string, val: string) => {
  mappings.value = { ...mappings.value, [dsId]: { ...(mappings.value[dsId] || {}), [cls]: val } }
}

const targetOptions = computed(() => Array.from(new Set([...props.datasets.flatMap(d => d.classes || []), ...Object.values(mappings.value).flatMap(m => Object.values(m)), 'car', 'motorcycle', 'truck', 'bus', 'cell-phone', 'ignore'])).filter(Boolean))
const totalImgs = computed(() => props.datasets.filter(d => selectedDsIds.value.includes(d.id || d.dataset_id)).reduce((acc, d) => acc + (d.train_count || d.train_images || 0) + (d.val_count || d.val_images || 0), 0))

const handleExecute = async () => {
  isMerging.value = true
  try {
    const activeTargets = Array.from(new Set(selectedDsIds.value.flatMap(id => Object.values(mappings.value[id] || {})))).filter(c => c && c !== 'ignore')
    const res = await fetch('/api/v1/training/datasets/merge', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_name: mergedName.value, dataset_ids: selectedDsIds.value, mappings: mappings.value, classes: activeTargets.length ? activeTargets : ['cell-phone'] })
    })
    if (res.ok) {
      emit('datasetsMerged', await res.json())
      if (typeof window !== 'undefined') window.location.reload()
    }
  } finally {
    isMerging.value = false
    emit('close')
  }
}
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click="emit('close')">
    <div class="import-modal-box" style="max-width: 720px;" @click.stop>
      <div class="card-header" style="margin-bottom: 0.5rem;">
        <span class="card-title">DATASET HARMONIZER & FUSION ENGINE (MERGE)</span>
        <button class="cyber-pill" style="padding: 0.2rem 0.5rem;" @click="emit('close')">✕</button>
      </div>

      <div style="display: flex; gap: 0.5rem; align-items: center; margin-bottom: 0.6rem;">
        <span style="font-size: 0.72rem; color: #94a3b8; font-family: var(--font-mono);">NAME:</span>
        <input type="text" :value="mergedName" style="flex: 1; background: #05070a; border: 1px solid rgba(0,240,255,0.3); color: #fff; padding: 0.25rem 0.5rem; font-size: 0.75rem; border-radius: 3px; font-family: var(--font-mono);" @input="(e) => mergedName = (e.target as HTMLInputElement).value" />
        <button type="button" class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.62rem;" @click="() => { const all = datasets.map(d => d.id || d.dataset_id); selectedDsIds = all; mergedName = getSuggestedName(all); }">ALL</button>
        <button type="button" class="cyber-pill" style="padding: 0.15rem 0.4rem; font-size: 0.62rem;" @click="() => { selectedDsIds = []; mergedName = 'fusion_custom_dataset'; }">NONE</button>
      </div>

      <div style="max-height: 240px; overflow-y: auto; padding-right: 4px; margin-bottom: 0.75rem;">
        <DatasetMergerRow v-for="ds in datasets" :key="ds.id || ds.dataset_id" :dataset="ds" :isSelected="selectedDsIds.includes(ds.id || ds.dataset_id)" :mappings="mappings" :targetOptions="targetOptions" @toggle="toggle" @setMap="setMap" />
      </div>

      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span style="font-size: 0.72rem; color: var(--cb-cyan); font-family: var(--font-mono);">{{ selectedDsIds.length }} repos selected • {{ totalImgs.toLocaleString() }} images</span>
        <div style="display: flex; gap: 0.5rem;">
          <button class="cyber-action-btn secondary" :disabled="isMerging" @click="emit('close')">CANCEL</button>
          <button class="cyber-action-btn" :disabled="isMerging || selectedDsIds.length === 0" @click="handleExecute">
            <DatasetCyberIcons name="fusion" :size="14" color="#07080c" /><span>{{ isMerging ? 'MERGING...' : 'FUSE DATASETS' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
