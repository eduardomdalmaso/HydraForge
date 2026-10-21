import cv2
import os
import time
from ultralytics import YOLO

os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

v1_path = "storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt"
v2_path = "runs/detect/runs/tesla_data_engine/v1_hard_cases/weights/best.pt"

print(f"=== COMPARING GOLDEN 89% vs TESLA HARD-CASE FINE-TUNED ===")
print(f"Model V1 (Golden 89%): {v1_path}")
print(f"Model V2 (Tesla Candidate): {v2_path}")

model_v1 = YOLO(v1_path)
model_v2 = YOLO(v2_path)

cameras = [
    ("CAM_01_LIVE (d38)", "rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1"),
    ("CAM_04_LIVE (g8q)", "rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt"),
    ("CAM_02_LIVE (v0m)", "rtsp://k-monitorcenter.app:8554/live/live_v0m7gz9przaz"),
    ("CAM_05_LIVE (k4x)", "rtsp://k-monitorcenter.app:8554/live/live_k4xvpbtdqji8"),
    ("CAM_06_LIVE (bjd)", "rtsp://k-monitorcenter.app:8554/live/live_bjdbradbpyrr"),
]

for name, url in cameras:
    print(f"\n" + "="*60)
    print(f"Evaluating {name}...")
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
    ret, frame = cap.read()
    cap.release()

    if not ret or frame is None:
        print(f"Could not connect to {name}")
        continue

    # V1 Golden
    r1 = model_v1(frame, conf=0.20, verbose=False)
    dets_v1 = []
    if len(r1) > 0 and len(r1[0].boxes) > 0:
        for b in r1[0].boxes:
            conf = float(b.conf[0].item())
            dets_v1.append(f"{conf*100:.1f}%")

    # V2 Tesla Candidate
    r2 = model_v2(frame, conf=0.20, verbose=False)
    dets_v2 = []
    if len(r2) > 0 and len(r2[0].boxes) > 0:
        for b in r2[0].boxes:
            conf = float(b.conf[0].item())
            dets_v2.append(f"{conf*100:.1f}%")

    print(f"  [GOLDEN 89% BASELINE]   Detecções ({len(dets_v1)}): {', '.join(dets_v1) if dets_v1 else 'Nenhuma'}")
    print(f"  [TESLA DATA ENGINE V2]  Detecções ({len(dets_v2)}): {', '.join(dets_v2) if dets_v2 else 'Nenhuma'}")
