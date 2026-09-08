<script setup lang="ts">
defineProps<{
  formats: string[]
  formatList: Array<{ id: string; label: string; gpuOnly?: boolean }>
  device: string
}>()

const emit = defineEmits<{
  (e: 'toggle', id: string, gpuOnly?: boolean): void
}>()
</script>

<template>
  <div style="margin-top: 0.75rem;">
    <div class="selector-label">TARGET RUNTIMES ({{ formats.length }} SELECTED)</div>
    <div class="format-checkbox-grid">
      <div
        v-for="f in formatList"
        :key="f.id"
        class="format-chip"
        :class="{ selected: formats.includes(f.id) }"
        :style="{ opacity: device === 'cpu' && f.gpuOnly ? 0.35 : 1, cursor: device === 'cpu' && f.gpuOnly ? 'not-allowed' : 'pointer' }"
        @click="emit('toggle', f.id, f.gpuOnly)"
      >
        <span>{{ formats.includes(f.id) ? '◈' : '◇' }}</span> <span>{{ f.label }}</span>
      </div>
    </div>
  </div>
</template>
