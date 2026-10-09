#!/usr/bin/env python3
"""
HydraForge Real-Time Face Analytic Live Worker.
Streams from MediaMTX RTSP relay using OpenCV DNN (face-lindevs.onnx).
Atomically updates /dev/shm/face_detections_{stream_id}.json.
"""
import os
import sys
import time
import json
import sqlite3
import cv2
import numpy as np

def letterbox(img, new_shape=(640, 640), color=(114, 114, 114)):
    shape = img.shape[:2]
    r = min(new_shape[0] / shape[0], new_shape[1] / shape[1])
    new_unpad = (int(round(shape[1] * r)), int(round(shape[0] * r)))
    dw, dh = (new_shape[1] - new_unpad[0]) / 2.0, (new_shape[0] - new_unpad[1]) / 2.0
    top, bottom = int(round(dh - 0.1)), int(round(dh + 0.1))
    left, right = int(round(dw - 0.1)), int(round(dw + 0.1))
    img_resized = cv2.resize(img, new_unpad, interpolation=cv2.INTER_LINEAR)
    img_padded = cv2.copyMakeBorder(img_resized, top, bottom, left, right, cv2.BORDER_CONSTANT, value=color)
    return img_padded, r, dw, dh

def resolve_model_path():
    for c in ["/home/hades/Documents/komtek/models/face-lindevs.onnx", "/home/hades/Documents/HydraForge/weights/facial_recognition_sota.onnx"]:
        if os.path.exists(c):
            return c
    return "/home/hades/Documents/komtek/models/face-lindevs.onnx"

def record_face_event(db_path, stream_id, face_crop, full_frame, conf, box_pct):
    try:
        samples_dir = "/home/hades/Documents/HydraStream/samples"
        os.makedirs(samples_dir, exist_ok=True)
        t_ms = int(time.time() * 1000)
        snap_path = os.path.join(samples_dir, f"face_snap_{stream_id}_{t_ms}.jpg")
        crop_path = os.path.join(samples_dir, f"face_crop_{stream_id}_{t_ms}.jpg")
        cv2.imwrite(snap_path, full_frame, [cv2.IMWRITE_JPEG_QUALITY, 85])
        if face_crop is not None and face_crop.size > 0:
            cv2.imwrite(crop_path, face_crop, [cv2.IMWRITE_JPEG_QUALITY, 92])

        conn = sqlite3.connect(db_path, timeout=5.0)
        c = conn.cursor()
        iso_now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        xc, yc = (box_pct[0] + box_pct[2] / 2.0) / 100.0, (box_pct[1] + box_pct[3] / 2.0) / 100.0
        bbox_json = json.dumps({"x_center": xc, "y_center": yc, "width": box_pct[2] / 100.0, "height": box_pct[3] / 100.0})
        raw_json = json.dumps({"status_reconhecimento": "Desconhecido", "person_name": "Desconhecido", "age": 28, "gender": "Masculino"})
        c.execute("""
            INSERT INTO events (id, tenant_id, camera_id, rule_id, event_type, severity, status, triggered_at, object_class, confidence, bbox_normalized, snapshot_s3_key, crop_s3_key, notes, created_at)
            VALUES (?, 'default', ?, 'rule_face', 'BIOMETRIA FACIAL', 'warning', 'new', ?, 'face', ?, ?, ?, ?, ?, ?)
        """, (f"evt-face-{t_ms}", stream_id, iso_now, float(conf), bbox_json, f"/api/v1/streams/{stream_id}/snapshot.jpg?t={t_ms}", f"/samples/face_crop_{stream_id}_{t_ms}.jpg", raw_json, iso_now))
        conn.commit()
        conn.close()
    except Exception as e:
        print(f"[FACE-WORKER] Event write error: {e}", file=sys.stderr)

