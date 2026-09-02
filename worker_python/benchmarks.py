#!/usr/bin/env python3
"""
HydraForge - Ultralytics YOLO Multi-Format Benchmark Runner
Invokes `ultralytics.utils.benchmarks.benchmark` to evaluate model latency, throughput,
size, and mAP across export formats (PyTorch, ONNX, TensorRT, OpenVINO, etc.).
"""

import argparse
import json
import os
import sys

def run_benchmark(model: str, data: str, imgsz: int, quantize: int, device: str, format_name: str, verbose: bool = False):
    try:
        from ultralytics.utils.benchmarks import benchmark
    except ImportError:
        print(json.dumps({
            "error": "ultralytics package not installed. Run: pip install 'ultralytics[export]'"
        }))
        sys.exit(1)

    # Clean quantization parameter (None if 0/FP32, 16 for FP16, 8 for INT8)
    quant_param = None
    if quantize in (8, 16):
        quant_param = quantize

    fmt_param = format_name if format_name else ""

    try:
        results_df = benchmark(
            model=model,
            data=data,
            imgsz=imgsz,
            quantize=quant_param,
            device=device,
            format=fmt_param,
            verbose=verbose
        )

        formatted_results = []
        if results_df is not None:
            for _, row in results_df.iterrows():
                formatted_results.append({
                    "format": str(row.get("Format", "")),
                    "status": "SUCCESS" if row.get("Status", False) else "FAILED",
                    "size_mb": float(row.get("Size (MB)", 0.0)),
                    "inference_time_ms": float(row.get("Inference time (ms/im)", 0.0)),
                    "fps": float(row.get("FPS", 0.0)),
                    "map50_95": float(row.get("mAP50-95(B)", row.get("mAP50-95(M)", 0.0))),
                    "accuracy_top1": float(row.get("accuracy_top1", 0.0)),
                    "export_args": str(row.get("Arguments", ""))
                })

        print(json.dumps({
            "status": "SUCCESS",
            "results": formatted_results
        }))
    except Exception as e:
        print(json.dumps({
            "status": "FAILED",
            "error": str(e)
        }))
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="Ultralytics YOLO Benchmark Suite CLI")
    parser.add_argument("--model", type=str, default="yolo26n.pt", help="Model path (.pt or .yaml)")
    parser.add_argument("--data", type=str, default="coco8.yaml", help="Dataset yaml path")
    parser.add_argument("--imgsz", type=int, default=640, help="Image size (square)")
    parser.add_argument("--quantize", type=int, default=0, help="Quantization (16 for FP16, 8 for INT8)")
    parser.add_argument("--device", type=str, default="0", help="Computation device (0, cuda:0, cpu)")
    parser.add_argument("--format", type=str, default="", help="Specific export format (onnx, engine, openvino, etc.)")
    parser.add_argument("--verbose", action="store_true", help="Verbose output")

    args = parser.parse_args()
    run_benchmark(
        model=args.model,
        data=args.data,
        imgsz=args.imgsz,
        quantize=args.quantize,
        device=args.device,
        format_name=args.format,
        verbose=args.verbose
    )

if __name__ == "__main__":
    main()
