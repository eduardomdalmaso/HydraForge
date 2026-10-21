import os
import shutil
from ultralytics import YOLO

base_model = 'storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt'
data_yaml = 'datasets/tesla_hard_cases/data.yaml'

print(f"=== INITIATING TESLA DATA ENGINE FINE-TUNING ===")
print(f"Base Golden Model: {base_model}")
print(f"Dataset: {data_yaml}")
print(f"Backbone Freeze: 15 layers")
print(f"Learning Rate (lr0): 0.00005 (Micro-adjustment)")

model = YOLO(base_model)

results = model.train(
    data=data_yaml,
    epochs=12,
    imgsz=640,
    batch=8,
    device='cpu',
    freeze=15,             # Freeze first 15 layers (backbone intact)
    lr0=0.00005,           # Ultra-low LR
    lrf=0.01,
    optimizer='AdamW',
    weight_decay=0.0005,
    warmup_epochs=1,
    mosaic=0.3,
    mixup=0.0,
    fliplr=0.5,
    project='runs/tesla_data_engine',
    name='v1_hard_cases',
    exist_ok=True,
    verbose=True
)

best_pt = 'runs/tesla_data_engine/v1_hard_cases/weights/best.pt'
print(f"\n=== TRAINING FINISHED ===")
print(f"Trained candidate saved at: {best_pt}")
