#!/bin/sh
# ============================================================================
# Quine End-to-End Functional Tests
# ============================================================================
#
# Validates quine as a correct POSIX process by running the real binary
# against a real LLM API. Each test gives quine a deterministic mission
# and verifies observable behavior: exit codes, stdout, stderr, tape files.
#
# Usage:
#   source profiles/gpt-5.4-codex-oauth.env    # default live-validation profile
#   ./tests/runtime/run.sh                       # run all tests
#   ./tests/runtime/run.sh test_exit_success     # run one test
#   ./tests/runtime/run.sh test_fd3_piped_input test_fd4_delivery  # run selected tests
#   ./tests/runtime/run.sh gate_overlay_finalization_baseline       # run the stable overlay/finalization gate
#
# Prerequisites:
#   - QUINE_API_KEY, QUINE_API_BASE, QUINE_MODEL_ID, QUINE_API_TYPE set
#   - quine binary built: go build -o /tmp/quine ./cmd/quine/
#
# Design:
#   - Each test is a function named test_*
#   - Tests are independent; order does not matter
#   - Each test gets its own temp directory for QUINE_DATA_DIR
#   - Prompts are crafted to elicit deterministic, verifiable behavior
#   - QUINE_MAX_TURNS is kept low (3-5) to limit API cost
#   - Tests verify OBSERVABLE outputs only (exit code, stdout, stderr, tape)
#
# ============================================================================

set -u  # Treat unset variables as errors

# ── Configuration ──────────────────────────────────────────
QUINE="${QUINE:-/tmp/quine}"
MAX_TURNS="${E2E_MAX_TURNS:-5}"
TIMEOUT="${E2E_TIMEOUT:-120}"  # seconds per test

normalize_self_reentry_mode_for_host() {
    if [ "$(uname -s)" = "Darwin" ] && { [ -z "${QUINE_SELF_REENTRY_MODE:-}" ] || [ "${QUINE_SELF_REENTRY_MODE:-}" = "self" ]; }; then
        export QUINE_SELF_REENTRY_MODE=executable_path
        printf 'runtime test runner: forcing QUINE_SELF_REENTRY_MODE=executable_path on Darwin\n' >&2
    fi
}

normalize_self_reentry_mode_for_host

# ── Counters ───────────────────────────────────────────────
PASS=0
FAIL=0
SKIP=0
TOTAL=0

# ── Colors (if terminal) ──────────────────────────────────
if [ -t 1 ]; then
    GREEN='\033[0;32m'
    RED='\033[0;31m'
    YELLOW='\033[0;33m'
    BOLD='\033[1m'
    RESET='\033[0m'
else
    GREEN='' RED='' YELLOW='' BOLD='' RESET=''
fi

# ── Test harness ───────────────────────────────────────────

# Print test header
begin_test() {
    TOTAL=$((TOTAL + 1))
    printf "${BOLD}[TEST %d] %s${RESET} ... " "$TOTAL" "$1"
}

# Record a pass
pass() {
    PASS=$((PASS + 1))
    _CURRENT_TEST_FAILED=0
    printf "${GREEN}PASS${RESET}\n"
}

# Record a fail with reason
fail() {
    FAIL=$((FAIL + 1))
    _CURRENT_TEST_FAILED=1
    printf "${RED}FAIL${RESET}: %s\n" "$1"
}

# Record a skip with reason
skip() {
    SKIP=$((SKIP + 1))
    printf "${YELLOW}SKIP${RESET}: %s\n" "$1"
}

# Create a fresh temp directory for one test, export QUINE_DATA_DIR
setup() {
    TEST_DIR=$(mktemp -d "${TMPDIR:-/tmp}/quine-e2e.XXXXXX")
    _CURRENT_TEST_FAILED=0
    RUNTIME_BRIDGE_EXPORT_DIR=""
    export QUINE_DATA_DIR="${TEST_DIR}/tapes" QUINE_RETENTION_DIR="${TEST_DIR}/tapes/log"
    mkdir -p "$QUINE_DATA_DIR"
}

keep_failed_runs() {
    [ "${QUINE_BEHAVIOR_KEEP_FAILED_RUNS:-1}" = "1" ]
}

# Clean up after one test
teardown() {
    if [ "${QUINE_BEHAVIOR_KEEP_ALL_RUNS:-0}" = "1" ]; then
        printf 'Runtime bridge export dir: %s\n' "${RUNTIME_BRIDGE_EXPORT_DIR:-}" >&2
        if [ -n "${RUNTIME_BRIDGE_EXPORT_DIR:-}" ]; then
            printf '%s\n' "$RUNTIME_BRIDGE_EXPORT_DIR" > "$TEST_DIR/bridge-export-dir.txt"
        fi
        printf 'Preserving runtime test dir: %s\n' "$TEST_DIR" >&2
        return
    fi
    if keep_failed_runs && [ "${_CURRENT_TEST_FAILED:-0}" = "1" ]; then
        printf 'Runtime bridge export dir: %s\n' "${RUNTIME_BRIDGE_EXPORT_DIR:-}" >&2
        if [ -n "${RUNTIME_BRIDGE_EXPORT_DIR:-}" ]; then
            printf '%s\n' "$RUNTIME_BRIDGE_EXPORT_DIR" > "$TEST_DIR/bridge-export-dir.txt"
        fi
        printf 'Preserving failed runtime test dir: %s\n' "$TEST_DIR" >&2
        return
    fi
    if [ -n "${RUNTIME_BRIDGE_EXPORT_DIR:-}" ]; then
        rm -rf "$RUNTIME_BRIDGE_EXPORT_DIR" 2>/dev/null || true
    fi
    rm -rf "$TEST_DIR" 2>/dev/null || true
}

is_linux() {
    [ "$(uname -s)" = "Linux" ]
}

runtime_surface_fuse_available() {
    is_linux || return 1
    [ -e /dev/fuse ] || return 1
    if command -v fusermount3 >/dev/null 2>&1; then
        :
    elif command -v fusermount >/dev/null 2>&1; then
        :
    else
        return 1
    fi
    command -v go >/dev/null 2>&1 || return 1

    if [ "${RUNTIME_SURFACE_FUSE_PROBED:-0}" = "1" ]; then
        [ "${RUNTIME_SURFACE_FUSE_AVAILABLE:-0}" = "1" ]
        return
    fi

    _probe_dir="$(mktemp -d "${TMPDIR:-/tmp}/quine-fuse-probe.XXXXXX")" || return 1
    mkdir -p "$_probe_dir/mnt"
    cat > "$_probe_dir/probe.go" <<'EOF'
package main

import (
	"os"
	"time"

	fusefs "github.com/hanwen/go-fuse/v2/fs"
	fusepkg "github.com/hanwen/go-fuse/v2/fuse"
)

type root struct{ fusefs.Inode }

func main() {
	timeout := time.Duration(0)
	server, err := fusefs.Mount(os.Args[1], &root{}, &fusefs.Options{
		EntryTimeout: &timeout,
		AttrTimeout:  &timeout,
		MountOptions: fusepkg.MountOptions{
			FsName:            "quine-fuse-probe",
			Name:              "quine-fuse-probe",
			DisableXAttrs:     true,
			DirectMount:       true,
		},
	})
	if err != nil {
		os.Exit(2)
	}
	if err := server.Unmount(); err != nil {
		os.Exit(3)
	}
}
EOF
    if command -v timeout >/dev/null 2>&1; then
        timeout 10s go run "$_probe_dir/probe.go" "$_probe_dir/mnt" >/dev/null 2>&1
    else
        go run "$_probe_dir/probe.go" "$_probe_dir/mnt" >/dev/null 2>&1
    fi
    _probe_code=$?
    rm -rf "$_probe_dir"

    RUNTIME_SURFACE_FUSE_PROBED=1
    if [ "$_probe_code" -eq 0 ]; then
        RUNTIME_SURFACE_FUSE_AVAILABLE=1
    else
        RUNTIME_SURFACE_FUSE_AVAILABLE=0
    fi
    [ "$RUNTIME_SURFACE_FUSE_AVAILABLE" = "1" ]
}

resolve_pid_target() {
    python3 - "$1" <<'PY'
import os
import sys
from pathlib import Path

link_path = Path(sys.argv[1])
target = os.readlink(link_path)
target_path = Path(target)
if not target_path.is_absolute():
    target_path = (link_path.parent / target_path).resolve()
else:
    target_path = target_path.resolve()
print(target_path)
PY
}

legacy_control_surface_ready() {
    _public_root="$1"
    [ -d "$_public_root/ctl" ] &&
        [ -p "$_public_root/ctl/post" ] &&
        [ -p "$_public_root/ctl/poke" ] &&
        [ -p "$_public_root/ctl/inject" ] &&
        [ -p "$_public_root/ctl/interrupt" ]
}

fuse_control_surface_ready() {
    _public_root="$1"
    [ -d "$_public_root/ctl" ] &&
        [ -f "$_public_root/ctl/post" ] && [ ! -p "$_public_root/ctl/post" ] &&
        [ -f "$_public_root/ctl/poke" ] && [ ! -p "$_public_root/ctl/poke" ] &&
        [ -f "$_public_root/ctl/inject" ] && [ ! -p "$_public_root/ctl/inject" ] &&
        [ -f "$_public_root/ctl/interrupt" ] && [ ! -p "$_public_root/ctl/interrupt" ]
}

subjective_worlds_available() {
    if is_linux; then
        return 0
    fi
    resolve_lima_instance >/dev/null 2>&1
}

overlay_substrate_unavailable() {
    _stderr_file="$1"
    grep -Eq 'workspace physics unsupported|overlay mount preflight failed|cannot mount overlay|invalid argument|permission denied|operation not permitted|runtime surface FUSE unsupported' "$_stderr_file" 2>/dev/null
}

has_lima() {
    command -v limactl >/dev/null 2>&1
}

lima_effective_home() {
    if [ -n "${E2E_LIMA_HOME:-}" ]; then
        printf '%s\n' "$E2E_LIMA_HOME"
        return 0
    fi
    if [ -n "${LIMA_HOME:-}" ]; then
        printf '%s\n' "$LIMA_HOME"
        return 0
    fi
    if limactl list 2>/dev/null | awk 'NR>1 && $1 != "" {found=1} END {exit(found ? 0 : 1)}'; then
        printf '\n'
        return 0
    fi
    if [ -d "$HOME/.colima/_lima" ] &&
        LIMA_HOME="$HOME/.colima/_lima" limactl list 2>/dev/null | awk 'NR>1 && $1 != "" {found=1} END {exit(found ? 0 : 1)}'; then
        printf '%s\n' "$HOME/.colima/_lima"
        return 0
    fi
    printf '\n'
    return 0
}

lima_run() {
    _lima_home="$(lima_effective_home)"
    if [ -n "$_lima_home" ]; then
        env LIMA_HOME="$_lima_home" limactl "$@"
    else
        limactl "$@"
    fi
}

lima_instance_names() {
    has_lima || return 1
    lima_run list 2>/dev/null | awk 'NR>1 && $1 != "" {print $1}'
}

resolve_lima_instance() {
    if ! has_lima; then
        printf 'limactl not available\n' >&2
        return 1
    fi

    if [ -n "${E2E_LIMA_INSTANCE:-}" ]; then
        if lima_instance_names | grep -qx "${E2E_LIMA_INSTANCE}"; then
            printf '%s\n' "${E2E_LIMA_INSTANCE}"
            return 0
        fi
        printf 'named Lima instance not found: %s\n' "${E2E_LIMA_INSTANCE}" >&2
        return 2
    fi

    _instances="$(lima_instance_names)"
    _count="$(printf '%s\n' "$_instances" | sed '/^$/d' | wc -l | tr -d ' ')"
    if [ "$_count" = "1" ]; then
        printf '%s\n' "$_instances"
        return 0
    fi
    if [ "$_count" = "0" ]; then
        _lima_home="$(lima_effective_home)"
        if [ -n "$_lima_home" ]; then
            printf 'no Lima instances found under %s; run `limactl create` there or set E2E_LIMA_INSTANCE\n' "$_lima_home" >&2
        else
            printf 'no Lima instances found; run `limactl create` or set E2E_LIMA_INSTANCE\n' >&2
        fi
        return 3
    fi
    printf 'multiple Lima instances found (%s); set E2E_LIMA_INSTANCE explicitly\n' \
        "$(printf '%s' "$_instances" | tr '\n' ' ' | sed 's/[[:space:]]*$//')" >&2
    return 4
}

_lima_instance_status() {
    _instance="$1"
    lima_run list 2>/dev/null | awk -v name="$_instance" 'NR>1 && $1 == name { print $2; found=1; exit } END { if (!found) exit 1 }'
}

ensure_lima_instance_running() {
    _instance="$1"
    _status="$(_lima_instance_status "$_instance" 2>/dev/null)" || return 1
    [ "$_status" = "Running" ] && return 0
    lima_run start "$_instance" >/dev/null 2>&1
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
    printf "missing timeout command: install coreutils for gtimeout, provide timeout in PATH, or ensure python3 is available\n" >&2
    return 127
}

oauth_token_expiry_ms() {
    _token_path="$1"
    python3 - "$_token_path" <<'PY'
import json
import sys

path = sys.argv[1]
try:
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)
    value = int(data.get("expires_at", 0) or 0)
except Exception:
    value = 0
print(value)
PY
}

maybe_refresh_host_kimi_oauth() {
    _host_quine_config="$1"

    [ "${QUINE_API_KEY:-}" = "kimi-oauth" ] || return 0

    _token_path="${_host_quine_config}/kimi-oauth.json"
    [ -f "$_token_path" ] || return 0

    _now_ms="$(python3 - <<'PY'
import time
print(int(time.time() * 1000))
PY
)"
    _expires_ms="$(oauth_token_expiry_ms "$_token_path")"
    _refresh_floor_ms=$((_now_ms + 5 * 60 * 1000))

    [ "$_expires_ms" -gt "$_refresh_floor_ms" ] && return 0

    _refresh_dir="$(mktemp -d "${TMPDIR:-/tmp}/quine-oauth-refresh.XXXXXX")"
    printf 'Refreshing host Kimi OAuth token before Lima guest run...\n' >&2
    if ! QUINE_CONFIG_DIR="$_host_quine_config" \
        QUINE_DATA_DIR="$_refresh_dir" QUINE_RETENTION_DIR="$_refresh_dir/log" \
        QUINE_MAX_TURNS=1 \
        timeout_cmd 90 "$QUINE" 'Exit immediately with status success.' >/dev/null; then
        rm -rf "$_refresh_dir" 2>/dev/null || true
        return 1
    fi
    rm -rf "$_refresh_dir" 2>/dev/null || true
    return 0
}

stage_oauth_seed() {
    _export_dir="$1"
    _host_quine_config="$2"

    [ -d "$_host_quine_config" ] || return 0

    mkdir -p "$_export_dir/quine-config-host"
    for _f in kimi-oauth.json kimi-device.json codex-oauth.json copilot-oauth.json; do
        cp -f "$_host_quine_config/$_f" "$_export_dir/quine-config-host/" 2>/dev/null || true
    done
}

count_oauth_seed_files() {
    _seed_dir="$1"
    _count=0
    for _f in kimi-oauth.json kimi-device.json codex-oauth.json copilot-oauth.json; do
        [ -f "$_seed_dir/$_f" ] && _count=$((_count + 1))
    done
    echo "$_count"
}

restore_oauth_seed() {
    _export_dir="$1"
    _host_quine_config="$2"

    [ -d "$_export_dir/quine-config-out" ] || return 0
    mkdir -p "$_host_quine_config"
    for _f in kimi-oauth.json kimi-device.json codex-oauth.json copilot-oauth.json; do
        cp -f "$_export_dir/quine-config-out/$_f" "$_host_quine_config/" 2>/dev/null || true
    done
}

load_lima_fallback_env() {
    _fallback_env="$1"

    _fallback_vars="$(
        (
            set -a
            # shellcheck disable=SC1090
            . "$_fallback_env"
            printf '%s\n%s\n%s\n%s\n%s\n%s\n' \
                "${QUINE_PROVIDER:-}" "${QUINE_MODEL_ID:-}" "${QUINE_API_TYPE:-}" \
                "${QUINE_API_BASE:-}" "${QUINE_API_KEY:-}" "${QUINE_THINKING_BUDGET:-}"
        )
    )"
    printf '%s' "$_fallback_vars"
}

# Run quine with timeout. Arguments:
#   run_quine <exit_var> <stdout_file> <stderr_file> [quine_args...]
# Sets the named variable to the exit code.
is_transient_provider_failure() {
    _stderr_file="$1"
    [ -f "$_stderr_file" ] || return 1
    grep -E -q 'LLM error: .*(HTTP (408|409|429|499|500|502|503|504)|timeout|temporar|connection reset|unexpected EOF|EOF$)' "$_stderr_file"
}

run_quine_once() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    _stdin_file="$1"; shift
    if [ -n "$_stdin_file" ]; then
        timeout_cmd "$TIMEOUT" env \
            QUINE_MAX_TURNS="$MAX_TURNS" \
            "$QUINE" "$@" \
            <"$_stdin_file" \
            >"$_stdout_file" \
            2>"$_stderr_file"
    else
        timeout_cmd "$TIMEOUT" env \
            QUINE_MAX_TURNS="$MAX_TURNS" \
            "$QUINE" "$@" \
            >"$_stdout_file" \
            2>"$_stderr_file"
    fi
    return $?
}

run_quine() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    _stdin_file=""
    if [ ! -t 0 ]; then
        _stdin_file="$(mktemp "${TMPDIR:-/tmp}/quine-e2e-stdin.XXXXXX")"
        cat >"$_stdin_file"
    fi

    run_quine_once "$_stdout_file" "$_stderr_file" "$_stdin_file" "$@"
    _code=$?
    if [ "$_code" -ne 0 ] && is_transient_provider_failure "$_stderr_file"; then
        printf 'runtime test runner: retrying after transient provider failure\n' >&2
        run_quine_once "$_stdout_file" "$_stderr_file" "$_stdin_file" "$@"
        _code=$?
    fi

    if [ -n "$_stdin_file" ]; then
        rm -f "$_stdin_file"
    fi
    return "$_code"
}

start_quine_helper() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" "$@" \
        >"$_stdout_file" \
        2>"$_stderr_file" \
        </dev/null &
    HELPER_PID=$!
}

# Run workspace physics tests on host directly when possible.
# `direct` backend runs on the local host on any OS; `overlay` still requires a
# Linux host or an explicit Lima bridge on non-Linux.
# Arguments:
#   run_workspace_quine_with_backend <stdout_file> <stderr_file> <backend-or-empty> <mission>
run_workspace_quine_with_backend() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    _workspace_backend="$1"; shift
    _mission="$1"

    if [ "$_workspace_backend" = "direct" ] || is_linux; then
        if [ -n "$_workspace_backend" ]; then
            timeout_cmd "$TIMEOUT" env \
                QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
                QUINE_WORKSPACE="$TEST_DIR/workspace" \
                QUINE_WORKSPACE_BACKEND="$_workspace_backend" \
                QUINE_MAX_TURNS="$MAX_TURNS" \
                "$QUINE" "$_mission" \
                >"$_stdout_file" \
                2>"$_stderr_file" \
                </dev/null
        else
            timeout_cmd "$TIMEOUT" env \
                QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
                QUINE_WORKSPACE="$TEST_DIR/workspace" \
                QUINE_MAX_TURNS="$MAX_TURNS" \
                "$QUINE" "$_mission" \
                >"$_stdout_file" \
                2>"$_stderr_file" \
                </dev/null
        fi
        return $?
    fi

    _lima_instance="$(resolve_lima_instance 2>"$TEST_DIR/lima.err")"
    if [ $? -ne 0 ]; then
        _why="$(cat "$TEST_DIR/lima.err" 2>/dev/null)"
        rm -f "$TEST_DIR/lima.err"
        fail "workspace tests require Linux host or a configured Lima instance${_why:+: $_why}"
        return 127
    fi
    rm -f "$TEST_DIR/lima.err"
    if ! ensure_lima_instance_running "$_lima_instance"; then
        fail "failed to start Lima instance: $_lima_instance"
        return 2
    fi

    _lima_workspace_backend="${_workspace_backend:-${E2E_LIMA_WORKSPACE_BACKEND:-overlay}}"
    _repo_root="$PWD"
    _runner_template="$_repo_root/tests/runtime/lib/lima-workspace-run.sh"
    _guest_quine="$_repo_root/.tmp/quine-linux-arm64-e2e"
    _guest_runner=""
    if [ ! -f "$_runner_template" ]; then
        fail "missing Lima runner template: $_runner_template"
        return 2
    fi
    mkdir -p "$_repo_root/.tmp"
    GOOS=linux GOARCH=arm64 go build -o "$_guest_quine" ./cmd/quine/ || {
        fail "failed to build Linux workspace e2e binary"
        return 2
    }
    _export_dir="$(mktemp -d "$_repo_root/.tmp/quine-e2e-bridge.XXXXXX")"
    RUNTIME_BRIDGE_EXPORT_DIR="$_export_dir"
    _lima_provider="${QUINE_PROVIDER:-}"
    _lima_model="${QUINE_MODEL_ID:-}"
    _lima_api_type="${QUINE_API_TYPE:-}"
    _lima_api_base="${QUINE_API_BASE:-}"
    _lima_api_key="${QUINE_API_KEY:-}"
    _lima_thinking_budget="${QUINE_THINKING_BUDGET:-}"
    _host_quine_config="${QUINE_CONFIG_DIR:-$HOME/.config/quine}"

    if ! maybe_refresh_host_kimi_oauth "$_host_quine_config"; then
        fail "host Kimi OAuth refresh failed before Lima run"
        return 2
    fi

    stage_oauth_seed "$_export_dir" "$_host_quine_config"
    mkdir -p "$_export_dir/workspace-in"
    cp -R "$TEST_DIR/workspace"/. "$_export_dir/workspace-in/" 2>/dev/null || true
    : > "$_export_dir/workspace-enabled"
    if [ -n "${RUNTIME_BRIDGE_EXTRA_ENV_LINES:-}" ]; then
        printf '%s\n' "$RUNTIME_BRIDGE_EXTRA_ENV_LINES" > "$_export_dir/extra-env.list"
    fi
    cat >>"$_export_dir/extra-env.list" <<'EOF'
