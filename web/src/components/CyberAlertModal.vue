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
  title: 'AVISO DO SISTEMA',
  type: 'danger',
  confirmText: 'CONFIRMAR',
  cancelText: 'CANCELAR',
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
      <div class="cyber-alert-header">
        <div class="cyber-alert-title">
          <span>{{ title }}</span>
        </div>
        <button class="pagination-btn" style="padding: 0.15rem 0.45rem; font-size: 0.75rem;" @click="emit('close')">
          ✕
        </button>
      </div>

      <div class="cyber-alert-body">
        {{ message }}
      </div>

      <div
        v-if="checkboxLabel"
        style="margin: 0.75rem 0 1.25rem; padding: 0.5rem 0.65rem; background: rgba(255, 255, 255, 0.04); border: 1px solid var(--vms-border); border-radius: var(--vms-radius-sm);"
      >
        <label style="display: flex; align-items: center; gap: 0.5rem; cursor: pointer; font-size: 0.78rem; color: #fff; font-family: var(--font-inter);">
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
          :class="{ danger: type === 'danger' }"
          :disabled="isProcessing"
          @click="emit('confirm')"
        >
          {{ isProcessing ? 'PROCESSANDO...' : confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>
