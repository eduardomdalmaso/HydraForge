<script setup lang="ts">
withDefaults(defineProps<{
  name: string
  color?: string
  size?: number
}>(), {
  color: 'currentColor',
  size: 16
})
</script>

<template>
  <svg v-if="name === 'zip'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M4 4h16v16H4z" />
    <path d="M10 4v4h4V4M10 8h4v4h-4M10 12v4h4v-4M12 16v2" />
  </svg>
  <svg v-else-if="name === 'tag'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M12 2l10 10-9 9L3 11V2h9z" />
    <circle cx="7.5" cy="6.5" r="1.5" :fill="color" />
  </svg>
  <svg v-else-if="name === 'fusion'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M4 6h6l4 6h6M4 18h6l4-6" />
    <circle cx="20" cy="12" r="2" :fill="color" />
    <circle cx="4" cy="6" r="2" :fill="color" />
    <circle cx="4" cy="18" r="2" :fill="color" />
  </svg>
  <svg v-else-if="name === 'kfold'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M12 2v20M2 12h20M4.93 4.93l14.14 14.14M4.93 19.07L19.07 4.93" stroke-dasharray="3 3" />
    <circle cx="12" cy="12" r="9" />
  </svg>
  <svg v-else-if="name === 'shield'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M12 2L3 6v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V6l-9-4z" />
    <path d="M9 12l2 2 4-4" />
  </svg>
  <svg v-else-if="name === 'launch'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" :fill="color" fill-opacity="0.2" />
  </svg>
  <svg v-else-if="name === 'copy'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <rect x="8" y="8" width="13" height="13" rx="1" />
    <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
  </svg>
  <svg v-else :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <polygon points="12 2 22 12 12 22 2 12" />
  </svg>
</template>
