#!/usr/bin/env bash
#
# VALIDATOR INTENTION
# purpose:
#   - require canonical control-plane surfaces to expose stable YAML frontmatter for query-oriented tooling
#   - keep paper dossiers and opted-in experiment surfaces mechanically queryable for Obsidian/Bases-like overlays
#   - support both repo-wide audit mode and targeted file validation for hooks
# repair-guidance:
#   - preserve "report drift, do not rewrite files" behavior
#   - preserve targeted-file mode because repo hooks may rely on it
#   - if frontmatter schemas change, update this script, Paper/README.md, experiments/PROTOCOL.md, and .githooks/README.md together

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
sys.path.insert(0, str(root / "scripts" / "lib"))

from paper_dossier_labels import list_paper_dossiers  # noqa: E402

issues = 0
docs_scanned = 0
frontmatter_re = re.compile(r"\A---\n(.*?)\n---\n", re.S)
kv_re = re.compile(r"^([A-Za-z0-9_-]+):\s*(.+?)\s*$")
allowed_state_modes = {
    "full-manuscript-home",
    "planning-skeleton",
    "delta-reservoir",
    "published-single-manuscript",
    "workshop-submission-home",
    "navigation-stub",
}


def report(msg: str) -> None:
    global issues
    print(msg)
    issues += 1


def parse_frontmatter(path: pathlib.Path):
    text = path.read_text()
    m = frontmatter_re.match(text)
    if not m:
        return None
    data = {}
    for line in m.group(1).splitlines():
        if not line.strip():
            continue
        km = kv_re.match(line)
        if km:
            data[km.group(1)] = km.group(2)
    return data


def resolve_repo_relative(base_path: pathlib.Path, rel: str):
    try:
        return (base_path.parent / rel).resolve().relative_to(root.resolve())
    except (ValueError, OSError):
        return None


known_paper_refs = set()
for dossier_rel in list_paper_dossiers(root):
    data = parse_frontmatter(dossier_rel)
    if not data:
        continue
    track = data.get("track")
    paper_id = data.get("paper_id")
    if track and paper_id:
        known_paper_refs.add(f"{track}/{paper_id}")


def scan_paper_dossier(path: pathlib.Path):
    global docs_scanned
    docs_scanned += 1
    data = parse_frontmatter(path)
    if data is None:
        report(f"MISSING frontmatter: {path}")
        return
    claim_root_required = [
        "paper_id",
        "track",
        "cluster",
        "family",
        "source_home",
        "plan_home",
    ]
    single_config_required = claim_root_required + [
        "status",
        "state_mode",
        "role",
        "publication_class",
        "maturity_target",
        "claim_bundle",
        "venue",
        "format-envelope",
    ]
    config_child_required = [
        "configuration_id",
        "claim_bundle",
        "venue",
        "format-envelope",
        "role",
        "maturity-target",
        "status",
        "state-mode",
    ]
    legacy_required = claim_root_required + [
        "status",
        "state_mode",
        "role",
        "publication_class",
        "maturity_target",
    ]
    paper_family_home = (
        len(path.parts) >= 4
        and path.parts[0] == "Paper"
        and path.parts[1] == "papers"
    )
    config_child_home = paper_family_home and "configuration_id" in data
    if config_child_home:
        required = config_child_required
    elif paper_family_home and "configurations" in data:
        required = claim_root_required
    elif paper_family_home:
        required = single_config_required
    else:
        required = legacy_required
    missing = [k for k in required if k not in data]
    if missing:
        report(f"INCOMPLETE paper frontmatter: {path} :: {' '.join(missing)}")
    if paper_family_home and "target_venue" in data:
        report(f"OBSOLETE paper frontmatter: {path} :: target_venue")
    if config_child_home:
        expected_config_id = path.parent.name
        if data.get("configuration_id") != expected_config_id:
            report(
                f"MISMATCH configuration_id: {path} :: "
                f"expected={expected_config_id} got={data.get('configuration_id')}"
            )
    if paper_family_home and "configurations" in data:
        for key in ("status", "state_mode", "publication_class", "maturity_target"):
            if key in data:
                report(f"INVALID multi-config root field: {path} :: {key}")
    state_mode = data.get("state_mode") or data.get("state-mode")
    if state_mode and state_mode not in allowed_state_modes:
        report(f"INVALID state_mode: {path} :: {state_mode}")
    if state_mode == "navigation-stub" and data.get("publication_class") != "redirect":
        report(f"INVALID navigation-stub class: {path}")
    if data.get("publication_class") == "redirect" and state_mode != "navigation-stub":
        report(f"INVALID redirect state_mode: {path}")
    expected_track = "blogs" if len(path.parts) >= 3 and path.parts[0:2] == ("Paper", "essays") else None
    if expected_track and data.get("track") != expected_track:
        report(f"MISMATCH paper track: {path} :: expected={expected_track} got={data.get('track')}")
    expected_paper_id = path.parent.name
    theme_role_home = (
        len(path.parts) >= 6
        and path.parts[0] == "Paper"
        and path.parts[1] == "themes"
        and path.parts[3] == "papers"
    )
    if not paper_family_home and not theme_role_home and data.get("paper_id") != expected_paper_id:
        report(f"MISMATCH paper_id: {path} :: expected={expected_paper_id} got={data.get('paper_id')}")
    theme = data.get("theme")
    theme_home = data.get("theme_home")
    if bool(theme) != bool(theme_home):
        report(f"INCOMPLETE theme metadata: {path} :: theme and theme_home must appear together")
    if theme and theme_home:
        resolved_theme_home = resolve_repo_relative(path, theme_home)
        if resolved_theme_home is None:
            report(f"INVALID theme_home path: {path} :: {theme_home}")
        else:
            abs_theme_home = root / resolved_theme_home
            if not abs_theme_home.is_file():
                report(f"MISSING theme_home target: {path} :: {resolved_theme_home}")
            if len(resolved_theme_home.parts) < 4 or resolved_theme_home.parts[:3] != ("Paper", "theory", "views"):
                report(f"INVALID theme_home location: {path} :: {resolved_theme_home}")
            elif resolved_theme_home.suffix != ".md":
                report(f"INVALID theme_home target: {path} :: {resolved_theme_home}")
    inherits_from = data.get("inherits_from")
    if inherits_from and inherits_from not in known_paper_refs:
        report(f"UNKNOWN inherits_from target: {path} :: {inherits_from}")
    if state_mode == "delta-reservoir" and not inherits_from:
        report(f"MISSING inherits_from for delta-reservoir: {path}")


