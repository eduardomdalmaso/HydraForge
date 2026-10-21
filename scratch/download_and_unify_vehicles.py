import os
import zipfile
import urllib.request
import glob
import shutil
import random
import yaml

raw_dir = "datasets/vehicles_raw"
unified_dir = "datasets/vehicles_unified"

os.makedirs(raw_dir, exist_ok=True)
os.makedirs(f"{unified_dir}/images/train", exist_ok=True)
os.makedirs(f"{unified_dir}/labels/train", exist_ok=True)
os.makedirs(f"{unified_dir}/images/val", exist_ok=True)
os.makedirs(f"{unified_dir}/labels/val", exist_ok=True)

downloads = [
    ("motorbike", "https://app.roboflow.com/ds/neuoQnIZ6q?key=IVn08oUQdd"),
    ("ua_detrac", "https://app.roboflow.com/ds/q20waK20G2?key=GgO1qRgxyN")
]

print("=== STEP 1: DOWNLOADING VEHICLE DATASETS ===")
for name, url in downloads:
    zip_path = os.path.join(raw_dir, f"{name}.zip")
    extract_path = os.path.join(raw_dir, name)
    
    if not os.path.exists(extract_path):
        if not os.path.exists(zip_path):
            print(f"Downloading {name} from Roboflow...")
            urllib.request.urlretrieve(url, zip_path)
            print(f"Downloaded {name}.zip: {os.path.getsize(zip_path)/(1024*1024):.2f} MB")
            
        print(f"Extracting {name}.zip...")
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_path)
        print(f"Extracted {name} successfully.")
    else:
        print(f"{name} already extracted at {extract_path}.")

# STEP 2: CLASS REMAPPING TO UNIFIED 4-CLASS SCHEMA:
# 0: car, 1: motorcycle, 2: bus, 3: truck
print("\n=== STEP 2: PARSING AND REMAPPING ANNOTATIONS ===")

# Check Motorbike dataset yaml
mb_yaml_path = os.path.join(raw_dir, "motorbike", "data.yaml")
with open(mb_yaml_path, 'r') as f:
    mb_yaml = yaml.safe_load(f)
mb_names = mb_yaml.get('names', [])
print(f"Motorbike dataset classes: {mb_names}")

# Map mb_names to our target IDs (0: car, 1: motorcycle, 2: bus, 3: truck)
mb_map = {}
for i, n in enumerate(mb_names):
    n_lower = n.lower()
    if 'car' in n_lower:
        mb_map[i] = 0
    elif 'motor' in n_lower or 'bike' in n_lower:
        mb_map[i] = 1
    elif 'bus' in n_lower:
        mb_map[i] = 2
    elif 'truck' in n_lower:
        mb_map[i] = 3

print(f"Motorbike class mapping: {mb_map}")

# Check UA-DETRAC dataset yaml
ua_yaml_path = os.path.join(raw_dir, "ua_detrac", "data.yaml")
with open(ua_yaml_path, 'r') as f:
    ua_yaml = yaml.safe_load(f)
ua_names = ua_yaml.get('names', [])
print(f"UA-DETRAC dataset classes: {ua_names}")

ua_map = {}
for i, n in enumerate(ua_names):
    n_lower = n.lower()
    if 'car' in n_lower or 'van' in n_lower:
        ua_map[i] = 0
    elif 'motor' in n_lower or 'bike' in n_lower:
        ua_map[i] = 1
    elif 'bus' in n_lower:
        ua_map[i] = 2
    elif 'truck' in n_lower:
        ua_map[i] = 3

print(f"UA-DETRAC class mapping: {ua_map}")

# STEP 3: ASSEMBLE BALANCED UNIFIED DATASET
print("\n=== STEP 3: ASSEMBLING BALANCED DATASET ===")

processed_samples = []

# Process Motorbike images (rich in motorcycles)
mb_train_images = glob.glob(os.path.join(raw_dir, "motorbike", "train", "images", "*.jpg"))
for img_path in mb_train_images:
    base = os.path.splitext(os.path.basename(img_path))[0]
    lbl_path = os.path.join(raw_dir, "motorbike", "train", "labels", f"{base}.txt")
    if not os.path.exists(lbl_path):
        continue
    new_lines = []
    with open(lbl_path, 'r') as f:
        for line in f:
            parts = line.strip().split()
            if len(parts) >= 5:
                old_cls = int(parts[0])
                if old_cls in mb_map:
                    new_lines.append(f"{mb_map[old_cls]} {' '.join(parts[1:5])}")
    if new_lines:
        processed_samples.append((img_path, new_lines))

print(f"Collected {len(processed_samples)} samples from Motorbike dataset.")

# Process UA-DETRAC images (sample up to 1,500 balanced images of cars, buses, trucks)
ua_train_images = glob.glob(os.path.join(raw_dir, "ua_detrac", "train", "images", "*.jpg"))
random.seed(42)
random.shuffle(ua_train_images)

ua_count = 0
for img_path in ua_train_images:
    if ua_count >= 1500:
        break
    base = os.path.splitext(os.path.basename(img_path))[0]
    lbl_path = os.path.join(raw_dir, "ua_detrac", "train", "labels", f"{base}.txt")
    if not os.path.exists(lbl_path):
        continue
    new_lines = []
    with open(lbl_path, 'r') as f:
        for line in f:
            parts = line.strip().split()
            if len(parts) >= 5:
                old_cls = int(parts[0])
                if old_cls in ua_map:
                    new_lines.append(f"{ua_map[old_cls]} {' '.join(parts[1:5])}")
    if new_lines:
        processed_samples.append((img_path, new_lines))
        ua_count += 1

print(f"Collected {ua_count} samples from UA-DETRAC dataset.")
print(f"Total Unified Samples: {len(processed_samples)}")

# Shuffle and split: 85% train, 15% val
random.shuffle(processed_samples)
split_idx = int(len(processed_samples) * 0.85)
train_set = processed_samples[:split_idx]
val_set = processed_samples[split_idx:]

for split_name, dataset in [("train", train_set), ("val", val_set)]:
    for idx, (img_path, lines) in enumerate(dataset):
        ext = os.path.splitext(img_path)[1]
        out_base = f"veh_{idx:05d}"
        dest_img = os.path.join(unified_dir, "images", split_name, f"{out_base}{ext}")
        dest_lbl = os.path.join(unified_dir, "labels", split_name, f"{out_base}.txt")
        shutil.copyfile(img_path, dest_img)
        with open(dest_lbl, 'w') as f:
            f.write("\n".join(lines) + "\n")

# Write data.yaml
unified_yaml = f"""path: {os.path.abspath(unified_dir).replace(chr(92), '/')}
train: images/train
val: images/val

names:
  0: car
  1: motorcycle
  2: bus
  3: truck
"""
with open(os.path.join(unified_dir, "data.yaml"), "w") as f:
    f.write(unified_yaml)

print("\n=== UNIFIED VEHICLE DATASET READY ===")
train_cnt = len(glob.glob(f"{unified_dir}/images/train/*.*"))
val_cnt = len(glob.glob(f"{unified_dir}/images/val/*.*"))
print(f"Train images: {train_cnt}")
print(f"Val images: {val_cnt}")
print(f"Config: {unified_dir}/data.yaml")
