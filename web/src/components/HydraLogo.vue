<script setup lang="ts">
import logoImg from '../assets/hydra-logo.jpg'

withDefaults(defineProps<{
  size?: number
  glow?: boolean
}>(), {
  size: 34,
  glow: false
})
</script>

<template>
  <div 
    class="hydra-logo-wrapper" 
    :style="{ 
      display: 'inline-flex', 
      alignItems: 'center', 
      justifyContent: 'center',
      position: 'relative',
      flexShrink: 0
    }"
  >
    <img 
      :src="logoImg" 
      alt="HydraForge Logo" 
      :style="{
        width: `${size}px`,
        height: `${size}px`,
        objectFit: 'cover',
        borderRadius: '6px',
        border: '1px solid rgba(255, 255, 255, 0.12)',
        boxShadow: glow ? '0 2px 8px rgba(255, 94, 58, 0.25)' : 'none'
      }"
    />
  </div>
</template>
