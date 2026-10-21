<script setup lang="ts">
import { ref } from 'vue'
import { registerHydraStream } from '../../api/hydrastream_client'

defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'streamAdded', streamId: string): void }>()

const streamId = ref('')
const sourceUrl = ref('')
const fps = ref(25)
const isSubmitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleRegister = async () => {
  if (!streamId.value.trim() || !sourceUrl.value.trim()) {
    errorMessage.value = 'PREENCHA O ID DA CAMERA E A URL RTSP'
    return
  }
  isSubmitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const res = await registerHydraStream(streamId.value.trim().toLowerCase(), sourceUrl.value.trim(), Number(fps.value) || 25)
    successMessage.value = `[OK] FLUXO ${res.stream_id.toUpperCase()} INGESTADO COM SUCESSO NO HYDRASTREAM`
    setTimeout(() => {
      emit('streamAdded', res.stream_id)
      emit('close')
      streamId.value = ''
      sourceUrl.value = ''
      successMessage.value = ''
    }, 900)
  } catch (err: any) {
    errorMessage.value = err.message || 'FALHA NA CONEXAO RTSP'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div v-if="isOpen" class="modal-overlay" @click.self="emit('close')">
    <div class="cyber-modal">
      <div class="modal-header">
        <span class="modal-title">[INGESTAO RTSP // HYDRASTREAM]</span>
        <button class="close-btn" @click="emit('close')">X</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label>IDENTIFICADOR DO STREAM (STREAM ID)</label>
          <input v-model="streamId" class="cyber-input" placeholder="ex: cam_portaria_01" :disabled="isSubmitting" />
        </div>

        <div class="form-group">
          <label>URL RTSP / RTMP DO SENSOR</label>
          <input v-model="sourceUrl" class="cyber-input" placeholder="rtsp://usuario:senha@ip:porta/h264" :disabled="isSubmitting" />
        </div>

        <div class="form-group">
          <label>FPS DE INGESTAO ALVO</label>
          <input v-model.number="fps" type="number" min="1" max="120" class="cyber-input" :disabled="isSubmitting" />
        </div>

        <div v-if="errorMessage" class="cyber-alert error">{{ errorMessage }}</div>
        <div v-if="successMessage" class="cyber-alert success">{{ successMessage }}</div>
      </div>

      <div class="modal-footer">
        <button class="cyber-btn secondary" :disabled="isSubmitting" @click="emit('close')">CANCELAR</button>
        <button class="cyber-btn primary" :disabled="isSubmitting" @click="handleRegister">
          {{ isSubmitting ? 'INGESTANDO...' : '[CONECTAR AO HYDRASTREAM]' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay { position: fixed; inset: 0; background: rgba(7,8,12,0.88); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 999; }
.cyber-modal { width: 90%; max-width: 520px; background: #0b0e14; border: 1px solid var(--cb-cyan, #00f0ff); box-shadow: 0 0 20px rgba(0,240,255,0.2); font-family: var(--font-mono, monospace); }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 0.8rem 1rem; border-bottom: 1px solid rgba(0,240,255,0.25); background: rgba(0,240,255,0.06); }
.modal-title { font-size: 0.85rem; font-weight: bold; color: var(--cb-cyan, #00f0ff); letter-spacing: 0.05em; }
.close-btn { background: transparent; border: none; color: #ff003c; font-size: 1rem; cursor: pointer; }
.modal-body { padding: 1.2rem; display: flex; flex-direction: column; gap: 0.9rem; }
.form-group { display: flex; flex-direction: column; gap: 0.35rem; }
.form-group label { font-size: 0.68rem; color: #8892b0; letter-spacing: 0.05em; }
.cyber-input { background: #121824; border: 1px solid #1f2b3e; color: #e6f1ff; padding: 0.55rem; font-size: 0.82rem; font-family: inherit; }
.cyber-input:focus { border-color: var(--cb-cyan, #00f0ff); outline: none; }
.cyber-alert { padding: 0.5rem; font-size: 0.75rem; border-left: 3px solid; }
.cyber-alert.error { background: rgba(255,0,60,0.1); border-color: #ff003c; color: #ff003c; }
.cyber-alert.success { background: rgba(0,255,157,0.1); border-color: #00ff9d; color: #00ff9d; }
.modal-footer { display: flex; justify-content: flex-end; gap: 0.6rem; padding: 0.8rem 1rem; border-top: 1px solid rgba(31,43,62,0.6); }
.cyber-btn { padding: 0.5rem 1rem; font-size: 0.75rem; font-family: inherit; cursor: pointer; border: 1px solid; }
.cyber-btn.secondary { background: transparent; border-color: #495670; color: #a8b2d1; }
.cyber-btn.primary { background: rgba(0,240,255,0.15); border-color: var(--cb-cyan, #00f0ff); color: var(--cb-cyan, #00f0ff); font-weight: bold; }
.cyber-btn.primary:hover { background: rgba(0,240,255,0.3); }
</style>
