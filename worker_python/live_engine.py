#!/usr/bin/env python3
"""
HydraForge Continuous Real-Time Analytics Engine (RTX 5090 / Zero-Latency Worker)
Loads YOLO weights once into VRAM and yields continuous detection JSON to stdout.
Reads live frames via OpenCV VideoCapture (RTSP/RTMP) or HydraStream live snapshots.
"""
import sys
import json
import time
import argparse
import os
import urllib.request
import numpy as np

def resolve_weights(weights_arg):
    if os.path.exists(weights_arg):
        return weights_arg

    clean_id = weights_arg
    if "runs/train/" in clean_id:
        clean_id = clean_id.split("runs/train/")[1].split("/")[0]
    else:
        clean_id = os.path.basename(clean_id).replace(".pt", "").replace(".engine", "")

    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    home_dir = os.path.expanduser("~")
    runs_env = os.environ.get("RUNS_DIR", "")

    candidates = [
        os.path.join(base_dir, "runs/train", clean_id, "weights", "best.pt"),
        os.path.join(base_dir, "storage/models", f"{clean_id}.pt"),
        os.path.join(base_dir, "storage/models", clean_id),
        os.path.join(base_dir, "storage/engines", f"{clean_id}.engine"),
        os.path.join(base_dir, "weights", f"{clean_id}.pt"),
        os.path.join(base_dir, "weights", f"{clean_id}.engine"),
    ]
    if runs_env:
        candidates.insert(0, os.path.join(runs_env, clean_id, "weights", "best.pt"))
    candidates.extend([
        os.path.join(home_dir, "runs/train", clean_id, "weights", "best.pt"),
        os.path.join(home_dir, "Documents/HydraForge/runs/train", clean_id, "weights", "best.pt"),
    ])
    for c in candidates:
        if os.path.exists(c):
            return c
    return weights_arg

