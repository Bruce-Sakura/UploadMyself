#!/usr/bin/env python3
"""
UploadMyself ML Service
Runs on host with GPU, called by Go backend via HTTP.
Endpoints:
  POST /generate-avatar  — CharacterGen pipeline
  POST /health           — Health check
"""

import json
import os
import sys
import io
import base64
import argparse
from http.server import HTTPServer, BaseHTTPRequestHandler
from pathlib import Path

# CharacterGen paths
CHARACTERGEN_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "charactergen")
# Only add CHARACTERGEN_DIR so root webui.py (with Inference2D_API / Inference3D_API) is found
# Do NOT add 2D_Stage or 3D_Stage — their webui.py would shadow the root one
sys.path.insert(0, CHARACTERGEN_DIR)

# Global model instances (lazy loaded)
infer2d = None
infer3d = None
remove_bg = None


def load_models():
    """Lazy load CharacterGen models."""
    global infer2d, infer3d, remove_bg
    import torch
    from omegaconf import OmegaConf

    # CharacterGen must run from its own directory
    original_dir = os.getcwd()
    os.chdir(CHARACTERGEN_DIR)

    if remove_bg is None:
        print("[ML] Loading background remover...", flush=True)
        from webui import rm_bg_api
        remove_bg = rm_bg_api()

    if infer2d is None:
        print("[ML] Loading 2D model...", flush=True)
        from webui import Inference2D_API
        config_path = os.path.join(CHARACTERGEN_DIR, "2D_Stage", "configs", "infer.yaml")
        cfg = OmegaConf.load(config_path)
        infer2d = Inference2D_API(**cfg)

    if infer3d is None:
        print("[ML] Loading 3D model...", flush=True)
        from webui import Inference3D_API
        infer3d = Inference3D_API()

    os.chdir(original_dir)


def generate_avatar(input_path: str, output_dir: str, seed: int = 2333, timestep: int = 40) -> dict:
    """Full CharacterGen pipeline."""
    import torch
    import numpy as np
    from PIL import Image
    from datetime import datetime

    # CharacterGen expects to run from its own directory
    original_dir = os.getcwd()
    os.chdir(CHARACTERGEN_DIR)

    try:
        return _generate_avatar_inner(input_path, output_dir, seed, timestep)
    finally:
        os.chdir(original_dir)


def _generate_avatar_inner(input_path: str, output_dir: str, seed: int, timestep: int) -> dict:
    """Inner function that runs in CharacterGen directory."""
    from webui import traverse
    from pygltflib import GLTF2
    import shutil

    os.makedirs(output_dir, exist_ok=True)

    # Step 1: Remove background
    print("[1/3] Removing background...", flush=True)
    img = Image.open(input_path).convert("RGBA")
    img_clean = remove_bg.remove_background(
        imgs=[np.array(img)], alpha_min=0.1, alpha_max=0.9
    )[0]
    clean_path = os.path.join(output_dir, "input_clean.png")
    img_clean.save(clean_path)

    # Step 2: Generate 4 views
    print("[2/3] Generating 4 views...", flush=True)
    views = infer2d.inference(img_clean, 512, 768, crop=True, seed=seed, timestep=timestep)
    views = remove_bg.remove_background(imgs=views, alpha_min=0.2, alpha_max=0.9)

    view_names = ["back", "front", "right", "left"]
    view_paths = []
    for name, view in zip(view_names, views):
        path = os.path.join(output_dir, f"view_{name}.png")
        view.save(path)
        view_paths.append(path)

    # Step 3: 3D reconstruction
    print("[3/3] Reconstructing 3D...", flush=True)
    save_dir, obj_path, glb_path = infer3d.process_images(
        views[0], views[1], views[2], views[3],
        back_proj=False, smooth_iter=5
    )

    # Copy to output dir
    import shutil
    final_glb = os.path.join(output_dir, "avatar.glb")
    final_obj = os.path.join(output_dir, "avatar.obj")
    shutil.copy2(glb_path, final_glb)
    shutil.copy2(obj_path, final_obj)

    return {
        "cartoon_image": clean_path,
        "views": view_paths,
        "glb_model": final_glb,
        "obj_model": final_obj,
    }


class MLHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == "/generate-avatar":
            self._handle_generate_avatar()
        else:
            self._respond(404, {"error": "not found"})

    def do_GET(self):
        if self.path == "/health":
            import torch
            self._respond(200, {
                "status": "ok",
                "gpu": torch.cuda.get_device_name(0) if torch.cuda.is_available() else "none",
                "cuda": torch.cuda.is_available(),
            })
        else:
            self._respond(404, {"error": "not found"})

    def _handle_generate_avatar(self):
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        data = json.loads(body)

        input_path = data.get("input_path", "")
        output_dir = data.get("output_dir", "")
        seed = data.get("seed", 2333)
        timestep = data.get("timestep", 40)

        if not input_path or not output_dir:
            self._respond(400, {"error": "input_path and output_dir required"})
            return

        try:
            load_models()
            result = generate_avatar(input_path, output_dir, seed, timestep)
            self._respond(200, result)
        except Exception as e:
            self._respond(500, {"error": str(e)})

    def _respond(self, code, data):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps(data).encode())

    def log_message(self, format, *args):
        print(f"[ML] {args[0]}", flush=True)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--port", type=int, default=8001)
    parser.add_argument("--host", default="0.0.0.0")
    args = parser.parse_args()

    print(f"[ML Service] Starting on {args.host}:{args.port}", flush=True)
    print(f"[ML Service] CharacterGen dir: {CHARACTERGEN_DIR}", flush=True)

    server = HTTPServer((args.host, args.port), MLHandler)
    server.serve_forever()


if __name__ == "__main__":
    main()
