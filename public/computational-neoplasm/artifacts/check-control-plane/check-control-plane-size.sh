#!/usr/bin/env bash
#
# VALIDATOR INTENTION
# purpose:
#   - surface control-plane file bloat: warn when any registered control-plane file exceeds the line threshold
#   - make accretion VISIBLE (a "pain signal") so bloat cannot accumulate silently, even when not immediately fixed
#   - enforce the TOTAL BUDGET RATCHET: the registry's `budget_lines:` frontmatter is a hard
#     ceiling on the summed line count of all registered files. Growth anywhere must be paid
#     for by compression elsewhere in the same change ("one in, one out"). This is the
#     forcing-function half the per-file pain line lacks: per-file thresholds incentivize
#     splitting (which resets counters without shrinking the surface); the total budget
#     cannot be gamed by splitting.
#   - keep the control-plane surface addressable via the registry
# repair-guidance:
#   - per-file threshold stays advisory: report drift, do not rewrite files; exit 0 unless
#     --strict AND something is over threshold
#   - the BUDGET breach is BLOCKING (exit 1) even without --strict — that is its point.
#     To land a change that grows the total: compress a registered file in the same change.
#     Raising `budget_lines:` itself is allowed only as a deliberate human-gated decision
#     (visible in the registry diff), never as a reflex to un-block a commit.
#   - when the total falls below the budget, ratchet DOWN: set `budget_lines:` to the new
#     total (the sensor prints the exact value). The budget only moves down or, by human
#     decision, up.
#   - the registry (control-plane-registry.md) is authoritative for the file list; paths are parsed from its table
#   - to change the pain line, pass --threshold N or set QUINE_CP_SIZE_THRESHOLD; do not hardcode it in logic
#   - this sensor does NOT prescribe what to compress — it only makes the size visible; the file owner decides

set -euo pipefail

REGISTRY="control-plane-registry.md"
THRESHOLD="${QUINE_CP_SIZE_THRESHOLD:-200}"
BUDGET="${QUINE_CP_BUDGET:-}"
STRICT=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --strict) STRICT=1; shift ;;
    --threshold) THRESHOLD="$2"; shift 2 ;;
    --budget) BUDGET="$2"; shift 2 ;;
    --registry) REGISTRY="$2"; shift 2 ;;
    -h|--help)
      echo "Usage: $0 [--strict] [--threshold N] [--budget N] [--registry FILE]"
      echo "  Warns when registered control-plane files exceed N lines (default 200)."
      echo "  Per-file check is advisory unless --strict (then exits 1 if any file is over)."
      echo "  Total-budget check is BLOCKING: exits 1 when the summed line count exceeds"
      echo "  the registry's budget_lines: frontmatter (override with --budget/QUINE_CP_BUDGET)."
      exit 0 ;;
    *) echo "unknown arg: $1" >&2; exit 2 ;;
  esac
done

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

if [[ ! -f "$REGISTRY" ]]; then
  echo "FAIL: registry not found: $REGISTRY" >&2
  exit 2
fi

# Budget = registry frontmatter `budget_lines: N` unless overridden.
if [[ -z "$BUDGET" ]]; then
  BUDGET="$(sed -n 's/^budget_lines:[[:space:]]*\([0-9][0-9]*\).*/\1/p' "$REGISTRY" | head -1)"
fi
if [[ -z "$BUDGET" ]]; then
  echo "FAIL: no budget_lines: frontmatter in $REGISTRY (and no --budget/QUINE_CP_BUDGET override)" >&2
  exit 2
fi

# Registered paths = first backtick-quoted token of each table data row.
mapfile -t FILES < <(sed -n 's/^| `\([^`]*\)`.*/\1/p' "$REGISTRY")

if [[ ${#FILES[@]} -eq 0 ]]; then
  echo "FAIL: no control-plane files parsed from $REGISTRY (table format drift?)" >&2
  exit 2
fi

over=0
missing=0
total=0
echo "Control-plane size — pain line ${THRESHOLD} lines (registry: ${REGISTRY}, ${#FILES[@]} files):"
for f in "${FILES[@]}"; do
  if [[ ! -f "$f" ]]; then
    echo "  MISSING (registry drift — file gone?): $f"
    missing=$((missing + 1))
    continue
  fi
  n=$(wc -l < "$f" | tr -d ' ')
  total=$((total + n))
  if (( n > THRESHOLD )); then
    echo "  OVER  ${n}  ${f}"
    over=$((over + 1))
  fi
done
echo "Total: ${#FILES[@]} files, ${total} lines. Budget: ${BUDGET} lines. Over threshold: ${over}. Missing: ${missing}."

if (( total > BUDGET )); then
  echo "FAIL (budget ratchet): control-plane total ${total} lines exceeds budget ${BUDGET}." >&2
  echo "  Growth must be paid for: compress a registered file in this same change." >&2
  echo "  Raising budget_lines: in ${REGISTRY} is a deliberate human-gated decision, not an un-block reflex." >&2
  exit 1
fi
if (( total < BUDGET )); then
  echo "  RATCHET-DOWN available: total ${total} < budget ${BUDGET} — set 'budget_lines: ${total}' in ${REGISTRY}."
fi

if [[ "$STRICT" == 1 && "$over" -gt 0 ]]; then
  echo "FAIL (--strict): ${over} control-plane file(s) over ${THRESHOLD} lines." >&2
  exit 1
fi
exit 0
