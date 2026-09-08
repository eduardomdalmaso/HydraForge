<script setup lang="ts">
import type { TrainingParams } from '../types/training'

const props = defineProps<{ params: TrainingParams }>()
const emit = defineEmits<{ (e: 'update:params', val: TrainingParams): void }>()
const updateField = <K extends keyof TrainingParams>(field: K, val: TrainingParams[K]) => {
  emit('update:params', { ...props.params, [field]: val })
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">2. HYPERPARAMETER COCKPIT</span>
      <span class="badge-yellow">PYTORCH CUDA</span>
    </div>

    <div class="selector-group">
      <div class="selector-label">
        <span>TRAINING EPOCHS & EARLY STOPPING</span>
        <span class="slider-val">{{ params.epochs }} EP (Patience: {{ params.patience || 20 }})</span>
      </div>
      <div class="slider-row">
        <input type="range" min="5" max="300" :value="params.epochs" @input="(e) => updateField('epochs', parseInt((e.target as HTMLInputElement).value))" />
      </div>
    </div>

    <div class="selector-group">
      <div class="selector-label">BATCH SIZE & RESOLUTION</div>
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.65rem;">
        <div>
          <label style="font-size: 0.72rem; color: #94a3b8;">BATCH SIZE</label>
          <select class="cyber-select" :value="params.batch" @change="(e) => updateField('batch', parseInt((e.target as HTMLSelectElement).value))">
            <option value="-1">[AUTO] Auto-Batch (-1)</option>
            <option value="16">[BATCH: 16] 16 Images</option>
            <option value="32">[BATCH: 32] 32 Images</option>
            <option value="64">[BATCH: 64] 64 Images (RTX 5090)</option>
          </select>
        </div>
        <div>
          <label style="font-size: 0.72rem; color: #94a3b8;">IMAGE SIZE (IMGSZ)</label>
          <select class="cyber-select" :value="params.imgsz" @change="(e) => updateField('imgsz', parseInt((e.target as HTMLSelectElement).value))">
            <option value="640">640 x 640 (Standard)</option>
            <option value="1280">1280 x 1280 (HD)</option>
          </select>
        </div>
      </div>
    </div>

    <div class="selector-group" style="margin-bottom: 0;">
      <div class="selector-label">OPTIMIZER & EARLY STOPPING (PATIENCE)</div>
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.65rem;">
        <div>
          <label style="font-size: 0.72rem; color: #94a3b8;">OPTIMIZER</label>
          <select class="cyber-select" :value="params.optimizer" @change="(e) => updateField('optimizer', (e.target as HTMLSelectElement).value)">
            <option value="AdamW">[OPT] AdamW</option>
            <option value="SGD">[OPT] SGD with Momentum</option>
          </select>
        </div>
        <div>
          <label style="font-size: 0.72rem; color: #94a3b8;">EARLY STOPPING</label>
          <select class="cyber-select" :value="params.patience || 20" @change="(e) => updateField('patience', parseInt((e.target as HTMLSelectElement).value))">
            <option value="15">[PATIENCE] 15 Epochs</option>
            <option value="20">[PATIENCE] 20 Epochs</option>
            <option value="30">[PATIENCE] 30 Epochs</option>
            <option value="50">[PATIENCE] 50 Epochs</option>
            <option value="0">[OFF] Desativado (0)</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>
