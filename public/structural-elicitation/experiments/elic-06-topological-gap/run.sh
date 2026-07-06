#!/bin/bash
# 13.01 Incompleteness Elicitation — minimal campaign.
#
# The minimal experiment for the Structural-Elicitation claim is ONE factor
# (incompleteness) with two MISSIONLESS conditions: gap vs complete. This script
# runs that pair x N replicates through the canonical single-condition runner.
#
# Each run is fully isolated (its own run dir, workspace, runtime tree, and
# container), so replicates may run concurrently. Use --jobs N to cap how many
# run at once (default 1 = sequential). Replicates are emitted in interleaved
# order (r01 gap, r01 complete, r02 gap, ...) so a provider hiccup does not land
# entirely on one condition.
#
# Usage:
#   ./run.sh [ENV_FILE] [REPLICATES] [WALLCLOCK_SECONDS] [--directed] [--jobs N]
#
# Examples:
#   ./run.sh .env.gpt-5.4-codex-medium 5 240 --jobs 5   # 10 runs, up to 5 at once
#   ./run.sh .env.gpt-5.4-codex-medium 10 240           # 20 runs, sequential
#   ./run.sh .env.gpt-5.4-codex-medium 5 240 --directed # + the optional ceiling pair
#
# Then score:
#   python3 analysis/score.py runs/

set -uo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Positional + flag parsing.
ENV_FILE=".env.gpt-5.4-codex-medium"
REPLICATES=10
WALLCLOCK_SECONDS=240
JOBS=1
DIRECTED=0
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --directed) DIRECTED=1; shift ;;
    --jobs)     JOBS="$2"; shift 2 ;;
    --jobs=*)   JOBS="${1#*=}"; shift ;;
    *)          positional+=("$1"); shift ;;
  esac
done
[[ ${#positional[@]} -ge 1 ]] && ENV_FILE="${positional[0]}"
[[ ${#positional[@]} -ge 2 ]] && REPLICATES="${positional[1]}"
[[ ${#positional[@]} -ge 3 ]] && WALLCLOCK_SECONDS="${positional[2]}"

CONDITIONS=(dangling closed)
: # (no directed variant)

total=$(( REPLICATES * ${#CONDITIONS[@]} ))
LOGDIR="$(mktemp -d /tmp/p13-wave.XXXXXX)"
echo "=== 13.01 minimal campaign ==="
echo "Env file:    ${ENV_FILE}"
echo "Conditions:  ${CONDITIONS[*]}"
echo "Replicates:  ${REPLICATES} each  (${total} runs)"
echo "Wallclock:   ${WALLCLOCK_SECONDS}s/run"
echo "Concurrency: ${JOBS}"
echo "Per-run logs: ${LOGDIR}/"
echo ""

# Build the (condition, replicate) task list in interleaved order.
TASKS=()
for r in $(seq 1 "${REPLICATES}"); do
  rid="$(printf '%02d' "${r}")"
  for cond in "${CONDITIONS[@]}"; do
    TASKS+=("${cond} ${rid}")
  done
done

active=0
launched=0
for task in "${TASKS[@]}"; do
  set -- ${task}; cond="$1"; rid="$2"
  launched=$(( launched + 1 ))
  echo "--- launch ${cond} r${rid}  (${launched}/${total}) ---"
  (
    REPLICATE="${rid}" "${SCRIPT_DIR}/run-container.sh" "${ENV_FILE}" "${cond}" "${WALLCLOCK_SECONDS}" \
      > "${LOGDIR}/${cond}-r${rid}.log" 2>&1
  ) &
  active=$(( active + 1 ))
  if (( active >= JOBS )); then
    wait -n 2>/dev/null || true   # a slot freed; one run finished (success or failure)
    active=$(( active - 1 ))
  fi
done
# Drain remaining.
wait || true

echo ""
echo "=== campaign complete (${total} runs) ==="
echo "Per-run logs: ${LOGDIR}/"
echo "Score with: python3 ${SCRIPT_DIR}/analysis/score.py ${SCRIPT_DIR}/runs/"
