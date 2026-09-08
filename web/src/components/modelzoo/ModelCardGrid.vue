<script setup lang="ts">
defineProps<{
  models: any[]
  selectedModel?: any
}>()

const emit = defineEmits<{
  (e: 'selectModel', m: any): void
}>()
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
        <span class="model-name">{{ m.name }}</span>
        <span class="model-task-badge">{{ m.task }}</span>
      </div>

      <div style="font-size: 0.72rem; color: #94a3b8; min-height: 32px;">
        {{ m.desc }}
      </div>

      <div class="model-stat-grid">
        <div class="model-stat-box">
          <div class="model-stat-val" style="color: var(--cb-green);">{{ m.map5095 }}%</div>
          <div class="model-stat-lbl">mAP 50-95</div>
        </div>
        <div class="model-stat-box">
          <div class="model-stat-val">{{ m.params }}M</div>
          <div class="model-stat-lbl">PARAMS</div>
        </div>
        <div class="model-stat-box">
          <div class="model-stat-val" style="color: var(--cb-yellow);">{{ m.trtLatency }}ms</div>
          <div class="model-stat-lbl">RTX 5090</div>
        </div>
      </div>
    </div>
  </div>
</template>
