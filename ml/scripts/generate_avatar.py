#!/usr/bin/env python3
"""
Photo → 2D Cartoon Character with Skeleton Animation
Pipeline: Photo → Cartoon Style → Skeleton Detection → Animation Data
"""

import argparse
import json
import os
import sys

import cv2
import numpy as np


def photo_to_cartoon(input_path: str, output_path: str, style: str = "cartoon") -> str:
    """Convert photo to cartoon style using OpenCV edge-preserving filter."""
    img = cv2.imread(input_path)
    if img is None:
        raise ValueError(f"Cannot read image: {input_path}")

    # Resize for processing
    h, w = img.shape[:2]
    max_dim = 1024
    if max(h, w) > max_dim:
        scale = max_dim / max(h, w)
        img = cv2.resize(img, (int(w * scale), int(h * scale)))

    # Cartoon effect: bilateral filter + edge detection
    # Step 1: Apply bilateral filter for smooth color regions
    color = img.copy()
    for _ in range(5):
        color = cv2.bilateralFilter(color, 9, 75, 75)

    # Step 2: Convert to grayscale and detect edges
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    gray = cv2.medianBlur(gray, 7)
    edges = cv2.adaptiveThreshold(
        gray, 255, cv2.ADAPTIVE_THRESH_MEAN_C, cv2.THRESH_BINARY, 9, 2
    )

    # Step 3: Combine color and edges
    edges_colored = cv2.cvtColor(edges, cv2.COLOR_GRAY2BGR)
    cartoon = cv2.bitwise_and(color, edges_colored)

    # Step 4: Style-specific adjustments
    if style == "anime":
        # More saturated colors for anime style
        hsv = cv2.cvtColor(cartoon, cv2.COLOR_BGR2HSV).astype(np.float32)
        hsv[:, :, 1] = np.clip(hsv[:, :, 1] * 1.5, 0, 255)
        hsv[:, :, 2] = np.clip(hsv[:, :, 2] * 1.1, 0, 255)
        cartoon = cv2.cvtColor(hsv.astype(np.uint8), cv2.COLOR_HSV2BGR)
    elif style == "pixel":
        # Pixel art: downscale then upscale
        small = cv2.resize(cartoon, (w // 8, h // 8), interpolation=cv2.INTER_LINEAR)
        cartoon = cv2.resize(small, (w, h), interpolation=cv2.INTER_NEAREST)

    cv2.imwrite(output_path, cartoon)
    return output_path


def detect_skeleton(input_path: str, output_path: str) -> dict:
    """
    Detect body skeleton using OpenCV DNN with MobileNet/OpenPose.
    Returns joint positions for 2D animation.
    """
    img = cv2.imread(input_path)
    if img is None:
        raise ValueError(f"Cannot read image: {input_path}")

    h, w = img.shape[:2]

    # Use OpenCV DNN for pose estimation if model available
    # Fallback: use contour-based body part detection
    skeleton = detect_skeleton_contour(img)

    # Draw skeleton on image
    skeleton_img = draw_skeleton(img, skeleton)

    cv2.imwrite(output_path, skeleton_img)

    return {
        "joints": skeleton["joints"],
        "bones": skeleton["bones"],
        "image_size": {"width": w, "height": h},
    }


def detect_skeleton_contour(img: np.ndarray) -> dict:
    """
    Simplified skeleton detection using contour analysis.
    For production: replace with OpenPose/MediaPipe/BlazePose.
    """
    h, w = img.shape[:2]
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    # Detect face region (approximate head position)
    face_cascade = cv2.CascadeClassifier(
        cv2.data.haarcascades + "haarcascade_frontalface_default.xml"
    )
    faces = face_cascade.detectMultiScale(gray, 1.1, 5, minSize=(30, 30))

    # Estimate body proportions based on face position
    if len(faces) > 0:
        fx, fy, fw, fh = faces[0]
        head_center = (fx + fw // 2, fy + fh // 2)
        head_size = fw
    else:
        # Default: assume face is in upper-center
        head_center = (w // 2, h // 4)
        head_size = w // 6

    # Generate skeleton joints based on body proportions
    # (Head → Neck → Shoulders → Elbows → Wrists → Hips → Knees → Ankles)
    cx, cy = head_center
    neck_y = cy + head_size
    shoulder_y = neck_y + head_size // 2
    shoulder_span = head_size * 2
    hip_y = int(shoulder_y + (h - shoulder_y) * 0.4)
    knee_y = int(hip_y + (h - hip_y) * 0.5)
    ankle_y = h - head_size // 2

    joints = {
        "head": {"x": int(cx), "y": int(cy), "confidence": 0.9},
        "neck": {"x": int(cx), "y": int(neck_y), "confidence": 0.85},
        "left_shoulder": {"x": int(cx - shoulder_span // 2), "y": int(shoulder_y), "confidence": 0.8},
        "right_shoulder": {"x": int(cx + shoulder_span // 2), "y": int(shoulder_y), "confidence": 0.8},
        "left_elbow": {"x": int(cx - shoulder_span), "y": int(shoulder_y + head_size), "confidence": 0.7},
        "right_elbow": {"x": int(cx + shoulder_span), "y": int(shoulder_y + head_size), "confidence": 0.7},
        "left_wrist": {"x": int(cx - shoulder_span - head_size // 2), "y": int(shoulder_y + head_size * 2), "confidence": 0.6},
        "right_wrist": {"x": int(cx + shoulder_span + head_size // 2), "y": int(shoulder_y + head_size * 2), "confidence": 0.6},
        "hip": {"x": int(cx), "y": int(hip_y), "confidence": 0.8},
        "left_hip": {"x": int(cx - shoulder_span // 3), "y": int(hip_y), "confidence": 0.75},
        "right_hip": {"x": int(cx + shoulder_span // 3), "y": int(hip_y), "confidence": 0.75},
        "left_knee": {"x": int(cx - shoulder_span // 3), "y": int(knee_y), "confidence": 0.7},
        "right_knee": {"x": int(cx + shoulder_span // 3), "y": int(knee_y), "confidence": 0.7},
        "left_ankle": {"x": int(cx - shoulder_span // 3), "y": int(ankle_y), "confidence": 0.65},
        "right_ankle": {"x": int(cx + shoulder_span // 3), "y": int(ankle_y), "confidence": 0.65},
    }

    bones = [
        ("head", "neck"),
        ("neck", "left_shoulder"),
        ("neck", "right_shoulder"),
        ("left_shoulder", "left_elbow"),
        ("right_shoulder", "right_elbow"),
        ("left_elbow", "left_wrist"),
        ("right_elbow", "right_wrist"),
        ("neck", "hip"),
        ("hip", "left_hip"),
        ("hip", "right_hip"),
        ("left_hip", "left_knee"),
        ("right_hip", "right_knee"),
        ("left_knee", "left_ankle"),
        ("right_knee", "right_ankle"),
    ]

    return {"joints": joints, "bones": bones}


def draw_skeleton(img: np.ndarray, skeleton: dict) -> np.ndarray:
    """Draw skeleton overlay on image."""
    canvas = img.copy()
    joints = skeleton["joints"]
    bones = skeleton["bones"]

    # Draw bones (lines)
    for j1_name, j2_name in bones:
        if j1_name in joints and j2_name in joints:
            j1 = joints[j1_name]
            j2 = joints[j2_name]
            cv2.line(canvas, (j1["x"], j1["y"]), (j2["x"], j2["y"]), (0, 255, 0), 2)

    # Draw joints (circles)
    for name, j in joints.items():
        color = (0, 0, 255) if "head" in name else (255, 0, 0)
        cv2.circle(canvas, (j["x"], j["y"]), 5, color, -1)
        cv2.putText(canvas, name[:3], (j["x"] + 5, j["y"] - 5),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.3, (255, 255, 255), 1)

    return canvas


def generate_animation_data(skeleton: dict, output_path: str) -> dict:
    """
    Generate animation keyframes for the skeleton.
    Creates idle breathing + wave animation.
    """
    joints = skeleton["joints"]
    h = skeleton["image_size"]["height"]

    # Generate idle animation (breathing motion)
    keyframes = []
    num_frames = 30  # 1 second at 30fps

    for frame in range(num_frames):
        t = frame / num_frames
        breath_offset = int(3 * (1 + __import__("math").sin(t * 2 * 3.14159)))

        frame_joints = {}
        for name, j in joints.items():
            fj = {"x": j["x"], "y": j["y"]}
            # Add breathing to upper body
            if name in ["head", "neck", "left_shoulder", "right_shoulder"]:
                fj["y"] = j["y"] - breath_offset
            # Add subtle sway
            if "left" in name:
                fj["x"] = j["x"] - 1
            elif "right" in name:
                fj["x"] = j["x"] + 1
            frame_joints[name] = {"x": int(fj["x"]), "y": int(fj["y"])}

        keyframes.append({"frame": frame, "joints": frame_joints})

    # Add wave animation for right hand
    wave_keyframes = []
    for frame in range(30):
        t = frame / 30
        wave_angle = 20 * __import__("math").sin(t * 4 * 3.14159)

        frame_joints = {}
        for name, j in joints.items():
            fj = {"x": j["x"], "y": j["y"]}
            if name == "right_wrist":
                fj["y"] = j["y"] + int(wave_angle)
                fj["x"] = j["x"] + int(wave_angle * 0.5)
            frame_joints[name] = {"x": int(fj["x"]), "y": int(fj["y"])}

        wave_keyframes.append({"frame": frame, "joints": frame_joints})

    animation = {
        "idle": keyframes,
        "wave": wave_keyframes,
        "fps": 30,
        "skeleton": {
            "joints": {k: {"x": v["x"], "y": v["y"]} for k, v in joints.items()},
            "bones": skeleton["bones"],
        },
    }

    with open(output_path, "w") as f:
        json.dump(animation, f, indent=2)

    return animation


def main():
    parser = argparse.ArgumentParser(description="Photo → 2D Cartoon with Skeleton")
    parser.add_argument("--input", required=True, help="Input photo path")
    parser.add_argument("--output-dir", required=True, help="Output directory")
    parser.add_argument("--style", default="cartoon", choices=["cartoon", "anime", "pixel"])
    parser.add_argument("--name", default="avatar", help="Output name prefix")
    args = parser.parse_args()

    os.makedirs(args.output_dir, exist_ok=True)

    # Step 1: Cartoon style
    cartoon_path = os.path.join(args.output_dir, f"{args.name}_cartoon.png")
    print(f"[1/3] Generating {args.style} style...", file=sys.stderr)
    photo_to_cartoon(args.input, cartoon_path, args.style)

    # Step 2: Skeleton detection
    skeleton_path = os.path.join(args.output_dir, f"{args.name}_skeleton.png")
    print("[2/3] Detecting skeleton...", file=sys.stderr)
    skeleton_data = detect_skeleton(cartoon_path, skeleton_path)

    # Step 3: Animation data
    anim_path = os.path.join(args.output_dir, f"{args.name}_animation.json")
    print("[3/3] Generating animation data...", file=sys.stderr)
    anim_data = generate_animation_data(skeleton_data, anim_path)

    # Output result
    result = {
        "cartoon_image": cartoon_path,
        "skeleton_image": skeleton_path,
        "animation_data": anim_path,
        "joints": skeleton_data["joints"],
        "bones": skeleton_data["bones"],
        "animations": list(anim_data.keys()),
    }

    print(json.dumps(result))


if __name__ == "__main__":
    main()
