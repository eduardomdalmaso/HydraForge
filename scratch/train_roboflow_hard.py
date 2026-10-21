import os
from ultralytics import YOLO

base_model = "storage/models/person_detection_reinforced.pt"
data_yaml = "datasets/tesla_roboflow_hard/data.yaml"

print("=== INITIATING ROBOFLOW HARD CASES ELEVATION TRAINING ===")
print(f"Base: {base_model}")
print(f"Data: {data_yaml}")
print(f"Freeze: 15 layers (Backbone preserved)")
print(f"Learning Rate (lr0): 0.00004")

model = YOLO(base_model)

results = model.train(
    data=data_yaml,
    epochs=15,
    imgsz=640,
    batch=8,
    device='cpu',
    freeze=15,             # Freeze backbone layers
    lr0=0.00004,           # Surgical micro-LR
    lrf=0.01,
    optimizer='AdamW',
    weight_decay=0.0005,
    warmup_epochs=1,
    mosaic=0.2,
    mixup=0.0,
    fliplr=0.5,
    project='runs/tesla_roboflow_elevation',
    name='hard_cases_elevation',
    exist_ok=True,
    verbose=True
)

print("\n=== ELEVATION TRAINING COMPLETED ===")
