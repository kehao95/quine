#!/usr/bin/env python3
"""
Batch run baseline experiments on MRCR samples.

Usage:
    python run_baseline_batch.py --model moonshotai/kimi-k2.5 --provider openrouter
"""

import argparse
import json
import sys
from pathlib import Path
from datetime import datetime

# Import baseline runner
sys.path.insert(0, str(Path(__file__).parent))
from baseline import run_baseline

# Samples to test
SAMPLES = [
    # ~4K tokens (4-needle)
    "4needle_0000",
    "4needle_0002",
    # ~7K tokens (2-needle)  
    "2needle_0400",
    "2needle_0401",
    "2needle_0402",
    # ~178K-278K tokens (8-needle)
    "8needle_0000",
    "8needle_0002",
    "8needle_0004",
]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--model", default="moonshotai/kimi-k2.5")
    parser.add_argument("--provider", default="openrouter")
    parser.add_argument("--samples", nargs="+", default=SAMPLES)
    parser.add_argument("--output-dir", default=None)
    args = parser.parse_args()
    
    script_dir = Path(__file__).parent.parent  # Go up from eval/ to experiment dir
    data_dir = script_dir / "data"
    
    # Output directory
    if args.output_dir:
        output_dir = Path(args.output_dir)
    else:
        timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
        model_short = args.model.split("/")[-1]
        output_dir = script_dir / "baseline_results" / f"{timestamp}-{model_short}"
    
    output_dir.mkdir(parents=True, exist_ok=True)
    
    print(f"=" * 50)
    print(f"MRCR Baseline Batch Run")
    print(f"Model: {args.model}")
    print(f"Provider: {args.provider}")
    print(f"Samples: {len(args.samples)}")
    print(f"Output: {output_dir}")
    print(f"=" * 50)
    
    results = []
    
    for sample_id in args.samples:
        sample_dir = data_dir / sample_id
        
        if not sample_dir.exists():
            print(f"\n[SKIP] {sample_id} - not found")
            continue
        
        # Get token count
        meta = json.loads((sample_dir / "meta.json").read_text())
        tokens = meta.get("estimated_tokens", "?")
        
        print(f"\n[{sample_id}] ~{tokens} tokens ...")
        
        try:
            result = run_baseline(sample_dir, args.model, args.provider)
            results.append(result)
            
            # Save individual result
            result_file = output_dir / f"{sample_id}.json"
            result_file.write_text(json.dumps(result, indent=2, ensure_ascii=False))
            
            print(f"    Score: {result['score']:.3f}")
            print(f"    Time: {result['elapsed_seconds']:.1f}s")
            print(f"    Tokens: {result.get('total_tokens', '?')}")
            
        except Exception as e:
            print(f"    ERROR: {e}")
            results.append({
                "sample_id": sample_id,
                "error": str(e),
                "score": None,
            })
    
    # Summary
    print(f"\n{'=' * 50}")
    print("SUMMARY")
    print(f"{'=' * 50}")
    
    successful = [r for r in results if r.get("score") is not None]
    if successful:
        avg_score = sum(r["score"] for r in successful) / len(successful)
        print(f"Successful: {len(successful)}/{len(results)}")
        print(f"Average score: {avg_score:.3f}")
        
        print(f"\nBy sample:")
        for r in results:
            score = r.get("score", "ERROR")
            if isinstance(score, float):
                score = f"{score:.3f}"
            print(f"  {r['sample_id']}: {score}")
    
    # Save summary
    summary = {
        "model": args.model,
        "provider": args.provider,
        "timestamp": datetime.now().isoformat(),
        "results": results,
    }
    (output_dir / "summary.json").write_text(json.dumps(summary, indent=2, ensure_ascii=False))
    print(f"\nResults saved to: {output_dir}")


if __name__ == "__main__":
    main()