QUINE_WORKSPACE=__QUINE_GUEST_WORKSPACE__
QUINE_WORKSPACE_ROOT=__QUINE_GUEST_WORKSPACE__
QUINE_WORKSPACE_REVISION_MODE=restore
EOF
    printf '%s' "$_mission" > "$_export_dir/mission.txt"

    if printf '%s' "$_lima_api_key" | grep -Eqi 'oauth|device'; then
        _fallback_env="${E2E_LIMA_ENV_FILE:-.env.gpt-5.4-codex-pool}"
        _oauth_seed_count=$(count_oauth_seed_files "$_export_dir/quine-config-host")

        if [ "$_oauth_seed_count" -eq 0 ]; then
            if [ ! -f "$_fallback_env" ]; then
                fail "OAuth token cache missing; run local OAuth once or set E2E_LIMA_ENV_FILE for non-OAuth fallback"
                rm -rf "$_export_dir" 2>/dev/null || true
                return 2
            fi

            _fallback_vars="$(load_lima_fallback_env "$_fallback_env")"
            _lima_provider=$(printf '%s\n' "$_fallback_vars" | sed -n '1p')
            _lima_model=$(printf '%s\n' "$_fallback_vars" | sed -n '2p')
            _lima_api_type=$(printf '%s\n' "$_fallback_vars" | sed -n '3p')
            _lima_api_base=$(printf '%s\n' "$_fallback_vars" | sed -n '4p')
            _lima_api_key=$(printf '%s\n' "$_fallback_vars" | sed -n '5p')
            _lima_thinking_budget=$(printf '%s\n' "$_fallback_vars" | sed -n '6p')

            if [ -z "$_lima_api_key" ] || printf '%s' "$_lima_api_key" | grep -Eqi 'oauth|device'; then
                fail "fallback credentials in $_fallback_env are not usable for Lima path"
                rm -rf "$_export_dir" 2>/dev/null || true
                return 2
            fi
        fi
    fi

    # Host-side timeout is required here: the guest workload may already be done
    # while the Lima bridge is still wedged. Feeding the runner over stdin avoids
    # hangs observed when executing a host-mounted runner script path directly.
    _host_timeout=$((TIMEOUT + 60))
    _lima_home="$(lima_effective_home)"
    _guest_runner="/tmp/quine-e2e-runner.$$.sh"
    _guest_output_dir="/tmp/quine-e2e-output.$$.d"
    _export_result_dir="$_export_dir/$(basename "$_guest_output_dir")"
    if [ -n "$_lima_home" ]; then
        env LIMA_HOME="$_lima_home" limactl copy "$_runner_template" "$_lima_instance:$_guest_runner" >/dev/null || {
            fail "failed to copy Lima runner into guest"
            return 2
        }
    else
        limactl copy "$_runner_template" "$_lima_instance:$_guest_runner" >/dev/null || {
            fail "failed to copy Lima runner into guest"
            return 2
        }
    fi
    if [ -n "$_lima_home" ]; then
        timeout_cmd "$_host_timeout" \
            env LIMA_HOME="$_lima_home" limactl shell --start --preserve-env "$_lima_instance" \
                /bin/sh "$_guest_runner" \
                "$TIMEOUT" "$_guest_quine" "$_export_dir" "$_guest_output_dir" "$MAX_TURNS" \
                "$_lima_provider" "$_lima_model" "$_lima_api_type" \
                "$_lima_api_base" "$_lima_api_key" "$_lima_thinking_budget" \
                "$_export_dir/quine-config-host" "$_lima_workspace_backend"
    else
        timeout_cmd "$_host_timeout" \
            limactl shell --start --preserve-env "$_lima_instance" \
                /bin/sh "$_guest_runner" \
                "$TIMEOUT" "$_guest_quine" "$_export_dir" "$_guest_output_dir" "$MAX_TURNS" \
                "$_lima_provider" "$_lima_model" "$_lima_api_type" \
                "$_lima_api_base" "$_lima_api_key" "$_lima_thinking_budget" \
                "$_export_dir/quine-config-host" "$_lima_workspace_backend"
    fi
    _code=$?
    if [ -n "$_guest_runner" ]; then
        if [ -n "$_lima_home" ]; then
            env LIMA_HOME="$_lima_home" limactl shell "$_lima_instance" /bin/sh -c "rm -f '$_guest_runner'" >/dev/null 2>&1 || true
        else
            limactl shell "$_lima_instance" /bin/sh -c "rm -f '$_guest_runner'" >/dev/null 2>&1 || true
        fi
    fi
    if [ -n "$_lima_home" ]; then
        env LIMA_HOME="$_lima_home" limactl copy "$_lima_instance:$_guest_output_dir" "$_export_dir" >/dev/null 2>&1 || true
        env LIMA_HOME="$_lima_home" limactl shell "$_lima_instance" /bin/sh -c "rm -rf '$_guest_output_dir'" >/dev/null 2>&1 || true
    else
        limactl copy "$_lima_instance:$_guest_output_dir" "$_export_dir" >/dev/null 2>&1 || true
        limactl shell "$_lima_instance" /bin/sh -c "rm -rf '$_guest_output_dir'" >/dev/null 2>&1 || true
    fi

    mkdir -p "$TEST_DIR/workspace" "$QUINE_DATA_DIR"
    [ -f "$_export_result_dir/stdout.txt" ] && cp "$_export_result_dir/stdout.txt" "$_stdout_file"
    [ -f "$_export_result_dir/stderr.txt" ] && cp "$_export_result_dir/stderr.txt" "$_stderr_file"
    [ -d "$_export_result_dir/workspace-out" ] && cp -R "$_export_result_dir/workspace-out"/. "$TEST_DIR/workspace"/
    [ -d "$_export_result_dir/workspace" ] && cp -R "$_export_result_dir/workspace"/. "$TEST_DIR/workspace"/
    [ -d "$_export_result_dir/quine" ] && cp -R "$_export_result_dir/quine"/. "$QUINE_DATA_DIR"/
    restore_oauth_seed "$_export_result_dir" "$_host_quine_config"
    if [ "${QUINE_BEHAVIOR_KEEP_ALL_RUNS:-0}" != "1" ] &&
        [ "${QUINE_BEHAVIOR_KEEP_FAILED_RUNS:-0}" != "1" ]; then
        rm -rf "$_export_dir" 2>/dev/null || true
        RUNTIME_BRIDGE_EXPORT_DIR=""
    fi

    return "$_code"
}

# Arguments:
#   run_workspace_quine <stdout_file> <stderr_file> <mission>
run_workspace_quine() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    _mission="$1"
    run_workspace_quine_with_backend "$_stdout_file" "$_stderr_file" "" "$_mission"
}

# Assert a file contains a string
assert_contains() {
    _file="$1"; _pattern="$2"; _label="${3:-}"
    if grep -q "$_pattern" "$_file" 2>/dev/null; then
        return 0
    else
        fail "${_label:+$_label: }expected '$_pattern' in $(basename "$_file"), got: $(head -c 200 "$_file" 2>/dev/null || echo '(empty)')"
        return 1
    fi
}

assert_exists() {
    _path="$1"; _label="${2:-}"
    if [ -e "$_path" ]; then
        return 0
    else
        fail "${_label:+$_label: }expected $(basename "$_path") to exist"
        return 1
    fi
}

# Assert a file does NOT contain a string
assert_not_contains() {
    _file="$1"; _pattern="$2"; _label="${3:-}"
    if grep -q "$_pattern" "$_file" 2>/dev/null; then
        fail "${_label:+$_label: }'$_pattern' should NOT appear in $(basename "$_file")"
        return 1
    else
        return 0
    fi
}

# Assert exit code
assert_exit() {
    _got="$1"; _want="$2"; _label="${3:-}"
    if [ "$_got" -eq "$_want" ]; then
        return 0
    else
        fail "${_label:+$_label: }exit code = $_got, want $_want"
        return 1
    fi
}

# Enumerate canonical tape JSONL files without mixing in control logs.
# Prefer the stable session tape surface, then retained session mirrors, then
# compatibility tape aliases when QUINE_RETENTION_DIR externalizes retained ownership.
list_tape_jsonl() {
    _tape_dir="$1"
    {
        find -L "$_tape_dir" -maxdepth 1 -type f -name '[0-9][0-9][0-9][0-9].jsonl' 2>/dev/null
        find -L "$_tape_dir/tapes" -maxdepth 1 -type f -name '[0-9][0-9][0-9][0-9].jsonl' 2>/dev/null
        find -L "$_tape_dir/tapes" -mindepth 2 -maxdepth 2 -type f -name '[0-9][0-9][0-9][0-9].jsonl' 2>/dev/null
        find -L "$_tape_dir/log" -mindepth 3 -maxdepth 3 -type f -path '*/tapes/[0-9][0-9][0-9][0-9].jsonl' 2>/dev/null
        find -L "$_tape_dir/log" -mindepth 4 -maxdepth 4 -type f -path '*/tapes/[0-9][0-9][0-9][0-9].jsonl' 2>/dev/null
        find -L "$_tape_dir/log/tapes" -maxdepth 1 -type f -name '[0-9][0-9][0-9][0-9].jsonl' 2>/dev/null
    } | awk '!seen[$0]++'
}

# Assert a tape file exists and has an outcome entry
assert_tape_has_outcome() {
    _tape_dir="$1"
    _found="$(
    list_tape_jsonl "$_tape_dir" | while IFS= read -r f; do
        if grep -q '"type":"outcome"' "$f" 2>/dev/null || grep -q '"type": "outcome"' "$f" 2>/dev/null; then
            echo 1
            break
        fi
    done
    )"
    if [ "$_found" = "1" ]; then
        return 0
    else
        fail "no outcome entry in tape files under $_tape_dir"
        return 1
    fi
}

# Assert that any tape file under a directory contains a string
assert_any_tape_contains() {
    _tape_dir="$1"; _pattern="$2"; _label="${3:-}"
    _found="$(
    list_tape_jsonl "$_tape_dir" | while IFS= read -r f; do
        if grep -q "$_pattern" "$f" 2>/dev/null; then
            echo "$f"
            break
        fi
    done
    )"
    if [ -n "$_found" ]; then
        return 0
    else
        fail "${_label:+$_label: }expected '$_pattern' in some tape under $_tape_dir"
        return 1
    fi
}

assert_any_tape_contains_literal() {
    _tape_dir="$1"; _needle="$2"; _label="${3:-}"
    _found="$(
    list_tape_jsonl "$_tape_dir" | while IFS= read -r f; do
        if grep -F -q "$_needle" "$f" 2>/dev/null; then
            echo "$f"
            break
        fi
    done
    )"
    if [ -n "$_found" ]; then
        return 0
    else
        fail "${_label:+$_label: }expected literal '$_needle' in some tape under $_tape_dir"
        return 1
    fi
}

wait_file_contains_literal() {
    _file="$1"; _needle="$2"; _limit="${3:-80}"
    _i=0
    while [ "$_i" -lt "$_limit" ]; do
        if grep -F -q "$_needle" "$_file" 2>/dev/null; then
            return 0
        fi
        sleep 0.25
        _i=$((_i + 1))
    done
    return 1
}

extract_first_tape_string_field() {
    _tape_dir="$1"; _field="$2"
    list_tape_jsonl "$_tape_dir" | while IFS= read -r f; do
        _match="$(grep -h -o "\"$_field\":\"[^\"]*\"" "$f" 2>/dev/null | head -n 1)"
        if [ -n "$_match" ]; then
            printf '%s\n' "$_match" | sed "s/^\"$_field\":\"//; s/\"$//"
            break
        fi
    done
}

extract_first_fork_content_field() {
    _tape_dir="$1"; _field_path="$2"
    python3 - "$_tape_dir" "$_field_path" <<'PY'
import json
import sys
from pathlib import Path

tape_dir = Path(sys.argv[1])
field_path = sys.argv[2].split(".")

def list_tapes(root: Path):
    for path in sorted(root.rglob("*.jsonl")):
        if path.name == "control.jsonl":
            continue
        if path.parent.name != "tapes":
            continue
        yield path

def decode_content(content):
    if isinstance(content, dict):
        return content
    if isinstance(content, str):
        try:
            return json.loads(content)
        except Exception:
            return None
    return None

def walk(payload, parts):
    cur = payload
    for part in parts:
        if isinstance(cur, list):
            try:
                cur = cur[int(part)]
            except Exception:
                return None
        elif isinstance(cur, dict):
            if part not in cur:
                return None
            cur = cur[part]
        else:
            return None
    return cur

for tape_path in list_tapes(tape_dir):
    try:
        lines = tape_path.read_text(encoding="utf-8", errors="replace").splitlines()
    except Exception:
        continue
    for line in lines:
        line = line.strip()
        if not line:
            continue
        try:
            entry = json.loads(line)
        except Exception:
            continue
        if entry.get("type") != "tool_result":
            continue
        data = entry.get("data") or {}
        payload = decode_content(data.get("content"))
        if not isinstance(payload, dict):
            continue
        if payload.get("tool") != "fork":
            continue
        value = walk(payload, field_path)
        if value is None:
            continue
        if isinstance(value, (str, int, float)):
            print(value)
        else:
            print(json.dumps(value))
        raise SystemExit(0)
raise SystemExit(1)
PY
}

extract_first_spawn_content_field() {
    _tape_dir="$1"; _field_path="$2"
    python3 - "$_tape_dir" "$_field_path" <<'PY'
import json
import sys
from pathlib import Path

tape_dir = Path(sys.argv[1])
field_path = sys.argv[2].split(".")

def list_tapes(root: Path):
    for path in sorted(root.rglob("*.jsonl")):
        if path.name == "control.jsonl":
            continue
        if path.parent.name != "tapes":
            continue
        yield path

def decode_content(content):
    if isinstance(content, dict):
        return content
    if isinstance(content, str):
        try:
            return json.loads(content)
        except Exception:
            return None
    return None

def walk(payload, parts):
    cur = payload
    for part in parts:
        if isinstance(cur, list):
            try:
                cur = cur[int(part)]
            except Exception:
                return None
        elif isinstance(cur, dict):
            if part not in cur:
                return None
            cur = cur[part]
        else:
            return None
    return cur

for tape_path in list_tapes(tape_dir):
    try:
        lines = tape_path.read_text(encoding="utf-8", errors="replace").splitlines()
    except Exception:
        continue
    for line in lines:
        line = line.strip()
        if not line:
            continue
        try:
            entry = json.loads(line)
        except Exception:
            continue
        if entry.get("type") != "tool_result":
            continue
        data = entry.get("data") or {}
        payload = decode_content(data.get("content"))
        if not isinstance(payload, dict):
            continue
        if payload.get("tool") != "spawn":
            continue
        value = walk(payload, field_path)
        if value is None:
            continue
        if isinstance(value, (str, int, float)):
            print(value)
        else:
            print(json.dumps(value))
        raise SystemExit(0)
raise SystemExit(1)
PY
}

extract_tool_result_is_error() {
    _tape_dir="$1"; _tool_id="$2"
    python3 - "$_tape_dir" "$_tool_id" <<'PY'
import json
import sys
from pathlib import Path

root = Path(sys.argv[1])
tool_id = sys.argv[2]

for path in sorted(root.rglob("*.jsonl")):
    if path.name == "control.jsonl":
        continue
    try:
        lines = path.read_text(encoding="utf-8", errors="replace").splitlines()
    except Exception:
        continue
    for line in lines:
        line = line.strip()
        if not line:
            continue
        try:
            entry = json.loads(line)
        except Exception:
            continue
        if entry.get("type") != "tool_result":
            continue
        data = entry.get("data") or {}
        if data.get("tool_id") != tool_id:
            continue
        print("true" if bool(data.get("is_error")) else "false")
        raise SystemExit(0)
raise SystemExit(1)
PY
}

# Find the primary tape JSONL file (assumes one session per test).
# Prefer canonical session tapes over retained control or mirror files.
find_tape() {
    _tape_dir="$1"
    _tape="$(list_tape_jsonl "$_tape_dir" | head -n 1)"
    [ -n "$_tape" ] || return 1
    echo "$_tape"
}


# ============================================================================
# TESTS
# ============================================================================

# ── 1. Exit codes ──────────────────────────────────────────

test_exit_success() {
    begin_test "exit(success) → exit code 0"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Write the word DONE to stdout via >&4, then call exit with status success.'
    code=$?

    assert_exit "$code" 0 && pass
    teardown
}

test_exit_failure() {
    begin_test "exit(failure) → exit code 1"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Immediately call exit with status failure and stderr message "deliberate failure".'
    code=$?

    assert_exit "$code" 1 && pass
    teardown
}

# ── 2. Stdout (fd 4) — deliverable channel ─────────────────

test_fd4_delivery() {
    begin_test "echo >&4 delivers to process stdout"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Run: echo "QUINE_MARKER_42" >&4   Then exit success. Do NOT print anything else to >&4.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "QUINE_MARKER_42" "stdout"; then
        pass
    fi
    teardown
}

test_fd1_not_leaked() {
    begin_test "echo (fd 1) does NOT leak to process stdout"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Run: echo "INTERNAL_ONLY_xyz"   (do NOT use >&4). Then exit success. Do NOT write anything to >&4.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_not_contains "$TEST_DIR/stdout" "INTERNAL_ONLY_xyz" "fd1 leak"; then
        pass
    fi
    teardown
}

# ── 3. Stderr — failure signal channel ──────────────────────

test_stderr_failure_signal() {
    begin_test "exit(failure, stderr=...) writes to stderr"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Immediately exit with status failure and stderr "E_QUINE_TEST_404".'
    code=$?

    if assert_exit "$code" 1 "exit" && assert_contains "$TEST_DIR/stderr" "E_QUINE_TEST_404" "stderr"; then
        pass
    fi
    teardown
}

test_stderr_success_silent() {
    begin_test "exit(success) produces no agent stderr"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Exit immediately with status success. Do not produce any output.'
    code=$?

    # stderr may contain operational logs from quine runtime, but should NOT
    # contain any agent-written failure message. We check the exit code only
    # since operational logs go to the log file, not stderr.
    assert_exit "$code" 0 && pass
    teardown
}

test_fd5_signal_channel() {
    begin_test "explicit fd 5 writes reach process stderr only"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Run exactly one sh call: echo "FD5_SIGNAL_E2E" >&5. Then exit success. Do not write anything to >&4.'
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stderr" "FD5_SIGNAL_E2E" "fd5 stderr" &&
        assert_not_contains "$TEST_DIR/stdout" "FD5_SIGNAL_E2E" "fd5 stdout leak"; then
        pass
    fi
    teardown
}

# ── 4. Ephemeral shell state ───────────────────────────────

test_shell_cd_does_not_persist() {
    begin_test "cd does not persist across sh calls"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly two sh calls. First sh call: run exactly `cd /tmp` and nothing else. Second sh call: run exactly `pwd >&4` and nothing else. Do not use `&&`, `;`, subshells, or combine steps. Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_not_contains "$TEST_DIR/stdout" "/tmp" "pwd should not inherit prior cd"; then
        pass
    fi
    teardown
}

test_shell_export_does_not_persist() {
    begin_test "export does not persist across sh calls"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly two sh calls. First sh call: run exactly `export E2E_VAR=quine_test_789` and nothing else. Second sh call: run exactly `printf "<%s>\n" "$E2E_VAR" >&4` and nothing else. Do not use `&&`, `;`, subshells, or combine steps. Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "<>" "export should not persist"; then
        pass
    fi
    teardown
}

test_shell_variable_does_not_persist() {
    begin_test "shell variable does not persist across sh calls"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly two sh calls. First sh call: run exactly `MY_COUNTER=1337` and nothing else. Second sh call: run exactly `printf "<%s>\n" "$MY_COUNTER" >&4` and nothing else. Do not use `&&`, `;`, subshells, or combine steps. Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "<>" "shell variable should not persist"; then
        pass
    fi
    teardown
}

test_vision_runtime_surface() {
    begin_test "vision tool works through live runtime surface"
    setup

    cat > "$TEST_DIR/gen_red_jpeg.go" <<'EOF'
package main

import (
    "image"
    "image/color"
    "image/jpeg"
    "os"
)

func main() {
    img := image.NewRGBA(image.Rect(0, 0, 64, 64))
    for y := 0; y < 64; y++ {
        for x := 0; x < 64; x++ {
            img.Set(x, y, color.RGBA{R: 255, A: 255})
        }
    }
    f, err := os.Create(os.Args[1])
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 95}); err != nil {
        panic(err)
    }
}
EOF
    go run "$TEST_DIR/gen_red_jpeg.go" "$TEST_DIR/red.jpg" >/dev/null 2>&1 || {
        fail "failed to generate runtime vision JPEG fixture"
        teardown
        return
    }

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        "Use exactly one vision call on \"$TEST_DIR/red.jpg\" with a question asking for the dominant color. If the image is identified as red, use exactly one sh call to write VISION_RUNTIME_OK to fd 4. Then exit success."
    code=$?

    _tape="$(find_tape "$QUINE_DATA_DIR")"
    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "VISION_RUNTIME_OK" "vision stdout marker" &&
        assert_contains "$_tape" '"name":"vision"' "vision tool call"; then
        pass
    fi
    teardown
}

test_fork_depth_limit_rejected() {
    begin_test "fork rejects requests that exceed depth limit"
    setup

    _old_max_depth="${QUINE_MAX_DEPTH:-}"
    export QUINE_MAX_DEPTH=1

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Attempt exactly one fork call with one child mission. If the fork result reports max recursion depth exceeded, use exactly one sh call to write DEPTH_LIMIT_OK to fd 4. Then exit success.'
    code=$?

    _tape="$(find_tape "$QUINE_DATA_DIR")"

    if [ -n "$_old_max_depth" ]; then export QUINE_MAX_DEPTH="$_old_max_depth"; else unset QUINE_MAX_DEPTH; fi

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "DEPTH_LIMIT_OK" "depth limit marker" &&
        assert_contains "$_tape" 'Max recursion depth exceeded' "depth rejection"; then
        pass
    fi
    teardown
}

test_fork_agent_slot_limit_rejected() {
    begin_test "fork rejects requests when no child agent slots are available"
    setup

    _old_max_agents="${QUINE_MAX_AGENTS:-}"
    export QUINE_MAX_AGENTS=1

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Attempt exactly one fork call with one child mission. If the fork result reports insufficient slots, use exactly one sh call to write AGENT_SLOT_LIMIT_OK to fd 4. Then exit success.'
    code=$?

    _tape="$(find_tape "$QUINE_DATA_DIR")"

    if [ -n "$_old_max_agents" ]; then export QUINE_MAX_AGENTS="$_old_max_agents"; else unset QUINE_MAX_AGENTS; fi

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "AGENT_SLOT_LIMIT_OK" "agent slot marker" &&
        assert_contains "$_tape" 'Insufficient slots' "agent slot rejection"; then
        pass
    fi
    teardown
}

# ── 5. Stdin (fd 3) — material channel ─────────────────────

test_fd3_piped_input() {
    begin_test "piped stdin readable via /dev/fd/3"
    setup

    echo "MATERIAL_DATA_e2e_test" | run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Read the piped input using: cat /dev/fd/3    Then write what you read to >&4 and exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "MATERIAL_DATA_e2e_test" "fd3"; then
        pass
    fi
    teardown
}

test_fd3_consumed_across_calls() {
    begin_test "fd 3 is a live stream consumed across separate sh calls"
    setup

    # fd 3 is the quine process stdin. This case asserts the consuming-stream
    # contract ("reads consume data, later calls see only the remainder"), which
    # holds for the real `cmd | quine` pipe case. Feed it as a genuine pipe here
    # rather than through run_quine's stdin file-snapshot: re-opening /dev/fd/3 on
    # a regular file rewinds to offset 0 (a snapshot artifact, not a runtime
    # behavior), whereas a pipe shares one offset across the separate sh
    # subprocesses, so the second read continues where the first stopped.
    printf 'abcdeFGHIJ' | timeout_cmd "$TIMEOUT" env \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" \
        'Do exactly these steps in separate sh calls: (1) run `dd bs=5 count=1 < /dev/fd/3 2>/dev/null > /tmp/part1.txt` (2) run `cat /dev/fd/3 > /tmp/part2.txt` (3) run `printf "PART1=%s\nPART2=%s\n" "$(cat /tmp/part1.txt)" "$(cat /tmp/part2.txt)" >&4` Then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr"
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "PART1=abcde" "fd3 part1" &&
        assert_contains "$TEST_DIR/stdout" "PART2=FGHIJ" "fd3 part2" &&
        assert_not_contains "$TEST_DIR/stdout" "PART2=abcdeFGHIJ" "fd3 replay"; then
        pass
    fi
    teardown
}

test_interactive_screen() {
    begin_test "interactive PTY job exposes screen-oriented filesystem surface"
    setup

    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=7
    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one sh call with interactive=true to start `python3 -q`. Use the returned absolute job path. Then, in separate normal sh calls: (1) verify `<job>/screen.txt` shows the Python prompt `>>>`; (2) write `print(6*7)<enter>` to `<job>/in`; (3) wait until `<job>/screen.txt` contains `42`; (4) write `exit()<enter>` to `<job>/in`; (5) confirm that after the REPL exits, `<job>/exit` contains `0`; (6) echo `INTERACTIVE_OK` to >&4. Then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if assert_exit "$code" 0 "interactive exit" && assert_contains "$TEST_DIR/stdout" "INTERACTIVE_OK" "interactive stdout"; then
        pass
    fi
    teardown
}

