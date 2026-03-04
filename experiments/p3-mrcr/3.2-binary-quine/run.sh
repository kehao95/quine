#!/bin/bash
# Experiment 3.2: Binary Quine
#
# Can quine output a working binary implementation of itself?
#
# Usage: ./run.sh [MODEL]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
MODEL="${1:-claude-opus-4-6}"
MODEL_SHORT="${MODEL##*-}"

RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}"
RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"

mkdir -p "${RUN_DIR}/workspace" "${RUN_DIR}/meta"

cp "${SCRIPT_DIR}/prompt.md" "${RUN_DIR}/meta/prompt-used.md"

echo "=== Experiment 3.2: Binary Quine ==="
echo "Run ID:  ${RUNID}"
echo "Model:   ${MODEL}"
echo ""

[[ -f "${PROJECT_ROOT}/.env" ]] && source "${PROJECT_ROOT}/.env"

QUINE_BIN="${PROJECT_ROOT}/quine"
if [[ ! -x "${QUINE_BIN}" ]]; then
    echo "Building quine..."
    (cd "${PROJECT_ROOT}" && go build -o quine ./cmd/quine)
fi

# Copy quine to workspace so agent can examine it
cp "${QUINE_BIN}" "${RUN_DIR}/workspace/quine"

cd "${RUN_DIR}/workspace"

# Run quine with binary output mode expectation
QUINE_DATA_DIR=".quine" \
QUINE_MODEL_ID="${MODEL}" \
QUINE_MAX_TURNS="${QUINE_MAX_TURNS:-30}" \
QUINE_MAX_DEPTH="${QUINE_MAX_DEPTH:-3}" \
  ./quine "$(cat ../meta/prompt-used.md)" \
  > "../meta/output.bin" \
  2> "../meta/stderr.txt" || true

echo ""
echo "=== Complete ==="
echo "Workspace: ${RUN_DIR}/workspace/"
echo ""

# Check if output is a valid binary
if file "../meta/output.bin" | grep -q "Mach-O"; then
    echo "SUCCESS: Output is a Mach-O binary!"
    cp "../meta/output.bin" "${RUN_DIR}/workspace/q"
    chmod +x "${RUN_DIR}/workspace/q"
    echo "Binary saved to: ${RUN_DIR}/workspace/q"
    echo ""
    echo "Testing with: ./q 'say hello to the world'"
    cd "${RUN_DIR}/workspace"
    timeout 60 ./q "say hello to the world" 2>&1 || echo "(test completed or timed out)"
else
    echo "Output is not a binary. Content type:"
    file "../meta/output.bin"
    echo ""
    echo "First 500 bytes:"
    head -c 500 "../meta/output.bin" | xxd | head -30
fi
