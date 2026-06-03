#!/usr/bin/env python3
"""Face detection utility using OpenCV Haar cascades with blur quality scoring."""

import argparse
import json
import logging
import os
import sys

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)


def detect_faces(image_path: str) -> dict:
    """Detect faces and compute quality metrics."""
    try:
        import cv2
    except ImportError:
        logger.error("opencv-python-headless not installed")
        sys.exit(1)

    if not os.path.isfile(image_path):
        raise FileNotFoundError(f"Image not found: {image_path}")

    img = cv2.imread(image_path)
    if img is None:
        raise ValueError(f"Failed to read image: {image_path}")

    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    h, w = gray.shape[:2]

    # Global blur metric (Laplacian variance)
    laplacian = cv2.Laplacian(gray, cv2.CV_64F)
    global_blur = float(laplacian.var())

    # Load Haar cascade
    cascade_path = cv2.data.haarcascades + "haarcascade_frontalface_default.xml"
    if not os.path.isfile(cascade_path):
        raise FileNotFoundError(f"Haar cascade not found at {cascade_path}")

    face_cascade = cv2.CascadeClassifier(cascade_path)
    detections = face_cascade.detectMultiScale(
        gray,
        scaleFactor=1.1,
        minNeighbors=5,
        minSize=(30, 30),
        flags=cv2.CASCADE_SCALE_IMAGE,
    )

    faces = []
    for x, y, fw, fh in detections:
        # Per-face blur metric
        face_roi = gray[y : y + fh, x : x + fw]
        face_lap = cv2.Laplacian(face_roi, cv2.CV_64F)
        face_blur = float(face_lap.var())

        # Confidence proxy: sharper face = higher quality
        quality = min(1.0, face_blur / 500.0)  # normalize; 500+ is sharp

        faces.append({
            "bbox": {"x": int(x), "y": int(y), "width": int(fw), "height": int(fh)},
            "blur_score": round(face_blur, 2),
            "quality_score": round(quality, 4),
        })

    result = {
        "image_path": os.path.abspath(image_path),
        "image_size": {"width": w, "height": h},
        "face_count": len(faces),
        "global_blur_score": round(global_blur, 2),
        "faces": faces,
    }
    return result


def main():
    parser = argparse.ArgumentParser(description="Detect faces and assess image quality")
    parser.add_argument("--input", required=True, help="Path to input image")
    parser.add_argument("--output", required=True, help="Output JSON path")
    args = parser.parse_args()

    try:
        result = detect_faces(args.input)
    except Exception as e:
        logger.error("Detection failed: %s", e)
        # Output error as JSON so Go can parse it
        error_result = {
            "error": str(e),
            "image_path": os.path.abspath(args.input) if os.path.exists(args.input) else args.input,
        }
        os.makedirs(os.path.dirname(os.path.abspath(args.output)), exist_ok=True)
        with open(args.output, "w") as f:
            json.dump(error_result, f, indent=2)
        print(json.dumps(error_result, indent=2))
        sys.exit(1)

    os.makedirs(os.path.dirname(os.path.abspath(args.output)), exist_ok=True)
    with open(args.output, "w") as f:
        json.dump(result, f, indent=2)

    logger.info("Detected %d face(s) in %s -> %s", result["face_count"], args.input, args.output)
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
