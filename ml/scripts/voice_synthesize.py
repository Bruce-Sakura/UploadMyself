#!/usr/bin/env python3
"""Voice synthesis wrapper. Uses edge-tts if available, otherwise generates a placeholder WAV."""

import argparse
import json
import logging
import math
import os
import struct
import subprocess
import sys
import wave
from pathlib import Path

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)

SAMPLE_RATE = 22050


def load_voice_profile(model_dir: str, voice_id: str) -> dict | None:
    """Load voice profile if it exists."""
    profile_path = os.path.join(model_dir, voice_id, "profile.json")
    if os.path.isfile(profile_path):
        with open(profile_path) as f:
            return json.load(f)
    return None


def synthesize_with_edge_tts(text: str, output_path: str, voice: str = "en-US-AriaNeural") -> bool:
    """Try synthesis with edge-tts CLI. Returns True on success."""
    try:
        cmd = ["edge-tts", "--voice", voice, "--text", text, "--write-media", output_path]
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
        if result.returncode == 0 and os.path.isfile(output_path):
            logger.info("Synthesized with edge-tts")
            return True
        logger.warning("edge-tts failed: %s", result.stderr)
    except FileNotFoundError:
        logger.info("edge-tts not installed")
    except Exception as e:
        logger.warning("edge-tts error: %s", e)
    return False


def generate_placeholder_wav(text: str, output_path: str, duration: float = 2.0) -> None:
    """Generate a simple sine wave WAV as a placeholder."""
    n_samples = int(SAMPLE_RATE * duration)
    freq = 440.0  # A4

    with wave.open(output_path, "w") as wf:
        wf.setnchannels(1)
        wf.setsampwidth(2)  # 16-bit
        wf.setframerate(SAMPLE_RATE)

        frames = bytearray()
        for i in range(n_samples):
            t = i / SAMPLE_RATE
            # Simple envelope: fade in/out
            env = min(t * 10, 1.0) * min((duration - t) * 10, 1.0)
            value = int(32767 * 0.5 * env * math.sin(2 * math.pi * freq * t))
            frames.extend(struct.pack("<h", max(-32768, min(32767, value))))

        wf.writeframes(bytes(frames))

    logger.info("Generated placeholder WAV: %s (%.1fs)", output_path, duration)


def estimate_duration(text: str) -> float:
    """Rough estimate of speech duration from text (~150 words/min)."""
    words = len(text.split())
    return max(1.0, words / 2.5)  # ~150 wpm


def main():
    parser = argparse.ArgumentParser(description="Voice synthesis wrapper")
    parser.add_argument("--voice-id", required=True, help="Voice ID to use")
    parser.add_argument("--text", required=True, help="Text to synthesize")
    parser.add_argument("--output", required=True, help="Output audio file path")
    parser.add_argument("--model-dir", required=True, help="Directory containing voice models")
    parser.add_argument("--engine", default="edge-tts", choices=["edge-tts", "placeholder"],
                        help="Synthesis engine (default: edge-tts with placeholder fallback)")
    args = parser.parse_args()

    os.makedirs(os.path.dirname(os.path.abspath(args.output)), exist_ok=True)

    profile = load_voice_profile(args.model_dir, args.voice_id)
    if profile:
        logger.info("Loaded voice profile for %s (status=%s)", args.voice_id, profile.get("status"))
    else:
        logger.warning("No voice profile found for %s, using defaults", args.voice_id)

    # Try engines in order
    success = False
    if args.engine == "edge-tts":
        success = synthesize_with_edge_tts(args.text, args.output)

    if not success:
        duration = estimate_duration(args.text)
        generate_placeholder_wav(args.text, args.output, duration=duration)

    # Verify output exists
    if not os.path.isfile(args.output):
        logger.error("Output file was not created: %s", args.output)
        sys.exit(1)

    result = {
        "voice_id": args.voice_id,
        "text": args.text,
        "output": os.path.abspath(args.output),
        "engine_used": "edge-tts" if success else "placeholder",
        "duration_estimate": round(estimate_duration(args.text), 2),
    }
    logger.info("Synthesis complete: %s", json.dumps(result))
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
