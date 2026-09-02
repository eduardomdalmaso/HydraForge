import React, { useState, useEffect } from 'react';
import BenchmarkLauncherCard from '../components/benchmarks/BenchmarkLauncherCard';
import BenchmarkSummaryCards from '../components/benchmarks/BenchmarkSummaryCards';
import BenchmarkResultsTable from '../components/benchmarks/BenchmarkResultsTable';
import BenchmarkThroughputChart from '../components/benchmarks/BenchmarkThroughputChart';
import { fetchBenchmarksAPI, launchBenchmarkAPI, fetchBenchmarkByIDAPI } from '../api/benchmark_client';

export default function BenchmarkStudio() {
  const [activeJob, setActiveJob] = useState(null);
  const [isRunning, setIsRunning] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  const loadLatestBenchmark = async () => {
    const list = await fetchBenchmarksAPI();
    if (list && list.length > 0) {
      const latest = list[list.length - 1];
      setActiveJob(latest);
      setIsRunning(latest.status === 'RUNNING' || latest.status === 'QUEUED');
    }
  };

  useEffect(() => {
    loadLatestBenchmark();
    const interval = setInterval(async () => {
      if (activeJob?.job_id && isRunning) {
        const updated = await fetchBenchmarkByIDAPI(activeJob.job_id);
        if (updated) {
          setActiveJob(updated);
          if (updated.status !== 'RUNNING' && updated.status !== 'QUEUED') {
            setIsRunning(false);
          }
        }
      }
    }, 1000);
    return () => clearInterval(interval);
  }, [activeJob?.job_id, isRunning]);

  const handleLaunch = async (config) => {
    setErrorMsg('');
    try {
      setIsRunning(true);
      const created = await launchBenchmarkAPI(config);
      setActiveJob(created);
    } catch (err) {
      setErrorMsg(err.message || 'Failed to start benchmark');
      setIsRunning(false);
    }
  };

  return (
    <div className="view-container benchmark-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">BENCHMARK & EXPORT STUDIO</h1>
        <p className="cockpit-main-subtitle">
          HARDWARE ACCELERATION • TENSORRT 10.x • ONNX • OPENVINO & THROUGHPUT PROFILING
        </p>
      </div>

      {errorMsg && (
        <div className="cyber-alert alert-critical">
          <span>⚠️ {errorMsg}</span>
        </div>
      )}

      <BenchmarkSummaryCards results={activeJob?.results} />

      <div className="benchmark-grid-top">
        <BenchmarkLauncherCard onLaunch={handleLaunch} isRunning={isRunning} />
        <BenchmarkThroughputChart results={activeJob?.results} />
      </div>

      <BenchmarkResultsTable results={activeJob?.results} />
    </div>
  );
}
