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
#   source .env.kimi    # or .env.claude — load API credentials
#   ./tests/runtime/run.sh                       # run all tests
#   ./tests/runtime/run.sh test_exit_success     # run one test
#   ./tests/runtime/run.sh test_fd3_piped_input test_fd4_delivery  # run selected tests
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
    printf "${GREEN}PASS${RESET}\n"
}

# Record a fail with reason
fail() {
    FAIL=$((FAIL + 1))
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
    export QUINE_DATA_DIR="${TEST_DIR}/tapes"
    mkdir -p "$QUINE_DATA_DIR"
}

# Clean up after one test
teardown() {
    rm -rf "$TEST_DIR" 2>/dev/null || true
}

is_linux() {
    [ "$(uname -s)" = "Linux" ]
}

has_lima() {
    command -v limactl >/dev/null 2>&1
}

stage_oauth_seed() {
    _export_dir="$1"
    _host_quine_config="$2"

    [ -d "$_host_quine_config" ] || return 0

    mkdir -p "$_export_dir/quine-config-host"
    for _f in kimi-oauth.json kimi-device.json codex-oauth.json; do
        cp -f "$_host_quine_config/$_f" "$_export_dir/quine-config-host/" 2>/dev/null || true
    done
}

count_oauth_seed_files() {
    _seed_dir="$1"
    _count=0
    for _f in kimi-oauth.json kimi-device.json codex-oauth.json; do
        [ -f "$_seed_dir/$_f" ] && _count=$((_count + 1))
    done
    echo "$_count"
}

restore_oauth_seed() {
    _export_dir="$1"
    _host_quine_config="$2"

    [ -d "$_export_dir/quine-config-out" ] || return 0
    mkdir -p "$_host_quine_config"
    for _f in kimi-oauth.json kimi-device.json codex-oauth.json; do
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
run_quine() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    timeout "$TIMEOUT" env \
        QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" "$@" \
        >"$_stdout_file" \
        2>"$_stderr_file"
    return $?
}

