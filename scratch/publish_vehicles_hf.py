import os
from huggingface_hub import HfApi

token = None
with open('.env') as f:
    for line in f:
        if line.startswith('HF_TOKEN='):
            token = line.strip().split('=', 1)[1].strip('"\' ')

if not token:
    raise ValueError("HF_TOKEN missing in .env")

repo_id = "hydravision/yolo26m-object"
api = HfApi(token=token)

print(f"=== PUBLISHING VEHICLE ANALYTIC TO HUGGING FACE ===")
print(f"Repository: {repo_id}")

# 1. Upload model weights as yolo26m.pt
weights_path = "storage/models/vehicle_detection.pt"
print(f"Uploading {weights_path} -> yolo26m.pt ...")
api.upload_file(
    path_or_fileobj=weights_path,
    path_in_repo="yolo26m.pt",
    repo_id=repo_id,
    repo_type="model",
    commit_message="Update YOLO26m vehicle model weights (91.1% mAP50)"
)

# 2. Upload manifest as plugin.json
manifest_path = "storage/models/vehicle_detection_manifest.json"
print(f"Uploading {manifest_path} -> plugin.json ...")
api.upload_file(
    path_or_fileobj=manifest_path,
    path_in_repo="plugin.json",
    repo_id=repo_id,
    repo_type="model",
    commit_message="Add HydraVMS plugin.json manifest for vehicle detection"
)

# 3. Model Card README.md
readme_content = """---
language:
- en
- pt
tags:
- vision
- yolo
- yolo26
- vehicle-detection
- object-detection
- tensorrt
- hydra-vision
- traffic-surveillance
license: apache-2.0
datasets:
- custom-cctv-ua-detrac
metrics:
- map
---

# YOLO26m // Traffic & Vehicle Analytic

High-precision vehicle detection and tracking model based on **YOLO26m**, trained and calibrated for urban traffic management, highway monitoring, and security CCTV video streams.

## Model Summary
- **Architecture:** YOLO26m (280 layers, 21.8M params, 75.0 GFLOPs)
- **Framework:** Ultralytics YOLO / PyTorch 2.14.1 + CUDA 13.0
- **Hardware Calibration:** NVIDIA GeForce RTX 5090 (32GB GDDR7)
- **Input Resolution:** 640x640
- **Inference Latency:** 1.2 ms / frame on RTX 5090

## Classes & Benchmark Metrics (mAP50)

| Class | Precision (P) | Recall (R) | mAP50 | mAP50-95 |
| :--- | :---: | :---: | :---: | :---: |
| **Overall (All Classes)** | **89.4%** | **84.4%** | **91.1%** | **70.3%** |
| `car` | 91.5% | 90.2% | **95.8%** | 73.0% |
| `motorcycle` | 93.6% | 88.4% | **94.9%** | 66.6% |
| `bus` | 85.3% | 82.6% | **87.7%** | 73.3% |
| `truck` | 87.2% | 76.3% | **86.1%** | 68.5% |

## Deployment in HydraVMS
Integrated with the Hydra Ecosystem via `plugin.json` zero-friction installation.
"""

print("Uploading README.md ...")
api.upload_file(
    path_or_fileobj=readme_content.encode('utf-8'),
    path_in_repo="README.md",
    repo_id=repo_id,
    repo_type="model",
    commit_message="Update model card with RTX 5090 benchmark results"
)

print("\n=== VEHICLE MODEL PUBLISHED SUCCESSFULLY! ===")
print(f"Hugging Face URL: https://huggingface.co/{repo_id}")
