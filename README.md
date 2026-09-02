# HydraForge

> **High-Performance Cyberpunk-themed AI Training Studio for YOLOv8, YOLO11, and YOLO26. Built with Go Control Plane, PyTorch/CUDA 13.3 Worker, and React SPA with 1-Click TensorRT Export.**

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
2. **Disconnected Monitoring:** Relying on heavy external SaaS tracking platforms that introduce latency and lack direct GPU hardware telemetry (VRAM, Tensor Core load, temperature).
3. **Complex Model Compilation & Export:** Converting raw PyTorch `.pt` checkpoints into optimized **TensorRT (`.engine`)** with mixed precision (FP16/FP8/INT8) requires error-prone C++ toolchains and manual TensorRT engine building.

---

## The Solution

**HydraForge** provides an all-in-one, high-performance **AI Model Training Cockpit & Compilation Studio**:

- **Native YOLOv8, YOLO11 & YOLO26 Support:** Full architecture matrix covering **Nano (n)**, **Small (s)**, **Medium (m)**, **Large (l)**, and **XLarge (x)** models across **Detection**, **Segmentation**, **Pose**, **Classification**, and **OBB (Oriented Bounding Box)** tasks.
- **Go Control Plane + Python CUDA Worker:** High-concurrency Go server (:8081) managing the training queue and broadcasting live WebSockets telemetry, backed by an isolated Python PyTorch/CUDA 13.3 worker executing directly on the **NVIDIA GeForce RTX 5090 (32GB VRAM)**.
- **Real-Time Cyberpunk HUD:** Live SVG Bézier curves for loss functions (`box_loss`, `cls_loss`, `dfl_loss`), precision metrics (`mAP@50`, `mAP@50-95`), and hardware sensors (VRAM usage, Tensor Cores, Watts).
- **Dataset Studio & Class Balance Matrix:** Ingests local datasets or `data.yaml`, inspects bounding boxes visually, and detects dataset class imbalance before training.
- **1-Click TensorRT & ONNX Exporter:** Compiles completed training runs directly into high-throughput **TensorRT `.engine`** files optimized for sub-millisecond edge inference in [HydraStream](https://github.com/eduardomdalmaso/HydraStream).

---

## Architecture Overview

```mermaid
flowchart TD
    subgraph Frontend [React SPA + Cyberpunk CSS]
        UI_Cockpit[Training Cockpit & Hyperparameter Launcher]
        UI_LiveHUD["Live Training Curves: mAP50, Box/Cls Loss, VRAM"]
        UI_Datasets[Dataset Studio & Annotation Viewer]
        UI_Export["Model Zoo & 1-Click TensorRT / ONNX Exporter"]
    end

    subgraph GoControlPlane [Go Control Plane Engine :8081]
        RESTRouter[REST API: /api/v1/training/*]
        WSHub[WebSockets Telemetry Hub]
        JobQueue[Training Job Queue & State Machine]
        ArtifactsDB[Model Checkpoint & Dataset Registry]
    end

    subgraph PythonWorker [Python PyTorch Training Worker]
        Trainer[Ultralytics Engine YOLOv8 / YOLO11 / YOLO26]
        PyTorchCUDA["PyTorch 2.x + CUDA 13.3 (RTX 5090 32GB)"]
        Callbacks[Epoch Callback & Metric Emitter]
        Exporter[TensorRT FP8/FP16 & ONNX Precision Exporter]
    end

    Frontend <-->|REST & WebSockets| GoControlPlane
    GoControlPlane <-->|IPC & Process Control| PythonWorker
    PythonWorker -->|Streams loss, mAP, GPU metrics| GoControlPlane
    GoControlPlane -->|Broadcasts live telemetry| Frontend
```

---

## Key Studio Views

```
┌───────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                     HYDRAFORGE - STUDIO VIEWS                                     │
├───────────────┬─────────────────┬──────────────────┬──────────────────┬─────────────────┬─────────┤
│  1. COCKPIT   │  2. LIVE HUD    │  3. BENCHMARKS   │  4. DATASETS     │  5. MODEL ZOO   │ 6. TEST │
│   (LAUNCHER)  │  (MONITORING)   │   (SPEED & MAP)  │   (DATA & YAML)  │   (EXPORTER)    │ (PLAYG) │
└───────────────┴─────────────────┴──────────────────┴──────────────────┴─────────────────┴─────────┘
```

1. **Training Cockpit (`#cockpit`):** Model architecture selector, task variant picker, optimizer selection (`AdamW`, `SGD`, `RMSprop`), Cosine LR scheduling, AMP FP16/BF16 toggles, and dynamic VRAM estimation.
2. **Live Training HUD (`#live-hud`):** Real-time loss and mAP curves, visual validation gallery with predicted bounding boxes per epoch, and RTX 5090 GPU hardware monitors.
3. **Benchmark Studio (`#benchmarks`):** Ultralytics multi-format benchmark launcher with live throughput comparison charts and speedup multipliers.
4. **Dataset Studio (`#datasets`):** `data.yaml` validation, train/val/test split sliders, and class distribution frequency bar chart.
5. **Model Zoo & Exporter (`#model-zoo`):** Checkpoint comparison (`best.pt` vs `last.pt`), size vs mAP trade-off matrix, and 1-click **TensorRT (`.engine`)** compilation.
6. **Inference Playground (`#playground`):** Drag-and-drop image/video testing with real-time confidence and IoU threshold sliders.

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

## License

This project is licensed under the [MIT License](LICENSE).
