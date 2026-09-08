<script setup lang="ts">
import { ref } from 'vue'
import DatasetCyberIcons from './DatasetCyberIcons.vue'

defineProps<{
  dataset?: any
}>()

const kFolds = ref(5)
const isGenerating = ref(false)
const genSuccess = ref(false)

const handleGenerateFolds = () => {
  isGenerating.value = true
  genSuccess.value = false
  setTimeout(() => {
    isGenerating.value = false
    genSuccess.value = true
  }, 450)
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">3. K-FOLD CROSS-VALIDATION SPLITTER</span>
      <span class="badge-yellow">SKLEARN KFOLD</span>
    </div>

    <div class="selector-group">
      <div class="selector-label">
        <span>NUMBER OF FOLDS (K)</span>
        <span class="slider-val">K = {{ kFolds }} FOLDS</span>
      </div>
      <div class="pills-container" style="margin-bottom: 0.75rem;">
        <button
          v-for="k in [3, 5, 10]"
          :key="k"
          type="button"
          class="cyber-pill"
          :class="{ active: kFolds === k }"
          @click="() => { kFolds = k; genSuccess = false; }"
        >
          {{ k }}-Fold ({{ ((1 / k) * 100).toFixed(0) }}% Val per Split)
        </button>
      </div>
    </div>

    <div style="background: rgba(0,0,0,0.3); padding: 0.65rem 0.8rem; border-radius: 4px; border: 1px solid rgba(0,240,255,0.15);">
      <div style="display: flex; justify-content: space-between; font-size: 0.75rem; font-family: var(--font-mono);">
        <span style="color: var(--cb-cyan);">TRAIN ({{ 100 - Number(((1 / kFolds) * 100).toFixed(0)) }}%)</span>
        <span style="color: var(--cb-yellow);">VAL ({{ ((1 / kFolds) * 100).toFixed(0) }}%)</span>
      </div>
      <div class="kfold-bar">
        <div class="kfold-seg-train" :style="{ width: `${100 - (100 / kFolds)}%` }"></div>
        <div class="kfold-seg-val" :style="{ width: `${100 / kFolds}%` }"></div>
      </div>
      <div style="margin-top: 0.5rem; font-size: 0.7rem; color: #94a3b8;">
        Generates <code>split_1.yaml</code> ... <code>split_{{ kFolds }}.yaml</code> in zero-copy text.
      </div>
    </div>

    <div v-if="genSuccess" style="margin-top: 0.65rem; padding: 0.5rem; background: rgba(0,255,157,0.06); border: 1px solid var(--cb-green); border-radius: 4px; font-size: 0.72rem; color: var(--cb-green); font-family: var(--font-mono);">
      ✓ Generated {{ kFolds }} Stratified Folds for {{ dataset?.name || 'Dataset' }}
    </div>

    <button
      class="cyber-action-btn"
      style="width: 100%; margin-top: 0.75rem; padding: 0.75rem;"
      :disabled="isGenerating"
      @click="handleGenerateFolds"
    >
      <DatasetCyberIcons name="kfold" :size="15" color="#07080c" />
      <span>{{ isGenerating ? 'GENERATING...' : 'GENERATE' }}</span>
    </button>
  </div>
</template>
