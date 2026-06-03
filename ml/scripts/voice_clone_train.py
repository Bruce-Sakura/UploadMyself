#!/usr/bin/env python3
"""Voice clone training wrapper. Currently simulates training and writes a voice profile."""

import argparse
import json
import logging
import os
import sys
import time
from datetime import datetime, timezone
from pathlib import Path

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)


def load_config(config_path: str) -> dict:
    """Load and validate training config JSON."""
    with open(config_path) as f:
        cfg = json.load(f)

    required = ["voice_id", "audio_dir", "model_dir"]
    missing = [k for k in required if k not in cfg]
    if missing:
        raise ValueError(f"Missing required config keys: {missing}")

    return cfg


def collect_audio_files(audio_dir: str) -> list[str]:
    """List audio files in the training directory."""
    exts = {".wav", ".mp3", ".flac", ".ogg", ".m4a", ".webm"}
    files = []
    for p in sorted(Path(audio_dir).iterdir()):
        if p.suffix.lower() in exts and p.is_file():
            files.append(str(p))
    return files


def simulate_training(cfg: dict) -> dict:
    """Simulate training and return a voice profile dict."""
    voice_id = cfg["voice_id"]
    audio_dir = cfg["audio_dir"]
    text = cfg.get("text", "")
    epochs = cfg.get("epochs", 5)

    audio_files = collect_audio_files(audio_dir)
    logger.info("Voice %s: found %d audio files in %s", voice_id, len(audio_files), audio_dir)

    # Simulate training time
    logger.info("Simulating training for %d epochs...", epochs)
    time.sleep(0.5)  # brief pause to simulate work

    total_duration = 0.0
    file_details = []
    for af in audio_files:
        try:
            import soundfile as sf
            info = sf.info(af)
            dur = info.duration
        except Exception:
            dur = 0.0
        total_duration += dur
        file_details.append({"path": af, "duration": round(dur, 2)})

    profile = {
        "voice_id": voice_id,
        "created_at": datetime.now(timezone.utc).isoformat(),
        "status": "ready",
        "training": {
            "epochs": epochs,
            "audio_files": file_details,
            "total_audio_duration": round(total_duration, 2),
            "reference_text": text,
        },
        "model": {
            "type": "voice_clone",
            "version": "0.1.0-simulated",
            "embedding_path": os.path.join(cfg["model_dir"], voice_id, "embedding.npy"),
        },
        "metadata": {
            "sample_count": len(audio_files),
            "engine": "simulated",
        },
    }
    return profile


def main():
    parser = argparse.ArgumentParser(description="Voice clone training wrapper")
    parser.add_argument("--config", required=True, help="Path to training config JSON")
    args = parser.parse_args()

    if not os.path.isfile(args.config):
        logger.error("Config file not found: %s", args.config)
        sys.exit(1)

    try:
        cfg = load_config(args.config)
    except Exception as e:
        logger.error("Failed to load config: %s", e)
        sys.exit(1)

    voice_id = cfg["voice_id"]
    model_dir = cfg["model_dir"]
    profile_dir = os.path.join(model_dir, voice_id)
    os.makedirs(profile_dir, exist_ok=True)

    try:
        profile = simulate_training(cfg)
    except Exception as e:
        logger.error("Training failed: %s", e)
        sys.exit(1)

    # Write profile
    profile_path = os.path.join(profile_dir, "profile.json")
    with open(profile_path, "w") as f:
        json.dump(profile, f, indent=2)

    logger.info("Voice profile saved: %s", profile_path)
    print(json.dumps(profile, indent=2))


if __name__ == "__main__":
    main()