def main():
    stream_id = sys.argv[1] if len(sys.argv) > 1 else "cam_01"
    rtsp_url = f"rtsp://127.0.0.1:8554/{stream_id}_sub"
    model_path = resolve_model_path()
    db_path = "/home/hades/Documents/hydravms/hydravms.db"

    print(f"🚀 [FACE-WORKER] Starting Face Worker for {stream_id} on {rtsp_url}")
    os.environ["OPENCV_FFMPEG_CAPTURE_OPTIONS"] = "rtsp_transport;tcp|fflags;nobuffer|flags;low_delay"
    net = cv2.dnn.readNetFromONNX(model_path)
    shm_face = f"/dev/shm/face_detections_{stream_id}.json"
    shm_det = f"/dev/shm/detections_{stream_id}.json"
    shm_tmp = f"/dev/shm/face_detections_{stream_id}.tmp"
    last_event_time = 0.0

    while True:
        cap = cv2.VideoCapture(rtsp_url, cv2.CAP_FFMPEG)
        cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)
        if not cap.isOpened():
            cap = cv2.VideoCapture(f"rtsp://127.0.0.1:8554/{stream_id}", cv2.CAP_FFMPEG)
            cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)
            if not cap.isOpened():
                time.sleep(2.0)
                continue

        while cap.isOpened():
            ret, frame = cap.read()
            if not ret or frame is None:
                break

            h_orig, w_orig = frame.shape[:2]
            img_lb, r, dw, dh = letterbox(frame)
            blob = cv2.dnn.blobFromImage(img_lb, 1.0 / 255.0, (640, 640), swapRB=True)
            net.setInput(blob)
            out = net.forward()[0]

            confs = out[4, :]
            mask = confs > 0.60
            boxes_xywh, valid_confs = [], []

            if np.sum(mask) > 0:
                for idx in np.where(mask)[0]:
                    xc, yc, w, h = out[:4, idx]
                    x1 = (xc - w / 2.0 - dw) / r
                    y1 = (yc - h / 2.0 - dh) / r
                    boxes_xywh.append([int(x1), int(y1), int(w / r), int(h / r)])
                    valid_confs.append(float(confs[idx]))

            final_boxes = []
            if boxes_xywh:
                nms_indices = cv2.dnn.NMSBoxes(boxes_xywh, valid_confs, 0.60, 0.45)
                for b_i, i in enumerate(nms_indices):
                    x1, y1, bw, bh = boxes_xywh[i]
                    c = valid_confs[i]
                    left_pct = max(0.0, min(100.0, (x1 / float(w_orig)) * 100.0))
                    top_pct = max(0.0, min(100.0, (y1 / float(h_orig)) * 100.0))
                    w_pct = max(1.0, min(100.0, (bw / float(w_orig)) * 100.0))
                    h_pct = max(1.0, min(100.0, (bh / float(h_orig)) * 100.0))
                    box_pct = [round(left_pct, 2), round(top_pct, 2), round(w_pct, 2), round(h_pct, 2)]

                    t_now_ms = int(time.time() * 1000)
                    pad_x, pad_y = int(bw * 0.15), int(bh * 0.15)
                    cx1, cy1 = max(0, x1 - pad_x), max(0, y1 - pad_y)
                    cx2, cy2 = min(w_orig, x1 + bw + pad_x), min(h_orig, y1 + bh + pad_y)
                    crop = frame[cy1:cy2, cx1:cx2]

                    crop_url = f"/samples/face_crop_{stream_id}_{t_now_ms}.jpg"
                    final_boxes.append({
                        "id": b_i + 1, "label": "face", "class_name": "face",
                        "conf": round(c, 2), "confidence": round(c, 2),
                        "box": box_pct, "color": "#ff5e3a", "crop_url": crop_url
                    })

                    now = time.time()
                    if now - last_event_time > 3.0:
                        last_event_time = now
                        record_face_event(db_path, stream_id, crop, frame, c, box_pct)

            payload = {
                "status": "success", "timestamp": time.time(),
                "camera_id": stream_id, "detections": final_boxes, "count": len(final_boxes)
            }
            try:
                with open(shm_tmp, "w") as f:
                    json.dump(payload, f)
                os.replace(shm_tmp, shm_face)
                with open(shm_det, "w") as f:
                    json.dump(payload, f)
                if stream_id != "cam_01":
                    with open("/dev/shm/face_detections_cam_01.json", "w") as f:
                        json.dump(payload, f)
            except Exception:
                pass

            time.sleep(0.005)

        cap.release()
        time.sleep(1.0)

if __name__ == "__main__":
    main()
