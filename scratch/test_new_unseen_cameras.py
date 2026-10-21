import cv2
import os
import time
from ultralytics import YOLO

os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

v1_path = "storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt"
v2_path = "runs/detect/runs/tesla_data_engine/v1_hard_cases/weights/best.pt"

print(f"Loading Model V1 (Golden 89%): {v1_path}")
model_v1 = YOLO(v1_path)

print(f"Loading Model V2 (Tesla Candidate): {v2_path}")
model_v2 = YOLO(v2_path)

new_cameras = [
    ("NEW_CAM_01 (live_98w3)", "rtsp://k-monitorcenter.app:8554/live/live_98w3g4fvjl1z", False),
    ("NEW_CAM_02 (cam-29x5)", "rtsp://k-monitorcenter.app:8554/live/cam-29x5ldb506me", False),
    ("NEW_CAM_03_EMPTY (akq74s5r)", "rtsp://k-monitorcenter.app:8554/live/akq74s5r", True), # True = No person expected
    ("NEW_CAM_04 (208_main)", "rtsp://k-monitorcenter.app:8554/208_main", False),
    ("NEW_CAM_05 (216_main)", "rtsp://k-monitorcenter.app:8554/216_main", False),
    ("NEW_CAM_06 (live_81pz)", "rtsp://k-monitorcenter.app:8554/live/live_81pzxh8idvi0", False),
]

os.makedirs("scratch/unseen_test", exist_ok=True)

print("\n" + "="*70)
print("=== TESTING ON 100% UNSEEN REAL CAMERAS (OUT-OF-DISTRIBUTION) ===")
print("="*70)

for name, url, is_empty_expected in new_cameras:
    print(f"\n[STREAM] {name} | Expected Empty: {is_empty_expected}")
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
    ret, frame = cap.read()
    cap.release()

    if not ret or frame is None:
        print(f"  --> OFFLINE ou Timeout no RTSP")
        continue

    # Save frame for debug
    safe_name = name.split()[0]
    cv2.imwrite(f"scratch/unseen_test/{safe_name}.jpg", frame)

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
    
    if is_empty_expected:
        v1_fp = len(dets_v1) > 0
        v2_fp = len(dets_v2) > 0
        print(f"  --> Validação Falso Positivo: V1={'FALSO POSITIVO!' if v1_fp else 'PERFEITO (0)'} | V2={'FALSO POSITIVO!' if v2_fp else 'PERFEITO (0)'}")
