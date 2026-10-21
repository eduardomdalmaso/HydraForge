import os
import glob
import shutil
import random
import yaml
from ultralytics import YOLO

# Production model
model = YOLO("storage/models/person_detection_reinforced.pt")

output_dir = "datasets/tesla_roboflow_hard"
os.makedirs(f"{output_dir}/images/train", exist_ok=True)
os.makedirs(f"{output_dir}/labels/train", exist_ok=True)
os.makedirs(f"{output_dir}/images/val", exist_ok=True)
os.makedirs(f"{output_dir}/labels/val", exist_ok=True)

# 1. LIVE CAM ANCHORS (1 FRAME PER CAM)
cam_anchors = [
    ('scratch/tesla_mining/cams/live_d38ygmd5gzm1.jpg', True),
    ('scratch/tesla_mining/cams/live_g8q6t3d6jgjt.jpg', True),
    ('scratch/unseen_test/NEW_CAM_01.jpg', True), # live_98w3
    ('scratch/unseen_test/NEW_CAM_02.jpg', True), # cam-29x5 (352x240)
    ('scratch/unseen_test/NEW_CAM_03_EMPTY.jpg', False), # akq74s5r (negative empty)
]

print("Adding live camera anchors...")
for img_path, has_person in cam_anchors:
    if not os.path.exists(img_path):
        continue
    base = os.path.splitext(os.path.basename(img_path))[0]
    dest_img = f"{output_dir}/images/train/cam_{base}.jpg"
    dest_lbl = f"{output_dir}/labels/train/cam_{base}.txt"
    shutil.copyfile(img_path, dest_img)
    
    with open(dest_lbl, "w") as f:
        if has_person:
            res = model.predict(source=img_path, conf=0.45, classes=[0], verbose=False)[0]
            for b in res.boxes:
                x, y, w, h = b.xywhn[0].tolist()
                f.write(f"0 {x:.6f} {y:.6f} {w:.6f} {h:.6f}\n")

# 2. MINE HARD CASES (< 45%) FROM ROBOFLOW WITH GROUND TRUTH MAPPING
leo_human_ids = {6, 7, 8, 9, 10, 25, 41, 42, 43}

configs = [
    {
        "name": "leo_ueno",
        "imgs": "datasets/roboflow/leo_ueno/train/images/*.jpg",
        "lbls_dir": "datasets/roboflow/leo_ueno/train/labels",
        "human_cls_filter": lambda c: c in leo_human_ids
    },
    {
        "name": "human_v2",
        "imgs": "datasets/roboflow/human_v2/train/images/*.jpg",
        "lbls_dir": "datasets/roboflow/human_v2/train/labels",
        "human_cls_filter": lambda c: c in {0, 1} # 'Human ', 'Human'
    },
    {
        "name": "people_kahye",
        "imgs": "datasets/roboflow/people_kahye/train/images/*.jpg",
        "lbls_dir": "datasets/roboflow/people_kahye/train/labels",
        "human_cls_filter": lambda c: c == 0 # 'person'
    },
]

collected_hard_cases = []

for cfg in configs:
    all_imgs = glob.glob(cfg["imgs"])
    random.seed(42)
    sample = random.sample(all_imgs, min(120, len(all_imgs)))
    
    print(f"Scanning {cfg['name']} for hard cases (< 45%)...")
    for img_path in sample:
        base = os.path.splitext(os.path.basename(img_path))[0]
        lbl_path = os.path.join(cfg["lbls_dir"], f"{base}.txt")
        
        # Parse ground-truth boxes
        gt_boxes = []
        if os.path.exists(lbl_path):
            with open(lbl_path, "r") as f:
                for line in f:
                    parts = line.strip().split()
                    if len(parts) >= 5 and cfg["human_cls_filter"](int(parts[0])):
                        gt_boxes.append([float(p) for p in parts[1:5]])
                        
        if not gt_boxes:
            continue
            
        # Predict with production model
        res = model.predict(source=img_path, conf=0.15, classes=[0], verbose=False)[0]
        confs = [float(b.conf[0]) for b in res.boxes]
        
        # Check hard case condition: low confidence detection (< 45%) OR missed person
        has_low_conf = any(0.20 <= c <= 0.45 for c in confs)
        missed = len(confs) < len(gt_boxes)
        
        if has_low_conf or missed:
            collected_hard_cases.append((img_path, gt_boxes))

print(f"\nTotal Hard Cases Mined with Ground Truth: {len(collected_hard_cases)}")

# Shuffle and split into train and val (80% train, 20% val)
random.seed(42)
random.shuffle(collected_hard_cases)

split_idx = int(len(collected_hard_cases) * 0.8)
train_cases = collected_hard_cases[:split_idx]
val_cases = collected_hard_cases[split_idx:]

for split_name, cases in [("train", train_cases), ("val", val_cases)]:
    for img_path, gt_boxes in cases:
        base = os.path.splitext(os.path.basename(img_path))[0]
        dest_img = f"{output_dir}/images/{split_name}/{base}.jpg"
        dest_lbl = f"{output_dir}/labels/{split_name}/{base}.txt"
        shutil.copyfile(img_path, dest_img)
        
        with open(dest_lbl, "w") as f:
            for b in gt_boxes:
                # Class 0 = person
                f.write(f"0 {b[0]:.6f} {b[1]:.6f} {b[2]:.6f} {b[3]:.6f}\n")

# Write data.yaml
yaml_content = f"""path: {os.path.abspath(output_dir).replace(chr(92), '/')}
train: images/train
val: images/val

names:
  0: person
"""
with open(f"{output_dir}/data.yaml", "w") as f:
    f.write(yaml_content)

tr_cnt = len(glob.glob(f"{output_dir}/images/train/*.jpg"))
val_cnt = len(glob.glob(f"{output_dir}/images/val/*.jpg"))
print(f"\n=== TESLA ROBOFLOW HARD DATASET READY ===")
print(f"Train images: {tr_cnt}")
print(f"Val images: {val_cnt}")
print(f"data.yaml path: {output_dir}/data.yaml")
