import cv2
import os
import glob
from ultralytics import YOLO

# Golden 89% model
model = YOLO('storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt')

cams = {
    "live_d38ygmd5gzm1": "rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1",
    "live_g8q6t3d6jgjt": "rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt",
    "live_v0m7gz9przaz": "rtsp://k-monitorcenter.app:8554/live/live_v0m7gz9przaz",
    "live_k4xvpbtdqji8": "rtsp://k-monitorcenter.app:8554/live/live_k4xvpbtdqji8",
    "live_bjdbradbpyrr": "rtsp://k-monitorcenter.app:8554/live/live_bjdbradbpyrr",
}

os.makedirs('scratch/tesla_mining/cams', exist_ok=True)

print("=== CAPTURING EXACTLY 1 FRAME PER LIVE CAMERA ===")
for cam_id, url in cams.items():
    cap = cv2.VideoCapture(url)
    cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)
    ret, frame = cap.read()
    cap.release()
    if ret and frame is not None:
        save_path = f'scratch/tesla_mining/cams/{cam_id}.jpg'
        cv2.imwrite(save_path, frame)
        results = model.predict(source=frame, conf=0.25, classes=[0], verbose=False)
        boxes = results[0].boxes
        confs = [round(float(b.conf[0]), 3) for b in boxes] if len(boxes) > 0 else []
        print(f"[{cam_id}] Captured: {frame.shape[1]}x{frame.shape[0]} | Detected {len(boxes)} persons | Confs: {confs}")
    else:
        print(f"[{cam_id}] Failed to capture frame")
