#!/bin/bash
# Exp 2.3: Silent Stigmergy
# Usage: ./run.sh [MODEL] [NUM_AGENTS]
# Example: ./run.sh claude-sonnet-4-5 2

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MODEL="${1:-claude-sonnet-4-5}"
NUM_AGENTS="${2:-2}"
MODEL_SHORT="${MODEL##*-}"

# Use shared library from 1.1-ariadne
LIBRARY_SRC="${SCRIPT_DIR}/../../early-existence/exst-01-ariadne/library"
if [[ ! -d "$LIBRARY_SRC" ]]; then
    echo "ERROR: Library not found at $LIBRARY_SRC"
    exit 1
fi

RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}-${NUM_AGENTS}agents"
RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"
mkdir -p "${RUN_DIR}/meta"

# Shared workspace for all agents
WORKSPACE="${RUN_DIR}/workspace"
mkdir -p "${WORKSPACE}/coordination" "${WORKSPACE}/results"
ln -s "$LIBRARY_SRC" "${WORKSPACE}/library"

# Snapshot prompt
cp "${SCRIPT_DIR}/prompt.md" "${RUN_DIR}/meta/prompt-used.md"

echo "═══════════════════════════════════════════════════════════"
echo "Exp 2.3: Silent Stigmergy"
echo "═══════════════════════════════════════════════════════════"
echo "Run ID:     ${RUNID}"
echo "Model:      ${MODEL}"
echo "Agents:     ${NUM_AGENTS}"
echo "Run dir:    ${RUN_DIR}"
echo "Workspace:  ${WORKSPACE} (shared)"
echo "═══════════════════════════════════════════════════════════"

# Load API credentials
ENV_FILE="${SCRIPT_DIR}/../../../../.env.${MODEL}"
if [[ -f "$ENV_FILE" ]]; then
    echo "Loading credentials from: ${ENV_FILE}"
    source "$ENV_FILE"
else
    echo "WARNING: No env file found at ${ENV_FILE}"
fi

# Build quine if needed
QUINE_BIN="/tmp/quine"
if [[ ! -x "$QUINE_BIN" ]] || [[ "$(find "${SCRIPT_DIR}/../../../../cmd/quine" -newer "$QUINE_BIN" 2>/dev/null | head -1)" ]]; then
    echo "Building quine..."
    (cd "${SCRIPT_DIR}/../../../.." && go build -o "$QUINE_BIN" ./cmd/quine/)
fi

echo ""
echo "Launching ${NUM_AGENTS} agents simultaneously..."
echo ""

# Launch agents in parallel
PIDS=()
for i in $(seq 1 $NUM_AGENTS); do
    AGENT_DIR="${RUN_DIR}/agent-${i}"
    mkdir -p "${AGENT_DIR}/quine"
    
    echo "  Agent $i: starting..."
    
    (
        cd "${WORKSPACE}"
        QUINE_DATA_DIR="${AGENT_DIR}/quine" QUINE_RETENTION_DIR="${AGENT_DIR}/quine/log" \
        QUINE_MODEL_ID="${QUINE_MODEL_ID:-$MODEL}" \
        QUINE_MAX_TURNS=30 \
        AGENT_ID="agent-${i}" \
          "$QUINE_BIN" "$(cat "${RUN_DIR}/meta/prompt-used.md")" \
          > "${AGENT_DIR}/stdout.txt" \
          2> "${AGENT_DIR}/stderr.txt"
    ) &
    PIDS+=($!)
done

echo ""
echo "Waiting for all agents to complete..."

# Wait for all agents
for pid in "${PIDS[@]}"; do
    wait "$pid" || true
done

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "Run complete: ${RUN_DIR}"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Artifacts:"
for i in $(seq 1 $NUM_AGENTS); do
    echo "  Agent $i:"
    echo "    - Tapes:  ${RUN_DIR}/agent-${i}/quine/"
    echo "    - Stdout: ${RUN_DIR}/agent-${i}/stdout.txt"
done
echo ""
echo "Shared workspace:"
echo "  - Coordination: ${WORKSPACE}/coordination/"
echo "  - Results:      ${WORKSPACE}/results/"
echo ""

# Run analysis
if [[ -x "${SCRIPT_DIR}/analysis/analyze.sh" ]]; then
    echo "Running analysis..."
    "${SCRIPT_DIR}/analysis/analyze.sh" "${RUN_DIR}"
fi
