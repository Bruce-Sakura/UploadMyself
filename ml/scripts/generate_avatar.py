#!/usr/bin/env python3
"""
UploadMyself Avatar Pipeline (CharacterGen)
Photo → Remove BG → 4-View Generation → 3D Reconstruction → GLB/VRM

Usage:
  python generate_avatar.py --input photo.jpg --output-dir ./output --style cartoon
"""

import argparse
import json
import os
import sys
import subprocess
import shutil
from pathlib import Path

# CharacterGen paths
CHARACTERGEN_DIR = os.path.join(os.path.dirname(__file__), "..", "charactergen")
STAGE_2D_DIR = os.path.join(CHARACTERGEN_DIR, "2D_Stage")
STAGE_3D_DIR = os.path.join(CHARACTERGEN_DIR, "3D_Stage")


def check_environment():
    """Check if CharacterGen dependencies are available."""
    try:
        import torch
        if not torch.cuda.is_available():
            return False, "CUDA not available"
        return True, f"GPU: {torch.cuda.get_device_name(0)}"
    except ImportError:
        return False, "PyTorch not installed"


def remove_background(input_path: str, output_path: str) -> str:
    """Remove background from input image using anime-seg."""
    try:
        from rm_anime_bg.cli import get_mask, SCALE
        import onnxruntime as rt
        from huggingface_hub import hf_hub_download
        import cv2
        import numpy as np
        from PIL import Image

        # Load ONNX model
        session_path = hf_hub_download(
            repo_id="skytnt/anime-seg", filename="isnetis.onnx"
        )
        providers = ["CPUExecutionProvider"]
        if "CUDAExecutionProvider" in rt.get_available_providers():
            providers = ["CUDAExecutionProvider"]
        session = rt.InferenceSession(session_path, providers=providers)

        # Process image
        img = cv2.imread(input_path)
        if img is None:
            raise ValueError(f"Cannot read image: {input_path}")
        img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)

        mask = get_mask(session, img)
        mask[mask < 0.1] = 0.0
        mask[mask > 0.9] = 1.0

        img_rgba = np.concatenate([
            (mask[..., None] * img).astype(np.uint8),
            (mask * SCALE).astype(np.uint8)[..., None]
        ], axis=2)

        Image.fromarray(img_rgba).save(output_path)
        return output_path
    except Exception as e:
        print(f"Background removal failed: {e}, using original", file=sys.stderr)
        shutil.copy(input_path, output_path)
        return output_path


def generate_4views(input_path: str, output_dir: str, seed: int = 2333, timestep: int = 40) -> list:
    """Generate 4-view images using CharacterGen 2D stage."""
    sys.path.insert(0, STAGE_2D_DIR)
    sys.path.insert(0, CHARACTERGEN_DIR)

    from omegaconf import OmegaConf
    from webui import Inference2D_API, rm_bg_api, process_image
    from PIL import Image
    import numpy as np

    # Load config
    config_path = os.path.join(STAGE_2D_DIR, "configs", "infer.yaml")
    if not os.path.exists(config_path):
        raise FileNotFoundError(f"CharacterGen 2D config not found: {config_path}")

    cfg = OmegaConf.load(config_path)

    # Initialize 2D inference
    print("[2D] Loading model...", file=sys.stderr)
    infer2d = Inference2D_API(**cfg)
    remove_api = rm_bg_api()

    # Load and process input image
    input_img = Image.open(input_path).convert("RGBA")

    # Remove background
    print("[2D] Removing background...", file=sys.stderr)
    input_img = remove_api.remove_background(
        imgs=[np.array(input_img)], alpha_min=0.1, alpha_max=0.9
    )[0]

    # Generate 4 views
    print("[2D] Generating 4 views...", file=sys.stderr)
    views = infer2d.inference(
        input_img, 512, 768, crop=True, seed=seed, timestep=timestep
    )

    # Remove background from generated views
    print("[2D] Post-processing views...", file=sys.stderr)
    views = remove_api.remove_background(imgs=views, alpha_min=0.2, alpha_max=0.9)

    # Save views
    view_names = ["back", "front", "right", "left"]
    view_paths = []
    for i, (name, view) in enumerate(zip(view_names, views)):
        path = os.path.join(output_dir, f"view_{name}.png")
        view.save(path)
        view_paths.append(path)
        print(f"[2D] Saved {name} view: {path}", file=sys.stderr)

    return view_paths


