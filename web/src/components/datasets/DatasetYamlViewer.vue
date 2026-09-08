<script setup lang="ts">
import { ref, computed } from 'vue'
import ClassSnapshotModal from './ClassSnapshotModal.vue'
import { autoSuggestCategory } from './yolo_coco_classes'
import type { DatasetInfo } from '../../types/dataset'

const props = defineProps<{ dataset: DatasetInfo | null }>()
const emit = defineEmits<{ (e: 'saveMappings', datasetId: string, mappings: Record<string, string>): void }>()
const classMappings = ref<Record<string, string>>({})
const activeModalClass = ref<string | null>(null)
const saveStatus = ref<string | null>(null)

const handleRemapClass = (rawClass: string, targetCategory: string) => {
  classMappings.value[rawClass] = targetCategory
  emit('saveMappings', props.dataset?.id || props.dataset?.dataset_id || '', classMappings.value)
  saveStatus.value = 'OK'
  setTimeout(() => saveStatus.value = null, 2500)
}

const rawClasses = computed(() => props.dataset?.classes || [])
const uniqueTargetClasses = computed(() => {
  const targetNames = rawClasses.value
    .map((cls: string) => classMappings.value[cls] || autoSuggestCategory(cls))
    .filter((c: string) => c !== 'ignore' && c !== 'ignorar')
  return Array.from(new Set(targetNames))
})
const numClasses = computed(() => uniqueTargetClasses.value.length > 0 ? uniqueTargetClasses.value.length : 1)

const yamlString = computed(() => {
  const id = props.dataset?.id || props.dataset?.dataset_id || 'dataset'
  const pathStr = props.dataset?.yaml_path ? props.dataset.yaml_path.replace('/data.yaml', '') : `/home/hades/datasets/${id}`
  return `# Ultralytics YOLO26 Dataset Config\npath: ${pathStr}\ntrain: train/images\nval: valid/images\ntest: test/images\n\nnc: ${numClasses.value}\nnames:\n${uniqueTargetClasses.value.map((c, i) => `  ${i}: ${c}`).join('\n')}`
})

const copyYaml = () => { if (typeof navigator !== 'undefined' && navigator.clipboard) navigator.clipboard.writeText(yamlString.value) }
const applyMappings = () => {
  if (!props.dataset) return
  emit('saveMappings', props.dataset.id || props.dataset.dataset_id || '', classMappings.value)
  saveStatus.value = 'OK'
  setTimeout(() => saveStatus.value = null, 2000)
}
</script>

<template>
  <div v-if="!dataset" class="cyber-card">
    <div class="card-header"><span class="card-title">2. CLASS MAPPER</span></div>
    <p style="color: #64748b; font-size: 0.8rem;">No dataset selected.</p>
  </div>

  <div v-else class="cyber-card" style="height: 380px; display: flex; flex-direction: column;">
    <div class="card-header" style="margin-bottom: 0.4rem; padding-bottom: 0.35rem; flex-shrink: 0;">
      <span class="card-title">2. CLASS MAPPER ({{ rawClasses.length }} CLASSES)</span>
      <span class="badge-cyan">CLICK CLASS FOR SNAPSHOT</span>
    </div>

    <div style="flex-shrink: 0; margin-bottom: 0.3rem;">
      <div style="display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.3rem; max-height: 120px; overflow-y: auto; padding-right: 2px;">
        <button v-for="cls in rawClasses" :key="cls" type="button" class="class-mapper-row" style="width: 100%; cursor: pointer; text-align: left;" @click="activeModalClass = cls">
          <span class="class-mapper-label" :title="cls">[CLS] "{{ cls }}"</span>
          <span class="badge-yellow" style="font-size: 0.62rem; padding: 0.1rem 0.35rem;">➔ {{ classMappings[cls] || autoSuggestCategory(cls) }}</span>
        </button>
      </div>
    </div>

    <pre class="yaml-code-box" style="height: 90px; max-height: 90px; font-size: 0.68rem; flex-shrink: 0; margin-top: 0.3rem;"><code>{{ yamlString }}</code></pre>

    <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 0.4rem; flex-shrink: 0;">
      <span :style="{ fontSize: '0.68rem', color: saveStatus ? 'var(--cb-green)' : 'var(--cb-cyan)', fontFamily: 'var(--font-mono)' }">
        {{ saveStatus ? '[SAVED] MAPPING SAVED' : `nc=${numClasses} (${uniqueTargetClasses.join(', ')})` }}
      </span>
      <div style="display: flex; gap: 0.4rem;">
        <button class="cyber-action-btn secondary" style="padding: 0.35rem 0.8rem; font-size: 0.7rem;" @click="copyYaml">COPY</button>
        <button class="cyber-action-btn" style="padding: 0.35rem 0.8rem; font-size: 0.7rem;" @click="applyMappings">APPLY</button>
      </div>
    </div>

    <ClassSnapshotModal :isOpen="Boolean(activeModalClass)" :datasetId="dataset.id || dataset.dataset_id" :className="activeModalClass" :currentTarget="activeModalClass ? (classMappings[activeModalClass] || autoSuggestCategory(activeModalClass)) : undefined" @close="activeModalClass = null" @remap="handleRemapClass" />
  </div>
</template>
