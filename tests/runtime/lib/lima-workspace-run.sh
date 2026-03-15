#!/bin/sh
set -u

TIMEOUT_SECS="$1"
REPO_ROOT="$2"
QUINE_BIN="$3"
MISSION="$4"
EXPORT_DIR="$5"
MAX_TURNS_IN="$6"
PROVIDER_IN="$7"
MODEL_IN="$8"
API_TYPE_IN="$9"
API_BASE_IN="${10}"
API_KEY_IN="${11}"
THINKING_BUDGET_IN="${12}"
USE_SUDO_IN="${13}"
HOST_QUINE_CONFIG_IN="${14}"

WORK_ROOT="$(mktemp -d /tmp/quine-e2e-lima.XXXXXX)"
WORKSPACE_DIR="$WORK_ROOT/workspace"
TAPE_DIR="$WORK_ROOT/quine"
CONFIG_DIR="$WORK_ROOT/quine-config"
mkdir -p "$WORKSPACE_DIR" "$TAPE_DIR" "$CONFIG_DIR"

if [ -d "$HOST_QUINE_CONFIG_IN" ]; then
    cp -f "$HOST_QUINE_CONFIG_IN/kimi-oauth.json" "$CONFIG_DIR/" 2>/dev/null || true
    cp -f "$HOST_QUINE_CONFIG_IN/kimi-device.json" "$CONFIG_DIR/" 2>/dev/null || true
    cp -f "$HOST_QUINE_CONFIG_IN/codex-oauth.json" "$CONFIG_DIR/" 2>/dev/null || true
fi

if [ ! -x "$QUINE_BIN" ]; then
    echo "missing guest quine binary: $QUINE_BIN" >&2
    exit 97
fi

if [ "$USE_SUDO_IN" = "1" ]; then
    sudo -n env \
        QUINE_PROVIDER="$PROVIDER_IN" \
        QUINE_MODEL_ID="$MODEL_IN" \
        QUINE_API_TYPE="$API_TYPE_IN" \
        QUINE_API_BASE="$API_BASE_IN" \
        QUINE_API_KEY="$API_KEY_IN" \
        QUINE_THINKING_BUDGET="$THINKING_BUDGET_IN" \
        QUINE_CONFIG_DIR="$CONFIG_DIR" \
        QUINE_DATA_DIR="$TAPE_DIR" \
        QUINE_WORKSPACE="$WORKSPACE_DIR" \
        QUINE_WORKSPACE_BACKEND="direct" \
        QUINE_MAX_TURNS="$MAX_TURNS_IN" \
        timeout "$TIMEOUT_SECS" "$QUINE_BIN" "$MISSION" \
        >"$WORK_ROOT/stdout.txt" \
        2>"$WORK_ROOT/stderr.txt" \
        </dev/null
    RUN_CODE=$?
else
    timeout "$TIMEOUT_SECS" env \
        QUINE_PROVIDER="$PROVIDER_IN" \
        QUINE_MODEL_ID="$MODEL_IN" \
        QUINE_API_TYPE="$API_TYPE_IN" \
        QUINE_API_BASE="$API_BASE_IN" \
        QUINE_API_KEY="$API_KEY_IN" \
        QUINE_THINKING_BUDGET="$THINKING_BUDGET_IN" \
        QUINE_CONFIG_DIR="$CONFIG_DIR" \
        QUINE_DATA_DIR="$TAPE_DIR" \
        QUINE_WORKSPACE="$WORKSPACE_DIR" \
        QUINE_WORKSPACE_BACKEND="direct" \
        QUINE_MAX_TURNS="$MAX_TURNS_IN" \
        "$QUINE_BIN" "$MISSION" \
        >"$WORK_ROOT/stdout.txt" \
        2>"$WORK_ROOT/stderr.txt" \
        </dev/null
    RUN_CODE=$?
fi

mkdir -p "$EXPORT_DIR"
mkdir -p "$EXPORT_DIR/workspace" "$EXPORT_DIR/quine" "$EXPORT_DIR/quine-config-out"
cp -R "$WORKSPACE_DIR"/. "$EXPORT_DIR/workspace/" 2>/dev/null || true
cp -R "$TAPE_DIR"/. "$EXPORT_DIR/quine/" 2>/dev/null || true
cp -f "$CONFIG_DIR/kimi-oauth.json" "$EXPORT_DIR/quine-config-out/" 2>/dev/null || true
cp -f "$CONFIG_DIR/kimi-device.json" "$EXPORT_DIR/quine-config-out/" 2>/dev/null || true
cp -f "$CONFIG_DIR/codex-oauth.json" "$EXPORT_DIR/quine-config-out/" 2>/dev/null || true
cp "$WORK_ROOT/stdout.txt" "$EXPORT_DIR/stdout.txt" 2>/dev/null || true
cp "$WORK_ROOT/stderr.txt" "$EXPORT_DIR/stderr.txt" 2>/dev/null || true

rm -rf "$WORK_ROOT" 2>/dev/null || true
exit "$RUN_CODE"
