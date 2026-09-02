#!/usr/bin/env python3
"""
HydraForge Real-Time Inference Worker (Ultralytics YOLO26 / NVIDIA RTX 5090)
Outputs structured JSON detection coordinates and physical telemetry.
"""
import sys
import json
import time
import argparse
import os
import glob

def resolve_weights(weights_arg):
    if os.path.exists(weights_arg):
        return weights_arg

    clean_id = os.path.basename(weights_arg).replace(".pt", "").replace(".engine", "")
    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    candidates = [
        os.path.join(base_dir, "runs/train", clean_id, "weights", "best.pt"),
        os.path.join("/home/hades/Documents/HydraForge/runs/train", clean_id, "weights", "best.pt"),
        os.path.join("/home/hades/runs/train", clean_id, "weights", "best.pt"),
        os.path.join(base_dir, "weights", f"{clean_id}.pt"),
        os.path.join(base_dir, "weights", f"{clean_id}.engine"),
    ]
    for c in candidates:
        if os.path.exists(c):
            return c
    return weights_arg

def resolve_source(source_arg):
    if os.path.exists(source_arg):
        return source_arg

    # Check if stream ID from HydraStream
    shm_sample = f"/home/hades/Documents/HydraStream/samples/{source_arg}.jpg"
    if os.path.exists(shm_sample):
        return shm_sample

    clean_id = source_arg.replace("dataset:", "")
    ds_pattern = f"/home/hades/datasets/{clean_id}/train/images/*.jpg"
    matched = glob.glob(ds_pattern)
    if matched:
        return matched[0]

    # Global fallback to urban fleet sample
    fallback = "/home/hades/datasets/frota_urbana_4classes/train/images/rf100_adit_mp4-100_jpg.rf.1ceee906e811b590f48e8e4decda7380.jpg"
    if os.path.exists(fallback):
        return fallback
    return source_arg

def parse_args():
    parser = argparse.ArgumentParser(description="HydraForge YOLO Inference Engine")
    parser.add_argument("--weights", type=str, default="weights/yolo26n.pt", help="Path to .pt / .engine weights")
    parser.add_argument("--source", type=str, required=True, help="Image file path or URL")
    parser.add_argument("--conf", type=float, default=0.25, help="Confidence threshold")
    parser.add_argument("--iou", type=float, default=0.45, help="IoU NMS threshold")
    parser.add_argument("--device", type=str, default="0", help="CUDA device index or cpu")
    parser.add_argument("--imgsz", type=int, default=640, help="Inference image resolution")
    return parser.parse_args()

def main():
    args = parse_args()
    weights_path = resolve_weights(args.weights)
    source_path = resolve_source(args.source)

    from ultralytics import YOLO
    import torch

    try:
        model = YOLO(weights_path)
        device = 0 if (torch.cuda.is_available() and args.device != "cpu") else "cpu"

        t_start = time.perf_counter()
        results = model.predict(
            source=source_path,
            conf=args.conf,
            iou=args.iou,
            device=device,
            imgsz=args.imgsz,
            verbose=False
        )
        t_total_ms = (time.perf_counter() - t_start) * 1000.0

        detections = []
        speed = {"preprocess": 0.0, "inference": 0.0, "postprocess": 0.0}

        if len(results) > 0:
            res = results[0]
            if hasattr(res, "speed"):
                speed = res.speed

            boxes = res.boxes
            for idx, box in enumerate(boxes):
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

        inf_ms = speed.get("inference", t_total_ms)
        fps = round(1000.0 / inf_ms, 1) if inf_ms > 0 else 0

        vram_mb = 0
        if torch.cuda.is_available():
            vram_mb = round(torch.cuda.memory_allocated(0) / (1024 * 1024), 1)

        payload = {
            "status": "success",
            "model": os.path.basename(weights_path),
            "model_path": weights_path,
            "source": source_path,
            "detections": detections,
            "count": len(detections),
            "telemetry": {
                "preprocess_ms": f"{speed.get('preprocess', 0.0):.2f}",
                "inference_ms": f"{inf_ms:.2f}",
                "postprocess_ms": f"{speed.get('postprocess', 0.0):.2f}",
                "fps": f"{fps:,.0f}",
                "vram_mb": f"{vram_mb:,.0f}",
                "device": torch.cuda.get_device_name(0) if torch.cuda.is_available() else "CPU"
            }
        }
        print(json.dumps(payload))

    except Exception as exc:
        print(json.dumps({"status": "error", "error": str(exc), "detections": []}))
        sys.exit(1)

if __name__ == "__main__":
    main()
