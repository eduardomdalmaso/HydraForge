import os
import shutil
import random
from ultralytics import YOLO

# Golden 89% model
model = YOLO('storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt')

dataset_dir = 'datasets/tesla_hard_cases'
os.makedirs(f'{dataset_dir}/images/train', exist_ok=True)
os.makedirs(f'{dataset_dir}/labels/train', exist_ok=True)
os.makedirs(f'{dataset_dir}/images/val', exist_ok=True)
os.makedirs(f'{dataset_dir}/labels/val', exist_ok=True)

# 1. LIVE CAMS (EXACTLY 1 FRAME PER CAMERA)
cam_files = [
    'scratch/tesla_mining/cams/live_d38ygmd5gzm1.jpg',
    'scratch/tesla_mining/cams/live_g8q6t3d6jgjt.jpg',
    'scratch/tesla_mining/cams/live_v0m7gz9przaz.jpg',
    'scratch/tesla_mining/cams/live_k4xvpbtdqji8.jpg',
    'scratch/tesla_mining/cams/live_bjdbradbpyrr.jpg'
]

print("=== PROCESSING 1 FRAME PER CAMERA ===")
for img_path in cam_files:
    if not os.path.exists(img_path):
        continue
    base = os.path.splitext(os.path.basename(img_path))[0]
    # Predict with golden model to get clean bounding boxes
    res = model.predict(source=img_path, conf=0.50, classes=[0], verbose=False)[0]
    
    # Save image to train
    dest_img = f'{dataset_dir}/images/train/cam_{base}.jpg'
    shutil.copyfile(img_path, dest_img)
    
    # Save YOLO format label (class 0 = person)
    dest_lbl = f'{dataset_dir}/labels/train/cam_{base}.txt'
    with open(dest_lbl, 'w') as f:
        for b in res.boxes:
            # xywh normalized
            x, y, w, h = b.xywhn[0].tolist()
            f.write(f"0 {x:.6f} {y:.6f} {w:.6f} {h:.6f}\n")
    print(f"Cam frame added: cam_{base}.jpg with {len(res.boxes)} person boxes")

# 2. INTERNET DATASET (COCO128 HARD CASES + ANCHORS)
# We mine from coco128:
import glob
coco_images = glob.glob('datasets/coco128/images/train2017/*.jpg')

hard_positives = []
hard_negatives = []
easy_anchors = []

for img_path in coco_images:
    base = os.path.splitext(os.path.basename(img_path))[0]
    lbl_path = f'datasets/coco128/labels/train2017/{base}.txt'
    gt_persons = []
    if os.path.exists(lbl_path):
        with open(lbl_path, 'r') as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) >= 5 and parts[0] == '0':
                    gt_persons.append(line.strip())
                    
    res = model.predict(source=img_path, conf=0.15, classes=[0], verbose=False)[0]
    pred_confs = [float(b.conf[0]) for b in res.boxes]
    
    if len(gt_persons) > 0:
        has_low_conf = any(c < 0.65 for c in pred_confs)
        missed_some = len(pred_confs) < len(gt_persons)
        if has_low_conf or missed_some:
            hard_positives.append((img_path, gt_persons))
        else:
            easy_anchors.append((img_path, gt_persons))
    else:
        # No persons in GT
        if len(pred_confs) > 0:
            hard_negatives.append(img_path)

print(f"\nFound {len(hard_positives)} Hard Positives, {len(hard_negatives)} Hard Negatives, {len(easy_anchors)} Easy Anchors.")

# Select: all 29 hard positives + 3 hard negatives + 10 easy anchors
selected_train = []
selected_val = []

# Split hard positives: 24 train, 5 val
random.seed(42)
random.shuffle(hard_positives)
for img, gts in hard_positives[:24]:
    selected_train.append((img, gts))
for img, gts in hard_positives[24:]:
    selected_val.append((img, gts))

# Hard negatives (empty background): 2 train, 1 val
for img in hard_negatives[:2]:
    selected_train.append((img, []))
for img in hard_negatives[2:]:
    selected_val.append((img, []))

# Easy anchors: 8 train, 4 val
random.shuffle(easy_anchors)
for img, gts in easy_anchors[:8]:
    selected_train.append((img, gts))
for img, gts in easy_anchors[8:12]:
    selected_val.append((img, gts))

# Copy to datasets/tesla_hard_cases
for split, items in [('train', selected_train), ('val', selected_val)]:
    for img_path, gts in items:
        base = os.path.splitext(os.path.basename(img_path))[0]
        dest_img = f'{dataset_dir}/images/{split}/net_{base}.jpg'
        shutil.copyfile(img_path, dest_img)
        dest_lbl = f'{dataset_dir}/labels/{split}/net_{base}.txt'
        with open(dest_lbl, 'w') as f:
            for line in gts:
                # keep only class 0
                parts = line.split()
                if parts[0] == '0':
                    f.write(f"0 {' '.join(parts[1:5])}\n")

# Write data.yaml
yaml_content = f"""path: {os.path.abspath(dataset_dir).replace(chr(92), '/')}
train: images/train
val: images/val

names:
  0: person
"""
with open(f'{dataset_dir}/data.yaml', 'w') as f:
    f.write(yaml_content)

print("\n=== TESLA DATASET READY ===")
train_cnt = len(glob.glob(f'{dataset_dir}/images/train/*.jpg'))
val_cnt = len(glob.glob(f'{dataset_dir}/images/val/*.jpg'))
print(f"Train images: {train_cnt} (5 live cams + 24 hard cases + 2 hard negs + 8 anchors)")
print(f"Val images: {val_cnt} (5 hard cases + 1 hard neg + 4 anchors)")
print(f"data.yaml path: {dataset_dir}/data.yaml")
