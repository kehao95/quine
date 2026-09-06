#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF' >&2
usage: ./tests/validate.sh --change <substrate|runtime|usage|discovery|necessity> [options]

options:
  --from-git                  infer the minimum required layer from current git changes
  --runtime <test>            add a runtime test name (repeatable); use "all" for the full runtime suite
  --usage <evaluation>        add an L3 usage evaluation id (repeatable); use "all" for that full layer
  --discovery <evaluation>    add an L4 discovery evaluation id (repeatable); use "all" for that full layer
  --necessity <evaluation>    add an L5 necessity evaluation id (repeatable); use "all" for that full layer
  --model <id>                override QUINE_MODEL_ID for model-layer runs
  --dry-run                   print the plan without executing commands

examples:
  ./tests/validate.sh --change substrate
  ./tests/validate.sh --change runtime --runtime test_exit_success
  ./tests/validate.sh --change usage --runtime test_fd4_delivery --usage stdin-explicit-handoff
  ./tests/validate.sh --change discovery --runtime test_workspace_overlay_commit --usage workspace-overlay-relative-path-explicit --discovery sandbox-unknown-format-boldness
  ./tests/validate.sh --change necessity --runtime <test> --usage <evaluation> --discovery <evaluation> --necessity <evaluation>
EOF
  exit 2
}

CHANGE=""
FROM_GIT=0
DRY_RUN=0
MODEL_OVERRIDE=""
declare -a RUNTIME_TESTS=()
declare -a USAGE_EVALUATIONS=()
declare -a DISCOVERY_EVALUATIONS=()
declare -a NECESSITY_EVALUATIONS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --change)
      [[ $# -ge 2 ]] || usage
      CHANGE="$2"
      shift 2
      ;;
    --from-git)
      FROM_GIT=1
      shift
      ;;
    --runtime)
      [[ $# -ge 2 ]] || usage
      RUNTIME_TESTS+=("$2")
      shift 2
      ;;
    --usage)
      [[ $# -ge 2 ]] || usage
      USAGE_EVALUATIONS+=("$2")
      shift 2
      ;;
    --discovery)
      [[ $# -ge 2 ]] || usage
      DISCOVERY_EVALUATIONS+=("$2")
      shift 2
      ;;
    --necessity)
      [[ $# -ge 2 ]] || usage
      NECESSITY_EVALUATIONS+=("$2")
      shift 2
      ;;
    --model)
      [[ $# -ge 2 ]] || usage
      MODEL_OVERRIDE="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    *)
      usage
      ;;
  esac
done

if [[ -z "$CHANGE" && "$FROM_GIT" -ne 1 ]]; then
  usage
fi

if [[ "$FROM_GIT" -eq 1 ]]; then
  if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "FATAL: --from-git requires a git worktree" >&2
    exit 1
  fi
  CHANGE="$(
    python3 - <<'PY'
from pathlib import Path
import subprocess

out = subprocess.run(
    ["git", "status", "--porcelain"],
    check=True,
    capture_output=True,
    text=True,
).stdout.splitlines()
paths = []
for line in out:
    if not line:
        continue
    path = line[3:]
    if " -> " in path:
        path = path.split(" -> ", 1)[1]
    paths.append(path)

layer = "substrate"
for path in paths:
    if path.startswith(("tests/model/", "scripts/check-model-evaluations.sh")):
        layer = "discovery"
        break
    if path.startswith(("tests/runtime/",)):
        layer = "runtime" if layer == "substrate" else layer
    if path in {"EVALUATION.md", "DEVELOPMENT.md", "README.md", "tests/README.md", "Makefile"}:
        layer = "discovery"
    if path.startswith(("cmd/quine/", "internal/runtime/", "internal/tools/", "internal/llm/")):
        if layer == "substrate":
            layer = "runtime"
    if path.endswith(".go") and not path.endswith("_test.go") and layer == "substrate":
        layer = "runtime"
print(layer)
PY
  )"
fi

case "$CHANGE" in
  substrate) REQUIRED_LEVEL=1 ;;
  runtime) REQUIRED_LEVEL=2 ;;
  usage) REQUIRED_LEVEL=3 ;;
  discovery) REQUIRED_LEVEL=4 ;;
  necessity) REQUIRED_LEVEL=5 ;;
  *) echo "FATAL: unknown change layer: $CHANGE" >&2; exit 1 ;;
esac

need_runtime=0
need_usage=0
need_discovery=0
need_necessity=0
[[ $REQUIRED_LEVEL -ge 2 ]] && need_runtime=1
[[ $REQUIRED_LEVEL -ge 3 ]] && need_usage=1
[[ $REQUIRED_LEVEL -ge 4 ]] && need_discovery=1
[[ $REQUIRED_LEVEL -ge 5 ]] && need_necessity=1

