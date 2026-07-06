#!/bin/bash
set -uo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

ENV_FILE=".env.gpt-5.4-codex-medium"
REPLICATES=3
WALLCLOCK_SECONDS=300
JOBS=1
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --jobs) JOBS="$2"; shift 2 ;;
    --jobs=*) JOBS="${1#*=}"; shift ;;
    *) positional+=("$1"); shift ;;
  esac
done
[[ ${#positional[@]} -ge 1 ]] && ENV_FILE="${positional[0]}"
[[ ${#positional[@]} -ge 2 ]] && REPLICATES="${positional[1]}"
[[ ${#positional[@]} -ge 3 ]] && WALLCLOCK_SECONDS="${positional[2]}"

CONDITIONS=(passive_vague active_vague active_conflict)
total=$(( REPLICATES * ${#CONDITIONS[@]} ))
LOGDIR="$(mktemp -d /tmp/p14-06-override.XXXXXX)"
echo "=== 14.06 override pilot ==="
echo "Env file:    ${ENV_FILE}"
echo "Conditions:  ${CONDITIONS[*]}"
echo "Replicates:  ${REPLICATES} each (${total} runs)"
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
  echo "--- launch ${cond} r${rid} (${launched}/${total}) ---"
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
