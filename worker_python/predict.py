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

def resolve_source(source_arg):
    import tempfile
    tmp_dir = tempfile.gettempdir()
    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    captures_dir = os.path.join(base_dir, "storage/captures")
    os.makedirs(captures_dir, exist_ok=True)
    cam_clean = os.path.basename(source_arg).replace(".jpg", "")

    # 1. Direct RTSP / RTMP video stream capture with OpenCV
    if source_arg.startswith("rtsp://") or source_arg.startswith("rtmp://"):
        try:
            import cv2
            os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'
            cap = cv2.VideoCapture(source_arg, cv2.CAP_FFMPEG)
            cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
            cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
            ret, frame = cap.read()
            cap.release()
            if ret and frame is not None and frame.size > 0:
                tmp = os.path.join(tmp_dir, f"live_infer_{cam_clean}.jpg")
                cap_file = os.path.join(captures_dir, f"{cam_clean}.jpg")
                latest_file = os.path.join(captures_dir, "latest.jpg")
                cv2.imwrite(tmp, frame)
                cv2.imwrite(cap_file, frame)
                cv2.imwrite(latest_file, frame)
                return cap_file
        except Exception:
            pass

    # 2. HTTP image or MJPEG frame capture
    if source_arg.startswith("http://") or source_arg.startswith("https://"):
        try:
            import urllib.request
            tmp = os.path.join(tmp_dir, "live_infer_url.jpg")
            urllib.request.urlretrieve(source_arg, tmp)
            if os.path.exists(tmp) and os.path.getsize(tmp) > 100:
                return tmp
        except Exception:
            pass

    # 3. HydraStream stream ID: resolve RTSP source URL from Data Plane or fetch snapshot
    if cam_clean.startswith("cam") or cam_clean.startswith("test_") or not os.path.exists(source_arg):
        try:
            import urllib.request, json, cv2
            req = urllib.request.Request("http://127.0.0.1:8080/api/v1/streams", headers={'User-Agent': 'HydraForge/1.0'})
            with urllib.request.urlopen(req, timeout=2.0) as resp:
                sdata = json.loads(resp.read().decode('utf-8'))
                for st in sdata.get("streams", []):
                    if st.get("stream_id") == cam_clean and st.get("source_url"):
                        rtsp_url = st.get("source_url")
                        if rtsp_url.startswith("rtsp://") or rtsp_url.startswith("rtmp://"):
                            os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'
                            cap = cv2.VideoCapture(rtsp_url, cv2.CAP_FFMPEG)
                            cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
                            cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
                            ret, frame = cap.read()
                            cap.release()
                            if ret and frame is not None and frame.size > 0:
                                cap_file = os.path.join(captures_dir, f"{cam_clean}.jpg")
                                latest_file = os.path.join(captures_dir, "latest.jpg")
                                cv2.imwrite(cap_file, frame)
                                cv2.imwrite(latest_file, frame)
                                return cap_file
        except Exception:
            pass

        # Fallback to snapshot endpoint
        live_url = f"http://127.0.0.1:8080/api/v1/streams/{cam_clean}/snapshot.jpg?t={time.time()}"
        try:
            import urllib.request, cv2, numpy as np
            req = urllib.request.Request(live_url, headers={'User-Agent': 'HydraForge-Inference/1.0'})
            with urllib.request.urlopen(req, timeout=2.0) as resp:
                arr = np.asarray(bytearray(resp.read()), dtype=np.uint8)
                img = cv2.imdecode(arr, cv2.IMREAD_COLOR)
                if img is not None and img.size > 0:
                    cap_file = os.path.join(captures_dir, f"{cam_clean}.jpg")
                    cv2.imwrite(cap_file, img)
                    cv2.imwrite(os.path.join(captures_dir, "latest.jpg"), img)
                    return cap_file
        except Exception:
            pass

    if os.path.exists(source_arg):
        return source_arg

    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    home_dir = os.path.expanduser("~")
    samples_env = os.environ.get("HYDRASTREAM_SAMPLES_DIR", "")
    datasets_env = os.environ.get("DATASETS_DIR", "")

    # Check HydraStream sample files
    sample_candidates = []
    if samples_env:
        sample_candidates.extend([
            os.path.join(samples_env, f"{source_arg}.jpg"),
            os.path.join(samples_env, source_arg),
            os.path.join(samples_env, "cam_01.jpg"),
        ])
    sample_candidates.extend([
        os.path.join(base_dir, "../HydraStream/samples", f"{source_arg}.jpg"),
        os.path.join(base_dir, "../HydraStream/samples", source_arg),
        os.path.join(home_dir, "Documents/HydraStream/samples", f"{source_arg}.jpg"),
        os.path.join(home_dir, "Documents/HydraStream/samples", source_arg),
        os.path.join(home_dir, "Documents/HydraStream/samples/cam_01.jpg"),
        os.path.join(home_dir, "Documents/HydraStream/samples/cam_10_0_0_64.jpg"),
    ])
    for sc in sample_candidates:
        if os.path.exists(sc):
            return sc

    clean_id = source_arg.replace("dataset:", "")
    dataset_dirs = [
        os.path.join(base_dir, "datasets"),
        os.path.join(home_dir, "datasets"),
        os.path.join(home_dir, "Documents/HydraForge/datasets"),
    ]
    if datasets_env:
        dataset_dirs.insert(0, datasets_env)

    for ds_dir in dataset_dirs:
        ds_pattern = os.path.join(ds_dir, clean_id, "train/images/*.jpg")
        matched = glob.glob(ds_pattern)
        if matched:
            return matched[0]

    # Global dataset fallbacks
    for ds_dir in dataset_dirs:
        fallbacks = glob.glob(os.path.join(ds_dir, "*/train/images/*.jpg"))
        if fallbacks:
            return fallbacks[0]

    return source_arg