if [[ $need_runtime -eq 1 && ${#RUNTIME_TESTS[@]} -eq 0 ]]; then
  echo "FATAL: runtime validation requires at least one --runtime test or --runtime all" >&2
  exit 1
fi
if [[ $need_usage -eq 1 && ${#USAGE_EVALUATIONS[@]} -eq 0 ]]; then
  echo "FATAL: usage validation requires at least one --usage evaluation or --usage all" >&2
  exit 1
fi
if [[ $need_discovery -eq 1 && ${#DISCOVERY_EVALUATIONS[@]} -eq 0 ]]; then
  echo "FATAL: discovery validation requires at least one --discovery evaluation or --discovery all" >&2
  exit 1
fi
if [[ $need_necessity -eq 1 && ${#NECESSITY_EVALUATIONS[@]} -eq 0 ]]; then
  echo "FATAL: necessity validation requires at least one --necessity evaluation or --necessity all" >&2
  exit 1
fi

print_layer() {
  local level="$1"
  local name="$2"
  local needed="$3"
  if [[ "$needed" -eq 1 ]]; then
    printf '  L%s %-14s required\n' "$level" "$name"
  else
    printf '  L%s %-14s visible, not required\n' "$level" "$name"
  fi
}

run_substrate_suite() {
  local packages=()
  local pkg
  while IFS= read -r pkg; do
    [[ -n "$pkg" ]] || continue
    packages+=("$pkg")
  done < <(./tests/active-go-packages.sh)

  if [[ "${#packages[@]}" -eq 0 ]]; then
    echo "FATAL: active substrate package selector returned no packages" >&2
    exit 1
  fi

  go test ${GOFLAGS:-} -count=1 "${packages[@]}"
}

echo "Validation ladder for change layer: $CHANGE"
print_layer 1 substrate 1
print_layer 2 runtime "$need_runtime"
print_layer 3 usage "$need_usage"
print_layer 4 discovery "$need_discovery"
print_layer 5 necessity "$need_necessity"
echo ""

if [[ ${#RUNTIME_TESTS[@]} -gt 0 ]]; then
  echo "Runtime selection: ${RUNTIME_TESTS[*]}"
fi
if [[ ${#USAGE_EVALUATIONS[@]} -gt 0 ]]; then
  echo "Usage selection: ${USAGE_EVALUATIONS[*]}"
fi
if [[ ${#DISCOVERY_EVALUATIONS[@]} -gt 0 ]]; then
  echo "Discovery selection: ${DISCOVERY_EVALUATIONS[*]}"
fi
if [[ ${#NECESSITY_EVALUATIONS[@]} -gt 0 ]]; then
  echo "Necessity selection: ${NECESSITY_EVALUATIONS[*]}"
fi
if [[ -n "$MODEL_OVERRIDE" ]]; then
  echo "Model override: $MODEL_OVERRIDE"
fi

if [[ "$DRY_RUN" -eq 1 ]]; then
  exit 0
fi

run_substrate_suite

if [[ $need_runtime -eq 1 ]]; then
  if [[ ${#RUNTIME_TESTS[@]} -eq 1 && "${RUNTIME_TESTS[0]}" == "all" ]]; then
    ./tests/runtime/run.sh
  else
    ./tests/runtime/run.sh "${RUNTIME_TESTS[@]}"
  fi
fi

run_model_layer() {
  local layer="$1"
  shift
  local evaluation
  for evaluation in "$@"; do
    if [[ "$evaluation" == "all" ]]; then
      if [[ -n "$MODEL_OVERRIDE" ]]; then
        ./tests/model/run.sh "$layer" "$MODEL_OVERRIDE"
      else
        ./tests/model/run.sh "$layer"
      fi
      return
    fi
    if [[ -n "$MODEL_OVERRIDE" ]]; then
      ./tests/model/run.sh "$evaluation" "$MODEL_OVERRIDE"
    else
      ./tests/model/run.sh "$evaluation"
    fi
  done
}

if [[ $need_usage -eq 1 ]]; then
  run_model_layer usage "${USAGE_EVALUATIONS[@]}"
fi

if [[ $need_discovery -eq 1 ]]; then
  run_model_layer discovery "${DISCOVERY_EVALUATIONS[@]}"
fi

if [[ $need_necessity -eq 1 ]]; then
  run_model_layer necessity "${NECESSITY_EVALUATIONS[@]}"
fi
