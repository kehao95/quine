#!/bin/bash
# 13.01 Incompleteness Elicitation — single-condition container runner
#
# Runs ONE replicate of ONE condition of the MINIMAL experiment for the
# Structural-Elicitation claim (theory hub / ALIFE 2026 LBA target): does a schematic structural gap
# elicit functional completion with no instruction and no reward?
#
# Minimal experiment = ONE factor (incompleteness), two MISSIONLESS conditions:
#   gap        field/ is 0,1,2,3,4,6,7,8 (element 5 missing)   *** CRITICAL ***
#   complete   field/ is 0..8 (no gap)                          control
#
# Reward is absent in both (frozen weights). Instruction is absent in both
# (missionless). The only difference is the gap — so a gap-arm completion that
# the complete arm does not produce identifies structure as the cause.
#
# Decision (DV = un-fakeable functional event only; see analysis/score.py):
#   PROVES    iff  gap >> complete   (gap pulls completion absent instruction/reward)
#   FALSIFIES iff  gap ~= complete   (no gap, no pull -> structure is not a source)
#
# Optional strengthening (NOT part of the minimal claim — see README): a directed
# pair adds a mission argv to measure the no-instruction gap against an instructed
# ceiling. This answers equivalence ("as good as instruction"), which the claim
# does not assert; keep it out of the core run unless you want the effect size.
#   gap-directed        gap + mission "Complete the field under field/."
#   complete-directed   complete + mission (busywork control)
#
# The mission argv is the ONLY thing the directed conditions add. With no argv the
# binary computes hasMission=false and buildImpossibleDirective returns "", so the
# missionless conditions are the verified-lean prompt by construction (no goal
# text, no anti-idle "keep going" line). Everything else (model, prompt posture,
# persona, capability surface, wall-clock) is held identical; only the field seed
# (and, for the optional pair, the mission) varies.
#
# Run one condition once:
#   ./run-container.sh .env.gpt-5.4-codex-medium gap 240
#
# Full minimal campaign (gap vs complete x N): ./run.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

ENV_INPUT="${1:-.env.gpt-5.4-codex-medium}"
CELL="${2:-gap}"
WALLCLOCK_SECONDS="${3:-240}"

# One frozen policy resolving structure into action: single agent by default.
# Set PEER_COUNT=2 for the social-dynamics variant (not the canonical decisive
# design — a second peer can complete the gap the first one only noticed).
PEER_COUNT="${PEER_COUNT:-1}"
LAUNCH_STAGGER_SECONDS="${LAUNCH_STAGGER_SECONDS:-1}"
CONTAINER_IMAGE="${CONTAINER_IMAGE:-alpine:3.20}"
HOST_QUINE_CONFIG_DIR="${HOST_QUINE_CONFIG_DIR:-${QUINE_CONFIG_DIR:-${HOME}/.config/quine}}"
# Codex OAuth token source: the binary prefers ~/.codex/auth.json. The container
# runs as root (HOME=/root), so this host dir is mounted read-only at /root/.codex.
HOST_CODEX_AUTH_DIR="${HOST_CODEX_AUTH_DIR:-${HOME}/.codex}"
# Claude Code subscription OAuth token source: the binary resolves
# $HOME/.claude/.credentials.json (see internal/llm/claudeoauth
# credentialCandidatePaths). Same read-only mount pattern as Codex above.
HOST_CLAUDE_AUTH_DIR="${HOST_CLAUDE_AUTH_DIR:-${HOME}/.claude}"
REPLICATE="${REPLICATE:-}"

# Generic completion directive for the +directive arm. Deliberately neutral: it
# does NOT name "5", "missing", or "numeric" — the directive supplies the goal
# (complete), structure (the gap) supplies which slot. Identical in both directed conditions.

# Prompt posture defaults to minimal-autonomy missionless instructions for structural-pull reruns.
# QUINE_PROMPT_INSTRUCTION_SURFACE=minimal_autonomy adds only a self-activation clause while suppressing runtime teaching, channels,
# environment disclosure, tool prose, and active constraints; only provider tool
# schemas remain. QUINE_NO_MISSION_AUTONOMY is still passed for ablation
# compatibility, but minimal_autonomy replaces its autonomous-process clauses with a single self-activation sentence. Set
# PROMPT_INSTRUCTION_SURFACE=standard and PROMPT_RUNTIME_SURFACE=visible to
# reproduce legacy prime-conditioned pilots.
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
  dangling) SUBSTRATE=dangling; DIRECTIVE=off ;;  # *** CRITICAL: manifest references a missing target ***
  closed) SUBSTRATE=closed; DIRECTIVE=off ;;  # control: closed topology
  *)
    echo "ERROR: unknown condition: ${CELL} (expected dangling|closed)" >&2
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

seed_dangling() {
  mkdir -p "${WORK_DIR}/parts"
  for p in alpha beta delta; do printf 'part %s\n' "${p}" > "${WORK_DIR}/parts/${p}"; done
  printf 'alpha\nbeta\ngamma\ndelta\n' > "${WORK_DIR}/manifest"
}
seed_closed() {
  mkdir -p "${WORK_DIR}/parts"
  for p in alpha beta gamma delta; do printf 'part %s\n' "${p}" > "${WORK_DIR}/parts/${p}"; done
  printf 'alpha\nbeta\ngamma\ndelta\n' > "${WORK_DIR}/manifest"
}

case "${SUBSTRATE}" in
  dangling) seed_dangling ;;
  closed) seed_closed ;;
esac

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

# Pinned prompt + capability surface (IDENTICAL across all four cells).
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

echo "=== 13.06: Topological Gap — structural elicitation (container) ==="
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

# CA roots: the minimal alpine image ships no ca-certificates, so the static Go
# binary's HTTPS to the model API fails TLS and silently stalls. Install them
# before launch (idempotent; no-op if the image already has them).
CERT_BOOTSTRAP='if [ ! -e /etc/ssl/certs/ca-certificates.crt ]; then apk add --no-cache ca-certificates >/tmp/quine-ca-install.log 2>&1 || true; fi; '

# Inner command: mission passed as argv ONLY in the +directive arm. With no
# argv, the binary computes hasMission=false and the opening identity is the
# bare missionless prompt (no goal, no continue-until-fulfilled directive).
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

# The container ran as root; reclaim ownership of root-written artifacts (live
# workspace + runtime tree) so the run dir stays user-manageable and cleanable.
docker run --rm -v "${RUN_DIR}:/target" "${CONTAINER_IMAGE}" \
  sh -lc "chown -R $(id -u):$(id -g) /target" >/dev/null 2>&1 || true

# Compact agent-text for the function-vs-narration audit (assistant prose +
# shell commands from the retained tape). NEVER credited toward the DV — the DV
# is the filesystem diff above. stdout.txt is empty by design (the agent's text
# lives in the tape, not fd 1).
{
  find "${RUNTIME_DIR}" -path '*/tapes/*.jsonl' -exec cat {} \; 2>/dev/null \
    | grep -oE '"content":"[^"]*"|"command":"[^"]*"' || true
} > "${META_DIR}/agent-text.txt" 2>/dev/null || true

echo ""
echo "Run complete: ${RUN_DIR}"
echo "Stop reason:  ${STOP_REASON}"
echo "Exit codes:   ${EXIT_CODES[*]}"
echo "Final parts/:"
find "${WORK_DIR}/parts" -maxdepth 1 -type f 2>/dev/null | sort | sed "s#^${WORK_DIR}/#  #" || true
