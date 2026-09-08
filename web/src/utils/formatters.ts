export function formatLoss(val?: number | null): string {
  if (val == null || isNaN(val)) return '-'
  if (val < 0.001 && val > 0) return val.toExponential(2)
  return val.toFixed(4)
}

export function formatNumber(num?: number | null): string {
  if (num == null || isNaN(num)) return '0'
  return num.toLocaleString()
}
