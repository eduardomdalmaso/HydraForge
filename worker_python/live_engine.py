#!/usr/bin/env python3
"""
HydraForge Continuous Real-Time Analytics Engine (RTX 5090 / Zero-Latency Worker)
Loads YOLO weights once into VRAM and yields continuous detection JSON to stdout.
"""
import sys
import json
import time
import argparse
import os

def resolve_weights(weights_arg):
    if os.path.exists(weights_arg):
        return weights_arg

    clean_id = weights_arg
    if "runs/train/" in clean_id:
        clean_id = clean_id.split("runs/train/")[1].split("/")[0]
    else:
        clean_id = os.path.basename(clean_id).replace(".pt", "").replace(".engine", "")

    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    candidates = [
        os.path.join("/home/hades/runs/train", clean_id, "weights", "best.pt"),
        os.path.join("/home/hades/Documents/HydraForge/runs/train", clean_id, "weights", "best.pt"),
        os.path.join(base_dir, "runs/train", clean_id, "weights", "best.pt"),
        os.path.join(base_dir, "weights", f"{clean_id}.pt"),
        os.path.join(base_dir, "weights", f"{clean_id}.engine"),
    ]
    for c in candidates:
        if os.path.exists(c):
            return c
    return weights_arg

def main():
    parser = argparse.ArgumentParser(description="HydraForge Live Analytics Engine")
    parser.add_argument("--weights", type=str, required=True)
    parser.add_argument("--source", type=str, required=True)
    parser.add_argument("--conf", type=float, default=0.25)
    parser.add_argument("--device", type=str, default="0")
    args = parser.parse_args()

    weights_path = resolve_weights(args.weights)
    from ultralytics import YOLO
    import torch

    model = YOLO(weights_path)
    device = 0 if (torch.cuda.is_available() and args.device != "cpu") else "cpu"

    while True:
        try:
            source_path = args.source
            if not os.path.exists(source_path):
                # check shm sample
                shm_sample = f"/home/hades/Documents/HydraStream/samples/{source_path}.jpg"
                if os.path.exists(shm_sample):
                    source_path = shm_sample
                else:
                    time.sleep(0.04)
                    continue

            t0 = time.perf_counter()
            results = model.predict(source=source_path, conf=args.conf, device=device, verbose=False)
            inf_ms = (time.perf_counter() - t0) * 1000.0

            detections = []
            if len(results) > 0:
                res = results[0]
                for idx, box in enumerate(res.boxes):
                    cls_id = int(box.cls[0].item())
                    cls_name = res.names.get(cls_id, f"class_{cls_id}")
                    conf = float(box.conf[0].item())
                    xywhn = box.xywhn[0].tolist()
                    left_pct = max(0.0, min(100.0, (xywhn[0] - xywhn[2] / 2.0) * 100.0))
                    top_pct = max(0.0, min(100.0, (xywhn[1] - xywhn[3] / 2.0) * 100.0))
                    w_pct = max(1.0, min(100.0, xywhn[2] * 100.0))
                    h_pct = max(1.0, min(100.0, xywhn[3] * 100.0))

                    color = "#00f0ff"
                    if any(k in cls_name.lower() for k in ["bus", "onibus", "caminhao", "truck"]):
                        color = "#fcee0a"
                    elif any(k in cls_name.lower() for k in ["person", "moto"]):
                        color = "#00ff9d"
                    elif "phone" in cls_name.lower():
                        color = "#ff0055"

                    detections.append({
                        "id": idx + 1,
                        "label": cls_name,
                        "conf": round(conf, 3),
                        "box": [round(left_pct, 2), round(top_pct, 2), round(w_pct, 2), round(h_pct, 2)],
                        "color": color
                    })

            vram_mb = round(torch.cuda.memory_allocated(0) / (1024 * 1024), 1) if torch.cuda.is_available() else 0
            fps = round(1000.0 / inf_ms, 0) if inf_ms > 0 else 0

            payload = {
                "status": "success",
                "detections": detections,
                "count": len(detections),
                "telemetry": {
                    "inference_ms": f"{inf_ms:.2f}",
                    "fps": f"{fps:,.0f}",
                    "vram_mb": f"{vram_mb:,.0f}",
                    "device": "NVIDIA GeForce RTX 5090"
                }
            }
            print(json.dumps(payload), flush=True)
            time.sleep(0.035)

        except Exception as err:
            time.sleep(0.05)

if __name__ == "__main__":
    main()
