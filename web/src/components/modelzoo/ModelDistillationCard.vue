<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'launchDistill', payload: { teacher: string; student: string; disLoss: number }): void
}>()

const teacher = ref('yolo26x')
const student = ref('yolo26n')
const disLoss = ref(1.5)
</script>

<template>
  <div class="cyber-card" style="margin-top: 1.25rem;">
    <div class="card-header">
      <span class="card-title">DESTILACAO DE CONHECIMENTO</span>
      <span class="badge-yellow">TEACHER ➔ STUDENT</span>
    </div>

    <div style="font-size: 0.75rem; color: var(--vms-text-muted); margin-bottom: 0.65rem;">
      Transfira acuracia do modelo <strong>Teacher ({{ teacher }})</strong> para o <strong>Student ({{ student }})</strong>, elevando mAP sem aumentar latencia.
    </div>

    <div class="distill-connector-box">
      <div style="text-align: center;">
        <div style="font-size: 0.65rem; color: var(--vms-text-muted); font-family: var(--font-mono);">TEACHER</div>
        <div style="font-family: var(--font-inter); font-weight: 700; color: #fff; font-size: 0.9rem;">{{ teacher }}</div>
        <div style="font-size: 0.65rem; color: var(--vms-text-muted);">56.1% mAP</div>
      </div>

      <div style="color: var(--vms-primary); font-size: 0.9rem; font-weight: 700;">➔ DESTILAR ➔</div>

      <div style="text-align: center;">
        <div style="font-size: 0.65rem; color: var(--vms-text-muted); font-family: var(--font-mono);">STUDENT</div>
        <div style="font-family: var(--font-inter); font-weight: 700; color: #fff; font-size: 0.9rem;">{{ student }}</div>
        <div style="font-size: 0.65rem; color: var(--vms-text-muted);">0.38ms Latencia</div>
      </div>
    </div>

    <div class="selector-group" style="margin-bottom: 0.75rem;">
      <div class="selector-label">
        <span>PESO DO LOSS DE DESTILACAO (dis)</span>
        <span class="slider-val">{{ disLoss.toFixed(1) }}</span>
      </div>
      <div class="slider-row">
        <input
          type="range"
          min="0.5"
          max="3.0"
          step="0.1"
          :value="disLoss"
          @input="(e) => disLoss = parseFloat((e.target as HTMLInputElement).value)"
        />
      </div>
    </div>

    <button
      class="cyber-action-btn"
      style="width: 100%;"
      @click="emit('launchDistill', { teacher, student, disLoss })"
    >
      EXECUTAR DESTILACAO NA RTX 5090
    </button>
  </div>
</template>
