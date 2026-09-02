import React, { useState, useEffect, useRef } from 'react';
import PlaygroundControlBar from '../components/playground/PlaygroundControlBar';
import PlaygroundTuningCard from '../components/playground/PlaygroundTuningCard';
import PlaygroundViewport from '../components/playground/PlaygroundViewport';
import PlaygroundTelemetryCard from '../components/playground/PlaygroundTelemetryCard';
import { useWebcamStream } from '../components/playground/useWebcamStream';
import { fetchHydraStreams, fetchHydraTelemetry } from '../api/hydrastream_client';
import { fetchModelsAPI, fetchTelemetryAPI } from '../api/client';
import { runRealInferenceAPI } from '../api/inference_client';

export default function PlaygroundStudio() {
  const [config, setConfig] = useState({
    model: 'yolo26n', source: 'cam_entrance_01', runtime: 'pytorch',
    conf: 0.25, iou: 0.45, nmsFree: true, sahi: false, isolateBg: false
  });
  const [modelsList, setModelsList] = useState([]);
  const [hydraStreams, setHydraStreams] = useState([]);
  const [hydraTelemetry, setHydraTelemetry] = useState(null);
  const [gpuStats, setGpuStats] = useState(null);
  const [detections, setDetections] = useState([]);
  const [selectedEntity, setSelectedEntity] = useState(null);
  const [isScanning, setIsScanning] = useState(false);
  const [isContinuous, setIsContinuous] = useState(false);
  const [telemetry, setTelemetry] = useState(null);
  const [imageSrc, setImageSrc] = useState('/api/v1/hydrastream/api/v1/streams/cam_entrance_01/snapshot.jpg');

  const isWebcam = config.source === 'webcam';
  const activeStream = hydraStreams.find(s => s.stream_id === config.source) || (hydraStreams[0] || null);
  const { videoRef, captureFrame } = useWebcamStream(isWebcam);

  useEffect(() => {
    (async () => {
      const [models, streams, stats, hw] = await Promise.all([
        fetchModelsAPI(), fetchHydraStreams(), fetchHydraTelemetry(), fetchTelemetryAPI()
      ]);
      if (models?.length) { setModelsList(models); const b = models.find(m => m.isCustom) || models[0]; setConfig(c => ({ ...c, model: b.id })); }
      if (streams?.length) setHydraStreams(streams);
      if (stats) setHydraTelemetry(stats);
      if (hw?.gpu_stats) setGpuStats(hw.gpu_stats);
    })();
    const int = setInterval(async () => {
      const [s, h] = await Promise.all([fetchHydraTelemetry(), fetchTelemetryAPI()]);
      if (s) setHydraTelemetry(s);
      if (h?.gpu_stats) setGpuStats(h.gpu_stats);
    }, 3000);
    return () => clearInterval(int);
  }, []);

  const handleInference = async () => {
    const b64 = captureFrame();
    const res = await runRealInferenceAPI(config, b64);
    if (res) {
      if (res.displayImageUrl && !isWebcam && !isContinuous) setImageSrc(res.displayImageUrl);
      if (Array.isArray(res.detections)) setDetections(res.detections);
      if (res.telemetry) setTelemetry(res.telemetry);
    }
  };

  useEffect(() => {
    if (!isContinuous) {
      if (!isWebcam) setImageSrc(`/api/v1/hydrastream/api/v1/streams/${config.source}/snapshot.jpg?t=${Date.now()}`);
      return;
    }
    if (!isWebcam) setImageSrc(`/api/v1/hydrastream/api/v1/streams/${config.source}/mjpeg?t=${Date.now()}`);
    let intId = isWebcam ? setInterval(() => {
      const b64 = captureFrame();
      if (b64) fetch('/api/v1/inference/frame', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ image_base64: b64 }) }).catch(() => {});
    }, 40) : null;
    const es = new EventSource(`/api/v1/inference/live?model=${encodeURIComponent(config.model)}&source=${encodeURIComponent(config.source)}&conf=${config.conf}`);
    es.onmessage = (e) => {
      try {
        const d = JSON.parse(e.data);
        if (Array.isArray(d.detections)) setDetections(d.detections);
        if (d.telemetry) setTelemetry(d.telemetry);
      } catch {}
    };
    return () => { es.close(); if (intId) clearInterval(intId); };
  }, [isContinuous, config.source, config.model, config.conf]);

  return (
    <div className="view-container playground-container">
      <div className="cockpit-full-header">
        <h1 className="cockpit-main-title">KIROSHI PLAYGROUND // REALTIME TRACKER</h1>
        <p className="cockpit-main-subtitle">HARDWARE NATIVO RTX 5090 • FLUXO CONTÍNUO HYDRASTREAM • BBOX REALTIME</p>
      </div>
      <div className="playground-layout">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <PlaygroundControlBar config={config} setConfig={setConfig} modelsList={modelsList} hydraStreams={hydraStreams} isHydraOnline={hydraStreams.length > 0} onRunInference={async () => { setIsScanning(true); await handleInference(); setIsScanning(false); }} isRunning={isScanning} isContinuous={isContinuous} setIsContinuous={setIsContinuous} />
          <PlaygroundTuningCard config={config} setConfig={setConfig} />
        </div>
        <PlaygroundViewport imageSrc={imageSrc} detections={detections} selectedEntity={selectedEntity} setSelectedEntity={setSelectedEntity} isScanning={isScanning} isContinuous={isContinuous} isHydraLinked={!isWebcam && hydraStreams.length > 0} activeStream={activeStream} isWebcam={isWebcam} videoRef={videoRef} />
        <PlaygroundTelemetryCard telemetry={telemetry} detections={detections} selectedEntity={selectedEntity} onSelectEntity={setSelectedEntity} hydraTelemetry={hydraTelemetry} activeStream={activeStream} gpuStats={gpuStats} />
      </div>
    </div>
  );
}
