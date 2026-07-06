#!/bin/sh
# Shared infra (plan §2.3): opaque, uniform-length / uniform-charset filename tokens.
#
# A homogeneous *set* (which invites "is one missing?") that carries NO readable ordinal,
# slug, or condition cue. Deterministic given a seed integer (the harness must not use
# Math.random / $RANDOM so replicates are reproducible). Used by E4 (and E6/E5).
#
# Usage: neutral_names.sh <count> [seed]   # prints <count> 6-hex-char tokens, one per line
#
# Collisions: 6 hex chars = 16.7M space; for the small counts here a collision is
# vanishingly unlikely, but the caller SHOULD assert uniqueness (sort -u | wc -l == count)
# and bump the seed if it ever trips.
set -eu
count="${1:?usage: neutral_names.sh <count> [seed]}"
seed="${2:-0}"
i=0
while [ "$i" -lt "$count" ]; do
  # opaque token = first 6 hex of md5(seed:index); uniform length+charset, no ordinal.
  printf '%s:%s' "$seed" "$i" | md5sum | cut -c1-6
  i=$((i + 1))
done
