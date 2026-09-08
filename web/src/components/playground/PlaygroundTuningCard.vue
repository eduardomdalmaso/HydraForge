<script setup lang="ts">
const props = defineProps<{
  config: Record<string, any>
}>()

const emit = defineEmits<{
  (e: 'update:config', cfg: Record<string, any>): void
}>()

const updateField = (field: string, val: any) => {
  emit('update:config', { ...props.config, [field]: val })
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header"><span class="card-title">2. HYPER-TUNING & FILTERS</span><span class="badge-yellow">NMS & SAHI</span></div>

    <div class="selector-group">
      <div class="selector-label">
        <span>CONFIDENCE THRESHOLD</span>
        <span class="slider-val">{{ (config.conf * 100).toFixed(0) }}%</span>
      </div>
      <div class="slider-row">
        <input type="range" min="0.05" max="0.95" step="0.05" :value="config.conf" @input="(e) => updateField('conf', parseFloat((e.target as HTMLInputElement).value))" />
      </div>
    </div>

    <div class="selector-group">
      <div class="selector-label">
        <span>IOU NMS THRESHOLD</span>
        <span class="slider-val">{{ (config.iou * 100).toFixed(0) }}%</span>
      </div>
      <div class="slider-row">
        <input type="range" min="0.10" max="0.90" step="0.05" :value="config.iou" :disabled="config.nmsFree" @input="(e) => updateField('iou', parseFloat((e.target as HTMLInputElement).value))" />
      </div>
      <div style="display: flex; align-items: center; gap: 0.5rem; margin-top: 0.4rem;">
        <input id="nmsFreeCheck" type="checkbox" :checked="config.nmsFree" @change="(e) => updateField('nmsFree', (e.target as HTMLInputElement).checked)" />
        <label for="nmsFreeCheck" style="font-size: 0.75rem; color: #cbd5e1; cursor: pointer;">NMS-FREE END-TO-END OUTPUT (YOLO26)</label>
      </div>
    </div>

    <div class="selector-group" style="border-top: 1px solid rgba(0,240,255,0.15); padding-top: 0.75rem;">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <div>
          <div style="font-family: var(--font-oxanium); font-size: 0.85rem; font-weight: 700;" :style="{ color: config.sahi ? 'var(--cb-cyan)' : '#fff' }">SAHI 4K SLICED INFERENCE</div>
          <div style="font-size: 0.7rem; color: #94a3b8;">Detect tiny objects in ultra-HD</div>
        </div>
        <button type="button" class="cyber-pill" :class="{ active: config.sahi }" @click="updateField('sahi', !config.sahi)">{{ config.sahi ? 'ON' : 'OFF' }}</button>
      </div>
    </div>

    <div class="selector-group" style="border-top: 1px solid rgba(0,240,255,0.15); padding-top: 0.75rem; margin-bottom: 0;">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <div>
          <div style="font-family: var(--font-oxanium); font-size: 0.85rem; font-weight: 700;" :style="{ color: config.isolateBg ? 'var(--cb-green)' : '#fff' }">ALPHA BACKGROUND REMOVER (PNG)</div>
          <div style="font-size: 0.7rem; color: #94a3b8;">Extract transparent alpha masks</div>
        </div>
        <button type="button" class="cyber-pill" :class="{ active: config.isolateBg }" @click="updateField('isolateBg', !config.isolateBg)">{{ config.isolateBg ? 'ON' : 'OFF' }}</button>
      </div>
    </div>
  </div>
</template>
