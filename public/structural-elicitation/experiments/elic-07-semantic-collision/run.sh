#!/bin/bash
# 13.07 Semantic Collision — minimal campaign (collision vs consistent x N).
#
# Runs the matched pair x N replicates through the single-condition runner. Each
# run is fully isolated (own run dir, workspace, runtime tree, container), so
# replicates may run concurrently. --jobs N caps concurrency (default 1).
#
# Usage:   ./run.sh [ENV_FILE] [REPLICATES] [WALLCLOCK_SECONDS] [--jobs N]
# Example: ./run.sh .env.gpt-5.4-codex-medium 3 240 --jobs 4
# Score:   python3 analysis/score.py runs/

set -uo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

ENV_FILE=".env.gpt-5.4-codex-medium"
REPLICATES=3
WALLCLOCK_SECONDS=240
JOBS=1
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --jobs)   JOBS="$2"; shift 2 ;;
    --jobs=*) JOBS="${1#*=}"; shift ;;
    *)        positional+=("$1"); shift ;;
  esac
done
[[ ${#positional[@]} -ge 1 ]] && ENV_FILE="${positional[0]}"
[[ ${#positional[@]} -ge 2 ]] && REPLICATES="${positional[1]}"
[[ ${#positional[@]} -ge 3 ]] && WALLCLOCK_SECONDS="${positional[2]}"

CONDITIONS=(collision_a collision_b consistent)
total=$(( REPLICATES * ${#CONDITIONS[@]} ))
LOGDIR="$(mktemp -d /tmp/p13-07-wave.XXXXXX)"
echo "=== 13.07 semantic-collision campaign ==="
echo "Env file:    ${ENV_FILE}"
echo "Conditions:  ${CONDITIONS[*]}"
echo "Replicates:  ${REPLICATES} each  (${total} runs)"
echo "Wallclock:   ${WALLCLOCK_SECONDS}s/run"
echo "Concurrency: ${JOBS}"
echo "Per-run logs: ${LOGDIR}/"
echo ""

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
    wait -n 2>/dev/null || true
    active=$(( active - 1 ))
  fi
done
wait || true

echo ""
echo "=== campaign complete (${total} runs) ==="
echo "Per-run logs: ${LOGDIR}/"
echo "Score with: python3 ${SCRIPT_DIR}/analysis/score.py ${SCRIPT_DIR}/runs/"
