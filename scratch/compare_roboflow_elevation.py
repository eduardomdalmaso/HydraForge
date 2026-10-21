import cv2
import os
import time
from ultralytics import YOLO

os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

v1_path = "storage/models/person_detection_reinforced.pt"
v2_path = "runs/detect/runs/tesla_roboflow_elevation/hard_cases_elevation/weights/best.pt"

print(f"=== TESTING ROBOFLOW ELEVATION CANDIDATE vs PRODUCTION ===")
print(f"Model V1 (Current Prod): {v1_path}")
print(f"Model V2 (Roboflow Elev): {v2_path}")

model_v1 = YOLO(v1_path)
model_v2 = YOLO(v2_path)

cameras = [
    ("NEW_CAM_02 (cam-29x5 352x240)", "rtsp://k-monitorcenter.app:8554/live/cam-29x5ldb506me", False),
    ("NEW_CAM_01 (live_98w3 720p)", "rtsp://k-monitorcenter.app:8554/live/live_98w3g4fvjl1z", False),
    ("NEW_CAM_03_EMPTY (akq74s5r)", "rtsp://k-monitorcenter.app:8554/live/akq74s5r", True), # Negative
    ("CAM_01_LIVE (d38ygmd5gzm1)", "rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1", False),
    ("CAM_04_LIVE (g8q6t3d6jgjt)", "rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt", False),
]

for name, url, is_empty in cameras:
    print(f"\n" + "="*65)
    print(f"Testing on {name}...")
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
    ret, frame = cap.read()
    cap.release()

    if not ret or frame is None:
        print("  --> Offline/Timeout")
        continue

    def get_dets(m):
        res = m(frame, conf=0.20, verbose=False)
        if len(res) > 0 and len(res[0].boxes) > 0:
            return [f"{float(b.conf[0])*100:.1f}%" for b in res[0].boxes]
        return []

    d1 = get_dets(model_v1)
    d2 = get_dets(model_v2)

    print(f"  [CURRENT PROD]      Detecções ({len(d1)}): {', '.join(d1) if d1 else 'Nenhuma'}")
    print(f"  [ROBOFLOW ELEV]     Detecções ({len(d2)}): {', '.join(d2) if d2 else 'Nenhuma'}")
    if is_empty:
        print(f"  --> Falso Positivo: Prod={'OK' if not d1 else 'FP'} | Elev={'OK' if not d2 else 'FP'}")
