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
  return `# Ultralytics Architecture (${props.model.id}.yaml)\nscales:\n  ${props.model.id[6] || 'n'}: [${props.model.depth || 0.33}, ${props.model.width || 0.25}, 1024]\n\nbackbone:\n  - [-1, 1, Conv, [64, 3, 2]]\n  - [-1, 1, Conv, [128, 3, 2]]\n  - [-1, 2, C3k2, [256, False, 0.25]]\nhead:\n  - [-1, 1, Detect, [nc, [256, 512, 1024]]]`
})
</script>

<template>
  <div v-if="model" class="cyber-card" style="min-width: 0; overflow: hidden;">
    <div class="card-header" style="gap: 0.5rem; min-width: 0;">
      <span class="card-title" style="min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="model.name">
        {{ model.name }}
      </span>
      <span class="badge-cyan" style="flex-shrink: 0;">{{ model.task }}</span>
    </div>

    <div style="margin-bottom: 0.85rem; min-width: 0;">
      <div class="telemetry-row" style="gap: 0.5rem;">
        <span class="k" style="flex-shrink: 0;">ARCHITECTURE HEAD</span>
        <span class="v" style="color: var(--vms-primary); text-align: right; word-break: break-word;">
          {{ model.nmsFree ? '[NMS-FREE] End-to-End' : 'Standard Decoupled' }}
        </span>
      </div>
      <div class="telemetry-row" style="gap: 0.5rem;">
        <span class="k" style="flex-shrink: 0;">FLOPs / COMPLEXITY</span>
        <span class="v">{{ model.flops }} GFLOPs</span>
      </div>
      <div class="telemetry-row" style="gap: 0.5rem;">
        <span class="k" style="flex-shrink: 0;">TENSORRT LATENCY</span>
        <span class="v" style="color: var(--vms-primary);">{{ model.trtLatency }} ms</span>
      </div>
      <div class="telemetry-row" style="border-bottom: none; gap: 0.5rem;">
        <span class="k" style="flex-shrink: 0;">CHECKPOINT</span>
        <span class="v" style="color: #ffffff; text-align: right; word-break: break-all; max-width: 60%; font-size: 0.75rem;" :title="`${model.id}.pt`">
          {{ model.id }}.pt
        </span>
      </div>
    </div>

    <pre class="yaml-code-box" style="max-height: 140px; font-size: 0.72rem; margin-bottom: 0.85rem; overflow-x: auto; white-space: pre-wrap; word-break: break-word;">
      <code>{{ yamlConfig }}</code>
    </pre>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem; margin-bottom: 0.5rem;">
      <button class="cyber-pill active" style="padding: 0.55rem; text-align: center;" @click="emit('sendToCockpit', model)">
        [TREINAR]
      </button>
      <button class="cyber-pill" style="padding: 0.55rem; text-align: center;" @click="emit('sendToBenchmark', model)">
        [BENCHMARK]
      </button>
    </div>

    <button
      class="cyber-pill"
      style="width: 100%; padding: 0.55rem; text-align: center; border-color: var(--vms-primary); color: var(--vms-primary);"
      @click="emit('sendToPlayground', model)"
    >
      [TESTAR NO PLAYGROUND] ➔
    </button>
  </div>
</template>
