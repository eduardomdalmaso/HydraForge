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
  <svg v-if="name === 'pipeline' || name === 'fusion'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
    <!-- Data Source Cylinder 1 -->
    <ellipse cx="5" cy="5.5" rx="3.5" ry="1.8" />
    <path d="M1.5 5.5v3.5c0 .99 1.57 1.8 3.5 1.8s3.5-.81 3.5-1.8V5.5" />
    <!-- Data Source Cylinder 2 -->
    <ellipse cx="5" cy="15.5" rx="3.5" ry="1.8" />
    <path d="M1.5 15.5v3c0 .99 1.57 1.8 3.5 1.8s3.5-.81 3.5-1.8v-3" />
    <!-- Pipeline Flow Manifold -->
    <path d="M8.5 7.3h3.5a2 2 0 0 1 2 2v1.2" />
    <path d="M8.5 17h3.5a2 2 0 0 0 2-2v-1.2" />
    <path d="M14 12h6" />
    <polygon points="20 9.5 23.5 12 20 14.5" :fill="color" />
  </svg>
  <svg v-else-if="name === 'zip'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M4 4h16v16H4z" />
    <path d="M10 4v4h4V4M10 8h4v4h-4M10 12v4h4v-4M12 16v2" />
  </svg>
  <svg v-else-if="name === 'tag'" :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <path d="M12 2l10 10-9 9L3 11V2h9z" />
    <circle cx="7.5" cy="6.5" r="1.5" :fill="color" />
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
  <svg v-else :width="size" :height="size" viewBox="0 0 24 24" fill="none" :stroke="color" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter">
    <polygon points="12 2 22 12 12 22 2 12" />
  </svg>
</template>
