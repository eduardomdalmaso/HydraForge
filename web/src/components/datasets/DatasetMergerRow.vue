<script setup lang="ts">
defineProps<{
  dataset: any
  isSelected: boolean
  mappings: Record<string, any>
  targetOptions: string[]
}>()

const emit = defineEmits<{
  (e: 'toggle', id: string): void
  (e: 'setMap', dsId: string, cls: string, val: string): void
}>()
</script>

<template>
  <div
    :style="{
      background: isSelected ? 'rgba(0,240,255,0.05)' : 'rgba(0,0,0,0.2)',
      border: isSelected ? '1px solid rgba(0,240,255,0.3)' : '1px solid rgba(255,255,255,0.08)',
      padding: '0.5rem',
      borderRadius: '4px',
      marginBottom: '0.5rem',
      opacity: isSelected ? 1 : 0.5
    }"
  >
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.35rem;">
      <label style="display: flex; align-items: center; gap: 0.4rem; cursor: pointer;">
        <input type="checkbox" :checked="isSelected" @change="emit('toggle', dataset.id || dataset.dataset_id)" />
        <span style="font-size: 0.78rem; font-family: var(--font-oxanium); font-weight: 700;" :style="{ color: isSelected ? 'var(--cb-yellow)' : '#94a3b8' }">
          [REPO] {{ dataset.name }} ({{ ((dataset.train_count || dataset.train_images || 0) + (dataset.val_count || dataset.val_images || 0)).toLocaleString() }} imgs)
        </span>
      </label>
      <span style="font-size: 0.65rem; color: #64748b;">{{ isSelected ? 'INCLUDED' : 'EXCLUDED' }}</span>
    </div>

    <div v-if="isSelected" style="display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 0.35rem;">
      <div v-for="cls in (dataset.classes || [])" :key="cls" class="class-mapper-row" style="height: 28px; padding: 0.15rem 0.4rem;">
        <span class="class-mapper-label" style="font-size: 0.65rem; max-width: 110px;">"{{ cls }}" ➔</span>
        <select class="cyber-select" style="width: 100px; padding: 0.1rem 0.2rem; font-size: 0.65rem;" :value="mappings[dataset.id || dataset.dataset_id]?.[cls] || ''" @change="(e) => emit('setMap', dataset.id || dataset.dataset_id, cls, (e.target as HTMLSelectElement).value)">
          <option v-for="t in targetOptions" :key="t" :value="t">{{ t }}</option>
        </select>
      </div>
    </div>
  </div>
</template>
