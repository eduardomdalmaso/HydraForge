<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  model?: any
}>()

const emit = defineEmits<{
  (e: 'sendToCockpit', m: any): void
  (e: 'sendToBenchmark', m: any): void
  (e: 'sendToPlayground', m: any): void
}>()

const yamlConfig = computed(() => {
  if (!props.model) return ''
  return `# Ultralytics Model Architecture (${props.model.id}.yaml)\nscales:\n  ${props.model.id[6] || 'n'}: [${props.model.depth || 0.33}, ${props.model.width || 0.25}, 1024]\n\n# Dual-Head End-to-End NMS-Free Backbone\nbackbone:\n  - [-1, 1, Conv, [64, 3, 2]]       # 0-P1/2\n  - [-1, 1, Conv, [128, 3, 2]]      # 1-P2/4\n  - [-1, 2, C3k2, [256, False, 0.25]] # 2\n  - [-1, 1, C2PSA, [512]]           # 9 Attention\nhead:\n  - [-1, 1, Detect, [nc, [256, 512, 1024]]] # Dual One-to-One / One-to-Many`
})
</script>

<template>
  <div v-if="model" class="cyber-card">
    <div class="card-header">
      <span class="card-title">SPECIFICATION // {{ model.name }}</span>
      <span class="badge-cyan">{{ model.task }}</span>
    </div>

    <div style="margin-bottom: 0.85rem;">
      <div class="telemetry-row">
        <span class="k">ARCHITECTURE HEAD</span>
        <span class="v" style="color: var(--cb-yellow);">{{ model.nmsFree ? '[NMS-FREE] End-to-End' : 'Standard Decoupled Head' }}</span>
      </div>
      <div class="telemetry-row">
        <span class="k">FLOPs / COMPLEXITY</span>
        <span class="v">{{ model.flops }} GFLOPs</span>
      </div>
      <div class="telemetry-row">
        <span class="k">TENSORRT LATENCY (RTX 5090)</span>
        <span class="v" style="color: var(--cb-cyan);">{{ model.trtLatency }} ms</span>
      </div>
      <div class="telemetry-row" style="border-bottom: none;">
        <span class="k">WEIGHT FILE</span>
        <span class="v" style="color: var(--cb-green);">{{ model.id }}.pt (PyTorch Checkpoint)</span>
      </div>
    </div>

    <pre class="yaml-code-box" style="max-height: 140px; font-size: 0.72rem; margin-bottom: 0.85rem;">
      <code>{{ yamlConfig }}</code>
    </pre>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem; margin-bottom: 0.5rem;">
      <button class="cyber-pill active" style="padding: 0.55rem; text-align: center;" @click="emit('sendToCockpit', model)">
        [TRAIN] COCKPIT
      </button>
      <button class="cyber-pill" style="padding: 0.55rem; text-align: center;" @click="emit('sendToBenchmark', model)">
        [BENCHMARK] TRT
      </button>
    </div>

    <button
      class="cyber-pill"
      style="width: 100%; padding: 0.55rem; text-align: center; border-color: var(--cb-yellow); color: var(--cb-yellow);"
      @click="emit('sendToPlayground', model)"
    >
      [TEST] PLAYGROUND ❯❯
    </button>
  </div>
</template>
