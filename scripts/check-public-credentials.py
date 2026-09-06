#!/usr/bin/env python3
"""Fail when a projected tree or its manifest-addressed DVC data contains secrets."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import subprocess
import tempfile


ALLOWED = {
    # Literal environment-variable name in a retained source snapshot.
    "12ccdeebc28425d95b639cc46376b22c388bd618e4dc6302f2edf8790c055355",
    # Deliberately inert key in the hostile-script-survival fixture.
    "3f2cd8e57b096fe7e4a78a5627e34ca3f885ad65a56e61c287cf4211bbc5949f",
}
TREE_ALLOWED = {
    (
        "internal/llm/claudeoauth/claudeoauth.go",
        "577ac1adb016904259eff7b7cb12a01d97cd591d34d9ebd20780f580d9a107d7",
    )
}


def run_gitleaks(path: Path, report: Path) -> list[dict]:
    result = subprocess.run(
        [
            "gitleaks", "dir", str(path), "--no-banner", "--no-color",
            "--report-format", "json", "--report-path", str(report),
        ],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.PIPE,
        text=True,
    )
    if result.returncode not in (0, 1):
        raise SystemExit(f"gitleaks failed: {result.stderr.strip()}")
    return json.loads(report.read_text()) if report.exists() else []


def digest(finding: dict) -> str:
    return hashlib.sha256(finding["Secret"].encode()).hexdigest()


def cache_object(roots: list[Path], md5: str) -> Path:
    for root in roots:
        path = root / md5[:2] / md5[2:]
        if path.is_file():
            return path
    return roots[0] / md5[:2] / md5[2:]


def pointer(pointer_path: Path) -> tuple[str, str]:
    text = pointer_path.read_text()
    md5 = re.search(r"^\s*-?\s*md5:\s*([0-9a-f]{32})(\.dir)?\s*$", text, re.M)
    out = re.search(r"^\s*path:\s*(\S.*?)\s*$", text, re.M)
    if not md5 or not out:
        raise SystemExit(f"unsupported DVC pointer: {pointer_path}")
    return md5.group(1) + (md5.group(2) or ""), out.group(1)


def safe_relpath(value: str) -> PurePosixPath:
    path = PurePosixPath(value)
    if path.is_absolute() or ".." in path.parts:
        raise SystemExit(f"unsafe DVC relpath: {value}")
    return path


def stage_dvc(manifest: Path, roots: list[Path], corpus: Path) -> dict[str, list[tuple[str, str]]]:
    root = manifest.parent.parent
    refs: dict[str, list[tuple[str, str]]] = {}
    pointers = [
        line.split("#", 1)[0].strip()
        for line in manifest.read_text().splitlines()
        if line.split("#", 1)[0].strip()
    ]
    for pointer_name in dict.fromkeys(pointers):
        pointer_path = root / safe_relpath(pointer_name)
        if not pointer_path.is_file():
            raise SystemExit(f"manifest entry missing: {pointer_name}")
        md5, _ = pointer(pointer_path)
        if md5.endswith(".dir"):
            index = cache_object(roots, md5[:-4] + ".dir")
            if not index.is_file():
                raise SystemExit(f"DVC index missing from cache: {pointer_name}")
            entries = json.loads(index.read_text())
        else:
            entries = [{"md5": md5, "relpath": Path(pointer_name).stem}]
        for entry in entries:
            rel = str(safe_relpath(entry["relpath"]))
            source = cache_object(roots, entry["md5"])
            if not source.is_file():
                raise SystemExit(f"DVC object missing from cache: {pointer_name}:{rel}")
            suffix = "".join(PurePosixPath(rel).suffixes)[-24:]
            name = entry["md5"] + suffix
            target = corpus / name
            refs.setdefault(name, []).append((pointer_name, rel))
            if not target.exists():
                try:
                    os.link(source, target)
                except OSError:
                    shutil.copyfile(source, target)
    return refs


def allowed_dvc(finding: dict, refs: dict[str, list[tuple[str, str]]]) -> bool:
    sha = digest(finding)
    if sha not in ALLOWED:
        return False
    uses = refs.get(Path(finding["File"]).name, [])
    if sha.startswith("12ccde"):
        return bool(uses) and all(rel.endswith("internal/config/config.go") for _, rel in uses)
    return bool(uses) and all("hostile-script-survival" in ptr for ptr, _ in uses)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--tree", type=Path, required=True)
    parser.add_argument("--manifest", type=Path, required=True)
    args = parser.parse_args()
    tree = args.tree.resolve()
    manifest = args.manifest.resolve()
    cache = Path(subprocess.check_output(["dvc", "cache", "dir"], text=True).strip()).resolve()
    roots = [cache / "files" / "md5"]
    localstore = subprocess.run(
        ["dvc", "config", "remote.localstore.url"], capture_output=True, text=True
    )
    if localstore.returncode == 0 and localstore.stdout.strip():
        roots.append(Path(localstore.stdout.strip()).resolve() / "files" / "md5")
    failures: list[tuple[str, str, str]] = []
    scratch = Path.cwd() / ".tmp"
    scratch.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="public-credential-scan-", dir=scratch) as temp:
        temp = Path(temp)
        tree_report = temp / "tree.json"
        for finding in run_gitleaks(tree, tree_report):
            rel = str(Path(finding["File"]).resolve().relative_to(tree))
            sha = digest(finding)
            if (rel, sha) not in TREE_ALLOWED:
                failures.append((finding["RuleID"], sha[:12], rel))
        corpus = temp / "dvc"
        corpus.mkdir()
        refs = stage_dvc(manifest, roots, corpus)
        dvc_report = temp / "dvc.json"
        for finding in run_gitleaks(corpus, dvc_report):
            if not allowed_dvc(finding, refs):
                name = Path(finding["File"]).name
                where = refs.get(name, [("unknown", name)])[0]
                failures.append((finding["RuleID"], digest(finding)[:12], f"{where[0]}:{where[1]}"))
    if failures:
        for rule, sha, where in failures:
            print(f"FAIL: {rule} sha256={sha} at {where}")
        print(f"credential scan failed: {len(failures)} finding(s)")
        return 1
    print("credential scan passed: projected tree and manifest-addressed DVC data")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
