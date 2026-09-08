<script setup lang="ts">
withDefaults(defineProps<{
  isOpen: boolean
  title?: string
  message: string
  type?: 'danger' | 'cyan' | 'yellow'
  confirmText?: string
  cancelText?: string
  checkboxLabel?: string
  checkboxChecked?: boolean
  isProcessing?: boolean
}>(), {
  title: 'SYSTEM ALERT',
  type: 'danger',
  confirmText: 'CONFIRM',
  cancelText: 'ABORT',
  checkboxChecked: false,
  isProcessing: false
})

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'close'): void
  (e: 'toggleCheckbox', checked: boolean): void
}>()
</script>

<template>
  <div v-if="isOpen" class="cyber-alert-overlay" @click="emit('close')">
    <div
      class="cyber-alert-box"
      :class="{ cyan: type === 'cyan', yellow: type === 'yellow' }"
      @click.stop
    >
      <div class="hud-corner-tl" />
      <div class="hud-corner-br" />

      <div class="cyber-alert-header">
        <div class="cyber-alert-title">
          <span class="cyber-hud-tag">
            {{ type === 'danger' ? '[ALERT]' : (type === 'yellow' ? '[WARN]' : '[INFO]') }}
          </span>
          <span>{{ title }}</span>
        </div>
        <button class="cyber-pill" style="padding: 0.15rem 0.45rem; font-size: 0.7rem;" @click="emit('close')">
          ✕
        </button>
      </div>

      <div class="cyber-alert-body">
        {{ message }}
      </div>

      <div
        v-if="checkboxLabel"
        style="margin: 0.75rem 0 1.25rem; padding: 0.5rem 0.65rem; background: rgba(255, 0, 85, 0.08); border: 1px dashed rgba(255, 0, 85, 0.3); border-radius: 3px;"
      >
        <label style="display: flex; align-items: center; gap: 0.5rem; cursor: pointer; font-size: 0.78rem; color: #fff; font-family: var(--font-mono);">
          <input
            type="checkbox"
            :checked="checkboxChecked"
            @change="(e) => emit('toggleCheckbox', (e.target as HTMLInputElement).checked)"
          />
          <span>{{ checkboxLabel }}</span>
        </label>
      </div>

      <div class="cyber-alert-actions">
        <button class="cyber-action-btn secondary" :disabled="isProcessing" @click="emit('close')">
          {{ cancelText }}
        </button>
        <button
          class="cyber-action-btn"
          :style="type === 'danger' ? { background: 'var(--cb-magenta)', borderColor: 'var(--cb-magenta)', color: '#fff' } : {}"
          :disabled="isProcessing"
          @click="emit('confirm')"
        >
          {{ isProcessing ? 'PROCESSING...' : confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>
