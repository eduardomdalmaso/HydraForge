/* UNIFIED UNIQUE YOLO / COCO 80-CLASS ONTOLOGY (TypeScript) */
export interface YoloClassDef {
  id: string
  label: string
  en: string
}

export const YOLO_CLASSES: YoloClassDef[] = [
  { id: 'carro', label: 'carro', en: 'car' },
  { id: 'moto', label: 'moto', en: 'motorcycle' },
  { id: 'caminhao', label: 'caminhao', en: 'truck' },
  { id: 'onibus', label: 'onibus', en: 'bus' },
  { id: 'pessoa', label: 'pessoa', en: 'person' },
  { id: 'bicicleta', label: 'bicicleta', en: 'bicycle' },
  { id: 'veiculo', label: 'veiculo', en: 'vehicle' },
  { id: 'aviao', label: 'aviao', en: 'airplane' },
  { id: 'barco', label: 'barco', en: 'boat' },
  { id: 'trem', label: 'trem', en: 'train' },
  { id: 'semaforo', label: 'semaforo', en: 'traffic light' },
  { id: 'placa_pare', label: 'placa_pare', en: 'stop sign' },
  { id: 'parquimetro', label: 'parquimetro', en: 'parking meter' },
  { id: 'hidrante', label: 'hidrante', en: 'fire hydrant' },
  { id: 'cachorro', label: 'cachorro', en: 'dog' },
  { id: 'gato', label: 'gato', en: 'cat' },
  { id: 'passaro', label: 'passaro', en: 'bird' },
  { id: 'cavalo', label: 'cavalo', en: 'horse' },
  { id: 'vaca', label: 'vaca', en: 'cow' },
  { id: 'ovelha', label: 'ovelha', en: 'sheep' },
  { id: 'mochila', label: 'mochila', en: 'backpack' },
  { id: 'guarda_chuva', label: 'guarda_chuva', en: 'umbrella' },
  { id: 'bolsa', label: 'bolsa', en: 'handbag' },
  { id: 'mala', label: 'mala', en: 'suitcase' },
  { id: 'garrafa', label: 'garrafa', en: 'bottle' },
  { id: 'copo', label: 'copo', en: 'cup' },
  { id: 'cadeira', label: 'cadeira', en: 'chair' },
  { id: 'sofa', label: 'sofa', en: 'couch' },
  { id: 'cama', label: 'cama', en: 'bed' },
  { id: 'mesa', label: 'mesa', en: 'table' },
  { id: 'tv', label: 'tv', en: 'tv' },
  { id: 'notebook', label: 'notebook', en: 'laptop' },
  { id: 'mouse', label: 'mouse', en: 'mouse' },
  { id: 'teclado', label: 'teclado', en: 'keyboard' },
  { id: 'celular', label: 'celular', en: 'cell phone' },
  { id: 'relogio', label: 'relogio', en: 'clock' },
  { id: 'livro', label: 'livro', en: 'book' },
  { id: 'objeto', label: 'objeto', en: 'object' },
  { id: 'ignorar', label: 'ignorar', en: 'drop' }
]

export function decodeYoloClass(raw: any): string {
  if (raw === null || raw === undefined) return 'objeto'
  const num = parseInt(raw, 10)
  if (!isNaN(num) && num >= 0 && num < YOLO_CLASSES.length) {
    const item = YOLO_CLASSES[num]
    return `#${num} ${item.label} [${item.en}]`
  }
  return String(raw)
}

export function autoSuggestCategory(raw: any, contextName?: string): string {
  const str = String(raw || '').toLowerCase().trim()
  const ctx = String(contextName || '').toLowerCase().trim()
  const num = parseInt(raw, 10)

  // Contextual dataset hints when class is raw numeric from specialized datasets
  if (ctx.includes('truck') && (num === 2 || num === 0 || str === '2' || str === '0')) return 'truck'
  if ((ctx.includes('motorcycle') || ctx.includes('moto') || ctx.includes('bike')) && (!isNaN(num) || str === 'bike')) return 'motorcycle'
  if (ctx.includes('cell') && !isNaN(num)) return 'cell-phone'

  if (str.includes('cell') || str.includes('phone') || str.includes('celular')) return 'cell-phone'
  if (str.includes('truck') || str.includes('caminh') || str.includes('lorry') || str.includes('pickup') || num === 7) return 'truck'
  if (str.includes('bus') || str.includes('onibus') || str.includes('ônibus') || num === 5) return 'bus'
  if (str.includes('bike') || str.includes('moto') || str.includes('cycle') || str.includes('scooter') || num === 3) return 'motorcycle'
  if (str.includes('person') || str.includes('pedestrian') || str.includes('pessoa') || str.includes('people') || num === 0) return 'person'
  if (str.includes('bicycle') || str.includes('bicicleta') || str.includes('ciclista') || num === 1) return 'bicycle'
  if (str.includes('car') || str.includes('van') || str.includes('auto') || str.includes('carro') || str.includes('sedan') || str.includes('suv') || num === 2) return 'car'
  return str || 'object'
}

