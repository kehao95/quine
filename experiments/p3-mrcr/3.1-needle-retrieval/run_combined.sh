#!/bin/bash
# Combined 16-needle test: Find 2 needles from 2 different 8-needle haystacks
#
# Task 1: Find the 6th short essay about distance (hash: l4d2BA2kq8)
# Task 2: Find the 2nd social media post about defense (hash: PC4EUlCZBQ)
#
# Usage: ./run_combined.sh [MODEL]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SAMPLE_DIR="${SCRIPT_DIR}/data/combined_16needle_test"
MODEL="${1:-claude-sonnet-4-20250514}"
MODEL_SHORT="${MODEL##*-}"

# Create run directory
RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}"
RUN_DIR="${SCRIPT_DIR}/runs/combined_16needle_test/${RUNID}"
mkdir -p "${RUN_DIR}/quine" "${RUN_DIR}/meta"

echo "=========================================="
echo "Combined 16-Needle Test"
echo "=========================================="
echo "Model: $MODEL"
echo "Run ID: $RUNID"
echo "Run dir: $RUN_DIR"
echo ""
echo "Tasks:"
echo "  1. Find 6th 'short essay' about 'distance' → prepend l4d2BA2kq8"
echo "  2. Find 2nd 'social media post' about 'defense' → prepend PC4EUlCZBQ"
echo ""
echo "Estimated tokens: ~356K (2 x 178K)"
echo "=========================================="

# Build mission
MISSION='You are given a streaming conversation transcript (~356K tokens). Your mission is to find and **extract** TWO specific pieces of content.

## Tasks

You must find BOTH of these:

### Task A: Find the **6th** (1-indexed) **short essay** about **distance**
- Search for: `[USER]` asking for "short essay about distance"
- Extract: The `[ASSISTANT]` response that follows
- Output format: `l4d2BA2kq8{content}` (hash directly followed by content, no space)

### Task B: Find the **2nd** (1-indexed) **social media post** about **defense**
- Search for: `[USER]` asking for "social media post about defense"
- Extract: The `[ASSISTANT]` response that follows
- Output format: `PC4EUlCZBQ{content}` (hash directly followed by content, no space)

## CRITICAL: This is an EXTRACTION task, not a GENERATION task

The content you need **already exists** in the transcript. You must:
1. Read through and find the exact patterns
2. Extract the ASSISTANT responses verbatim
3. Output both results with their respective hashes

## Output Format

Output BOTH results, each on its own section:
```
l4d2BA2kq8{exact content of 6th short essay about distance}
---
PC4EUlCZBQ{exact content of 2nd social media post about defense}
```

## Strategy

Use `read` to stream through the conversation. Track TWO counters:
- `essay_count`: for "short essay about distance" 
- `post_count`: for "social media post about defense"

When you find both targets (essay #6 and post #2), output them and exit.

Use `exec` with `wisdom` to preserve your progress when context fills up:
```
exec(wisdom={
  "essay_count": "3",
  "essay_target": "6", 
  "post_count": "1",
  "post_target": "2",
  "essay_6_content": "...",  // if found
  "post_2_content": "..."    // if found
}, reason="context full, continuing search")
```

Begin reading now. Remember: EXTRACT, do not generate.'

# Save mission
echo "$MISSION" > "${RUN_DIR}/meta/mission.md"
cp "$SAMPLE_DIR/meta.json" "${RUN_DIR}/meta/sample-meta.json"

# Find quine binary
QUINE_BIN="${SCRIPT_DIR}/../../../quine"
if [[ ! -x "$QUINE_BIN" ]]; then
    echo "Building quine..."
    (cd "${SCRIPT_DIR}/../../.." && go build -o quine ./cmd/quine)
fi

echo ""
echo "Running Quine..."
START_TIME=$(date +%s.%N)

cat "$SAMPLE_DIR/conversation.txt" | \
    QUINE_DATA_DIR="${RUN_DIR}/quine" \
    QUINE_MODEL_ID="${MODEL}" \
    QUINE_MAX_TURNS=8 \
    "$QUINE_BIN" -t "$MISSION" \
    > "${RUN_DIR}/meta/stdout.txt" \
    2> "${RUN_DIR}/meta/stderr.txt" || true

END_TIME=$(date +%s.%N)
ELAPSED=$(echo "$END_TIME - $START_TIME" | bc)

echo "Completed in ${ELAPSED}s"
echo ""

# Display results
echo "=========================================="
echo "Results"
echo "=========================================="
echo ""
echo "Output:"
cat "${RUN_DIR}/meta/stdout.txt"
echo ""
echo "=========================================="

# Check for both hashes
echo ""
echo "Verification:"
if grep -q "l4d2BA2kq8" "${RUN_DIR}/meta/stdout.txt"; then
    echo "  ✓ Task A (essay about distance): Hash found"
else
    echo "  ✗ Task A (essay about distance): Hash NOT found"
fi

if grep -q "PC4EUlCZBQ" "${RUN_DIR}/meta/stdout.txt"; then
    echo "  ✓ Task B (social media about defense): Hash found"
else
    echo "  ✗ Task B (social media about defense): Hash NOT found"
fi

echo ""
echo "Run dir: $RUN_DIR"
