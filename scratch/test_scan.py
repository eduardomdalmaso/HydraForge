import urllib.request
import json

url = 'http://127.0.0.1:8081/api/v1/inference/predict'
payload = {
    'model': 'person_detection_reinforced',
    'source': 'cam_live_d38',
    'conf': 0.25,
    'iou': 0.45,
    'track': True
}
data = json.dumps(payload).encode('utf-8')
req = urllib.request.Request(url, data=data, headers={'Content-Type': 'application/json'})
with urllib.request.urlopen(req, timeout=15.0) as resp:
    res = json.loads(resp.read().decode('utf-8'))
    print('Status:', res.get('status'))
    print('DisplayImageUrl:', res.get('displayImageUrl'))
    print('Detections count:', len(res.get('detections', [])))
    for d in res.get('detections', []):
        print(f"  - {d.get('label')} conf={d.get('conf')} box={d.get('box')}")
