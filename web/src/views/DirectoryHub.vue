<script setup lang="ts">
import DirectoryFolderCard from '../components/DirectoryFolderCard.vue'

const DIRECTORIES = [
  { id: 'cockpit', code: 'DIR_01', label: 'ESTUDIO DE TREINO', tag: 'LAUNCHER', desc: 'Configure arquiteturas YOLOv8, YOLO11, YOLO26, hiperparametros e execute treinos na GPU RTX 5090.', stats: 'YOLO MATRIX', status: 'ATIVO' },
  { id: 'live-hud', code: 'DIR_02', label: 'TELEMETRIA AO VIVO', tag: 'METRICS', desc: 'Monitoramento em tempo real com curvas de loss (box/cls/dfl), precisao mAP50 e telemetria da GPU.', stats: 'WS STREAMING', status: 'ONLINE' },
  { id: 'benchmarks', code: 'DIR_03', label: 'ESTUDIO DE BENCHMARKS', tag: 'SPEED_MAP', desc: 'Avaliacao de velocidade e exportacao. Compare PyTorch vs TensorRT vs ONNX com medicao de FPS.', stats: 'RUNTIMES', status: 'OTIMIZADO' },
  { id: 'datasets', code: 'DIR_04', label: 'ESTUDIO DE DATASETS', tag: 'DATA_YAML', desc: 'Valide configuracoes data.yaml, inspecione distribuicao de classes e gerencie splits de treino/val.', stats: 'COCO / YOLO', status: 'PRONTO' },
  { id: 'model-zoo', code: 'DIR_05', label: 'REPOSITORIO MODEL ZOO', tag: 'TENSORRT', desc: 'Gerencie checkpoints (best.pt vs last.pt), compare retencao de acuracia e compile engines TensorRT.', stats: 'CHECKPOINTS', status: 'EXPORT' },
  { id: 'playground', code: 'DIR_06', label: 'PLAYGROUND DE INFERENCIA', tag: 'INFERENCE', desc: 'Inferencia em imagens e streams com ajuste de confianca e IoU em tempo real via PyTorch/TensorRT.', stats: 'LIVE TEST', status: 'PRONTO' }
]

defineProps<{ gpuStats?: any }>()
const emit = defineEmits<{ (e: 'navigate', tabId: string): void }>()
</script>

<template>
  <div class="hub-container">
    <div class="hub-hero-banner">
      <div>
        <div class="hub-title">CENTRAL DE MODULOS HYDRAFORGE</div>
        <div class="hub-subtitle">SISTEMA INTEGRADO DE TREINAMENTO E INFERENCIA YOLO</div>
      </div>
      <div class="header-status-group">
        <div class="header-metric-pill">
          <span class="pill-label">GPU:</span>
          <span class="pill-val">{{ gpuStats?.model || gpuStats?.name || 'RTX 5090' }}</span>
        </div>
        <div class="header-metric-pill">
          <span class="pill-label">VRAM:</span>
          <span class="pill-val">{{ (gpuStats?.total_vram_mb || gpuStats?.vram_total_mb || 32607).toFixed(0) }} MB</span>
        </div>
      </div>
    </div>

    <div class="folder-grid">
      <DirectoryFolderCard v-for="dir in DIRECTORIES" :key="dir.id" :folder="dir" @open="(id) => emit('navigate', id)" />
    </div>
  </div>
</template>
