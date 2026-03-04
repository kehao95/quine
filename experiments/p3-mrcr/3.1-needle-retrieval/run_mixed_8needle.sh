#!/bin/bash
# MRCR Mixed 8-Needle Retrieval Experiment
#
# Usage: ./run_mixed_8needle.sh [MODEL] [MAX_TURNS]
# 
# Tests wisdom passing across exec cycles with large interleaved data.
# Default MAX_TURNS=2 to force frequent exec and observe wisdom behavior.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SAMPLE_DIR="${SCRIPT_DIR}/data/mixed_8needle_2doc"
MODEL="${1:-claude-sonnet-4-20250514}"
MAX_TURNS="${2:-5}"
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

EST_TOKENS=$(jq -r '.estimated_tokens' "$SAMPLE_DIR/meta.json")
TOTAL_LINES=$(jq -r '.total_lines' "$SAMPLE_DIR/meta.json")

# Create run directory
RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}-t${MAX_TURNS}"
RUN_DIR="${SCRIPT_DIR}/runs/mixed_8needle_2doc/${RUNID}"
mkdir -p "${RUN_DIR}/quine" "${RUN_DIR}/meta"

echo "=========================================="
echo "MRCR Mixed 8-Needle Experiment"
echo "=========================================="
echo "Model: $MODEL"
echo "Max Turns: $MAX_TURNS (low to force exec)"
echo "Run ID: $RUNID"
echo ""
echo "Data: ~${EST_TOKENS} tokens, ${TOTAL_LINES} lines"
echo ""
echo "Task A: Find the ${TASK_A_NTH}th ${TASK_A_FORMAT} about ${TASK_A_TOPIC}"
echo "        Hash: $TASK_A_HASH"
echo ""
echo "Task B: Find the ${TASK_B_NTH}th ${TASK_B_FORMAT} about ${TASK_B_TOPIC}"
echo "        Hash: $TASK_B_HASH"
echo "=========================================="

# Build mission
MISSION="You are processing a streaming conversation transcript to find TWO pieces of content.

## Tasks

**Task A:** Find the ${TASK_A_NTH}th (1-indexed) ${TASK_A_FORMAT} about ${TASK_A_TOPIC}
**Task B:** Find the ${TASK_B_NTH}th (1-indexed) ${TASK_B_FORMAT} about ${TASK_B_TOPIC}

## This is an EXTRACTION task

Both pieces already exist in the conversation. Find them and copy verbatim.

## Output Format (when BOTH are found)

\`\`\`
${TASK_A_HASH}<content_A>
${TASK_B_HASH}<content_B>
\`\`\`

## Strategy

1. Use \`read\` to stream through (50 lines per call)
2. Track TWO counters:
   - Count of '${TASK_A_FORMAT} about ${TASK_A_TOPIC}' 
   - Count of '${TASK_B_FORMAT} about ${TASK_B_TOPIC}'
3. When you find a target (Nth match), SAVE the [ASSISTANT] response content
4. Once BOTH targets are saved, output them and exit

## CRITICAL: Using exec for long contexts

This is a ~${EST_TOKENS} token document. You WILL run out of context.
When you need to \`exec\` to reset context, pass wisdom like:

\`\`\`json
{
  \"progress\": {
    \"essay_count\": 3,
    \"social_count\": 1,
    \"task_A_found\": false,
    \"task_B_found\": true,
    \"task_B_content\": \"...the actual content...\"
  }
}
\`\`\`

Your next incarnation will receive this wisdom and continue from where you left off.

Begin reading now."

# Save mission
echo "$MISSION" > "${RUN_DIR}/meta/mission.md"
cp "$SAMPLE_DIR/meta.json" "${RUN_DIR}/meta/sample-meta.json"

# Find quine binary
QUINE_BIN="${SCRIPT_DIR}/../../../quine"
if [[ ! -x "$QUINE_BIN" ]]; then
    echo "Building quine..."
    (cd "${SCRIPT_DIR}/../../.." && go build -o quine ./cmd/quine)
fi

# Run Quine with low max_turns
echo ""
echo "Running Quine (MAX_TURNS=$MAX_TURNS)..."
START_TIME=$(date +%s.%N)

cat "$SAMPLE_DIR/conversation.txt" | \
    QUINE_DATA_DIR="${RUN_DIR}/quine" \
    QUINE_MODEL_ID="${MODEL}" \
    QUINE_MAX_TURNS="${MAX_TURNS}" \
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
head -50 "${RUN_DIR}/meta/stdout.txt"
echo ""
echo "=========================================="

OUTPUT=$(cat "${RUN_DIR}/meta/stdout.txt")

if echo "$OUTPUT" | grep -q "${TASK_A_HASH}"; then
    echo "✅ Task A hash found: $TASK_A_HASH"
else
    echo "❌ Task A hash NOT found: $TASK_A_HASH"
fi

if echo "$OUTPUT" | grep -q "${TASK_B_HASH}"; then
    echo "✅ Task B hash found: $TASK_B_HASH"
else
    echo "❌ Task B hash NOT found: $TASK_B_HASH"
fi

# Stats
SESSIONS=$(ls "${RUN_DIR}/quine/"*.jsonl 2>/dev/null | wc -l | tr -d ' ')
READS=$(grep -h "read completed" "${RUN_DIR}/quine/"*.log 2>/dev/null | wc -l | tr -d ' ')
EXECS=$(grep -h "called exec" "${RUN_DIR}/quine/"*.log 2>/dev/null | wc -l | tr -d ' ')

echo ""
echo "Stats:"
echo "  Sessions: $SESSIONS"
echo "  Execs: $EXECS"
echo "  Reads: $READS"
echo "  Time: ${ELAPSED}s"
echo ""

# Show wisdom usage
echo "Wisdom usage (grep from logs):"
grep -h "wisdom" "${RUN_DIR}/quine/"*.log 2>/dev/null | head -10 || echo "  (none found)"

echo ""
echo "Results saved to: $RUN_DIR"
