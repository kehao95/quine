#!/usr/bin/env bash

set -euo pipefail

STRICT=0
REQUIRE_BASELINES=0
PRUNE_RUN_TREE=0
SCENARIO_FILTER=""
LAYER_FILTER=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --strict)
      STRICT=1
      shift
      ;;
    --require-baselines)
      REQUIRE_BASELINES=1
      shift
      ;;
    --prune-run-tree)
      PRUNE_RUN_TREE=1
      shift
      ;;
    --scenario)
      [[ $# -ge 2 ]] || { echo "usage: $0 [--strict] [--require-baselines] [--prune-run-tree] [--scenario <id>] [--layer <instructional|emergent>]" >&2; exit 2; }
      SCENARIO_FILTER="$2"
      shift 2
      ;;
    --layer)
      [[ $# -ge 2 ]] || { echo "usage: $0 [--strict] [--require-baselines] [--prune-run-tree] [--scenario <id>] [--layer <instructional|emergent>]" >&2; exit 2; }
      LAYER_FILTER="$2"
      shift 2
      ;;
    *)
      echo "usage: $0 [--strict] [--require-baselines] [--prune-run-tree] [--scenario <id>] [--layer <instructional|emergent>]" >&2
      exit 2
      ;;
  esac
done

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

STRICT="$STRICT" \
REQUIRE_BASELINES="$REQUIRE_BASELINES" \
PRUNE_RUN_TREE="$PRUNE_RUN_TREE" \
SCENARIO_FILTER="$SCENARIO_FILTER" \
LAYER_FILTER="$LAYER_FILTER" \
python3 - <<'PY'
from pathlib import Path
import os
import re
import shutil
import sys
try:
    import tomllib
except ModuleNotFoundError:
    import tomli as tomllib

root = Path(".")
registry_path = root / "tests" / "model" / "scenarios.toml"
runner_path = root / "tests" / "model" / "run.sh"
require_baselines = os.environ.get("REQUIRE_BASELINES") == "1"
prune_run_tree = os.environ.get("PRUNE_RUN_TREE") == "1"
scenario_filter = os.environ.get("SCENARIO_FILTER", "").strip()
layer_filter = os.environ.get("LAYER_FILTER", "").strip()

issues = 0
warnings = 0
required_run_files = (
    "prompt-used.md",
    "score.txt",
    "stdout.txt",
    "stderr.txt",
    "tape.jsonl",
    "tape.log",
)

def report(msg: str) -> None:
    global issues
    print(msg)
    issues += 1

def warn(msg: str) -> None:
    global warnings
    print(msg)
    warnings += 1

def prune(msg: str) -> None:
    print(msg)

if not registry_path.exists():
    report(f"MISSING registry: {registry_path}")
    print(f"Model-scenario audit: scenarios=0 scored=0 issues={issues} warnings={warnings}")
    sys.exit(1)

data = tomllib.loads(registry_path.read_text(encoding="utf-8"))
entries = data.get("scenario", [])
if not isinstance(entries, list):
    report(f"INVALID registry root: {registry_path}")
    print(f"Model-scenario audit: scenarios=0 scored=0 issues={issues} warnings={warnings}")
    sys.exit(1)

seen_ids = set()
scenarios = []
for entry in entries:
    sid = entry.get("id", "").strip()
    layer = entry.get("layer", "").strip()
    if not sid:
      report("REGISTRY scenario missing id")
      continue
    if sid in seen_ids:
      report(f"DUPLICATE registry scenario id: {sid}")
      continue
    seen_ids.add(sid)
    if layer not in {"instructional", "emergent"}:
      report(f"INVALID layer for scenario {sid}: {layer}")
      continue
    scenarios.append(entry)

if layer_filter and layer_filter not in {"instructional", "emergent"}:
    report(f"INVALID layer filter: {layer_filter}")

if scenario_filter and scenario_filter not in seen_ids:
    report(f"UNKNOWN scenario filter: {scenario_filter}")

if scenario_filter:
    scenarios = [entry for entry in scenarios if entry["id"] == scenario_filter]
if layer_filter:
    scenarios = [entry for entry in scenarios if entry["layer"] == layer_filter]

runner_text = runner_path.read_text(encoding="utf-8")
score_match = re.search(
    r"^score_scenario\(\) \{\n(.*?)(?=^# Check if a marker string appears in stdout)",
    runner_text,
    re.S | re.M,
)
if not score_match:
    report(f"MISSING score_scenario body: {runner_path}")
    scored_ids = set()
else:
    scored_ids = set()
    for label_group in re.findall(r"^\s*([a-z0-9-]+(?:\|[a-z0-9-]+)*)\)\s*$", score_match.group(1), re.M):
        for label in label_group.split("|"):
            scored_ids.add(label)

def scenario_prompt_path(entry: dict) -> Path:
    return root / "tests" / "model" / entry["layer"] / "prompts" / f'{entry["id"]}.md'

def scenario_run_root(entry: dict) -> Path:
    return root / "tests" / "model" / entry["layer"] / "runs" / entry["id"]

def layer_runs_dir(layer: str) -> Path:
    return root / "tests" / "model" / layer / "runs"

def run_is_complete(run_dir: Path) -> tuple[bool, list[str]]:
    missing = [name for name in required_run_files if not (run_dir / name).exists()]
    return len(missing) == 0, missing

def score_is_passing(score_path: Path) -> bool:
    if not score_path.exists():
        return False
    text = score_path.read_text(encoding="utf-8", errors="replace")
    if "FAIL" in text:
        return False
    match = re.search(r"Result:\s+(\d+)/(\d+)\s+passed", text)
    return bool(match and match.group(1) == match.group(2))

def run_should_keep(run_dir: Path) -> bool:
    complete, _ = run_is_complete(run_dir)
    return complete and score_is_passing(run_dir / "score.txt")

def refresh_latest_symlink(layer: str) -> None:
    runs_dir = layer_runs_dir(layer)
    latest = runs_dir / "latest"
    newest_rel = None
    newest_mtime = None
    if not runs_dir.exists():
        return
    for scenario_dir in sorted(runs_dir.iterdir()):
        if not scenario_dir.is_dir() or scenario_dir.name == "latest":
            continue
        for run_dir in sorted(p for p in scenario_dir.iterdir() if p.is_dir()):
            rel = run_dir.relative_to(runs_dir)
            mtime = run_dir.stat().st_mtime
            if newest_rel is None or newest_mtime is None or mtime > newest_mtime or (mtime == newest_mtime and str(rel) > str(newest_rel)):
                newest_rel = rel
                newest_mtime = mtime
    if newest_rel is None:
        if latest.exists() or latest.is_symlink():
            latest.unlink()
        return
    if latest.exists() or latest.is_symlink():
        latest.unlink()
    latest.symlink_to(newest_rel)

registry_ids = {entry["id"] for entry in scenarios}

for entry in scenarios:
    prompt_path = scenario_prompt_path(entry)
    if not prompt_path.exists():
        report(f"MISSING prompt file: {prompt_path}")
    if entry["id"] not in scored_ids:
        report(f"UNSCORED scenario: {entry['id']}")

for scored_id in sorted(scored_ids):
    if scored_id not in seen_ids:
        report(f"STALE scorer with no registry scenario: {scored_id}")

if prune_run_tree:
    layers = {entry["layer"] for entry in scenarios} if scenarios else {"instructional", "emergent"}
    for layer in sorted(layers):
        runs_dir = layer_runs_dir(layer)
        if not runs_dir.exists():
            continue
        valid_ids = {entry["id"] for entry in entries if entry["layer"] == layer}
        for scenario_dir in sorted(runs_dir.iterdir()):
            if not scenario_dir.is_dir() or scenario_dir.name == "latest":
                continue
            if scenario_dir.name not in valid_ids:
                prune(f"PRUNE stale run directory: {scenario_dir}")
                shutil.rmtree(scenario_dir)
                continue
            for run_dir in sorted(p for p in scenario_dir.iterdir() if p.is_dir()):
                if run_should_keep(run_dir):
                    continue
                prune(f"PRUNE non-canonical run directory: {run_dir}")
                shutil.rmtree(run_dir)
            if not any(p.is_dir() for p in scenario_dir.iterdir()):
                prune(f"PRUNE empty scenario directory: {scenario_dir}")
                shutil.rmtree(scenario_dir)
        refresh_latest_symlink(layer)

for layer in ("instructional", "emergent"):
    runs_dir = layer_runs_dir(layer)
    if not runs_dir.exists():
        continue
    valid_ids = {entry["id"] for entry in entries if entry["layer"] == layer}
    for scenario_dir in sorted(runs_dir.iterdir()):
        if not scenario_dir.is_dir() or scenario_dir.name == "latest":
            continue
        if scenario_dir.name not in valid_ids:
            warn(f"STALE run directory with no registry scenario: {scenario_dir}")

if require_baselines:
    for entry in scenarios:
        scenario_run_dir = scenario_run_root(entry)
        if not scenario_run_dir.exists():
            report(f"MISSING baseline run directory: {scenario_run_dir}")
            continue
        passing_baselines = []
        incomplete_runs = []
        for run_dir in sorted(p for p in scenario_run_dir.iterdir() if p.is_dir()):
            complete, missing = run_is_complete(run_dir)
            if not complete:
                incomplete_runs.append((run_dir.name, missing))
                continue
            if score_is_passing(run_dir / "score.txt"):
                passing_baselines.append(run_dir.name)
        if incomplete_runs:
            preview = ", ".join(
                f"{name} [{', '.join(missing)}]"
                for name, missing in incomplete_runs[:3]
            )
            extra = ""
            if len(incomplete_runs) > 3:
                extra = f", +{len(incomplete_runs) - 3} more"
            warn(f"INCOMPLETE run artifacts: {scenario_run_dir} ({preview}{extra})")
        if not passing_baselines:
            report(f"MISSING passing baseline evidence: {scenario_run_dir}")

print(
    "Model-scenario audit: "
    f"scenarios={len(scenarios)} "
    f"scored={len(scored_ids)} "
    f"issues={issues} "
    f"warnings={warnings}"
)

sys.exit(1 if issues else 0)
PY
status=$?

if [[ "$status" -ne 0 && "$STRICT" -eq 0 ]]; then
  exit 0
fi

exit "$status"
