<script setup lang="ts">
import { computed } from 'vue'

const FAMILIES = ['yolov8', 'yolo11', 'yolo26']
const SCALES = [
  { id: 'n', label: 'Nano (n)' },
  { id: 's', label: 'Small (s)' },
  { id: 'm', label: 'Medium (m)' },
  { id: 'l', label: 'Large (l)' },
  { id: 'x', label: 'XLarge (x)' }
]
const TASKS = [
  { id: 'detect', label: '[DETECT]' },
  { id: 'segment', label: '[SEGMENT]' },
  { id: 'pose', label: '[POSE]' },
  { id: 'classify', label: '[CLASSIFY]' },
  { id: 'obb', label: '[OBB]' }
]

const props = defineProps<{
  family: string
  scale: string
  task: string
}>()

const emit = defineEmits<{
  (e: 'update:family', val: string): void
  (e: 'update:scale', val: string): void
  (e: 'update:task', val: string): void
}>()

const modelTag = computed(() => {
  return `${props.family}${props.scale}${props.task !== 'detect' ? '-' + props.task : ''}.pt`
})
</script>

<template>
  <div class="cyber-card" style="position: relative;">
    <div class="hud-corner-tl" />
    <div class="hud-corner-br" />
    <div class="card-header">
      <span class="card-title">1. MODEL ARCHITECTURE MATRIX</span>
      <span class="badge-cyan">{{ modelTag }}</span>
    </div>

    <div class="selector-group">
      <div class="selector-label">MODEL FAMILY</div>
      <div class="pills-container">
        <button
          v-for="f in FAMILIES"
          :key="f"
          class="cyber-pill"
          :class="{ active: family === f }"
          @click="emit('update:family', f)"
        >
          {{ f.toUpperCase() }}
        </button>
      </div>
    </div>

    <div class="selector-group">
      <div class="selector-label">SCALE & COMPLEXITY</div>
      <div class="pills-container">
        <button
          v-for="s in SCALES"
          :key="s.id"
          class="cyber-pill"
          :class="{ active: scale === s.id }"
          @click="emit('update:scale', s.id)"
        >
          {{ s.label }}
        </button>
      </div>
    </div>

    <div class="selector-group" style="margin-bottom: 0;">
      <div class="selector-label">VISION TASK</div>
      <div class="pills-container">
        <button
          v-for="t in TASKS"
          :key="t.id"
          class="cyber-pill"
          :class="{ active: task === t.id }"
          @click="emit('update:task', t.id)"
        >
          {{ t.label }}
        </button>
      </div>
    </div>
  </div>
</template>
