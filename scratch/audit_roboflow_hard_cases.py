import os
import glob
import random
import cv2
from ultralytics import YOLO

# Load production model
model = YOLO("storage/models/person_detection_reinforced.pt")

print("=== STARTING TESLA HARD-CASE AUDITOR (CRITERIA: CONF < 45%) ===")

datasets = {
    "leo_ueno": "datasets/roboflow/leo_ueno/train/images/*.jpg",
    "human_v2": "datasets/roboflow/human_v2/train/images/*.jpg",
    "people_kahye": "datasets/roboflow/people_kahye/train/images/*.jpg"
}

mined_hard_cases = []
total_scanned = 0

for ds_name, pattern in datasets.items():
    all_imgs = glob.glob(pattern)
    if not all_imgs:
        continue
    random.seed(42)
    sample = random.sample(all_imgs, min(100, len(all_imgs)))
    print(f"\n[AUDITING] {ds_name} (Sampling {len(sample)} of {len(all_imgs)} images)...")
    
    ds_hard_cases = 0
    for img_path in sample:
        total_scanned += 1
        res = model.predict(source=img_path, conf=0.15, classes=[0], verbose=False)[0]
        boxes = res.boxes
        confs = [float(b.conf[0]) for b in boxes]
        
        # Hard case criteria: model detected someone with confidence < 45% (hesitation/occlusion)
        low_confs = [c for c in confs if 0.20 <= c <= 0.45]
        
        if len(low_confs) > 0:
            ds_hard_cases += 1
            mined_hard_cases.append({
                'dataset': ds_name,
                'path': img_path,
                'confs': [round(c, 3) for c in confs],
                'hard_confs': [round(c, 3) for c in low_confs],
                'boxes_xywhn': [b.xywhn[0].tolist() for b in boxes]
            })

    print(f"  --> Found {ds_hard_cases} HARD cases (conf < 45%) in {ds_name}")

print("\n" + "="*60)
print(f"AUDIT SUMMARY:")
print(f"Total Images Scanned: {total_scanned}")
print(f"Total Hard Cases (< 45%): {len(mined_hard_cases)} ({(len(mined_hard_cases)/total_scanned)*100:.1f}%)")
print("="*60)

print("\nTop 5 Hardest Mined Cases Sample:")
for hc in mined_hard_cases[:5]:
    fname = os.path.basename(hc['path'])
    print(f" - [{hc['dataset']}] {fname}: All Confs={hc['confs']} | Under 45%: {hc['hard_confs']}")
