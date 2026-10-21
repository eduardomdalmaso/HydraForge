import os
import glob
import shutil
import random
from ultralytics import YOLO

# Golden 89% model
model = YOLO('storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt')

dataset_dir = 'datasets/tesla_hard_cases_v2'
os.makedirs(f'{dataset_dir}/images/train', exist_ok=True)
os.makedirs(f'{dataset_dir}/labels/train', exist_ok=True)
os.makedirs(f'{dataset_dir}/images/val', exist_ok=True)
os.makedirs(f'{dataset_dir}/labels/val', exist_ok=True)

# 1. ADD ALL UNIQUE CAMERA FRAMES (1 FRAME PER CAM)
# Old cams
old_cams = glob.glob('scratch/tesla_mining/cams/*.jpg')
for c in old_cams:
    base = os.path.splitext(os.path.basename(c))[0]
    dest_img = f'{dataset_dir}/images/train/{base}.jpg'
    shutil.copyfile(c, dest_img)
    # Get box
    r = model.predict(source=c, conf=0.45, classes=[0], verbose=False)[0]
    with open(f'{dataset_dir}/labels/train/{base}.txt', 'w') as f:
        for b in r.boxes:
            x, y, w, h = b.xywhn[0].tolist()
            f.write(f"0 {x:.6f} {y:.6f} {w:.6f} {h:.6f}\n")

# New cams
new_cams_info = [
    # (image_file, [boxes_xywhn], is_val)
    ('scratch/unseen_test/NEW_CAM_01.jpg', [[0.498251, 0.775388, 0.343872, 0.382249]], False), # live_98w3
    ('scratch/unseen_test/NEW_CAM_02.jpg', [[0.059338, 0.777769, 0.115823, 0.182090]], False), # cam-29x5 (352x240)
    ('scratch/unseen_test/NEW_CAM_05.jpg', [[0.293777, 0.104014, 0.017419, 0.039383]], False), # 216_main (far person)
    ('scratch/unseen_test/NEW_CAM_03_EMPTY.jpg', [], False), # akq74s5r (empty room negative)
    ('scratch/unseen_test/NEW_CAM_04.jpg', [], True), # 208_main (val empty)
]

for img_path, boxes, is_val in new_cams_info:
    split = 'val' if is_val else 'train'
    base = os.path.splitext(os.path.basename(img_path))[0]
    dest_img = f'{dataset_dir}/images/{split}/new_{base}.jpg'
    shutil.copyfile(img_path, dest_img)
    with open(f'{dataset_dir}/labels/{split}/new_{base}.txt', 'w') as f:
        for b in boxes:
            f.write(f"0 {b[0]:.6f} {b[1]:.6f} {b[2]:.6f} {b[3]:.6f}\n")

print("All unique camera frames registered.")

# 2. MINE HARD CASES FROM PASCAL VOC (INTERNET DATASET)
# In VOC, person is class 14
voc_label_files = glob.glob('datasets/VOC/labels/train2007/*.txt') + glob.glob('datasets/VOC/labels/val2007/*.txt')
print(f"Scanning {len(voc_label_files)} VOC label files for hard person cases...")

voc_hard_cases = []
random.seed(42)
random.shuffle(voc_label_files)

for lbl_path in voc_label_files:
    if len(voc_hard_cases) >= 30:
        break
    
    # Check if contains class 14 (person)
    person_boxes = []
    with open(lbl_path, 'r') as f:
        for line in f:
            parts = line.strip().split()
            if len(parts) >= 5 and parts[0] == '14':
                person_boxes.append([float(p) for p in parts[1:5]])
                
    if not person_boxes:
        continue
        
    base = os.path.splitext(os.path.basename(lbl_path))[0]
    # Find matching image in train2007 or val2007
    img_path = None
    for folder in ['train2007', 'val2007']:
        candidate = f'datasets/VOC/images/{folder}/{base}.jpg'
        if os.path.exists(candidate):
            img_path = candidate
            break
            
    if not img_path:
        continue
        
    # Run model with low threshold to find uncertainty
    res = model.predict(source=img_path, conf=0.15, classes=[0], verbose=False)[0]
    pred_confs = [float(b.conf[0]) for b in res.boxes]
    
    # Hard criteria: low conf (0.20 to 0.55) or missed person
    has_low_conf = any(0.20 <= c <= 0.55 for c in pred_confs)
    missed_some = len(pred_confs) < len(person_boxes)
    
    if has_low_conf or missed_some:
        voc_hard_cases.append((img_path, person_boxes, pred_confs))

print(f"Mined {len(voc_hard_cases)} HARD person cases from Pascal VOC!")

# Save VOC hard cases (remap class 14 -> 0)
for idx, (img_path, boxes, confs) in enumerate(voc_hard_cases):
    split = 'val' if idx % 5 == 0 else 'train'
    base = os.path.splitext(os.path.basename(img_path))[0]
    dest_img = f'{dataset_dir}/images/{split}/voc_{base}.jpg'
    shutil.copyfile(img_path, dest_img)
    with open(f'{dataset_dir}/labels/{split}/voc_{base}.txt', 'w') as f:
        for b in boxes:
            f.write(f"0 {b[0]:.6f} {b[1]:.6f} {b[2]:.6f} {b[3]:.6f}\n")

# Write data.yaml
yaml_content = f"""path: {os.path.abspath(dataset_dir).replace(chr(92), '/')}
train: images/train
val: images/val

names:
  0: person
"""
with open(f'{dataset_dir}/data.yaml', 'w') as f:
    f.write(yaml_content)

train_imgs = glob.glob(f'{dataset_dir}/images/train/*.jpg')
val_imgs = glob.glob(f'{dataset_dir}/images/val/*.jpg')
print(f"\n=== TESLA DATASET V2 READY ===")
print(f"Train images: {len(train_imgs)}")
print(f"Val images: {len(val_imgs)}")
print(f"data.yaml path: {dataset_dir}/data.yaml")
