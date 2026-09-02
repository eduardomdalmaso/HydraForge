/* UNIFIED UNIQUE YOLO / COCO 80-CLASS ONTOLOGY (< 100 LINES) */
export const YOLO_CLASSES = [
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
];

export function decodeYoloClass(raw) {
  if (raw === null || raw === undefined) return 'objeto';
  const num = parseInt(raw, 10);
  if (!isNaN(num) && num >= 0 && num < YOLO_CLASSES.length) {
    const item = YOLO_CLASSES[num];
    return `#${num} ${item.label} [${item.en}]`;
  }
  return String(raw);
}

export function autoSuggestCategory(raw) {
  const str = String(raw || '').toLowerCase().trim();
  const num = parseInt(raw, 10);
  if (str.includes('cell') || str.includes('phone')) return 'cell-phone';
  if (num === 7 || str.includes('truck') || str === 'caminhao') return 'truck';
  if (num === 5 || str.includes('bus') || str === 'onibus') return 'bus';
  if (num === 3 || str.includes('bike') || str.includes('moto') || str.includes('cycle')) return 'motorcycle';
  if (num === 0 || str.includes('person') || str.includes('pedestrian') || str === 'pessoa') return 'person';
  if (num === 1 || str.includes('bicycle') || str === 'bicicleta') return 'bicycle';
  if (num === 2 || str.includes('car') || str.includes('van') || str.includes('auto') || str === 'carro') return 'car';
  return str || 'object';
}
