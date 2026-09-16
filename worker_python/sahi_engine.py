#!/usr/bin/env python3
"""
HydraForge SAHI (Slicing Aided Hyper Inference) Engine
Dynamic window slicing (640x640 @ 20% overlap) with batch GPU execution and Non-Maximum Merging (NMM).
"""
import numpy as np
import cv2

def calculate_slices(img_w, img_h, slice_w=640, slice_h=640, overlap_w=0.20, overlap_h=0.20):
    step_x = int(slice_w * (1.0 - overlap_w))
    step_y = int(slice_h * (1.0 - overlap_h))

    x_starts = list(range(0, max(1, img_w - slice_w + 1), max(1, step_x)))
    if not x_starts or x_starts[-1] + slice_w < img_w:
        x_starts.append(max(0, img_w - slice_w))

    y_starts = list(range(0, max(1, img_h - slice_h + 1), max(1, step_y)))
    if not y_starts or y_starts[-1] + slice_h < img_h:
        y_starts.append(max(0, img_h - slice_h))

    slices = []
    for y1 in y_starts:
        for x1 in x_starts:
            x2 = min(img_w, x1 + slice_w)
            y2 = min(img_h, y1 + slice_h)
            slices.append((x1, y1, x2, y2))
    return slices

def box_iou(box1, box2):
    # box format: [x1, y1, x2, y2]
    xA = max(box1[0], box2[0])
    yA = max(box1[1], box2[1])
    xB = min(box1[2], box2[2])
    yB = min(box1[3], box2[3])

    interArea = max(0, xB - xA) * max(0, yB - yA)
    if interArea == 0:
        return 0.0

    box1Area = (box1[2] - box1[0]) * (box1[3] - box1[1])
    box2Area = (box2[2] - box2[0]) * (box2[3] - box2[1])
    return interArea / float(box1Area + box2Area - interArea)

def non_max_merging(detections, iou_threshold=0.45):
    """
    Non-Maximum Merging: merges multiple slice detections for the same object
    weighted by confidence scores.
    """
    if not detections:
        return []

    by_class = {}
    for d in detections:
        cls = d["class_name"]
        by_class.setdefault(cls, []).append(d)

    merged_all = []

    for cls, dets in by_class.items():
        dets.sort(key=lambda x: x["conf"], reverse=True)
        used = [False] * len(dets)

        for i in range(len(dets)):
            if used[i]:
                continue
            used[i] = True
            main_det = dets[i]
            matched = [main_det]

            for j in range(i + 1, len(dets)):
                if used[j]:
                    continue
                if box_iou(main_det["xyxy"], dets[j]["xyxy"]) >= iou_threshold:
                    matched.append(dets[j])
                    used[j] = True

            if len(matched) == 1:
                merged_all.append(main_det)
            else:
                total_conf = sum(m["conf"] for m in matched)
                x1 = sum(m["xyxy"][0] * m["conf"] for m in matched) / total_conf
                y1 = sum(m["xyxy"][1] * m["conf"] for m in matched) / total_conf
                x2 = sum(m["xyxy"][2] * m["conf"] for m in matched) / total_conf
                y2 = sum(m["xyxy"][3] * m["conf"] for m in matched) / total_conf
                max_conf = max(m["conf"] for m in matched)

                merged_all.append({
                    "id": main_det["id"],
                    "track_id": main_det.get("track_id"),
                    "label": main_det["label"],
                    "class_name": main_det["class_name"],
                    "conf": max_conf,
                    "xyxy": [x1, y1, x2, y2],
                    "color": main_det["color"]
                })

    merged_all.sort(key=lambda x: x["conf"], reverse=True)
    return merged_all

