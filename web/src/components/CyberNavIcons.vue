<script setup lang="ts">
defineProps<{
  type: string
}>()
</script>

<template>
  <!-- COCKPIT -->
  <svg v-if="type === 'cockpit'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" style="display: inline-block; vertical-align: middle;">
    <polygon points="12 2 22 8.5 22 15.5 12 22 2 15.5 2 8.5" stroke="currentColor" />
    <circle cx="12" cy="12" r="3" fill="currentColor" opacity="0.8" />
    <line x1="12" y1="2" x2="12" y2="7" />
    <line x1="12" y1="17" x2="12" y2="22" />
  </svg>

  <!-- LIVE HUD -->
  <svg v-else-if="type === 'live-hud'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" style="display: inline-block; vertical-align: middle;">
    <polyline points="2 12 6 12 9 4 15 20 18 12 22 12" stroke="currentColor" />
    <circle cx="9" cy="4" r="1.5" fill="currentColor" />
    <circle cx="15" cy="20" r="1.5" fill="currentColor" />
  </svg>

  <!-- BENCHMARKS -->
  <svg v-else-if="type === 'benchmarks'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" style="display: inline-block; vertical-align: middle;">
    <path d="M12 2v4M4.93 4.93l2.83 2.83M2 12h4M4.93 19.07l2.83-2.83M19.07 4.93l-2.83 2.83M22 12h-4" />
    <polyline points="12 12 16 8" stroke="var(--cb-yellow)" stroke-width="2.5" />
    <circle cx="12" cy="12" r="2" fill="currentColor" />
  </svg>

  <!-- DATASETS -->
  <svg v-else-if="type === 'datasets'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" style="display: inline-block; vertical-align: middle;">
    <rect x="3" y="3" width="7" height="7" stroke="currentColor" />
    <rect x="14" y="3" width="7" height="7" stroke="currentColor" />
    <rect x="14" y="14" width="7" height="7" stroke="currentColor" />
    <rect x="3" y="14" width="7" height="7" stroke="currentColor" />
    <line x1="10" y1="6.5" x2="14" y2="6.5" />
    <line x1="6.5" y1="10" x2="6.5" y2="14" />
  </svg>

  <!-- MODEL ZOO -->
  <svg v-else-if="type === 'model-zoo'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" style="display: inline-block; vertical-align: middle;">
    <polygon points="12 2 22 7.5 22 16.5 12 22 2 16.5 2 7.5" />
    <polyline points="2 7.5 12 13 22 7.5" />
    <line x1="12" y1="13" x2="12" y2="22" />
  </svg>

  <!-- PLAYGROUND -->
  <svg v-else-if="type === 'playground'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" style="display: inline-block; vertical-align: middle;">
    <circle cx="12" cy="12" r="8" stroke="currentColor" stroke-dasharray="3 2" />
    <circle cx="12" cy="12" r="2" fill="currentColor" />
    <line x1="12" y1="2" x2="12" y2="6" />
    <line x1="12" y1="18" x2="12" y2="22" />
    <line x1="2" y1="12" x2="6" y2="12" />
    <line x1="18" y1="12" x2="22" y2="12" />
  </svg>

  <span v-else>◈</span>
</template>
