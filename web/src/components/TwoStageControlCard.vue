<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ params: Record<string, any> }>()
const emit = defineEmits<{ (e: 'update:params', val: Record<string, any>): void }>()
const isTwoStage = computed(() => !!props.params.two_stage)

const handleToggle = () => {
  emit('update:params', {
    ...props.params,
    two_stage: !props.params.two_stage,
    stage1_epochs: props.params.stage1_epochs || 20,
    stage1_freeze: props.params.stage1_freeze || 10,
    stage2_epochs: props.params.stage2_epochs || 30,
    close_mosaic: props.params.close_mosaic || 10
  })
}

const applyPreset = (presetKey: string) => {
  if (presetKey === 'yolo26') {
    emit('update:params', { ...props.params, recipe_preset: 'yolo26_recipe', close_mosaic: 10, optimizer: 'AdamW', lr0: 0.00038, two_stage: false })
  } else if (presetKey === 'small_data') {
    emit('update:params', { ...props.params, recipe_preset: 'small_dataset', freeze: 10, patience: 20, epochs: 50, lr0: 0.001, two_stage: false })
  } else if (presetKey === 'two_stage') {
    emit('update:params', { ...props.params, recipe_preset: 'two_stage', two_stage: true, stage1_epochs: 20, stage1_freeze: 10, stage2_epochs: 30, close_mosaic: 10 })
  } else {
    emit('update:params', { ...props.params, recipe_preset: 'default', two_stage: false, freeze: 0 })
  }
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">4. RECEITAS & TWO-STAGE TUNING</span>
      <span class="badge-green">TRANSFER LEARNING</span>
    </div>

    <div class="selector-group">
      <div class="selector-label">PRESETS DE RECEITA</div>
      <div class="pills-container">
        <button type="button" class="cyber-pill" :class="{ active: !isTwoStage && !params.freeze }" @click="applyPreset('default')">[PADRAO]</button>
        <button type="button" class="cyber-pill" :class="{ active: params.recipe_preset === 'yolo26_recipe' }" @click="applyPreset('yolo26')">[RECIPE: YOLO26]</button>
        <button type="button" class="cyber-pill" :class="{ active: params.freeze === 10 }" @click="applyPreset('small_data')">[FREEZE: 10]</button>
        <button type="button" class="cyber-pill" :class="{ active: isTwoStage }" @click="applyPreset('two_stage')">[TWO-STAGE]</button>
      </div>
    </div>

    <div style="display: flex; align-items: center; justify-content: space-between; background: var(--vms-bg-elevated); padding: 0.65rem 0.8rem; border-radius: var(--vms-radius-sm); border: 1px solid var(--vms-border);">
      <div>
        <div style="font-family: var(--font-inter); font-size: 0.8125rem; font-weight: 600; color: #ffffff;">[PIPELINE] TWO-STAGE FINE-TUNING</div>
        <div style="font-size: 0.72rem; color: var(--vms-text-muted);">Estagio 1 (Congela Backbone) // Estagio 2 (Descongela Global)</div>
      </div>
      <button type="button" class="cyber-pill" :class="{ active: isTwoStage }" style="min-width: 80px; text-align: center;" @click="handleToggle">{{ isTwoStage ? 'ATIVO' : 'ATIVAR' }}</button>
    </div>

    <div v-if="isTwoStage" style="margin-top: 0.75rem; display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; background: var(--vms-bg-elevated); padding: 0.75rem; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-sm);">
      <div>
        <div class="selector-label" style="font-size: 0.75rem;"><span>ESTAGIO 1: HEAD ADAPT</span><span class="slider-val">{{ params.stage1_epochs || 20 }} ep</span></div>
        <input type="range" min="5" max="50" :value="params.stage1_epochs || 20" style="width: 100%;" @input="(e) => emit('update:params', { ...params, stage1_epochs: parseInt((e.target as HTMLInputElement).value) })" />
        <div style="margin-top: 0.2rem; font-size: 0.7rem; color: var(--vms-primary);">Freeze: <strong>Backbone (0-10)</strong></div>
      </div>
      <div>
        <div class="selector-label" style="font-size: 0.75rem;"><span>ESTAGIO 2: FULL REFINE</span><span class="slider-val">{{ params.stage2_epochs || 30 }} ep</span></div>
        <input type="range" min="10" max="100" :value="params.stage2_epochs || 30" style="width: 100%;" @input="(e) => emit('update:params', { ...params, stage2_epochs: parseInt((e.target as HTMLInputElement).value) })" />
        <div style="margin-top: 0.2rem; font-size: 0.7rem; color: var(--vms-success);">Unfreeze: <strong>All // lr0=0.001</strong></div>
      </div>
    </div>
  </div>
</template>
