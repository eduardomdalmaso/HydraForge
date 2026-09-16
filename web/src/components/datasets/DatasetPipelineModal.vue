<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import DatasetMergerRow from './DatasetMergerRow.vue'
import { autoSuggestCategory } from './yolo_coco_classes'

const props = withDefaults(defineProps<{ isOpen: boolean; datasets?: any[]; savedMappings?: Record<string, any> }>(), {
  datasets: () => [], savedMappings: () => ({})
})
const emit = defineEmits<{ (e: 'close'): void; (e: 'pipelineComplete', result: any): void }>()

const activeStep = ref(1)
const mergedName = ref('frota_urbana_fusion')
const isProcessing = ref(false)
const selectedDsIds = ref<string[]>([])
const mappings = ref<Record<string, any>>({})
const kfoldCount = ref(5)
const splitTrainRatio = ref(80)
const kfoldSeed = ref(42)
const enableLeakageDefense = ref(true)
const enableBboxGuard = ref(true)
const enableEmptyLabelsClean = ref(true)

watch(() => props.isOpen, (open) => {
  if (!open) return
  activeStep.value = 1
  const initial: Record<string, any> = {}
    props.datasets.forEach(ds => {
    const id = ds.id || ds.dataset_id
    initial[id] = {}
    ;(ds.classes || []).forEach((cls: string) => { initial[id][cls] = props.savedMappings[id]?.[cls] || autoSuggestCategory(cls, ds.name || id) })
  })
  mappings.value = initial
  selectedDsIds.value = props.datasets.map(d => d.id || d.dataset_id)
})

const toggleDs = (id: string) => { selectedDsIds.value = selectedDsIds.value.includes(id) ? selectedDsIds.value.filter(x => x !== id) : [...selectedDsIds.value, id] }
const setMap = (dsId: string, cls: string, val: string) => { mappings.value = { ...mappings.value, [dsId]: { ...(mappings.value[dsId] || {}), [cls]: val } } }

const targetOptions = computed(() => Array.from(new Set([...props.datasets.flatMap(d => d.classes || []), ...Object.values(mappings.value).flatMap(m => Object.values(m)), 'car', 'motorcycle', 'truck', 'bus', 'person', 'bicycle', 'cell-phone', 'ignore'])).filter(Boolean))
const totalImgs = computed(() => props.datasets.filter(d => selectedDsIds.value.includes(d.id || d.dataset_id)).reduce((acc, d) => acc + (d.train_count || d.train_images || 0) + (d.val_count || d.val_images || 0), 0))

