<script setup lang="ts">
defineProps<{
  models: any[]
  selectedModel?: any
}>()

const emit = defineEmits<{
  (e: 'selectModel', m: any): void
}>()

const formatMap = (val: any) => {
  if (val === undefined || val === null || val === '') return '-'
  const num = Number(val)
  return isNaN(num) ? val : num.toFixed(2)
}
</script>

<template>
  <div class="model-cards-grid">
    <div
      v-for="m in models"
      :key="m.id"
      class="model-item-card"
      :class="{ selected: selectedModel?.id === m.id }"
      @click="emit('selectModel', m)"
    >
      <div class="model-header-row">
        <span class="model-name" :title="m.name">{{ m.name }}</span>
        <span class="model-task-badge">{{ m.task }}</span>
      </div>

      <div style="font-size: 0.72rem; color: var(--vms-text-muted); min-height: 32px; line-height: 1.35; overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2; -webkit-box-orient: vertical;">
        {{ m.desc }}
      </div>

      <div class="model-stat-grid">
        <div class="model-stat-box">
          <div class="model-stat-val" style="color: #ffffff;">{{ formatMap(m.map5095) }}%</div>
          <div class="model-stat-lbl">mAP 50-95</div>
        </div>
        <div class="model-stat-box">
          <div class="model-stat-val">{{ m.params }}M</div>
          <div class="model-stat-lbl">PARAMS</div>
        </div>
        <div class="model-stat-box">
          <div class="model-stat-val" style="color: var(--vms-primary);">{{ m.trtLatency }}ms</div>
          <div class="model-stat-lbl">RTX 5090</div>
        </div>
      </div>
    </div>
  </div>
</template>
