#!/usr/bin/env python3
"""
Grade Quine output against expected answer.

Usage:
    python grade.py --sample-dir ./data/2needle_0000 --output ./result.txt
"""

import argparse
import json
from pathlib import Path
from difflib import SequenceMatcher


def grade(response: str, answer: str, hash_prefix: str) -> float:
    """
    Grade response against expected answer.
    
    Returns SequenceMatcher ratio if hash prefix is correct, 0 otherwise.
    """
    response = response.strip()
    
    if not response.startswith(hash_prefix):
        return 0.0
    
    response_stripped = response.removeprefix(hash_prefix).strip()
    answer_stripped = answer.removeprefix(hash_prefix).strip()
    
    return SequenceMatcher(None, response_stripped, answer_stripped).ratio()


def main():
    parser = argparse.ArgumentParser(description="Grade Quine output on MRCR sample")
    parser.add_argument("--sample-dir", type=str, required=True,
                        help="Path to sample directory")
    parser.add_argument("--output", type=str, required=True,
                        help="Path to Quine's output file (stdout)")
    parser.add_argument("--result", type=str, default=None,
                        help="Output JSON file for result")
    
    args = parser.parse_args()
    
    sample_dir = Path(args.sample_dir)
    meta = json.loads((sample_dir / "meta.json").read_text())
    
    response = Path(args.output).read_text().strip()
    
    score = grade(response, meta["answer"], meta["random_string_to_prepend"])
    
    result = {
        "sample_id": meta["sample_id"],
        "method": "quine",
        "score": score,
        "response": response[:500],
        "expected_prefix": meta["random_string_to_prepend"],
        "expected_answer_preview": meta["answer"][:200],
    }
    
    print(f"Sample: {result['sample_id']}")
    print(f"Score: {score:.3f}")
    print(f"Expected prefix: {meta['random_string_to_prepend']}")
    print(f"Response starts with: {response[:50]}...")
    
    if args.result:
        Path(args.result).write_text(json.dumps(result, indent=2))
        print(f"Saved to {args.result}")
    
    return result


if __name__ == "__main__":
    main()
