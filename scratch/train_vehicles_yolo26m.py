import os
from ultralytics import YOLO

def main():
    data_yaml = "datasets/vehicles_unified/data.yaml"
    model_name = "yolo26m.pt"

    print("==================================================")
    print("=== INITIATING VEHICLE DETECTION TRAINING (YOLO26m) ===")
    print("==================================================")
    print(f"Base Architecture: {model_name}")
    print(f"Dataset: {data_yaml}")
    print(f"Target Classes: 0: car, 1: motorcycle, 2: bus, 3: truck")

    # Load YOLO26m pre-trained checkpoint
    model = YOLO(model_name)

    results = model.train(
        data=data_yaml,
        epochs=20,
        imgsz=640,
        batch=32,
        device=0,
        workers=0,
        optimizer='AdamW',
        lr0=0.001,
        lrf=0.01,
        warmup_epochs=2,
        mosaic=0.5,
        mixup=0.1,
        fliplr=0.5,
        project='runs/vehicles_yolo26m',
        name='v1_all_vehicles',
        exist_ok=True,
        verbose=True
    )

    best_path = "runs/vehicles_yolo26m/v1_all_vehicles/weights/best.pt"
    print(f"\n=== VEHICLE MODEL TRAINING FINISHED ===")
    print(f"Best model saved to: {best_path}")

if __name__ == '__main__':
    main()
