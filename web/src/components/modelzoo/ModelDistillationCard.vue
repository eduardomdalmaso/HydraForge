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
      <span class="card-title">KNOWLEDGE DISTILLATION STUDIO</span>
      <span class="badge-yellow">TEACHER ➔ STUDENT</span>
    </div>

    <div style="font-size: 0.75rem; color: #94a3b8; margin-bottom: 0.65rem;">
      Transfer predictive power from the <strong>Teacher ({{ teacher }})</strong> model to the <strong>Student ({{ student }})</strong> model, boosting mAP with zero inference latency overhead.
    </div>

    <div class="distill-connector-box">
      <div style="text-align: center;">
        <div style="font-size: 0.65rem; color: var(--cb-yellow); font-family: var(--font-mono);">TEACHER (ORACLE)</div>
        <div style="font-family: var(--font-oxanium); font-weight: 700; color: #fff; fontSize: 0.9rem;">{{ teacher }}</div>
        <div style="font-size: 0.65rem; color: #94a3b8;">56.1% mAP</div>
      </div>

      <div style="color: var(--cb-cyan); font-size: 1.2rem; font-weight: 700;">➔ [DISTILL] ➔</div>

      <div style="text-align: center;">
        <div style="font-size: 0.65rem; color: var(--cb-cyan); font-family: var(--font-mono);">STUDENT (TARGET)</div>
        <div style="font-family: var(--font-oxanium); font-weight: 700; color: #fff; font-size: 0.9rem;">{{ student }}</div>
        <div style="font-size: 0.65rem; color: #94a3b8;">0.38ms Latency</div>
      </div>
    </div>

    <div class="selector-group" style="margin-bottom: 0.75rem;">
      <div class="selector-label">
        <span>DISTILLATION LOSS WEIGHT (dis)</span>
        <span class="slider-val">{{ disLoss.toFixed(1) }}</span>
      </div>
      <input
        type="range"
        min="0.5"
        max="3.0"
        step="0.1"
        :value="disLoss"
        @input="(e) => disLoss = parseFloat((e.target as HTMLInputElement).value)"
      />
    </div>

    <button
      class="cyber-pill active"
      style="width: 100%; padding: 0.6rem; text-align: center;"
      @click="emit('launchDistill', { teacher, student, disLoss })"
    >
      [LAUNCH] DISTILLATION ON RTX 5090
    </button>
  </div>
</template>
