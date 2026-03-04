#!/usr/bin/env python3
"""
Baseline: Pure LLM evaluation on MRCR task.

Sends the full conversation as messages to the API and grades the response.

Usage:
    python baseline.py --sample-dir ./data/2needle_0000 --model gpt-4o
    python baseline.py --sample-dir ./data/8needle_0000 --model moonshotai/kimi-k2.5 --provider openrouter
"""

import argparse
import json
import os
import time
from pathlib import Path
from difflib import SequenceMatcher

try:
    from openai import OpenAI
except ImportError:
    print("Please install: pip install openai")
    exit(1)


# Provider configurations
PROVIDERS = {
    "openai": {
        "base_url": None,  # Use default
        "api_key_env": "OPENAI_API_KEY",
    },
    "openrouter": {
        "base_url": "https://openrouter.ai/api/v1",
        "api_key_env": "OPENROUTER_API_KEY",
    },
}


def load_sample(sample_dir: Path) -> tuple[str, dict]:
    """Load conversation and metadata from sample directory."""
    conversation = (sample_dir / "conversation.txt").read_text()
    meta = json.loads((sample_dir / "meta.json").read_text())
    return conversation, meta


def grade(response: str, answer: str, hash_prefix: str) -> float:
    """
    Grade response against expected answer.
    
    Returns SequenceMatcher ratio if hash prefix is correct, 0 otherwise.
    """
    # Strip leading/trailing whitespace
    response = response.strip()
    
    if not response.startswith(hash_prefix):
        return 0.0
    
    response_stripped = response.removeprefix(hash_prefix).strip()
    answer_stripped = answer.removeprefix(hash_prefix).strip()
    
    return SequenceMatcher(None, response_stripped, answer_stripped).ratio()


def run_baseline(sample_dir: Path, model: str, provider: str = "openai") -> dict:
    """Run baseline LLM on a sample."""
    
    conversation, meta = load_sample(sample_dir)
    task = meta["task"]
    
    # Reconstruct messages array for API call
    # Parse the conversation text back into messages
    messages = []
    current_role = None
    current_content = []
    
    for line in conversation.split("\n"):
        if line == "[USER]":
            if current_role and current_content:
                messages.append({
                    "role": current_role,
                    "content": "\n".join(current_content).strip()
                })
            current_role = "user"
            current_content = []
        elif line == "[ASSISTANT]":
            if current_role and current_content:
                messages.append({
                    "role": current_role,
                    "content": "\n".join(current_content).strip()
                })
            current_role = "assistant"
            current_content = []
        else:
            current_content.append(line)
    
    # Don't forget the last message
    if current_role and current_content:
        messages.append({
            "role": current_role,
            "content": "\n".join(current_content).strip()
        })
    
    # Add the final question
    messages.append({
        "role": "user",
        "content": task["original_question"]
    })
    
    # Set up client based on provider
    provider_config = PROVIDERS.get(provider, PROVIDERS["openai"])
    client_kwargs = {}
    
    if provider_config["base_url"]:
        client_kwargs["base_url"] = provider_config["base_url"]
    
    api_key = os.environ.get(provider_config["api_key_env"])
    if api_key:
        client_kwargs["api_key"] = api_key
    
    client = OpenAI(**client_kwargs)
    
    start_time = time.time()
    completion = client.chat.completions.create(
        model=model,
        messages=messages,
    )
    elapsed = time.time() - start_time
    
    response = completion.choices[0].message.content
    
    # Grade
    score = grade(response, meta["answer"], meta["random_string_to_prepend"])
    
    # Handle token usage (may be None for some providers)
    usage = completion.usage
    input_tokens = usage.prompt_tokens if usage else None
    output_tokens = usage.completion_tokens if usage else None
    total_tokens = usage.total_tokens if usage else None
    
    result = {
        "sample_id": meta["sample_id"],
        "model": model,
        "provider": provider,
        "method": "baseline",
        "score": score,
        "elapsed_seconds": elapsed,
        "input_tokens": input_tokens,
        "output_tokens": output_tokens,
        "total_tokens": total_tokens,
        "response": response,  # Full response for accurate grading
        "response_preview": response[:500],  # Preview for logging
        "expected_prefix": meta["random_string_to_prepend"],
    }
    
    return result


def main():
    parser = argparse.ArgumentParser(description="Run baseline LLM on MRCR sample")
    parser.add_argument("--sample-dir", type=str, required=True,
                        help="Path to sample directory")
    parser.add_argument("--model", type=str, default="gpt-4o",
                        help="Model to use")
    parser.add_argument("--provider", type=str, default="openai",
                        choices=["openai", "openrouter"],
                        help="API provider")
    parser.add_argument("--output", type=str, default=None,
                        help="Output JSON file")
    
    args = parser.parse_args()
    
    sample_dir = Path(args.sample_dir)
    result = run_baseline(sample_dir, args.model, args.provider)
    
    print(f"Sample: {result['sample_id']}")
    print(f"Score: {result['score']:.3f}")
    print(f"Tokens: {result['total_tokens']}")
    print(f"Time: {result['elapsed_seconds']:.2f}s")
    
    if args.output:
        Path(args.output).write_text(json.dumps(result, indent=2))
        print(f"Saved to {args.output}")
    
    return result
    
    return result


if __name__ == "__main__":
    main()
