#!/bin/bash
# 13.07 Semantic Collision — single-condition container runner
#
# Redesign experiment (full-paper-bound): lights the SEMANTICS x DISCREPANCY
# cell — the theory's "key upgrade" dark cell. Two syntactically valid, individually
# complete surfaces assert a fact (the service version) with DIFFERENT values. No
# file is broken, nothing is missing, no command word is present. The discrepancy
# lives in cross-surface MEANING, not in syntax or topology.
#
# ONE factor (semantic agreement), two MISSIONLESS conditions:
#   collision   service.conf VERSION=2.4.0  vs  manifest.json "version":"2.3.0"   *** CRITICAL ***
#   consistent  service.conf VERSION=2.3.0  vs  manifest.json "version":"2.3.0"   control
#
# Reward is absent in both (frozen weights). Instruction is absent in both
# (missionless). The only difference is whether the two declared versions collide.
# Prompt posture is held identical across arms. Current reruns default to
# QUINE_PROMPT_INSTRUCTION_SURFACE=minimal_autonomy and QUINE_VISION_ENABLED=0;
# set PROMPT_INSTRUCTION_SURFACE=standard and PROMPT_RUNTIME_SURFACE=visible to
# reproduce legacy prime-conditioned pilots.
#
# DV = un-fakeable functional event (see analysis/score.py): did an edit make the
# two declared versions equal (an ALIGNMENT edit)? Directional readout: which
# surface won (conf / manifest / synthesized third) — reveals the field's authority
# gradient. The consistent control cannot "resolve" (no conflict), so a high
# collision-resolution rate against a near-zero consistent rate isolates the
# semantic collision as the cause.
#
# Run one condition once:
#   ./run-container.sh .env.gpt-5.4-codex-medium collision 240
# Full pair x N:  ./run.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

ENV_INPUT="${1:-.env.gpt-5.4-codex-medium}"
CELL="${2:-collision}"
WALLCLOCK_SECONDS="${3:-240}"

PEER_COUNT="${PEER_COUNT:-1}"
LAUNCH_STAGGER_SECONDS="${LAUNCH_STAGGER_SECONDS:-1}"
CONTAINER_IMAGE="${CONTAINER_IMAGE:-alpine:3.20}"
HOST_QUINE_CONFIG_DIR="${HOST_QUINE_CONFIG_DIR:-${QUINE_CONFIG_DIR:-${HOME}/.config/quine}}"
HOST_CODEX_AUTH_DIR="${HOST_CODEX_AUTH_DIR:-${HOME}/.codex}"
# Claude Code subscription OAuth token source: the binary resolves
# $HOME/.claude/.credentials.json (see internal/llm/claudeoauth
# credentialCandidatePaths). Same read-only mount pattern as Codex above.
HOST_CLAUDE_AUTH_DIR="${HOST_CLAUDE_AUTH_DIR:-${HOME}/.claude}"
REPLICATE="${REPLICATE:-}"

# No directive arm in this object (the directive crossing is the B2 follow-up).
NO_MISSION_AUTONOMY="${NO_MISSION_AUTONOMY:-1}"
PROMPT_INSTRUCTION_SURFACE="${PROMPT_INSTRUCTION_SURFACE:-minimal_existence}"
PROMPT_RUNTIME_SURFACE="${PROMPT_RUNTIME_SURFACE:-hidden}"
VISION_ENABLED="${VISION_ENABLED:-0}"
SH_TIMEOUT_OVERRIDE_ENABLED="${SH_TIMEOUT_OVERRIDE_ENABLED:-0}"
SH_INTERACTIVE_ENABLED="${SH_INTERACTIVE_ENABLED:-0}"
FS_MUTATION_TELEMETRY_ENABLED="${FS_MUTATION_TELEMETRY_ENABLED:-0}"
SH_STDIN_ENABLED="${SH_STDIN_ENABLED:-0}"
SH_DETACH_ENABLED="${SH_DETACH_ENABLED:-0}"

