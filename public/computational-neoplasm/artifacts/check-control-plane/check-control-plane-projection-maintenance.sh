#!/usr/bin/env bash
#
# VALIDATOR INTENTION
# purpose:
#   - require active control-plane indexes to expose a stale-state boundary
#   - make "control plane as maintained projection index" machine-checkable
#   - support both repo-wide audit mode and targeted file validation for hooks
# repair-guidance:
#   - preserve "report drift, do not rewrite files" behavior
#   - keep the candidate set to active routing/index/control surfaces, not raw
#     evidence, retained run snapshots, or ordinary manuscript prose
#   - keep the public root README out of this active-control-plane set; its
#     human-facing skeleton is checked by check-root-entry-docs.sh
#   - if the maintained-index contract changes, update this script,
#     DEVELOPMENT.md, development/contracts/control-plane-causal-path.md,
#     scripts/README.md, .githooks/README.md, and the affected owner surfaces

set -euo pipefail

STRICT=0
TARGET_FILES=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --strict)
      STRICT=1
      shift
      ;;
    --)
      shift
      while [[ $# -gt 0 ]]; do
        TARGET_FILES+=("$1")
        shift
      done
      ;;
    *)
      TARGET_FILES+=("$1")
      shift
      ;;
  esac
done

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

PY_ARGS=("$ROOT_DIR" "$STRICT")
if [[ "${#TARGET_FILES[@]}" -gt 0 ]]; then
  PY_ARGS+=("${TARGET_FILES[@]}")
fi

python3 - "${PY_ARGS[@]}" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1]).resolve()
strict = sys.argv[2] == "1"
targets = [pathlib.Path(p) for p in sys.argv[3:]]

issues = 0
docs_scanned = 0

required_fields = [
    "source owner",
    "projection role",
    "freshness trigger",
    "drift check",
    "absorption owner",
]

section_heading_re = re.compile(r"^## (Index Projection Maintenance|Projection Maintenance)\s*$", re.M)
next_h2_re = re.compile(r"^## ", re.M)
projection_card_re = re.compile(
    r"(Index projection card|Control-plane projection index|maintained projection-index card|Index Projection Cards)",
    re.I,
)
constitution_marker_re = re.compile(r"Control Planes As Maintained Indexes")


def report(msg: str) -> None:
    global issues
    print(msg)
    issues += 1


def normalize_target(raw: pathlib.Path):
    rel = raw
    if rel.is_absolute():
        try:
            rel = rel.relative_to(root)
        except ValueError:
            return None
    return rel


def add_if_file(paths: set[pathlib.Path], rel: str) -> None:
    path = root / rel
    if path.is_file():
        paths.add(path.relative_to(root))


def add_glob(paths: set[pathlib.Path], pattern: str) -> None:
    for path in root.glob(pattern):
        if path.is_file():
            rel = path.relative_to(root)
            if should_exclude(rel):
                continue
            paths.add(rel)


def should_exclude(rel: pathlib.Path) -> bool:
    parts = rel.parts
    joined = rel.as_posix()
    if "/runs/" in joined or "/campaigns/" in joined:
        return True
    if "/prompts/" in joined or rel.name.startswith("prompt"):
        return True
    if len(parts) >= 2 and parts[0] == "experiments" and parts[1] == "archive":
        return rel != pathlib.Path("experiments/archive/README.md")
    if len(parts) >= 2 and parts[0] == "Paper" and parts[1] == "archive":
        return True
    return False


