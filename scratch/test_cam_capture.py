import urllib.request
import json
import cv2
import os

cam_clean = 'cam_live_d38'
req = urllib.request.Request('http://127.0.0.1:8080/api/v1/streams', headers={'User-Agent': 'HydraForge/1.0'})
with urllib.request.urlopen(req, timeout=2.0) as resp:
    sdata = json.loads(resp.read().decode('utf-8'))
    for st in sdata.get('streams', []):
        if st.get('stream_id') == cam_clean and st.get('source_url'):
            rtsp_url = st.get('source_url')
            print('Found RTSP URL:', rtsp_url)
            os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'
            cap = cv2.VideoCapture(rtsp_url, cv2.CAP_FFMPEG)
            cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
            cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)
            ret, frame = cap.read()
            cap.release()
            print('Capture ret:', ret, 'shape:', frame.shape if ret else None)
            if ret and frame is not None:
                cv2.imwrite('storage/captures/cam_live_d38.jpg', frame)
                print('Wrote storage/captures/cam_live_d38.jpg successfully!')
