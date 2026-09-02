# HydraForge

> **High-Performance Cyberpunk-themed AI Training Studio for YOLOv8, YOLO11, and YOLO26. Built with Go Control Plane, PyTorch/CUDA 13.3 Worker, Real-Time Intra-Epoch Telemetry, and React SPA with 1-Click TensorRT Export.**

[**English**] | [**Português do Brasil**](README.pt-BR.md)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![PyTorch](https://img.shields.io/badge/PyTorch-2.x-EE4C2C?logo=pytorch&logoColor=white)](https://pytorch.org/)
[![CUDA](https://img.shields.io/badge/CUDA-13.3%20%7C%20RTX%205090-76B900?logo=nvidia&logoColor=white)](https://developer.nvidia.com/cuda-toolkit)
[![Ultralytics](https://img.shields.io/badge/Ultralytics-YOLOv8%20%7C%2011%20%7C%2026-blue)](https://github.com/ultralytics/ultralytics)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![TensorRT](https://img.shields.io/badge/TensorRT-10.x%20Export-76B900?logo=nvidia&logoColor=white)](https://developer.nvidia.com/tensorrt)

---

## The Problem

Training, fine-tuning, and optimizing state-of-the-art YOLO computer vision models for production edge deployments is typically hindered by fragmented tooling:

1. **Scattered CLI Scripts & Fragile Workflows:** Manual orchestration of dataset formatting (`data.yaml`), batch sizing, and learning rate scheduling across CLI flags.
2. **Disconnected Monitoring & Delayed Telemetry:** Most platforms only report metrics at epoch boundaries, leaving operators blind during long 70+ second epochs on large datasets.
3. **Complex Model Compilation & Export:** Converting raw PyTorch `.pt` checkpoints into optimized **TensorRT (`.engine`)** with mixed precision (FP16/FP8/INT8) requires error-prone C++ toolchains and manual TensorRT engine building.

---

## The Solution

**HydraForge** provides an all-in-one, high-performance **AI Model Training Cockpit & Compilation Studio**:

- **Native YOLOv8, YOLO11 & YOLO26 Support:** Full architecture matrix covering **Nano (n)**, **Small (s)**, **Medium (m)**, **Large (l)**, and **XLarge (x)** models across **Detection**, **Segmentation**, **Pose**, **Classification**, and **OBB (Oriented Bounding Box)** tasks.
- **Go Control Plane + Python CUDA Worker:** High-concurrency Go server (`:8081`) managing the training queue and broadcasting live WebSockets/REST telemetry, backed by an isolated Python PyTorch/CUDA 13.3 worker executing directly on the **NVIDIA GeForce RTX 5090 (32GB VRAM)**.
- **Intra-Epoch Batch Streaming:** Live telemetry emitted every 2–5 batches (~150ms) tracking fractional progress, real-time milliwatt energy integration ($kWh$), and train throughput ($FPS$).
- **Dual Progress & Lifecycle Controls:** Independent bars for **Current Epoch Progress (0–100%)** and **Total Session Progress**, with instant **ABORT**, **RESTART**, and **RESUME** actions.
- **Multi-GPU Telemetry Cycler:** Aggregates total cluster VRAM while automatically cycling between individual device sensors (Temp, Watts, SM load) every 3 seconds.
- **1-Click TensorRT & ONNX Exporter:** Compiles completed training runs directly into high-throughput **TensorRT `.engine`** files optimized for sub-millisecond edge inference in [HydraStream](https://github.com/eduardomdalmaso/HydraStream).

---

## Architecture Overview

```mermaid
flowchart TD
    subgraph Frontend [React SPA + Cyberpunk CSS Design System]
        UI_Cockpit[Training Cockpit & Hyperparameter Launcher]
        UI_LiveHUD["Live Dual Progress Bars, Loss/mAP Charts, GPU HUD"]
        UI_Datasets[Dataset Studio & Annotation Viewer]
        UI_Export["Model Zoo & 1-Click TensorRT / ONNX Exporter"]
    end

    subgraph GoControlPlane [Go Control Plane Engine :8081]
        RESTRouter[REST API: /api/v1/training/*]
        WSHub[WebSockets Telemetry Hub]
        JobQueue[Training Job Queue & In-Memory Active Job State]
        SQLiteDB[SQLite Database in WAL Mode]
    end

    subgraph PythonWorker [Python PyTorch Training Worker]
        Trainer[Ultralytics Engine YOLOv8 / YOLO11 / YOLO26]
        PyTorchCUDA["PyTorch 2.x + CUDA 13.3 (RTX 5090 32GB)"]
        BatchCallback["on_train_batch_end: Intra-Batch FPS, kWh & Loss"]
        EpochCallback["on_fit_epoch_end: mAP50, Precision, Recall"]
        Exporter[TensorRT FP8/FP16 & ONNX Precision Exporter]
    end

    Frontend <-->|REST Polling & WebSockets| GoControlPlane
    GoControlPlane <-->|IPC Process Control & Scanner| PythonWorker
    PythonWorker -->|Streams HYDRA_BATCH & HYDRA_METRIC| GoControlPlane
    GoControlPlane -->|Persists & Broadcasts Live Telemetry| Frontend
```

---

## Key Studio Views

```text
┌───────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                     HYDRAFORGE - STUDIO VIEWS                                     │
├───────────────┬─────────────────┬──────────────────┬──────────────────┬─────────────────┬─────────┤
│  1. COCKPIT   │  2. LIVE HUD    │  3. BENCHMARKS   │  4. DATASETS     │  5. MODEL ZOO   │ 6. TEST │
│   (LAUNCHER)  │  (MONITORING)   │   (SPEED & MAP)  │   (DATA & YAML)  │   (EXPORTER)    │ (PLAYG) │
└───────────────┴─────────────────┴──────────────────┴──────────────────┴─────────────────┴─────────┘
```

1. **Training Cockpit (`#cockpit`):** Model architecture selector, task variant picker, optimizer selection (`AdamW`, `SGD`), AMP FP16/BF16 toggles, two-stage fine-tuning switches, and dynamic VRAM estimation.
2. **Live Training HUD (`#live-hud`):** Dual progress bars (Epoch & Session), real-time energy budget ($kWh$), train throughput ($FPS$), vector loss/mAP charts, live PyTorch terminal logs, and multi-GPU sensor cycling.
3. **Benchmark Studio (`#benchmarks`):** Ultralytics multi-format benchmark launcher with live throughput comparison charts and speedup multipliers.
4. **Dataset Studio (`#datasets`):** `data.yaml` validation, train/val/test split sliders, class distribution frequency bar chart, and integrated bounding-box visualizer.
5. **Model Zoo & Exporter (`#model-zoo`):** Checkpoint comparison (`best.pt` vs `last.pt`), size vs mAP trade-off matrix, official & custom model repository, and 1-click **TensorRT (`.engine`)** compilation.
6. **Inference Playground (`#playground`):** Drag-and-drop image/video testing with real-time confidence and IoU threshold sliders.

---

## Engineering Guidelines & Strict Modularity

- **Hexagonal Architecture (DDD):** Pure domain logic in `internal/domain/` with zero HTTP/network dependencies.
- **Modular Frontend:** Strict `<100 lines per file` rule across all CSS, JS, and JSX files in `web/`.
- **Hardware Acceleration:** Native PyTorch with CUDA 13.3 (`sm_120` Blackwell support on RTX 5090).

---

## Quick Start & Development

### 1. Start HydraForge in Live Dev Mode
```bash
make dev
```
> Starts the Go Control Plane & React Training Studio on **`http://localhost:8081`**.

### 2. Run Test Suite
```bash
make test
```

### 3. Build Production Binaries
```bash
make build
```

---

## Roadmap

- [ ] **Fix Current Epoch Progress Bar Stream:** Correct intra-epoch batch callback telemetry payload synchronization so the Current Epoch progress bar animates live during intra-epoch batches.
- [ ] **Quantization-Aware Training (QAT) Engine:** Fine-tuning loop emulating INT8/FP8 quantization noise during training for zero-accuracy-loss edge deployment on NVIDIA TensorRT and edge NPUs.
- [ ] **Interactive Confusion Matrix & Class Diagnostics:** In-HUD drawer for granular per-class precision, recall, and F1-confidence curve optimization.
- [ ] **Overfitting Early-Warning Sentinel:** Real-time heuristic anomaly detection alerting on validation loss divergence.
- [ ] **HydraVault Active Learning Pipeline:** 1-click dataset synchronization with [HydraVault](https://github.com/eduardomdalmaso/HydraVault) for automated hard-negative retraining.

---

## License

This project is licensed under the [MIT License](LICENSE).