def candidate_paths() -> set[pathlib.Path]:
    paths: set[pathlib.Path] = set()
    for rel in [
        "AGENTS.md",
        "CLOSURE.md",
        "DEVELOPMENT.md",
        "EVALUATION.md",
        "EVOLUTION.md",
        "development/README.md",
        "development/status/README.md",
        "development/open-source/README.md",
        "development/contracts/README.md",
        "development/contracts/control-plane-causal-path.md",
        "evolution/README.md",
        "evolution/patterns.md",
        "evolution/arc-index.md",
        "experiments/README.md",
        "experiments/PROTOCOL.md",
        "experiments/protocol/README.md",
        "experiments/archive/README.md",
        "experiments/incubator/README.md",
        # (active experiment-surface entries redacted in the public copy)
        "Paper/README.md",
        "Paper/program-map.md",
        "Paper/research-map.md",
        "Paper/evidence-map.md",
        "Paper/TRACK/README.md",
        "Paper/TRACK/PORTFOLIO.md",
        "Paper/venues/README.md",
        "Paper/archive/track-system-meta/README.md",
        "Paper/archive/track-alife-presentations/README.md",
        "Paper/venues/objects/alife-2026/references/README.md",
        "Paper/essays/README.md",
        "Paper/essays/STATUS.md",
        "Paper/essays/PORTFOLIO.md",
        "Paper/essays/_meta/README.md",
        "Paper/themes/README.md",
        "Paper/core/README.md",
        "Paper/core/registries/README.md",
        "Paper/core/registries/runtime-surface.md",
        "Paper/core/registries/env-controls.md",
        "Paper/core/evolution/19-runtime-physics-map.md",
        "Paper/themes/repository-self-model/plans/control-plane-causal-map.md",
        "tests/README.md",
        "tests/runtime/README.md",
        "tests/runtime/COVERAGE_MAP.md",
        "tests/model/README.md",
        "tests/model/l5-necessity/README.md",
        "tests/model/pilot-l5-necessity/process-surface/peer-callback-cleanroom/README.md",
        "scripts/README.md",
        ".githooks/README.md",
        "obsidian/README.md",
        "obsidian/bases/README.md",
        "operator/qctl/README.md",
        "internal/world/README.md",
    ]:
        add_if_file(paths, rel)

    for pattern in [
        # (active experiment-surface entry redacted in the public copy)
        "Paper/TRACK/venues/README.md",
        "Paper/TRACK/audiences/README.md",
        "Paper/TRACK/routes/README.md",
        "Paper/papers/quine-notes-series/plans/README.md",
        "Paper/core/*/README.md",
        "Paper/themes/*/README.md",
        "Paper/themes/*/papers/README.md",
        "Paper/themes/*/plans/README.md",
        "Paper/themes/*/references/README.md",
        "Paper/themes/*/analysis/README.md",
        "Paper/themes/*/adjacent/README.md",
    ]:
        add_glob(paths, pattern)
    return paths


def projection_section(text: str) -> str | None:
    match = section_heading_re.search(text)
    if not match:
        return None
    start = match.start()
    next_match = next_h2_re.search(text, match.end())
    end = next_match.start() if next_match else len(text)
    return text[start:end]


def has_required_fields(text: str) -> bool:
    lower = text.lower()
    return all(field in lower for field in required_fields)


def check_path(rel: pathlib.Path) -> None:
    global docs_scanned
    path = root / rel
    if not path.is_file():
        return
    docs_scanned += 1
    text = path.read_text()

    section = projection_section(text)
    if section is not None:
        if not has_required_fields(section):
            missing = [field for field in required_fields if field not in section.lower()]
            report(f"INCOMPLETE projection maintenance: {rel} :: missing {', '.join(missing)}")
        return

    if projection_card_re.search(text):
        if not has_required_fields(text):
            missing = [field for field in required_fields if field not in text.lower()]
            report(f"INCOMPLETE projection-index card: {rel} :: missing {', '.join(missing)}")
        return

    if rel == pathlib.Path("AGENTS.md") and constitution_marker_re.search(text):
        return

    report(f"MISSING projection maintenance: {rel}")


all_candidates = candidate_paths()
if targets:
    selected = set()
    for raw in targets:
        rel = normalize_target(raw)
        if rel is None:
            continue
        if should_exclude(rel):
            continue
        if rel in all_candidates:
            selected.add(rel)
else:
    selected = all_candidates

for rel in sorted(selected):
    check_path(rel)

print(f"Control-plane projection-maintenance audit: docs_scanned={docs_scanned} issues={issues}")
if strict and issues:
    sys.exit(1)
PY
