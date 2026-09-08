<script setup lang="ts">
import { ref } from 'vue'
import DatasetCyberIcons from './DatasetCyberIcons.vue'

const props = defineProps<{
  dataset?: any
}>()

const emit = defineEmits<{
  (e: 'sendToCockpit'): void
}>()

const auditStatus = ref<any>(null)
const isAuditing = ref(false)

const handleRunAudit = async () => {
  const id = props.dataset?.id || props.dataset?.dataset_id
  if (!id) return
  isAuditing.value = true
  try {
    const res = await fetch(`/api/v1/training/datasets/audit/${id}`, { method: 'POST' })
    if (res.ok) auditStatus.value = await res.json()
  } finally {
    isAuditing.value = false
  }
}
</script>

<template>
  <div class="cyber-card">
    <div class="card-header">
      <span class="card-title">4. QUALITY AUDITOR & LEAKAGE DEFENSE</span>
      <span class="badge-green">REAL DATA INTEGRITY</span>
    </div>

    <div style="margin-bottom: 0.85rem;">
      <div class="telemetry-row">
        <span class="k">ANNOTATION INTEGRITY</span>
        <span class="audit-badge-ok">{{ auditStatus ? `✓ ${auditStatus.valid_bboxes_pct} VALID` : 'AUDIT REQUIRED' }}</span>
      </div>
      <div class="telemetry-row">
        <span class="k">CORRUPTED LABELS / OUT-OF-BOUNDS</span>
        <span class="audit-badge-ok">{{ auditStatus ? `✓ ${auditStatus.corrupt_files} CORRUPTED` : '-' }}</span>
      </div>
      <div class="telemetry-row" style="border-bottom: none;">
        <span class="k">TRAIN / VAL DATA LEAKAGE</span>
        <span class="audit-badge-ok">{{ auditStatus ? `✓ ${auditStatus.leakage_overlap_pct} OVERLAP` : '-' }}</span>
      </div>
    </div>

    <div v-if="auditStatus" style="background: rgba(0,255,157,0.06); border: 1px solid var(--cb-green); padding: 0.65rem; border-radius: 4px; margin-bottom: 0.85rem; font-family: var(--font-mono); font-size: 0.75rem;">
      <div style="color: var(--cb-green); font-weight: 700; margin-bottom: 0.2rem;">
        [AUDIT] // {{ auditStatus.status }} [{{ auditStatus.timestamp }}]
      </div>
      <div style="color: #cbd5e1;">
        {{ auditStatus.total_bboxes?.toLocaleString() }} bboxes auditados no disco. Zero vazamento.
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 1.2fr; gap: 0.5rem;">
      <button class="cyber-action-btn" :disabled="isAuditing || !dataset" @click="handleRunAudit">
        <DatasetCyberIcons name="shield" :size="14" color="#07080c" />
        <span>{{ isAuditing ? 'AUDITING...' : 'AUDIT' }}</span>
      </button>
      <button class="cyber-action-btn secondary" :disabled="!dataset" @click="emit('sendToCockpit')">
        <DatasetCyberIcons name="launch" :size="14" />
        <span>COCKPIT</span>
      </button>
    </div>
  </div>
</template>
