#!/bin/sh
set -eu

TIMEOUT_SECS="$1"
QUINE_BIN="$2"
INPUT_DIR="$3"
OUTPUT_DIR="$4"
MAX_TURNS_IN="$5"
PROVIDER_IN="$6"
MODEL_IN="$7"
API_TYPE_IN="$8"
API_BASE_IN="$9"
API_KEY_IN="${10}"
THINKING_BUDGET_IN="${11}"
HOST_QUINE_CONFIG_IN="${12}"
WORKSPACE_BACKEND_IN="${13:-overlay}"

WORK_ROOT_BASE="/tmp"
if [ "$WORKSPACE_BACKEND_IN" = "overlay" ] && [ -d /dev/shm ] && [ -w /dev/shm ]; then
    WORK_ROOT_BASE="/dev/shm"
fi

WORK_ROOT="$(mktemp -d "$WORK_ROOT_BASE/quine-e2e-lima.XXXXXX")"
WORKSPACE_DIR="$WORK_ROOT/workspace"
TAPE_DIR="$WORK_ROOT/quine"
CONFIG_DIR="$WORK_ROOT/quine-config"
mkdir -p "$WORKSPACE_DIR" "$TAPE_DIR" "$CONFIG_DIR"

run_maybe_sudo() {
    if command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
        sudo -n "$@"
    else
        "$@"
    fi
}

timeout_cmd() {
    if command -v timeout >/dev/null 2>&1; then
        timeout "$@"
        return $?
    fi
    if command -v gtimeout >/dev/null 2>&1; then
        gtimeout "$@"
        return $?
    fi
    if command -v python3 >/dev/null 2>&1; then
        _timeout_secs="$1"; shift
        python3 -c '
import subprocess
import sys

timeout = float(sys.argv[1])
cmd = sys.argv[2:]

try:
    proc = subprocess.Popen(cmd)
except FileNotFoundError:
    sys.exit(127)

try:
    sys.exit(proc.wait(timeout=timeout))
except subprocess.TimeoutExpired:
    proc.kill()
    proc.wait()
    sys.exit(124)
' "$_timeout_secs" "$@"
        return $?
    fi
    echo "missing timeout command: install coreutils for gtimeout, provide timeout in PATH, or ensure python3 is available" >&2
    return 127
}

if [ -d "$HOST_QUINE_CONFIG_IN" ]; then
    cp -f "$HOST_QUINE_CONFIG_IN/kimi-oauth.json" "$CONFIG_DIR/" 2>/dev/null || true
    cp -f "$HOST_QUINE_CONFIG_IN/kimi-device.json" "$CONFIG_DIR/" 2>/dev/null || true
    cp -f "$HOST_QUINE_CONFIG_IN/codex-oauth.json" "$CONFIG_DIR/" 2>/dev/null || true
    cp -f "$HOST_QUINE_CONFIG_IN/copilot-oauth.json" "$CONFIG_DIR/" 2>/dev/null || true
fi

if [ -d "$INPUT_DIR/workspace-in" ]; then
    cp -R "$INPUT_DIR/workspace-in"/. "$WORKSPACE_DIR/" 2>/dev/null || true
fi

if [ ! -x "$QUINE_BIN" ]; then
    echo "missing guest quine binary: $QUINE_BIN" >&2
    exit 97
fi

MISSION_PATH="$INPUT_DIR/mission.txt"
MISSION="$(cat "$MISSION_PATH")"
EXTRA_ENV_FILE="$INPUT_DIR/extra-env.list"
WORKSPACE_ENABLED_MARKER="$INPUT_DIR/workspace-enabled"
WORKSPACE_ESCAPED="$(printf '%s\n' "$WORKSPACE_DIR" | sed 's/[\/&]/\\&/g')"

set -- env \
    QUINE_PROVIDER="$PROVIDER_IN" \
    QUINE_MODEL_ID="$MODEL_IN" \
    QUINE_API_TYPE="$API_TYPE_IN" \
    QUINE_API_BASE="$API_BASE_IN" \
    QUINE_API_KEY="$API_KEY_IN" \
    QUINE_THINKING_BUDGET="$THINKING_BUDGET_IN" \
    QUINE_CONFIG_DIR="$CONFIG_DIR" \
    QUINE_DATA_DIR="$TAPE_DIR" \
    QUINE_RETENTION_DIR="$TAPE_DIR/log" \
    QUINE_MAX_TURNS="$MAX_TURNS_IN"
if [ -f "$EXTRA_ENV_FILE" ]; then
    while IFS= read -r entry; do
        [ -n "$entry" ] || continue
        entry="$(printf '%s\n' "$entry" | sed "s/__QUINE_GUEST_WORKSPACE__/${WORKSPACE_ESCAPED}/g")"
        set -- "$@" "$entry"
    done < "$EXTRA_ENV_FILE"
fi
if [ -f "$WORKSPACE_ENABLED_MARKER" ]; then
    set -- "$@" \
        "QUINE_WORKSPACE=$WORKSPACE_DIR" \
        "QUINE_WORKSPACE_ROOT=$WORKSPACE_DIR" \
        "QUINE_WORKSPACE_REVISION_MODE=restore" \
        "QUINE_WORKSPACE_BACKEND=$WORKSPACE_BACKEND_IN"
fi

if command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
    if run_maybe_sudo "$@" timeout "$TIMEOUT_SECS" "$QUINE_BIN" "$MISSION" \
        >"$WORK_ROOT/stdout.txt" \
        2>"$WORK_ROOT/stderr.txt" \
        </dev/null; then
        RUN_CODE=0
    else
        RUN_CODE=$?
    fi
elif timeout_cmd "$TIMEOUT_SECS" "$@" "$QUINE_BIN" "$MISSION" \
    >"$WORK_ROOT/stdout.txt" \
    2>"$WORK_ROOT/stderr.txt" \
    </dev/null; then
    RUN_CODE=0
else
    RUN_CODE=$?
fi

mkdir -p "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR/workspace-out" "$OUTPUT_DIR/quine" "$OUTPUT_DIR/quine-config-out"
cp -R "$WORKSPACE_DIR"/. "$OUTPUT_DIR/workspace-out/" 2>/dev/null || true
cp -R "$TAPE_DIR"/. "$OUTPUT_DIR/quine/" 2>/dev/null || true
cp -f "$CONFIG_DIR/kimi-oauth.json" "$OUTPUT_DIR/quine-config-out/" 2>/dev/null || true
cp -f "$CONFIG_DIR/kimi-device.json" "$OUTPUT_DIR/quine-config-out/" 2>/dev/null || true
cp -f "$CONFIG_DIR/codex-oauth.json" "$OUTPUT_DIR/quine-config-out/" 2>/dev/null || true
cp -f "$CONFIG_DIR/copilot-oauth.json" "$OUTPUT_DIR/quine-config-out/" 2>/dev/null || true
cp "$WORK_ROOT/stdout.txt" "$OUTPUT_DIR/stdout.txt" 2>/dev/null || true
cp "$WORK_ROOT/stderr.txt" "$OUTPUT_DIR/stderr.txt" 2>/dev/null || true

rm -rf "$WORK_ROOT" 2>/dev/null || true
exit "$RUN_CODE"
