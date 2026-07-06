#!/bin/sh
# Shared infra (plan §2.4): embedded-directive pre-flight linter. Run as a LAUNCH GATE.
#
# Exits 1 if the seed tree carries an imperative / task-leak token or a readable ordinal
# name scheme that would smuggle an implicit instruction or a linguistic ordinal into a
# fixture that claims to be "structure-only". This makes the reflexive-capture discipline
# a check, not a reviewer note.
#
# Usage: lint_seed.sh <workspace_dir> [extra_forbidden_regex]
set -eu
ws="${1:?usage: lint_seed.sh <workspace_dir> [extra_regex]}"
extra="${2:-}"

# Imperative / task-leak tokens. Case-insensitive. NOTE: experiment *condition* names
# (gap/affordance/funnel/...) must never appear in seed CONTENT — they are caught here too.
# NOTE: directive/task-cue tokens only. Deliberately NOT general English words like
# "should"/"must"/"please" — those appear in benign natural-language fixtures (E5 diary
# prose) and would false-positive. Binaries are skipped (grep -I), since a compiled
# binary's internal string table is not a readable cue the agent perceives.
FORBID='missing|todo|fixme|fill[ _-]?in|reconcile|inconsistent|required|run[ _-]?me|complete[ _-]the|draft|placeholder|\bgap\b|fix[ _-]me|resolve'

# validated warrant wordings — elic-18/19/32 lineage
#
# Narrow, checked-in exemption: these are notice.txt bodies that were reviewed
# and confirmed to carry no imperative/task-leak content, but that legitimately
# contain a FORBID token (e.g. "required" in "no approval is required") as part
# of validated warrant/predicate wording. Exempted by exact sha256 of file
# content, one file at a time -- everything else still scans as before.
ALLOWLIST_SHA256='
70ccfacdd96bdd46bd5547e7683ed804da2b694538ed320d4ea1e5fc8b55c4b9
17f69c31c400c6e367fb00441e2cd140e471ff7d8d4bb07a30cc97658fc50046
'

is_allowlisted_notice() {  # $1 = path to a seed file named notice.txt
  [ "$(basename "$1")" = "notice.txt" ] || return 1
  h=$(sha256sum "$1" 2>/dev/null | awk '{print $1}')
  [ -n "$h" ] || return 1
  printf '%s\n' "$ALLOWLIST_SHA256" | grep -qx "$h"
}

fail=0
hits=$(grep -rEilI "$FORBID" "$ws" 2>/dev/null || true)
if [ -n "$(printf '%s' "$hits" | tr -d '[:space:]')" ]; then
  kept=""
  while IFS= read -r f; do
    [ -n "$f" ] || continue
    if is_allowlisted_notice "$f"; then
      continue
    fi
    kept="${kept}${f}
"
  done <<EOF
$hits
EOF
  hits="$kept"
fi
if [ -n "$extra" ]; then
  ehits=$(grep -rEilI "$extra" "$ws" 2>/dev/null || true)
  hits=$(printf '%s\n%s' "$hits" "$ehits")
fi
# Readable ordinal name SCHEMES in filenames: bare digits, part_a, step1, section2, item3.
ord=$(find "$ws" -maxdepth 4 \( \
        -name '[0-9]' -o -name '[0-9][0-9]' \
        -o -iname 'part[_-]?[a-z0-9]' -o -iname 'step[0-9]*' \
        -o -iname 'section[0-9]*' -o -iname 'item[0-9]*' \
        -o -iname 'chapter[0-9]*' \) 2>/dev/null || true)

if [ -n "$(printf '%s' "$hits" | tr -d '[:space:]')" ]; then
  echo "LINT FAIL: embedded-directive token in seed content:" >&2
  printf '%s\n' "$hits" | grep -v '^$' | sed 's/^/    /' >&2
  fail=1
fi
if [ -n "$(printf '%s' "$ord" | tr -d '[:space:]')" ]; then
  echo "LINT FAIL: readable ordinal filename scheme in seed:" >&2
  printf '%s\n' "$ord" | grep -v '^$' | sed 's/^/    /' >&2
  fail=1
fi
[ "$fail" -eq 0 ] && echo "LINT OK: no embedded-directive / ordinal cue in $ws"
exit "$fail"