def run_sahi_inference(model, img, conf=0.25, iou=0.45, device="0", imgsz=640, nms_free=False):
    """
    Runs SAHI high-resolution inference on an image matrix (BGR).
    """
    img_h, img_w = img.shape[:2]

    # If image is already smaller or equal to slice, run standard inference
    if img_w <= imgsz and img_h <= imgsz:
        results = model.predict(source=img, conf=conf, iou=(1.0 if nms_free else iou), device=device, imgsz=imgsz, verbose=False)
        return extract_standard_detections(results, img_w, img_h)

    # 1. Generate Slices
    slice_boxes = calculate_slices(img_w, img_h, slice_w=imgsz, slice_h=imgsz, overlap_w=0.20, overlap_h=0.20)
    slice_images = [img[y1:y2, x1:x2] for (x1, y1, x2, y2) in slice_boxes]

    # 2. Batch Inference on GPU for all slices
    slice_results = model.predict(
        source=slice_images,
        conf=conf,
        iou=(1.0 if nms_free else iou),
        device=device,
        imgsz=imgsz,
        batch=len(slice_images),
        verbose=False
    )

    all_detections = []
    det_counter = 1

    # Map slice detections to global coordinates
    for s_idx, res in enumerate(slice_results):
        x_offset, y_offset, _, _ = slice_boxes[s_idx]
        names = res.names
        for box in res.boxes:
            cls_id = int(box.cls[0].item())
            cls_name = names.get(cls_id, f"class_{cls_id}")
            c_score = float(box.conf[0].item())
            local_xyxy = box.xyxy[0].tolist()

            global_xyxy = [
                local_xyxy[0] + x_offset,
                local_xyxy[1] + y_offset,
                local_xyxy[2] + x_offset,
                local_xyxy[3] + y_offset,
            ]

            color = get_class_color(cls_name)
            all_detections.append({
                "id": det_counter,
                "track_id": None,
                "label": cls_name,
                "class_name": cls_name,
                "conf": c_score,
                "xyxy": global_xyxy,
                "color": color
            })
            det_counter += 1

    # 3. Full-Frame (Downscaled) Inference to catch large objects spanning multiple slices
    full_res = model.predict(
        source=img,
        conf=conf,
        iou=(1.0 if nms_free else iou),
        device=device,
        imgsz=imgsz,
        verbose=False
    )
    if len(full_res) > 0:
        names = full_res[0].names
        for box in full_res[0].boxes:
            cls_id = int(box.cls[0].item())
            cls_name = names.get(cls_id, f"class_{cls_id}")
            c_score = float(box.conf[0].item())
            global_xyxy = box.xyxy[0].tolist()
            all_detections.append({
                "id": det_counter,
                "track_id": None,
                "label": cls_name,
                "class_name": cls_name,
                "conf": c_score,
                "xyxy": global_xyxy,
                "color": get_class_color(cls_name)
            })
            det_counter += 1

    # 4. Non-Maximum Merging (NMM)
    merged = non_max_merging(all_detections, iou_threshold=iou)

    # 5. Format to percentage Bounding Boxes for UI Reticle
    formatted = []
    for idx, d in enumerate(merged):
        gx1, gy1, gx2, gy2 = d["xyxy"]
        w = max(1.0, gx2 - gx1)
        h = max(1.0, gy2 - gy1)

        left_pct = max(0.0, min(100.0, (gx1 / img_w) * 100.0))
        top_pct = max(0.0, min(100.0, (gy1 / img_h) * 100.0))
        w_pct = max(0.5, min(100.0, (w / img_w) * 100.0))
        h_pct = max(0.5, min(100.0, (h / img_h) * 100.0))

        formatted.append({
            "id": idx + 1,
            "track_id": d.get("track_id"),
            "label": d["label"],
            "class_name": d["class_name"],
            "conf": round(d["conf"], 3),
            "confidence": round(d["conf"], 3),
            "box": [round(left_pct, 2), round(top_pct, 2), round(w_pct, 2), round(h_pct, 2)],
            "color": d["color"]
        })

    return formatted

def extract_standard_detections(results, img_w, img_h):
    if not results or len(results) == 0:
        return []
    res = results[0]
    detections = []
    for idx, box in enumerate(res.boxes):
        cls_id = int(box.cls[0].item())
        cls_name = res.names.get(cls_id, f"class_{cls_id}")
        conf = float(box.conf[0].item())
        xywhn = box.xywhn[0].tolist()

        track_id = int(box.id[0].item()) if box.id is not None else None
        left_pct = max(0.0, min(100.0, (xywhn[0] - xywhn[2] / 2.0) * 100.0))
        top_pct = max(0.0, min(100.0, (xywhn[1] - xywhn[3] / 2.0) * 100.0))
        w_pct = max(0.5, min(100.0, xywhn[2] * 100.0))
        h_pct = max(0.5, min(100.0, xywhn[3] * 100.0))

        detections.append({
            "id": track_id if track_id is not None else (idx + 1),
            "track_id": track_id,
            "label": cls_name,
            "class_name": cls_name,
            "conf": round(conf, 3),
            "confidence": round(conf, 3),
            "box": [round(left_pct, 2), round(top_pct, 2), round(w_pct, 2), round(h_pct, 2)],
            "color": get_class_color(cls_name)
        })
    return detections

def get_class_color(cls_name):
    low = cls_name.lower()
    if any(k in low for k in ["bus", "onibus", "caminhao", "truck"]):
        return "#fcee0a"
    elif any(k in low for k in ["person", "moto", "pedestre"]):
        return "#00ff9d"
    elif "phone" in low or "celular" in low:
        return "#ff0055"
    return "#00f0ff"