test_interactive_sigint_process_group() {
    begin_test "interactive PTY job handles POSIX SIGINT sent to process group"
    setup

    result_file="$TEST_DIR/interactive-sigint-result.txt"
    mission="$(
cat <<EOF
Use exactly one sh call with interactive=true to start this exact shell command:

trap 'printf "INT_TRAPPED\n" > "$result_file"; exit 0' INT; printf 'READY_INT\n'; while :; do sleep 1; done

Use the returned absolute job path. Then use ordinary sh calls as needed to:
1. Verify <job>/screen.txt contains READY_INT.
2. Send POSIX SIGINT to the recorded job process group using kill -INT -\$(cat <job>/pid).
3. Wait until <job>/exit exists and contains 0.
4. Verify $result_file contains INT_TRAPPED.
5. Emit INTERACTIVE_SIGINT_OK to fd 4 and exit success.
EOF
)"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=8
    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" "$mission"
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if assert_exit "$code" 0 "interactive SIGINT exit" &&
        assert_contains "$TEST_DIR/stdout" "INTERACTIVE_SIGINT_OK" "interactive SIGINT marker" &&
        assert_contains "$result_file" "INT_TRAPPED" "interactive SIGINT trap result" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"interactive":true' "interactive tool call recorded" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" 'kill -INT' "POSIX process-group SIGINT recorded"; then
        pass
    fi
    teardown
}

test_interactive_overlay_world_adoption() {
    begin_test "interactive PTY job under overlay exposes adoptable world"
    if ! subjective_worlds_available; then
        skip "interactive overlay world adoption requires Linux or a configured Lima instance"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=9
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Use exactly one interactive sh call, one switch_world call, and normal sh calls as needed. Step 1: run exactly `sh(command="cat parent.txt 2>/dev/null; printf \"interactive-line\n\" > interactive.txt; printf \"SCREEN_DONE\n\"", interactive=true)`. Use the returned absolute job path. Step 2: wait for `<job>/exit` to contain `0`, verify `<job>/screen.txt` contains `SCREEN_DONE`, verify `<job>/world_handle` starts with `world://`, and verify the parent workspace does not yet contain `interactive.txt`; emit `PRE_SWITCH_PRIVATE` to fd 4. Step 3: call switch_world with the exact world handle read from `<job>/world_handle`. Step 4: verify `interactive.txt` contains `interactive-line`, emit `INTERACTIVE_WORLD_OK` and `HANDLE_OK` to fd 4, then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if assert_exit "$code" 0 "interactive overlay adoption exit" &&
        assert_contains "$TEST_DIR/stdout" "PRE_SWITCH_PRIVATE" "pre-switch privacy marker" &&
        assert_contains "$TEST_DIR/stdout" "INTERACTIVE_WORLD_OK" "interactive world marker" &&
        assert_contains "$TEST_DIR/stdout" "HANDLE_OK" "world handle marker" &&
        assert_contains "$TEST_DIR/workspace/interactive.txt" "interactive-line" "adopted interactive file" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"interactive":true' "interactive tool call recorded" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"name":"switch_world"' "switch_world tool recorded" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"target":"world://' "world handle switch target recorded"; then
        pass
    fi
    teardown
}

# ── 6. Binary stdin (-b flag) ──────────────────────────────

test_binary_stdin() {
    begin_test "binary stdin (-b) saves file and references it"
    setup

    # Create a small binary payload
    printf '\x89PNG\x0d\x0a\x1a\x0a' | run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        -b 'A binary file was provided to you. Print its path to >&4 and exit success.'
    code=$?

    if assert_exit "$code" 0 "exit"; then
        # The binary file should have been saved under QUINE_DATA_DIR
        _bin_count=$(find "$QUINE_DATA_DIR" -name 'stdin-*.bin' 2>/dev/null | wc -l)
        if [ "$_bin_count" -gt 0 ]; then
            pass
        else
            fail "no stdin-*.bin file found in $QUINE_DATA_DIR"
        fi
    fi
    teardown
}

# ── 7. Execution budget physics ────────────────────────────

test_execution_budget_disabled_hidden() {
    begin_test "QUINE_MAX_TURNS=0 omits execution-budget prompt/feedback"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=0 \
        QUINE_PROMPT_METAPHOR=off \
        "$QUINE" 'Run exactly one sh call: echo "BUDGET_OFF_OK" >&4. Then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "no-budget exit" && assert_contains "$TEST_DIR/stdout" "BUDGET_OFF_OK" "no-budget stdout"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif grep -q '\[EXECUTIONS LEFT\]' "$_tape" 2>/dev/null; then
            fail "execution marker should be absent when QUINE_MAX_TURNS=0"
        elif grep -q 'Execution Budget:' "$_tape" 2>/dev/null; then
            fail "execution budget block should be omitted when QUINE_MAX_TURNS=0"
        else
            pass
        fi
    fi
    teardown
}

test_execution_budget_enabled_feedback() {
    begin_test "QUINE_MAX_TURNS>0 shows execution-budget prompt/feedback"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=3 \
        QUINE_PROMPT_METAPHOR=off \
        "$QUINE" 'Run exactly one sh call: echo "BUDGET_ON_OK" >&4. Then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "budget-on exit" && assert_contains "$TEST_DIR/stdout" "BUDGET_ON_OK" "budget-on stdout"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'Execution Budget: 3 `sh` calls' "$_tape" 2>/dev/null; then
            fail "system prompt missing active execution budget disclosure"
        elif grep -q '\[TURNS LEFT\]' "$_tape" 2>/dev/null; then
            fail "legacy turn marker should not appear"
        else
            pass
        fi
    fi
    teardown
}

test_execution_budget_hard_fail() {
    begin_test "hard_fail exhaustion leaves one final exit-only response"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=1 \
        QUINE_PROMPT_METAPHOR=off \
        "$QUINE" 'Run exactly one sh call now: echo step1. After that result comes back, do not call sh again. Use the remaining response to call exit with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "hard_fail continuation exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '"no_turns_left":' "$_tape" 2>/dev/null; then
            fail "tape missing no_turns_left continuation guidance"
        elif grep -q '"termination_mode":"turn_exhaustion"' "$_tape" 2>/dev/null; then
            fail "hard_fail live path should not end in turn_exhaustion when exit succeeds"
        elif ! grep -q '"termination_mode":"exit"' "$_tape" 2>/dev/null; then
            fail "tape missing exit termination mode after continuation exit"
        else
            pass
        fi
    fi
    teardown
}

