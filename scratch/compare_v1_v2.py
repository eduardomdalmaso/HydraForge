import cv2
import os
import time
from ultralytics import YOLO

os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

v1_path = "storage/models/person_detection_reinforced_v1_backup.pt"
v2_path = "storage/models/person_detection_reinforced.pt"

print(f"Loading Model V1 (Original): {v1_path}")
model_v1 = YOLO(v1_path)

print(f"Loading Model V2 (Fine-Tuned): {v2_path}")
model_v2 = YOLO(v2_path)

cameras = [
    ("CAM_01_LIVE (d38)", "rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1"),
    ("CAM_04_LIVE (g8q)", "rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt"),
]

for name, url in cameras:
    print(f"\n========================================================")
    print(f"Testing on {name}...")
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
    ret, frame = cap.read()
    cap.release()

    if not ret or frame is None:
        print(f"Failed to grab frame from {name}")
        continue

    # V1 Predict
    r1 = model_v1(frame, conf=0.20, verbose=False)
    dets_v1 = []
    if len(r1) > 0 and len(r1[0].boxes) > 0:
        for b in r1[0].boxes:
            c_name = r1[0].names[int(b.cls[0].item())]
            conf = float(b.conf[0].item())
            dets_v1.append(f"{c_name}: {conf*100:.1f}%")

    # V2 Predict
    r2 = model_v2(frame, conf=0.20, verbose=False)
    dets_v2 = []
    if len(r2) > 0 and len(r2[0].boxes) > 0:
        for b in r2[0].boxes:
            c_name = r2[0].names[int(b.cls[0].item())]
            conf = float(b.conf[0].item())
            dets_v2.append(f"{c_name}: {conf*100:.1f}%")

    print(f"  [MODELO V1 - ORIGINAL]   Detecções: {', '.join(dets_v1) if dets_v1 else 'Nenhuma'}")
    print(f"  [MODELO V2 - FINE-TUNED] Detecções: {', '.join(dets_v2) if dets_v2 else 'Nenhuma'}")
