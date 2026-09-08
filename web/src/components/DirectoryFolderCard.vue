<script setup lang="ts">
import CyberNavIcons from './CyberNavIcons.vue'

defineProps<{
  folder: {
    id: string
    label: string
    code: string
    tag: string
    desc: string
    stats: string
    status?: string
  }
}>()

const emit = defineEmits<{
  (e: 'open', id: string): void
}>()
</script>

<template>
  <div class="cyber-folder-card" @click="emit('open', folder.id)">
    <div class="folder-tab-ear">
      // {{ folder.code }} [{{ folder.tag }}]
    </div>
    
    <div class="folder-header">
      <div class="folder-title-group">
        <div class="folder-icon cyber-icon-glow">
          <CyberNavIcons :type="folder.id" />
        </div>
        <div>
          <div class="folder-title">{{ folder.label }}</div>
          <div class="text-mono" style="font-size: 0.7rem; color: var(--cb-yellow);">
            PATH: /{{ folder.id }}
          </div>
        </div>
      </div>
      <span class="card-badge">{{ folder.status || 'READY' }}</span>
    </div>

    <div class="folder-desc">
      {{ folder.desc }}
    </div>

    <div class="folder-footer">
      <span style="color: #64748b;">
        METRIC: <strong style="color: var(--cb-green);">{{ folder.stats }}</strong>
      </span>
      <span class="folder-action-btn">
        MOUNT DIRECTORY ❯❯
      </span>
    </div>
  </div>
</template>