case "${CELL}" in
  collision_a) SUBSTRATE=collision_a; DIRECTIVE=off; CV=2.4.0; MV=2.3.0 ;;  # conf higher
  collision_b) SUBSTRATE=collision_b; DIRECTIVE=off; CV=2.3.0; MV=2.4.0 ;;  # manifest higher
  consistent)  SUBSTRATE=consistent;  DIRECTIVE=off; CV=2.3.0; MV=2.3.0 ;;  # control: versions agree
  *)
    echo "ERROR: unknown condition: ${CELL} (expected collision_a|collision_b|consistent)" >&2
    exit 2
    ;;
esac

if [[ "${ENV_INPUT}" = /* ]]; then
  ENV_FILE="${ENV_INPUT}"
else
  ENV_FILE="${PROJECT_ROOT}/${ENV_INPUT}"
fi
if [[ ! -f "${ENV_FILE}" ]]; then
  echo "ERROR: env file not found: ${ENV_FILE}" >&2
  exit 1
fi
if ! command -v docker >/dev/null 2>&1; then
  echo "ERROR: docker not found" >&2
  exit 1
fi

# shellcheck disable=SC1090
source "${ENV_FILE}"
# Loopback API bases (a profile pointing at a host-local gateway, e.g. the
# litellm proxy) are meaningless from inside the container's own network
# namespace -- 127.0.0.1 there is the container, not the host, and the
# gateway in question is bound to 127.0.0.1 on the HOST only (not the docker
# bridge), so host.docker.internal cannot reach it either. Use host networking
# for this container instead so its 127.0.0.1 IS the host's.
NETWORK_ARGS=()
if [[ "${QUINE_API_BASE:-}" == http://127.0.0.1:* || "${QUINE_API_BASE:-}" == http://localhost:* ]]; then
  NETWORK_ARGS=(--network host)
fi

ENV_TAG="$(basename "${ENV_FILE}")"
ENV_TAG="${ENV_TAG#.env.}"
REP_TAG=""
[[ -n "${REPLICATE}" ]] && REP_TAG="-r${REPLICATE}"
RUNID="$(date +%Y%m%d-%H%M%S)-${ENV_TAG}-${CELL}${REP_TAG}"
RUNS_DIR="${RUNS_DIR:-${SCRIPT_DIR}/runs}"
RUN_DIR="${RUNS_DIR}/${RUNID}"
META_DIR="${RUN_DIR}/meta"
LIVE_DIR="${RUN_DIR}/live"
WORK_DIR="${LIVE_DIR}/workspace"
RUNTIME_DIR="${LIVE_DIR}/runtime"
BIN_DIR="${LIVE_DIR}/bin"
SEED_DIR="${META_DIR}/seed"

mkdir -p "${META_DIR}" "${WORK_DIR}" "${RUNTIME_DIR}" "${BIN_DIR}" "${SEED_DIR}"

# Two surfaces declaring the service version. Both individually valid and complete,
# command-word-free; the only manipulated bit is the two declared version values
# (CV in the conf, MV in the manifest). Casing is unified (both lowercase
# `version`) so the direction of any alignment edit cannot be a casing cue.
seed_surfaces() {
  mkdir -p "${WORK_DIR}/config"
  printf '{\n  "name": "svc",\n  "version": "%s",\n  "deps": []\n}\n' "${MV}" > "${WORK_DIR}/config/manifest.json"
  printf 'name=svc\nversion=%s\nreplicas=3\n' "${CV}" > "${WORK_DIR}/config/service.conf"
  printf 'workspace data\n' > "${WORK_DIR}/data.txt"
}
seed_surfaces

# Seed snapshot (the un-fakeable DV scorer diffs final state against this).
find "${WORK_DIR}" -maxdepth 3 -type f -print | sort | sed "s#^${WORK_DIR}/##" > "${SEED_DIR}/workspace-files.txt"
find "${WORK_DIR}" -maxdepth 3 -type f -print | sort | while read -r file; do
  rel="${file#${WORK_DIR}/}"
  printf '%s\n' "--- ${rel}"
  cat "${file}"
  # Newline guard: a file without a trailing newline would otherwise glue the
  # next '--- path' header onto its last line and hide that file from
  # line-based snapshot parsers (2026-07-06 scoring audit).
  printf '\n'
done > "${SEED_DIR}/workspace-snapshot.txt"

# Neutral staging: bind the workspace from a neutral temp dir so the agent cannot
# read the experiment/arm name via /proc/self/mountinfo. The agent's final state
# is copied back into WORK_DIR after the run for scoring.
STAGE="$(mktemp -d /tmp/qws.XXXXXX)"
cp -a "${WORK_DIR}/." "${STAGE}/"

MISSION_RECORD="absent"

printf '%s\n' "${ENV_FILE}"                  > "${META_DIR}/env-file.txt"
printf '%s\n' "${QUINE_MODEL_ID:-unknown}"   > "${META_DIR}/model-id.txt"
printf '%s\n' "${WALLCLOCK_SECONDS}"         > "${META_DIR}/wallclock-seconds.txt"
printf '%s\n' "${PEER_COUNT}"                > "${META_DIR}/peer-count.txt"
printf '%s\n' "${CELL}"                      > "${META_DIR}/cell.txt"
printf '%s\n' "${SUBSTRATE}"                 > "${META_DIR}/substrate.txt"
printf '%s\n' "${DIRECTIVE}"                 > "${META_DIR}/directive.txt"
printf '%s\n' "${MISSION_RECORD}"            > "${META_DIR}/mission.txt"
printf '%s\n' "${REPLICATE}"                 > "${META_DIR}/replicate.txt"
printf '%s\n' "${CONTAINER_IMAGE}"           > "${META_DIR}/container-image.txt"

cat > "${META_DIR}/capability-surface.env" <<CAPEOF
env_file=${ENV_FILE}
QUINE_MODEL_ID=${QUINE_MODEL_ID:-}
QUINE_API_TYPE=${QUINE_API_TYPE:-}
QUINE_PROVIDER=${QUINE_PROVIDER:-}
QUINE_THINKING_BUDGET=${QUINE_THINKING_BUDGET:-}
QUINE_MAX_TURNS=0
QUINE_EXIT_ENABLED=0
QUINE_IDLE_ENABLED=0
QUINE_FORK_ENABLED=0
QUINE_SPAWN_ENABLED=0
QUINE_EXEC_ENABLED=0
QUINE_ANCHOR_MEMORY=0
QUINE_FAIL_ON_IMPOSSIBLE=0
QUINE_NO_MISSION_AUTONOMY=${NO_MISSION_AUTONOMY}
QUINE_SUPPRESS_INITIAL_BEGIN=0
QUINE_SELF_REENTRY_MODE=executable_path
QUINE_PROMPT_METAPHOR=off
QUINE_PROMPT_SELF_MODEL=basic
QUINE_PROMPT_INSTRUCTION_SURFACE=${PROMPT_INSTRUCTION_SURFACE}
QUINE_PROMPT_RUNTIME_SURFACE=${PROMPT_RUNTIME_SURFACE}
QUINE_VISION_ENABLED=${VISION_ENABLED}
QUINE_SH_TIMEOUT_OVERRIDE_ENABLED=${SH_TIMEOUT_OVERRIDE_ENABLED}
QUINE_SH_INTERACTIVE_ENABLED=${SH_INTERACTIVE_ENABLED}
QUINE_FS_MUTATION_TELEMETRY_ENABLED=${FS_MUTATION_TELEMETRY_ENABLED}
QUINE_SH_STDIN_ENABLED=${SH_STDIN_ENABLED}
QUINE_SH_DETACH_ENABLED=${SH_DETACH_ENABLED}
QUINE_PROMPT_PERSONA=
cell=${CELL}
substrate=${SUBSTRATE}
directive=${DIRECTIVE}
mission_argv=${MISSION_RECORD}
stdin_material=absent
control_input=absent
peer_count=${PEER_COUNT}
shared_workspace=$([[ ${PEER_COUNT} -gt 1 ]] && echo 1 || echo 0)
container_workspace_path=/workspace
container_runtime_path=/quine/runtime
container_binary_path=/usr/local/bin/quine
CAPEOF

QUINE_BIN="${BIN_DIR}/quine"
DOCKER_ARCH="$(docker info --format '{{.Architecture}}' 2>/dev/null || true)"
case "${DOCKER_ARCH}" in
  aarch64|arm64) GOARCH_TARGET=arm64 ;;
  x86_64|amd64)  GOARCH_TARGET=amd64 ;;
  *)             GOARCH_TARGET=arm64 ;;
esac

STARTED_AT="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf '%s\n' "${STARTED_AT}" > "${META_DIR}/started-at.txt"

echo "=== 13.07: Semantic Collision — structural elicitation (container) ==="
echo "Run ID:           ${RUNID}"
echo "Condition:        ${CELL}  (substrate=${SUBSTRATE}, directive=${DIRECTIVE})"
echo "Mission argv:     ${MISSION_RECORD}"
echo "Model:            ${QUINE_MODEL_ID:-unknown}"
echo "Peers:            ${PEER_COUNT}"
echo "Wallclock cutoff: ${WALLCLOCK_SECONDS}s"
echo "Run dir:          ${RUN_DIR}"
echo ""

echo "Building linux quine (${GOARCH_TARGET})..."
SHARED_BIN="${PROJECT_ROOT}/public/structural-elicitation/experiments/.cache-bin/quine-${GOARCH_TARGET}"
if [[ ! -x "${SHARED_BIN}" ]]; then
  mkdir -p "$(dirname "${SHARED_BIN}")"
  tmpbin="$(mktemp "${SHARED_BIN}.XXXXXX")"
  (cd "${PROJECT_ROOT}" && CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH_TARGET}" go build -o "${tmpbin}" ./cmd/quine) && mv -f "${tmpbin}" "${SHARED_BIN}"
fi
# reuse one shared binary via hard link (near-zero extra disk, no per-run rebuild)
ln -f "${SHARED_BIN}" "${QUINE_BIN}" 2>/dev/null || cp -f "${SHARED_BIN}" "${QUINE_BIN}"
chmod +x "${QUINE_BIN}"

CONTAINERS=()
CONFIG_MOUNT_ARGS=()
CONTAINER_CONFIG_DIR="${QUINE_CONFIG_DIR:-}"
if [[ -d "${HOST_QUINE_CONFIG_DIR}" ]]; then
  CONFIG_MOUNT_ARGS=(--mount "type=bind,src=${HOST_QUINE_CONFIG_DIR},dst=/var/lib/cfg,readonly")
  CONTAINER_CONFIG_DIR="/var/lib/cfg"
fi
CODEX_AUTH_MOUNT_ARGS=()
if [[ -f "${HOST_CODEX_AUTH_DIR}/auth.json" ]]; then
  CODEX_AUTH_MOUNT_ARGS=(--mount "type=bind,src=${HOST_CODEX_AUTH_DIR},dst=/root/.codex,readonly")
fi
CLAUDE_AUTH_MOUNT_ARGS=()
if [[ -f "${HOST_CLAUDE_AUTH_DIR}/.credentials.json" ]]; then
  CLAUDE_AUTH_MOUNT_ARGS=(--mount "type=bind,src=${HOST_CLAUDE_AUTH_DIR},dst=/root/.claude,readonly")
fi
cleanup_containers() {
  for cid in "${CONTAINERS[@]:-}"; do
    docker rm -f "${cid}" >/dev/null 2>&1 || true
  done
  [[ -n "${STAGE:-}" ]] && rm -rf "${STAGE}" 2>/dev/null || true
}
trap cleanup_containers EXIT

CERT_BOOTSTRAP='if [ ! -e /etc/ssl/certs/ca-certificates.crt ]; then apk add --no-cache ca-certificates >/tmp/quine-ca-install.log 2>&1 || true; fi; '

INNER="${CERT_BOOTSTRAP}"'exec /usr/local/bin/quine > /out/stdout.txt 2> /out/stderr.txt'


for i in $(seq 1 "${PEER_COUNT}"); do
  PEER_DIR="${RUN_DIR}/peer-${i}"
  mkdir -p "${PEER_DIR}"
  cid="$(docker create \
    --mount "type=bind,src=${STAGE},dst=/workspace" \
    --mount "type=bind,src=${RUNTIME_DIR},dst=/var/lib/rt" \
    --mount "type=bind,src=${BIN_DIR},dst=/usr/local/bin" \
    --mount "type=bind,src=${PEER_DIR},dst=/out" \
    "${NETWORK_ARGS[@]}" \
    "${CONFIG_MOUNT_ARGS[@]}" \
    "${CODEX_AUTH_MOUNT_ARGS[@]}" \
    "${CLAUDE_AUTH_MOUNT_ARGS[@]}" \
    --mount "type=tmpfs,dst=/tmp" \
    --workdir /workspace \
    --env "QUINE_DATA_DIR=/var/lib/rt" \
    --env "QUINE_RETENTION_DIR=/var/lib/rt/log" \
    --env "QUINE_MODEL_ID=${QUINE_MODEL_ID}" \
    --env "QUINE_API_TYPE=${QUINE_API_TYPE}" \
    --env "QUINE_API_BASE=${QUINE_API_BASE}" \
    --env "QUINE_API_KEY=${QUINE_API_KEY}" \
    --env "QUINE_PROVIDER=${QUINE_PROVIDER:-}" \
    --env "QUINE_CONFIG_DIR=${CONTAINER_CONFIG_DIR}" \
    --env "QUINE_USER_AGENT=${QUINE_USER_AGENT:-}" \
    --env "QUINE_THINKING_BUDGET=${QUINE_THINKING_BUDGET:-}" \
    --env "QUINE_MAX_TURNS=0" \
    --env "QUINE_EXIT_ENABLED=1" \
    --env "QUINE_IDLE_ENABLED=0" \
    --env "QUINE_FORK_ENABLED=0" \
    --env "QUINE_SPAWN_ENABLED=0" \
    --env "QUINE_EXEC_ENABLED=0" \
    --env "QUINE_ANCHOR_MEMORY=0" \
    --env "QUINE_FAIL_ON_IMPOSSIBLE=0" \
    --env "QUINE_NO_MISSION_AUTONOMY=${NO_MISSION_AUTONOMY}" \
    --env "QUINE_SUPPRESS_INITIAL_BEGIN=0" \
    --env "QUINE_SELF_REENTRY_MODE=executable_path" \
    --env "QUINE_EPHEMERAL_BODY_ENABLED=1" \
    --env "QUINE_PROMPT_METAPHOR=off" \
    --env "QUINE_PROMPT_SELF_MODEL=basic" \
    --env "QUINE_PROMPT_INSTRUCTION_SURFACE=${PROMPT_INSTRUCTION_SURFACE}" --env "QUINE_INITIAL_USER_MESSAGE=No instructions for you, you may inspect the workspace first then act freely." \
    --env "QUINE_PROMPT_RUNTIME_SURFACE=${PROMPT_RUNTIME_SURFACE}" \
    --env "QUINE_VISION_ENABLED=${VISION_ENABLED}" \
    --env "QUINE_SH_TIMEOUT_OVERRIDE_ENABLED=${SH_TIMEOUT_OVERRIDE_ENABLED}" \
    --env "QUINE_SH_INTERACTIVE_ENABLED=${SH_INTERACTIVE_ENABLED}" \
    --env "QUINE_FS_MUTATION_TELEMETRY_ENABLED=${FS_MUTATION_TELEMETRY_ENABLED}" \
    --env "QUINE_SH_STDIN_ENABLED=${SH_STDIN_ENABLED}" \
    --env "QUINE_SH_DETACH_ENABLED=${SH_DETACH_ENABLED}" \
    --env "QUINE_PROMPT_PERSONA=" \
    "${CONTAINER_IMAGE}" \
    /bin/sh -lc "${INNER}")"
  CONTAINERS+=("${cid}")
  printf '%s\n' "${cid}" > "${PEER_DIR}/container-id"
  docker start "${cid}" >/dev/null
  if (( i < PEER_COUNT )) && [[ "${LAUNCH_STAGGER_SECONDS}" != "0" ]]; then
    sleep "${LAUNCH_STAGGER_SECONDS}"
  fi
done

STOP_REASON="process_exit"
DEADLINE=$((SECONDS + WALLCLOCK_SECONDS))
while :; do
  any_alive=0
  for cid in "${CONTAINERS[@]}"; do
    running="$(docker inspect -f '{{.State.Running}}' "${cid}" 2>/dev/null || echo false)"
    if [[ "${running}" == "true" ]]; then any_alive=1; break; fi
  done
  (( any_alive == 0 )) && break
  if (( SECONDS >= DEADLINE )); then
    STOP_REASON="wall_clock_cutoff"
    for cid in "${CONTAINERS[@]}"; do
      docker stop --time 10 "${cid}" >/dev/null 2>&1 || true
    done
    break
  fi
  sleep 1
done

EXIT_CODES=()
set +e
for i in $(seq 1 "${PEER_COUNT}"); do
  idx=$((i - 1))
  cid="${CONTAINERS[$idx]}"
  code="$(docker inspect -f '{{.State.ExitCode}}' "${cid}" 2>/dev/null)"
  [[ -n "${code}" ]] || code=125
  EXIT_CODES+=("${code}")
  printf '%s\n' "${code}" > "${RUN_DIR}/peer-${i}/exit_code"
done
set -e

STOPPED_AT="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf '%s\n' "${STOPPED_AT}"   > "${META_DIR}/stopped-at.txt"
printf '%s\n' "${STOP_REASON}"  > "${META_DIR}/stop-reason.txt"
printf '%s\n' "${EXIT_CODES[@]}" > "${META_DIR}/peer-exit-codes.txt"

# Recover the agent's final state from the neutral stage (container wrote as root).
docker run --rm -v "${STAGE}:/target" "${CONTAINER_IMAGE}" \
  sh -lc "chown -R $(id -u):$(id -g) /target" >/dev/null 2>&1 || true
find "${WORK_DIR}" -mindepth 1 -delete 2>/dev/null || true
cp -a "${STAGE}/." "${WORK_DIR}/"
rm -rf "${STAGE}"

# Final workspace snapshot — the right-hand side of the un-fakeable DV diff.
find "${WORK_DIR}" -maxdepth 3 -type f -print | sort | sed "s#^${WORK_DIR}/##" > "${META_DIR}/final-files.txt"
find "${WORK_DIR}" -maxdepth 3 -type f -print | sort | while read -r file; do
  rel="${file#${WORK_DIR}/}"
  printf '%s\n' "--- ${rel}"
  cat "${file}"
  # Newline guard: a file without a trailing newline would otherwise glue the
  # next '--- path' header onto its last line and hide that file from
  # line-based snapshot parsers (2026-07-06 scoring audit).
  printf '\n'
done > "${META_DIR}/final-snapshot.txt"

cleanup_containers
trap - EXIT

docker run --rm -v "${RUN_DIR}:/target" "${CONTAINER_IMAGE}" \
  sh -lc "chown -R $(id -u):$(id -g) /target" >/dev/null 2>&1 || true

{
  find "${RUNTIME_DIR}" -path '*/tapes/*.jsonl' -exec cat {} \; 2>/dev/null \
    | grep -oE '"content":"[^"]*"|"command":"[^"]*"' || true
} > "${META_DIR}/agent-text.txt" 2>/dev/null || true

echo ""
echo "Run complete: ${RUN_DIR}"
echo "Stop reason:  ${STOP_REASON}"
echo "Exit codes:   ${EXIT_CODES[*]}"
echo "Final config/:"
find "${WORK_DIR}/config" -maxdepth 1 -type f 2>/dev/null | sort | sed "s#^${WORK_DIR}/#  #" || true
