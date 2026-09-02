import React, { useState, useEffect } from 'react';
import LiveTrainingStatusCard from '../components/livehud/LiveTrainingStatusCard';
import LiveMetricsChartsCard from '../components/livehud/LiveMetricsChartsCard';
import LiveGpuTelemetryCard from '../components/livehud/LiveGpuTelemetryCard';
import LiveTerminalLogsCard from '../components/livehud/LiveTerminalLogsCard';
import ExperimentsComparisonCard from '../components/livehud/ExperimentsComparisonCard';
import { fetchJobsAPI, fetchTelemetryAPI, stopTrainingJobAPI, launchTrainingJobAPI } from '../api/client';

export default function LiveHudStudio() {
  const [jobs, setJobs] = useState([]);
  const [activeJob, setActiveJob] = useState(null);
  const [telemetry, setTelemetry] = useState(null);

  const loadLiveData = async () => {
    const [list, tel] = await Promise.all([fetchJobsAPI(), fetchTelemetryAPI()]);
    if (list) setJobs(list);
    if (tel) {
      setTelemetry(tel);
      if (tel.active_job) {
        setActiveJob(tel.active_job);
      } else if (list) {
        const running = list.find(j => j.status === 'TRAINING' || j.status === 'RUNNING' || j.status === 'QUEUED');
        setActiveJob(running || list[0] || null);
      }
    } else if (list) {
      const running = list.find(j => j.status === 'TRAINING' || j.status === 'RUNNING' || j.status === 'QUEUED');
      setActiveJob(running || list[0] || null);
    }
  };

  useEffect(() => {
    loadLiveData();
    const interval = setInterval(loadLiveData, 1000);
    return () => clearInterval(interval);
  }, []);

  const handleAbort = async () => {
    if (activeJob?.job_id) {
      await stopTrainingJobAPI(activeJob.job_id);
      loadLiveData();
    }
  };

  const handleRestart = async (job) => {
    if (!job) return;
    if (job.status === 'TRAINING' || job.status === 'RUNNING') {
      await stopTrainingJobAPI(job.job_id);
    }
    await launchTrainingJobAPI({
      model_architecture: job.model_architecture || 'yolo26m',
      task: job.task || 'detect',
      dataset_id: job.dataset_id,
      hyperparameters: { ...job.hyperparameters, pretrained: true }
    });
    loadLiveData();
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
          <LiveTrainingStatusCard
            job={activeJob}
            onAbortJob={handleAbort}
            onRestartJob={handleRestart}
            onResumeJob={handleRestart}
          />
          <LiveGpuTelemetryCard gpuStats={telemetry?.gpu_stats} job={activeJob} />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <LiveMetricsChartsCard job={activeJob} recentMetrics={telemetry?.recent_metrics} />
          <LiveTerminalLogsCard job={activeJob} rawLogs={telemetry?.raw_logs} recentMetrics={telemetry?.recent_metrics} gpuStats={telemetry?.gpu_stats} />
        </div>
      </div>

      <ExperimentsComparisonCard jobs={jobs} onTestInPlayground={() => { window.location.hash = 'playground'; }} />
    </div>
  );
}
