import os
from ultralytics import YOLO

base_model = 'storage/models/golden_backup/person_detection_reinforced_89pct_golden.pt'
data_yaml = 'datasets/tesla_hard_cases_v2/data.yaml'

print("=== INITIATING TESLA DATA ENGINE V2 TRAINING ===")
print(f"Base: {base_model}")
print(f"Data: {data_yaml}")
print(f"Freeze: 15 layers (Backbone preserved)")
print(f"Learning Rate: 0.00004")

model = YOLO(base_model)

results = model.train(
    data=data_yaml,
    epochs=12,
    imgsz=640,
    batch=8,
    device='cpu',
    freeze=15,             # Freeze backbone
    lr0=0.00004,           # Micro learning rate
    lrf=0.01,
    optimizer='AdamW',
    weight_decay=0.0005,
    warmup_epochs=1,
    mosaic=0.3,
    mixup=0.0,
    fliplr=0.5,
    project='runs/tesla_data_engine_v2',
    name='iter2_voc_hard',
    exist_ok=True,
    verbose=True
)

print("\n=== TRAINING V2 COMPLETED ===")
