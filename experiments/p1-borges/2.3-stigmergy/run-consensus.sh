#!/bin/bash
# Exp 2.3b: Stigmergy with Consensus Requirement
# Usage: ./run-consensus.sh [MODEL] [NUM_AGENTS]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MODEL="${1:-claude-sonnet-4-5}"
NUM_AGENTS="${2:-4}"
MODEL_SHORT="${MODEL##*-}"
MAX_TURNS="${3:-50}"

# Use shared library from 1.1-ariadne
LIBRARY_SRC="${SCRIPT_DIR}/../1.1-ariadne/library"
if [[ ! -d "$LIBRARY_SRC" ]]; then
    echo "ERROR: Library not found at $LIBRARY_SRC"
    exit 1
fi

RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}-${NUM_AGENTS}agents-consensus"
RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"
mkdir -p "${RUN_DIR}/meta"

# Shared workspace for all agents
WORKSPACE="${RUN_DIR}/workspace"
mkdir -p "${WORKSPACE}/coordination/proposals"
mkdir -p "${WORKSPACE}/coordination/claims"
mkdir -p "${WORKSPACE}/coordination/findings"
mkdir -p "${WORKSPACE}/results"
ln -s "$LIBRARY_SRC" "${WORKSPACE}/library"

# Initialize active_processes file
touch "${WORKSPACE}/coordination/active_processes"

# Snapshot prompt
cp "${SCRIPT_DIR}/prompt-consensus.md" "${RUN_DIR}/meta/prompt-used.md"

echo "═══════════════════════════════════════════════════════════"
echo "Exp 2.3b: Stigmergy with Consensus"
echo "═══════════════════════════════════════════════════════════"
echo "Run ID:     ${RUNID}"
echo "Model:      ${MODEL}"
echo "Agents:     ${NUM_AGENTS}"
echo "Max Turns:  ${MAX_TURNS}"
echo "Run dir:    ${RUN_DIR}"
echo "═══════════════════════════════════════════════════════════"

# Load API credentials
ENV_FILE="${SCRIPT_DIR}/../../../.env.${MODEL}"
if [[ -f "$ENV_FILE" ]]; then
    echo "Loading credentials from: ${ENV_FILE}"
    source "$ENV_FILE"
else
    echo "WARNING: No env file found at ${ENV_FILE}"
fi

# Build quine if needed
QUINE_BIN="/tmp/quine"
if [[ ! -x "$QUINE_BIN" ]] || [[ "$(find "${SCRIPT_DIR}/../../../cmd/quine" -newer "$QUINE_BIN" 2>/dev/null | head -1)" ]]; then
    echo "Building quine..."
    (cd "${SCRIPT_DIR}/../../.." && go build -o "$QUINE_BIN" ./cmd/quine/)
fi

echo ""
echo "Launching ${NUM_AGENTS} agents (consensus mode)..."
echo ""

# Launch agents in parallel
PIDS=()
for i in $(seq 1 $NUM_AGENTS); do
    AGENT_DIR="${RUN_DIR}/agent-${i}"
    mkdir -p "${AGENT_DIR}/quine"
    
    echo "  Agent $i: starting..."
    
    (
        cd "${WORKSPACE}"
        QUINE_DATA_DIR="${AGENT_DIR}/quine" \
        QUINE_MODEL_ID="${QUINE_MODEL_ID:-$MODEL}" \
        QUINE_MAX_TURNS="${MAX_TURNS}" \
        AGENT_ID="agent-${i}" \
          "$QUINE_BIN" "$(cat "${RUN_DIR}/meta/prompt-used.md")" \
          > "${AGENT_DIR}/stdout.txt" \
          2> "${AGENT_DIR}/stderr.txt"
    ) &
    PIDS+=($!)
done

echo ""
echo "Waiting for all agents to reach consensus..."

# Wait for all agents
for pid in "${PIDS[@]}"; do
    wait "$pid" || true
done

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "Run complete!"
echo "═══════════════════════════════════════════════════════════"

# Check for consensus
echo ""
echo "=== Consensus Check ==="
if [[ -f "${WORKSPACE}/results/consensus.txt" ]]; then
    echo "✅ consensus.txt EXISTS"
    echo ""
    echo "Content:"
    cat "${WORKSPACE}/results/consensus.txt"
else
    echo "❌ consensus.txt NOT FOUND"
fi

echo ""
echo "=== Coordination Artifacts ==="
echo "Proposals:"
ls -la "${WORKSPACE}/coordination/proposals/" 2>/dev/null || echo "  (none)"
echo ""
echo "Claims:"
ls -la "${WORKSPACE}/coordination/claims/" 2>/dev/null || echo "  (none)"
echo ""
echo "Findings:"
ls -la "${WORKSPACE}/coordination/findings/" 2>/dev/null || echo "  (none)"
echo ""
echo "Active Processes:"
cat "${WORKSPACE}/coordination/active_processes" 2>/dev/null || echo "  (empty)"

echo ""
echo "=== Agent Summaries ==="
for i in $(seq 1 $NUM_AGENTS); do
    echo "Agent $i:"
    grep "session ended" "${RUN_DIR}/agent-${i}/quine/"*.log 2>/dev/null || echo "  (no log)"
done
