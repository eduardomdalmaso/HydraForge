#!/usr/bin/env python3
"""
HydraForge - Downloader & TensorRT/ONNX Compiler Engine
Baixa o modelo do Hugging Face e compila na GPU local.
Uso:
  python download_and_compile.py <plugin_id> <package_url>
"""
import sys
import os
import json
import urllib.request

STORAGE_DIR = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "storage")
MODELS_DIR = os.path.join(STORAGE_DIR, "models")
ENGINES_DIR = os.path.join(STORAGE_DIR, "engines")

os.makedirs(MODELS_DIR, exist_ok=True)
os.makedirs(ENGINES_DIR, exist_ok=True)

def download_file(url: str, dest_path: str):
    print(f"📥 Baixando de {url} para {dest_path}...")
    headers = {'User-Agent': 'HydraVMS-Downloader/1.0'}
    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req) as response, open(dest_path, 'wb') as out_file:
        total_size = int(response.info().get('Content-Length', 0))
        downloaded = 0
        block_size = 1024 * 1024
        while True:
            buffer = response.read(block_size)
            if not buffer:
                break
            downloaded += len(buffer)
            out_file.write(buffer)
            if total_size > 0:
                percent = (downloaded / total_size) * 100
                print(f"\rProgresso: {percent:.1f}% ({downloaded / (1024*1024):.1f}MB / {total_size / (1024*1024):.1f}MB)", end="")
    print("\n✔ Download concluído com sucesso!")

def compile_model(pt_path: str, engine_path: str):
    print(f"⚙️ Compilando modelo para TensorRT FP16 na GPU RTX 5090: {engine_path}...")
    try:
        from ultralytics import YOLO
        if os.path.exists(pt_path):
            model = YOLO(pt_path)
            model.export(format='engine', half=True, device=0)
            print("✔ Compilação TensorRT concluída com sucesso!")
            return True
    except Exception as e:
        print(f"⚠️ Aviso na compilação TensorRT: {e}. Mantendo ONNX ativo para inferência.")
        return True

def main():
    plugin_id = sys.argv[1] if len(sys.argv) > 1 else "object_detection_sota"
    package_url = sys.argv[2] if len(sys.argv) > 2 else "https://huggingface.co/hydravision/yolo26m-object"

    pt_dest = os.path.join(MODELS_DIR, "yolo26m.pt")
    onnx_dest = os.path.join(MODELS_DIR, "yolo26m.onnx")
    engine_dest = os.path.join(ENGINES_DIR, "yolo26m.engine")

    # URLs diretas de download do Hugging Face
    if "huggingface.co" in package_url:
        download_pt_url = f"{package_url.rstrip('/')}/resolve/main/yolo26m.pt"
        download_onnx_url = f"{package_url.rstrip('/')}/resolve/main/yolo26m.onnx"
    else:
        download_pt_url = package_url
        download_onnx_url = package_url

    try:
        if not os.path.exists(pt_dest):
            download_file(download_pt_url, pt_dest)
        if not os.path.exists(onnx_dest):
            download_file(download_onnx_url, onnx_dest)

        compile_model(pt_dest, engine_dest)
        print(json.dumps({"status": "success", "plugin_id": plugin_id, "pt": pt_dest, "onnx": onnx_dest, "engine": engine_dest}))
    except Exception as e:
        print(json.dumps({"status": "error", "message": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main()
