#!/bin/bash
# Batch run MRCR experiments: both Quine and Baseline
#
# Usage: ./run_batch.sh [MODEL]
# Example: ./run_batch.sh claude-sonnet-4-20250514

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MODEL="${1:-claude-sonnet-4-20250514}"

# Samples to run - covering different context lengths
# Excluding mixed/combined tests which have different structure
SAMPLES=(
    # ~4K tokens (4-needle)
    "4needle_0000"
    "4needle_0001"
    "4needle_0002"
    # ~7K tokens (2-needle)  
    "2needle_0400"
    "2needle_0401"
    "2needle_0402"
    # ~178K-278K tokens (8-needle)
    "8needle_0000"
    "8needle_0002"
    "8needle_0003"
    "8needle_0004"
    # Note: 8needle_0001 excluded - known model miscount issue
)

echo "=========================================="
echo "MRCR Batch Experiment"
echo "Model: $MODEL"
echo "Samples: ${#SAMPLES[@]}"
echo "=========================================="

# Results collection
RESULTS_DIR="${SCRIPT_DIR}/batch_results/$(date +%Y%m%d-%H%M%S)"
mkdir -p "$RESULTS_DIR"

for sample in "${SAMPLES[@]}"; do
    SAMPLE_DIR="${SCRIPT_DIR}/data/${sample}"
    
    if [[ ! -d "$SAMPLE_DIR" ]]; then
        echo "SKIP: $sample (not found)"
        continue
    fi
    
    echo ""
    echo ">>> Running: $sample"
    
    # Run Quine
    echo "    [Quine] ..."
    if "${SCRIPT_DIR}/run.sh" "$SAMPLE_DIR" "$MODEL" > "${RESULTS_DIR}/${sample}_quine.log" 2>&1; then
        # Find the latest result
        LATEST_RUN=$(ls -td "${SCRIPT_DIR}/runs/${sample}/"*/ 2>/dev/null | head -1)
        if [[ -f "${LATEST_RUN}/meta/result.json" ]]; then
            cp "${LATEST_RUN}/meta/result.json" "${RESULTS_DIR}/${sample}_quine_result.json"
            SCORE=$(jq -r '.score' "${RESULTS_DIR}/${sample}_quine_result.json")
            echo "    [Quine] Score: $SCORE"
        fi
    else
        echo "    [Quine] FAILED"
    fi
    
    # Run Baseline
    echo "    [Baseline] ..."
    if python3 "${SCRIPT_DIR}/eval/baseline.py" \
        --sample-dir "$SAMPLE_DIR" \
        --model "gpt-4o" \
        --output "${RESULTS_DIR}/${sample}_baseline_result.json" \
        > "${RESULTS_DIR}/${sample}_baseline.log" 2>&1; then
        SCORE=$(jq -r '.score' "${RESULTS_DIR}/${sample}_baseline_result.json")
        echo "    [Baseline] Score: $SCORE"
    else
        echo "    [Baseline] FAILED (check ${RESULTS_DIR}/${sample}_baseline.log)"
    fi
done

echo ""
echo "=========================================="
echo "Results saved to: $RESULTS_DIR"
echo "=========================================="

# Generate summary
echo ""
echo "Summary:"
echo "Sample,Quine_Score,Baseline_Score,Tokens"
for sample in "${SAMPLES[@]}"; do
    TOKENS=$(jq -r '.estimated_tokens // "?"' "${SCRIPT_DIR}/data/${sample}/meta.json" 2>/dev/null || echo "?")
    QUINE_SCORE=$(jq -r '.score // "N/A"' "${RESULTS_DIR}/${sample}_quine_result.json" 2>/dev/null || echo "N/A")
    BASELINE_SCORE=$(jq -r '.score // "N/A"' "${RESULTS_DIR}/${sample}_baseline_result.json" 2>/dev/null || echo "N/A")
    echo "$sample,$QUINE_SCORE,$BASELINE_SCORE,$TOKENS"
done | tee "${RESULTS_DIR}/summary.csv"
