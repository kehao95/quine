#!/usr/bin/env python3
"""
Prepare MRCR dataset samples for Quine experiments.

Downloads samples from HuggingFace and transforms them into:
1. A streaming text file (conversation.txt)
2. A metadata JSON file (meta.json) with answer and grading info

Usage:
    python prepare-data.py --needles 2 --samples 10 --output ./data
"""

import argparse
import json
import os
from pathlib import Path

try:
    from huggingface_hub import hf_hub_download
    import pandas as pd
except ImportError:
    print("Please install dependencies: pip install huggingface_hub pandas pyarrow")
    exit(1)


def download_dataset(needles: int) -> pd.DataFrame:
    """Download MRCR dataset for specified needle count."""
    print(f"Downloading {needles}needle dataset...")
    
    # Download both parquet files and concatenate
    df0 = pd.read_parquet(
        hf_hub_download(
            repo_id="openai/mrcr",
            filename=f"{needles}needle/{needles}needle_0.parquet",
            repo_type="dataset"
        )
    )
    df1 = pd.read_parquet(
        hf_hub_download(
            repo_id="openai/mrcr",
            filename=f"{needles}needle/{needles}needle_1.parquet",
            repo_type="dataset"
        )
    )
    
    return pd.concat([df0, df1], ignore_index=True)


def messages_to_text(messages: list[dict]) -> tuple[str, dict]:
    """
    Convert MRCR message array to streaming text format.
    
    Returns:
        - conversation_text: The conversation without the final question
        - task_info: Dict with nth, format, topic, hash extracted from final question
    """
    # The last user message is the actual question
    # Format: "Prepend {hash} to the {nth} (1 indexed) {format} about {topic}. Do not include any other text in your response."
    
    conversation_messages = messages[:-1]  # All but last
    final_question = messages[-1]["content"]
    
    # Convert conversation to text
    lines = []
    for msg in conversation_messages:
        role = msg["role"].upper()
        content = msg["content"]
        lines.append(f"[{role}]")
        lines.append(content)
        lines.append("")  # Blank line separator
    
    conversation_text = "\n".join(lines)
    
    # Parse the final question to extract task parameters
    # Example: "Prepend aYooSG8CQg to the 2nd (1 indexed) poem about tapirs. Do not include any other text in your response."
    import re
    
    # Extract hash (alphanumeric string after "Prepend")
    hash_match = re.search(r'Prepend\s+(\w+)\s+to', final_question)
    hash_value = hash_match.group(1) if hash_match else ""
    
    # Extract ordinal (1st, 2nd, 3rd, etc.)
    ordinal_match = re.search(r'the\s+(\d+)(?:st|nd|rd|th)', final_question)
    nth = int(ordinal_match.group(1)) if ordinal_match else 1
    
    # Extract format and topic
    # Pattern: "{format} about {topic}"
    format_topic_match = re.search(r'(\w+(?:\s+\w+)*?)\s+about\s+(\w+(?:\s+\w+)*?)\.', final_question)
    if format_topic_match:
        format_type = format_topic_match.group(1)
        topic = format_topic_match.group(2)
    else:
        format_type = "content"
        topic = "unknown"
    
    task_info = {
        "hash": hash_value,
        "nth": nth,
        "format": format_type,
        "topic": topic,
        "original_question": final_question
    }
    
    return conversation_text, task_info


def estimate_tokens(text: str) -> int:
    """Rough token estimate (chars / 4)."""
    return len(text) // 4


def prepare_samples(
    needles: int,
    num_samples: int,
    output_dir: Path,
    max_tokens: int = None
):
    """Prepare samples for the experiment."""
    
    df = download_dataset(needles)
    
    # Filter by token count if specified
    if max_tokens:
        # Estimate tokens from prompt length
        df = df[df["prompt"].apply(lambda x: estimate_tokens(x)) <= max_tokens]
        print(f"Filtered to {len(df)} samples under {max_tokens} tokens")
    
    # Take requested number of samples
    samples = df.head(num_samples)
    
    output_dir.mkdir(parents=True, exist_ok=True)
    
    prepared = []
    for idx, row in samples.iterrows():
        messages = json.loads(row["prompt"])
        conversation_text, task_info = messages_to_text(messages)
        
        sample_id = f"{needles}needle_{idx:04d}"
        sample_dir = output_dir / sample_id
        sample_dir.mkdir(exist_ok=True)
        
        # Write conversation text
        conv_file = sample_dir / "conversation.txt"
        conv_file.write_text(conversation_text)
        
        # Write metadata
        meta = {
            "sample_id": sample_id,
            "needles": needles,
            "answer": row["answer"],
            "random_string_to_prepend": row["random_string_to_prepend"],
            "task": task_info,
            "estimated_tokens": estimate_tokens(conversation_text),
        }
        meta_file = sample_dir / "meta.json"
        meta_file.write_text(json.dumps(meta, indent=2))
        
        prepared.append(meta)
        print(f"Prepared {sample_id}: ~{meta['estimated_tokens']} tokens")
    
    # Write index
    index_file = output_dir / "index.json"
    index_file.write_text(json.dumps(prepared, indent=2))
    
    print(f"\nPrepared {len(prepared)} samples in {output_dir}")
    return prepared


def main():
    parser = argparse.ArgumentParser(description="Prepare MRCR samples for Quine experiments")
    parser.add_argument("--needles", type=int, default=2, choices=[2, 4, 8],
                        help="Number of needles (2, 4, or 8)")
    parser.add_argument("--samples", type=int, default=10,
                        help="Number of samples to prepare")
    parser.add_argument("--max-tokens", type=int, default=None,
                        help="Maximum token count filter")
    parser.add_argument("--output", type=str, default="./data",
                        help="Output directory")
    
    args = parser.parse_args()
    
    prepare_samples(
        needles=args.needles,
        num_samples=args.samples,
        output_dir=Path(args.output),
        max_tokens=args.max_tokens
    )


if __name__ == "__main__":
    main()
