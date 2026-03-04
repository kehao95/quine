#!/bin/bash
# MRCR Mixed Needle Retrieval Experiment
#
# Usage: ./run_mixed.sh [MODEL]
# 
# This tests Quine's ability to find TWO needles in an interleaved conversation.
# Key observations:
# - Early exit behavior (does it wait for both or exit after first?)
# - Wisdom usage (how does it track progress across exec cycles?)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SAMPLE_DIR="${SCRIPT_DIR}/data/mixed_4needle_test"
MODEL="${1:-claude-sonnet-4-20250514}"
MODEL_SHORT="${MODEL##*-}"

# Load metadata
TASK_A_NTH=$(jq -r '.tasks[0].nth' "$SAMPLE_DIR/meta.json")
TASK_A_FORMAT=$(jq -r '.tasks[0].format' "$SAMPLE_DIR/meta.json")
TASK_A_TOPIC=$(jq -r '.tasks[0].topic' "$SAMPLE_DIR/meta.json")
TASK_A_HASH=$(jq -r '.tasks[0].hash' "$SAMPLE_DIR/meta.json")

TASK_B_NTH=$(jq -r '.tasks[1].nth' "$SAMPLE_DIR/meta.json")
TASK_B_FORMAT=$(jq -r '.tasks[1].format' "$SAMPLE_DIR/meta.json")
TASK_B_TOPIC=$(jq -r '.tasks[1].topic' "$SAMPLE_DIR/meta.json")
TASK_B_HASH=$(jq -r '.tasks[1].hash' "$SAMPLE_DIR/meta.json")

# Create run directory
RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}"
RUN_DIR="${SCRIPT_DIR}/runs/mixed_4needle_test/${RUNID}"
mkdir -p "${RUN_DIR}/quine" "${RUN_DIR}/meta"

echo "=========================================="
echo "MRCR Mixed Needle Retrieval Experiment"
echo "=========================================="
echo "Model: $MODEL"
echo "Run ID: $RUNID"
echo "Run dir: $RUN_DIR"
echo ""
echo "Task A: Find the ${TASK_A_NTH}th ${TASK_A_FORMAT} about ${TASK_A_TOPIC}"
echo "        Hash: $TASK_A_HASH"
echo ""
echo "Task B: Find the ${TASK_B_NTH}th ${TASK_B_FORMAT} about ${TASK_B_TOPIC}"
echo "        Hash: $TASK_B_HASH"
echo "=========================================="

# Build mission - TWO targets!
MISSION="You are given a streaming conversation transcript. You must find and extract TWO pieces of content.

## Tasks

**Task A:** Find the ${TASK_A_NTH}th (1-indexed) ${TASK_A_FORMAT} about ${TASK_A_TOPIC}
**Task B:** Find the ${TASK_B_NTH}th (1-indexed) ${TASK_B_FORMAT} about ${TASK_B_TOPIC}

## CRITICAL: This is an EXTRACTION task, not a GENERATION task

Both pieces of content **already exist** in the conversation. You must:
1. Read through to find '[USER]' requests matching each target
2. Extract the '[ASSISTANT]' response that immediately follows
3. Output BOTH results

## Output Format

Your response must be EXACTLY (two lines):
\`\`\`
${TASK_A_HASH}<content_A>
${TASK_B_HASH}<content_B>
\`\`\`

No explanations, no markdown. Just the two hashes followed by extracted content.

## Strategy

Use \`read\` to stream through the conversation. Track TWO counters:
- Count of '${TASK_A_FORMAT} about ${TASK_A_TOPIC}' occurrences
- Count of '${TASK_B_FORMAT} about ${TASK_B_TOPIC}' occurrences

When you find Task A's target (${TASK_A_NTH}th match), save it.
When you find Task B's target (${TASK_B_NTH}th match), save it.
Once you have BOTH, output them and exit.

**If using exec/wisdom:** Pass your current counts and any found content in the wisdom to preserve progress.

Begin reading now."

# Save mission
echo "$MISSION" > "${RUN_DIR}/meta/mission.md"

# Copy sample metadata  
cp "$SAMPLE_DIR/meta.json" "${RUN_DIR}/meta/sample-meta.json"

# Find quine binary
QUINE_BIN="${SCRIPT_DIR}/../../../quine"
if [[ ! -x "$QUINE_BIN" ]]; then
    echo "Building quine..."
    (cd "${SCRIPT_DIR}/../../.." && go build -o quine ./cmd/quine)
fi

# Run Quine
echo ""
echo "Running Quine..."
START_TIME=$(date +%s.%N)

cat "$SAMPLE_DIR/conversation.txt" | \
    QUINE_DATA_DIR="${RUN_DIR}/quine" \
    QUINE_MODEL_ID="${MODEL}" \
    "$QUINE_BIN" "$MISSION" \
    > "${RUN_DIR}/meta/stdout.txt" \
    2> "${RUN_DIR}/meta/stderr.txt" || true

END_TIME=$(date +%s.%N)
ELAPSED=$(echo "$END_TIME - $START_TIME" | bc)

echo "Completed in ${ELAPSED}s"
echo ""

# Check output
echo "=========================================="
echo "Output:"
echo "=========================================="
cat "${RUN_DIR}/meta/stdout.txt"
echo ""
echo "=========================================="

# Check for both hashes
OUTPUT=$(cat "${RUN_DIR}/meta/stdout.txt")

if echo "$OUTPUT" | grep -q "^${TASK_A_HASH}"; then
    echo "✅ Task A hash found: $TASK_A_HASH"
else
    echo "❌ Task A hash NOT found: $TASK_A_HASH"
fi

if echo "$OUTPUT" | grep -q "^${TASK_B_HASH}"; then
    echo "✅ Task B hash found: $TASK_B_HASH"
else
    echo "❌ Task B hash NOT found: $TASK_B_HASH"
fi

# Count sessions (exec usage)
SESSIONS=$(ls "${RUN_DIR}/quine/"*.jsonl 2>/dev/null | wc -l | tr -d ' ')
echo ""
echo "Sessions (exec calls + 1): $SESSIONS"

# Count read calls
READS=$(grep -h "read completed" "${RUN_DIR}/quine/"*.log 2>/dev/null | wc -l | tr -d ' ')
echo "Total read calls: $READS"

echo ""
echo "Results saved to: $RUN_DIR"
