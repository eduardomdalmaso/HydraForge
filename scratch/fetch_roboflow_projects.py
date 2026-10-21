import os
import json
import zipfile
import urllib.request

api_key = "nGTrUsOuJRHkAPQhc3ng"
target_dir = "datasets/roboflow"
os.makedirs(target_dir, exist_ok=True)

projects = [
    ("leo-ueno", "people-detection-o4rdr", "leo_ueno"),
    ("human-v2", "human-dataset-v2", "human_v2"),
    ("lennys-workspace-gabu9", "people-kahye", "people_kahye"),
]

for ws, proj, local_name in projects:
    print(f"\n==================================================")
    print(f"Connecting to Roboflow project: {ws}/{proj}...")
    
    # 1. Get project info
    info_url = f"https://api.roboflow.com/{ws}/{proj}?api_key={api_key}"
    req = urllib.request.Request(info_url, headers={'User-Agent': 'Mozilla/5.0'})
    try:
        res = urllib.request.urlopen(req)
        data = json.loads(res.read().decode('utf-8'))
        p = data.get('project', {})
        print(f"  Project: {p.get('name')} | Type: {p.get('type')}")
        
        # Check versions
        versions_dict = p.get('versions')
        version_num = 1
        if isinstance(versions_dict, list) and len(versions_dict) > 0:
            version_num = versions_dict[0].get('version', 1)
        elif isinstance(versions_dict, int):
            version_num = versions_dict
        elif isinstance(versions_dict, dict):
            version_num = list(versions_dict.keys())[-1]
            
        print(f"  Using version: {version_num}")
        
        # 2. Request YOLOv8 export link
        export_url = f"https://api.roboflow.com/{ws}/{proj}/{version_num}/yolov8?api_key={api_key}"
        exp_req = urllib.request.Request(export_url, headers={'User-Agent': 'Mozilla/5.0'})
        exp_res = urllib.request.urlopen(exp_req)
        exp_data = json.loads(exp_res.read().decode('utf-8'))
        
        zip_link = exp_data.get('export', {}).get('link')
        if not zip_link:
            print(f"  Could not get export link: {exp_data}")
            continue
            
        print(f"  Export link ready. Downloading zip...")
        zip_path = os.path.join(target_dir, f"{local_name}.zip")
        extract_dir = os.path.join(target_dir, local_name)
        os.makedirs(extract_dir, exist_ok=True)
        
        # Download zip
        urllib.request.urlretrieve(zip_link, zip_path)
        print(f"  Downloaded: {os.path.getsize(zip_path) / (1024*1024):.2f} MB")
        
        # Unzip
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_dir)
        print(f"  Extracted to {extract_dir}")
        
    except Exception as e:
        print(f"  Failed for {ws}/{proj}: {e}")

print("\n=== ALL AVAILABLE ROBOFLOW DATASETS READY ===")
