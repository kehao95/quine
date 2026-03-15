#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF' >&2
usage: ./tests/validate.sh --change <substrate|runtime|instructional|emergent> [options]

options:
  --from-git                  infer the minimum required layer from current git changes
  --runtime <test>            add a runtime test name (repeatable); use "all" for the full runtime suite
  --instructional <scenario>  add an instructional scenario id (repeatable); use "all" for that full layer
  --emergent <scenario>       add an emergent scenario id (repeatable); use "all" for that full layer
  --model <id>                override QUINE_MODEL_ID for model-layer runs
  --dry-run                   print the plan without executing commands

examples:
  ./tests/validate.sh --change substrate
  ./tests/validate.sh --change runtime --runtime test_exit_success
  ./tests/validate.sh --change instructional --runtime test_fd4_delivery --instructional stdin
  ./tests/validate.sh --change emergent --runtime test_workspace_overlay_commit --instructional workspace-shadow --emergent workspace-shadow-emergent
EOF
  exit 2
}

CHANGE=""
FROM_GIT=0
DRY_RUN=0
MODEL_OVERRIDE=""
declare -a RUNTIME_TESTS=()
declare -a INSTRUCTIONAL_SCENARIOS=()
declare -a EMERGENT_SCENARIOS=()

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
    --instructional)
      [[ $# -ge 2 ]] || usage
      INSTRUCTIONAL_SCENARIOS+=("$2")
      shift 2
      ;;
    --emergent)
      [[ $# -ge 2 ]] || usage
      EMERGENT_SCENARIOS+=("$2")
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
    if path.startswith(("tests/model/", "scripts/check-model-scenarios.sh")):
        layer = "emergent"
        break
    if path.startswith(("tests/runtime/",)):
        layer = "runtime" if layer == "substrate" else layer
    if path in {"TESTING.md", "DEVELOPMENT.md", "README.md", "tests/README.md", "Makefile"}:
        layer = "emergent"
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
  instructional) REQUIRED_LEVEL=3 ;;
  emergent) REQUIRED_LEVEL=4 ;;
  *) echo "FATAL: unknown change layer: $CHANGE" >&2; exit 1 ;;
esac

need_runtime=0
need_instructional=0
need_emergent=0
[[ $REQUIRED_LEVEL -ge 2 ]] && need_runtime=1
[[ $REQUIRED_LEVEL -ge 3 ]] && need_instructional=1
[[ $REQUIRED_LEVEL -ge 4 ]] && need_emergent=1

if [[ $need_runtime -eq 1 && ${#RUNTIME_TESTS[@]} -eq 0 ]]; then
  echo "FATAL: runtime validation requires at least one --runtime test or --runtime all" >&2
  exit 1
fi
if [[ $need_instructional -eq 1 && ${#INSTRUCTIONAL_SCENARIOS[@]} -eq 0 ]]; then
  echo "FATAL: instructional validation requires at least one --instructional scenario or --instructional all" >&2
  exit 1
fi
if [[ $need_emergent -eq 1 && ${#EMERGENT_SCENARIOS[@]} -eq 0 ]]; then
  echo "FATAL: emergent validation requires at least one --emergent scenario or --emergent all" >&2
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

echo "Validation ladder for change layer: $CHANGE"
print_layer 1 substrate 1
print_layer 2 runtime "$need_runtime"
print_layer 3 instructional "$need_instructional"
print_layer 4 emergent "$need_emergent"
echo ""

if [[ ${#RUNTIME_TESTS[@]} -gt 0 ]]; then
  echo "Runtime selection: ${RUNTIME_TESTS[*]}"
fi
if [[ ${#INSTRUCTIONAL_SCENARIOS[@]} -gt 0 ]]; then
  echo "Instructional selection: ${INSTRUCTIONAL_SCENARIOS[*]}"
fi
if [[ ${#EMERGENT_SCENARIOS[@]} -gt 0 ]]; then
  echo "Emergent selection: ${EMERGENT_SCENARIOS[*]}"
fi
if [[ -n "$MODEL_OVERRIDE" ]]; then
  echo "Model override: $MODEL_OVERRIDE"
fi

if [[ "$DRY_RUN" -eq 1 ]]; then
  exit 0
fi

go test ./...

if [[ $need_runtime -eq 1 ]]; then
  ./tests/runtime/run.sh "${RUNTIME_TESTS[@]}"
fi

run_model_layer() {
  local layer="$1"
  shift
  local scenario
  for scenario in "$@"; do
    if [[ "$scenario" == "all" ]]; then
      if [[ -n "$MODEL_OVERRIDE" ]]; then
        ./tests/model/run.sh "$layer" "$MODEL_OVERRIDE"
      else
        ./tests/model/run.sh "$layer"
      fi
      return
    fi
    if [[ -n "$MODEL_OVERRIDE" ]]; then
      ./tests/model/run.sh "$scenario" "$MODEL_OVERRIDE"
    else
      ./tests/model/run.sh "$scenario"
    fi
  done
}

if [[ $need_instructional -eq 1 ]]; then
  run_model_layer instructional "${INSTRUCTIONAL_SCENARIOS[@]}"
fi

if [[ $need_emergent -eq 1 ]]; then
  run_model_layer emergent "${EMERGENT_SCENARIOS[@]}"
fi