def scan_experiment_surface(path: pathlib.Path):
    global docs_scanned
    data = parse_frontmatter(path)
    if data is None:
        return
    if "surface_kind" not in data:
        return
    docs_scanned += 1
    expected_id = "/".join(path.parent.parts[1:])
    if len(path.parts) >= 3 and path.parts[1] == "active":
        # migrated family home (experiments restructure): a family/<id-slug>
        # README under the active experiment tree is validated by family (not
        # phase); `phase`/`lineage_phase` are retained only as lineage.
        required = ["surface_kind", "family", "experiment_id", "status"]
        missing = [k for k in required if k not in data]
        if missing:
            report(f"INCOMPLETE experiment frontmatter: {path} :: {' '.join(missing)}")
        expected_family = path.parts[2]
        if data.get("family") != expected_family:
            report(f"MISMATCH experiment family: {path} :: "
                   f"expected={expected_family} got={data.get('family')}")
        if data.get("experiment_id") != expected_id:
            report(f"MISMATCH experiment_id: {path} :: "
                   f"expected={expected_id} got={data.get('experiment_id')}")
        return
    required = ["surface_kind", "phase", "experiment_id", "status"]
    missing = [k for k in required if k not in data]
    if missing:
        report(f"INCOMPLETE experiment frontmatter: {path} :: {' '.join(missing)}")
    expected_phase = path.parts[1]
    if data.get("phase") != expected_phase:
        report(f"MISMATCH experiment phase: {path} :: expected={expected_phase} got={data.get('phase')}")
    if data.get("experiment_id") != expected_id:
        report(f"MISMATCH experiment_id: {path} :: expected={expected_id} got={data.get('experiment_id')}")


candidate_paths = set()
if targets:
    for path in list_paper_dossiers(root, targets):
        candidate_paths.add(path)
    for raw in targets:
        rel = raw
        if rel.is_absolute():
            try:
                rel = rel.relative_to(root)
            except ValueError:
                continue
        if rel.name == "README.md":
            candidate_paths.add(rel)
else:
    for path in list_paper_dossiers(root):
        candidate_paths.add(path)
    for path in root.glob("Paper/papers/*/*/README.md"):
        rel_path = path.relative_to(root)
        data = parse_frontmatter(rel_path)
        if data and "configuration_id" in data:
            candidate_paths.add(rel_path)
    for path in root.glob("experiments/p*/**/README.md"):
        if any(part in {"runs", "archive", "analysis", "habitat"} for part in path.parts):
            continue
        candidate_paths.add(path.relative_to(root))

for rel_path in sorted(candidate_paths):
    parts = rel_path.parts
    if len(parts) >= 4 and parts[0] == "Paper":
        scan_paper_dossier(rel_path)
        continue
    if len(parts) >= 3 and parts[0] == "experiments" and parts[1].startswith("p"):
        scan_experiment_surface(rel_path)

print(f"Control-plane-frontmatter audit: docs_scanned={docs_scanned} issues={issues}")
if strict and issues:
    sys.exit(1)
PY
