#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

FILES=()

for path in \
  README.md \
  DEVELOPMENT.md \
  TESTING.md \
  AGENTS.md \
  STATUS.md \
  ROADMAP.md \
  QUICKSTART.md \
  tests/README.md \
  tests/model/README.md; do
  if [[ -f "$path" ]]; then
    FILES+=("$path")
  fi
done

if [[ "${#FILES[@]}" -eq 0 ]]; then
  echo "Public surface check skipped: no public-surface files found."
  exit 0
fi

CHECKS=(
  'tests/llm-behavior/:::stale model-test path; use tests/model/'
  'test/llm-behavior/:::stale legacy behavior-test path; use tests/model/'
  '\| `main` \| Curated public subset chosen from lab \|:::stale public-branch wording; main is now a public-facing repo'
  '`main` branch = curated public:::stale public-branch wording; main is now a public-facing repo'
  'Public branch \(selective sync from `lab`\):::stale branch description; main now carries code, tests, and root docs'
  'approved public writing:::stale public-surface wording; Paper/ is not part of the public-facing contract'
  'approved public `Paper/` materials:::stale public-surface wording; Paper/ is not part of the public-facing contract'
  'selected manuscript source:::stale public-surface wording; manuscript source is not part of the public-facing contract'
  'smallest clean subset that expresses the contribution:::stale teaser-only release wording'
)

failures=0

for entry in "${CHECKS[@]}"; do
  pattern="${entry%%:::*}"
  message="${entry#*:::}"
  if matches="$(rg -n --color=never "$pattern" "${FILES[@]}" || true)"; then
    if [[ -n "$matches" ]]; then
      printf 'FAIL: %s\n' "$message" >&2
      printf '%s\n\n' "$matches" >&2
      failures=1
    fi
  fi
done

if [[ "$failures" -ne 0 ]]; then
  echo "Public surface check failed." >&2
  exit 1
fi

echo "Public surface check passed."
