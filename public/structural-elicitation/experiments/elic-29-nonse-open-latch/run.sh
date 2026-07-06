#!/bin/bash
# 14.02 Non-SE Substrate Generalization (E5 / Gate 2) — campaign driver.
#
# Runs a matched-arm family x N replicates through the single-condition runner. Each
# run is fully isolated (own run dir, workspace, runtime tree, container), so
# replicates may run concurrently. --jobs N caps concurrency (default 1).
#
# A FAMILY selects which set of matched arms to run:
#   conflict  (PRIMARY, pre-registered):  conflict_a conflict_b consistent          (n=8-10/arm)
#   nonchore  (REQUIRED valence control):  nonchore_a nonchore_b nonchore_consistent (n=8-10/arm)
#   script    (SECONDARY, 13.04 port):     script_affordance script_inert script_pointless (n>=5/arm)
#   all:      conflict + nonchore + script
#
# Usage:   ./run.sh [ENV_FILE] [REPLICATES] [WALLCLOCK_SECONDS] [--family F] [--jobs N]
# Example: ./run.sh .env.gpt-5.4-codex-medium 10 240 --family conflict --jobs 4
# Score:   python3 analysis/score.py runs/

set -uo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

ENV_FILE=".env.gpt-5.4-codex-medium"
REPLICATES=10
WALLCLOCK_SECONDS=240
JOBS=1
FAMILY="conflict"
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --jobs)     JOBS="$2"; shift 2 ;;
    --jobs=*)   JOBS="${1#*=}"; shift ;;
    --family)   FAMILY="$2"; shift 2 ;;
    --family=*) FAMILY="${1#*=}"; shift ;;
    *)          positional+=("$1"); shift ;;
  esac
done
[[ ${#positional[@]} -ge 1 ]] && ENV_FILE="${positional[0]}"
[[ ${#positional[@]} -ge 2 ]] && REPLICATES="${positional[1]}"
[[ ${#positional[@]} -ge 3 ]] && WALLCLOCK_SECONDS="${positional[2]}"

DANGLING_ARMS=(dangling no_ref)            # F1 PRIMARY: dangling cross-ref vs symmetric no-ref control
BIB_ARMS=(bib_dangling bib_complete)       # F2 SECONDARY: dangling citation
CSV_ARMS=(csv_broken csv_complete)         # F3 EXPLORATORY: broken column-count row

case "${FAMILY}" in
  dangling) CONDITIONS=("${DANGLING_ARMS[@]}") ;;
  bib)      CONDITIONS=("${BIB_ARMS[@]}") ;;
  csv)      CONDITIONS=("${CSV_ARMS[@]}") ;;
  all)      CONDITIONS=("${DANGLING_ARMS[@]}" "${BIB_ARMS[@]}" "${CSV_ARMS[@]}") ;;
  *) echo "ERROR: unknown --family: ${FAMILY} (expected dangling|bib|csv|all)" >&2; exit 2 ;;
esac

total=$(( REPLICATES * ${#CONDITIONS[@]} ))
LOGDIR="$(mktemp -d /tmp/p14-02-wave.XXXXXX)"
echo "=== 14.05 non-SE open-latch campaign (E9 / Gate 2 fair retest) ==="
echo "Env file:    ${ENV_FILE}"
echo "Family:      ${FAMILY}"
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
