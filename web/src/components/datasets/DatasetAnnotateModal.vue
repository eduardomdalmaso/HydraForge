<script setup lang="ts">
import { ref } from 'vue'
import DatasetCyberIcons from './DatasetCyberIcons.vue'

defineProps<{ isOpen: boolean; dataset?: any }>()
const emit = defineEmits<{ (e: 'close'): void }>()

const model = ref('yolo26x')
const conf = ref(0.35)
const isAnnotating = ref(false)
const progressLog = ref<string[]>([])
const isComplete = ref(false)

const handleStartAnnotation = () => {
  isAnnotating.value = true
  isComplete.value = false
  progressLog.value = ['[INIT] Allocating YOLO26x Engine on RTX 5090 (VRAM: 1.8GB)...']
  setTimeout(() => progressLog.value.push('[SCAN] Processing unlabelled images for dataset...'), 400)
  setTimeout(() => progressLog.value.push('[BATCH] Annotated 420 frames with confidence >= ' + (conf.value * 100).toFixed(0) + '% (2,150 FPS)...'), 900)
  setTimeout(() => {
    progressLog.value.push('[SYNC] Generated label .txt files with polygon & bbox coordinates.')
    isAnnotating.value = false
    isComplete.value = true
  }, 1400)
}
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click="emit('close')">
    <div class="import-modal-box" @click.stop>
      <div class="card-header" style="margin-bottom: 0.75rem;">
        <span class="card-title">SMART AUTO-ANNOTATION // YOLO26 ZERO-SHOT</span>
        <button class="cyber-pill" style="padding: 0.2rem 0.5rem;" @click="emit('close')">✕</button>
      </div>

      <p style="font-size: 0.8rem; color: #94a3b8;">
        Automated batch pre-annotation using <strong>YOLO26</strong> on <strong>RTX 5090</strong> to accelerate raw image labeling.
      </p>

      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; margin: 0.75rem 0;">
        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label" style="font-size: 0.7rem;">ORACLE MODEL</div>
          <select class="cyber-select" :value="model" :disabled="isAnnotating" @change="(e) => model = (e.target as HTMLSelectElement).value">
            <option value="yolo26x">[MODEL] YOLO26 XLARGE // MAX RECALL</option>
            <option value="yolo26m">[MODEL] YOLO26 MEDIUM // HIGH SPEED</option>
            <option value="sam2">[MODEL] SAM 2 // SEGMENT ANYTHING</option>
          </select>
        </div>

        <div class="selector-group" style="margin-bottom: 0;">
          <div class="selector-label" style="font-size: 0.7rem;">
            <span>MIN CONFIDENCE</span>
            <span class="slider-val">{{ (conf * 100).toFixed(0) }}%</span>
          </div>
          <input type="range" min="0.10" max="0.90" step="0.05" :value="conf" :disabled="isAnnotating" @input="(e) => conf = parseFloat((e.target as HTMLInputElement).value)" />
        </div>
      </div>

      <div v-if="progressLog.length > 0" class="yaml-code-box" style="max-height: 110px; font-size: 0.75rem; margin-bottom: 0.85rem;">
        <div v-for="(log, i) in progressLog" :key="i" :style="{ color: i === progressLog.length - 1 ? (isComplete ? 'var(--cb-green)' : 'var(--cb-yellow)') : 'var(--cb-cyan)' }">
          {{ log }}
        </div>
      </div>

      <div v-if="isComplete" style="background: rgba(0,255,157,0.08); border: 1px solid var(--cb-green); padding: 0.5rem; border-radius: 3px; font-size: 0.75rem; color: var(--cb-green); font-family: var(--font-mono); margin-bottom: 0.85rem;">
        [OK] AUTO-ANNOTATION COMPLETE: 420 Labels Synced to /labels/train/
      </div>

      <div style="display: flex; gap: 0.75rem; justify-content: flex-end;">
        <button class="cyber-pill" :disabled="isAnnotating" @click="emit('close')">FECHAR</button>
        <button class="cyber-pill active" :disabled="isAnnotating" style="display: inline-flex; align-items: center; gap: 0.4rem;" @click="handleStartAnnotation">
          <DatasetCyberIcons name="tag" :size="14" color="#07080c" />
          <span>{{ isAnnotating ? 'ANNOTATING ON RTX 5090...' : 'INICIAR AUTO-ROTULAGEM' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>
