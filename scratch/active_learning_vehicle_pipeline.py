import os
import cv2
import glob
import json
import time
import shutil
from typing import Any, cast
import numpy as np
from ultralytics import YOLO

def run_pipeline():
    video_path = "C:\\Users\\hades\\Documents\\HydraForge\\hiv00022.mp4"
    base_model_path = "storage/models/vehicle_detection.pt"
    backup_model_path = "storage/models/golden_backup/vehicle_detection_v1_base.pt"
    
    os.makedirs("storage/models/golden_backup", exist_ok=True)
    if not os.path.exists(backup_model_path) and os.path.exists(base_model_path):
        shutil.copyfile(base_model_path, backup_model_path)
        print(f"Backed up base model to: {backup_model_path}")

    reinforced_dir = "datasets/vehicles_reinforced"
    os.makedirs(f"{reinforced_dir}/images/train", exist_ok=True)
    os.makedirs(f"{reinforced_dir}/labels/train", exist_ok=True)
    os.makedirs(f"{reinforced_dir}/images/val", exist_ok=True)
    os.makedirs(f"{reinforced_dir}/labels/val", exist_ok=True)
    
    vis_dir = "storage/audit_vehicles/hard_cases_visual"
    os.makedirs(vis_dir, exist_ok=True)

    print("==========================================================")
    print("=== HYDRA VISION // TESLA ACTIVE LEARNING & TUNING LOOP ===")
    print("==========================================================")
    print(f"Target Video: {video_path}")
    print(f"Base Model: {base_model_path}")
    print(f"Hardware: NVIDIA GeForce RTX 5090 (CUDA:0)")

    # ----------------------------------------------------
    # STEP 1: Copy base dataset into vehicles_reinforced
    # ----------------------------------------------------
    base_dir = "datasets/vehicles_unified"
    print("\n[STEP 1/4] Assembling base dataset foundation...")
    for split in ["train", "val"]:
        img_srcs = glob.glob(f"{base_dir}/images/{split}/*.*")
        for img_p in img_srcs:
            fname = os.path.basename(img_p)
            dest_img = f"{reinforced_dir}/images/{split}/{fname}"
            if not os.path.exists(dest_img):
                shutil.copyfile(img_p, dest_img)
            
            base_name = os.path.splitext(fname)[0]
            lbl_p = f"{base_dir}/labels/{split}/{base_name}.txt"
            dest_lbl = f"{reinforced_dir}/labels/{split}/{base_name}.txt"
            if os.path.exists(lbl_p) and not os.path.exists(dest_lbl):
                shutil.copyfile(lbl_p, dest_lbl)

    print(f"Base train images loaded: {len(glob.glob(f'{reinforced_dir}/images/train/*.*'))}")
    print(f"Base val images loaded: {len(glob.glob(f'{reinforced_dir}/images/val/*.*'))}")

    # ----------------------------------------------------
    # STEP 2: Harvest Unique Hard Cases from Video
    # ----------------------------------------------------
    print("\n[STEP 2/4] Harvesting unique hard cases (< 45% conf) from video...")
    model = YOLO(base_model_path)
    cap = cv2.VideoCapture(video_path)
    stride = 3
    max_hard_target = 60

    vehicles: dict[int, dict[str, Any]] = {}
    frame_idx = 0
    class_names = {0: 'car', 1: 'motorcycle', 2: 'bus', 3: 'truck'}

    while cap.isOpened():
        ret, frame = cap.read()
        if not ret:
            break

        frame_idx += 1
        if frame_idx % stride != 0:
            continue

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
                        "worst_frame": None,
                        "worst_bbox": None
                    }

                v = vehicles[tid]
                v["confs"].append(c_val)
                v["classes"][cname] = v["classes"].get(cname, 0) + 1

                if c_val <= v["min_conf"] or v["worst_frame"] is None:
                    v["min_conf"] = c_val
                    v["worst_frame"] = frame.copy()
                    v["worst_bbox"] = [int(c) for c in box]

        hard_so_far = sum(1 for v in vehicles.values() if len(v["confs"]) >= 3 and v["min_conf"] < 0.45)
        if hard_so_far >= max_hard_target:
            print(f"Reached harvest target of {hard_so_far} unique hard cases at frame {frame_idx}.")
            break

    cap.release()

    # Filter valid vehicles and generate auto-labeled samples
    valid_hard = {tid: v for tid, v in vehicles.items() if len(v["confs"]) >= 3 and v["min_conf"] < 0.45}
    print(f"Total Unique Hard-Case Vehicles Identified: {len(valid_hard)}")

    injected_count = 0
    h_w = None

    for tid, v in valid_hard.items():
        clean_img = v["worst_frame"]
        if clean_img is None:
            continue

        h, w = clean_img.shape[:2]
        h_w = (h, w)
        assigned_class = max(v["classes"], key=v["classes"].get)

        # Run high-sensitivity inference on the clean frame to label all vehicles
        preds = model.predict(source=clean_img, conf=0.20, device=0, verbose=False)[0]

        # Distribute: 80% to train, 20% to val (for honest tuning validation)
        split = "val" if (injected_count % 5 == 0) else "train"
        base_id = f"cctv_hard_id_{tid}_{assigned_class}"
        img_out = f"{reinforced_dir}/images/{split}/{base_id}.jpg"
        lbl_out = f"{reinforced_dir}/labels/{split}/{base_id}.txt"

        cv2.imwrite(img_out, clean_img)

        with open(lbl_out, "w") as f_lbl:
            if preds.boxes is not None:
                for b in preds.boxes:
                    bx_norm = b.xywhn[0].tolist()
                    b_cls = int(b.cls[0].item())
                    f_lbl.write(f"{b_cls} {bx_norm[0]:.6f} {bx_norm[1]:.6f} {bx_norm[2]:.6f} {bx_norm[3]:.6f}\n")

        # Save visual debug
        vis_img = clean_img.copy()
        bx = v["worst_bbox"]
        min_c = v["min_conf"]
        cv2.rectangle(vis_img, (bx[0], bx[1]), (bx[2], bx[3]), (0, 0, 255), 2)
        cv2.putText(
            vis_img,
            f"ID:{tid} {assigned_class} {min_c:.1%}",
            (bx[0], max(bx[1] - 8, 15)),
            cv2.FONT_HERSHEY_SIMPLEX,
            0.6,
            (0, 0, 255),
            2
        )
        cv2.imwrite(os.path.join(vis_dir, f"{base_id}_{int(min_c*100)}pct.jpg"), vis_img)
        injected_count += 1

    print(f"Injected {injected_count} clean hard-case frames into dataset ({reinforced_dir}).")

    # Create data.yaml
    data_yaml_path = f"{reinforced_dir}/data.yaml"
    with open(data_yaml_path, "w") as f:
        f.write(f"path: {os.path.abspath(reinforced_dir).replace(chr(92), '/')}\n")
        f.write("train: images/train\n")
        f.write("val: images/val\n\n")
        f.write("names:\n")
        f.write("  0: car\n")
        f.write("  1: motorcycle\n")
        f.write("  2: bus\n")
        f.write("  3: truck\n")

    # ----------------------------------------------------
    # STEP 3: Run Tuning (Active Learning Fine-Tuning)
    # ----------------------------------------------------
    print("\n[STEP 3/4] Launching Fine-Tuning on RTX 5090 (freeze=10, lr0=0.0003)...")
    tune_model = YOLO(base_model_path)
    
    tune_results = tune_model.train(
        data=data_yaml_path,
        epochs=10,
        imgsz=640,
        batch=32,
        device=0,
        workers=0,
        freeze=10,
        optimizer='AdamW',
        lr0=0.0003,
        lrf=0.01,
        warmup_epochs=1,
        project='runs/vehicles_reinforced',
        name='v1_active_learning',
        exist_ok=True,
        verbose=True
    )

    new_weights = "runs/vehicles_reinforced/v1_active_learning/weights/best.pt"
    if not os.path.exists(new_weights):
        # Check detect subdirectory if nested
        alt_weights = "runs/detect/runs/vehicles_reinforced/v1_active_learning/weights/best.pt"
        if os.path.exists(alt_weights):
            new_weights = alt_weights

    print(f"\n[STEP 4/4] Evaluating Baseline vs Tuned Model...")
    
    # 1. Validate Baseline on Reinforced Val Set
    print("\n--- Validating BASELINE Model ---")
    base_val = model.val(data=data_yaml_path, device=0, verbose=False)
    base_map = float(base_val.box.map50)

    # 2. Validate Tuned Model on Reinforced Val Set
    print("\n--- Validating TUNED Reinforced Model ---")
    tuned_eval = YOLO(new_weights)
    tuned_val = tuned_eval.val(data=data_yaml_path, device=0, verbose=False)
    tuned_map = float(tuned_val.box.map50)

    # 3. Test on Hard Cases Directory specifically
    hard_img_paths = glob.glob(f"{vis_dir}/*.jpg")
    base_confs = []
    tuned_confs = []

    for hp in hard_img_paths[:20]:
        im = cv2.imread(hp)
        if im is None:
            continue
        p_base = model.predict(source=im, conf=0.15, device=0, verbose=False)[0]
        p_tuned = tuned_eval.predict(source=im, conf=0.15, device=0, verbose=False)[0]
        
        if len(p_base.boxes) > 0:
            base_confs.append(float(np.mean(p_base.boxes.conf.cpu().numpy())))
        if len(p_tuned.boxes) > 0:
            tuned_confs.append(float(np.mean(p_tuned.boxes.conf.cpu().numpy())))

    avg_base_conf = float(np.mean(base_confs)) if base_confs else 0.0
    avg_tuned_conf = float(np.mean(tuned_confs)) if tuned_confs else 0.0

    print("\n==========================================================")
    print("=== ACTIVE LEARNING BENCHMARK // COMPARATIVE VERDICT ===")
    print("==========================================================")
    print(f"Baseline Model mAP50:  {base_map * 100:.2f}%")
    print(f"Tuned Model mAP50:     {tuned_map * 100:.2f}%")
    diff_map = (tuned_map - base_map) * 100
    print(f"Global mAP Delta:      {'+' if diff_map >= 0 else ''}{diff_map:.2f}%")
    print("----------------------------------------------------------")
    print(f"Hard-Cases Mean Conf (Baseline): {avg_base_conf * 100:.1f}%")
    print(f"Hard-Cases Mean Conf (Tuned):    {avg_tuned_conf * 100:.1f}%")
    diff_conf = (avg_tuned_conf - avg_base_conf) * 100
    print(f"Hard-Cases Gain:                 {'+' if diff_conf >= 0 else ''}{diff_conf:.1f}%")
    print("----------------------------------------------------------")

    improved = (tuned_map >= (base_map - 0.01)) and (avg_tuned_conf >= avg_base_conf)
    if improved:
        print("VERDICT: PROMOTION APPROVED! Model successfully reinforced without regression.")
        shutil.copyfile(new_weights, base_model_path)
        print(f"Promoted {new_weights} -> {base_model_path}")
    else:
        print("VERDICT: BASELINE RETAINED (No significant lift or regression detected).")
    print("==========================================================")

if __name__ == "__main__":
    run_pipeline()
