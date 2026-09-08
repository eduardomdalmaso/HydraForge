<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { YOLO_CLASSES, decodeYoloClass } from './yolo_coco_classes'
import DatasetCyberIcons from './DatasetCyberIcons.vue'

const props = defineProps<{ isOpen: boolean; datasetId?: string; className?: string | null; currentTarget?: string }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'saveTarget', origClass: string, targetClass: string): void }>()

const target = ref(props.currentTarget || 'carro')
const search = ref('')
const sample = ref<any>(null)
const isLoading = ref(false)

const filtered = computed(() => YOLO_CLASSES.filter(c => c.label.toLowerCase().includes(search.value.toLowerCase()) || c.en.toLowerCase().includes(search.value.toLowerCase())))

const fetchSample = async () => {
  if (!props.datasetId || !props.className) return
  isLoading.value = true
  try {
    const res = await fetch(`/api/v1/training/datasets/sample?id=${props.datasetId}&class=${encodeURIComponent(props.className)}&t=${Date.now()}`)
    if (res.ok) sample.value = await res.json()
  } finally {
    isLoading.value = false
  }
}

watch(() => [props.isOpen, props.className, props.currentTarget], () => {
  if (props.isOpen && props.className) {
    target.value = props.currentTarget || 'carro'
    search.value = ''
    fetchSample()
  }
})

const bx = computed(() => sample.value?.bbox ? (sample.value.bbox[0] - sample.value.bbox[2] / 2) * 1000 : 0)
const by = computed(() => sample.value?.bbox ? (sample.value.bbox[1] - sample.value.bbox[3] / 2) * 1000 : 0)
const bw = computed(() => sample.value?.bbox ? sample.value.bbox[2] * 1000 : 0)
const bh = computed(() => sample.value?.bbox ? sample.value.bbox[3] * 1000 : 0)
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click="emit('close')">
    <div class="import-modal-box" style="max-width: 680px;" @click.stop>
      <div class="card-header" style="margin-bottom: 0.65rem;">
        <span class="card-title">SNAPSHOT: "{{ className }}" ({{ decodeYoloClass(className) }})</span>
        <button class="cyber-pill" style="padding: 0.2rem 0.5rem;" @click="emit('close')">✕</button>
      </div>

      <div style="display: grid; grid-template-columns: 1.2fr 1fr; gap: 0.85rem;">
        <div>
          <div style="width: 100%; height: 215px; background: #05070a; border: 1px solid rgba(0,240,255,0.3); border-radius: 4px; display: flex; align-items: center; justify-content: center; overflow: hidden;">
            <div v-if="isLoading" style="color: var(--cb-yellow); font-family: var(--font-mono); font-size: 0.75rem;">CARREGANDO...</div>
            <svg v-else-if="sample?.image_url" viewBox="0 0 1000 1000" preserveAspectRatio="xMidYMid meet" style="width: 100%; height: 100%; display: block;">
              <image :href="sample.image_url" x="0" y="0" width="1000" height="1000" preserveAspectRatio="none" />
              <template v-if="sample.bbox">
                <rect :x="bx" :y="by" :width="bw" :height="bh" fill="rgba(0, 240, 255, 0.2)" stroke="#00f0ff" stroke-width="8" />
                <text :x="bx + 10" :y="Math.max(by - 15, 35)" fill="#00f0ff" font-size="42" font-family="monospace" font-weight="bold">{{ className }}</text>
              </template>
            </svg>
            <div v-else style="color: #64748b; font-size: 0.75rem;">Sem preview.</div>
          </div>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.4rem; margin-top: 0.4rem;">
            <button class="cyber-action-btn secondary" style="padding: 0.45rem;" :disabled="isLoading" @click="fetchSample">NEXT</button>
            <button class="cyber-action-btn danger" style="padding: 0.45rem;" @click="() => { if (className) emit('saveTarget', className, 'ignorar'); emit('close'); }">IGNORE</button>
          </div>
        </div>

        <div style="display: flex; flex-direction: column; justify-content: space-between;">
          <div>
            <div class="selector-label" style="margin-bottom: 0.3rem;">ESCOLHA A CLASSE DESTINO:</div>
            <input type="text" placeholder="Filtrar classe..." :value="search" style="width: 100%; background: #05070a; border: 1px solid rgba(0,240,255,0.25); color: #fff; padding: 0.25rem 0.5rem; font-size: 0.7rem; border-radius: 3px; margin-bottom: 0.35rem; font-family: var(--font-mono);" @input="(e) => search = (e.target as HTMLInputElement).value" />
            <div style="height: 135px; overflow-y: auto; display: flex; flex-direction: column; gap: 0.25rem; padding-right: 3px;">
              <div v-for="c in filtered" :key="c.id" class="class-mapper-row" :style="{ cursor: 'pointer', height: '28px', padding: '0.2rem 0.5rem', borderColor: target === c.id ? 'var(--cb-yellow)' : 'rgba(0,240,255,0.15)', background: target === c.id ? 'rgba(252,238,10,0.15)' : 'rgba(0,0,0,0.4)' }" @click="target = c.id">
                <span :style="{ color: target === c.id ? 'var(--cb-yellow)' : '#fff', fontSize: '0.7rem', fontFamily: 'var(--font-mono)', fontWeight: target === c.id ? '700' : '400' }">{{ target === c.id ? '[X] ' : '[-] ' }}{{ c.label }} <span :style="{ color: target === c.id ? '#fef08a' : '#64748b', fontSize: '0.62rem' }">[{{ c.en }}]</span></span>
                <span v-if="target === c.id" class="badge-yellow" style="font-size: 0.55rem; padding: 0.05rem 0.3rem;">SELECTED</span>
              </div>
            </div>
          </div>
          <button class="cyber-action-btn" style="width: 100%; padding: 0.6rem; margin-top: 0.5rem;" @click="() => { if (className) emit('saveTarget', className, target); emit('close'); }">
            <DatasetCyberIcons name="launch" :size="14" color="#07080c" /><span>SALVAR</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
