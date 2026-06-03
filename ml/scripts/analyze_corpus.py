#!/usr/bin/env python3
"""Corpus analysis for skill cloning: chunk text, compute stats, extract keywords."""

import argparse
import json
import logging
import os
import re
import sys
from collections import Counter
from pathlib import Path

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)

# Common English stop words to filter from keywords
STOP_WORDS = frozenset(
    "a an the and or but is are was were be been being have has had do does did "
    "will would shall should may might can could to of in for on with at by from "
    "as into through during before after above below between out off over under "
    "again further then once here there when where why how all each every both "
    "few more most other some such no nor not only own same so than too very "
    "it its he she they them their his her my your our i me we us you".split()
)


def chunk_paragraphs(text: str) -> list[str]:
    """Split text into paragraphs (double-newline separated), filtering empty."""
    chunks = re.split(r"\n\s*\n", text.strip())
    return [c.strip() for c in chunks if c.strip()]


def sentences(text: str) -> list[str]:
    """Split text into sentences."""
    sents = re.split(r"[.!?]+", text)
    return [s.strip() for s in sents if s.strip()]


def extract_keywords(text: str, top_n: int = 20) -> list[dict]:
    """Extract top keywords by frequency, filtering stop words."""
    words = re.findall(r"[a-zA-Z]{3,}", text.lower())
    filtered = [w for w in words if w not in STOP_WORDS]
    counter = Counter(filtered)
    return [{"word": word, "count": count} for word, count in counter.most_common(top_n)]


def analyze(text: str) -> dict:
    """Full corpus analysis."""
    paragraphs = chunk_paragraphs(text)
    all_sentences = sentences(text)
    words = re.findall(r"[a-zA-Z]+", text)
    unique = set(w.lower() for w in words)

    avg_sentence_len = 0.0
    if all_sentences:
        sent_lengths = [len(s.split()) for s in all_sentences]
        avg_sentence_len = round(sum(sent_lengths) / len(sent_lengths), 2)

    return {
        "chunks": paragraphs,
        "chunk_count": len(paragraphs),
        "word_count": len(words),
        "unique_words": len(unique),
        "sentence_count": len(all_sentences),
        "avg_sentence_length": avg_sentence_len,
        "top_keywords": extract_keywords(text),
    }


def main():
    parser = argparse.ArgumentParser(description="Analyze text corpus for skill cloning")
    parser.add_argument("--input", required=True, help="Path to corpus text file")
    parser.add_argument("--output", required=True, help="Output JSON path")
    parser.add_argument("--top-n", type=int, default=20, help="Number of top keywords (default: 20)")
    args = parser.parse_args()

    if not os.path.isfile(args.input):
        logger.error("Input file not found: %s", args.input)
        sys.exit(1)

    try:
        text = Path(args.input).read_text(encoding="utf-8")
    except Exception as e:
        logger.error("Failed to read input: %s", e)
        sys.exit(1)

    if not text.strip():
        logger.error("Input file is empty")
        sys.exit(1)

    result = analyze(text)

    os.makedirs(os.path.dirname(os.path.abspath(args.output)), exist_ok=True)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(result, f, indent=2, ensure_ascii=False)

    logger.info("Analysis complete: %d words, %d unique, %d chunks -> %s",
                result["word_count"], result["unique_words"], result["chunk_count"], args.output)
    print(json.dumps(result, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
