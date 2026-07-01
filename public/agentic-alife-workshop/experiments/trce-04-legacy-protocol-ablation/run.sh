#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

ENV_INPUT="${1:-.env.gpt-5.4-pool}"
CONDITION="${2:-A}"
MAX_TURNS="${3:-20}"
MAX_GENS="${4:-12}"
SAFETY_TIMEOUT="${5:-1800}"
PROMPT_PROFILE="${6:-minimal}"
KEEP_WORKSPACES="${KEEP_WORKSPACES:-0}"
WORLD_ITEM_COUNT="${WORLD_ITEM_COUNT:-24}"

if [[ "${ENV_INPUT}" = /* ]]; then
    ENV_FILE="${ENV_INPUT}"
else
    ENV_FILE="${PROJECT_ROOT}/${ENV_INPUT}"
fi

if [[ ! -f "${ENV_FILE}" ]]; then
    echo "ERROR: env file not found: ${ENV_FILE}"
    exit 1
fi

# shellcheck disable=SC1090
source "${ENV_FILE}"

if [[ -z "${QUINE_API_TYPE:-}" ]]; then
    echo "ERROR: QUINE_API_TYPE not set after sourcing ${ENV_FILE}"
    exit 1
fi

TIMEOUT_BIN="$(command -v gtimeout || true)"
if [[ -z "${TIMEOUT_BIN}" ]]; then
    TIMEOUT_BIN="$(command -v timeout || true)"
fi
TIMEOUT_MODE="bin"
if [[ -z "${TIMEOUT_BIN}" ]]; then
    TIMEOUT_BIN="$(command -v python3 || true)"
    TIMEOUT_MODE="python"
fi
if [[ -z "${TIMEOUT_BIN}" ]]; then
    echo "ERROR: timeout command not found (expected gtimeout, timeout, or python3)"
    exit 1
fi

case "${CONDITION}" in
    A|B|C) ;;
    *)
        echo "ERROR: condition must be A, B, or C"
        exit 1
        ;;
esac

ENV_TAG="$(basename "${ENV_FILE}")"
ENV_TAG="${ENV_TAG#.env.}"
RUN_BASENAME="$(date +%Y%m%d-%H%M%S)-${ENV_TAG}-cond${CONDITION}-${MAX_TURNS}turns"
RUNID="${RUN_BASENAME}"
RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"
attempt=1
while ! mkdir "${RUN_DIR}" 2>/dev/null; do
    attempt=$((attempt + 1))
    RUNID="${RUN_BASENAME}-r${attempt}"
    RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"
done
BASE_WORKSPACE="${RUN_DIR}/base-workspace"
GENERATIONS_DIR="${RUN_DIR}/generations"
META_DIR="${RUN_DIR}/meta"
SHARED_WORKSPACE="${RUN_DIR}/shared-workspace"

mkdir -p "${BASE_WORKSPACE}" "${GENERATIONS_DIR}" "${META_DIR}"

PROMPT_SOURCE="${SCRIPT_DIR}/prompt.md"
if [[ "${PROMPT_PROFILE}" != "minimal" ]]; then
    PROMPT_SOURCE="${SCRIPT_DIR}/prompt.${PROMPT_PROFILE}.md"
fi
if [[ ! -f "${PROMPT_SOURCE}" ]]; then
    echo "ERROR: prompt profile not found: ${PROMPT_PROFILE}"
    exit 1
fi

PROMPT_SURFACE_KIND="prompt-control"
case "${PROMPT_PROFILE}" in
    neutral)
        PROMPT_SURFACE_KIND="neutral-lineage"
        ;;
    sham)
        PROMPT_SURFACE_KIND="sham-artifact"
        ;;
    birth-uptake)
        PROMPT_SURFACE_KIND="birth-uptake-split"
        ;;
esac

cp "${PROMPT_SOURCE}" "${META_DIR}/prompt-template.md"
printf '%s\n' "${ENV_FILE}" > "${META_DIR}/env-file.txt"
printf '%s\n' "${QUINE_MODEL_ID:-unknown}" > "${META_DIR}/model-id.txt"
printf '%s\n' "${CONDITION}" > "${META_DIR}/condition.txt"
printf '%s\n' "${MAX_TURNS}" > "${META_DIR}/max-turns.txt"
printf '%s\n' "${MAX_GENS}" > "${META_DIR}/max-gens.txt"
printf '%s\n' "${SAFETY_TIMEOUT}" > "${META_DIR}/safety-timeout.txt"
printf '%s\n' "${PROMPT_PROFILE}" > "${META_DIR}/prompt-profile.txt"
printf '%s\n' "${PROMPT_SURFACE_KIND}" > "${META_DIR}/surface-kind.txt"
printf '%s\n' "0" > "${META_DIR}/exec-enabled.txt"
printf '%s\n' "0" > "${META_DIR}/fork-enabled.txt"
printf '%s\n' "disabled-via-QUINE_FORK_ENABLED=0" > "${META_DIR}/fork-policy.txt"
printf '%s\n' "${WORLD_ITEM_COUNT}" > "${META_DIR}/world-item-count.txt"

WORLD_SPEC_BUILD="$(mktemp "${TMPDIR:-/tmp}/world-spec.XXXXXX")"
trap 'rm -f "${WORLD_SPEC_BUILD}"' EXIT
"${SCRIPT_DIR}/setup/init-world.sh" "${BASE_WORKSPACE}" "${WORLD_SPEC_BUILD}" "${WORLD_ITEM_COUNT}"

QUINE_BIN="${META_DIR}/quine"
WORLD_BIN="${META_DIR}/world"
echo "Building quine..."
(cd "${PROJECT_ROOT}" && go build -o "${QUINE_BIN}" ./cmd/quine/)
echo "Building world..."
WORLD_SPEC_B64="$(base64 < "${WORLD_SPEC_BUILD}" | tr -d '\n')"
(cd "${PROJECT_ROOT}" && go build -ldflags "-X github.com/kehao95/quine/internal/world.embeddedSpecBase64=${WORLD_SPEC_B64}" -o "${WORLD_BIN}" ./cmd/world/)

make_workspace_copy() {
    local target="$1"
    rm -rf "${target}"
    mkdir -p "${target}"
    rsync -a --delete "${BASE_WORKSPACE}/" "${target}/"
}

apply_read_only() {
    local target="$1"
    find "${target}" -type d -exec chmod a-w '{}' +
    find "${target}" -type f -exec chmod a-w '{}' +
}

make_writable_for_cleanup() {
    local target="$1"
    if [[ ! -e "${target}" ]]; then
        return
    fi
    find "${target}" -type d -exec chmod u+w '{}' + 2>/dev/null || true
    find "${target}" -type f -exec chmod u+w '{}' + 2>/dev/null || true
}

temporarily_lock_run_tree() {
    local target="$1"
    find "${target}" -type d -exec chmod a-w '{}' + 2>/dev/null || true
    find "${target}" -type f -exec chmod a-w '{}' + 2>/dev/null || true
}

list_artifacts() {
    local workspace="$1"
    local output="$2"
    (
        cd "${workspace}"
        find . -type f | sed 's#^\./##' | sort
    ) > "${output}"
}

snapshot_artifacts() {
    local workspace="$1"
    local target="$2"
    rm -rf "${target}"
    mkdir -p "${target}"
    while IFS= read -r rel; do
        [[ -z "${rel}" ]] && continue
        mkdir -p "${target}/$(dirname "${rel}")"
        cp "${workspace}/${rel}" "${target}/${rel}"
    done < <(cd "${workspace}" && find . -type f | sed 's#^\./##' | sort)
}

install_runtime_world_spec() {
    local quine_data_dir="$1"
    mkdir -p "${quine_data_dir}/world"
    cp "${WORLD_SPEC_BUILD}" "${quine_data_dir}/world/world.json"
}

apply_sham_scramble() {
    local workspace="$1"
    local output="$2"
    python3 - "${workspace}" "${output}" <<'PY'
from pathlib import Path
import shutil
import sys

workspace = Path(sys.argv[1])
output = Path(sys.argv[2])
rows = []
rewrites = []
preserve_roots = {"archive", "lineage_state"}

targets = []
for path in sorted(workspace.rglob("*")):
    if not path.is_file():
        continue
    rel = path.relative_to(workspace).as_posix()
    parts = rel.split("/")
    if parts[0] == "archive":
        continue
    targets.append((rel, path))

if targets:
    for index, (rel, path) in enumerate(targets, start=1):
        original = path.read_text(encoding="utf-8", errors="ignore")
        replacement_rel = f"lineage_state/artifact-{index:03d}.txt"
        payload = (
            f"artifact_slot={index:03d}\n"
            f"bytes={len(original.encode('utf-8'))}\n"
            f"lines={len(original.splitlines())}\n"
            f"tokens={len(original.split())}\n"
        )
        rewrites.append((replacement_rel, payload))
        rows.append(f"{rel}\t{replacement_rel}\t{len(original.encode('utf-8'))}\t{len(original.splitlines())}\t{len(original.split())}")

if targets:
    for child in workspace.iterdir():
        if child.name in preserve_roots:
            continue
        if child.is_dir():
            shutil.rmtree(child)
        else:
            child.unlink()
    target_root = workspace / "lineage_state"
    if target_root.exists():
        shutil.rmtree(target_root)
    target_root.mkdir(parents=True, exist_ok=True)
    for replacement_rel, payload in rewrites:
        replacement_path = workspace / replacement_rel
        replacement_path.parent.mkdir(parents=True, exist_ok=True)
        replacement_path.write_text(payload, encoding="utf-8")

output.write_text("\n".join(rows) + ("\n" if rows else ""), encoding="utf-8")
PY
}

prepare_birth_uptake_surface() {
    local workspace="$1"
    mkdir -p "${workspace}/lineage_state/birth" "${workspace}/lineage_state/uptake"
}

render_prompt() {
    local target="$1"
    local gen="$2"
    local prev="$3"
    local pid_value="$$"
    sed \
        -e "s|{{GEN}}|${gen}|g" \
        -e "s|{{PREV_GEN}}|${prev}|g" \
        -e "s|{{PID}}|${pid_value}|g" \
        -e "s|{{MAX_TURNS}}|${MAX_TURNS}|g" \
        "${META_DIR}/prompt-template.md" > "${target}"
}

if [[ "${CONDITION}" = "A" ]]; then
    make_workspace_copy "${SHARED_WORKSPACE}"
fi

echo "============================================================"
echo "Exp 2.6: Legacy Protocol Ablation"
echo "Run ID:      ${RUNID}"
echo "Env file:    ${ENV_FILE}"
echo "Model:       ${QUINE_MODEL_ID:-unknown}"
echo "Condition:   ${CONDITION}"
echo "Max turns:   ${MAX_TURNS}"
echo "Max Gens:    ${MAX_GENS}"
echo "Safety time: ${SAFETY_TIMEOUT}s"
echo "Prompt:      ${PROMPT_PROFILE}"
echo "Surface:     ${PROMPT_SURFACE_KIND}"
echo "Exec:        disabled"
echo "Fork:        disabled via QUINE_FORK_ENABLED=0"
echo "Run dir:     ${RUN_DIR}"
echo "============================================================"

for gen in $(seq 1 "${MAX_GENS}"); do
    slot="$(printf '%02d' "${gen}")"
    prev="$((gen - 1))"
    gen_dir="${GENERATIONS_DIR}/${slot}"
    mkdir -p "${gen_dir}/meta" "${gen_dir}/quine" "${gen_dir}/snapshots"
    install_runtime_world_spec "${gen_dir}/quine"
    prompt_file="${gen_dir}/meta/prompt.md"
    render_prompt "${prompt_file}" "${gen}" "${prev}"

    case "${CONDITION}" in
        A)
            workspace="${SHARED_WORKSPACE}"
            ;;
        B|C)
            workspace="${gen_dir}/workspace"
            make_workspace_copy "${workspace}"
            if [[ "${CONDITION}" = "C" ]]; then
                apply_read_only "${workspace}"
            fi
            ;;
    esac

    if [[ "${PROMPT_PROFILE}" = "birth-uptake" ]]; then
        prepare_birth_uptake_surface "${workspace}"
    fi

    list_artifacts "${workspace}" "${gen_dir}/meta/pre-artifacts.txt"
    snapshot_artifacts "${workspace}" "${gen_dir}/snapshots/pre"

    start_epoch="$(date +%s)"
    set +e
    (
        cd "${workspace}"
        exec 3> "${gen_dir}/meta/stdout.txt"
        exec 4> "${gen_dir}/meta/stderr.txt"
        if [[ "${CONDITION}" = "C" ]]; then
            temporarily_lock_run_tree "${RUN_DIR}"
            make_writable_for_cleanup "${gen_dir}/quine"
        fi
        if [[ "${TIMEOUT_MODE}" = "bin" ]]; then
            PATH="${META_DIR}:${PATH}" \
            QUINE_DATA_DIR="${gen_dir}/quine" QUINE_RETENTION_DIR="${gen_dir}/quine/log" \
            QUINE_MODEL_ID="${QUINE_MODEL_ID}" \
            QUINE_MAX_TURNS="${MAX_TURNS}" \
            QUINE_TURN_EXHAUSTION_POLICY="hard_fail" \
            QUINE_OUTPUT_TRUNCATE="4096" \
            QUINE_FORK_ENABLED="0" \
            QUINE_EXEC_ENABLED="0" \
            "${TIMEOUT_BIN}" "${SAFETY_TIMEOUT}s" \
                "${QUINE_BIN}" "$(cat "${prompt_file}")" \
                < /dev/null \
                >&3 \
                2>&4
        else
            PATH="${META_DIR}:${PATH}" \
            QUINE_DATA_DIR="${gen_dir}/quine" QUINE_RETENTION_DIR="${gen_dir}/quine/log" \
            QUINE_MODEL_ID="${QUINE_MODEL_ID}" \
            QUINE_MAX_TURNS="${MAX_TURNS}" \
            QUINE_TURN_EXHAUSTION_POLICY="hard_fail" \
            QUINE_OUTPUT_TRUNCATE="4096" \
            QUINE_FORK_ENABLED="0" \
            QUINE_EXEC_ENABLED="0" \
            "${TIMEOUT_BIN}" -c '
import os
import subprocess
import sys

timeout = int(sys.argv[1])
cmd = sys.argv[2:]
try:
    result = subprocess.run(cmd, stdin=subprocess.DEVNULL, stdout=sys.stdout, stderr=sys.stderr, timeout=timeout, check=False)
except subprocess.TimeoutExpired:
    sys.exit(124)
sys.exit(result.returncode)
' "${SAFETY_TIMEOUT}" "${QUINE_BIN}" "$(cat "${prompt_file}")" \
                >&3 \
                2>&4
        fi
        exit_code=$?
        if [[ "${CONDITION}" = "C" ]]; then
            make_writable_for_cleanup "${RUN_DIR}"
        fi
        exit "${exit_code}"
    )
    exit_code=$?
    set -e
    end_epoch="$(date +%s)"
    duration="$((end_epoch - start_epoch))"

    timed_out=0
    if [[ "${exit_code}" -eq 124 ]]; then
        timed_out=1
    fi

    printf '%s\n' "${exit_code}" > "${gen_dir}/meta/exit-code.txt"
    printf '%s\n' "${duration}" > "${gen_dir}/meta/duration-seconds.txt"
    printf '%s\n' "${timed_out}" > "${gen_dir}/meta/timed-out.txt"
    printf '%s\n' "${workspace}" > "${gen_dir}/meta/workspace.txt"

    list_artifacts "${workspace}" "${gen_dir}/meta/post-artifacts.txt"
    snapshot_artifacts "${workspace}" "${gen_dir}/snapshots/post"
    comm -13 "${gen_dir}/meta/pre-artifacts.txt" "${gen_dir}/meta/post-artifacts.txt" > "${gen_dir}/meta/new-artifacts.txt" || true
    comm -23 "${gen_dir}/meta/pre-artifacts.txt" "${gen_dir}/meta/post-artifacts.txt" > "${gen_dir}/meta/removed-artifacts.txt" || true

    if [[ "${timed_out}" -eq 1 ]]; then
        echo "Generation ${slot}: safety-timeout after ${duration}s"
    else
        echo "Generation ${slot}: exit ${exit_code} after ${duration}s"
    fi

    if [[ "${PROMPT_PROFILE}" = "sham" && "${CONDITION}" = "A" ]]; then
        apply_sham_scramble "${workspace}" "${gen_dir}/meta/sham-scramble.tsv"
    fi
done

"${SCRIPT_DIR}/analysis/analyze.sh" "${RUN_DIR}" | tee "${META_DIR}/summary.md"
cp "${WORLD_SPEC_BUILD}" "${META_DIR}/world.json"

if [[ "${KEEP_WORKSPACES}" != "1" ]]; then
    make_writable_for_cleanup "${BASE_WORKSPACE}"
    make_writable_for_cleanup "${SHARED_WORKSPACE}"
    find "${GENERATIONS_DIR}" -type d -name workspace -print0 | while IFS= read -r -d '' workspace_dir; do
        make_writable_for_cleanup "${workspace_dir}"
    done
    rm -rf "${BASE_WORKSPACE}" "${SHARED_WORKSPACE}"
    find "${GENERATIONS_DIR}" -type d -name workspace -prune -exec rm -rf '{}' +
    find "${GENERATIONS_DIR}" -type d -path '*/quine/agent' -prune -exec rm -rf '{}' +
fi

echo
echo "Run complete: ${RUN_DIR}"
