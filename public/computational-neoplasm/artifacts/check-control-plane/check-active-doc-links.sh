#!/usr/bin/env bash
#
# VALIDATOR INTENTION
# purpose:
#   - protect active control-plane markdown docs from non-portable local links
#   - prevent active docs from advertising broken clickable evidence
#   - require run-like evidence linked from active docs to be repo-retained rather than merely present locally
#   - ignore archived/generated evidence snapshots where historical paths are the artifact
#   - support both repo-wide audit mode and targeted file validation for hooks
# repair-guidance:
#   - preserve "report drift, do not rewrite files" behavior
#   - preserve targeted-file mode because repo hooks may rely on it
#   - if the active-control-plane boundary changes, update this script, DEVELOPMENT.md, experiments/PROTOCOL.md, and .githooks/README.md together

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

if [[ "${#TARGET_FILES[@]}" -eq 0 ]]; then
  while IFS= read -r path; do
    [[ -n "$path" ]] || continue
    TARGET_FILES+=("$path")
  done < <(git ls-files '*.md')
fi

python3 - "$ROOT_DIR" "$STRICT" "${TARGET_FILES[@]}" <<'PY'
import os
import pathlib
import re
import subprocess
import sys

root = pathlib.Path(sys.argv[1]).resolve()
strict = sys.argv[2] == "1"
targets = [pathlib.Path(p) for p in sys.argv[3:]]

excluded_segments = {
    "runs",
    "archive",
    "artifacts",
    "jobs",
    "workspaces",
    "sections",
    "analysis",
    "external-sources",
    "templates",
}
markdown_link_re = re.compile(r"\[([^\]]*)\]\(([^)]+)\)")
issues = 0
advisory_issues = 0
docs_scanned = 0
tracked_files = {
    line
    for line in subprocess.run(
        ["git", "ls-files", "--"],
        capture_output=True,
        check=True,
        text=True,
    ).stdout.splitlines()
    if line
}
tracked_prefixes = set()
for tracked_path in tracked_files:
    path = pathlib.PurePosixPath(tracked_path)
    for parent in path.parents:
        parent_str = parent.as_posix()
        if parent_str != ".":
            tracked_prefixes.add(parent_str)


def is_active_doc(path: pathlib.Path) -> bool:
    if path.suffix != ".md":
        return False
    if any(part in excluded_segments for part in path.parts):
        return False
    return path.exists()


def report(msg: str) -> None:
    global issues
    print(msg)
    issues += 1


def advisory(msg: str) -> None:
    global advisory_issues
    print(f"ADVISORY {msg}")
    advisory_issues += 1


def github_slug(heading: str) -> str:
    """Convert a heading string to its GitHub-style anchor slug."""
    slug = heading.lower()
    slug = re.sub(r"[^\w\s-]", "", slug)
    slug = re.sub(r"[\s]+", "-", slug.strip())
    return slug


def extract_headings(md_path: pathlib.Path):
    """Return a set of GitHub-style slugs for all headings in a markdown file."""
    heading_re = re.compile(r"^#{1,6}\s+(.+?)\s*$")
    fence_re = re.compile(r"^(```|~~~)")
    slugs = set()
    in_fence = False
    try:
        text = md_path.read_text(encoding="utf-8", errors="replace")
    except OSError:
        return slugs
    for line in text.splitlines():
        if fence_re.match(line):
            in_fence = not in_fence
            continue
        if in_fence:
            continue
        m = heading_re.match(line)
        if m:
            slugs.add(github_slug(m.group(1)))
    return slugs


def git_tracked(rel_path: str) -> bool:
    normalized = rel_path.rstrip("/")
    if normalized in tracked_files or normalized in tracked_prefixes:
        return True
    candidates = {f"{normalized}.dvc"}
    parts = pathlib.PurePosixPath(normalized).parts
    for marker in ("runs", "batches", "campaigns", "habitat"):
        for idx, part in enumerate(parts):
            if part != marker or idx + 1 >= len(parts):
                continue
            bundle = pathlib.PurePosixPath(*parts[: idx + 2]).as_posix()
            candidates.add(f"{bundle}.dvc")
            candidates.add(f"{bundle}.tar.zst.dvc")
    return any(candidate in tracked_files for candidate in candidates)


for rel_path in targets:
    rel_path = pathlib.Path(rel_path)
    if not is_active_doc(rel_path):
        continue

    docs_scanned += 1
    text = rel_path.read_text()
    for lineno, line in enumerate(text.splitlines(), 1):
        for _, raw_target in markdown_link_re.findall(line):
            if raw_target.startswith(("http://", "https://", "mailto:")):
                continue

            # Pure same-page anchor: skip broken check, but do advisory slug check.
            if raw_target.startswith("#"):
                frag = raw_target[1:]
                if frag:  # non-empty fragment
                    slugs = extract_headings(rel_path)
                    if slugs and frag not in slugs:
                        advisory(
                            f"anchor-not-found (same-page): {rel_path}:{lineno} :: {raw_target}"
                        )
                continue

            # Split file path from anchor fragment.
            if "#" in raw_target:
                file_part, frag = raw_target.split("#", 1)
            else:
                file_part, frag = raw_target, ""

            target_path = pathlib.Path(file_part)
            if target_path.is_absolute():
                report(f"NONPORTABLE absolute link: {rel_path}:{lineno} :: {raw_target}")
                if target_path.exists():
                    if frag:
                        slugs = extract_headings(target_path)
                        if slugs and frag not in slugs:
                            advisory(
                                f"anchor-not-found: {rel_path}:{lineno} :: {raw_target}"
                            )
                    continue
                report(f"BROKEN markdown link: {rel_path}:{lineno} :: {raw_target}")
                continue

            resolved = (rel_path.parent / target_path).resolve()
            if not resolved.exists():
                try:
                    rel_target = resolved.relative_to(root).as_posix()
                except ValueError:
                    rel_target = ""
                if rel_target and git_tracked(rel_target):
                    continue
                report(f"BROKEN markdown link: {rel_path}:{lineno} :: {raw_target}")
                continue

            # Advisory anchor check for valid paths.
            if frag:
                if resolved.suffix == ".md":
                    slugs = extract_headings(resolved)
                    if slugs and frag not in slugs:
                        advisory(
                            f"anchor-not-found: {rel_path}:{lineno} :: {raw_target}"
                        )

            try:
                rel_target = resolved.relative_to(root).as_posix()
            except ValueError:
                continue

            if "/runs/" in f"/{rel_target}" or "/batches/" in f"/{rel_target}":
                if not git_tracked(rel_target):
                    report(
                        f"UNRETAINED evidence link: {rel_path}:{lineno} :: {raw_target}"
                    )

print(
    f"Active-doc-link audit: docs_scanned={docs_scanned} issues={issues}"
    + (f" advisory={advisory_issues}" if advisory_issues else "")
)
if strict and issues:
    sys.exit(1)
PY
