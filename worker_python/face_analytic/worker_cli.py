#!/usr/bin/env python3
"""
HydraForge Face Recognition Live Streaming Worker (RTX 5090).
Executes Person-Gated Face Analytics pipeline on RTSP / Zero-Copy streams.
"""
import sys
import os
import time
import json
import argparse
import numpy as np

from .config import FaceAnalyticConfig
from .pipeline import FaceAnalyticPipeline


def parse_args():
    parser = argparse.ArgumentParser(description="HydraForge Live Face Analytic Worker")
    parser.add_argument("--stream-id", type=str, default="cam_01")
    parser.add_argument("--rtsp-url", type=str, default="rtsp://127.0.0.1:8554/cam_01")
    parser.add_argument("--min-height", type=int, default=40)
    parser.add_argument("--max-height", type=int, default=0, help="0 = unbounded / infinity")
    parser.add_argument("--min-faces", type=int, default=3)
    parser.add_argument("--blur-threshold", type=float, default=60.0)
    parser.add_argument("--unknown-score", type=float, default=0.60)
    parser.add_argument("--recognition-threshold", type=float, default=0.75)
    parser.add_argument("--adaptive-blur", action="store_true", default=True)
    return parser.parse_args()


def main():
    args = parse_args()
    config = FaceAnalyticConfig(
        camera_id=args.stream_id,
        min_height=args.min_height,
        max_height=args.max_height,
        min_faces=args.min_faces,
        blur_threshold=args.blur_threshold,
        unknown_score=args.unknown_score,
        recognition_threshold=args.recognition_threshold,
        adaptive_blur=args.adaptive_blur
    )
    config.validate()

    pipeline = FaceAnalyticPipeline(config)
    shm_out_path = f"/dev/shm/face_detections_{args.stream_id}.json"
    shm_tmp_path = f"/dev/shm/face_detections_{args.stream_id}.json.tmp"

    print(
        f"🚀 [HYDRA-FACE] Worker online for {args.stream_id} "
        f"[min_h={config.min_height}, max_h={config.max_height}, min_faces={config.min_faces}]",
        flush=True
    )

    # In production, this loop connects to HydraStream's Zero-Copy /dev/shm ring buffer
    # or MediaMTX RTSP stream.
    last_tick = time.time()
    while True:
        try:
            # Heartbeat telemetry
            now = time.time()
            if now - last_tick >= 1.0:
                telemetry = {
                    "status": "running",
                    "timestamp": now,
                    "stream_id": args.stream_id,
                    "active_tracks": len(pipeline.tracker.tracks),
                    "known_gallery_count": len(pipeline.matcher.known_gallery),
                    "active_unknowns_count": len(pipeline.matcher.active_unknowns)
                }
                with open(shm_tmp_path, "w") as f:
                    json.dump(telemetry, f)
                os.replace(shm_tmp_path, shm_out_path)
                last_tick = now
            time.sleep(0.05)
        except KeyboardInterrupt:
            print("\n🛑 [HYDRA-FACE] Worker stopping cleanly...", flush=True)
            break
        except Exception as e:
            print(f"⚠️ [HYDRA-FACE] Loop exception: {e}", file=sys.stderr, flush=True)
            time.sleep(0.1)


if __name__ == "__main__":
    main()
