import os
import glob
import cv2
import numpy as np
from ultralytics import YOLO

# Golden 89% model
model = YOLO('storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt')

output_dir = 'datasets/tesla_hard_cases'
os.makedirs(f'{output_dir}/images/train', exist_ok=True)
os.makedirs(f'{output_dir}/labels/train', exist_ok=True)
os.makedirs(f'{output_dir}/images/val', exist_ok=True)
os.makedirs(f'{output_dir}/labels/val', exist_ok=True)

coco_images = glob.glob('datasets/coco128/images/train2017/*.jpg')

hard_positives = [] # persons with low confidence or missed
hard_negatives = [] # no persons in gt, but low/noisy conf
easy_cases = []

print(f"Mining hard cases from {len(coco_images)} diverse internet images...")

for img_path in coco_images:
    base = os.path.splitext(os.path.basename(img_path))[0]
    lbl_path = f'datasets/coco128/labels/train2017/{base}.txt'
    
    gt_persons = []
    if os.path.exists(lbl_path):
        with open(lbl_path, 'r') as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) >= 5 and parts[0] == '0':
                    gt_persons.append([float(p) for p in parts[1:5]])
                    
    # Run golden model with low threshold (0.15) to see uncertainty
    res = model.predict(source=img_path, conf=0.15, classes=[0], verbose=False)[0]
    pred_confs = [float(b.conf[0]) for b in res.boxes]
    
    # Classification:
    if len(gt_persons) > 0:
        # Has person in GT
        # If model predicted with low confidence (< 0.65) or fewer detections than GT
        has_low_conf = any(c < 0.65 for c in pred_confs)
        missed_some = len(pred_confs) < len(gt_persons)
        
        if has_low_conf or missed_some:
            hard_positives.append({
                'img_path': img_path,
                'lbl_path': lbl_path,
                'gt_count': len(gt_persons),
                'pred_confs': [round(c, 3) for c in pred_confs],
                'reason': f'low_conf={has_low_conf}, missed={missed_some}'
            })
        else:
            easy_cases.append(img_path)
    else:
        # No person in GT
        if len(pred_confs) > 0:
            # False alarm (hard negative)
            hard_negatives.append({
                'img_path': img_path,
                'pred_confs': [round(c, 3) for c in pred_confs]
            })

print(f"\n=== MINING RESULTS ===")
print(f"Total Images Analyzed: {len(coco_images)}")
print(f"Hard Positives (Low Conf / Missed Persons): {len(hard_positives)}")
print(f"Hard Negatives (False Alarms in Empty/Non-Person): {len(hard_negatives)}")
print(f"Easy Cases (High Conf 0.85+): {len(easy_cases)}")

for hp in hard_positives[:5]:
    print(f" - {os.path.basename(hp['img_path'])}: GT={hp['gt_count']} | Pred Confs={hp['pred_confs']} | {hp['reason']}")
