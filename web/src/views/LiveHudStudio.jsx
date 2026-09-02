import React, { useState, useEffect } from 'react';
import LiveTrainingStatusCard from '../components/livehud/LiveTrainingStatusCard';
import LiveMetricsChartsCard from '../components/livehud/LiveMetricsChartsCard';
import LiveGpuTelemetryCard from '../components/livehud/LiveGpuTelemetryCard';
import LiveTerminalLogsCard from '../components/livehud/LiveTerminalLogsCard';
import ExperimentsComparisonCard from '../components/livehud/ExperimentsComparisonCard';
import { fetchJobsAPI, fetchTelemetryAPI, stopTrainingJobAPI } from '../api/client';

export default function LiveHudStudio() {
  const [jobs, setJobs] = useState([]);
  const [activeJob, setActiveJob] = useState(null);
  const [telemetry, setTelemetry] = useState(null);

  const loadLiveData = async () => {
    const list = await fetchJobsAPI();
    if (list) {
      setJobs(list);
      const running = list.find(j => j.status === 'RUNNING' || j.status === 'QUEUED');
      setActiveJob(running || list[0] || null);
    }
    const tel = await fetchTelemetryAPI();
    if (tel) setTelemetry(tel);
  };

  useEffect(() => {
    loadLiveData();
    const interval = setInterval(loadLiveData, 1500);
    return () => clearInterval(interval);
  }, []);

  const handleAbort = async () => {
    if (activeJob?.job_id) {
      await stopTrainingJobAPI(activeJob.job_id);
      loadLiveData();
    }
  };

  return (
    <div className="view-container livehud-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">LIVE HUD TELEMETRY & VALIDATION STUDIO</h1>
        <p className="cockpit-main-subtitle">
          REAL-TIME PYTORCH RUNTIMES • LOSS & mAP DYNAMICS • GPU HARDWARE TELEMETRY & EXPERIMENT MATRIX
        </p>
      </div>

      <div className="livehud-grid">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <LiveTrainingStatusCard job={activeJob} onAbortJob={handleAbort} />
          <LiveGpuTelemetryCard gpuStats={telemetry?.gpu_stats} job={activeJob} />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <LiveMetricsChartsCard job={activeJob} />
          <LiveTerminalLogsCard job={activeJob} />
        </div>
      </div>

      <ExperimentsComparisonCard jobs={jobs} onTestInPlayground={() => { window.location.hash = 'playground'; }} />
    </div>
  );
}
