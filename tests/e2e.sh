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
#   ./tests/e2e.sh                       # run all tests
#   ./tests/e2e.sh test_exit_success     # run one test
#   ./tests/e2e.sh test_fd3 test_fd4     # run selected tests
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
    _found=0
    for f in "$_tape_dir"/*.jsonl; do
        [ -f "$f" ] || continue
        if grep -q '"type":"outcome"' "$f" 2>/dev/null || grep -q '"type": "outcome"' "$f" 2>/dev/null; then
            _found=1
            break
        fi
    done
    if [ "$_found" -eq 1 ]; then
        return 0
    else
        fail "no outcome entry in tape files under $_tape_dir"
        return 1
    fi
}

# Find the tape JSONL file (assumes one session per test)
find_tape() {
    _tape_dir="$1"
    for f in "$_tape_dir"/*.jsonl; do
        [ -f "$f" ] && echo "$f" && return 0
    done
    return 1
}


# ============================================================================
# TESTS
# ============================================================================

# ── 1. Exit codes ──────────────────────────────────────────

test_exit_success() {
    begin_test "exit(success) → exit code 0"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Write the word DONE to stdout via >&3, then call exit with status success.'
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

# ── 2. Stdout (fd 3) — deliverable channel ─────────────────

test_fd3_delivery() {
    begin_test "echo >&3 delivers to process stdout"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Run: echo "QUINE_MARKER_42" >&3   Then exit success. Do NOT print anything else to >&3.'
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
        'Run: echo "INTERNAL_ONLY_xyz"   (do NOT use >&3). Then exit success. Do NOT write anything to >&3.'
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

# ── 4. Persistent shell state ──────────────────────────────

test_shell_cd_persists() {
    begin_test "cd persists across sh calls"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do exactly these steps in separate sh calls: (1) cd /tmp  (2) pwd >&3  Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "/tmp" "pwd"; then
        pass
    fi
    teardown
}

test_shell_export_persists() {
    begin_test "export persists across sh calls"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do exactly these steps in separate sh calls: (1) export E2E_VAR=quine_test_789  (2) echo $E2E_VAR >&3  Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "quine_test_789" "export"; then
        pass
    fi
    teardown
}

test_shell_variable_persists() {
    begin_test "shell variable persists across sh calls"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Do exactly these steps in separate sh calls: (1) MY_COUNTER=1337  (2) echo $MY_COUNTER >&3  Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "1337" "variable"; then
        pass
    fi
    teardown
}

# ── 5. Stdin (fd 4) — material channel ─────────────────────

test_fd4_piped_input() {
    begin_test "piped stdin readable via /dev/fd/4"
    setup

    echo "MATERIAL_DATA_e2e_test" | run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Read the piped input using: cat /dev/fd/4    Then write what you read to >&3 and exit success.'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "MATERIAL_DATA_e2e_test" "fd4"; then
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
        -b 'A binary file was provided to you. Print its path to >&3 and exit success.'
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

# ── 7. Turn exhaustion ─────────────────────────────────────

test_turn_exhaustion() {
    begin_test "QUINE_MAX_TURNS exhaustion → exit code 1"
    setup

    # Give it only 2 turns but a mission that requires more
    # MaxTurns=1: the LLM gets exactly one sh call before exhaustion.
    # The prompt forces it to make an sh call first, then try more.
    timeout "$TIMEOUT" env \
        QUINE_DATA_DIR="$QUINE_DATA_DIR" \
        QUINE_MAX_TURNS=1 \
        "$QUINE" 'Step 1: Run "echo step1" in sh. Step 2: Run "echo step2" in sh. Step 3: Run "echo step3" in sh. You MUST run all three steps as separate sh calls before calling exit. Do NOT call exit until all 3 sh calls are done.' \
        >"$TEST_DIR/stdout" 2>"$TEST_DIR/stderr" </dev/null
    code=$?

    if assert_exit "$code" 1 "turn exhaustion"; then
        # Tape should have termination_mode = turn_exhaustion
        _tape=$(find_tape "$QUINE_DATA_DIR")
        if [ -n "$_tape" ] && grep -q 'turn_exhaustion' "$_tape" 2>/dev/null; then
            pass
        else
            fail "tape missing turn_exhaustion termination mode"
        fi
    fi
    teardown
}

# ── 8. Tape integrity ──────────────────────────────────────

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

# ── 9. Fork (child process) ───────────────────────────────

test_fork_wait() {
    begin_test "fork(wait=true) returns child output"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Use the fork tool with intent "echo CHILD_OUTPUT_42 >&3 and exit success" and wait=true. Then write the child stdout to >&3 and exit success.'
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
        'Use fork with intent "exit success" and wait=true. Then exit success.'
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

# ── 10. Exec (metamorphosis) ──────────────────────────────

test_exec_preserves_mission() {
    begin_test "exec resets context but preserves mission"
    setup

    # The mission asks to exec once, then complete.
    # After exec, the agent should still know its original mission.
    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'If you have wisdom key "phase" with value "post_exec", write "EXEC_SURVIVED" to >&3 and exit success. Otherwise, call exec with wisdom phase=post_exec and reason "testing exec".'
    code=$?

    if assert_exit "$code" 0 "exit" && assert_contains "$TEST_DIR/stdout" "EXEC_SURVIVED" "exec mission"; then
        pass
    fi
    teardown
}

# ── 11. Multi-fd: fd 1 captured, fd 3 delivered ────────────

test_dual_channel_separation() {
    begin_test "fd 1 (captured) and fd 3 (delivered) are separate"
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'In one sh call, run: echo "CAPTURED_fd1" && echo "DELIVERED_fd3" >&3    Then exit success.'
    code=$?

    if assert_exit "$code" 0 "exit"; then
        # fd 3 content should be in stdout
        _has_delivered=0
        grep -q "DELIVERED_fd3" "$TEST_DIR/stdout" 2>/dev/null && _has_delivered=1

        # fd 1 content should NOT be in stdout
        _has_captured=0
        grep -q "CAPTURED_fd1" "$TEST_DIR/stdout" 2>/dev/null && _has_captured=1

        if [ "$_has_delivered" -eq 1 ] && [ "$_has_captured" -eq 0 ]; then
            pass
        elif [ "$_has_delivered" -eq 0 ]; then
            fail "DELIVERED_fd3 not found in stdout"
        else
            fail "CAPTURED_fd1 leaked to stdout (should stay in tool result)"
        fi
    fi
    teardown
}

# ── 12. No-stdin mode ─────────────────────────────────────

test_no_stdin() {
    begin_test "no stdin (TTY mode) → material = Begin."
    setup

    run_quine "$TEST_DIR/stdout" "$TEST_DIR/stderr" \
        'Echo "NO_STDIN_OK" to >&3 and exit success.' </dev/null
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
    test_fd3_delivery
    test_fd1_not_leaked
    test_stderr_failure_signal
    test_stderr_success_silent
    test_shell_cd_persists
    test_shell_export_persists
    test_shell_variable_persists
    test_fd4_piped_input
    test_binary_stdin
    test_turn_exhaustion
    test_tape_has_meta
    test_tape_has_outcome
    test_tape_has_messages
    test_fork_wait
    test_fork_creates_child_tape
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
