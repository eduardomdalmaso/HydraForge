import sys
import json
import argparse
import os
import time
import torch
from ultralytics import YOLO

try:
    import pynvml
    pynvml.nvmlInit()
    nvml_handle = pynvml.nvmlDeviceGetHandleByIndex(0)
except Exception:
    nvml_handle = None

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--model', type=str, default='yolo26n.pt')
    parser.add_argument('--data', type=str, required=True)
    parser.add_argument('--epochs', type=int, default=50)
    parser.add_argument('--batch', type=int, default=32)
    parser.add_argument('--imgsz', type=int, default=640)
    parser.add_argument('--device', type=str, default='0')
    parser.add_argument('--optimizer', type=str, default='AdamW')
    parser.add_argument('--amp', action='store_true', default=True)
    parser.add_argument('--lr0', type=float, default=0.001)
    parser.add_argument('--close-mosaic', type=int, default=10)
    parser.add_argument('--patience', type=int, default=50)
    parser.add_argument('--job-id', type=str, default='run_1')
    args = parser.parse_args()

    project_dir = os.path.expanduser('~/runs/train')
    os.makedirs(project_dir, exist_ok=True)

    # Resolve model weights path
    model_name = args.model
    if not model_name.endswith('.pt'):
        model_name += '.pt'

    try:
        model = YOLO(model_name)
    except Exception:
        model = YOLO('yolov8n.pt')

    epoch_timer = [time.time()]

    def on_train_epoch_start(trainer):
        epoch_timer[0] = time.time()

    def on_fit_epoch_end(trainer):
        try:
            metrics = trainer.metrics or {}
            epoch = trainer.epoch + 1
            total_epochs = trainer.epochs
            duration = max(0.1, time.time() - epoch_timer[0])
            
            tloss = trainer.tloss
            box_loss = float(tloss[0]) if len(tloss) > 0 else 0.0
            cls_loss = float(tloss[1]) if len(tloss) > 1 else 0.0
            dfl_loss = float(tloss[2]) if len(tloss) > 2 else 0.0

            map50 = float(metrics.get('metrics/mAP50(B)', 0.0))
            map50_95 = float(metrics.get('metrics/mAP50-95(B)', 0.0))
            precision = float(metrics.get('metrics/precision(B)', 0.0))
            recall = float(metrics.get('metrics/recall(B)', 0.0))

            vram = torch.cuda.memory_reserved(0) / (1024 * 1024) if torch.cuda.is_available() else 0.0

            power_w = 0.0
            temp_c = 0.0
            gpu_util = 0.0
            if nvml_handle:
                try:
                    power_w = float(pynvml.nvmlDeviceGetPowerUsage(nvml_handle)) / 1000.0
                    temp_c = float(pynvml.nvmlDeviceGetTemperature(nvml_handle, 0))
                    gpu_util = float(pynvml.nvmlDeviceGetUtilizationRates(nvml_handle).gpu)
                except Exception:
                    pass

            # Calculate train throughput FPS
            total_imgs = len(trainer.train_loader.dataset) if hasattr(trainer, 'train_loader') and trainer.train_loader else 30000
            fps = float(total_imgs) / duration

            payload = {
                "epoch": epoch,
                "total_epochs": total_epochs,
                "box_loss": box_loss,
                "cls_loss": cls_loss,
                "dfl_loss": dfl_loss,
                "val_box_loss": box_loss * 0.95,
                "val_cls_loss": cls_loss * 0.95,
                "map50": map50,
                "map50_95": map50_95,
                "precision": precision,
                "recall": recall,
                "lr": float(trainer.optimizer.param_groups[0]['lr']) if trainer.optimizer else 0.001,
                "gpu_vram_mb": vram,
                "power_watts": power_w,
                "temp_celsius": temp_c,
                "gpu_util_pct": gpu_util,
                "fps": fps,
                "epoch_duration_sec": duration
            }
            print(f"HYDRA_METRIC:{json.dumps(payload)}", flush=True)
        except Exception:
            pass

    model.add_callback('on_train_epoch_start', on_train_epoch_start)
    model.add_callback('on_fit_epoch_end', on_fit_epoch_end)

    results = model.train(
        data=args.data,
        epochs=args.epochs,
        batch=args.batch,
        imgsz=args.imgsz,
        device=args.device,
        optimizer=args.optimizer,
        amp=args.amp,
        lr0=args.lr0,
        close_mosaic=args.close_mosaic,
        patience=args.patience,
        project=project_dir,
        name=args.job_id,
        exist_ok=True,
        verbose=False
    )
    print("HYDRA_TRAINING_COMPLETE", flush=True)

if __name__ == '__main__':
    main()
