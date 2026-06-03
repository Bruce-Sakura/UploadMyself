#!/usr/bin/env python3
"""Audio preprocessing CLI: convert, denoise, VAD-segment, emit JSON manifest."""

import argparse
import json
import logging
import os
import subprocess
import sys
import tempfile
from pathlib import Path

import numpy as np
import soundfile as sf

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)

SAMPLE_RATE = 22050
# Energy-based VAD parameters
FRAME_MS = 30          # frame length in ms
ENERGY_THRESHOLD = 0.01  # RMS threshold for voiced frame
MIN_SEGMENT_MS = 500     # minimum segment duration
MERGE_GAP_MS = 300       # merge segments closer than this


def convert_to_wav(input_path: str, output_path: str) -> None:
    """Convert any audio to WAV 22050 Hz mono via ffmpeg."""
    cmd = [
        "ffmpeg", "-y", "-i", input_path,
        "-ar", str(SAMPLE_RATE),
        "-ac", "1",
        "-sample_fmt", "s16",
        output_path,
    ]
    logger.info("Converting %s -> %s", input_path, output_path)
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise RuntimeError(f"ffmpeg failed: {result.stderr}")


def reduce_noise(audio: np.ndarray, sr: int) -> np.ndarray:
    """Apply noise reduction using noisereduce library."""
    try:
        import noisereduce as nr
        logger.info("Applying noise reduction")
        return nr.reduce_noise(y=audio, sr=sr, prop_decrease=0.8)
    except ImportError:
        logger.warning("noisereduce not installed, skipping noise reduction")
        return audio


def vad_segment(audio: np.ndarray, sr: int) -> list[dict]:
    """Energy-based voice activity detection. Returns list of {start, end} in seconds."""
    frame_len = int(sr * FRAME_MS / 1000)
    hop = frame_len
    n_frames = max(1, (len(audio) - frame_len) // hop + 1)

    # Compute per-frame RMS energy
    energies = []
    for i in range(n_frames):
        start = i * hop
        frame = audio[start : start + frame_len]
        energies.append(float(np.sqrt(np.mean(frame ** 2))))

    # Label frames as voiced/unvoiced
    voiced = [e >= ENERGY_THRESHOLD for e in energies]

    # Group contiguous voiced frames into segments
    segments = []
    in_seg = False
    seg_start = 0
    for i, v in enumerate(voiced):
        if v and not in_seg:
            seg_start = i
            in_seg = True
        elif not v and in_seg:
            segments.append((seg_start, i))
            in_seg = False
    if in_seg:
        segments.append((seg_start, len(voiced)))

    # Convert frame indices to seconds
    segs_sec = []
    for s, e in segments:
        start_s = s * FRAME_MS / 1000
        end_s = min(e * FRAME_MS / 1000, len(audio) / sr)
        dur = end_s - start_s
        if dur * 1000 >= MIN_SEGMENT_MS:
            segs_sec.append({"start": round(start_s, 4), "end": round(end_s, 4)})

    # Merge segments separated by small gaps
    merged = []
    for seg in segs_sec:
        if merged and (seg["start"] - merged[-1]["end"]) * 1000 < MERGE_GAP_MS:
            merged[-1]["end"] = seg["end"]
        else:
            merged.append(seg)

    logger.info("VAD found %d segments", len(merged))
    return merged


def write_segments(audio: np.ndarray, sr: int, segments: list[dict], out_dir: str) -> list[dict]:
    """Write each segment to a WAV file and return manifest entries."""
    entries = []
    for i, seg in enumerate(segments):
        start_sample = int(seg["start"] * sr)
        end_sample = int(seg["end"] * sr)
        chunk = audio[start_sample:end_sample]
        filename = f"segment_{i:04d}.wav"
        path = os.path.join(out_dir, filename)
        sf.write(path, chunk, sr)
        entries.append({
            "path": path,
            "start": seg["start"],
            "end": seg["end"],
            "duration": round(seg["end"] - seg["start"], 4),
        })
    return entries


def main():
    parser = argparse.ArgumentParser(description="Preprocess audio: convert, denoise, segment")
    parser.add_argument("--input", required=True, help="Input audio file path")
    parser.add_argument("--output", required=True, help="Output directory for segments and manifest")
    args = parser.parse_args()

    input_path = args.input
    output_dir = args.output

    if not os.path.isfile(input_path):
        logger.error("Input file not found: %s", input_path)
        sys.exit(1)

    os.makedirs(output_dir, exist_ok=True)

    # Step 1: convert to WAV 22050 mono
    converted_path = os.path.join(output_dir, "_converted.wav")
    try:
        convert_to_wav(input_path, converted_path)
    except Exception as e:
        logger.error("Conversion failed: %s", e)
        sys.exit(1)

    # Step 2: load audio
    try:
        audio, sr = sf.read(converted_path, dtype="float32")
    except Exception as e:
        logger.error("Failed to read converted audio: %s", e)
        sys.exit(1)

    # Step 3: noise reduction
    audio = reduce_noise(audio, sr)

    # Step 4: VAD segmentation
    segments = vad_segment(audio, sr)

    if not segments:
        logger.warning("No voiced segments found; using full audio as single segment")
        segments = [{"start": 0.0, "end": round(len(audio) / sr, 4)}]

    # Step 5: write segment WAV files
    entries = write_segments(audio, sr, segments, output_dir)

    # Step 6: write manifest
    manifest = {
        "source": os.path.abspath(input_path),
        "sample_rate": sr,
        "total_duration": round(len(audio) / sr, 4),
        "segment_count": len(entries),
        "segments": entries,
    }
    manifest_path = os.path.join(output_dir, "manifest.json")
    with open(manifest_path, "w") as f:
        json.dump(manifest, f, indent=2)

    # Clean up temp converted file
    try:
        os.remove(converted_path)
    except OSError:
        pass

    logger.info("Done. Manifest: %s (%d segments)", manifest_path, len(entries))
    print(json.dumps(manifest, indent=2))


if __name__ == "__main__":
    main()
