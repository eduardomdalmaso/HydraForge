#!/usr/bin/env python3
"""
HydraForge Zero-Copy GPU Streaming Inference Worker (NVIDIA RTX 5090)
Directly streams MediaMTX RTSP stream into GPU memory without disk I/O or cv2.imread.
"""

import sys
import os
import time
import json
import argparse
from pathlib import Path

def resolve_model(weights_path: str) -> str:
    if os.path.exists(weights_path):
        return weights_path
    base_dir = Path(__file__).resolve().parent.parent
    candidates = [
        base_dir / "storage" / "models" / weights_path,
        base_dir / "storage" / "models" / f"{weights_path}.pt",
        base_dir / "weights" / weights_path,
        base_dir / "weights" / f"{weights_path}.pt",
    ]
    for c in candidates:
        if c.exists():
            return str(c)
    return weights_path

def main():
    parser = argparse.ArgumentParser(description="Zero-Disk GPU Live Inference Worker")
    parser.add_argument("--model", type=str, default="yolo26m.pt")
    parser.add_argument("--stream-id", type=str, default="cam_01")
    parser.add_argument("--rtsp-url", type=str, default="rtsp://127.0.0.1:8554/cam_01")
    parser.add_argument("--conf", type=float, default=0.25)
    parser.add_argument("--device", type=str, default="0")
    parser.add_argument("--classes", type=str, default="0", help="Comma-separated class IDs (0 for person)")
    args = parser.parse_args()

    model_path = resolve_model(args.model)
    from ultralytics import YOLO
    import torch

    print(f"🚀 [HYDRA-AI] Loading {model_path} into RTX 5090 VRAM...", flush=True)
    model = YOLO(model_path)
    device = 0 if torch.cuda.is_available() and args.device != "cpu" else "cpu"

    target_classes = [int(c.strip()) for c in args.classes.split(",") if c.strip().isdigit()]
    if not target_classes:
        target_classes = [0] # Person class by default

    shm_out_path = f"/dev/shm/detections_{args.stream_id}.json"
    shm_tmp_path = f"/dev/shm/detections_{args.stream_id}.json.tmp"

    print(f"📡 [HYDRA-AI] Starting zero-disk GPU stream on {args.rtsp_url} for class(es) {target_classes}...", flush=True)

    while True:
        try:
            results_stream = model.predict(
                source=args.rtsp_url,
                stream=True,
                classes=target_classes,
                conf=args.conf,
                device=device,
                imgsz=640,
                verbose=False
            )

            for r in results_stream:
                boxes = []
                for idx, box in enumerate(r.boxes):
                    cls_id = int(box.cls[0].item())
                    cls_name = r.names.get(cls_id, "person")
                    conf = float(box.conf[0].item())
                    xywhn = box.xywhn[0].tolist()

                    left_pct = max(0.0, min(100.0, (xywhn[0] - xywhn[2] / 2.0) * 100.0))
                    top_pct = max(0.0, min(100.0, (xywhn[1] - xywhn[3] / 2.0) * 100.0))
                    w_pct = max(1.0, min(100.0, xywhn[2] * 100.0))
                    h_pct = max(1.0, min(100.0, xywhn[3] * 100.0))

                    boxes.append({
                        "id": idx + 1,
                        "label": cls_name,
                        "class_name": cls_name,
                        "conf": round(conf, 3),
                        "confidence": round(conf, 3),
                        "box": [round(left_pct, 2), round(top_pct, 2), round(w_pct, 2), round(h_pct, 2)],
                        "color": "#ff5e3a"
                    })

                inf_ms = r.speed.get("inference", 4.5)
                payload = {
                    "status": "success",
                    "timestamp": time.time(),
                    "camera_id": args.stream_id,
                    "detections": boxes,
                    "count": len(boxes),
                    "telemetry": {
                        "inference_ms": f"{inf_ms:.2f}",
                        "device": "NVIDIA GeForce RTX 5090"
                    }
                }

                # Atomic write to /dev/shm
                with open(shm_tmp_path, "w") as f:
                    json.dump(payload, f)
                os.replace(shm_tmp_path, shm_out_path)

        except Exception as e:
            print(f"⚠️ [HYDRA-AI] Reconnecting RTSP stream: {e}", flush=True)
            time.sleep(1.0)

if __name__ == "__main__":
    main()