def main():
    parser = argparse.ArgumentParser(description="HydraForge Live Analytics Engine")
    parser.add_argument("--weights", type=str, required=True)
    parser.add_argument("--source", type=str, required=True)
    parser.add_argument("--classes", type=str, default="")
    parser.add_argument("--conf", type=float, default=0.25)
    parser.add_argument("--iou", type=float, default=0.45)
    parser.add_argument("--device", type=str, default="0")
    parser.add_argument("--sahi", action="store_true")
    parser.add_argument("--nms-free", action="store_true")
    args = parser.parse_args()

    weights_path = resolve_weights(args.weights)
    from ultralytics import YOLO
    import torch
    import cv2
    from sahi_engine import run_sahi_inference

    model = YOLO(weights_path)
    device = 0 if (torch.cuda.is_available() and args.device != "cpu") else "cpu"
    effective_iou = 1.0 if args.nms_free else args.iou
    target_classes = [c.strip().lower() for c in args.classes.split(",") if c.strip()]

    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    home_dir = os.path.expanduser("~")
    samples_env = os.environ.get("HYDRASTREAM_SAMPLES_DIR", "")

    source = args.source
    is_rtsp = source.startswith("rtsp://") or source.startswith("rtmp://")
    is_http = source.startswith("http://") or source.startswith("https://")
    cam_id = os.path.basename(source).replace(".jpg", "")

    captures_dir = os.path.join(base_dir, "storage/captures")
    os.makedirs(captures_dir, exist_ok=True)

    # If camera ID, resolve RTSP source URL from HydraStream
    if not is_rtsp and not is_http:
        try:
            req = urllib.request.Request("http://127.0.0.1:8080/api/v1/streams", headers={"User-Agent": "HydraForge/1.0"})
            with urllib.request.urlopen(req, timeout=1.5) as resp:
                sdata = json.loads(resp.read().decode('utf-8'))
                for st in sdata.get("streams", []):
                    if st.get("stream_id") == cam_id and st.get("source_url"):
                        resolved_url = st.get("source_url")
                        if resolved_url.startswith("rtsp://") or resolved_url.startswith("rtmp://"):
                            source = resolved_url
                            is_rtsp = True
                            break
        except Exception:
            pass

    rtsp_cap = None
    if is_rtsp:
        os.environ["OPENCV_FFMPEG_CAPTURE_OPTIONS"] = "rtsp_transport;tcp|fflags;nobuffer|flags;low_delay"
        try:
            rtsp_cap = cv2.VideoCapture(source, cv2.CAP_FFMPEG)
            rtsp_cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)
        except Exception:
            rtsp_cap = None

    while True:
        try:
            frame = None
            # 1. Direct RTSP stream via OpenCV
            if is_rtsp:
                if rtsp_cap is None or not rtsp_cap.isOpened():
                    rtsp_cap = cv2.VideoCapture(source, cv2.CAP_FFMPEG)
                    rtsp_cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)
                ret, frame = rtsp_cap.read()
                if not ret or frame is None:
                    time.sleep(0.1)
                    continue

            # 2. HTTP/HTTPS snapshot or stream
            elif is_http:
                try:
                    req = urllib.request.Request(source, headers={"User-Agent": "HydraForge/1.0"})
                    with urllib.request.urlopen(req, timeout=1.5) as resp:
                        arr = np.frombuffer(resp.read(), dtype=np.uint8)
                        frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)
                except Exception:
                    frame = None

            # 3. HydraStream camera ID: live snapshot from HydraStream Data Plane
            if frame is None and (cam_id.startswith("cam") or cam_id.startswith("test_") or not os.path.exists(source)):
                live_url = f"http://127.0.0.1:8080/api/v1/streams/{cam_id}/snapshot.jpg"
                try:
                    req = urllib.request.Request(live_url, headers={"User-Agent": "HydraForge/1.0"})
                    with urllib.request.urlopen(req, timeout=1.5) as resp:
                        arr = np.frombuffer(resp.read(), dtype=np.uint8)
                        frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)
                except Exception:
                    frame = None

            # 4. Fallback to local files/samples
            if frame is None:
                source_path = source
                if not os.path.exists(source_path):
                    sample_candidates = []
                    if samples_env:
                        sample_candidates.extend([
                            os.path.join(samples_env, f"{source_path}.jpg"),
                            os.path.join(samples_env, source_path),
                            os.path.join(samples_env, "cam_01.jpg"),
                        ])
                    sample_candidates.extend([
                        os.path.join(base_dir, "../HydraStream/samples", f"{source_path}.jpg"),
                        os.path.join(base_dir, "../HydraStream/samples", source_path),
                        os.path.join(home_dir, "Documents/HydraStream/samples", f"{source_path}.jpg"),
                        os.path.join(home_dir, "Documents/HydraStream/samples", source_path),
                        os.path.join(home_dir, "Documents/HydraStream/samples/cam_01.jpg"),
                    ])
                    for sc in sample_candidates:
                        if os.path.exists(sc):
                            source_path = sc
                            break
                if os.path.exists(source_path):
                    frame = cv2.imread(source_path)

            if frame is None:
                time.sleep(0.04)
                continue

            # Save live frame for browser display
            try:
                cv2.imwrite(os.path.join(captures_dir, f"{cam_id}.jpg"), frame)
                cv2.imwrite(os.path.join(captures_dir, "latest.jpg"), frame)
            except Exception:
                pass

            t0 = time.perf_counter()
            if args.sahi:
                detections = run_sahi_inference(
                    model=model,
                    img=frame,
                    conf=args.conf,
                    iou=effective_iou,
                    device=device,
                    imgsz=640,
                    nms_free=args.nms_free
                )
                inf_ms = (time.perf_counter() - t0) * 1000.0
            else:
                results = model.track(
                    source=frame,
                    conf=args.conf,
                    iou=effective_iou,
                    device=device,
                    persist=True,
                    tracker="bytetrack.yaml",
                    verbose=False
                )
                inf_ms = (time.perf_counter() - t0) * 1000.0

                detections = []
                if len(results) > 0:
                    res = results[0]
                    for idx, box in enumerate(res.boxes):
                        cls_id = int(box.cls[0].item())
                        cls_name = res.names.get(cls_id, f"class_{cls_id}")
                        conf = float(box.conf[0].item())

                        if target_classes:
                            is_target = False
                            for tc in target_classes:
                                if tc in cls_name.lower() or (tc in ["pessoa", "person", "human"] and cls_name.lower() in ["person", "human"]):
                                    is_target = True
                                    break
                            if not is_target:
                                continue

                        xywhn = box.xywhn[0].tolist()
                        track_id = None
                        if box.id is not None:
                            track_id = int(box.id[0].item())

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

                        det_id = track_id if track_id is not None else (idx + 1)
                        detections.append({
                            "id": det_id,
                            "track_id": track_id,
                            "label": cls_name,
                            "class_name": cls_name,
                            "conf": round(conf, 3),
                            "confidence": round(conf, 3),
                            "box": [round(left_pct, 2), round(top_pct, 2), round(w_pct, 2), round(h_pct, 2)],
                            "color": color
                        })

            vram_mb = round(torch.cuda.memory_allocated(0) / (1024 * 1024), 1) if torch.cuda.is_available() else 0
            fps = round(1000.0 / inf_ms, 0) if inf_ms > 0 else 0

            payload = {
                "status": "success",
                "displayImageUrl": f"/api/v1/inference/frame?source={cam_id}&t={int(time.time()*1000)}",
                "image_url": f"/api/v1/inference/frame?source={cam_id}&t={int(time.time()*1000)}",
                "detections": detections,
                "count": len(detections),
                "telemetry": {
                    "inference_ms": f"{inf_ms:.2f}",
                    "fps": f"{fps:,.0f}",
                    "vram_mb": f"{vram_mb:,.0f}",
                    "device": "NVIDIA GeForce RTX 5090" if torch.cuda.is_available() else "CPU"
                }
            }
            print(json.dumps(payload), flush=True)
            time.sleep(0.035)

        except Exception as err:
            time.sleep(0.05)

if __name__ == "__main__":
    main()
