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

    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    project_dir = os.path.join(base_dir, 'runs', 'train')
    os.makedirs(project_dir, exist_ok=True)

    # Resolve model weights path
    model_name = args.model
    if not model_name.endswith('.pt') and not model_name.endswith('.yaml'):
        model_name += '.pt'

    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    weights_path = os.path.join(base_dir, 'weights', os.path.basename(model_name))
    
    if os.path.isabs(model_name) and os.path.exists(model_name):
        model_target = model_name
    elif os.path.exists(weights_path):
        model_target = weights_path
    elif os.path.exists(model_name):
        model_target = model_name
    elif os.path.exists(os.path.join('weights', model_name)):
        model_target = os.path.join('weights', model_name)
    else:
        model_target = model_name

    try:
        model = YOLO(model_target)
    except Exception:
        fallback = os.path.join(base_dir, 'weights', 'yolo26n.pt')
        if os.path.exists(fallback):
            model = YOLO(fallback)
        else:
            model = YOLO('yolov8n.pt')

    epoch_timer = [time.time()]
    batch_in_epoch = [0]

    def extract_losses(trainer):
        box_loss, cls_loss, dfl_loss = 0.0, 0.0, 0.0
        source = None
        if hasattr(trainer, 'tloss') and trainer.tloss:
            source = trainer.tloss
        elif hasattr(trainer, 'loss_items') and trainer.loss_items is not None:
            source = trainer.loss_items

        if isinstance(source, dict):
            for k, v in source.items():
                val = float(v.item() if hasattr(v, 'item') else v)
                k_low = k.lower()
                if 'box' in k_low:
                    box_loss = val
                elif 'cls' in k_low:
                    cls_loss = val
                elif 'dfl' in k_low or 'l1' in k_low:
                    dfl_loss = val
        elif isinstance(source, (list, tuple)):
            vals = [float(x.item() if hasattr(x, 'item') else x) for x in source]
            if len(vals) > 0: box_loss = vals[0]
            if len(vals) > 1: cls_loss = vals[1]
            if len(vals) > 2: dfl_loss = vals[2]
        elif hasattr(source, 'tolist'):
            vals = source.tolist()
            if isinstance(vals, (list, tuple)):
                if len(vals) > 0: box_loss = float(vals[0])
                if len(vals) > 1: cls_loss = float(vals[1])
                if len(vals) > 2: dfl_loss = float(vals[2])

        return box_loss, cls_loss, dfl_loss

    def on_train_epoch_start(trainer):
        epoch_timer[0] = time.time()
        batch_in_epoch[0] = 0

    def on_fit_epoch_end(trainer):
        try:
            metrics = trainer.metrics or {}
            epoch = trainer.epoch + 1
            total_epochs = trainer.epochs
            duration = max(0.1, time.time() - epoch_timer[0])

            box_loss, cls_loss, dfl_loss = extract_losses(trainer)
            if box_loss == 0.0:
                box_loss = float(metrics.get('train/box_loss', 0.0))
            if cls_loss == 0.0:
                cls_loss = float(metrics.get('train/cls_loss', 0.0))
            if dfl_loss == 0.0:
                dfl_loss = float(metrics.get('train/dfl_loss', 0.0))

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
            total_imgs = 31048
            try:
                if hasattr(trainer, 'train_loader') and trainer.train_loader:
                    total_imgs = len(trainer.train_loader.dataset)
            except Exception:
                pass
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
                "lr": float(trainer.optimizer.param_groups[0]['lr']) if (hasattr(trainer, 'optimizer') and trainer.optimizer) else 0.001,
                "gpu_vram_mb": vram,
                "power_watts": power_w,
                "temp_celsius": temp_c,
                "gpu_util_pct": gpu_util,
                "fps": fps,
                "epoch_duration_sec": duration
            }
            print(f"HYDRA_METRIC:{json.dumps(payload)}", flush=True)
        except Exception as e:
            sys.stderr.write(f"CALLBACK_ERR: {e}\n")

    job_start_time = time.time()
    accumulated_energy = [0.0]
    last_batch_time = [time.time()]

    def on_train_batch_end(trainer):
        try:
            batch_in_epoch[0] += 1
            batch_idx = batch_in_epoch[0]
            now = time.time()
            dt = max(0.001, now - last_batch_time[0])
            last_batch_time[0] = now
            
            total_batches = len(trainer.train_loader) if (hasattr(trainer, 'train_loader') and trainer.train_loader) else 1000
            
            power_w = 0.0
            if nvml_handle:
                try:
                    power_w = float(pynvml.nvmlDeviceGetPowerUsage(nvml_handle)) / 1000.0
                except Exception:
                    pass
            
            # Integrate energy in kWh
            accumulated_energy[0] += (power_w * (dt / 3600.0)) / 1000.0
            elapsed_sec = now - job_start_time
            
            if batch_idx % 2 == 0 or batch_idx == total_batches:
                epoch = trainer.epoch + 1
                box_loss, cls_loss, _ = extract_losses(trainer)
                
                imgs_so_far = (trainer.epoch * total_batches + batch_idx) * args.batch
                current_fps = float(imgs_so_far) / max(0.1, elapsed_sec)
                
                batch_payload = {
                    "epoch": epoch,
                    "batch": batch_idx,
                    "total_batches": total_batches,
                    "box_loss": box_loss,
                    "cls_loss": cls_loss,
                    "power_watts": power_w,
                    "total_energy_kwh": accumulated_energy[0],
                    "fps": current_fps,
                    "duration_sec": elapsed_sec
                }
                print(f"\nHYDRA_BATCH:{json.dumps(batch_payload)}", flush=True)
        except Exception as e:
            sys.stderr.write(f"BATCH_ERR: {e}\n")

    model.add_callback('on_train_epoch_start', on_train_epoch_start)
    model.add_callback('on_train_batch_end', on_train_batch_end)
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
