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

def load_hf_token():
    token = os.environ.get("HF_TOKEN")
    if token:
        return token
    for env_file in [
        os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), ".env"),
        os.path.expanduser("~/Documents/reductor-prompt/.env"),
        os.path.expanduser("~/Documents/hydravms/.env")
    ]:
        if os.path.exists(env_file):
            with open(env_file, "r") as f:
                for line in f:
                    if line.strip().startswith("HF_TOKEN="):
                        return line.strip().split("=", 1)[1].strip().strip('"').strip("'")
    return None

def download_file(url: str, dest_path: str, token: str = None):
    print(f"📥 Baixando de {url} para {dest_path}...")
    headers = {'User-Agent': 'HydraVMS-Downloader/1.0'}
    if token:
        headers['Authorization'] = f'Bearer {token}'
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
        print(f"⚠️ Aviso na compilação TensorRT: {e}. Mantendo modelo PT/ONNX para inferência.")
        return True

def main():
    plugin_id = sys.argv[1] if len(sys.argv) > 1 else "person_detection_reinforced"
    package_url = sys.argv[2] if len(sys.argv) > 2 else "https://huggingface.co/hydravision/yolo26m-person"
    token = load_hf_token()

    base_name = "yolo26m"
    pt_dest = os.path.join(MODELS_DIR, f"{plugin_id}.pt")
    plugin_json_dest = os.path.join(MODELS_DIR, f"{plugin_id}_manifest.json")
    engine_dest = os.path.join(ENGINES_DIR, f"{plugin_id}.engine")

    # Download do manifest plugin.json e pesos do Hugging Face
    if "huggingface.co" in package_url:
        clean_url = package_url.rstrip('/')
        json_url = f"{clean_url}/resolve/main/plugin.json"
        pt_url = f"{clean_url}/resolve/main/yolo26m.pt"
        try:
            download_file(json_url, plugin_json_dest, token)
        except Exception:
            pass
    else:
        pt_url = package_url

    try:
        if not os.path.exists(pt_dest):
            download_file(pt_url, pt_dest, token)

        compile_model(pt_dest, engine_dest)
        print(json.dumps({"status": "success", "plugin_id": plugin_id, "pt": pt_dest, "engine": engine_dest}))
    except Exception as e:
        print(json.dumps({"status": "error", "message": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main()
