#!/bin/bash
# generate-world.sh — Create a budgeted world.json with random cell values.
#
# Usage: generate-world.sh <world_json_path> <state_dir> [cells] [budget] [agent_get_limit] [spec_state_dir]
#   world_json_path: where to write the spec (should be OUTSIDE workspace)
#   state_dir:       host directory for runtime state (budget, events)
#   cells:           number of cells (default: 20)
#   budget:          shared budget (default: 22)
#   agent_get_limit: optional per-agent get call cap per reset epoch
#   spec_state_dir:  optional state_dir value written into world.json; defaults
#                    to <state_dir>. Use this when the runtime should resolve
#                    the state through a sandbox-visible shared mount instead of
#                    the host absolute path.
set -euo pipefail

WORLD_JSON="${1:?world_json_path required}"
STATE_DIR="${2:?state_dir required}"
CELLS="${3:-20}"
BUDGET="${4:-22}"
AGENT_GET_LIMIT="${5:-0}"
SPEC_STATE_DIR="${6:-${STATE_DIR}}"

mkdir -p "$(dirname "${WORLD_JSON}")" "${STATE_DIR}"

# Generate random hex values for each cell.
items="{"
for i in $(seq 1 "${CELLS}"); do
  id=$(printf 'c%02d' "$i")
  val=$(openssl rand -hex 4)
  if [ "$i" -gt 1 ]; then items+=","; fi
  items+="\"${id}\":\"${val}\""
done
items+="}"

# Build the spec JSON.
config_body=$(cat <<EOF
    "budget": ${BUDGET},
    "cells": ${CELLS},
    "state_dir": "${SPEC_STATE_DIR}"
EOF
)
if [[ "${AGENT_GET_LIMIT}" -gt 0 ]]; then
  config_body+=$(printf ',\n    "agent_get_limit": %s' "${AGENT_GET_LIMIT}")
fi
cat > "${WORLD_JSON}" <<EOF
{
  "items": ${items},
  "config": {
${config_body}
  }
}
EOF

echo "world spec written to ${WORLD_JSON}"
echo "  cells:  ${CELLS}"
echo "  budget: ${BUDGET}"
if [[ "${AGENT_GET_LIMIT}" -gt 0 ]]; then
  echo "  per-agent get limit: ${AGENT_GET_LIMIT}"
fi
echo "  state (host): ${STATE_DIR}"
echo "  state (spec): ${SPEC_STATE_DIR}"