def reconstruct_3d(view_paths: list, output_dir: str, smooth_iter: int = 5) -> dict:
    """Reconstruct 3D mesh from 4 views using CharacterGen 3D stage."""
    sys.path.insert(0, STAGE_3D_DIR)
    sys.path.insert(0, CHARACTERGEN_DIR)

    from omegaconf import OmegaConf
    from webui import Inference3D_API
    from PIL import Image

    # Initialize 3D inference
    print("[3D] Loading model...", file=sys.stderr)
    infer3d = Inference3D_API()

    # Load views
    views = [Image.open(p).convert("RGBA") for p in view_paths]

    # Reconstruct
    print("[3D] Reconstructing 3D mesh...", file=sys.stderr)
    save_dir, obj_path, glb_path = infer3d.process_images(
        views[0], views[1], views[2], views[3],
        back_proj=True, smooth_iter=smooth_iter
    )

    # Copy output to our output_dir
    final_glb = os.path.join(output_dir, "avatar.glb")
    final_obj = os.path.join(output_dir, "avatar.obj")
    shutil.copy2(glb_path, final_glb)
    shutil.copy2(obj_path, final_obj)

    return {
        "glb_path": final_glb,
        "obj_path": final_obj,
        "save_dir": save_dir,
    }


def main():
    parser = argparse.ArgumentParser(description="UploadMyself Avatar Generation (CharacterGen)")
    parser.add_argument("--input", required=True, help="Input photo path")
    parser.add_argument("--output-dir", required=True, help="Output directory")
    parser.add_argument("--style", default="realistic", choices=["realistic", "cartoon", "anime"])
    parser.add_argument("--seed", type=int, default=2333)
    parser.add_argument("--timestep", type=int, default=40, help="Diffusion timesteps (10-70)")
    parser.add_argument("--smooth", type=int, default=5, help="Mesh smoothing iterations")
    parser.add_argument("--skip-views", action="store_true", help="Skip 4-view generation (use existing views)")
    args = parser.parse_args()

    os.makedirs(args.output_dir, exist_ok=True)

    # Check environment
    ok, msg = check_environment()
    print(f"Environment: {msg}", file=sys.stderr)
    if not ok:
        # Fallback: output just the photo with a note
        result = {
            "error": msg,
            "fallback": True,
            "photo_path": args.input,
        }
        print(json.dumps(result))
        sys.exit(1)

    # Step 1: Remove background
    print("[1/3] Removing background...", file=sys.stderr)
    clean_path = os.path.join(args.output_dir, "input_clean.png")
    remove_background(args.input, clean_path)

    # Step 2: Generate 4 views
    if args.skip_views:
        print("[2/3] Using existing views...", file=sys.stderr)
        view_paths = sorted([
            os.path.join(args.output_dir, f)
            for f in os.listdir(args.output_dir)
            if f.startswith("view_") and f.endswith(".png")
        ])
    else:
        print("[2/3] Generating 4 views (this may take a few minutes)...", file=sys.stderr)
        view_paths = generate_4views(clean_path, args.output_dir, args.seed, args.timestep)

    # Step 3: 3D reconstruction
    print("[3/3] Reconstructing 3D model...", file=sys.stderr)
    result = reconstruct_3d(view_paths, args.output_dir, args.smooth)

    # Output result
    output = {
        "cartoon_image": clean_path,
        "views": view_paths,
        "glb_model": result["glb_path"],
        "obj_model": result["obj_path"],
        "has_skeleton": False,  # CharacterGen outputs mesh only, rigging is separate
        "animations": [],
    }
    print(json.dumps(output))


if __name__ == "__main__":
    main()