const handleExecutePipeline = async () => {
  isProcessing.value = true
  try {
    const activeTargets = Array.from(new Set(selectedDsIds.value.flatMap(id => Object.values(mappings.value[id] || {})))).filter(c => c && c !== 'ignore')
    const res = await fetch('/api/v1/training/datasets/merge', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_name: mergedName.value, dataset_ids: selectedDsIds.value, mappings: mappings.value, classes: activeTargets.length ? activeTargets : ['car', 'motorcycle', 'truck', 'bus'] })
    })
    if (res.ok) {
      emit('pipelineComplete', await res.json())
      if (typeof window !== 'undefined') window.location.reload()
    }
  } finally {
    isProcessing.value = false
    emit('close')
  }
}
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click="emit('close')">
    <div class="import-modal-box" style="max-width: 780px; width: 95%;" @click.stop>
      <div class="card-header" style="margin-bottom: 0.75rem;">
        <span class="card-title">PIPELINE DE CRIAÇÃO E PROCESSAMENTO DE DATASET</span>
        <button class="cyber-pill" style="padding: 0.2rem 0.5rem;" @click="emit('close')">✕</button>
      </div>

      <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.4rem; margin-bottom: 1rem;">
        <button v-for="(sName, sIdx) in ['1. MERGE', '2. CLASS MAPPER', '3. K-FOLD SPLIT', '4. QUALITY AUDIT']" :key="sIdx" class="cyber-pill" :class="{ active: activeStep === (sIdx + 1) }" style="padding: 0.45rem 0.2rem; font-size: 0.68rem; text-align: center;" @click="activeStep = sIdx + 1">{{ sName }}</button>
      </div>

      <!-- STEP 1: MERGE & DATASETS SELECTION -->
      <div v-if="activeStep === 1" style="display: flex; flex-direction: column; gap: 0.75rem;">
        <div style="display: flex; gap: 0.5rem; align-items: center;">
          <span class="text-mono" style="font-size: 0.75rem; color: #8b94a0;">NOME FINAL:</span>
          <input v-model="mergedName" type="text" style="flex: 1; background: #040609; border: 1px solid var(--vms-border); color: #fff; padding: 0.35rem 0.65rem; font-size: 0.8rem; border-radius: 4px; font-family: var(--font-mono);" />
        </div>
        <div class="text-mono" style="font-size: 0.7rem; color: #8b94a0;">SELECIONE OS DATASETS BASE PARA A FUSÃO:</div>
        <div style="max-height: 220px; overflow-y: auto; display: flex; flex-direction: column; gap: 0.35rem;">
          <div v-for="ds in datasets" :key="ds.id || ds.dataset_id" class="dataset-card-item" :class="{ active: selectedDsIds.includes(ds.id || ds.dataset_id) }" style="margin-bottom: 0; padding: 0.5rem 0.75rem;" @click="toggleDs(ds.id || ds.dataset_id)">
            <div style="display: flex; align-items: center; gap: 0.5rem;">
              <input type="checkbox" :checked="selectedDsIds.includes(ds.id || ds.dataset_id)" style="accent-color: var(--vms-primary);" />
              <span class="ds-name">{{ ds.name || ds.id || ds.dataset_id }}</span>
            </div>
            <span class="ds-tag">{{ (ds.train_count || ds.train_images || 0) + (ds.val_count || ds.val_images || 0) }} imgs</span>
          </div>
        </div>
      </div>

      <!-- STEP 2: CLASS MAPPER -->
      <div v-else-if="activeStep === 2" style="max-height: 280px; overflow-y: auto; padding-right: 4px;">
        <div class="text-mono" style="font-size: 0.7rem; color: #8b94a0; margin-bottom: 0.5rem;">UNIFICAÇÃO E MAPEAMENTO DE CLASSES ENTRE AS FONTES:</div>
        <DatasetMergerRow v-for="ds in datasets.filter(d => selectedDsIds.includes(d.id || d.dataset_id))" :key="ds.id || ds.dataset_id" :dataset="ds" :isSelected="true" :mappings="mappings" :targetOptions="targetOptions" @toggle="toggleDs" @setMap="setMap" />
      </div>

      <!-- STEP 3: K-FOLD CROSS-VALIDATION -->
      <div v-else-if="activeStep === 3" style="display: flex; flex-direction: column; gap: 0.85rem;">
        <div>
          <div class="text-mono" style="font-size: 0.7rem; color: #8b94a0; margin-bottom: 0.35rem;">NÚMERO DE FOLDS (CROSS-VALIDATION):</div>
          <div style="display: flex; gap: 0.4rem;">
            <button v-for="k in [1, 3, 5, 10]" :key="k" class="cyber-pill" :class="{ active: kfoldCount === k }" style="flex: 1; padding: 0.4rem 0.2rem; font-size: 0.72rem; text-align: center;" @click="kfoldCount = k">
              {{ k === 1 ? '[SEM FOLDS / DIRETO]' : `[${k}-FOLDS ESTRATIFICADO]` }}
            </button>
          </div>
        </div>

        <div>
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.35rem;">
            <span class="text-mono" style="font-size: 0.7rem; color: #8b94a0;">DIVISÃO TREINO / VALIDAÇÃO:</span>
            <span class="text-mono" style="font-size: 0.75rem; color: var(--vms-primary); font-weight: 700;">{{ splitTrainRatio }}% TREINO / {{ 100 - splitTrainRatio }}% VALIDAÇÃO</span>
          </div>
          <input v-model.number="splitTrainRatio" type="range" min="50" max="95" step="5" style="width: 100%; accent-color: var(--vms-primary); cursor: pointer;" />
          <div class="kfold-bar" style="height: 10px; margin-top: 0.4rem;">
            <div class="kfold-seg-train" :style="{ width: `${splitTrainRatio}%` }"></div>
            <div class="kfold-seg-val" :style="{ width: `${100 - splitTrainRatio}%` }"></div>
          </div>
        </div>

        <div style="display: flex; gap: 0.5rem; align-items: center;">
          <span class="text-mono" style="font-size: 0.72rem; color: #8b94a0;">RANDOM SEED:</span>
          <input v-model.number="kfoldSeed" type="number" style="width: 90px; background: #040609; border: 1px solid var(--vms-border); color: #fff; padding: 0.25rem 0.5rem; font-size: 0.75rem; border-radius: 4px; font-family: var(--font-mono);" />
          <button class="cyber-pill" style="padding: 0.25rem 0.6rem; font-size: 0.7rem;" @click="kfoldSeed = Math.floor(Math.random() * 9000 + 1000)">[GERAR SEED]</button>
        </div>
      </div>

      <!-- STEP 4: QUALITY AUDITOR & LEAKAGE DEFENSE -->
      <div v-else-if="activeStep === 4" style="display: flex; flex-direction: column; gap: 0.65rem;">
        <label class="telemetry-row" style="cursor: pointer;"><span class="k">LEAKAGE DEFENSE (PHASH / COSINE SIMILARITY)</span><input v-model="enableLeakageDefense" type="checkbox" style="accent-color: var(--vms-primary);" /></label>
        <label class="telemetry-row" style="cursor: pointer;"><span class="k">BBOX BOUNDS GUARD ([0.0, 1.0])</span><input v-model="enableBboxGuard" type="checkbox" style="accent-color: var(--vms-primary);" /></label>
        <label class="telemetry-row" style="cursor: pointer;"><span class="k">EMPTY LABELS FILTER (NEGATIVE BACKGROUND SAMPLES)</span><input v-model="enableEmptyLabelsClean" type="checkbox" style="accent-color: var(--vms-primary);" /></label>
      </div>

      <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 1rem; border-top: 1px solid var(--vms-border); padding-top: 0.75rem;">
        <span class="text-mono" style="font-size: 0.72rem; color: var(--vms-primary);">{{ selectedDsIds.length }} repositórios • {{ totalImgs.toLocaleString() }} imagens</span>
        <div style="display: flex; gap: 0.5rem;">
          <button v-if="activeStep > 1" class="cyber-action-btn secondary" @click="activeStep--">VOLTAR</button>
          <button v-if="activeStep < 4" class="cyber-action-btn" @click="activeStep++">AVANÇAR ➔</button>
          <button v-else class="cyber-action-btn" :disabled="isProcessing || selectedDsIds.length === 0" @click="handleExecutePipeline">{{ isProcessing ? 'PROCESSANDO...' : 'EXECUTAR PIPELINE' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
