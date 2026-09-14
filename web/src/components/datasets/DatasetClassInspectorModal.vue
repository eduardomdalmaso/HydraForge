<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = defineProps<{ isOpen: boolean; dataset: any | null }>()
const emit = defineEmits<{(e: 'close'): void; (e: 'train', dataset: any): void}>()

const selectedClass = ref<string>('')
const sampleIndex = ref(0)
const sampleList = ref<any[]>([])
const isLoading = ref(false)

const classesList = computed(() => props.dataset?.classes || [])

const fetchClassSamples = async () => {
  if (!props.dataset || !selectedClass.value) return
  isLoading.value = true
  sampleIndex.value = 0
  try {
    const dsId = props.dataset.id || props.dataset.dataset_id
    const res = await fetch(`/api/v1/training/datasets/sample?id=${encodeURIComponent(dsId)}&class=${encodeURIComponent(selectedClass.value)}&t=${Date.now()}`)
    if (res.ok) {
      const data = await res.json()
      sampleList.value = Array.isArray(data) ? data : (data ? [data] : [])
    } else {
      sampleList.value = []
    }
  } catch {
    sampleList.value = []
  } finally {
    isLoading.value = false
  }
}

watch(() => [props.isOpen, props.dataset], () => {
  if (props.isOpen && props.dataset) {
    selectedClass.value = props.dataset.classes?.[0] || '0'
    sampleList.value = []
    sampleIndex.value = 0
    fetchClassSamples()
  }
})

watch(() => selectedClass.value, () => { if (props.isOpen) fetchClassSamples() })

const currentSample = computed(() => sampleList.value[sampleIndex.value] || null)

const getBBoxStyle = (bbox: number[], isTarget = true) => {
  if (!bbox || bbox.length < 4) return {}
  const [x, y, w, h] = bbox
  return {
    position: 'absolute' as const,
    left: `${(x - w / 2) * 100}%`,
    top: `${(y - h / 2) * 100}%`,
    width: `${w * 100}%`,
    height: `${h * 100}%`,
    border: isTarget ? '2px solid var(--vms-primary, #ff5e3a)' : '1px dashed rgba(0, 240, 255, 0.45)',
    backgroundColor: isTarget ? 'rgba(255, 94, 58, 0.25)' : 'rgba(0, 240, 255, 0.05)',
    boxSizing: 'border-box' as const,
    pointerEvents: 'none' as const,
    zIndex: isTarget ? 3 : 1
  }
}

const getBadgeStyle = (isTarget = true) => ({
  position: 'absolute' as const,
  top: '-18px',
  left: '-2px',
  background: isTarget ? 'var(--vms-primary)' : 'rgba(0, 240, 255, 0.75)',
  color: '#fff',
  fontSize: '0.62rem',
  fontFamily: 'monospace',
  fontWeight: 'bold' as const,
  padding: '1px 4px',
  borderRadius: '2px',
  whiteSpace: 'nowrap' as const
})

const nextSample = () => {
  if (sampleList.value.length > 1) sampleIndex.value = (sampleIndex.value + 1) % sampleList.value.length
  else fetchClassSamples()
}
const prevSample = () => {
  if (sampleList.value.length > 1) sampleIndex.value = (sampleIndex.value - 1 + sampleList.value.length) % sampleList.value.length
}
</script>

