import os
import shutil
import cv2
import time
import random
import glob
from ultralytics import YOLO

base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
out_dir = os.path.join(base_dir, "datasets", "cctv_person_mix")
os.makedirs(os.path.join(out_dir, "images", "train"), exist_ok=True)
os.makedirs(os.path.join(out_dir, "images", "val"), exist_ok=True)
os.makedirs(os.path.join(out_dir, "labels", "train"), exist_ok=True)
os.makedirs(os.path.join(out_dir, "labels", "val"), exist_ok=True)

print("1. Extracting diverse person images from COCO128...")
coco_labels = glob.glob(os.path.join(base_dir, "datasets/coco128/labels/train2017/*.txt"))
coco_person_samples = []
coco_bg_samples = []

for l_path in coco_labels:
    base_name = os.path.basename(l_path).replace(".txt", "")
    img_path = os.path.join(base_dir, "datasets/coco128/images/train2017", f"{base_name}.jpg")
    if not os.path.exists(img_path):
        continue

    with open(l_path, "r") as f:
        lines = [line.strip() for line in f if line.strip()]

    # Check if contains person (class 0)
    person_boxes = []
    for line in lines:
        parts = line.split()
        if parts[0] == "0":
            # keep as class 0
            person_boxes.append(" ".join(["0"] + parts[1:]))

    if len(person_boxes) > 0:
        coco_person_samples.append((img_path, person_boxes))
    else:
        coco_bg_samples.append((img_path, []))

print(f"Found {len(coco_person_samples)} COCO person images and {len(coco_bg_samples)} background images.")

# Select 50 person images and 15 background images
random.seed(42)
random.shuffle(coco_person_samples)
random.shuffle(coco_bg_samples)
selected_coco_person = coco_person_samples[:50]
selected_coco_bg = coco_bg_samples[:15]

# Add to dataset (80% train, 20% val)
for i, (img_p, boxes) in enumerate(selected_coco_person + selected_coco_bg):
    split = "val" if (i % 5 == 0) else "train"
    fn = f"coco_{os.path.basename(img_p)}"
    shutil.copy(img_p, os.path.join(out_dir, "images", split, fn))
    lbl_file = os.path.join(out_dir, "labels", split, fn.replace(".jpg", ".txt"))
    with open(lbl_file, "w") as f:
        if boxes:
            f.write("\n".join(boxes) + "\n")

print(f"Added {len(selected_coco_person) + len(selected_coco_bg)} COCO diverse images.")

# 2. Capture live frames from security cameras
print("\n2. Capturing frames from real security center cameras...")
os.environ['OPENCV_FFMPEG_CAPTURE_OPTIONS'] = 'rtsp_transport;tcp|fflags;nobuffer|flags;low_delay'

# Load current model to assist with auto-labeling
detector = YOLO("storage/models/person_detection_reinforced.pt")

active_cams = [
    ("live_d38", "rtsp://k-monitorcenter.app:8554/live/live_d38ygmd5gzm1", True, 8),
    ("live_g8q", "rtsp://k-monitorcenter.app:8554/live/live_g8q6t3d6jgjt", True, 8),
    ("live_v0m", "rtsp://k-monitorcenter.app:8554/live/live_v0m7gz9przaz", False, 5), # Background
    ("live_k4x", "rtsp://k-monitorcenter.app:8554/live/live_k4xvpbtdqji8", False, 5), # Background
]

for cam_id, url, is_positive, n_frames in active_cams:
    print(f"Capturing {n_frames} frames from {cam_id} ({url})...")
    cap = cv2.VideoCapture(url, cv2.CAP_FFMPEG)
    cap.set(cv2.CAP_PROP_OPEN_TIMEOUT_MSEC, 5000)
    cap.set(cv2.CAP_PROP_READ_TIMEOUT_MSEC, 5000)

    captured = 0
    attempts = 0
    while captured < n_frames and attempts < 25:
        attempts += 1
        ret, frame = cap.read()
        if not ret or frame is None:
            time.sleep(0.2)
            continue

        split = "val" if (captured == 0) else "train"
        fn = f"cam_{cam_id}_{captured+1}_{int(time.time()*1000)}.jpg"
        img_out = os.path.join(out_dir, "images", split, fn)
        lbl_out = os.path.join(out_dir, "labels", split, fn.replace(".jpg", ".txt"))

        cv2.imwrite(img_out, frame)

        if is_positive:
            # Detect person with model
            res = detector(frame, conf=0.25, verbose=False)
            boxes = []
            if len(res) > 0 and len(res[0].boxes) > 0:
                for b in res[0].boxes:
                    cls_name = res[0].names[int(b.cls[0].item())]
                    if cls_name == "person":
                        xywhn = b.xywhn[0].tolist()
                        boxes.append(f"0 {xywhn[0]:.6f} {xywhn[1]:.6f} {xywhn[2]:.6f} {xywhn[3]:.6f}")

            with open(lbl_out, "w") as f:
                if boxes:
                    f.write("\n".join(boxes) + "\n")
            print(f"  -> Captured positive frame {captured+1} ({len(boxes)} persons)")
        else:
            # Negative background image (empty label file)
            with open(lbl_out, "w") as f:
                f.write("")
            print(f"  -> Captured negative background frame {captured+1}")

        captured += 1
        time.sleep(0.4)

    cap.release()

# 3. Create data.yaml
clean_path = os.path.abspath(out_dir).replace('\\', '/')
data_yaml_content = f"""path: {clean_path}
train: images/train
val: images/val

names:
  0: person
"""
yaml_path = os.path.join(out_dir, "data.yaml")
with open(yaml_path, "w") as f:
    f.write(data_yaml_content)

print(f"\nDataset successfully built at: {out_dir}")
print(f"YAML config saved to: {yaml_path}")
train_imgs = glob.glob(os.path.join(out_dir, "images/train/*"))
val_imgs = glob.glob(os.path.join(out_dir, "images/val/*"))
print(f"Summary: {len(train_imgs)} train images, {len(val_imgs)} val images.")
