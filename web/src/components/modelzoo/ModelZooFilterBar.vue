<script setup lang="ts">
const FAMILIES = ['ALL', 'CUSTOM / TRAINED', 'YOLO26', 'YOLO11', 'YOLOv8']
const TASKS = ['ALL', 'DETECT', 'SEGMENT', 'POSE', 'OBB']

defineProps<{
  selectedFamily: string
  selectedTask: string
  totalCount: number
}>()

const emit = defineEmits<{
  (e: 'update:selectedFamily', fam: string): void
  (e: 'update:selectedTask', tsk: string): void
}>()
</script>

<template>
  <div class="cyber-card" style="padding: 0.85rem 1rem; margin-bottom: 0.5rem;">
    <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.75rem;">
      <div style="display: flex; alignItems: center; gap: 0.5rem; flex-wrap: wrap;">
        <span style="font-family: var(--font-oxanium); font-size: 0.75rem; color: var(--cb-cyan); font-weight: 700;">
          FAMILY:
        </span>
        <div class="pills-container" style="margin: 0;">
          <button
            v-for="fam in FAMILIES"
            :key="fam"
            type="button"
            class="cyber-pill"
            :class="{ active: selectedFamily === fam }"
            style="padding: 0.2rem 0.55rem; font-size: 0.72rem;"
            @click="emit('update:selectedFamily', fam)"
          >
            {{ fam }}
          </button>
        </div>
      </div>

      <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
        <span style="font-family: var(--font-oxanium); font-size: 0.75rem; color: var(--cb-yellow); font-weight: 700;">
          TASK:
        </span>
        <div class="pills-container" style="margin: 0;">
          <button
            v-for="tsk in TASKS"
            :key="tsk"
            type="button"
            class="cyber-pill"
            :class="{ active: selectedTask === tsk }"
            style="padding: 0.2rem 0.55rem; font-size: 0.72rem;"
            @click="emit('update:selectedTask', tsk)"
          >
            {{ tsk }}
          </button>
        </div>
      </div>

      <span class="badge-cyan" style="font-size: 0.75rem;">
        {{ totalCount }} CHECKPOINTS
      </span>
    </div>
  </div>
</template>
