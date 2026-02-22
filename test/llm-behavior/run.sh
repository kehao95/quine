#!/bin/bash
# LLM Behavior Test Runner
#
# Usage:
#   ./test/llm-behavior/run.sh jobs               # run one scenario
#   ./test/llm-behavior/run.sh all                # run all scenarios
#   ./test/llm-behavior/run.sh jobs gpt-4o        # specify model
#
# Requires:
#   - /tmp/quine built (go build -o /tmp/quine ./cmd/quine/)
#   - QUINE_MODEL_ID, QUINE_API_KEY, etc. in environment (source .env.*)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
QUINE="${QUINE:-/tmp/quine}"
MAX_TURNS="${QUINE_MAX_TURNS:-20}"

SCENARIO="${1:-all}"
MODEL_OVERRIDE="${2:-}"

# ── Helpers ──────────────────────────────────────────────────

die() { echo "FATAL: $*" >&2; exit 1; }

check_prereqs() {
    [[ -x "$QUINE" ]] || die "quine binary not found at $QUINE (run: go build -o /tmp/quine ./cmd/quine/)"
    [[ -n "${QUINE_MODEL_ID:-}" ]] || die "QUINE_MODEL_ID not set (run: source .env.*)"
    [[ -n "${QUINE_API_KEY:-}" ]] || die "QUINE_API_KEY not set (run: source .env.*)"
}

# Run one scenario, return 0 if all checks pass
run_scenario() {
    local name="$1"
    local prompt_file="$SCRIPT_DIR/prompts/${name}.md"
    [[ -f "$prompt_file" ]] || die "prompt not found: $prompt_file"

    local model="${MODEL_OVERRIDE:-$QUINE_MODEL_ID}"
    local model_short="${model##*-}"
    local runid="$(date +%Y%m%d-%H%M%S)-${model_short}"
    local run_dir="$SCRIPT_DIR/runs/${name}/${runid}"

    mkdir -p "$run_dir/workspace"

    echo "━━━ Scenario: ${name} ━━━"
    echo "  Model:  ${model}"
    echo "  Run ID: ${runid}"
    echo "  Dir:    ${run_dir}"
    echo ""

    # Copy prompt for traceability
    cp "$prompt_file" "$run_dir/prompt-used.md"

    # Run quine — mission from argv, no stdin (TTY mode → "Begin.")
    local mission
    mission="$(cat "$prompt_file")"
    cd "$run_dir/workspace"
    QUINE_DATA_DIR=".quine" \
    QUINE_MAX_TURNS="$MAX_TURNS" \
        "$QUINE" "$mission" < /dev/null \
        > "../stdout.txt" \
        2> "../stderr.txt" \
    || true  # don't fail on non-zero exit

    # Copy tape files out for easier access
    if [[ -d .quine ]]; then
        local jsonl=(.quine/*.jsonl)
        local logf=(.quine/*.log)
        [[ -f "${jsonl[0]}" ]] && cp "${jsonl[0]}" "$run_dir/tape.jsonl"
        [[ -f "${logf[0]}" ]] && cp "${logf[0]}" "$run_dir/tape.log"
    fi

    cd "$SCRIPT_DIR"

    # Score
    score_scenario "$name" "$run_dir"
    local score_result=$?

    # Update latest symlink
    ln -sfn "${name}/${runid}" "$SCRIPT_DIR/runs/latest"

    echo ""
    return $score_result
}

# Score a completed run. Prints report, returns 0 if all pass.
score_scenario() {
    local name="$1"
    local run_dir="$2"
    local stdout="$run_dir/stdout.txt"
    local score_file="$run_dir/score.txt"
    local pass=0
    local fail=0
    local total=0

    echo "── Scoring ──" | tee "$score_file"

    case "$name" in
        jobs)
            check_marker "$stdout" "NORMAL_OK"  "C1: Normal completion (no spurious PAUSED)"  "$score_file"
            check_marker "$stdout" "PAUSE_OK"   "C2: output_limit triggers [PAUSED]"          "$score_file"
            check_marker "$stdout" "RESUME_OK"  "C3: job(signal=cont) resumes to completion"  "$score_file"
            check_marker "$stdout" "KILL_OK"    "C4: job(signal=kill) terminates job"         "$score_file"
            check_marker "$stdout" "READ_OK"    "C5: job() reads output without resuming"     "$score_file"

            # C6: clean exit
            if [[ -f "$run_dir/stderr.txt" ]] && [[ ! -s "$run_dir/stderr.txt" ]]; then
                echo "  PASS  C6: Clean exit (no stderr)" | tee -a "$score_file"
            else
                echo "  WARN  C6: Stderr present (review tape)" | tee -a "$score_file"
            fi
            ;;
        *)
            echo "  No scoring rules for scenario: $name" | tee -a "$score_file"
            ;;
    esac

    # Count results
    pass=$(grep -c "PASS" "$score_file" || true)
    fail=$(grep -c "FAIL" "$score_file" || true)
    total=$((pass + fail))

    echo "" | tee -a "$score_file"
    echo "Result: ${pass}/${total} passed" | tee -a "$score_file"

    # Tape analysis hint
    if [[ -f "$run_dir/tape.log" ]]; then
        local turns=$(grep -c "^quine\[.*\]: turn " "$run_dir/tape.log" || true)
        local session_line=$(grep "session ended" "$run_dir/tape.log" || true)
        echo "Turns: ~${turns}" | tee -a "$score_file"
        echo "Session: ${session_line##*: }" | tee -a "$score_file"
    fi

    echo "" | tee -a "$score_file"
    echo "Review: cat $run_dir/tape.log" | tee -a "$score_file"

    [[ $fail -eq 0 ]]
}

# Check if a marker string appears in stdout
check_marker() {
    local file="$1"
    local marker="$2"
    local label="$3"
    local score_file="$4"

    if grep -q "$marker" "$file" 2>/dev/null; then
        echo "  PASS  ${label}" | tee -a "$score_file"
    else
        echo "  FAIL  ${label} (marker '${marker}' not found in stdout)" | tee -a "$score_file"
    fi
}

# ── Main ─────────────────────────────────────────────────────

check_prereqs

all_passed=true

if [[ "$SCENARIO" == "all" ]]; then
    for prompt in "$SCRIPT_DIR"/prompts/*.md; do
        name="$(basename "$prompt" .md)"
        run_scenario "$name" || all_passed=false
    done
else
    run_scenario "$SCENARIO" || all_passed=false
fi

if $all_passed; then
    echo "✓ All scenarios passed"
    exit 0
else
    echo "✗ Some scenarios failed — review tapes"
    exit 1
fi
