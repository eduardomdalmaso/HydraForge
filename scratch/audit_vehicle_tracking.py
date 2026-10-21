import os
import cv2
import json
import time
from typing import Any, cast
import numpy as np
from ultralytics import YOLO

def run_vehicle_audit(
    video_path: str,
    model_path: str = "storage/models/vehicle_detection.pt",
    stride: int = 3,
    max_frames: int | None = None,
    max_hard_cases: int = 50
) -> None:
    if not os.path.exists(video_path):
        print(f"Error: Video not found at {video_path}")
        return

    output_dir = "storage/audit_vehicles"
    hard_cases_dir = os.path.join(output_dir, "hard_cases_below_45")
    os.makedirs(hard_cases_dir, exist_ok=True)

    print("==================================================")
    print("=== HYDRA VISION // VEHICLE TRACKING & AUDITOR ===")
    print("==================================================")
    print(f"Video Source: {video_path}")
    print(f"Model: {model_path}")
    print(f"Hardware: NVIDIA GeForce RTX 5090 (CUDA:0)")
    print(f"Frame Stride: 1 frame every {stride} frames (Zero redundant tracking)")
    print(f"Target Hard Cases Limit: {max_hard_cases} unique frames (stops when reached)")

    # Load Model
    model = YOLO(model_path)

    cap = cv2.VideoCapture(video_path)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    video_fps = float(cap.get(cv2.CAP_PROP_FPS) or 15.0)
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))

    print(f"Resolution: {width}x{height} @ {video_fps:.1f} FPS | Total Video Frames: {total_frames:,}")

    # Track metadata per unique vehicle ID
    vehicles: dict[int, dict[str, Any]] = {}

    frame_idx = 0
    processed_count = 0
    start_time = time.time()

    class_names = {0: 'car', 1: 'motorcycle', 2: 'bus', 3: 'truck'}

    while cap.isOpened():
        ret, frame = cap.read()
        if not ret:
            break

        frame_idx += 1
        if frame_idx % stride != 0:
            continue

        processed_count += 1
        if max_frames and processed_count > max_frames:
            print(f"Reached max_frames limit ({max_frames}). Stopping scan.")
            break

        # Check if we already accumulated enough unique hard cases
        current_hard_count = sum(1 for v in vehicles.values() if len(v["confs"]) >= 3 and v["min_conf"] < 0.45)
        if current_hard_count >= max_hard_cases:
            print(f"\n[TARGET REACHED] Reached {current_hard_count} unique hard cases. Halting extraction.")
            break

        # Run tracking with ByteTrack
        raw_results = model.track(
            frame,
            device=0,
            tracker="bytetrack.yaml",
            conf=0.25,
            iou=0.6,
            persist=True,
            verbose=False
        )

        res_list = list(raw_results)
        if not res_list:
            continue
        res = res_list[0]

        boxes_obj = res.boxes
        if boxes_obj is not None and boxes_obj.id is not None:
            raw_boxes = cast(Any, boxes_obj.xyxy)
            raw_ids = cast(Any, boxes_obj.id)
            raw_confs = cast(Any, boxes_obj.conf)
            raw_clss = cast(Any, boxes_obj.cls)

            boxes_arr = np.asarray(raw_boxes.cpu() if hasattr(raw_boxes, 'cpu') else raw_boxes)
            track_ids_arr = np.asarray(raw_ids.cpu() if hasattr(raw_ids, 'cpu') else raw_ids, dtype=int)
            confs_arr = np.asarray(raw_confs.cpu() if hasattr(raw_confs, 'cpu') else raw_confs, dtype=float)
            clss_arr = np.asarray(raw_clss.cpu() if hasattr(raw_clss, 'cpu') else raw_clss, dtype=int)

            for box, track_id, conf, cls_id in zip(boxes_arr, track_ids_arr, confs_arr, clss_arr):
                cname = class_names.get(int(cls_id), str(cls_id))
                tid = int(track_id)
                c_val = float(conf)

                if tid not in vehicles:
                    vehicles[tid] = {
                        "classes": {},
                        "confs": [],
                        "min_conf": c_val,
                        "max_conf": c_val,
                        "worst_frame": None,
                        "worst_bbox": None,
                        "first_seen_frame": frame_idx,
                        "last_seen_frame": frame_idx
                    }

                v = vehicles[tid]
                v["confs"].append(c_val)
                v["classes"][cname] = v["classes"].get(cname, 0) + 1
                v["last_seen_frame"] = frame_idx
                if c_val > v["max_conf"]:
                    v["max_conf"] = c_val

                # Track worst frame (lowest confidence)
                if c_val <= v["min_conf"] or v["worst_frame"] is None:
                    v["min_conf"] = c_val
                    v["worst_frame"] = frame.copy()
                    v["worst_bbox"] = [int(coord) for coord in box]

        if processed_count % 300 == 0:
            elapsed = time.time() - start_time
            fps = processed_count / elapsed if elapsed > 0 else 0
            pct = (frame_idx / total_frames) * 100 if total_frames > 0 else 0
            hard_so_far = sum(1 for v in vehicles.values() if len(v["confs"]) >= 3 and v["min_conf"] < 0.45)
            print(f"[{pct:5.1f}%] Frame {frame_idx:6d}/{total_frames} | Speed: {fps:5.1f} FPS | Unique: {len(vehicles)} | Hard Cases: {hard_so_far}/{max_hard_cases}", flush=True)

    cap.release()
    total_time = time.time() - start_time
    avg_fps = processed_count / total_time if total_time > 0 else 0

    print("\n==================================================")
    print("=== AUDIT COMPLETED // COMPILING VEHICLE REPORT ===")
    print("==================================================")
    print(f"Processed Frames: {processed_count:,} (Sampled from {frame_idx:,})")
    print(f"Execution Time: {total_time:.1f}s | Average Speed: {avg_fps:.1f} FPS")

    # Filter out spurious single-frame noise (vehicles must appear in at least 3 sampled frames)
    valid_vehicles = {tid: v for tid, v in vehicles.items() if len(v["confs"]) >= 3}
    spurious_tracks = len(vehicles) - len(valid_vehicles)

    class_counts = {"car": 0, "motorcycle": 0, "bus": 0, "truck": 0}
    high_conf_vehicles = 0
    medium_conf_vehicles = 0
    hard_cases_count = 0

    export_records = []

    for tid, v in valid_vehicles.items():
        assigned_class = max(v["classes"], key=v["classes"].get)
        class_counts[assigned_class] = class_counts.get(assigned_class, 0) + 1

        mean_c = float(np.mean(v["confs"]))
        min_c = float(v["min_conf"])
        max_c = float(v["max_conf"])

        if mean_c >= 0.80:
            high_conf_vehicles += 1
        elif mean_c >= 0.50:
            medium_conf_vehicles += 1

        is_hard_case = min_c < 0.45
        if is_hard_case:
            hard_cases_count += 1
            if v["worst_frame"] is not None:
                img_to_save = v["worst_frame"].copy()
                bx = v["worst_bbox"]
                cv2.rectangle(img_to_save, (bx[0], bx[1]), (bx[2], bx[3]), (0, 0, 255), 2)
                cv2.putText(
                    img_to_save,
                    f"ID:{tid} {assigned_class} {min_c:.1%}",
                    (bx[0], max(bx[1] - 8, 15)),
                    cv2.FONT_HERSHEY_SIMPLEX,
                    0.6,
                    (0, 0, 255),
                    2
                )
                save_name = f"hard_id_{tid}_{assigned_class}_{int(min_c*100)}pct.jpg"
                cv2.imwrite(os.path.join(hard_cases_dir, save_name), img_to_save)

        export_records.append({
            "track_id": int(tid),
            "class": str(assigned_class),
            "frames_tracked": len(v["confs"]),
            "mean_conf": float(round(mean_c, 3)),
            "min_conf": float(round(min_c, 3)),
            "max_conf": float(round(max_c, 3)),
            "hard_case": bool(is_hard_case)
        })

    total_valid = len(valid_vehicles)
    high_pct = float(high_conf_vehicles / total_valid * 100) if total_valid > 0 else 0.0
    hard_pct = float(hard_cases_count / total_valid * 100) if total_valid > 0 else 0.0

    report = {
        "video": os.path.basename(video_path),
        "total_unique_vehicles": total_valid,
        "spurious_tracks_filtered": spurious_tracks,
        "classes": class_counts,
        "confidence_tiers": {
            "high_conf_ge_80pct": high_conf_vehicles,
            "medium_conf_50_to_80pct": medium_conf_vehicles,
            "hard_cases_lt_45pct": hard_cases_count,
            "pct_above_80": float(round(high_pct, 1)),
            "pct_hard_cases": float(round(hard_pct, 1))
        },
        "vehicles": export_records
    }

    report_path = os.path.join(output_dir, f"audit_{os.path.basename(video_path)}.json")
    with open(report_path, "w", encoding="utf-8") as f:
        json.dump(report, f, indent=2, default=str)

    print("\n---------------- AUDIT RESULTS ----------------")
    print(f"Total Unique Vehicles Verified: {total_valid}")
    for cname, cnt in class_counts.items():
        print(f" - {cname.upper():12s}: {cnt}")
    print("---------------- CONFIDENCE TIERS -------------")
    print(f"High Confidence (>= 80%): {high_conf_vehicles} vehicles ({high_pct:.1f}%)")
    print(f"Medium Confidence (50-79%): {medium_conf_vehicles} vehicles")
    print(f"Hard Cases (< 45%): {hard_cases_count} vehicles ({hard_pct:.1f}%)")
    print(f"Unique Hard Case Frames Saved to: {hard_cases_dir}")
    print(f"Full JSON Report Saved to: {report_path}")
    print("==================================================")

if __name__ == "__main__":
    import sys
    video = sys.argv[1] if len(sys.argv) > 1 else "C:\\Users\\hades\\Documents\\HydraForge\\hiv00022.mp4"
    run_vehicle_audit(video, max_hard_cases=50)