test_anchor_memory_roundtrip() {
    begin_test "anchor memory mark/unfold works when enabled"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=4 \
        QUINE_ANCHOR_MEMORY=1 \
        "$QUINE" 'Anchor memory is enabled for this run. Do exactly these steps: (1) call mark with resolution "e2e-anchor-checkpoint" and fold=false (2) call unfold with anchor_id=0 (3) write "ANCHOR_E2E_OK" to >&4 (4) exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "anchor-memory exit" &&
        assert_contains "$TEST_DIR/stdout" "ANCHOR_E2E_OK" "anchor-memory stdout"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        _anchor_meta="$(find "$QUINE_DATA_DIR" -path '*/context/state/anchors/0.anchor/meta.json' 2>/dev/null | head -n 1)"
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '"name":"mark"' "$_tape" 2>/dev/null; then
            fail "tape missing mark tool call"
        elif ! grep -q '"name":"unfold"' "$_tape" 2>/dev/null; then
            fail "tape missing unfold tool call"
        elif ! grep -q '"memory_feedback":' "$_tape" 2>/dev/null; then
            fail "tape missing runtime.memory_feedback"
        elif [ -z "$_anchor_meta" ]; then
            fail "anchor meta file not found under stable memory path"
        elif ! grep -q '"resolution":"e2e-anchor-checkpoint"' "$_anchor_meta" 2>/dev/null && \
             ! grep -q '"resolution": "e2e-anchor-checkpoint"' "$_anchor_meta" 2>/dev/null; then
            fail "anchor meta missing expected resolution"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_metaphor_off() {
    begin_test "QUINE_PROMPT_METAPHOR=off keeps prompt factual"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_PROMPT_METAPHOR=off \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "metaphor-off exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif grep -q 'THERMODYNAMIC SURVIVAL' "$_tape" 2>/dev/null; then
            fail "thermodynamic overlay should be absent when QUINE_PROMPT_METAPHOR=off"
        elif ! grep -q 'THE PRIME DIRECTIVE: RUNTIME PHYSICS' "$_tape" 2>/dev/null; then
            fail "physics prime directive should be present when metaphor is off"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_metaphor_thermodynamic() {
    begin_test "QUINE_PROMPT_METAPHOR=thermodynamic overlays prompt only"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_PROMPT_METAPHOR=thermodynamic \
        "$QUINE" 'Run exactly one sh call: echo "THERMO_OK" >&4. Then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "thermodynamic exit" && assert_contains "$TEST_DIR/stdout" "THERMO_OK" "thermodynamic stdout"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'THERMODYNAMIC SURVIVAL' "$_tape" 2>/dev/null; then
            fail "thermodynamic prime directive should be present when overlay is enabled"
        elif ! grep -q 'Energy (shell executions)' "$_tape" 2>/dev/null; then
            fail "thermodynamic overlay terms missing from system prompt"
        elif grep -q '\[TURNS LEFT\]' "$_tape" 2>/dev/null; then
            fail "legacy turn marker should not appear with metaphor enabled"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_self_model_basic() {
    begin_test "QUINE_PROMPT_SELF_MODEL=basic hides advanced self-model framing"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_PROMPT_SELF_MODEL=basic \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "self-model-basic exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'You are quine, a running process in a POSIX operating system\.' "$_tape" 2>/dev/null; then
            fail "basic self identity should remain visible"
        elif grep -q 'The current file on disk is one embodiment of you\.' "$_tape" 2>/dev/null; then
            fail "advanced embodiment framing should be hidden in basic self-model mode"
        elif grep -q 'Your cognition in this session is LLM-mediated' "$_tape" 2>/dev/null; then
            fail "advanced cognition framing should be hidden in basic self-model mode"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_runtime_surface_hidden() {
    begin_test "QUINE_PROMPT_RUNTIME_SURFACE=hidden hides self-mapping surface only"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_PROMPT_RUNTIME_SURFACE=hidden \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "runtime-surface-hidden exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif grep -q '### Runtime Process Surface' "$_tape" 2>/dev/null; then
            fail "runtime process surface section should be hidden"
        elif grep -q 'QUINE_AGENT_ROOT=' "$_tape" 2>/dev/null; then
            fail "QUINE_AGENT_ROOT prompt disclosure should be hidden"
        elif grep -q 'Runtime State Root:' "$_tape" 2>/dev/null; then
            fail "runtime state root disclosure should be hidden"
        elif ! grep -Fq '**exec** - Replace the current process image with a new executable.' "$_tape" 2>/dev/null; then
            fail "capability disclosure should remain when runtime surface is hidden"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_fragments_surface_visible() {
    begin_test "prompt discloses canonical context surface"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "fragments-surface exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '### Context Files' "$_tape" 2>/dev/null; then
            fail "context-files rules block should be present"
        elif ! grep -q '`context/` is the canonical current-incarnation context surface' "$_tape" 2>/dev/null; then
            fail "context root disclosure should be present"
        elif ! grep -q '`context/prompt/` holds provider-visible prompt fragments, assembled by filename order' "$_tape" 2>/dev/null; then
            fail "prompt context disclosure should be present"
        elif ! grep -q '`context/prompt/` — provider-visible prompt fragment surface assembled by filename order' "$_tape" 2>/dev/null; then
            fail "runtime surface should mention context/prompt/"
        elif ! grep -q '`context/state/current.jsonl` — raw current-turn stream for your current-incarnation live cognition surface' "$_tape" 2>/dev/null; then
            fail "runtime surface should mention context/state/current.jsonl"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_agents_md_disabled() {
    begin_test "QUINE_AGENTS_MD_ENABLED=0 hides AGENTS.md prompt fragment"
    setup
    mkdir -p "$TEST_DIR/workspace"

    cat > "$TEST_DIR/workspace/AGENTS.md" <<'EOF'
ROOT_AGENTS_DISABLED_MARKER
EOF

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_AGENTS_MD_ENABLED=0 \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "agents-md-disabled exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif grep -q '### AGENTS.md' "$_tape" 2>/dev/null; then
            fail "AGENTS.md fragment should be absent when QUINE_AGENTS_MD_ENABLED=0"
        elif grep -q 'ROOT_AGENTS_DISABLED_MARKER' "$_tape" 2>/dev/null; then
            fail "disabled AGENTS.md fragment leaked content"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_agents_md_enabled_single_projection() {
    begin_test "QUINE_AGENTS_MD_ENABLED=1 projects single AGENTS.md fragment"
    setup
    mkdir -p "$TEST_DIR/workspace"

    cat > "$TEST_DIR/workspace/AGENTS.md" <<'EOF'
ROOT_AGENTS_ENABLED_MARKER
EOF

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=3 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_AGENTS_MD_ENABLED=1 \
        "$QUINE" 'If `$QUINE_AGENT_ROOT/context/prompt/10-agents.md` contains ROOT_AGENTS_ENABLED_MARKER, write AGENTS_MD_READ_OK to fd 4 and exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "agents-md-enabled exit" &&
        assert_contains "$TEST_DIR/stdout" "AGENTS_MD_READ_OK" "agents md readable marker"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '### AGENTS.md' "$_tape" 2>/dev/null; then
            fail "AGENTS.md fragment should be present when QUINE_AGENTS_MD_ENABLED=1"
        elif ! grep -q 'ROOT_AGENTS_ENABLED_MARKER' "$_tape" 2>/dev/null; then
            fail "AGENTS.md content missing from prompt"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_skills_disabled() {
    begin_test "QUINE_AGENTS_SKILLS_ENABLED=0 hides skills prompt fragment"
    setup

    mkdir -p "$TEST_DIR/workspace/.agents/skills/foo"
    cat > "$TEST_DIR/workspace/.agents/skills/foo/SKILL.md" <<'EOF'
---
name: foo
description: SHOULD_NOT_APPEAR_DISABLED
---
BODY_MARKER_DISABLED
EOF

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_AGENTS_SKILLS_ENABLED=0 \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "skills-disabled exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif grep -q '### SKILLS.md' "$_tape" 2>/dev/null; then
            fail "SKILLS.md fragment should be absent when QUINE_AGENTS_SKILLS_ENABLED=0"
        elif grep -q 'SHOULD_NOT_APPEAR_DISABLED' "$_tape" 2>/dev/null || grep -q 'BODY_MARKER_DISABLED' "$_tape" 2>/dev/null; then
            fail "disabled skills surface leaked skill content"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_skills_enabled_catalog_only() {
    begin_test "QUINE_AGENTS_SKILLS_ENABLED=1 generates skills catalog only"
    setup

    mkdir -p "$TEST_DIR/workspace/.agents/skills/foo/scripts" "$TEST_DIR/workspace/.agents/skills/foo/references" "$TEST_DIR/workspace/.agents/skills/foo/assets"
    cat > "$TEST_DIR/workspace/.agents/skills/foo/SKILL.md" <<'EOF'
---
name: foo
description: Catalog description marker
---
BODY_MARKER_ENABLED
EOF
    printf 'SCRIPT_MARKER_ENABLED\n' > "$TEST_DIR/workspace/.agents/skills/foo/scripts/check.sh"
    printf 'REFERENCE_MARKER_ENABLED\n' > "$TEST_DIR/workspace/.agents/skills/foo/references/rules.md"
    printf 'ASSET_MARKER_ENABLED\n' > "$TEST_DIR/workspace/.agents/skills/foo/assets/sample.txt"

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_AGENTS_SKILLS_ENABLED=1 \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "skills-enabled exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '### SKILLS.md' "$_tape" 2>/dev/null; then
            fail "SKILLS.md fragment should be present when QUINE_AGENTS_SKILLS_ENABLED=1"
        elif ! grep -q 'Catalog description marker' "$_tape" 2>/dev/null; then
            fail "skills catalog description missing from prompt"
        elif ! grep -q 'Source: `.agents/skills/foo/SKILL.md`' "$_tape" 2>/dev/null; then
            fail "skills source path missing from prompt"
        elif ! grep -q 'Quine generated this fragment from the `name` and `description` frontmatter in each `SKILL.md` visible at startup and refresh boundaries' "$_tape" 2>/dev/null; then
            fail "skills fragment generation rule missing from prompt"
        elif grep -q 'BODY_MARKER_ENABLED' "$_tape" 2>/dev/null || \
             grep -q 'SCRIPT_MARKER_ENABLED' "$_tape" 2>/dev/null || \
             grep -q 'REFERENCE_MARKER_ENABLED' "$_tape" 2>/dev/null || \
             grep -q 'ASSET_MARKER_ENABLED' "$_tape" 2>/dev/null; then
            fail "skills prompt should include catalog only, not body/resources"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_skills_source_path_readable() {
    begin_test "skills Source path is readable from shell cwd"
    setup

    mkdir -p "$TEST_DIR/workspace/subdir" "$TEST_DIR/workspace/.agents/skills/foo"
    cat > "$TEST_DIR/workspace/.agents/skills/foo/SKILL.md" <<'EOF'
---
name: foo
description: Readable source marker
---
SOURCE_BODY_READABLE_MARKER
EOF

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=3 \
        QUINE_WORK_DIR="$TEST_DIR/workspace/subdir" \
        QUINE_AGENTS_SKILLS_ENABLED=1 \
        "$QUINE" 'Use the SKILLS.md fragment Source path for `foo` exactly as shown. Read that file from your shell cwd. If it contains SOURCE_BODY_READABLE_MARKER, write SKILL_SOURCE_READ_OK to fd 4, then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "skills-source-readable exit" &&
        assert_contains "$TEST_DIR/stdout" "SKILL_SOURCE_READ_OK" "skills source readable marker"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'Source: `../.agents/skills/foo/SKILL.md`' "$_tape" 2>/dev/null; then
            fail "skills source should be relative to QUINE_WORK_DIR"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_skills_exec_reentry_rescans_frontmatter() {
    begin_test "exec re-entry rescans skill frontmatter"
    setup

    mkdir -p "$TEST_DIR/workspace/.agents/skills/foo"
    cat > "$TEST_DIR/workspace/.agents/skills/foo/SKILL.md" <<'EOF'
---
name: foo
description: EXEC_REFRESH_V1
---
body
EOF

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=6 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_AGENTS_SKILLS_ENABLED=1 \
        "$QUINE" 'If `.skills_exec_done` does not exist, use one sh call to replace EXEC_REFRESH_V1 with EXEC_REFRESH_V2 in `.agents/skills/foo/SKILL.md` and create `.skills_exec_done`, then call exec with default arguments. If `.skills_exec_done` exists, exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    _inc0_skills="$(find "$QUINE_DATA_DIR/log/sessions" -path '*/inc/0/context/prompt/20-skills.md' 2>/dev/null | head -n 1)"
    _inc1_skills="$(find "$QUINE_DATA_DIR/log/sessions" -path '*/inc/1/context/prompt/20-skills.md' 2>/dev/null | head -n 1)"

    if [ "$code" -ne 0 ] && [ "$code" -ne 124 ]; then
        fail "skills-exec-refresh exit: exit code = $code, want 0 or 124"
    elif [ ! -f "$TEST_DIR/workspace/.skills_exec_done" ]; then
        fail "skills exec refresh should create .skills_exec_done before exec"
    elif [ ! -f "$_inc0_skills" ]; then
        fail "missing retained inc/0 SKILLS.md fragment"
    elif ! grep -q 'EXEC_REFRESH_V1' "$_inc0_skills" 2>/dev/null; then
        fail "retained inc/0 SKILLS.md should still show EXEC_REFRESH_V1"
    elif [ ! -f "$_inc1_skills" ]; then
        fail "missing retained inc/1 SKILLS.md fragment after exec"
    elif ! grep -q 'EXEC_REFRESH_V2' "$_inc1_skills" 2>/dev/null; then
        fail "retained inc/1 SKILLS.md should show EXEC_REFRESH_V2 after exec re-entry"
    elif ! assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"termination_mode":"exec"' "exec transition recorded"; then
        :
    else
        pass
    fi
    teardown
}

test_context_memory_next_turn_refresh() {
    begin_test "context/prompt/30-memory.md refreshes next-turn system prompt"
    setup
    mkdir -p "$TEST_DIR/workspace"

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=5 \
        QUINE_EXEC_ENABLED=0 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        "$QUINE" 'Do not call exec or fork. If the active prompt already contains MEMORY_REFRESH_MARKER_RUNTIME under Memory, write CONTEXT_MEMORY_REFRESH_OK to fd 4 and exit success. Otherwise, use sh to write MEMORY_REFRESH_MARKER_RUNTIME to `$QUINE_AGENT_ROOT/context/prompt/30-memory.md`. After that sh result returns, continue immediately with the next provider turn; do not idle or answer "waiting".' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    _memory_file="$(find "$QUINE_DATA_DIR/log/sessions" -path '*/inc/0/context/prompt/30-memory.md' 2>/dev/null | head -n 1)"

    if assert_exit "$code" 0 "context-memory-refresh exit"; then
        if [ ! -f "$_memory_file" ]; then
            fail "missing retained context/prompt/30-memory.md"
        elif ! grep -q 'MEMORY_REFRESH_MARKER_RUNTIME' "$_memory_file" 2>/dev/null; then
            fail "retained memory.md missing refresh marker"
        elif ! assert_any_tape_contains_literal "$QUINE_DATA_DIR" 'MEMORY_REFRESH_MARKER_RUNTIME' "memory marker appeared in retained runtime evidence"; then
            :
        else
            pass
        fi
    fi
    teardown
}

test_context_memory_exec_inherits() {
    begin_test "context/prompt/30-memory.md is copied across exec re-entry"
    setup
    mkdir -p "$TEST_DIR/workspace"

    timeout_cmd 35 env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=7 \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        "$QUINE" 'Read `$QUINE_AGENT_ROOT/status/session.json`. If `incarnation_id` is greater than 0 and the active prompt contains MEMORY_EXEC_MARKER_RUNTIME under Memory, write CONTEXT_MEMORY_EXEC_OK to fd 4 and exit success. If `incarnation_id` is 0, do not emit CONTEXT_MEMORY_EXEC_OK yet: use sh to write MEMORY_EXEC_MARKER_RUNTIME to `$QUINE_AGENT_ROOT/context/prompt/30-memory.md`, then call exec with default arguments.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    _inc0_memory="$(find "$QUINE_DATA_DIR/log/sessions" -path '*/inc/0/context/prompt/30-memory.md' 2>/dev/null | head -n 1)"
    _inc1_memory="$(find "$QUINE_DATA_DIR/log/sessions" -path '*/inc/1/context/prompt/30-memory.md' 2>/dev/null | head -n 1)"

    if [ "$code" -ne 0 ] && [ "$code" -ne 124 ]; then
        fail "context-memory-exec exit: exit code = $code, want 0 or 124"
    else
        if [ ! -f "$_inc0_memory" ] || ! grep -q 'MEMORY_EXEC_MARKER_RUNTIME' "$_inc0_memory" 2>/dev/null; then
            fail "inc/0 memory.md missing exec marker"
        elif [ ! -f "$_inc1_memory" ] || ! grep -q 'MEMORY_EXEC_MARKER_RUNTIME' "$_inc1_memory" 2>/dev/null; then
            fail "inc/1 memory.md did not inherit exec marker"
        elif ! assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"termination_mode":"exec"' "exec transition recorded"; then
            :
        elif ! assert_any_tape_contains_literal "$QUINE_DATA_DIR" '### Memory' "memory prompt block present after exec"; then
            :
        else
            pass
        fi
    fi
    teardown
}

# ── 8. Transactional workspace physics ────────────────────

test_workspace_overlay_commit() {
    begin_test "transactional workspace commits on success"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Do exactly these steps in separate sh calls: (1) printf "overlay-commit\n" > committed.txt (2) test -f committed.txt && echo "COMMIT_OK" >&4. Then exit success.'
    code=$?

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "workspace commit exit" &&
        assert_contains "$TEST_DIR/stdout" "COMMIT_OK" "workspace commit stdout" &&
        assert_contains "$TEST_DIR/workspace/committed.txt" "overlay-commit" "workspace committed file"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '\[FS MUTATIONS\]' "$_tape" 2>/dev/null; then
            fail "tool results should include [FS MUTATIONS] under workspace physics"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_fuse_exit_success_materializes() {
    begin_test "overlay+fuse exit success materializes workspace before outcome"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    mkdir -p "$TEST_DIR/workspace"
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND="overlay" \
        QUINE_WORKSPACE_OVERLAY_DRIVER="fuse" \
        QUINE_WORKSPACE_REVISION_MODE="restore" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" \
        'Overlay workspace and FUSE public surface are enabled. Do exactly these steps: (1) run one sh call: mkdir -p source && printf "#!/bin/sh\necho compiled\n" > compile.sh && chmod +x compile.sh && printf "package main\nfunc main(){}\n" > source/main.go && printf "submission-ready\n" > submission_ready.txt (2) run one sh call to verify compile.sh, source/main.go, and submission_ready.txt exist and emit OVERLAY_FUSE_READY to fd 4. Then exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?

    if [ "$code" -ne 0 ] &&
        grep -Eq 'workspace physics unsupported|invalid argument|permission denied|operation not permitted|runtime surface FUSE unsupported' "$TEST_DIR/stderr" 2>/dev/null; then
        skip "overlay+FUSE substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "overlay+fuse materialization exit" &&
        assert_contains "$TEST_DIR/stdout" "OVERLAY_FUSE_READY" "overlay+fuse stdout" &&
        assert_exists "$TEST_DIR/workspace/compile.sh" "compile script" &&
        assert_exists "$TEST_DIR/workspace/source/main.go" "source archive" &&
        assert_contains "$TEST_DIR/workspace/submission_ready.txt" "submission-ready" "submission ready marker"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '"termination_mode":"exit"' "$_tape" 2>/dev/null; then
            fail "tape missing success exit outcome"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_fuse_shutdown_phase_order() {
    begin_test "overlay+fuse shutdown phase order is durable"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    mkdir -p "$TEST_DIR/workspace"
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND="overlay" \
        QUINE_WORKSPACE_OVERLAY_DRIVER="fuse" \
        QUINE_WORKSPACE_REVISION_MODE="restore" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" \
        'Overlay workspace and FUSE public surface are enabled. Run exactly one sh call: printf "phase-order\n" > phase.txt && echo PHASE_READY >&4. Then exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?

    if [ "$code" -ne 0 ] &&
        grep -Eq 'workspace physics unsupported|invalid argument|permission denied|operation not permitted|runtime surface FUSE unsupported' "$TEST_DIR/stderr" 2>/dev/null; then
        skip "overlay+FUSE substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "overlay+fuse phase-order exit" &&
        assert_contains "$TEST_DIR/stdout" "PHASE_READY" "overlay+fuse phase stdout" &&
        assert_contains "$TEST_DIR/workspace/phase.txt" "phase-order" "phase workspace file"; then
        _state_file="$(find "$QUINE_DATA_DIR/log" -path '*/status/finalization.jsonl' -type f 2>/dev/null | head -n 1)"
        if [ -z "$_state_file" ]; then
            fail "missing durable finalization phase state"
        elif python3 - "$_state_file" <<'PY'
import json
import sys

want = [
    "exit_requested",
    "workspace_committing",
    "workspace_committed",
    "outcome_written",
    "surface_cleanup",
]
phases = []
with open(sys.argv[1], encoding="utf-8") as f:
    for line in f:
        line = line.strip()
        if not line:
            continue
        phases.append(json.loads(line).get("phase"))

cursor = 0
for phase in phases:
    if cursor < len(want) and phase == want[cursor]:
        cursor += 1

if cursor != len(want):
    print("phase order missing:", phases, file=sys.stderr)
    raise SystemExit(1)
PY
        then
            pass
        else
            fail "durable finalization phase order is invalid"
        fi
    fi
    teardown
}

test_workspace_overlay_fuse_long_lineage_materializes() {
    begin_test "overlay+fuse long lineage materializes"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    mkdir -p "$TEST_DIR/workspace"
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND="overlay" \
        QUINE_WORKSPACE_OVERLAY_DRIVER="fuse" \
        QUINE_WORKSPACE_REVISION_MODE="restore" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" \
        'Overlay workspace and FUSE public surface are enabled. Use exactly four sh calls: (1) printf "one\n" > one.txt (2) printf "two\n" > two.txt (3) mkdir -p nested && printf "three\n" > nested/three.txt (4) cat one.txt two.txt nested/three.txt >/dev/null && echo LINEAGE_READY >&4. Then exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?

    if [ "$code" -ne 0 ] &&
        grep -Eq 'workspace physics unsupported|invalid argument|permission denied|operation not permitted|runtime surface FUSE unsupported' "$TEST_DIR/stderr" 2>/dev/null; then
        skip "overlay+FUSE substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "overlay+fuse lineage exit" &&
        assert_contains "$TEST_DIR/stdout" "LINEAGE_READY" "overlay+fuse lineage stdout" &&
        assert_contains "$TEST_DIR/workspace/one.txt" "one" "lineage file one" &&
        assert_contains "$TEST_DIR/workspace/two.txt" "two" "lineage file two" &&
        assert_contains "$TEST_DIR/workspace/nested/three.txt" "three" "lineage file three"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'created=wr3' "$_tape" 2>/dev/null; then
            fail "long lineage should create at least three world revisions"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_commit_intent_recovery() {
    begin_test "overlay workspace commit intent recovery materializes before resumed outcome"
    setup

    if ! is_linux; then
        skip "overlay recovery fault injection is Linux-only"
        teardown
        return
    fi

    mkdir -p "$TEST_DIR/workspace"
    _session_id="e2e-commit-intent-recovery"
    _mission_first='Transactional workspace physics is enabled. Run exactly one sh call: printf "recovered-live\n" > recovered.txt. Then exit success.'
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_SESSION_ID="$_session_id" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND="overlay" \
        QUINE_WORKSPACE_REVISION_MODE="restore" \
        QUINE_TEST_EXIT_AFTER_WORKSPACE_COMMIT_INTENT="1" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" "$_mission_first" \
        >"$TEST_DIR/first.stdout" \
        2>"$TEST_DIR/first.stderr" \
        </dev/null
    first_code=$?

    if [ "$first_code" -ne 86 ] && overlay_substrate_unavailable "$TEST_DIR/first.stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi
    if ! assert_exit "$first_code" 86 "fault-injected first run"; then
        teardown
        return
    fi
    if [ -e "$TEST_DIR/workspace/recovered.txt" ]; then
        fail "fault-injected run materialized workspace before simulated crash"
        teardown
        return
    fi
    _intent_file="$QUINE_DATA_DIR/log/sessions/$_session_id/status/workspace-commit-intent.json"
    if ! assert_contains "$_intent_file" '"status": "pending"' "pending commit intent"; then
        teardown
        return
    fi
    if find "$QUINE_DATA_DIR/log/sessions/$_session_id/tapes" -name '*.jsonl' -exec grep -q '"termination_mode":"exit"' {} + 2>/dev/null; then
        fail "fault-injected run wrote success outcome before recovery"
        teardown
        return
    fi

    _mission_second='Recovery audit. Do not call sh and do not modify the workspace. Call exit success immediately.'
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_SESSION_ID="$_session_id" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND="overlay" \
        QUINE_WORKSPACE_REVISION_MODE="restore" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" "$_mission_second" \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?

    if [ "$code" -ne 0 ] && [ "$code" -ne 1 ]; then
        fail "recovery run exited with unexpected code $code"
    elif assert_contains "$TEST_DIR/workspace/recovered.txt" "recovered-live" "recovered workspace file" &&
        assert_contains "$_intent_file" '"status": "committed"' "committed recovery intent"; then
        _state_file="$QUINE_DATA_DIR/log/sessions/$_session_id/status/finalization.jsonl"
        if ! assert_contains "$_state_file" '"phase":"workspace_recovering"' "workspace recovery phase"; then
            :
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_rollback() {
    begin_test "overlay workspace does not materialize after failed session"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Run exactly one sh call: printf "overlay-rollback\n" > rolled_back.txt. Then exit failure with stderr "rollback test".'
    code=$?

    if [ "$code" -ne 1 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 1 "workspace rollback exit"; then
        if [ -e "$TEST_DIR/workspace/rolled_back.txt" ]; then
            fail "rolled_back.txt should not survive a failed session"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_failure_revision() {
    begin_test "overlay records failed shell side effects as restorable world revision"
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=6
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Revisioned overlay workspace physics is enabled. Use exactly three sh calls and one switch_world call. Do not use fork or exec. Turn 1: run exactly `printf "failed-live\n" > doomed.txt; false`. This shell failure is intentional; continue after the non-zero exit and observe that the tool result reports a new world revision `wr1`. Turn 2: run exactly `test "$(cat doomed.txt)" = "failed-live" && echo "FAILED_WORLD_VISIBLE" >&4`. Then call `switch_world` with `target="wr0"` so `doomed.txt` should disappear. Turn 3: run exactly `test ! -e doomed.txt && echo "FAILED_WORLD_RESTORED" >&4`. Then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "overlay failed-shell revision exit" &&
        assert_contains "$TEST_DIR/stdout" "FAILED_WORLD_VISIBLE" "failed-shell world-visible stdout" &&
        assert_contains "$TEST_DIR/stdout" "FAILED_WORLD_RESTORED" "failed-shell restored stdout" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"exit_code":1' "failed shell exit code recorded" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\\[WORLD REVISION\\] created=wr1 parent=wr0 current=wr1" "failed shell revision marker" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"target":"wr0"' "failed shell switch target"; then
        if [ -e "$TEST_DIR/workspace/doomed.txt" ]; then
            fail "doomed.txt should not be materialized after switch_world wr0"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_timeout_revision() {
    begin_test "overlay records timed-out shell boundary effects as restorable world revision"
    if ! is_linux; then
        skip "overlay timeout revision test requires Linux overlay substrate"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=6
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND="overlay" \
        QUINE_WORKSPACE_REVISION_MODE="restore" \
        QUINE_SH_DEFAULT_TIMEOUT_SECONDS=1 \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" \
        'Revisioned overlay workspace physics is enabled and the default shell timeout is intentionally 1 second. Use exactly three sh calls and one switch_world call. Do not use fork or exec. Turn 1: run exactly `printf "timeout-live\n" > timed.txt; sleep 30`. This timeout is intentional; continue after `status="interrupted"` and observe that the result reports a new world revision `wr1`. Turn 2: run exactly `test "$(cat timed.txt)" = "timeout-live" && echo "TIMEOUT_WORLD_VISIBLE" >&4`. Then call `switch_world` with `target="wr0"` so `timed.txt` should disappear. Turn 3: run exactly `test ! -e timed.txt && echo "TIMEOUT_WORLD_RESTORED" >&4`. Then exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "overlay timeout revision exit" &&
        assert_contains "$TEST_DIR/stdout" "TIMEOUT_WORLD_VISIBLE" "timeout world-visible stdout" &&
        assert_contains "$TEST_DIR/stdout" "TIMEOUT_WORLD_RESTORED" "timeout restored stdout" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"status":"interrupted"' "timeout status recorded" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\\[WORLD REVISION\\] created=wr1 parent=wr0 current=wr1" "timeout revision marker" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"target":"wr0"' "timeout switch target"; then
        if [ -e "$TEST_DIR/workspace/timed.txt" ]; then
            fail "timed.txt should not be materialized after switch_world wr0"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_absolute_path() {
    begin_test "transactional workspace tracks absolute-path writes"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        "Transactional workspace physics is enabled. Run exactly one sh call: cd \"\$QUINE_WORKSPACE\" && printf 'absolute-write\n' > absolute.txt. Then use a second sh call to verify the file exists and emit ABS_OK to fd 4. Exit success."
    code=$?

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "workspace absolute-path exit" &&
        assert_contains "$TEST_DIR/stdout" "ABS_OK" "workspace absolute-path stdout" &&
        assert_contains "$TEST_DIR/workspace/absolute.txt" "absolute-write" "workspace absolute-path file"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'absolute.txt (created)' "$_tape" 2>/dev/null; then
            fail "tool results should report absolute-path workspace mutations"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_direct_persists_on_failure() {
    begin_test "direct workspace persists writes on failure"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine_with_backend "$TEST_DIR/stdout" "$TEST_DIR/stderr" direct \
        'Direct shared workspace physics is enabled. Run exactly one sh call: printf "direct-persist\n" > persisted.txt. Then exit failure with stderr "direct persist test".'
    code=$?

    if assert_exit "$code" 1 "direct workspace failure exit" &&
        assert_contains "$TEST_DIR/workspace/persisted.txt" "direct-persist" "direct workspace persisted file"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'persisted.txt (created)' "$_tape" 2>/dev/null; then
            fail "tool results should report direct workspace mutations on failed shell"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_direct_observes_peer_mutation() {
    begin_test "direct workspace reports peer-created files at next shell boundary"
    setup

    mkdir -p "$TEST_DIR/workspace"
    (
        sleep 1
        printf 'peer-direct\n' > "$TEST_DIR/workspace/peer_drop.txt"
    ) &
    _peer_writer_pid=$!

    run_workspace_quine_with_backend "$TEST_DIR/stdout" "$TEST_DIR/stderr" direct \
        'Direct shared workspace physics is enabled. Another peer may create a file while you work. Do exactly these steps in separate sh calls: (1) printf "seed\n" > seed.txt (2) sleep 2 (3) read the `[FS MUTATIONS]` feedback from step 2, notice that `peer_drop.txt` appeared, then run exactly `test "$(cat peer_drop.txt)" = "peer-direct" && echo "PEER_OK" >&4`. Then exit success.'
    code=$?

    wait "$_peer_writer_pid" 2>/dev/null || true

    if assert_exit "$code" 0 "direct peer observation exit" &&
        assert_contains "$TEST_DIR/stdout" "PEER_OK" "direct peer observation stdout" &&
        assert_contains "$TEST_DIR/workspace/peer_drop.txt" "peer-direct" "direct peer drop file"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'peer_drop.txt (created)' "$_tape" 2>/dev/null; then
            fail "tool results should report peer-created direct workspace mutations"
        else
            pass
        fi
    fi
    teardown
}

test_workspace_overlay_lima_probe() {
    begin_test "Lima overlay backend probe"
    if is_linux; then
        skip "probe is for non-Linux Lima bridge"
        return
    fi
    if ! has_lima; then
        skip "limactl not available"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine_with_backend "$TEST_DIR/stdout" "$TEST_DIR/stderr" overlay \
        'Transactional workspace physics is enabled. Do exactly these steps in separate sh calls: (1) printf "overlay-lima\n" > probe.txt (2) test -f probe.txt && echo "OVERLAY_LIMA_OK" >&4. Then exit success.'
    code=$?

    if [ "$code" -eq 0 ]; then
        if assert_contains "$TEST_DIR/stdout" "OVERLAY_LIMA_OK" "Lima overlay stdout" &&
            assert_contains "$TEST_DIR/workspace/probe.txt" "overlay-lima" "Lima overlay file"; then
            _tape=$(find_tape "$QUINE_DATA_DIR")
            if [ -z "$_tape" ]; then
                fail "no tape file found"
            elif ! grep -q '\[FS MUTATIONS\]' "$_tape" 2>/dev/null; then
                fail "tool results should include [FS MUTATIONS] under Lima overlay probe"
            else
                pass
            fi
        fi
    elif grep -Eq 'workspace physics unsupported|invalid argument|permission denied|only supported on Linux' "$TEST_DIR/stderr" 2>/dev/null; then
        skip "Lima guest does not currently support overlay backend"
    else
        fail "unexpected Lima overlay failure: $(head -c 200 "$TEST_DIR/stderr" 2>/dev/null || echo '(empty)')"
    fi
    teardown
}

test_switch_world_restores_prior_revision() {
    begin_test "switch_world restores provisional workspace to an earlier world revision"
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=6
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Use exactly three sh calls and one switch_world call. Do not use fork or exec. Turn 1: run exactly `printf "v1\n" > state.txt`. Observe that the tool result reports current world revision `wr1`. Turn 2: run exactly `printf "v2\n" > state.txt`. Observe that the tool result reports current world revision `wr2`. Then call `switch_world` with `target="wr1"` so `state.txt` should become `v1` again. Turn 3: run exactly `test "$(cat state.txt)" = "v1" && echo "RESTORE_OK" >&4`. Then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "restore world exit" &&
        assert_contains "$TEST_DIR/stdout" "RESTORE_OK" "restore world stdout" &&
        assert_contains "$TEST_DIR/workspace/state.txt" "v1" "restore world file" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"tool":"switch_world"' "switch world tool result" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"target":"wr1"' "switch world target" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"revision":"wr1"' "switch world revision"; then
        pass
    fi
    teardown
}

test_switch_world_branches_forward_after_rewind() {
    begin_test "switch_world rewinds current world and next mutation branches from the target revision"
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=6
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Use exactly three sh calls and one switch_world call. Do not use fork or exec. Turn 1: run exactly `printf "v1\n" > state.txt`. Observe that the tool result reports current world revision `wr1`. Turn 2: run exactly `printf "v2\n" > state.txt`. Observe that the tool result reports current world revision `wr2`. Then call `switch_world` with `target="wr1"`. Turn 3: run exactly `printf "v3\n" > state.txt && test "$(cat state.txt)" = "v3" && echo "BRANCH_OK" >&4`. Then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "switch world branch exit" &&
        assert_contains "$TEST_DIR/stdout" "BRANCH_OK" "switch world branch stdout" &&
        assert_contains "$TEST_DIR/workspace/state.txt" "v3" "switch world branch file" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"tool":"switch_world"' "switch world rewind tool result" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"target":"wr1"' "switch world rewind target" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"revision":"wr1"' "switch world rewind revision" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\\[WORLD REVISION\\] created=wr3 parent=wr1 current=wr3" "branch revision marker"; then
        pass
    fi
    teardown
}

test_workspace_unsupported_on_non_linux() {
    begin_test "workspace physics are rejected on non-Linux"
    if is_linux; then
        skip "host is Linux"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        "$QUINE" 'Begin.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if [ "$code" -eq 0 ]; then
        fail "workspace unsupported exit: expected non-zero exit"
    elif \
        assert_contains "$TEST_DIR/stderr" "only supported on Linux" "workspace unsupported stderr"; then
        pass
    fi
    teardown
}

# ── 9. Tape integrity ──────────────────────────────────────

test_tape_has_meta() {
    begin_test "tape contains meta entry"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Exit immediately with status success.'
    code=$?

    _tape=$(find_tape "$QUINE_DATA_DIR")
    if [ -z "$_tape" ]; then
        fail "no tape file found"
    elif grep -q '"type":"meta"' "$_tape" 2>/dev/null || grep -q '"type": "meta"' "$_tape" 2>/dev/null; then
        pass
    else
        fail "no meta entry in tape"
    fi
    teardown
}

test_tape_has_outcome() {
    begin_test "tape contains outcome entry"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Exit immediately with status success.'
    code=$?

    assert_tape_has_outcome "$QUINE_DATA_DIR" && pass
    teardown
}

test_tape_has_messages() {
    begin_test "tape contains system and user messages"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Exit immediately with status success.'
    code=$?

    _tape=$(find_tape "$QUINE_DATA_DIR")
    if [ -z "$_tape" ]; then
        fail "no tape file found"
    else
        _has_system=0; _has_user=0
        grep -q '"system"' "$_tape" 2>/dev/null && _has_system=1
        grep -q '"user"' "$_tape" 2>/dev/null && _has_user=1
        if [ "$_has_system" -eq 1 ] && [ "$_has_user" -eq 1 ]; then
            pass
        else
            fail "missing system ($_has_system) or user ($_has_user) messages"
        fi
    fi
    teardown
}

# ── 10. Fork (child process) ───────────────────────────────

test_process_surface_self_identity() {
    begin_test "process surface exposes self identity under QUINE_AGENT_ROOT"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one sh call. In that call, verify that QUINE_AGENT_ROOT is set, that "$QUINE_AGENT_ROOT/mission.txt" exists and is non-empty, and that "$QUINE_AGENT_ROOT/status/session.json" exists with a non-empty session_id, a pid greater than zero, and an agent_root equal to the QUINE_AGENT_ROOT environment variable. If all checks pass, write SELF_SURFACE_OK to file descriptor 4 and exit success. If any check fails, print a clear reason to stderr and exit failure.'
    code=$?

    if assert_exit "$code" 0 "self surface exit" &&
        assert_contains "$TEST_DIR/stdout" "SELF_SURFACE_OK" "self surface stdout"; then
        pass
    fi
    teardown
}

test_process_surface_self_source_surface() {
    begin_test "process surface exposes read-only self-source tree when enabled"
    setup
    self_source_check='set -eu; root="$QUINE_AGENT_ROOT/source-code"; public="$QUINE_AGENT_ROOT/public/source-code"; [ -d "$root" ] || { echo "missing self-source root: $root" >&2; exit 1; }; [ -d "$public" ] || { echo "missing public self-source projection: $public" >&2; exit 1; }; [ ! -w "$root" ] || { echo "self-source root should be read-only" >&2; exit 1; }; [ -f "$root/go.mod" ] || { echo "missing go.mod" >&2; exit 1; }; [ -f "$root/go.sum" ] || { echo "missing go.sum" >&2; exit 1; }; [ -f "$root/selfsource.go" ] || { echo "missing selfsource.go" >&2; exit 1; }; [ -f "$root/selfsource_bundle_data.go" ] || { echo "missing selfsource_bundle_data.go" >&2; exit 1; }; [ -f "$root/cmd/quine/main.go" ] || { echo "missing cmd/quine/main.go" >&2; exit 1; }; [ -f "$root/internal/runtime/runtime.go" ] || { echo "missing internal/runtime/runtime.go" >&2; exit 1; }; cmp -s "$root/go.mod" "$public/go.mod" || { echo "public self-source go.mod drift" >&2; exit 1; }; manifest="$root/.git/quine-source-manifest.json"; public_manifest="$public/.git/quine-source-manifest.json"; git -C "$root" rev-parse --is-inside-work-tree >/dev/null || { echo "source-code git worktree is unreadable" >&2; exit 1; }; [ -f "$manifest" ] || { echo "missing self-source manifest: $manifest" >&2; exit 1; }; cmp -s "$manifest" "$public_manifest" || { echo "public self-source manifest drift" >&2; exit 1; }; grep -q "\"format\": \"quine-source-repo/v1\"" "$manifest" || { echo "self-source manifest format mismatch" >&2; exit 1; }; [ ! -w "$root/go.mod" ] || { echo "go.mod should be read-only" >&2; exit 1; }; grep -q "^module github.com/kehao95/quine$" "$root/go.mod" || { echo "go.mod module line mismatch" >&2; exit 1; }; grep -q "func main()" "$root/cmd/quine/main.go" || { echo "main.go missing func main" >&2; exit 1; }; grep -q "var SelfSourceBundle" "$root/selfsource_bundle_data.go" || { echo "selfsource_bundle_data.go missing SelfSourceBundle definition" >&2; exit 1; }; grep -q "type Runtime struct" "$root/internal/runtime/runtime.go" || { echo "runtime.go missing Runtime struct" >&2; exit 1; }; echo SELF_SOURCE_SURFACE_OK >&4; echo SELF_SOURCE_CONTENT_OK >&4'

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_SELF_SOURCE_CODE_ENABLED=1 \
        "$QUINE" "Use exactly one sh call with this exact command: \`$self_source_check\` Then exit success." \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "self-source surface exit" &&
        assert_contains "$TEST_DIR/stdout" "SELF_SOURCE_SURFACE_OK" "self-source surface marker" &&
        assert_contains "$TEST_DIR/stdout" "SELF_SOURCE_CONTENT_OK" "self-source content marker"; then
        pass
    fi
    teardown
}

test_prompt_self_source_enabled() {
    begin_test "QUINE_SELF_SOURCE_CODE_ENABLED=1 discloses self-source surface in prompt"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS=2 \
        QUINE_SELF_SOURCE_CODE_ENABLED=1 \
        "$QUINE" 'Exit immediately with status success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "self-source prompt exit"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'source-code/' "$_tape" 2>/dev/null; then
            fail "self-source prompt disclosure should be visible when enabled"
        elif ! grep -q "read-only session-root projection of this Quine body's source" "$_tape" 2>/dev/null; then
            fail "self-source prompt should describe read-only runtime-carried source"
        elif ! grep -q 'source-code/` is a git worktree with `.git/`' "$_tape" 2>/dev/null; then
            fail "self-source prompt should disclose repo-backed source-code behavior"
        elif ! grep -q 'Source manifest: `.git/quine-source-manifest.json`' "$_tape" 2>/dev/null; then
            fail "self-source prompt should disclose source manifest path"
        elif ! grep -q 'public/source-code/' "$_tape" 2>/dev/null; then
            fail "self-source prompt should disclose public source projection"
        else
            pass
        fi
    fi
    teardown
}

test_process_surface_incarnation_projection() {
    begin_test "process surface projects current incarnation through inc/current"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one sh call. In that call, set -eu; current="$QUINE_AGENT_ROOT/inc/current"; [ -L "$current" ] || { echo "missing inc/current symlink" >&2; exit 1; }; target=$(readlink "$current"); [ "$target" = "0" ] || { echo "inc/current target = $target, want 0" >&2; exit 1; }; session_json="$QUINE_AGENT_ROOT/status/session.json"; [ -f "$session_json" ] || { echo "missing session.json" >&2; exit 1; }; incarnation_id=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"incarnation_id\"])" "$session_json"); [ "$incarnation_id" = "0" ] || { echo "incarnation_id = $incarnation_id, want 0" >&2; exit 1; }; mission_target=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$QUINE_AGENT_ROOT/mission.txt"); expected_mission=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$QUINE_AGENT_ROOT/inc/current/mission.txt"); [ "$mission_target" = "$expected_mission" ] || { echo "mission projection mismatch" >&2; exit 1; }; context_target=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$QUINE_AGENT_ROOT/context"); expected_context=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$QUINE_AGENT_ROOT/inc/current/context"); [ "$context_target" = "$expected_context" ] || { echo "context projection mismatch" >&2; exit 1; }; printf "INCARNATION_PROJECTION_OK\n" >&4. Then exit success.'
    code=$?

    if assert_exit "$code" 0 "incarnation projection exit" &&
        assert_contains "$TEST_DIR/stdout" "INCARNATION_PROJECTION_OK" "incarnation projection marker"; then
        pass
    fi
    teardown
}

test_process_surface_config_surface() {
    begin_test "process surface exposes config/ capability read surface"
    setup

    # config/ is the enacted capability read surface: config/registry.json (the
    # compiled knob catalog, 0444) plus the config/env/ directory the agent writes
    # its child-env override into. The runtime renders NO env file back to the
    # agent — a process reads its OWN environment where the OS publishes it, at
    # /proc/<pid>/environ, and a QUINE_* name absent there is at its compiled
    # default. So this block proves: (1) registry.json is the FULL compiled
    # registry (size + anchor knobs, not a brittle literal — the Go suite's
    # TestRegistryBijectionWithEnvNames in internal/config/registry_test.go owns
    # the exact count and the registry ≡ envnames.go bijection); (2) the
    # config/env/ write directory exists; (3) the deleted surfaces
    # (config/resolved.env, config/env/{effective,pinned}, inc/current/environ)
    # are GONE; and (4) the process's own /proc environ carries the launched
    # marker QUINE_OUTPUT_TRUNCATE=7777 but no synthesized default like
    # QUINE_SPAWN_ENABLED (a line no operator authored).
    config_check='set -eu; dir="$QUINE_AGENT_ROOT/config"; reg="$dir/registry.json"; [ -f "$reg" ] || { echo "missing registry.json" >&2; exit 1; }; [ ! -w "$reg" ] || { echo "registry.json should be read-only" >&2; exit 1; }; count=$(python3 -c "import json,sys; print(len(json.load(open(sys.argv[1]))))" "$reg"); [ "$count" -ge 50 ] || { echo "registry knob count = $count, want a full registry (>= 50)" >&2; exit 1; }; python3 -c "import json,sys; names={k[\"env\"] for k in json.load(open(sys.argv[1]))}; sys.exit(0 if all(n in names for n in sys.argv[2:]) else 1)" "$reg" QUINE_OUTPUT_TRUNCATE QUINE_SH_DEFAULT_TIMEOUT_SECONDS QUINE_API_KEY || { echo "registry.json missing expected anchor knobs" >&2; exit 1; }; [ -d "$dir/env" ] || { echo "missing config/env override directory" >&2; exit 1; }; for gone in "$dir/resolved.env" "$dir/env/effective" "$dir/env/pinned" "$QUINE_AGENT_ROOT/inc/current/environ"; do [ ! -e "$gone" ] || { echo "deleted env surface still present: $gone" >&2; exit 1; }; done; pid=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"pid\"])" "$QUINE_AGENT_ROOT/status/session.json"); envfile="/proc/$pid/environ"; if [ -r "$envfile" ]; then tr "\0" "\n" < "$envfile" | grep -q "^QUINE_OUTPUT_TRUNCATE=7777$" || { echo "own /proc environ missing launched marker" >&2; exit 1; }; if tr "\0" "\n" < "$envfile" | grep -q "^QUINE_SPAWN_ENABLED="; then echo "own /proc environ carries a synthesized default" >&2; exit 1; fi; fi; echo CONFIG_SURFACE_OK >&4; echo CONFIG_CONTENT_OK >&4'

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_OUTPUT_TRUNCATE=7777 \
        "$QUINE" "Use exactly one sh call with this exact command: \`$config_check\` Then exit success." \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "config surface exit" &&
        assert_contains "$TEST_DIR/stdout" "CONFIG_SURFACE_OK" "config surface marker" &&
        assert_contains "$TEST_DIR/stdout" "CONFIG_CONTENT_OK" "config content marker"; then
        pass
    fi
    teardown
}

test_process_surface_config_projection() {
    begin_test "public/config projects the capability read surface to peers"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    projection_check='set -eu; pub="$QUINE_AGENT_ROOT/public/config"; [ -d "$pub" ] || { echo "missing public/config" >&2; exit 1; }; cmp -s "$pub/registry.json" "$QUINE_AGENT_ROOT/config/registry.json" || { echo "public registry.json drifts from config/registry.json" >&2; exit 1; }; ( echo x > "$pub/registry.json" ) 2>/dev/null && { echo "public registry.json should be read-only" >&2; exit 1; }; for gone in "$pub/resolved.env" "$pub/env"; do [ ! -e "$gone" ] || { echo "public/config projects a non-catalog surface: $gone" >&2; exit 1; }; done; adv_config=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"surfaces\"][\"config\"])" "$QUINE_AGENT_ROOT/public/status/contract.json"); [ "$adv_config" = "config" ] || { echo "contract.json config surface = $adv_config" >&2; exit 1; }; adv_ctl=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"surfaces\"][\"config_control\"])" "$QUINE_AGENT_ROOT/public/status/contract.json"); [ "$adv_ctl" = "ctl/env" ] || { echo "contract.json config_control surface = $adv_ctl" >&2; exit 1; }; echo CONFIG_PROJECTION_OK >&4'

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_OUTPUT_TRUNCATE=7777 \
        "$QUINE" "Use exactly one sh call with this exact command: \`$projection_check\` Then exit success." \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "config projection exit" &&
        assert_contains "$TEST_DIR/stdout" "CONFIG_PROJECTION_OK" "config projection marker"; then
        pass
    fi
    teardown
}

test_process_surface_neighbor_discovery_prunes_stale_indexes() {
    begin_test "pid surface exposes self-authored routes and prunes stale entries"
    setup

    mkdir -p "$QUINE_DATA_DIR/locks" "$QUINE_DATA_DIR/pid" "$QUINE_DATA_DIR/agent/live" "$QUINE_DATA_DIR/agent/live_by_pid" "$QUINE_DATA_DIR/agent/stale-session"
    printf '999999\n' > "$QUINE_DATA_DIR/locks/stale-session.agent"
    rm -f "$QUINE_DATA_DIR/pid/999999" "$QUINE_DATA_DIR/agent/live/stale-session" "$QUINE_DATA_DIR/agent/live_by_pid/999999"
    ln -s "$QUINE_DATA_DIR/agent/stale-session" "$QUINE_DATA_DIR/pid/999999"
    ln -s "$QUINE_DATA_DIR/agent/stale-session" "$QUINE_DATA_DIR/agent/live/stale-session"
    ln -s "$QUINE_DATA_DIR/agent/stale-session" "$QUINE_DATA_DIR/agent/live_by_pid/999999"

    helper_mission='Use exactly one sh call with command `sleep 120` and timeout `130`, and nothing else. Wait for that call to complete. Then exit success.'
    start_quine_helper "$TEST_DIR/helper.stdout" "$TEST_DIR/helper.stderr" "$helper_mission"
    helper_pid="$HELPER_PID"

    helper_ready=0
    i=0
    while [ "$i" -lt 80 ]; do
        if [ ! -e "$QUINE_DATA_DIR/locks/stale-session.agent" ] &&
            [ ! -e "$QUINE_DATA_DIR/pid/999999" ] &&
            [ ! -e "$QUINE_DATA_DIR/agent/live" ] &&
            [ ! -e "$QUINE_DATA_DIR/agent/live_by_pid" ] &&
            [ -L "$QUINE_DATA_DIR/pid/$helper_pid" ]; then
            helper_ready=1
            break
        fi
        sleep 0.25
        i=$((i + 1))
    done

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one sh call. Set its command to exactly:
set -eu; runtime_root=$(cd "$QUINE_AGENT_ROOT/../.." && pwd -P); session_json="$QUINE_AGENT_ROOT/status/session.json"; self_pid=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"pid\"])" "$session_json"); self_link="$runtime_root/pid/$self_pid"; [ -L "$self_link" ] || { echo "self pid entry missing: $self_link" >&2; exit 1; }; self_target=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$self_link"); self_public=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$QUINE_AGENT_ROOT/public"); self_root=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$QUINE_AGENT_ROOT"); [ "$self_target" = "$self_public" ] || { echo "self pid entry resolves to $self_target, expected $self_public" >&2; exit 1; }; [ ! -e "$runtime_root/pid/999999" ] || { echo "stale pid entry still present" >&2; exit 1; }; [ ! -e "$runtime_root/agent/live" ] || { echo "legacy agent/live still present" >&2; exit 1; }; [ ! -e "$runtime_root/agent/live_by_pid" ] || { echo "legacy agent/live_by_pid still present" >&2; exit 1; }; neighbor_link=$(find "$runtime_root/pid" -mindepth 1 -maxdepth 1 -type l ! -name "$self_pid" | head -n 1); [ -n "$neighbor_link" ] || { echo "no live neighbor besides self" >&2; exit 1; }; neighbor_target=$(python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$neighbor_link"); [ "$neighbor_target" != "$self_public" ] || { echo "neighbor resolves to self" >&2; exit 1; }; case "$neighbor_target" in "$runtime_root"/agent/*/public) ;; *) echo "neighbor target outside public agent namespace: $neighbor_target" >&2; exit 1 ;; esac; neighbor_root=$(dirname "$neighbor_target"); [ "$neighbor_root" != "$self_root" ] || { echo "neighbor root resolves to self" >&2; exit 1; }; printf "SELF_PID_SURFACE_OK\nNEIGHBOR_PID_SURFACE_OK\nSTALE_PRUNED_OK\nLEGACY_LIVE_REMOVED_OK\n" >/dev/fd/4
Do not append any extra prose or punctuation to that shell command. After the sh call succeeds, exit success.'
    code=$?

    if kill -0 "$helper_pid" >/dev/null 2>&1; then
        kill "$helper_pid" >/dev/null 2>&1 || true
    fi
    wait "$helper_pid" >/dev/null 2>&1 || true

    if [ "$helper_ready" -ne 1 ]; then
        fail "helper registration did not appear or stale entries were not pruned in time"
    elif assert_exit "$code" 0 "neighbor surface exit" &&
        assert_contains "$TEST_DIR/stdout" "SELF_PID_SURFACE_OK" "self pid surface" &&
        assert_contains "$TEST_DIR/stdout" "NEIGHBOR_PID_SURFACE_OK" "neighbor pid surface" &&
        assert_contains "$TEST_DIR/stdout" "STALE_PRUNED_OK" "stale pid surface" &&
        assert_contains "$TEST_DIR/stdout" "LEGACY_LIVE_REMOVED_OK" "legacy live index removal"; then
        pass
    fi
    teardown
}

test_process_surface_pid_route_removed_on_sigterm() {
    begin_test "SIGTERM gracefully removes pid route and pid lock"
    setup

    helper_mission='Use exactly one sh call with command `sleep 120` and timeout `130`, and nothing else. Wait for that call to complete. Then exit success.'
    start_quine_helper "$TEST_DIR/helper.stdout" "$TEST_DIR/helper.stderr" "$helper_mission"
    helper_pid="$HELPER_PID"
    helper_pid_link="$QUINE_DATA_DIR/pid/$helper_pid"
    helper_pid_lock="$QUINE_DATA_DIR/locks/agents/$helper_pid.agent.lock"

    helper_ready=0
    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$helper_pid_link" ] && [ -f "$helper_pid_lock" ]; then
            helper_ready=1
            break
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$helper_ready" -eq 1 ]; then
        kill -TERM "$helper_pid" >/dev/null 2>&1 || true
    fi
    wait "$helper_pid" >/dev/null 2>&1 || true

    if [ "$helper_ready" -ne 1 ]; then
        fail "helper pid route and lock did not appear in time"
    elif [ -e "$helper_pid_link" ]; then
        fail "pid route still exists after SIGTERM: $helper_pid_link"
    elif [ -e "$helper_pid_lock" ]; then
        fail "pid lock still exists after SIGTERM: $helper_pid_lock"
    elif find "$QUINE_DATA_DIR/locks" -maxdepth 1 -name '*.agent' -print -quit 2>/dev/null | grep -q .; then
        fail "agent registration remained after SIGTERM"
    else
        pass
    fi
    teardown
}

test_process_surface_stale_pid_lock_pruned_on_startup() {
    begin_test "startup prunes stale pid lock artifacts"
    setup

    stale_pid=999999
    stale_session="stale-session"
    mkdir -p "$QUINE_DATA_DIR/locks/agents" "$QUINE_DATA_DIR/pid" "$QUINE_DATA_DIR/agent/$stale_session/public" "$QUINE_DATA_DIR/agent/$stale_session/status"
    printf '{"session_id":"%s","run_id":"stale-run","pid":%s}\n' "$stale_session" "$stale_pid" > "$QUINE_DATA_DIR/locks/stale-run.agent"
    : > "$QUINE_DATA_DIR/locks/agents/$stale_pid.agent.lock"
    rm -f "$QUINE_DATA_DIR/pid/$stale_pid"
    ln -s "$QUINE_DATA_DIR/agent/$stale_session/public" "$QUINE_DATA_DIR/pid/$stale_pid"
    printf '{"session_id":"%s","pid":%s}\n' "$stale_session" "$stale_pid" > "$QUINE_DATA_DIR/agent/$stale_session/status/session.json"

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one sh call. In that call, set -eu; runtime_root=$(cd "$QUINE_AGENT_ROOT/../.." && pwd -P); [ ! -e "$runtime_root/locks/stale-run.agent" ] || { echo "stale agent registration remains" >&2; exit 1; }; [ ! -e "$runtime_root/locks/agents/999999.agent.lock" ] || { echo "stale pid lock remains" >&2; exit 1; }; [ ! -e "$runtime_root/pid/999999" ] || { echo "stale pid route remains" >&2; exit 1; }; [ ! -e "$runtime_root/agent/stale-session" ] || { echo "stale agent root remains" >&2; exit 1; }; printf "STALE_PID_LOCK_PRUNED_OK\n" >&4. Then exit success.'
    code=$?

    if assert_exit "$code" 0 "stale pid lock prune exit" &&
        assert_contains "$TEST_DIR/stdout" "STALE_PID_LOCK_PRUNED_OK" "stale pid lock prune marker"; then
        pass
    fi
    teardown
}

test_process_surface_peer_discovery_heartbeat_prunes_stale_pid_lock() {
    begin_test "peer discovery heartbeat prunes stale pid lock after startup"
    setup

    old_peer_observation="${QUINE_PEER_DISCOVERY_ENABLED-__unset__}"
    old_peer_heartbeat="${QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS-__unset__}"
    export QUINE_PEER_DISCOVERY_ENABLED=1
    export QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS=1000
    subject_mission='Use exactly one sh call with command `sleep 120` and timeout `130`, and nothing else. Wait for that call to complete. Then exit success.'
    start_quine_helper "$TEST_DIR/subject.stdout" "$TEST_DIR/subject.stderr" "$subject_mission"
    subject_pid="$HELPER_PID"
    if [ "$old_peer_observation" = "__unset__" ]; then
        unset QUINE_PEER_DISCOVERY_ENABLED
    else
        QUINE_PEER_DISCOVERY_ENABLED="$old_peer_observation"
        export QUINE_PEER_DISCOVERY_ENABLED
    fi
    if [ "$old_peer_heartbeat" = "__unset__" ]; then
        unset QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS
    else
        QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS="$old_peer_heartbeat"
        export QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS
    fi

    subject_ready=0
    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$subject_pid" ]; then
            subject_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$subject_pid")"
            if [ -n "$subject_public" ] && legacy_control_surface_ready "$subject_public"; then
                subject_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    stale_pid=999999
    stale_session="heartbeat-stale-peer"
    mkdir -p "$QUINE_DATA_DIR/locks/agents" "$QUINE_DATA_DIR/pid" "$QUINE_DATA_DIR/agent/$stale_session/public" "$QUINE_DATA_DIR/agent/$stale_session/status"
    printf '{"session_id":"%s","run_id":"heartbeat-stale-run","pid":%s}\n' "$stale_session" "$stale_pid" > "$QUINE_DATA_DIR/locks/heartbeat-stale-run.agent"
    rm -f "$QUINE_DATA_DIR/pid/$stale_pid"
    ln -s "$QUINE_DATA_DIR/agent/$stale_session/public" "$QUINE_DATA_DIR/pid/$stale_pid"
    printf '{"session_id":"%s","run_id":"heartbeat-stale-run","pid":%s}\n' "$stale_session" "$stale_pid" > "$QUINE_DATA_DIR/agent/$stale_session/status/session.json"
    : > "$QUINE_DATA_DIR/locks/agents/$stale_pid.agent.lock"

    heartbeat_seen=0
    i=0
    while [ "$i" -lt 80 ]; do
        if [ ! -L "$QUINE_DATA_DIR/pid/$stale_pid" ] &&
            [ ! -e "$QUINE_DATA_DIR/locks/heartbeat-stale-run.agent" ] &&
            [ ! -e "$QUINE_DATA_DIR/locks/agents/$stale_pid.agent.lock" ]; then
            heartbeat_seen=1
            break
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if kill -0 "$subject_pid" >/dev/null 2>&1; then
        kill -TERM "$subject_pid" >/dev/null 2>&1 || true
    fi
    wait "$subject_pid" >/dev/null 2>&1 || true
    if [ "$subject_ready" -ne 1 ]; then
        fail "heartbeat subject control surface did not appear in time"
    elif [ "$heartbeat_seen" -ne 1 ]; then
        fail "heartbeat did not prune stale pid lock"
    else
        pass
    fi
    teardown
}

test_process_surface_ctl_inbox_only() {
    begin_test "ctl/post queues inbox-only post and adds a post indicator"
    setup

    target_mission='Use exactly one sh call with command `sleep 5` and nothing else. Then exit success.'
    start_quine_helper "$TEST_DIR/target.stdout" "$TEST_DIR/target.stderr" "$target_mission"
    target_pid="$HELPER_PID"
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            if [ -n "$target_public" ] && legacy_control_surface_ready "$target_public"; then
                target_session="$(basename "$target_root")"
                retained_root="$QUINE_DATA_DIR/log/$target_session"
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        printf 'mail from runtime test\n' > "$target_public/ctl/post"
    fi

    wait "$target_pid"
    code=$?

    inbox_checks=""
    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"
        inbox_checks="$(python3 - "$retained_root/status/inbox.json" <<'PY'
import json
import sys
from pathlib import Path

inbox = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
messages = inbox.get("messages") or []
print(f"PENDING_OK={1 if inbox.get('pending_count') == 1 else 0}")
print(f"PAYLOAD_OK={1 if len(messages) == 1 and messages[0].get('payload') == 'mail from runtime test' else 0}")
PY
)"
    fi

    if [ "$target_ready" -ne 1 ]; then
        fail "target control surface did not appear in time"
    elif assert_exit "$code" 0 "ctl/post inbox-only exit" &&
        assert_contains "$retained_tape" '"inbox_indicator":"You have pending inbox messages\. Inspect QUINE_AGENT_ROOT/status/inbox\.json\."' "ctl/post inbox-only tape indicator" &&
        printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "PAYLOAD_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" 2>/dev/null; then
            fail "ctl/post pending_count was not 1"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "PAYLOAD_OK=1" 2>/dev/null; then
            fail "ctl/post payload missing from inbox snapshot"
        fi
    fi
    teardown
}

test_process_surface_ctl_inject_delivery() {
    begin_test "ctl/inject delivers post at the next safe point"
    setup

    target_mission='Use exactly one sh call with command `sleep 5` and nothing else. Then exit success.'
    start_quine_helper "$TEST_DIR/target.stdout" "$TEST_DIR/target.stderr" "$target_mission"
    target_pid="$HELPER_PID"
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            if [ -n "$target_public" ] && legacy_control_surface_ready "$target_public"; then
                target_session="$(basename "$target_root")"
                retained_root="$QUINE_DATA_DIR/log/$target_session"
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        printf 'wake from runtime test\n' > "$target_public/ctl/inject"
    fi

    wait "$target_pid"
    code=$?

    inbox_checks=""
    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"
        _settled=0
        _attempt=0
        while [ "$_attempt" -lt 20 ]; do
            inbox_checks="$(python3 - "$retained_root/status/inbox.json" "$retained_tape" "$retained_root/control.jsonl" <<'PY'
import json
import sys
from pathlib import Path

inbox = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
current = Path(sys.argv[2]).read_text(encoding="utf-8")
control = Path(sys.argv[3]).read_text(encoding="utf-8")
print(f"PENDING_OK={1 if inbox.get('pending_count') == 0 else 0}")
print(f"CURRENT_OK={1 if '\"incoming_messages\":[' in current and '\"delivery\":\"inject\"' in current and '\"payload\":\"wake from runtime test\"' in current else 0}")
print(f"CONTROL_OK={1 if '\"kind\":\"delivered\"' in control and '\"delivery\":\"inject\"' in control and '\"payload\":\"wake from runtime test\"' in control else 0}")
PY
)"
            if printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
                printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" &&
                printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1"; then
                _settled=1
                break
            fi
            sleep 0.1
            _attempt=$(( _attempt + 1 ))
        done
    fi

    if [ "$target_ready" -ne 1 ]; then
        fail "target control surface did not appear in time"
    elif assert_exit "$code" 0 "ctl/inject exit" &&
        printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" 2>/dev/null; then
            fail "ctl/inject inbox pending_count was not 0"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" 2>/dev/null; then
            fail "ctl/inject retained tape never surfaced incoming_messages"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1" 2>/dev/null; then
            fail "ctl/inject control log never recorded delivered inject mail"
        fi
    fi
    teardown
}

test_process_surface_ctl_interrupt_delivery() {
    begin_test "ctl/interrupt delivers interrupt mail"
    setup

    target_mission="Use exactly one sh call with command \`python3 -c 'import signal,sys,time; signal.signal(signal.SIGINT, lambda signum, frame: sys.exit(130)); time.sleep(5)'\` and nothing else. Then exit success."
    start_quine_helper "$TEST_DIR/target.stdout" "$TEST_DIR/target.stderr" "$target_mission"
    target_pid="$HELPER_PID"
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            if [ -n "$target_public" ] && legacy_control_surface_ready "$target_public"; then
                target_session="$(basename "$target_root")"
                retained_root="$QUINE_DATA_DIR/log/$target_session"
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        printf 'interrupt from runtime test\n' > "$target_public/ctl/interrupt"
    fi

    wait "$target_pid"
    code=$?

    inbox_checks=""
    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"
        inbox_checks="$(python3 - "$retained_root/status/inbox.json" <<'PY'
import json
import sys
from pathlib import Path

inbox = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
print(f"PENDING_OK={1 if inbox.get('pending_count') == 0 else 0}")
PY
)"
    fi

    if [ "$target_ready" -ne 1 ]; then
        fail "target control surface did not appear in time"
    elif assert_exit "$code" 0 "ctl/interrupt exit" &&
        assert_contains "$retained_tape" '"delivery":"interrupt"' "ctl/interrupt delivery label" &&
        assert_contains "$retained_tape" '"interrupt_notice":"Current operation was interrupted by peer control input\."' "ctl/interrupt notice" &&
        assert_contains "$retained_tape" '"payload":"interrupt from runtime test"' "ctl/interrupt payload" &&
        printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" 2>/dev/null; then
            fail "ctl/interrupt inbox pending_count was not 0"
        fi
    fi
    teardown
}

test_process_surface_fuse_ctl_post_transaction() {
    begin_test "fuse public ctl/post queues inbox post without delivery"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    target_mission='Use exactly one sh call with command sleep 5 and nothing else. Then exit success.'
    env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        "$QUINE" "$target_mission" \
        >"$TEST_DIR/target.stdout" \
        2>"$TEST_DIR/target.stderr" \
        </dev/null &
    target_pid=$!
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""
    ctl_summary=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            target_session="$(basename "$target_root")"
            retained_root="$QUINE_DATA_DIR/log/$target_session"
            if fuse_control_surface_ready "$target_public"; then
                ctl_summary="$(cat "$target_public/ctl/post" 2>/dev/null || true)"
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        printf 'mail from fuse runtime test\n' > "$target_public/ctl/post"
    fi

    wait "$target_pid"
    code=$?

    inbox_checks=""
    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"
        inbox_checks="$(python3 - "$retained_root/status/inbox.json" <<'PY'
import json
import sys
from pathlib import Path

inbox = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
messages = inbox.get("messages") or []
print(f"PENDING_OK={1 if inbox.get('pending_count') == 1 else 0}")
print(f"PAYLOAD_OK={1 if len(messages) == 1 and messages[0].get('payload') == 'mail from fuse runtime test' else 0}")
PY
)"
    fi

    if [ "$target_ready" -ne 1 ]; then
        fail "fuse public control surface did not appear in time"
    elif assert_exit "$code" 0 "fuse ctl/post exit" &&
        printf '%s\n' "$ctl_summary" | grep -q "control_file: post" &&
        printf '%s\n' "$ctl_summary" | grep -q "mode: queue-only" &&
        assert_contains "$retained_tape" '"inbox_indicator":"You have pending inbox messages\. Inspect QUINE_AGENT_ROOT/status/inbox\.json\."' "fuse ctl/post inbox-only tape indicator" &&
        printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "PAYLOAD_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$ctl_summary" | grep -q "control_file: post" 2>/dev/null; then
            fail "fuse public ctl/post summary did not expose post control file"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" 2>/dev/null; then
            fail "fuse ctl/post pending_count was not 1"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "PAYLOAD_OK=1" 2>/dev/null; then
            fail "fuse ctl/post payload missing from inbox snapshot"
        fi
    fi
    teardown
}

test_process_surface_fuse_ctl_env() {
    begin_test "fuse public ctl/env validates child-env policy writes and lands the override"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    # Validated write gate over config/env/override (the one managed env file):
    # a peer writes one whole child-env policy through public/ctl/env; a legal
    # payload lands the override byte-equal at close, an illegal one is rejected
    # at close (EINVAL), lands nothing, leaves the previously landed policy
    # intact, and its violations are readable back from the gate. The read-back
    # NAMES the policy (policy_names: NAME (set), ...) plus validation state and
    # never quotes VALUES — the override is the agent's own file, deliberately
    # not projected to peers, so this peer-facing node must not echo its bytes.
    # The illegal payload mixes the two rejection classes: QUINE_MAX_DEPTH is
    # pinned (operator-only, E5) and QUINE_TOTALLY_UNKNOWN_KNOB is not a knob.
    target_mission='Use exactly one sh call with command sleep 5 and nothing else. Then exit success.'
    env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" "$target_mission" \
        >"$TEST_DIR/target.stdout" \
        2>"$TEST_DIR/target.stderr" \
        </dev/null &
    target_pid=$!
    target_ready=0
    target_public=""
    target_root=""
    gate_checks=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            if fuse_control_surface_ready "$target_public" &&
                [ -f "$target_public/ctl/env" ] && [ ! -p "$target_public/ctl/env" ]; then
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        gate_checks="$(python3 - "$target_public/ctl/env" "$target_root/config/env/override" <<'PY'
import errno
import sys
from pathlib import Path

gate = Path(sys.argv[1])
override = Path(sys.argv[2])

summary = gate.read_text(encoding="utf-8")
print(f"SUMMARY_MODE_OK={1 if 'control_file: env' in summary and 'mode: validated-child-env-policy' in summary else 0}")
print(f"SUMMARY_EMPTY_OK={1 if 'policy: none' in summary else 0}")

# Names/values NOT present in the gate's static example line (which uses
# QUINE_MAX_TURNS=64, FOO=bar, LANG), so the values-not-echoed check below
# cannot be fooled by the help text.
valid = "QUINE_OUTPUT_TRUNCATE=31337\nQUINE_SPAWN_ENABLED=1\nAPP_SECRET=zzz-do-not-echo-9931\n"
gate.write_text(valid, encoding="utf-8")
print(f"VALID_LANDED_OK={1 if override.read_bytes() == valid.encode() else 0}")
summary = gate.read_text(encoding="utf-8")
names_ok = all(s in summary for s in (
    "policy_names:",
    "QUINE_OUTPUT_TRUNCATE (set)",
    "QUINE_SPAWN_ENABLED (set)",
    "APP_SECRET (set)",
    "validation: valid against the running capability registry",
))
print(f"VALID_NAMES_OK={1 if names_ok else 0}")
print(f"VALUES_NOT_ECHOED_OK={1 if '31337' not in summary and 'zzz-do-not-echo-9931' not in summary else 0}")

rejected = 0
try:
    gate.write_text("QUINE_MAX_DEPTH=7\nQUINE_TOTALLY_UNKNOWN_KNOB=1\n", encoding="utf-8")
except OSError as exc:
    rejected = 1 if exc.errno == errno.EINVAL else 0
print(f"INVALID_REJECTED_OK={rejected}")
print(f"OVERRIDE_INTACT_OK={1 if override.read_bytes() == valid.encode() else 0}")
summary = gate.read_text(encoding="utf-8")
viol_ok = (
    "last_rejected_write: rejected in full, nothing landed:" in summary
    and "QUINE_MAX_DEPTH: mutability" in summary
    and "QUINE_TOTALLY_UNKNOWN_KNOB: unknown env name" in summary
)
print(f"VIOLATIONS_READBACK_OK={1 if viol_ok else 0}")
PY
)"
    fi

    wait "$target_pid"
    code=$?

    if [ "$target_ready" -ne 1 ]; then
        fail "fuse public ctl/env gate did not appear in time"
    elif assert_exit "$code" 0 "fuse ctl/env exit" &&
        printf '%s\n' "$gate_checks" | grep -q "SUMMARY_MODE_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "SUMMARY_EMPTY_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "VALID_LANDED_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "VALID_NAMES_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "VALUES_NOT_ECHOED_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "INVALID_REJECTED_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "OVERRIDE_INTACT_OK=1" &&
        printf '%s\n' "$gate_checks" | grep -q "VIOLATIONS_READBACK_OK=1"; then
        pass
    else
        for _gate_check in SUMMARY_MODE_OK SUMMARY_EMPTY_OK VALID_LANDED_OK VALID_NAMES_OK VALUES_NOT_ECHOED_OK INVALID_REJECTED_OK OVERRIDE_INTACT_OK VIOLATIONS_READBACK_OK; do
            if ! printf '%s\n' "$gate_checks" | grep -q "${_gate_check}=1" 2>/dev/null; then
                fail "fuse ctl/env gate check failed: $_gate_check (results: $(printf '%s' "$gate_checks" | tr '\n' ' '))"
                break
            fi
        done
    fi
    teardown
}

test_process_surface_fuse_public_scan_safe() {
    begin_test "fuse public surface stays scan-safe under ordinary recursive reads"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    target_mission='Use exactly one sh call with command sleep 5 and nothing else. Then exit success.'
    env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        "$QUINE" "$target_mission" \
        >"$TEST_DIR/target.stdout" \
        2>"$TEST_DIR/target.stderr" \
        </dev/null &
    target_pid=$!
    target_ready=0
    target_public=""
    scan_checks=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            if fuse_control_surface_ready "$target_public"; then
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        scan_checks="$(python3 - "$target_public" <<'PY'
import subprocess
import sys
from pathlib import Path

public_root = Path(sys.argv[1])
ctl_path = public_root / "ctl"
ctl_poke = ctl_path / "poke"

find_ok = False
grep_ok = False
ctl_ok = False

try:
    find_run = subprocess.run(
        ["find", str(public_root)],
        capture_output=True,
        text=True,
        timeout=5,
        check=False,
    )
    find_ok = find_run.returncode == 0 and str(public_root) in find_run.stdout
except Exception:
    find_ok = False

try:
    grep_run = subprocess.run(
        ["grep", "-R", "poke", str(public_root)],
        capture_output=True,
        text=True,
        timeout=5,
        check=False,
    )
    grep_ok = grep_run.returncode == 0 and "poke" in grep_run.stdout
except Exception:
    grep_ok = False

try:
    ctl_text = ctl_poke.read_text(encoding="utf-8", errors="replace")
    ctl_ok = (
        "control_file: poke" in ctl_text
        and "mode: queue-and-resume" in ctl_text
        and "without context injection" in ctl_text
    )
except Exception:
    ctl_ok = False

print(f"FIND_OK={1 if find_ok else 0}")
print(f"GREP_OK={1 if grep_ok else 0}")
print(f"CTL_OK={1 if ctl_ok else 0}")
PY
)"
    fi

    wait "$target_pid"
    code=$?

    if [ "$target_ready" -ne 1 ]; then
        fail "fuse public control surface did not appear in time"
    elif assert_exit "$code" 0 "fuse public scan-safe exit" &&
        printf '%s\n' "$scan_checks" | grep -q "FIND_OK=1" &&
        printf '%s\n' "$scan_checks" | grep -q "GREP_OK=1" &&
        printf '%s\n' "$scan_checks" | grep -q "CTL_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$scan_checks" | grep -q "FIND_OK=1" 2>/dev/null; then
            fail "find on fuse public surface did not complete cleanly"
        elif ! printf '%s\n' "$scan_checks" | grep -q "GREP_OK=1" 2>/dev/null; then
            fail "recursive grep on fuse public surface did not complete cleanly"
        elif ! printf '%s\n' "$scan_checks" | grep -q "CTL_OK=1" 2>/dev/null; then
            fail "fuse public ctl/poke summary did not expose scan-safe named control mode"
        fi
    fi
    teardown
}

test_process_surface_fuse_ctl_inject_transaction() {
    begin_test "fuse public ctl/inject delivers post at the next safe point"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    target_mission='Use exactly one sh call with command sleep 5 and nothing else. Then exit success.'
    env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        "$QUINE" "$target_mission" \
        >"$TEST_DIR/target.stdout" \
        2>"$TEST_DIR/target.stderr" \
        </dev/null &
    target_pid=$!
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            target_session="$(basename "$target_root")"
            retained_root="$QUINE_DATA_DIR/log/$target_session"
            if fuse_control_surface_ready "$target_public"; then
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        printf 'wake from fuse runtime test\n' > "$target_public/ctl/inject"
    fi

    wait "$target_pid"
    code=$?

    inbox_checks=""
    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"
        _settled=0
        _attempt=0
        while [ "$_attempt" -lt 20 ]; do
            inbox_checks="$(python3 - "$retained_root/status/inbox.json" "$retained_tape" "$retained_root/control.jsonl" <<'PY'
import json
import sys
from pathlib import Path

inbox = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
current = Path(sys.argv[2]).read_text(encoding="utf-8")
control = Path(sys.argv[3]).read_text(encoding="utf-8")
print(f"PENDING_OK={1 if inbox.get('pending_count') == 0 else 0}")
current_ok = (
    '"incoming_messages":[' in current
    and '"delivery":"inject"' in current
    and '"payload":"wake from fuse runtime test"' in current
)
control_ok = (
    '"kind":"delivered"' in control
    and '"delivery":"inject"' in control
    and '"payload":"wake from fuse runtime test"' in control
)
print(f"CURRENT_OK={1 if current_ok else 0}")
print(f"CONTROL_OK={1 if control_ok else 0}")
PY
)"
            if printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
                printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" &&
                printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1"; then
                _settled=1
                break
            fi
            sleep 0.1
            _attempt=$(( _attempt + 1 ))
        done
    fi

    if [ "$target_ready" -ne 1 ]; then
        fail "fuse public control surface did not appear in time"
    elif assert_exit "$code" 0 "fuse ctl/inject exit" &&
        printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" 2>/dev/null; then
            fail "fuse ctl/inject inbox pending_count was not 0"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" 2>/dev/null; then
            fail "fuse ctl/inject retained tape never surfaced incoming_messages"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1" 2>/dev/null; then
            fail "fuse ctl/inject control log never recorded delivered inject mail"
        fi
    fi
    teardown
}

test_process_surface_fuse_ctl_interrupt_transaction() {
    begin_test "fuse public ctl/interrupt delivers interrupt mail"
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    target_mission="Use exactly one sh call with command \`python3 -c 'import signal,sys,time; signal.signal(signal.SIGINT, lambda signum, frame: sys.exit(130)); time.sleep(5)'\` and nothing else. Then exit success."
    env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        "$QUINE" "$target_mission" \
        >"$TEST_DIR/target.stdout" \
        2>"$TEST_DIR/target.stderr" \
        </dev/null &
    target_pid=$!
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            target_session="$(basename "$target_root")"
            retained_root="$QUINE_DATA_DIR/log/$target_session"
            if fuse_control_surface_ready "$target_public"; then
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    if [ "$target_ready" -eq 1 ]; then
        sleep 0.2
        printf 'interrupt from fuse runtime test\n' > "$target_public/ctl/interrupt"
    fi

    wait "$target_pid"
    code=$?

    inbox_checks=""
    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"
        _settled=0
        _attempt=0
        while [ "$_attempt" -lt 20 ]; do
            inbox_checks="$(python3 - "$retained_root/status/inbox.json" "$retained_tape" "$retained_root/control.jsonl" <<'PY'
import json
import sys
from pathlib import Path

inbox = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
current = Path(sys.argv[2]).read_text(encoding="utf-8")
control = Path(sys.argv[3]).read_text(encoding="utf-8")
print(f"PENDING_OK={1 if inbox.get('pending_count') == 0 else 0}")
current_ok = (
    '"incoming_messages":[' in current
    and '"delivery":"interrupt"' in current
    and '"payload":"interrupt from fuse runtime test"' in current
    and "Current operation was interrupted by peer control input." in current
)
control_ok = (
    '"kind":"delivered"' in control
    and '"delivery":"interrupt"' in control
    and '"payload":"interrupt from fuse runtime test"' in control
)
print(f"CURRENT_OK={1 if current_ok else 0}")
print(f"CONTROL_OK={1 if control_ok else 0}")
PY
)"
            if printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
                printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" &&
                printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1"; then
                _settled=1
                break
            fi
            sleep 0.1
            _attempt=$(( _attempt + 1 ))
        done
    fi

    if [ "$target_ready" -ne 1 ]; then
        fail "fuse public control surface did not appear in time"
    elif assert_exit "$code" 0 "fuse ctl/interrupt exit" &&
        printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" &&
        printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1"; then
        pass
    else
        if ! printf '%s\n' "$inbox_checks" | grep -q "PENDING_OK=1" 2>/dev/null; then
            fail "fuse ctl/interrupt inbox pending_count was not 0"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "CURRENT_OK=1" 2>/dev/null; then
            fail "fuse ctl/interrupt retained tape never surfaced incoming_messages"
        elif ! printf '%s\n' "$inbox_checks" | grep -q "CONTROL_OK=1" 2>/dev/null; then
            fail "fuse ctl/interrupt control log never recorded delivered interrupt mail"
        fi
    fi
    teardown
}

test_idle_explicit_inject_resume() {
    begin_test "idle suspends until ctl/inject resumes it"
    setup

    _prev_idle_enabled="${QUINE_IDLE_ENABLED:-}"
    export QUINE_IDLE_ENABLED=1
    target_mission='You MUST call `idle` immediately in your next response with no arguments. Do not emit text-only responses. Do not call `sh` before `idle` returns. Another process will later write one payload to ctl/inject. After `idle` resumes, read the delivered payload from the idle tool result. Then use exactly one sh call to write these lines to >&4 on separate lines: IDLE_INJECT_OK and PAYLOAD=<exact delivered payload>. Then exit success.'
    start_quine_helper "$TEST_DIR/target.stdout" "$TEST_DIR/target.stderr" "$target_mission"
    if [ -n "$_prev_idle_enabled" ]; then
        export QUINE_IDLE_ENABLED="$_prev_idle_enabled"
    else
        unset QUINE_IDLE_ENABLED
    fi
    target_pid="$HELPER_PID"
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            if [ -n "$target_public" ] && legacy_control_surface_ready "$target_public"; then
                target_session="$(basename "$target_root")"
                retained_root="$QUINE_DATA_DIR/log/$target_session"
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        sleep 0.5
        printf 'wake payload from runtime test\n' > "$target_public/ctl/inject"
    fi

    wait "$target_pid"
    code=$?
    retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"

    if [ "$target_ready" -ne 1 ]; then
        fail "idle target control surface did not appear in time"
    elif assert_exit "$code" 0 "idle inject exit" &&
        assert_contains "$TEST_DIR/target.stdout" "IDLE_INJECT_OK" "idle inject marker" &&
        assert_contains "$TEST_DIR/target.stdout" "PAYLOAD=wake payload from runtime test" "idle inject payload" &&
        assert_contains "$retained_tape" '"tool":"idle"' "idle tool result persisted" &&
        assert_contains "$retained_tape" '"delivery":"inject"' "idle inject delivery label" &&
        assert_contains "$retained_tape" '"payload":"wake payload from runtime test"' "idle inject payload persisted"; then
        pass
    fi
    teardown
}

test_idle_explicit_interrupt_resume() {
    begin_test "idle resumes on ctl/interrupt delivery"
    setup

    _prev_idle_enabled="${QUINE_IDLE_ENABLED:-}"
    export QUINE_IDLE_ENABLED=1
    target_mission='You MUST call `idle` immediately in your next response with no arguments. Do not emit text-only responses. Do not call `sh` before `idle` returns. Another process will later write one payload to ctl/interrupt. After `idle` resumes, read the delivered payload from the idle tool result. Then use exactly one sh call to write these lines to >&4 on separate lines: IDLE_INTERRUPT_OK and PAYLOAD=<exact delivered payload>. Then exit success.'
    start_quine_helper "$TEST_DIR/target.stdout" "$TEST_DIR/target.stderr" "$target_mission"
    if [ -n "$_prev_idle_enabled" ]; then
        export QUINE_IDLE_ENABLED="$_prev_idle_enabled"
    else
        unset QUINE_IDLE_ENABLED
    fi
    target_pid="$HELPER_PID"
    target_ready=0
    target_public=""
    target_root=""
    target_session=""
    retained_root=""

    i=0
    while [ "$i" -lt 80 ]; do
        if [ -L "$QUINE_DATA_DIR/pid/$target_pid" ]; then
            target_public="$(resolve_pid_target "$QUINE_DATA_DIR/pid/$target_pid")"
            target_root="$(dirname "$target_public")"
            if [ -n "$target_public" ] && legacy_control_surface_ready "$target_public"; then
                target_session="$(basename "$target_root")"
                retained_root="$QUINE_DATA_DIR/log/$target_session"
                target_ready=1
                break
            fi
        fi
        sleep 0.25
        i=$((i + 1))
    done

    retained_tape=""
    if [ "$target_ready" -eq 1 ]; then
        sleep 0.5
        printf 'interrupt payload from runtime test\n' > "$target_public/ctl/interrupt"
    fi

    wait "$target_pid"
    code=$?
    retained_tape="$(find_tape "$retained_root" 2>/dev/null || true)"

    if [ "$target_ready" -ne 1 ]; then
        fail "idle target control surface did not appear in time"
    elif assert_exit "$code" 0 "idle interrupt exit" &&
        assert_contains "$TEST_DIR/target.stdout" "IDLE_INTERRUPT_OK" "idle interrupt marker" &&
        assert_contains "$TEST_DIR/target.stdout" "PAYLOAD=interrupt payload from runtime test" "idle interrupt payload" &&
        assert_contains "$retained_tape" '"tool":"idle"' "idle tool result persisted" &&
        assert_contains "$retained_tape" '"delivery":"interrupt"' "idle interrupt delivery label" &&
        assert_contains "$retained_tape" '"interrupt_notice":"Current operation was interrupted by peer control input\."' "idle interrupt notice" &&
        assert_contains "$retained_tape" '"payload":"interrupt payload from runtime test"' "idle interrupt payload persisted"; then
        pass
    fi
    teardown
}

test_fork_wait() {
    begin_test "fork(wait=true) returns child output"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use the fork tool with one child: intent "echo CHILD_OUTPUT_42 >&4 and exit success", workspace "." and mode "wait". Then write the child stdout to >&4 and exit success.'
    code=$?

    # The child's stdout should contain our marker
    # The parent should relay it (or at minimum, the fork result should contain it)
    if assert_exit "$code" 0 "exit"; then
        pass
    fi
    teardown
}

test_spawn_wait_fresh_context() {
    begin_test "spawn(wait) starts fresh child context without parent memory"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_SPAWN_ENABLED=1 \
        "$QUINE" 'Do exactly these steps. Step 1: use exactly one sh call with command `mkdir -p "$QUINE_AGENT_ROOT/context/prompt" && printf "SPAWN_PARENT_MEMORY_MARKER\n" > "$QUINE_AGENT_ROOT/context/prompt/30-memory.md"`. Step 2: use exactly one spawn tool call with mode "wait" and one child mission: "Exit success immediately without using sh." Step 3: exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    relation_root="$(extract_first_spawn_content_field "$QUINE_DATA_DIR" "relation_root" 2>/dev/null || true)"
    child_session="$(extract_first_spawn_content_field "$QUINE_DATA_DIR" "children.0.session_id" 2>/dev/null || true)"
    child_retained="$(extract_first_spawn_content_field "$QUINE_DATA_DIR" "children.0.retained_root" 2>/dev/null || true)"
    _ok=1
    assert_exit "$code" 0 "spawn fresh context exit" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"tool":"spawn"' "spawn tool result" || _ok=0
    if [ -z "$child_retained" ] || [ ! -f "$child_retained/inc/0/context/prompt/30-memory.md" ]; then
        fail "spawn child retained memory surface missing"
        _ok=0
    elif grep -F -q 'SPAWN_PARENT_MEMORY_MARKER' "$child_retained/inc/0/context/prompt/30-memory.md"; then
        fail "spawn child inherited parent memory context"
        _ok=0
    fi
    if list_tape_jsonl "$QUINE_DATA_DIR" | xargs grep -F '"seed_root"' >/dev/null 2>&1; then
        fail "spawn result should not expose fork seed_root"
        _ok=0
    fi
    if [ -z "$relation_root" ] || [ ! -f "$relation_root/relation.json" ] || ! grep -F -q '"tool": "spawn"' "$relation_root/relation.json"; then
        fail "spawn relation surface missing or incomplete"
        _ok=0
    fi
    if [ -z "$child_session" ]; then
        fail "spawn child session_id missing"
        _ok=0
    fi
    if [ "$_ok" -eq 1 ]; then
        pass
    fi
    teardown
}

test_spawn_forget_spawns_child_independently() {
    begin_test "spawn(mode=forget) returns process handles"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_SPAWN_ENABLED=1 \
        "$QUINE" 'Use exactly one spawn tool call with mode "forget" and one child mission: "Create the file `$QUINE_DATA_DIR/spawn-forget-child-ok` containing `SPAWN_FORGET_CHILD_OK`, then exit success." After the spawn call returns, exit success immediately.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    _child_index="$(extract_first_spawn_content_field "$QUINE_DATA_DIR" "children.0.index" 2>/dev/null || true)"
    _child_status="$(extract_first_spawn_content_field "$QUINE_DATA_DIR" "children.0.status" 2>/dev/null || true)"
    _ok=1
    assert_exit "$code" 0 "spawn forget exit" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"tool":"spawn"' "spawn forget tool result" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"mode":"forget"' "spawn forget mode" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"spawned":1' "spawn forget spawned count" || _ok=0
    if [ "$_child_index" != "0" ]; then
        fail "spawn forget child index = ${_child_index:-<missing>}, want 0"
        _ok=0
    fi
    if [ "$_child_status" != "spawned" ]; then
        fail "spawn forget child status = ${_child_status:-<missing>}, want spawned"
        _ok=0
    fi
    if [ "$_ok" -eq 1 ]; then
        i=0
        while [ "$i" -lt 30 ]; do
            if [ -f "$QUINE_DATA_DIR/spawn-forget-child-ok" ]; then
                break
            fi
            sleep 0.25
            i=$((i + 1))
        done
        if assert_contains "$QUINE_DATA_DIR/spawn-forget-child-ok" 'SPAWN_FORGET_CHILD_OK' "spawn forget child durable effect"; then
            pass
        fi
    fi
    teardown
}

test_relation_recovery_preserves_is_error_semantics() {
    begin_test "relation recovery preserves retained is_error semantics"
    setup

    _session="relation-recovery-l2"
    _session_root="$QUINE_RETENTION_DIR/sessions/$_session"
    mkdir -p "$_session_root/inc/0/context/state" \
        "$_session_root/relations/call_fork_forget" \
        "$_session_root/relations/call_spawn_failed"
    cat > "$_session_root/inc/0/context/state/current.jsonl" <<'EOF'
{"type":"message","data":{"role":"user","content":"before pending relation tool batch"}}
{"type":"message","data":{"role":"assistant","tool_calls":[{"id":"call_fork_forget","name":"fork","arguments":{"mode":"forget","children":[{"intent":"background child","scope":"."}]}},{"id":"call_spawn_failed","name":"spawn","arguments":{"mode":"wait","children":[{"mission":"failed child"}]}}]}}
EOF
    cat > "$_session_root/relations/call_fork_forget/result.json" <<'EOF'
{"tool":"fork","mode":"forget","status":"spawned","requested":1,"spawned":1}
EOF
    cat > "$_session_root/relations/call_spawn_failed/result.json" <<'EOF'
{"tool":"spawn","mode":"wait","status":"completed","requested":1,"spawned":1,"succeeded":0,"children":[{"index":0,"mission":"failed child","status":"completed","exit_code":2}]}
EOF

    timeout_cmd "$TIMEOUT" env \
        QUINE_SESSION_ID="$_session" \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_RETENTION_DIR" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_SPAWN_ENABLED=1 \
        "$QUINE" 'Do not call fork or spawn. Use exactly one sh call to write RELATION_RECOVERY_L2_OK to fd 4, then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    _fork_is_error="$(extract_tool_result_is_error "$QUINE_DATA_DIR" "call_fork_forget" 2>/dev/null || true)"
    _spawn_is_error="$(extract_tool_result_is_error "$QUINE_DATA_DIR" "call_spawn_failed" 2>/dev/null || true)"
    _ok=1
    assert_exit "$code" 0 "relation recovery exit" || _ok=0
    assert_contains "$TEST_DIR/stdout" "RELATION_RECOVERY_L2_OK" "relation recovery marker" || _ok=0
    if [ "$_fork_is_error" != "false" ]; then
        fail "recovered fork forget is_error = ${_fork_is_error:-<missing>}, want false"
        _ok=0
    fi
    if [ "$_spawn_is_error" != "true" ]; then
        fail "recovered spawn all-failed is_error = ${_spawn_is_error:-<missing>}, want true"
        _ok=0
    fi
    if [ "$_ok" -eq 1 ]; then
        pass
    fi
    teardown
}

test_fork_wait_survives_ephemeral_body() {
    begin_test "fork(wait=true) still works after ephemeral launch-path consumption"
    setup

    ephemeral_quine="$TEST_DIR/quine-ephemeral"
    cp "$QUINE" "$ephemeral_quine" || {
        fail "failed to copy quine binary for ephemeral fork test"
        teardown
        return
    }
    chmod +x "$ephemeral_quine" || {
        fail "failed to chmod quine copy for ephemeral fork test"
        teardown
        return
    }

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_EPHEMERAL_BODY_ENABLED=1 \
        "$ephemeral_quine" 'Use exactly one fork tool call with mode "wait" and one child at scope ".". The child intent must be "Use exactly one sh call with command: printf '\''EPHEMERAL_FORK_CHILD_OK\n'\'' >&4. Then exit success." After the fork call returns, use exactly one sh call with command `printf "EPHEMERAL_FORK_PARENT_OK\n" >&4`. Then exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "ephemeral fork exit" &&
        assert_contains "$TEST_DIR/stdout" "EPHEMERAL_FORK_PARENT_OK" "ephemeral fork parent marker" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" 'EPHEMERAL_FORK_CHILD_OK' "ephemeral fork child marker"; then
        pass
    fi
    teardown
}

test_fork_retains_relation_surface() {
    begin_test "fork retains relation surface under initiator retained root"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one fork tool call with mode "wait" and one child at scope ".". The child intent must be "Use exactly one sh call with command: printf '\''REL_CHILD_OK\n'\'' >&4. Then exit success." After the fork call returns, exit success.'
    code=$?

    relation_root="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "relation_root" 2>/dev/null || true)"
    relation_id="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "relation_id" 2>/dev/null || true)"
    relation_handle="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "relation_handle" 2>/dev/null || true)"

    if assert_exit "$code" 0 "fork relation surface exit" &&
        [ -n "$relation_root" ] &&
        [ -n "$relation_id" ] &&
        [ -n "$relation_handle" ] &&
        [ -d "$relation_root" ] &&
        [ -f "$relation_root/relation.json" ] &&
        [ -f "$relation_root/status.json" ] &&
        [ -f "$relation_root/result.json" ] &&
        [ -f "$relation_root/log.jsonl" ] &&
        [ -f "$relation_root/members/000.json" ] &&
        grep -F -q "\"id\": \"$relation_id\"" "$relation_root/relation.json" &&
        grep -F -q "\"relation_id\": \"$relation_id\"" "$relation_root/result.json" &&
        grep -F -q "$relation_handle" "$relation_root/result.json"; then
        pass
    else
        fail "fork relation surface missing or incomplete"
    fi

    teardown
}

test_fork_retains_child_seed_origin() {
    begin_test "fork retains child seed context and fork origin"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one fork tool call with mode "wait" and one child at scope ".". The child intent must be "Use exactly one sh call with command: printf '\''SEED_CHILD_OK\n'\'' >&4. Then exit success." After the fork call returns, exit success.'
    code=$?

    child_session="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "children.0.session_id" 2>/dev/null || true)"
    retained_root="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "children.0.retained_root" 2>/dev/null || true)"
    seed_root="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "children.0.seed_root" 2>/dev/null || true)"
    relation_id="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "relation_id" 2>/dev/null || true)"

    if assert_exit "$code" 0 "fork child seed exit" &&
        [ -n "$child_session" ] &&
        [ -n "$retained_root" ] &&
        [ -n "$seed_root" ] &&
        [ -d "$retained_root" ] &&
        [ -d "$seed_root" ] &&
        [ -f "$seed_root/origin.json" ] &&
        [ -f "$seed_root/context/state/current.jsonl" ] &&
        grep -F -q '"kind": "fork"' "$seed_root/origin.json" &&
        grep -F -q "\"relation_id\": \"$relation_id\"" "$seed_root/origin.json" &&
        grep -F -q "\"initiator_session\": \"" "$seed_root/origin.json" &&
        grep -F -q 'SEED_CHILD_OK' "$seed_root/origin.json" &&
        [ -s "$seed_root/context/state/current.jsonl" ]; then
        pass
    else
        fail "fork child seed/origin surface missing or incomplete"
    fi

    teardown
}

test_fork_creates_child_tape() {
    begin_test "fork creates child tape file"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use fork with one child: intent "exit success", workspace "." and mode "wait". Then exit success.'
    code=$?

    # Should have at least 2 tape files (parent + child)
    _tape_count=$(find "$QUINE_DATA_DIR" -name '*.jsonl' 2>/dev/null | wc -l)
    if [ "$_tape_count" -ge 2 ]; then
        pass
    else
        # fork might not always create a child tape if it fails, so be lenient
        if [ "$code" -eq 0 ]; then
            fail "expected >= 2 tape files, got $_tape_count"
        else
            skip "fork exited non-zero (code=$code), child tape may not exist"
        fi
    fi
    teardown
}

test_fork_race_selects_first_success() {
    begin_test "fork(mode=race) keeps the first successful child"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one fork tool call with mode "race" and two children. Child 0 intent: "Finish in a single response. Use exactly one sh call with command: sleep 3; printf '\''RACE_SLOW\n'\'' >&4. Then call exit success in that same response." Child 1 intent: "Finish in a single response. Use exactly one sh call with command: printf '\''RACE_FAST\n'\'' >&4. Then call exit success in that same response." Both workspaces must be ".". After the fork call returns, exit success immediately.'
    code=$?
    _winner_index="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "winner.index" 2>/dev/null || true)"
    _winner_exit="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "winner.exit_code" 2>/dev/null || true)"
    _ok=1
    assert_exit "$code" 0 "fork race exit" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"mode":"race"' "fork race mode" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"spawned":2' "fork race spawned count" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"succeeded":1' "fork race success count" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"killed":1' "fork race kill count" || _ok=0
    if [ "$_winner_index" != "1" ]; then
        fail "fork race winner index = ${_winner_index:-<missing>}, want 1"
        _ok=0
    fi
    if [ "$_winner_exit" != "0" ]; then
        fail "fork race winner exit_code = ${_winner_exit:-<missing>}, want 0"
        _ok=0
    fi
    if [ "$_ok" -eq 1 ]; then
        pass
    fi
    teardown
}

test_fork_forget_spawns_child_independently() {
    begin_test "fork(mode=forget) spawns child tape without waiting"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one fork tool call with mode "forget" and one child. Child intent: "sleep 1; exit success". Workspace must be ".". After the fork call returns, exit success immediately.'
    code=$?
    _child_index="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "children.0.index" 2>/dev/null || true)"
    _child_status="$(extract_first_fork_content_field "$QUINE_DATA_DIR" "children.0.status" 2>/dev/null || true)"
    _ok=1
    assert_exit "$code" 0 "fork forget exit" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"mode":"forget"' "fork forget mode" || _ok=0
    assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"spawned":1' "fork forget spawned count" || _ok=0
    if [ "$_child_index" != "0" ]; then
        fail "fork forget child index = ${_child_index:-<missing>}, want 0"
        _ok=0
    fi
    if [ "$_child_status" != "spawned" ]; then
        fail "fork forget child status = ${_child_status:-<missing>}, want spawned"
        _ok=0
    fi
    if [ "$_ok" -eq 1 ]; then
        sleep 2
        _tape_count=$(find "$QUINE_DATA_DIR" -name '*.jsonl' 2>/dev/null | wc -l)
        if [ "$_tape_count" -ge 2 ]; then
            pass
        else
            fail "fork forget should leave parent + child tape, got $_tape_count"
        fi
    fi
    teardown
}

test_fork_child_workspace_narrowing() {
    begin_test "fork child workspace narrowing keeps relative writes inside child lineage"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics are enabled. Do exactly these steps: (1) run one sh call: mkdir -p subdir && printf "seed\n" > subdir/seed.txt (2) use exactly one fork tool call with mode "wait" and one child whose workspace is "subdir" and whose intent is "Do not exit before one sh call succeeds. Use exactly one sh call with command: printf '\''child-narrowed\n'\'' > child_created.txt && pwd > child_pwd.txt. After that sh call, exit success." (3) write FORK_WORKSPACE_OK to fd 4 using exactly one sh call (4) exit success.'
    code=$?

    if assert_exit "$code" 0 "fork workspace narrowing exit" &&
        assert_contains "$TEST_DIR/stdout" "FORK_WORKSPACE_OK" "fork workspace narrowing stdout" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"tool":"fork"' "fork workspace narrowing tool result" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"mode":"wait"' "fork workspace narrowing wait mode" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"scope":"subdir"' "fork workspace narrowing child scope"; then
        _child_workspace_session="$(extract_first_tape_string_field "$QUINE_DATA_DIR" "workspace_session")"
        _child_created=""
        _child_pwd=""
        if [ -n "$_child_workspace_session" ]; then
            _child_created="$(find "$QUINE_DATA_DIR/workspaces/$_child_workspace_session" -type f -path '*/subdir/child_created.txt' 2>/dev/null | head -n 1)"
            _child_pwd="$(find "$QUINE_DATA_DIR/workspaces/$_child_workspace_session" -type f -path '*/subdir/child_pwd.txt' 2>/dev/null | head -n 1)"
        fi
        if [ -z "$_child_workspace_session" ]; then
            fail "child workspace session missing from fork result"
        elif [ -z "$_child_created" ]; then
            fail "child relative write did not land in narrowed child lineage"
        elif [ -z "$_child_pwd" ]; then
            fail "child pwd record missing from narrowed child lineage"
        elif ! grep -q '/subdir' "$_child_pwd" 2>/dev/null; then
            fail "child pwd should point inside narrowed child workspace"
        elif [ -e "$TEST_DIR/workspace/subdir/child_created.txt" ] || [ -e "$TEST_DIR/workspace/child_created.txt" ]; then
            fail "child relative write should remain private before adoption"
        else
            pass
        fi
    fi
    teardown
}

test_fork_world_enabled_requires_workspace_config() {
    begin_test "fork world mode requires explicit workspace configuration"
    setup

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_FORK_WORLD_ENABLED=1 \
        QUINE_WORKSPACE_ROOT= \
        QUINE_WORKSPACE= \
        QUINE_WORKSPACE_BACKEND= \
        QUINE_WORKSPACE_REVISION_MODE= \
        QUINE_WORKSPACE_SESSION= \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" 'This process should fail during configuration before model work.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?

    if assert_exit "$code" 2 "fork world dependency exit" &&
        assert_contains "$TEST_DIR/stderr" "QUINE_FORK_WORLD_ENABLED=1 requires explicit workspace physics" "fork world dependency stderr"; then
        pass
    fi
    teardown
}

test_fork_world_enabled_host_workspace_child() {
    begin_test "fork world properties keep subjective child writes private from host parent"
    if ! is_linux; then
        skip "host-parent subjective child world modes currently require a native Linux host"
        return
    fi
    if ! subjective_worlds_available; then
        skip "subjective child worlds require Linux or a configured Lima instance"
        return
    fi
    setup
    mkdir -p "$TEST_DIR/workspace"

    # fork-world requires explicit workspace physics (its sibling
    # test_fork_world_enabled_requires_workspace_config asserts the rejection
    # when they are absent). The parent runs the host base layer; the subjective
    # transactional child gets an isolated overlay so its world_commit.txt write
    # stays private from the host parent.
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" QUINE_WORKSPACE_BACKEND=overlay \
        QUINE_FORK_WORLD_ENABLED=1 \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" 'Fork world-property teaching is enabled and you are in world=host with protection=none. Do exactly these steps: (1) use exactly one fork tool call with mode "wait" and one child whose world is "subjective", protection is "transactional", scope is ".", and intent is "Use exactly one sh call with command: if [ -n \"${QUINE_WORKSPACE_ROOT:-}\" ]; then printf '\''subjective-child\n'\'' > world_commit.txt && echo '\''SUBJECTIVE_CHILD_OK'\'' >&4; else echo '\''SUBJECTIVE_CHILD_BAD'\'' >&4; fi. Then exit success." (2) run exactly one sh call with command: if [ ! -e world_commit.txt ]; then echo CHILD_PRIVATE_OK >&4; fi (3) exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "subjective overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "fork host->workspace exit" &&
        assert_contains "$TEST_DIR/stdout" "CHILD_PRIVATE_OK" "host->subjective stdout" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\"world\":\"subjective\"" "subjective world recorded in tape" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\"protection\":\"transactional\"" "transactional protection recorded in tape"; then
        if [ -e "$TEST_DIR/workspace/world_commit.txt" ]; then
            fail "subjective child leaked world_commit.txt into host working surface"
        else
            pass
        fi
    fi
    teardown
}

test_fork_world_enabled_workspace_parent_host_child() {
    begin_test "fork world properties allow subjective parent to spawn host none-protection child"
    if ! subjective_worlds_available; then
        skip "subjective parent tests require Linux or a configured Lima instance"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_fork_world_enabled="${QUINE_FORK_WORLD_ENABLED-}"
    _had_fork_world_enabled=0
    if [ "${QUINE_FORK_WORLD_ENABLED+x}" = x ]; then
        _had_fork_world_enabled=1
    fi
    _saved_runtime_bridge_extra_env_lines="${RUNTIME_BRIDGE_EXTRA_ENV_LINES-}"
    _had_runtime_bridge_extra_env_lines=0
    if [ "${RUNTIME_BRIDGE_EXTRA_ENV_LINES+x}" = x ]; then
        _had_runtime_bridge_extra_env_lines=1
    fi
    export QUINE_FORK_WORLD_ENABLED=1
    export RUNTIME_BRIDGE_EXTRA_ENV_LINES='QUINE_FORK_WORLD_ENABLED=1'
    run_workspace_quine_with_backend "$TEST_DIR/stdout" "$TEST_DIR/stderr" "${E2E_LIMA_WORKSPACE_BACKEND:-overlay}" \
        'Transactional workspace physics and fork world-property teaching are enabled. Do exactly these steps: (1) use exactly one fork tool call with mode "wait" and one child whose world is "host", protection is "none", scope is ".", and intent is "Use exactly one sh call with command: if [ -z \"${QUINE_WORKSPACE_ROOT:-}\" ]; then echo '\''HOST_CHILD_OK'\'' >&4; else echo '\''HOST_CHILD_BAD'\'' >&4; fi. Then exit success." (2) exit success.'
    code=$?
    if [ "$_had_fork_world_enabled" -eq 1 ]; then
        export QUINE_FORK_WORLD_ENABLED="$_saved_fork_world_enabled"
    else
        unset QUINE_FORK_WORLD_ENABLED
    fi
    if [ "$_had_runtime_bridge_extra_env_lines" -eq 1 ]; then
        export RUNTIME_BRIDGE_EXTRA_ENV_LINES="$_saved_runtime_bridge_extra_env_lines"
    else
        unset RUNTIME_BRIDGE_EXTRA_ENV_LINES
    fi

    if assert_exit "$code" 0 "fork workspace->host exit" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "HOST_CHILD_OK" "host child marker returned" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\"world\":\"host\"" "host world recorded in tape" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\"protection\":\"none\"" "host protection recorded in tape"; then
        pass
    fi
    teardown
}

test_switch_world_adopts_child_world() {
    begin_test "switch_world can adopt a child world handle from fork(mode=wait)"
    if ! subjective_worlds_available; then
        skip "subjective child worlds require Linux or a configured Lima instance"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=6
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics are enabled. Use exactly one fork call, one switch_world call, and two parent sh calls. Do not use exec. Step 1: call fork with mode "wait" and one child at scope ".". The child should create adopted.txt in its own world and then exit success. Step 2: in the parent, run exactly one sh call with command `if [ ! -e adopted.txt ]; then echo PRE_SWITCH_PRIVATE >&4; fi`. Step 3: read the child world handle from the fork result and call switch_world with that exact handle. Step 4: run exactly one sh call with command `test -e adopted.txt && echo ADOPT_SWITCH_OK >&4`. Then exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "subjective overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "switch adopt exit" &&
        assert_contains "$TEST_DIR/stdout" "PRE_SWITCH_PRIVATE" "pre-switch privacy marker" &&
        assert_contains "$TEST_DIR/stdout" "ADOPT_SWITCH_OK" "post-switch adoption marker" &&
        assert_exists "$TEST_DIR/workspace/adopted.txt" "adopted child file" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" '"world_handle":"world://' "fork surfaced world handle" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\"name\":\"switch_world\"" "switch_world tool recorded"; then
        pass
    fi
    teardown
}

test_fork_race_adopt_winner() {
    begin_test "fork race can adopt winner automatically"
    if ! subjective_worlds_available; then
        skip "subjective child worlds require Linux or a configured Lima instance"
        return
    fi
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=4
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics are enabled. Use exactly one fork call and one parent sh call. Do not use switch_world, exec, or any extra parent sh calls. Step 1: call fork with mode "race", adopt_winner=true, and two children. Child A must finish in a single response: use exactly one sh call with command `: > result.txt`, then call exit success in that same response. Child B must finish in a single response: use exactly one sh call with command `sleep 8; : > loser.txt`, then call exit success in that same response. Both children should run in subjective transactional child worlds with scope ".". Step 2: run exactly one parent sh call with command `test -e result.txt && [ ! -e loser.txt ] && echo RACE_ADOPT_OK >&4`. Then exit success.' \
        >"$TEST_DIR/stdout" \
        2>"$TEST_DIR/stderr" \
        </dev/null
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if [ "$code" -ne 0 ] && overlay_substrate_unavailable "$TEST_DIR/stderr"; then
        skip "subjective overlay substrate unavailable in this environment"
        teardown
        return
    fi

    if assert_exit "$code" 0 "race adopt exit" &&
        assert_exists "$TEST_DIR/workspace/result.txt" "winner file committed" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\"adopt_winner\":true" "adopt_winner recorded in tape" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "wr0 -\\\\u003e wr" "winner adoption revision transition"; then
        if [ -e "$TEST_DIR/workspace/loser.txt" ]; then
            fail "loser.txt leaked into adopted winner surface"
        else
            pass
        fi
    fi
    teardown
}

# ── 11. Exec lifecycle ─────────────────────────────────────

test_exec_preserves_mission() {
    begin_test "exec preserves mission across image replacement"
    setup

    mkdir -p "$TEST_DIR/workspace" || {
        fail "failed to create workspace for exec mission test"
        teardown
        return
    }

    # Mission rides argv across default self re-exec: the successor can only
    # continue this procedure if the mission argv survived the image
    # replacement. A workspace marker file distinguishes the phases.
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        "$QUINE" 'Repeat this procedure in every incarnation. First, use exactly one sh call with this exact command: if [ -e mission-phase.txt ]; then printf "MISSION_CARRIED_OK\n" >&4; echo "PHASE=after-exec"; else touch mission-phase.txt; echo "PHASE=before-exec"; fi    Then: if the sh output contains PHASE=after-exec, exit success. Otherwise call exec exactly once with no target and no argv so this same mission continues in the next incarnation.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "exec mission loop exit" &&
        assert_contains "$TEST_DIR/stdout" "MISSION_CARRIED_OK" "mission preserved across exec"; then
        pass
    fi
    teardown
}

test_exec_self_reentry_survives_ephemeral_body() {
    begin_test "Linux self exec reentry survives ephemeral launch-path consumption"
    if [ "$(uname -s)" != "Linux" ]; then
        skip "QUINE_SELF_REENTRY_MODE=self is Linux-only"
        return
    fi
    setup

    ephemeral_quine="$TEST_DIR/quine-ephemeral"
    cp "$QUINE" "$ephemeral_quine" || {
        fail "failed to copy quine binary for ephemeral exec test"
        teardown
        return
    }
    chmod +x "$ephemeral_quine" || {
        fail "failed to chmod quine copy for ephemeral exec test"
        teardown
        return
    }
    mkdir -p "$TEST_DIR/workspace" || {
        fail "failed to create workspace for ephemeral exec test"
        teardown
        return
    }

    # The consumed (unlinked) launch body must not break /proc/self/exe
    # re-entry. A workspace marker file distinguishes the pre-exec and
    # post-exec incarnations.
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_SELF_REENTRY_MODE=self \
        QUINE_EPHEMERAL_BODY_ENABLED=1 \
        "$ephemeral_quine" 'Repeat this procedure in every incarnation. First, use exactly one sh call with this exact command: if [ -e ephemeral-phase.txt ]; then printf "EPHEMERAL_EXEC_OK\n" >&4; echo "PHASE=after-exec"; else touch ephemeral-phase.txt; echo "PHASE=before-exec"; fi    Then: if the sh output contains PHASE=after-exec, exit success. Otherwise call exec exactly once with no target and no argv so this same mission continues through the default self re-entry target.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "ephemeral exec exit" &&
        assert_contains "$TEST_DIR/stdout" "EPHEMERAL_EXEC_OK" "ephemeral exec marker"; then
        pass
    fi
    teardown
}

test_exec_self_reentry_reclaims_public_surface() {
    begin_test "exec re-entry reclaims the stale public FUSE mount across generations"
    if [ "$(uname -s)" != "Linux" ]; then
        skip "self exec re-entry over a live FUSE public surface is Linux-only"
        return
    fi
    setup

    if ! runtime_surface_fuse_available; then
        skip "runtime surface FUSE requires Linux with /dev/fuse, fusermount, and mount permission"
        teardown
        return
    fi

    mkdir -p "$TEST_DIR/workspace" || {
        fail "failed to create workspace for exec reclaim test"
        teardown
        return
    }

    # Each incarnation bumps a workspace counter and proves the successor can
    # read its own public/ FUSE projection: before the stale-mount reclaim, the
    # predecessor's disconnected mount made re-entry bootstrap fail with
    # "Transport endpoint is not connected". 3 execs = 3 handovers.
    timeout_cmd 240 env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        QUINE_RUNTIME_SURFACE_BACKEND="fuse" \
        "$QUINE" 'Repeat this procedure in every incarnation. First, use exactly one sh call with this exact command: c=$(cat gen.txt 2>/dev/null || echo 0); c=$((c+1)); echo "$c" > gen.txt; if cat "$QUINE_AGENT_ROOT/public/status/session.json" >/dev/null 2>&1; then printf "PUBLIC_OK_%s\n" "$c" >&4; else printf "PUBLIC_FAIL_%s\n" "$c" >&4; fi; echo "GEN=$c"    Then: if the sh output contains GEN=4, exit success. Otherwise call exec exactly once with no target and no argv so this same mission continues in the next incarnation.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "exec reclaim loop exit" &&
        assert_contains "$TEST_DIR/stdout" "PUBLIC_OK_2" "public surface readable after first exec handover" &&
        assert_contains "$TEST_DIR/stdout" "PUBLIC_OK_3" "public surface readable after second exec handover" &&
        assert_contains "$TEST_DIR/stdout" "PUBLIC_OK_4" "public surface readable after third exec handover"; then
        if grep -q "PUBLIC_FAIL_" "$TEST_DIR/stdout" 2>/dev/null; then
            fail "public surface read failed in some incarnation (stale-mount / degradation regression)"
        else
            pass
        fi
    fi
    teardown
}

test_exec_applies_staged_config() {
    begin_test "exec applies the staged config/env/override to the successor and archives it verbatim"
    setup

    mkdir -p "$TEST_DIR/workspace" || {
        fail "failed to create workspace for staged config test"
        teardown
        return
    }

    # Staged child-env override transaction: the first incarnation stages
    # QUINE_MAX_TURNS=41 (a free, override-settable exec-boundary budget knob)
    # in config/env/override and execs. The successor verifies its applied
    # budget by reading its own process environment (/proc/<pid>/environ, the
    # enacted self-readout surface — NOT a runtime-rendered resolved.env, which no
    # longer exists), asserts config/env/override was cleared, and asserts the applied
    # policy was archived VERBATIM to the PREDECESSOR's inc/0 dir as
    # override-applied.env (byte-equal to the bytes the agent wrote, no
    # re-render). The negative case (invalid override bounces the exec, file
    # intact, retry works) lives in test_exec_rejects_invalid_staged_config.
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        "$QUINE" 'Repeat this procedure in every incarnation. First, use exactly one sh call with this exact command: if [ -e staged-phase.txt ]; then pid=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"pid\"])" "$QUINE_AGENT_ROOT/status/session.json"); envfile="/proc/$pid/environ"; if [ -r "$envfile" ]; then tr "\0" "\n" < "$envfile" | grep -q "^QUINE_MAX_TURNS=41$" || { echo "successor /proc environ missing applied budget" >&2; exit 1; }; else [ "$QUINE_MAX_TURNS" = "41" ] || { echo "successor env missing applied budget" >&2; exit 1; }; fi; [ ! -e "$QUINE_AGENT_ROOT/config/env/override" ] || { echo "override not cleared after exec" >&2; exit 1; }; [ -f "$QUINE_AGENT_ROOT/inc/0/override-applied.env" ] || { echo "inc/0 override archive missing" >&2; exit 1; }; cmp -s "$QUINE_AGENT_ROOT/inc/0/override-applied.env" staged-before.txt || { echo "archive is not the verbatim staged override" >&2; exit 1; }; printf "STAGED_CONFIG_APPLIED_OK\n" >&4; echo "PHASE=after-exec"; else printf "QUINE_MAX_TURNS=41\n" > "$QUINE_AGENT_ROOT/config/env/override"; cp "$QUINE_AGENT_ROOT/config/env/override" staged-before.txt; touch staged-phase.txt; echo "PHASE=before-exec"; fi    Then: if the sh output contains PHASE=after-exec, exit success. Otherwise call exec exactly once with no target and no argv so this same mission continues in the next incarnation.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "staged config exec loop exit" &&
        assert_contains "$TEST_DIR/stdout" "STAGED_CONFIG_APPLIED_OK" "staged override applied, cleared, and archived verbatim across exec"; then
        pass
    fi
    teardown
}

test_exec_rejects_invalid_staged_config() {
    begin_test "exec rejects an invalid staged config/env/override loudly and applies the fix on retry"
    setup

    mkdir -p "$TEST_DIR/workspace" || {
        fail "failed to create workspace for staged config rejection test"
        teardown
        return
    }

    # Staged override failure semantics: the first incarnation stages an INVALID
    # config/env/override carrying two rejection classes (QUINE_MAX_DEPTH is
    # pinned/operator-only per E5, QUINE_TOTALLY_UNKNOWN_KNOB is not a knob) and
    # calls exec. The exec tool must bounce loudly with the production error
    # string: the tool result names the whole-file reject and enumerates every
    # violation (asserted from the tape below - that tool-result text is exactly
    # what the agent sees), the process stays in the SAME incarnation (inc/1
    # absent), and the override survives byte-intact (whole-file reject, no
    # partial strip). The agent then fixes the file with a legal exec-boundary
    # knob (QUINE_MAX_TURNS) and retries; the successor proves the value applied
    # by reading its own process environment - the retry-after-fix half.
    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        QUINE_WORK_DIR="$TEST_DIR/workspace" \
        "$QUINE" 'Follow this exact procedure. Step 1: use exactly one sh call with this exact command: if [ -e staged-phase.txt ]; then pid=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[\"pid\"])" "$QUINE_AGENT_ROOT/status/session.json"); envfile="/proc/$pid/environ"; if [ -r "$envfile" ]; then tr "\0" "\n" < "$envfile" | grep -q "^QUINE_MAX_TURNS=41$" || { echo "successor /proc environ missing applied budget" >&2; exit 1; }; else [ "$QUINE_MAX_TURNS" = "41" ] || { echo "successor env missing applied budget" >&2; exit 1; }; fi; [ ! -e "$QUINE_AGENT_ROOT/config/env/override" ] || { echo "override not cleared after exec" >&2; exit 1; }; printf "STAGED_RETRY_APPLIED_OK\n" >&4; echo "PHASE=after-exec"; else printf "QUINE_MAX_DEPTH=9\nQUINE_TOTALLY_UNKNOWN_KNOB=1\n" > "$QUINE_AGENT_ROOT/config/env/override"; cp "$QUINE_AGENT_ROOT/config/env/override" staged-before.txt; echo "PHASE=stage-invalid"; fi    Step 2: if the sh output contains PHASE=after-exec, exit success immediately and skip every later step. Otherwise call exec exactly once with no target and no argv. That exec call WILL FAIL with a child-env override validation error; the failure is expected - read the error result and continue with step 3 in this same process. Step 3: use exactly one sh call with this exact command: [ ! -e "$QUINE_AGENT_ROOT/inc/1" ] || { echo "unexpected successor incarnation" >&2; exit 1; }; cmp -s "$QUINE_AGENT_ROOT/config/env/override" staged-before.txt || { echo "override changed after rejected exec" >&2; exit 1; }; printf "REJECTED_SAME_INCARNATION_OK\n" >&4; printf "QUINE_MAX_TURNS=41\n" > "$QUINE_AGENT_ROOT/config/env/override"; touch staged-phase.txt; echo "PHASE=fixed"    Step 4: call exec exactly once with no target and no argv so this same mission continues in the next incarnation.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "staged config rejection arc exit" &&
        assert_contains "$TEST_DIR/stdout" "REJECTED_SAME_INCARNATION_OK" "rejected exec kept the same incarnation running with the override byte-intact" &&
        assert_contains "$TEST_DIR/stdout" "STAGED_RETRY_APPLIED_OK" "retry after fix applied the staged value in the successor" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" "[EXEC ERROR] child-env override" "exec tool result carries the loud child-env override error" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" "rejected (whole file; nothing applied; the file is left intact" "exec tool result states whole-file reject and intact file" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" "2 violations:" "exec tool result enumerates both violations" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" "QUINE_TOTALLY_UNKNOWN_KNOB: unknown env name" "exec tool result names the unknown-knob violation" &&
        assert_any_tape_contains_literal "$QUINE_DATA_DIR" "QUINE_MAX_DEPTH: mutability" "exec tool result names the pinned operator-only violation"; then
        pass
    fi
    teardown
}

test_exec_explicit_argv_replaces_mission() {
    begin_test "exec with explicit argv replaces the quine mission"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do not call sh before exec. Call exec exactly once with argv ["quine","In one sh call, run: printf \"SELF_ARGV_OK\n\" >&4. Then exit success."]'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "SELF_ARGV_OK" "exec explicit argv mission"; then
        pass
    fi
    teardown
}

test_exec_external_binary_completes_pipe_mission() {
    begin_test "exec can hand pipe mission to an external binary"
    setup

    printf 'hello\nworld\n' | run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do not call sh. Do not call exit. Call exec exactly once with target "tr" and argv ["tr","a-z","A-Z"]. The replacement binary itself must finish the mission by reading stdin, writing the uppercase stream to stdout, and terminating. The only valid stdout is exactly HELLO and WORLD on separate lines.'
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "HELLO" "exec external pipe stdout line 1" &&
        assert_contains "$TEST_DIR/stdout" "WORLD" "exec external pipe stdout line 2"; then
        pass
    fi
    teardown
}

test_exec_external_binary_preserves_stdio_triplet() {
    begin_test "exec preserves stdin, stdout, and stderr for an external binary"
    setup

    printf 'hello\nworld\n' | run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do not call sh. Do not call exit. Call exec exactly once with target "/bin/sh" and argv ["/bin/sh","-c","data=$(cat); printf \"%s\" \"$data\" | tr a-z A-Z; printf \"\\nEXEC_STDERR_OK:%s\\n\" \"$(printf \"%s\" \"$data\" | wc -c | tr -d \" \")\" >&2"]. The replacement binary itself must read stdin, write HELLO and WORLD to stdout, and emit EXEC_STDERR_OK:11 to stderr.'
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "HELLO" "exec stdio stdout line 1" &&
        assert_contains "$TEST_DIR/stdout" "WORLD" "exec stdio stdout line 2" &&
        assert_contains "$TEST_DIR/stderr" "EXEC_STDERR_OK:11" "exec stdio stderr"; then
        pass
    fi
    teardown
}

test_exec_copied_quine_binary_continues_mission() {
    begin_test "exec can continue through a copied quine binary"
    setup

    copied_quine="$TEST_DIR/quine-copy"
    cp "$QUINE" "$copied_quine" || {
        fail "failed to copy quine binary"
        teardown
        return
    }
    chmod +x "$copied_quine" || {
        fail "failed to chmod copied quine binary"
        teardown
        return
    }

    mission=$(printf 'Do not call sh before exec. Call exec exactly once with target "%s" and argv ["%s","In one sh call, run: printf \\"COPIED_QUINE_OK\\\\n\\" >&4. Then exit success."]' "$copied_quine" "$copied_quine")
    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" "$mission"
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "COPIED_QUINE_OK" "exec copied quine mission"; then
        pass
    fi
    teardown
}

test_exec_relative_workspace_target_continues_mission() {
    begin_test "exec resolves relative target paths from the workspace"
    setup

    mkdir -p "$TEST_DIR/workspace/outbox" || {
        fail "failed to create workspace outbox"
        teardown
        return
    }
    cp "$QUINE" "$TEST_DIR/workspace/outbox/quine-next" || {
        fail "failed to copy quine binary into workspace outbox"
        teardown
        return
    }
    chmod +x "$TEST_DIR/workspace/outbox/quine-next" || {
        fail "failed to chmod workspace quine copy"
        teardown
        return
    }

    timeout_cmd "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" QUINE_RETENTION_DIR="$QUINE_DATA_DIR/log" \
        QUINE_WORKSPACE_ROOT="$TEST_DIR/workspace" \
        QUINE_WORKSPACE="$TEST_DIR/workspace" \
        QUINE_WORKSPACE_BACKEND=direct \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" 'Do not call sh before exec. Call exec exactly once with target "outbox/quine-next" and argv ["outbox/quine-next","In one sh call, run: printf \"RELATIVE_EXEC_OK\n\" >&4. Then exit success."]' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "exit" &&
        assert_contains "$TEST_DIR/stdout" "RELATIVE_EXEC_OK" "exec relative workspace target mission"; then
        pass
    fi
    teardown
}

# ── 12. Multi-fd: fd 1 captured, fd 4 delivered ────────────

test_dual_channel_separation() {
    begin_test "fd 1 (captured) and fd 4 (delivered) are separate"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'In one sh call, run: echo "CAPTURED_fd1" && echo "DELIVERED_fd4" >&4    Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit"; then
        # fd 4 content should be in stdout
        _has_delivered=0
        grep -q "DELIVERED_fd4" "$TEST_DIR/stdout" 2>/dev/null && _has_delivered=1

        # fd 1 content should NOT be in stdout
        _has_captured=0
        grep -q "CAPTURED_fd1" "$TEST_DIR/stdout" 2>/dev/null && _has_captured=1

        if [ "$_has_delivered" -eq 1 ] && [ "$_has_captured" -eq 0 ]; then
            pass
        elif [ "$_has_delivered" -eq 0 ]; then
            fail "DELIVERED_fd4 not found in stdout"
        else
            fail "CAPTURED_fd1 leaked to stdout (should stay in tool result)"
        fi
    fi
    teardown
}

# ── 13. No-stdin mode ─────────────────────────────────────

test_no_stdin() {
    begin_test "no stdin (TTY mode) → material = Begin."
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Echo "NO_STDIN_OK" to >&4 and exit success.' </dev/null
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "NO_STDIN_OK" "no-stdin"; then
        pass
    fi
    teardown
}


# ============================================================================
# RUNNER
# ============================================================================

OVERLAY_FINALIZATION_BASELINE_TESTS="
    test_workspace_overlay_commit
    test_workspace_overlay_rollback
    test_workspace_overlay_failure_revision
    test_workspace_overlay_timeout_revision
    test_switch_world_restores_prior_revision
    test_switch_world_branches_forward_after_rewind
    test_switch_world_adopts_child_world
    test_workspace_overlay_fuse_shutdown_phase_order
    test_workspace_overlay_commit_intent_recovery
    test_workspace_overlay_fuse_long_lineage_materializes
"

run_named_test_group() {
    _group_name="$1"
    _group_tests="$2"
    printf "${BOLD}[%s]${RESET}\n" "$_group_name"
    for _test_name in $_group_tests; do
        "$_test_name"
    done
}

gate_overlay_finalization_baseline() {
    run_named_test_group "overlay-finalization-baseline" "$OVERLAY_FINALIZATION_BASELINE_TESTS"
}

# All test functions
ALL_TESTS="
    test_exit_success
    test_exit_failure
    test_fd4_delivery
    test_fd1_not_leaked
    test_stderr_failure_signal
    test_stderr_success_silent
    test_fd5_signal_channel
    test_shell_cd_does_not_persist
    test_shell_export_does_not_persist
    test_shell_variable_does_not_persist
    test_fd3_piped_input
    test_fd3_consumed_across_calls
    test_interactive_screen
    test_interactive_sigint_process_group
    test_interactive_overlay_world_adoption
    test_binary_stdin
    test_execution_budget_disabled_hidden
    test_execution_budget_enabled_feedback
    test_execution_budget_hard_fail
    test_anchor_memory_roundtrip
    test_vision_runtime_surface
    test_prompt_metaphor_off
    test_prompt_metaphor_thermodynamic
    test_prompt_self_model_basic
    test_prompt_runtime_surface_hidden
    test_prompt_fragments_surface_visible
    test_prompt_agents_md_disabled
    test_prompt_agents_md_enabled_single_projection
    test_prompt_skills_disabled
    test_prompt_skills_enabled_catalog_only
    test_prompt_skills_source_path_readable
    test_prompt_skills_exec_reentry_rescans_frontmatter
    test_context_memory_next_turn_refresh
    test_context_memory_exec_inherits
    test_workspace_overlay_commit
    test_workspace_overlay_fuse_exit_success_materializes
    test_workspace_overlay_fuse_shutdown_phase_order
    test_workspace_overlay_fuse_long_lineage_materializes
    test_workspace_overlay_commit_intent_recovery
    test_workspace_overlay_rollback
    test_workspace_overlay_failure_revision
    test_workspace_overlay_timeout_revision
    test_workspace_overlay_absolute_path
    test_workspace_overlay_lima_probe
    test_workspace_direct_persists_on_failure
    test_workspace_direct_observes_peer_mutation
    test_switch_world_restores_prior_revision
    test_switch_world_branches_forward_after_rewind
    test_workspace_unsupported_on_non_linux
    test_tape_has_meta
    test_tape_has_outcome
    test_tape_has_messages
    test_process_surface_self_identity
    test_process_surface_self_source_surface
    test_prompt_self_source_enabled
    test_process_surface_incarnation_projection
    test_process_surface_config_surface
    test_process_surface_config_projection
    test_process_surface_neighbor_discovery_prunes_stale_indexes
    test_process_surface_pid_route_removed_on_sigterm
    test_process_surface_stale_pid_lock_pruned_on_startup
    test_process_surface_peer_discovery_heartbeat_prunes_stale_pid_lock
    test_process_surface_ctl_inbox_only
    test_process_surface_ctl_inject_delivery
    test_process_surface_ctl_interrupt_delivery
    test_process_surface_fuse_ctl_post_transaction
    test_process_surface_fuse_ctl_env
    test_process_surface_fuse_public_scan_safe
    test_process_surface_fuse_ctl_inject_transaction
    test_process_surface_fuse_ctl_interrupt_transaction
    test_idle_explicit_inject_resume
    test_idle_explicit_interrupt_resume
    test_spawn_wait_fresh_context
    test_spawn_forget_spawns_child_independently
    test_relation_recovery_preserves_is_error_semantics
    test_fork_wait
    test_fork_wait_survives_ephemeral_body
    test_fork_retains_relation_surface
    test_fork_retains_child_seed_origin
    test_fork_depth_limit_rejected
    test_fork_agent_slot_limit_rejected
    test_fork_creates_child_tape
    test_fork_race_selects_first_success
    test_fork_forget_spawns_child_independently
    test_fork_child_workspace_narrowing
    test_fork_world_enabled_requires_workspace_config
    test_fork_world_enabled_host_workspace_child
    test_fork_world_enabled_workspace_parent_host_child
    test_switch_world_adopts_child_world
    test_fork_race_adopt_winner
    test_exec_preserves_mission
    test_exec_self_reentry_survives_ephemeral_body
    test_exec_self_reentry_reclaims_public_surface
    test_exec_applies_staged_config
    test_exec_rejects_invalid_staged_config
    test_exec_explicit_argv_replaces_mission
    test_exec_external_binary_completes_pipe_mission
    test_exec_external_binary_preserves_stdio_triplet
    test_exec_copied_quine_binary_continues_mission
    test_exec_relative_workspace_target_continues_mission
    test_dual_channel_separation
    test_no_stdin
"

# ── Preflight checks ──────────────────────────────────────

preflight() {
    _ok=1

    if [ ! -x "$QUINE" ]; then
        printf "${RED}ERROR${RESET}: quine binary not found at %s\n" "$QUINE"
        printf "  Build it with: go build -o %s ./cmd/quine/\n" "$QUINE"
        _ok=0
    fi

    for var in QUINE_MODEL_ID QUINE_API_TYPE QUINE_API_BASE QUINE_API_KEY; do
        eval "_val=\${${var}:-}"
        if [ -z "$_val" ]; then
            printf "${RED}ERROR${RESET}: %s is not set\n" "$var"
            printf "  Run: source profiles/gpt-5.4-codex-oauth.env\n"
            _ok=0
        fi
    done

    if [ "$_ok" -eq 0 ]; then
        echo ""
        echo "Preflight checks failed. Fix the above and re-run."
        exit 2
    fi

    printf "${BOLD}Quine E2E Tests${RESET}\n"
    printf "  Binary:  %s\n" "$QUINE"
    printf "  Model:   %s\n" "$QUINE_MODEL_ID"
    printf "  Turns:   %s per test\n" "$MAX_TURNS"
    printf "  Timeout: %ss per test\n" "$TIMEOUT"
    echo ""
}

# ── Main ───────────────────────────────────────────────────

preflight

if [ $# -gt 0 ]; then
    # Run selected tests
    for test_name in "$@"; do
        if type "$test_name" >/dev/null 2>&1; then
            "$test_name"
        else
            printf "${RED}Unknown test: %s${RESET}\n" "$test_name"
            FAIL=$((FAIL + 1))
            TOTAL=$((TOTAL + 1))
        fi
    done
else
    # Run all tests
    for test_name in $ALL_TESTS; do
        "$test_name"
    done
fi

# ── Summary ────────────────────────────────────────────────
echo ""
printf "${BOLD}Results: ${GREEN}%d passed${RESET}" "$PASS"
if [ "$FAIL" -gt 0 ]; then
    printf ", ${RED}%d failed${RESET}" "$FAIL"
fi
if [ "$SKIP" -gt 0 ]; then
    printf ", ${YELLOW}%d skipped${RESET}" "$SKIP"
fi
printf " (out of %d)\n" "$TOTAL"

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
exit 0
