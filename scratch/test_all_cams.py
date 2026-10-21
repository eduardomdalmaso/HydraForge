import urllib.request
import json

cameras = [
    'cam_live_d38',
    'cam_live_v0m',
    'cam_live_g8q',
    'cam_hd_bjd',
    'cam_hd_k4x'
]

url = 'http://127.0.0.1:8081/api/v1/inference/predict'

for cam in cameras:
    payload = {
        'model': 'person_detection_reinforced',
        'source': cam,
        'conf': 0.25,
        'iou': 0.45,
        'track': True
    }
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(url, data=data, headers={'Content-Type': 'application/json'})
    try:
        with urllib.request.urlopen(req, timeout=20.0) as resp:
            res = json.loads(resp.read().decode('utf-8'))
            print(f"[{cam}] Status: {res.get('status')} | Detections: {len(res.get('detections', []))}")
            for d in res.get('detections', []):
                print(f"   -> {d.get('label')} conf={d.get('conf')} box={d.get('box')}")
    except Exception as e:
        print(f"[{cam}] Error: {e}")