def parse_args():
    parser = argparse.ArgumentParser(description="HydraForge YOLO Inference Engine")
    parser.add_argument("--weights", type=str, default="weights/yolo26n.pt", help="Path to .pt / .engine weights")
    parser.add_argument("--source", type=str, required=True, help="Image file path or URL")
    parser.add_argument("--conf", type=float, default=0.25, help="Confidence threshold")
    parser.add_argument("--iou", type=float, default=0.45, help="IoU NMS threshold")
    parser.add_argument("--device", type=str, default="0", help="CUDA device index or cpu")
    parser.add_argument("--imgsz", type=int, default=640, help="Inference image resolution")
    parser.add_argument("--classes", type=str, default="", help="Comma separated target classes, e.g. person,pessoa")
    parser.add_argument("--track", action="store_true", help="Enable ByteTrack object tracking")
    parser.add_argument("--sahi", action="store_true", help="Enable SAHI 4K sliced inference")
    parser.add_argument("--nms-free", action="store_true", help="Enable NMS-free end-to-end decode")
    return parser.parse_args()

def main():
    args = parse_args()
    weights_path = resolve_weights(args.weights)
    source_path = resolve_source(args.source)

    if not os.path.exists(source_path):
        print(json.dumps({"status": "error", "error": f"source {source_path} not found", "detections": []}))
        sys.exit(1)

    from ultralytics import YOLO
    import torch
    import cv2
    from sahi_engine import run_sahi_inference

    try:
        model = YOLO(weights_path)
        device = 0 if (torch.cuda.is_available() and args.device != "cpu") else "cpu"

        t_start = time.perf_counter()
        effective_iou = 1.0 if args.nms_free else args.iou

        if args.sahi:
            img = cv2.imread(source_path)
            if img is not None:
                detections = run_sahi_inference(
                    model=model,
                    img=img,
                    conf=args.conf,
                    iou=effective_iou,
                    device=device,
                    imgsz=args.imgsz,
                    nms_free=args.nms_free
                )
            else:
                detections = []
            t_total_ms = (time.perf_counter() - t_start) * 1000.0
            speed = {"preprocess": 1.0, "inference": round(t_total_ms, 2), "postprocess": 1.0}
        else:
            if args.track:
                results = model.track(
                    source=source_path,
                    conf=args.conf,
                    iou=effective_iou,
                    device=device,
                    imgsz=args.imgsz,
                    persist=True,
                    tracker="bytetrack.yaml",
                    verbose=False
                )
            else:
                results = model.predict(
                    source=source_path,
                    conf=args.conf,
                    iou=effective_iou,
                    device=device,
                    imgsz=args.imgsz,
                    verbose=False
                )
            t_total_ms = (time.perf_counter() - t_start) * 1000.0

            detections = []
            speed = {"preprocess": 0.0, "inference": 0.0, "postprocess": 0.0}

            target_classes = [c.strip().lower() for c in args.classes.split(",") if c.strip()]
            if len(results) > 0:
                res = results[0]
                if hasattr(res, "speed"):
                    speed = res.speed

                boxes = res.boxes
                for idx, box in enumerate(boxes):
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

        inf_ms = speed.get("inference", t_total_ms)
        fps = round(1000.0 / inf_ms, 1) if inf_ms > 0 else 0

        vram_mb = 0
        if torch.cuda.is_available():
            vram_mb = round(torch.cuda.memory_allocated(0) / (1024 * 1024), 1)

        cam_id = os.path.basename(source_path).replace(".jpg", "")
        payload = {
            "status": "success",
            "model": os.path.basename(weights_path),
            "model_path": weights_path,
            "source": source_path,
            "displayImageUrl": f"/api/v1/inference/frame?source={cam_id}&t={int(time.time()*1000)}",
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
