#!/usr/bin/env bash
# Fixture runner for prediction-test dispatcher smoke. No live model.
set -euo pipefail
EXP="$(cd "$(dirname "$0")" && pwd)"
RUNID="${RUNID:-$(date -u +%Y%m%d-%H%M%S)-fixture}"
mkdir -p "$EXP/runs/$RUNID"
echo ok > "$EXP/runs/$RUNID/out.txt"
exit 0
