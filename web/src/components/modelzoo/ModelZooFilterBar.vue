<script setup lang="ts">
const FAMILIES = ['TODOS', 'CUSTOM / TREINADOS', 'YOLO26', 'YOLO11', 'YOLOv8']
const TASKS = ['TODOS', 'DETECT', 'SEGMENT', 'POSE', 'OBB']

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
  <div class="cyber-card" style="padding: 0.75rem 1rem; margin-bottom: 0.5rem;">
    <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.75rem;">
      <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
        <span style="font-family: var(--font-inter); font-size: 0.75rem; color: var(--vms-text-muted); font-weight: 600;">
          FAMILIA:
        </span>
        <div class="pills-container" style="margin: 0;">
          <button
            v-for="fam in FAMILIES"
            :key="fam"
            type="button"
            class="cyber-pill"
            :class="{ active: (selectedFamily === 'ALL' && fam === 'TODOS') || selectedFamily === fam }"
            style="padding: 0.2rem 0.55rem; font-size: 0.72rem;"
            @click="emit('update:selectedFamily', fam === 'TODOS' ? 'ALL' : fam)"
          >
            {{ fam }}
          </button>
        </div>
      </div>

      <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
        <span style="font-family: var(--font-inter); font-size: 0.75rem; color: var(--vms-text-muted); font-weight: 600;">
          TAREFA:
        </span>
        <div class="pills-container" style="margin: 0;">
          <button
            v-for="tsk in TASKS"
            :key="tsk"
            type="button"
            class="cyber-pill"
            :class="{ active: (selectedTask === 'ALL' && tsk === 'TODOS') || selectedTask === tsk }"
            style="padding: 0.2rem 0.55rem; font-size: 0.72rem;"
            @click="emit('update:selectedTask', tsk === 'TODOS' ? 'ALL' : tsk)"
          >
            {{ tsk }}
          </button>
        </div>
      </div>

      <span class="badge-cyan" style="font-size: 0.75rem;">
        {{ totalCount }} MODELOS
      </span>
    </div>
  </div>
</template>