<template>
  <div v-if="isOpen && dataset" class="modal-backdrop" @click="emit('close')">
    <div class="import-modal-box" style="max-width: 840px; width: 95%;" @click.stop>
      <div class="card-header" style="margin-bottom: 0.85rem;">
        <div>
          <span class="card-title">INSPEÇÃO DE CLASSES // {{ dataset.name || dataset.id || dataset.dataset_id }}</span>
          <div class="text-mono" style="font-size: 0.68rem; color: #8b94a0; margin-top: 2px;">
            {{ classesList.length }} CLASSES ANOTADAS • {{ (dataset.train_images || dataset.train_count || 0) + (dataset.val_images || dataset.val_count || 0) }} IMAGENS TOTAIS
          </div>
        </div>
        <button class="cyber-pill" style="padding: 0.2rem 0.5rem;" @click="emit('close')">✕</button>
      </div>

      <div style="display: grid; grid-template-columns: 240px 1fr; gap: 1rem; min-height: 380px;">
        <!-- LEFT: CLASSES SELECTOR -->
        <div style="display: flex; flex-direction: column; background: #040609; border: 1px solid var(--vms-border); border-radius: 4px; padding: 0.6rem;">
          <div class="text-mono" style="font-size: 0.7rem; color: var(--vms-primary); font-weight: 700; margin-bottom: 0.5rem;">
            SELECIONE A CLASSE:
          </div>
          <div style="max-height: 320px; overflow-y: auto; display: flex; flex-direction: column; gap: 0.35rem; padding-right: 3px;">
            <div
              v-for="(cls, idx) in classesList"
              :key="cls"
              class="dataset-card-item"
              :class="{ active: selectedClass === cls }"
              style="margin-bottom: 0; padding: 0.45rem 0.65rem; cursor: pointer;"
              @click="selectedClass = cls"
            >
              <div class="text-mono" style="font-size: 0.78rem; font-weight: 600;">#{{ idx }} {{ cls }}</div>
              <span v-if="selectedClass === cls" class="badge-yellow" style="font-size: 0.6rem; padding: 1px 4px;">FOTO ATIVA</span>
            </div>
          </div>
        </div>

        <!-- RIGHT: PHOTO PREVIEWER WITH RESPONSIVE BBOX OVERLAY -->
        <div style="display: flex; flex-direction: column; justify-content: space-between; background: #040609; border: 1px solid var(--vms-border); border-radius: 4px; padding: 0.85rem;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
            <span class="text-mono" style="font-size: 0.75rem; color: #cbd5e1;">
              AMOSTRA DA CLASSE: <strong style="color: var(--vms-primary);">{{ selectedClass.toUpperCase() }}</strong>
            </span>
            <span v-if="sampleList.length > 0" class="badge-cyan" style="font-size: 0.65rem;">
              FOTO {{ sampleIndex + 1 }} / {{ sampleList.length }}
            </span>
          </div>

          <!-- PHOTO CANVAS -->
          <div style="flex: 1; min-height: 280px; background: #010204; border: 1px solid rgba(255,255,255,0.06); border-radius: 4px; display: flex; align-items: center; justify-content: center; overflow: hidden; position: relative; padding: 0.5rem;">
            <div v-if="isLoading" style="color: var(--vms-primary); font-family: var(--font-mono); font-size: 0.8rem;">
              CARREGANDO FOTO DO DATASET...
            </div>
            <div v-else-if="currentSample?.image_url" style="position: relative; display: inline-block; max-width: 100%; max-height: 290px;">
              <img :src="currentSample.image_url" style="max-width: 100%; max-height: 290px; display: block; object-fit: contain; border-radius: 3px;" />
              <!-- ALL ANNOTATIONS -->
              <template v-if="currentSample.annotations && currentSample.annotations.length > 0">
                <div v-for="(ann, aIdx) in currentSample.annotations" :key="aIdx" :style="getBBoxStyle(ann.bbox, ann.is_target)">
                  <span :style="getBadgeStyle(ann.is_target)">{{ ann.class_name || (ann.is_target ? selectedClass : `class_${ann.class_id}`) }}</span>
                </div>
              </template>
              <template v-else-if="currentSample?.bbox">
                <div :style="getBBoxStyle(currentSample.bbox, true)">
                  <span :style="getBadgeStyle(true)">{{ selectedClass }}</span>
                </div>
              </template>
            </div>
            <div v-else style="color: #64748b; font-family: var(--font-mono); font-size: 0.75rem;">
              Nenhuma anotação disponível para esta classe.
            </div>
          </div>

          <!-- CYCLE BUTTONS -->
          <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 0.65rem;">
            <div style="display: flex; gap: 0.4rem;">
              <button class="cyber-pill" :disabled="isLoading" @click="prevSample">◀ ANTERIOR</button>
              <button class="cyber-pill" :disabled="isLoading" @click="nextSample">PRÓXIMA FOTO ▶</button>
            </div>
            <button class="cyber-action-btn" style="padding: 0.4rem 0.85rem; font-size: 0.75rem;" @click="emit('train', dataset)">
              INICIAR TREINO COM ESTE DATASET ➔
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
