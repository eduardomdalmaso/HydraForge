import cv2
import os
import time

urls = [
    'rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1',
    'rtsp://k-monitorcenter.app:8554/live/live_v0m7gz9przaz',
    'rtsp://k-monitorcenter.app:8554/live/live_9azek2m9mmg7',
    'rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt',
    'rtsp://k-monitorcenter.app:8554/live/live_bjdbradbpyrr',
    'rtsp://k-monitorcenter.app:8554/live/live_k4xvpbtdqji8'
]

os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

for i, url in enumerate(urls):
    t0 = time.time()
    stream_name = url.split('/')[-1]
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
    ret, frame = cap.read()
    elapsed = time.time() - t0
    cap.release()
    print(f"Cam {i+1} [{stream_name}]: ret={ret}, shape={frame.shape if ret and frame is not None else None}, time={elapsed:.2f}s")
