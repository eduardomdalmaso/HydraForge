import cv2
import os
import time
from ultralytics import YOLO

weights_path = "storage/models/person_detection_reinforced.pt"
print(f"Loading YOLO model: {weights_path}")
model = YOLO(weights_path)

cameras = [
    ("CAM_01_LIVE", "rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1"),
    ("CAM_02_LIVE", "rtsp://k-monitorcenter.app:8554/live/live_v0m7gz9przaz"),
    ("CAM_04_LIVE", "rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt"),
    ("CAM_05_HD",   "rtsp://k-monitorcenter.app:8554/live/live_bjdbradbpyrr"),
    ("CAM_06_HD",   "rtsp://k-monitorcenter.app:8554/live/live_k4xvpbtdqji8"),
]

os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

for name, url in cameras:
    print(f"\n==========================================")
    print(f"Connecting to {name}: {url}")
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
    ret, frame = cap.read()
    cap.release()

    if not ret or frame is None:
        print(f"FAILED to read frame from {name}")
        continue

    print(f"SUCCESS: Captured frame {frame.shape[1]}x{frame.shape[0]}")
    t0 = time.perf_counter()
    results = model.track(source=frame, persist=True, tracker="bytetrack.yaml", verbose=False)
    inf_ms = (time.perf_counter() - t0) * 1000.0

    dets = []
    if len(results) > 0:
        res = results[0]
        for box in res.boxes:
            cls_id = int(box.cls[0].item())
            cls_name = res.names.get(cls_id, str(cls_id))
            conf = float(box.conf[0].item())
            dets.append(f"{cls_name} ({conf*100:.1f}%)")

    print(f"Inference Time: {inf_ms:.1f}ms | Detections ({len(dets)}): {', '.join(dets) if dets else 'None'}")
