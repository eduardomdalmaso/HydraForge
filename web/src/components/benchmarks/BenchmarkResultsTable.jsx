import React from 'react';

export default function BenchmarkResultsTable({ results = [] }) {
  const hasResults = results && results.length > 0;

  return (
    <div className="cyber-card" style={{ marginTop: '0.5rem' }}>
      <div className="card-header">
        <span className="card-title">EXPORT RUNTIME EVALUATION MATRIX</span>
        <span className="badge-cyan">{hasResults ? `${results.length} RUNTIMES RECORDED` : '0 RUNTIMES'}</span>
      </div>
      <div style={{ padding: 0, overflowX: 'auto', marginTop: '0.5rem' }}>
        {!hasResults ? (
          <div style={{ textAlign: 'center', padding: '2rem', color: '#94a3b8', fontSize: '0.8rem' }}>
            Nenhum resultado de exportação gravado. Execute uma compilação de benchmark acima para gerar métricas de latência e tamanho de binário.
          </div>
        ) : (
          <table className="cyber-table">
            <thead>
              <tr>
                <th>RUNTIME FORMAT</th>
                <th>STATUS</th>
                <th>ENGINE SIZE</th>
                <th>LATENCY (GPU)</th>
                <th>THROUGHPUT</th>
                <th>mAP 50-95</th>
                <th>COMPILATION ARGS</th>
              </tr>
            </thead>
            <tbody>
              {results.map((r, idx) => {
                const isTensorRT = r.format?.toLowerCase().includes('tensorrt') || r.format?.toLowerCase().includes('engine');
                return (
                  <tr key={idx}>
                    <td>
                      <span className={`speedup-badge ${isTensorRT ? 'top-tier' : ''}`}>
                        {r.format}
                      </span>
                    </td>
                    <td>
                      <span className={r.status === 'SUCCESS' ? 'badge-online' : 'badge-offline'}>
                        {r.status === 'SUCCESS' ? '✓ SUCCESS' : r.status}
                      </span>
                    </td>
                    <td className="text-mono">{r.size_mb?.toFixed(1) || '0.0'} MB</td>
                    <td className="text-mono" style={{ color: 'var(--cb-green)', fontWeight: '700' }}>
                      {r.inference_time_ms?.toFixed(2) || '0.00'} ms
                    </td>
                    <td className="text-mono" style={{ color: 'var(--cb-cyan)', fontWeight: '700' }}>
                      {r.fps?.toFixed(0) || '0'} FPS
                    </td>
                    <td className="text-mono">
                      {r.map50_95 > 0 ? `${(r.map50_95 * 100).toFixed(1)}%` : '-'}
                    </td>
                    <td className="text-mono" style={{ fontSize: '0.72rem', color: '#94a3b8' }}>
                      {r.export_args || '-'}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