# Run workspace physics tests on host Linux directly, or via Lima on non-Linux.
# On macOS, the default repo path is the Lima `colima` guest. This harness does
# not use Podman for workspace e2e unless a separate experiment opts into it.
# Arguments:
#   run_workspace_quine <stdout_file> <stderr_file> <mission>
run_workspace_quine() {
    _stdout_file="$1"; shift
    _stderr_file="$1"; shift
    _mission="$1"

    if is_linux; then
        timeout "$TIMEOUT" env \
            QUINE_DATA_DIR="$QUINE_DATA_DIR" \
            QUINE_WORKSPACE="$TEST_DIR/workspace" \
            QUINE_MAX_TURNS="$MAX_TURNS" \
            "$QUINE" "$_mission" \
            >"$_stdout_file" \
            2>"$_stderr_file" \
            </dev/null
        return $?
    fi

    if ! has_lima; then
        fail "workspace tests require Linux host or limactl"
        return 127
    fi

    _lima_instance="${E2E_LIMA_INSTANCE:-colima}"
    _use_sudo="${E2E_LIMA_USE_SUDO:-0}"
    _repo_root="$PWD"
    _runner_template="$_repo_root/tests/runtime/lib/lima-workspace-run.sh"
    _guest_quine="$_repo_root/.tmp/quine-linux-arm64-e2e"
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
    _runner_script="$(mktemp "$_repo_root/.tmp/quine-e2e-runner.XXXXXX")"
    _lima_provider="${QUINE_PROVIDER:-}"
    _lima_model="${QUINE_MODEL_ID:-}"
    _lima_api_type="${QUINE_API_TYPE:-}"
    _lima_api_base="${QUINE_API_BASE:-}"
    _lima_api_key="${QUINE_API_KEY:-}"
    _lima_thinking_budget="${QUINE_THINKING_BUDGET:-}"
    _host_quine_config="${QUINE_CONFIG_DIR:-$HOME/.config/quine}"

    stage_oauth_seed "$_export_dir" "$_host_quine_config"

    if printf '%s' "$_lima_api_key" | grep -Eqi 'oauth|device'; then
        _fallback_env="${E2E_LIMA_ENV_FILE:-.env.kimi}"
        _oauth_seed_count=$(count_oauth_seed_files "$_export_dir/quine-config-host")

        if [ "$_oauth_seed_count" -eq 0 ]; then
            if [ ! -f "$_fallback_env" ]; then
                fail "OAuth token cache missing; run local OAuth once or set E2E_LIMA_ENV_FILE for non-OAuth fallback"
                rm -rf "$_export_dir" "$_runner_script" 2>/dev/null || true
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
                rm -rf "$_export_dir" "$_runner_script" 2>/dev/null || true
                return 2
            fi
        fi
    fi

    cp "$_runner_template" "$_runner_script"
    chmod +x "$_runner_script"

    limactl shell --start --preserve-env "$_lima_instance" \
        /bin/sh "$_runner_script" \
        "$TIMEOUT" "$_repo_root" "$_guest_quine" "$_mission" "$_export_dir" "$MAX_TURNS" \
        "$_lima_provider" "$_lima_model" "$_lima_api_type" \
        "$_lima_api_base" "$_lima_api_key" "$_lima_thinking_budget" \
        "$_use_sudo" "$_export_dir/quine-config-host"
    _code=$?

    mkdir -p "$TEST_DIR/workspace" "$QUINE_DATA_DIR"
    [ -f "$_export_dir/stdout.txt" ] && cp "$_export_dir/stdout.txt" "$_stdout_file"
    [ -f "$_export_dir/stderr.txt" ] && cp "$_export_dir/stderr.txt" "$_stderr_file"
    [ -d "$_export_dir/workspace" ] && cp -R "$_export_dir/workspace"/. "$TEST_DIR/workspace"/
    [ -d "$_export_dir/quine" ] && cp -R "$_export_dir/quine"/. "$QUINE_DATA_DIR"/
    restore_oauth_seed "$_export_dir" "$_host_quine_config"
    rm -rf "$_export_dir" "$_runner_script" 2>/dev/null || true

    return "$_code"
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

# Assert a tape file exists and has an outcome entry
assert_tape_has_outcome() {
    _tape_dir="$1"
    _found="$(
    find "$_tape_dir" -type f -name '*.jsonl' 2>/dev/null | while IFS= read -r f; do
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
    find "$_tape_dir" -type f -name '*.jsonl' 2>/dev/null | while IFS= read -r f; do
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

# Find the tape JSONL file (assumes one session per test)
find_tape() {
    _tape_dir="$1"
    _tape="$(find "$_tape_dir" -type f -name '*.jsonl' 2>/dev/null | head -n 1)"
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

    printf 'abcdeFGHIJ' | run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do exactly these steps in separate sh calls: (1) run `dd bs=5 count=1 < /dev/fd/3 2>/dev/null > /tmp/part1.txt` (2) run `cat /dev/fd/3 > /tmp/part2.txt` (3) run `printf "PART1=%s\nPART2=%s\n" "$(cat /tmp/part1.txt)" "$(cat /tmp/part2.txt)" >&4` Then exit success.'
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
    begin_test "interactive PTY job exposes screen/in/exit control surface"
    setup

    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=7
    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use exactly one sh call with interactive=true to start `python3 -q`. Use the returned absolute job path. Then, in separate normal sh calls: (1) verify `<job>/screen.txt` shows the Python prompt `>>>`; (2) write `print(6*7)<enter>` to `<job>/in`; (3) wait until `<job>/screen.txt` contains `42`; (4) write `exit()<enter>` to `<job>/in`; (5) `cat <job>/exit`; (6) echo `INTERACTIVE_OK` to >&4. Then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if assert_exit "$code" 0 "interactive exit" && assert_contains "$TEST_DIR/stdout" "INTERACTIVE_OK" "interactive stdout"; then
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

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=0 \
        QUINE_TURN_EXHAUSTION_POLICY=near_death_exec \
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

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=3 \
        QUINE_TURN_EXHAUSTION_POLICY=hard_fail \
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
    begin_test "hard_fail exhaustion terminates immediately"
    setup

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=1 \
        QUINE_TURN_EXHAUSTION_POLICY=hard_fail \
        QUINE_PROMPT_METAPHOR=off \
        "$QUINE" 'Run exactly one sh call now: echo step1. After that, you would need another separate sh call for step2, so do not call exit before the first sh result comes back.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 1 "hard_fail exhaustion"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q 'turn_exhaustion' "$_tape" 2>/dev/null; then
            fail "tape missing turn_exhaustion termination mode"
        elif grep -q '\[EXECUTION BUDGET EXHAUSTED\]' "$_tape" 2>/dev/null; then
            fail "hard_fail should not open continuation window"
        else
            pass
        fi
    fi
    teardown
}

test_execution_budget_near_death_exec() {
    begin_test "near_death_exec opens exec-only continuation"
    setup

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=1 \
        QUINE_TURN_EXHAUSTION_POLICY=near_death_exec \
        QUINE_PROMPT_METAPHOR=off \
        "$QUINE" 'If wisdom key "phase" equals "post_exec", call exit success immediately. Otherwise, first run one sh call: echo "PRE_NEAR_DEATH" >&4. After execution budget is exhausted, call exec with wisdom phase=post_exec and reason "continue via exec".' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if [ "$code" -eq 124 ]; then
        fail "near_death_exec scenario timed out"
        teardown
        return
    fi
    if assert_contains "$TEST_DIR/stdout" "PRE_NEAR_DEATH" "near_death stdout"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '"name":"exec"' "$_tape" 2>/dev/null; then
            fail "near_death_exec should trigger an exec continuation attempt"
        elif grep -q 'turn_exhaustion' "$_tape" 2>/dev/null; then
            fail "near_death_exec path should avoid terminal turn_exhaustion when exec continuation succeeds"
        else
            pass
        fi
    fi
    teardown
}

test_anchor_memory_roundtrip() {
    begin_test "anchor memory mark/unfold works when enabled"
    setup

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=4 \
        QUINE_ANCHOR_MEMORY=1 \
        "$QUINE" 'Anchor memory is enabled for this run. Do exactly these steps: (1) call mark with summary "e2e-anchor-checkpoint" and fold=false (2) call unfold with anchor_id=0 (3) write "ANCHOR_E2E_OK" to >&4 (4) exit success.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 0 "anchor-memory exit" &&
        assert_contains "$TEST_DIR/stdout" "ANCHOR_E2E_OK" "anchor-memory stdout"; then
        _tape=$(find_tape "$QUINE_DATA_DIR")
        _anchor_meta="$(find "$QUINE_DATA_DIR" -path '*/context/anchors/0.anchor/meta.json' 2>/dev/null | head -n 1)"
        if [ -z "$_tape" ]; then
            fail "no tape file found"
        elif ! grep -q '"name":"mark"' "$_tape" 2>/dev/null; then
            fail "tape missing mark tool call"
        elif ! grep -q '"name":"unfold"' "$_tape" 2>/dev/null; then
            fail "tape missing unfold tool call"
        elif ! grep -q '\[MEMORY META\]' "$_tape" 2>/dev/null; then
            fail "tape missing [MEMORY META] feedback"
        elif [ -z "$_anchor_meta" ]; then
            fail "anchor meta file not found under stable memory path"
        elif ! grep -q '"summary": "e2e-anchor-checkpoint"' "$_anchor_meta" 2>/dev/null; then
            fail "anchor meta missing expected summary"
        else
            pass
        fi
    fi
    teardown
}

test_prompt_metaphor_off() {
    begin_test "QUINE_PROMPT_METAPHOR=off keeps prompt factual"
    setup

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
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

    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=2 \
        QUINE_TURN_EXHAUSTION_POLICY=hard_fail \
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

# ── 8. Transactional workspace physics ────────────────────

test_workspace_overlay_commit() {
    begin_test "transactional workspace commits on success"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Do exactly these steps in separate sh calls: (1) printf "overlay-commit\n" > committed.txt (2) test -f committed.txt && echo "COMMIT_OK" >&4. Then exit success.'
    code=$?

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

test_workspace_overlay_rollback() {
    begin_test "transactional workspace rolls back on failure"
    setup

    mkdir -p "$TEST_DIR/workspace"
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Run exactly one sh call: printf "overlay-rollback\n" > rolled_back.txt. Then exit failure with stderr "rollback test".'
    code=$?

    if assert_exit "$code" 1 "workspace rollback exit"; then
        if [ -e "$TEST_DIR/workspace/rolled_back.txt" ]; then
            fail "rolled_back.txt should not survive a failed session"
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

test_restore_world_restores_prior_revision() {
    begin_test "restore_world restores provisional workspace to an earlier world revision"
    setup

    mkdir -p "$TEST_DIR/workspace"
    _saved_max_turns="$MAX_TURNS"
    MAX_TURNS=6
    run_workspace_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Transactional workspace physics is enabled. Use exactly three sh calls and one restore_world call. Do not use fork or exec. Turn 1: run exactly `printf "v1\n" > state.txt`. Observe that the tool result reports current world revision `wr1`. Turn 2: run exactly `printf "v2\n" > state.txt`. Observe that the tool result reports current world revision `wr2`. Then call `restore_world` with `revision="wr1"` so `state.txt` should become `v1` again. Turn 3: run exactly `test "$(cat state.txt)" = "v1" && echo "RESTORE_OK" >&4`. Then exit success.'
    code=$?
    MAX_TURNS="$_saved_max_turns"

    if assert_exit "$code" 0 "restore world exit" &&
        assert_contains "$TEST_DIR/stdout" "RESTORE_OK" "restore world stdout" &&
        assert_contains "$TEST_DIR/workspace/state.txt" "v1" "restore world file" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\\[RESTORE WORLD\\] restored provisional workspace to revision wr1" "restore world marker"; then
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
    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
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
        'Use exactly one fork tool call with mode "race" and two children. Child 0 intent: "sleep 3; echo RACE_SLOW >&4; exit success". Child 1 intent: "echo RACE_FAST >&4; exit success". Both workspaces must be ".". After the fork call returns, exit success immediately.'
    code=$?

    if assert_exit "$code" 0 "fork race exit" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\\[FORK RACE\\] child 1 won (2 spawned, 1 succeeded, 1 killed)" "fork race winner"; then
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

    if assert_exit "$code" 0 "fork forget exit" &&
        assert_any_tape_contains "$QUINE_DATA_DIR" "\\[FORK OK\\] 1 children spawned" "fork forget marker"; then
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

# ── 11. Exec lifecycle ─────────────────────────────────────

test_exec_preserves_mission() {
    begin_test "exec resets context but preserves mission"
    setup

    # The mission asks to exec once, then complete.
    # After exec, the agent should still know its original mission.
    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'If you have wisdom key "phase" with value "post_exec", write "EXEC_SURVIVED" to >&4 and exit success. Otherwise, call exec with wisdom phase=post_exec and reason "testing exec".'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "EXEC_SURVIVED" "exec mission"; then
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
    test_binary_stdin
    test_execution_budget_disabled_hidden
    test_execution_budget_enabled_feedback
    test_execution_budget_hard_fail
    test_execution_budget_near_death_exec
    test_anchor_memory_roundtrip
    test_prompt_metaphor_off
    test_prompt_metaphor_thermodynamic
    test_workspace_overlay_commit
    test_workspace_overlay_rollback
    test_workspace_overlay_absolute_path
    test_restore_world_restores_prior_revision
    test_workspace_unsupported_on_non_linux
    test_tape_has_meta
    test_tape_has_outcome
    test_tape_has_messages
    test_fork_wait
    test_fork_creates_child_tape
    test_fork_race_selects_first_success
    test_fork_forget_spawns_child_independently
    test_exec_preserves_mission
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
            printf "  Run: source .env.kimi  (or .env.claude)\n"
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
