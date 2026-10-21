import os
import time
import shutil
from ultralytics import YOLO

base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
weights_init = os.path.join(base_dir, "storage/models/person_detection_reinforced.pt")
data_yaml = os.path.join(base_dir, "datasets/cctv_person_mix/data.yaml")

print(f"Loading base weights: {weights_init}")
model = YOLO(weights_init)

import torch
device = 0 if torch.cuda.is_available() else "cpu"
amp_enabled = True if torch.cuda.is_available() else False
print(f"Target device: {device} (AMP: {amp_enabled})")

print("Starting Fine-Tuning...")
t0 = time.time()
results = model.train(
    data=data_yaml,
    epochs=20,
    imgsz=640,
    batch=8,
    device=device,
    optimizer="AdamW",
    lr0=0.001,
    lrf=0.01,
    amp=amp_enabled,
    workers=2,
    project=os.path.join(base_dir, "runs/train"),
    name="person_cctv_finetuned",
    exist_ok=True,
    save=True,
    val=True,
    plots=True,
    verbose=True,
    # Augmentation
    hsv_h=0.015,
    hsv_s=0.7,
    hsv_v=0.4,
    degrees=5.0,
    translate=0.1,
    scale=0.5,
    fliplr=0.5,
    mosaic=1.0,
    mixup=0.1
)
elapsed = time.time() - t0
print(f"\nFine-Tuning completed in {elapsed:.1f}s ({elapsed/60:.1f} minutes)!")

best_weights = os.path.join(base_dir, "runs/train/person_cctv_finetuned/weights/best.pt")
dest_model = os.path.join(base_dir, "storage/models/person_detection_reinforced.pt")

if os.path.exists(best_weights):
    print(f"Copying best weights from {best_weights} to {dest_model}...")
    shutil.copy(best_weights, dest_model)
    print("Updated storage/models/person_detection_reinforced.pt with fine-tuned weights!")
else:
    print(f"Warning: {best_weights} not found!")
