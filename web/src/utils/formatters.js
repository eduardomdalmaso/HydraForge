/**
 * Adaptive numerical formatters for HydraForge HUD & Telemetry.
 */

export function formatLoss(val) {
  if (val === undefined || val === null || val === '') return '-';
  const num = Number(val);
  if (isNaN(num) || num < 0) return '-';
  if (num === 0) return '0.000';
  if (num >= 10) return num.toFixed(2);
  if (num >= 1) return num.toFixed(3);
  if (num >= 0.01) return num.toFixed(4);
  if (num >= 0.0001) return num.toFixed(5);
  if (num >= 0.00001) return num.toFixed(6);
  return num.toExponential(2);
}

export function formatPercent(val, decimals = 1) {
  if (val === undefined || val === null || isNaN(val) || Number(val) <= 0) return '-';
  return `${(Number(val) * 100).toFixed(decimals)}%`;
}
