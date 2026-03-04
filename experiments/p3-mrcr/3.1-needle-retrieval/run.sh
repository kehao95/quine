#!/bin/bash
# MRCR Needle Retrieval Experiment
#
# Usage: ./run.sh [SAMPLE_DIR] [MODEL]
# Example: ./run.sh ./data/2needle_0000 claude-sonnet-4-20250514
#
# This script runs Quine on a single MRCR sample and grades the output.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SAMPLE_DIR="${1:?Usage: $0 SAMPLE_DIR [MODEL]}"
MODEL="${2:-claude-sonnet-4-20250514}"
MODEL_SHORT="${MODEL##*-}"

# Validate sample directory
if [[ ! -d "$SAMPLE_DIR" ]]; then
    echo "Error: Sample directory not found: $SAMPLE_DIR" >&2
    exit 1
fi

if [[ ! -f "$SAMPLE_DIR/conversation.txt" ]] || [[ ! -f "$SAMPLE_DIR/meta.json" ]]; then
    echo "Error: Missing conversation.txt or meta.json in $SAMPLE_DIR" >&2
    exit 1
fi

# Extract task parameters from meta.json
HASH=$(jq -r '.random_string_to_prepend' "$SAMPLE_DIR/meta.json")
NTH=$(jq -r '.task.nth' "$SAMPLE_DIR/meta.json")
FORMAT=$(jq -r '.task.format' "$SAMPLE_DIR/meta.json")
TOPIC=$(jq -r '.task.topic' "$SAMPLE_DIR/meta.json")
SAMPLE_ID=$(jq -r '.sample_id' "$SAMPLE_DIR/meta.json")

# Create run directory
RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}"
RUN_DIR="${SCRIPT_DIR}/runs/${SAMPLE_ID}/${RUNID}"
mkdir -p "${RUN_DIR}/quine" "${RUN_DIR}/meta"

echo "=========================================="
echo "MRCR Needle Retrieval Experiment"
echo "=========================================="
echo "Sample: $SAMPLE_ID"
echo "Model: $MODEL"
echo "Run ID: $RUNID"
echo "Run dir: $RUN_DIR"
echo ""
echo "Task: Find the ${NTH}th ${FORMAT} about ${TOPIC}"
echo "Hash prefix: $HASH"
echo "=========================================="

# Build mission from prompt template
# Substitute variables in prompt.md
MISSION=$(cat "$SCRIPT_DIR/prompt.md" | \
    sed "s/{nth}/${NTH}/g" | \
    sed "s/{format}/${FORMAT}/g" | \
    sed "s/{topic}/${TOPIC}/g" | \
    sed "s/{hash}/${HASH}/g" | \
    sed "s/{n}/${NTH}/g")

# Save mission for reference
echo "$MISSION" > "${RUN_DIR}/meta/mission.md"

# Copy sample metadata
cp "$SAMPLE_DIR/meta.json" "${RUN_DIR}/meta/sample-meta.json"

# Find quine binary
QUINE_BIN="${SCRIPT_DIR}/../../../quine"
if [[ ! -x "$QUINE_BIN" ]]; then
    # Try to build it
    echo "Building quine..."
    (cd "${SCRIPT_DIR}/../../.." && go build -o quine ./cmd/quine)
fi

# Run Quine with streaming input
echo ""
echo "Running Quine..."
START_TIME=$(date +%s.%N)

cat "$SAMPLE_DIR/conversation.txt" | \
    QUINE_DATA_DIR="${RUN_DIR}/quine" \
    QUINE_MODEL_ID="${MODEL}" \
    QUINE_MAX_TURNS=8 \
    "$QUINE_BIN" "$MISSION" \
    > "${RUN_DIR}/meta/stdout.txt" \
    2> "${RUN_DIR}/meta/stderr.txt" || true

END_TIME=$(date +%s.%N)
ELAPSED=$(echo "$END_TIME - $START_TIME" | bc)

echo "Completed in ${ELAPSED}s"
echo ""

# Grade the output
echo "Grading..."
python3 "${SCRIPT_DIR}/eval/grade.py" \
    --sample-dir "$SAMPLE_DIR" \
    --output "${RUN_DIR}/meta/stdout.txt" \
    --result "${RUN_DIR}/meta/result.json"

# Add timing to result
jq --arg elapsed "$ELAPSED" '. + {elapsed_seconds: ($elapsed | tonumber)}' \
    "${RUN_DIR}/meta/result.json" > "${RUN_DIR}/meta/result.json.tmp"
mv "${RUN_DIR}/meta/result.json.tmp" "${RUN_DIR}/meta/result.json"

echo ""
echo "=========================================="
echo "Results saved to: ${RUN_DIR}/meta/result.json"
echo "=========================================="
