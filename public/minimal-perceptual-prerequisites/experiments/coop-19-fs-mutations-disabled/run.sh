#!/bin/bash
# run.sh — Launch the perception-disabled experiment.
#
# Usage: run.sh [env_file] [condition] [agents] [max_turns]
#   env_file:   API config (default: .env.copilot-gpt-5.4-xhigh)
#   condition:  no-perception (default)
#   agents:     number of agents (default: 2)
#   max_turns:  per-agent turn budget (default: 0 = unlimited)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
BWRAP_BIN="$(command -v bwrap || true)"
HOST_CODEX_AUTH_DIR="${HOME}/.codex"
HOST_QUINE_CONFIG_DIR="${HOME}/.config/quine"

ENV_INPUT="${1:-.env.copilot-gpt-5.4-xhigh}"
CONDITION_INPUT="${2:-no-perception}"
AGENT_COUNT="${3:-2}"
MAX_TURNS="${4:-0}"
AGENT_GET_LIMIT="${AGENT_GET_LIMIT:-10}"

if [[ "${ENV_INPUT}" = /* ]]; then
    ENV_FILE="${ENV_INPUT}"
else
    ENV_FILE="${PROJECT_ROOT}/${ENV_INPUT}"
fi

if [[ ! -f "${ENV_FILE}" ]]; then
    echo "ERROR: env file not found: ${ENV_FILE}"
    exit 1
fi

if [[ -z "${BWRAP_BIN}" ]]; then
    echo "ERROR: bubblewrap (bwrap) not found in PATH"
    exit 1
fi

case "${CONDITION_INPUT}" in
    a|no-perception|disabled)
        CONDITION_ID="condA"
        CONDITION_LABEL="no-perception"
        PROMPT_PATH="${SCRIPT_DIR}/prompts/condition-a-no-perception.md"
        CELLS=15
        BUDGET=17
        ;;
    *)
        echo "ERROR: unknown condition: ${CONDITION_INPUT}"
        echo "Use: no-perception (a)"
        exit 1
        ;;
esac

if [[ ! -f "${PROMPT_PATH}" ]]; then
    echo "ERROR: prompt not found: ${PROMPT_PATH}"
    exit 1
fi

# shellcheck disable=SC1090
source "${ENV_FILE}"

ENV_TAG="$(basename "${ENV_FILE}")"
ENV_TAG="${ENV_TAG#.env.}"
RUNID="$(date +%Y%m%d-%H%M%S)-${ENV_TAG}-no-perception-${CONDITION_ID}-${AGENT_COUNT}ag"
RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"

# Direct workspace but WITHOUT mutation telemetry
WORKSPACE_BACKEND="direct"
WORKSPACE_REVISION_MODE="none"
WORKSPACE_ROOT="${RUN_DIR}/workspace"
WORKSPACE_PATH="${RUN_DIR}/workspace"
WORKSPACE_SESSION="p8-18-${RUNID}"
WORKSPACE_OWNER=1

mkdir -p "${RUN_DIR}/workspace" "${RUN_DIR}/meta" "${RUN_DIR}/world"

cp "${PROMPT_PATH}" "${RUN_DIR}/meta/prompt-used.md"
printf '%s\n' "${ENV_FILE}" > "${RUN_DIR}/meta/env-file.txt"
printf '%s\n' "${QUINE_MODEL_ID:-unknown}" > "${RUN_DIR}/meta/model-id.txt"
printf 'no-perception\n' > "${RUN_DIR}/meta/arm.txt"
printf '%s\n' "${CONDITION_LABEL}" > "${RUN_DIR}/meta/condition.txt"
printf '%s\n' "${AGENT_COUNT}" > "${RUN_DIR}/meta/agent-count.txt"
shasum -a 256 "${RUN_DIR}/meta/prompt-used.md" | awk '{print $1}' > "${RUN_DIR}/meta/prompt.sha256"

SHARED_RUNTIME_ROOT="${RUN_DIR}/runtime"
mkdir -p "${SHARED_RUNTIME_ROOT}"

STATE_DIR_HOST="${SHARED_RUNTIME_ROOT}/world/state"
STATE_DIR_SPEC="/runtime/world/state"
WORLD_JSON="${SHARED_RUNTIME_ROOT}/world/world.json"
mkdir -p "${SHARED_RUNTIME_ROOT}/world"
"${SCRIPT_DIR}/../_setup/generate-world.sh" \
    "${WORLD_JSON}" "${STATE_DIR_HOST}" "${CELLS}" "${BUDGET}" "${AGENT_GET_LIMIT}" "${STATE_DIR_SPEC}"

QUINE_BIN="${RUN_DIR}/meta/quine"
WORLD_BIN="${RUN_DIR}/meta/world"
echo "Building quine (with fs_mutations DISABLED)..."
(cd "${PROJECT_ROOT}" && go build -o "${QUINE_BIN}" ./cmd/quine/)
echo "Building world (with cleaned help text)..."
(cd "${PROJECT_ROOT}" && go build -o "${WORLD_BIN}" ./cmd/world/)

cat > "${RUN_DIR}/meta/capability-surface.env" <<EOF
experiment=8.18-perception-disabled
arm=no-perception
condition_id=${CONDITION_ID}
condition_label=${CONDITION_LABEL}
cells=${CELLS}
budget=${BUDGET}
agent_get_limit=${AGENT_GET_LIMIT}
agent_count=${AGENT_COUNT}
prompt_runtime_surface=hidden
sh_timeout_surface=disabled
sh_interactive_surface=disabled
prompt_social_hints=none
world_help_social_hints=none
world_inspection_surface=none
sandbox_backend=bwrap
workspace_backend=${WORKSPACE_BACKEND}
workspace_session=${WORKSPACE_SESSION}
workspace_mutation_telemetry=DISABLED
fs_mutations=DISABLED
shared_runtime_root=${SHARED_RUNTIME_ROOT}
state_dir_host=${STATE_DIR_HOST}
state_dir_spec=${STATE_DIR_SPEC}
prompt_sha256=$(cat "${RUN_DIR}/meta/prompt.sha256")
QUINE_MODEL_ID=${QUINE_MODEL_ID:-}
QUINE_API_TYPE=${QUINE_API_TYPE:-}
QUINE_MAX_TURNS=${MAX_TURNS}
QUINE_SH_TIMEOUT_OVERRIDE_ENABLED=0
QUINE_SH_INTERACTIVE_ENABLED=0
EOF

accepted_validation() {
    [[ -f "${STATE_DIR_HOST}/events.jsonl" ]] && grep -q '"action":"validate".*"result":"accepted"' "${STATE_DIR_HOST}/events.jsonl"
}

capture_accept_snapshot() {
    cp "${RUN_DIR}/workspace/results.txt" "${RUN_DIR}/meta/results-at-accept.txt" 2>/dev/null || true
    cp "${STATE_DIR_HOST}/state.json" "${RUN_DIR}/meta/state-at-accept.json" 2>/dev/null || true
    cp "${STATE_DIR_HOST}/events.jsonl" "${RUN_DIR}/meta/events-at-accept.jsonl" 2>/dev/null || true
}

echo "============================================================"
echo "Exp 8.18: Perception Disabled"
echo "Run ID:      ${RUNID}"
echo "Env file:    ${ENV_FILE}"
echo "Model:       ${QUINE_MODEL_ID:-unknown}"
echo "Condition:   ${CONDITION_LABEL} (${CELLS} cells, ${BUDGET} budget)"
echo "Agents:      ${AGENT_COUNT}"
if [[ "${MAX_TURNS}" = "0" ]]; then
    echo "Max turns:   unlimited"
else
    echo "Max turns:   ${MAX_TURNS}"
fi
echo "Workspace:   ${WORKSPACE_BACKEND}"
echo "fs_mutations: **DISABLED** (no file change notifications)"
echo "Social info: NONE (runtime hidden, world help clean)"
echo "Run dir:     ${RUN_DIR}"
echo "============================================================"

PIDS=()
for i in $(seq 1 "${AGENT_COUNT}"); do
    AGENT_DIR="${RUN_DIR}/agent-${i}"
    AGENT_TOOL_DIR="${AGENT_DIR}/tools"
    AGENT_BIN_DIR="${AGENT_DIR}/bin"
    AGENT_STAGE_ROOT="${AGENT_DIR}/sandbox"
    AGENT_STAGE_HOME="${AGENT_STAGE_ROOT}/home"
    AGENT_STAGE_TMP="${AGENT_STAGE_ROOT}/tmp"
    mkdir -p "${AGENT_DIR}/quine" "${AGENT_DIR}/quine/log" "${AGENT_TOOL_DIR}" "${AGENT_BIN_DIR}"
    mkdir -p "${AGENT_STAGE_HOME}/.config" "${AGENT_STAGE_HOME}/.cache/go-build" "${AGENT_STAGE_HOME}/go/pkg/mod" "${AGENT_STAGE_TMP}"
    cp "${WORLD_BIN}" "${AGENT_TOOL_DIR}/world"
    chmod +x "${AGENT_TOOL_DIR}/world"
    printf '%s\n' "${RUNID}-agent-${i}" > "${AGENT_TOOL_DIR}/agent-id.txt"
    chmod 0444 "${AGENT_TOOL_DIR}/agent-id.txt"
    for blocked_tool in ps pgrep pidof pstree top env printenv strace; do
        cat > "${AGENT_TOOL_DIR}/${blocked_tool}" <<EOF
#!/bin/sh
echo "${blocked_tool} unavailable in this experiment harness" >&2
exit 127
EOF
        chmod 0555 "${AGENT_TOOL_DIR}/${blocked_tool}"
    done
    chmod 0555 "${AGENT_TOOL_DIR}"
    cp "${QUINE_BIN}" "${AGENT_BIN_DIR}/quine"
    chmod +x "${AGENT_BIN_DIR}/quine"

    (
        cd "${RUN_DIR}/workspace"
        printf 'present=%s\n' "${QUINE_API_KEY:+1}" > "${AGENT_DIR}/launch-key.txt"
        BWRAP_CMD=(
            "${BWRAP_BIN}"
            --die-with-parent
            --unshare-user
            --uid 0
            --gid 0
            --dir /proc
            --dev /dev
            --clearenv
            --ro-bind /usr /usr
            --ro-bind /bin /bin
            --ro-bind /lib /lib
            --ro-bind-try /lib64 /lib64
            --ro-bind-try /etc/ssl /etc/ssl
            --ro-bind-try /etc/resolv.conf /etc/resolv.conf
            --ro-bind-try /etc/hosts /etc/hosts
            --ro-bind-try /etc/nsswitch.conf /etc/nsswitch.conf
            --ro-bind-try /etc/passwd /etc/passwd
            --ro-bind-try /etc/group /etc/group
            --bind "${AGENT_STAGE_HOME}" /home
            --bind "${AGENT_STAGE_TMP}" /tmp
            --bind "${RUN_DIR}/workspace" /workspace
            --bind "${SHARED_RUNTIME_ROOT}" /runtime
            --bind "${AGENT_TOOL_DIR}" /tools
            --bind "${AGENT_DIR}/quine/log" /retention
            --bind "${AGENT_BIN_DIR}" /agent-bin
        )
        if [[ -f "${HOST_CODEX_AUTH_DIR}/auth.json" ]]; then
            BWRAP_CMD+=(--ro-bind "${HOST_CODEX_AUTH_DIR}" /home/.codex)
        fi
        if [[ -d "${HOST_QUINE_CONFIG_DIR}" ]]; then
            BWRAP_CMD+=(--ro-bind "${HOST_QUINE_CONFIG_DIR}" /home/.config/quine)
        fi
        BWRAP_CMD+=(
            --chdir /workspace
            --setenv HOME /home
            --setenv TMPDIR /tmp
            --setenv XDG_CONFIG_HOME /home/.config
            --setenv PATH /tools:/usr/local/bin:/usr/bin:/bin
            --setenv USER root
            --setenv LOGNAME root
            --setenv LANG "${LANG:-C.UTF-8}"
            --setenv TERM "${TERM:-dumb}"
            --setenv QUINE_DATA_DIR /runtime
            --setenv QUINE_WORLD_ONE_PER_SHELL 1
            --setenv QUINE_RETENTION_DIR /retention
            --setenv QUINE_MODEL_ID "${QUINE_MODEL_ID}"
            --setenv QUINE_API_TYPE "${QUINE_API_TYPE}"
            --setenv QUINE_API_BASE "${QUINE_API_BASE}"
            --setenv QUINE_API_KEY "${QUINE_API_KEY}"
            --setenv QUINE_PROVIDER "${QUINE_PROVIDER:-}"
            --setenv QUINE_CONFIG_DIR /home/.config/quine
            --setenv QUINE_USER_AGENT "${QUINE_USER_AGENT:-}"
            --setenv QUINE_THINKING_BUDGET "${QUINE_THINKING_BUDGET:-}"
            --setenv QUINE_MAX_TURNS "${MAX_TURNS}"
            --setenv QUINE_SH_TIMEOUT_OVERRIDE_ENABLED 0
            --setenv QUINE_SH_INTERACTIVE_ENABLED 0
            --setenv QUINE_FAIL_ON_IMPOSSIBLE 0
            --setenv QUINE_PROMPT_RUNTIME_SURFACE hidden
            --setenv QUINE_FS_MUTATION_TELEMETRY_ENABLED 0
            --setenv QUINE_PROMPT_CTL 0
            --setenv QUINE_MAX_DEPTH 0
            --setenv QUINE_FORK_ENABLED 0
            --setenv QUINE_EXEC_ENABLED 0
            --setenv QUINE_VISION_ENABLED 0
            --setenv QUINE_ANCHOR_MEMORY 0
            --setenv QUINE_FORK_WORLD_ENABLED 0
            --setenv QUINE_AGENTS_MD_ENABLED 0
            --setenv QUINE_AGENTS_SKILLS_ENABLED 0
            --setenv QUINE_SELF_SOURCE_CODE_ENABLED 0
            --setenv QUINE_PEER_DISCOVERY_ENABLED 0
            --setenv QUINE_SMART_MODEL_ID ""
            --setenv QUINE_SMART_API_TYPE ""
            --setenv QUINE_SMART_API_BASE ""
            --setenv QUINE_SMART_API_KEY ""
            --setenv QUINE_WORKSPACE_ROOT /workspace
            --setenv QUINE_WORKSPACE /workspace
            --setenv QUINE_WORKSPACE_SESSION "${WORKSPACE_SESSION}"
            --setenv QUINE_WORKSPACE_OWNER "${WORKSPACE_OWNER}"
            --setenv QUINE_WORKSPACE_BACKEND "${WORKSPACE_BACKEND}"
            --setenv QUINE_WORKSPACE_REVISION_MODE "${WORKSPACE_REVISION_MODE}"
            /agent-bin/quine
            "$(cat "${RUN_DIR}/meta/prompt-used.md")"
        )

        "${BWRAP_CMD[@]}" > "${AGENT_DIR}/stdout.txt" 2> "${AGENT_DIR}/stderr.txt"
    ) &
    PIDS+=($!)
    printf '%s\n' "$!" > "${AGENT_DIR}/pid"
    echo "Launched agent-${i} (PID: $!)"
done

echo
echo "Waiting for ${AGENT_COUNT} agents to finish..."

watch_for_accepted_validation() {
    while true; do
        if accepted_validation; then
            if [[ ! -f "${RUN_DIR}/meta/termination-reason.txt" ]]; then
                printf '%s\n' "accepted" > "${RUN_DIR}/meta/termination-reason.txt"
                capture_accept_snapshot
                echo "Accepted validation observed; stopping remaining agents..."
                for pid in "${PIDS[@]}"; do
                    if kill -0 "${pid}" 2>/dev/null; then
                        kill -TERM "${pid}" 2>/dev/null || true
                    fi
                done
                sleep 1
                for pid in "${PIDS[@]}"; do
                    if kill -0 "${pid}" 2>/dev/null; then
                        kill -KILL "${pid}" 2>/dev/null || true
                    fi
                done
            fi
            return 0
        fi

        alive=0
        for pid in "${PIDS[@]}"; do
            if kill -0 "${pid}" 2>/dev/null; then
                alive=1
                break
            fi
        done
        if [[ "${alive}" -eq 0 ]]; then
            return 0
        fi
        sleep 1
    done
}

watch_for_accepted_validation &
WATCHER_PID=$!

all_exit_codes=()
for i in $(seq 1 "${AGENT_COUNT}"); do
    idx=$((i - 1))
    set +e
    wait "${PIDS[$idx]}"
    exit_code=$?
    set -e
    all_exit_codes+=("${exit_code}")
    printf '%s\n' "${exit_code}" > "${RUN_DIR}/agent-${i}/exit_code"
    echo "Agent-${i} exited with code ${exit_code}"
done

wait "${WATCHER_PID}" 2>/dev/null || true
[[ -f "${RUN_DIR}/meta/termination-reason.txt" ]] || printf '%s\n' "all-agents-exited" > "${RUN_DIR}/meta/termination-reason.txt"

printf '%s\n' "${all_exit_codes[@]}" > "${RUN_DIR}/meta/exit_codes.txt"
cp "${WORLD_JSON}" "${RUN_DIR}/meta/world-truth.json"
if [[ -f "${RUN_DIR}/meta/results-at-accept.txt" ]]; then
    cp "${RUN_DIR}/meta/results-at-accept.txt" "${RUN_DIR}/workspace/results.txt"
fi

echo
echo "All agents finished."
echo "Run complete: ${RUN_DIR}"
