#!/usr/bin/env python3
"""
Shared parsing/graph library for the idea-layer edge graph readouts.

GENERATED-VIEW SUPPORT SCRIPT — not a canonical store, per read-layer-prior.
Reused by same_referent_detector.py and under_attention_frontier.py.

Parses every in-scope Paper/incubator/**/*.md file's `edges:` frontmatter
block via the same line-regex approach validated in the QA agent's
m1_tightening.py / full_m1.py (development/frontmatter-object-graph-feasibility
phase2/phase3 scratchpad scripts) — re-derived here (not copy-pasted) so this
directory is self-contained and reproducible without depending on any
/tmp scratchpad path.
"""
import re
import os
import glob
from pathlib import Path

REPO = str(Path(__file__).resolve().parents[3])

EDGE_LINE_RE = re.compile(
    r'-\s*\{\s*rel:\s*([a-zA-Z0-9\-]+)\s*,\s*to:\s*"?([^",}]+)"?\s*\}'
)


def in_scope_incubator_files():
    """Every Paper/incubator/**/*.md file, minus the 2 known index/README files
    that never carry edges: (Paper/incubator/README.md,
    Paper/incubator/divergent-directions/README.md)."""
    if not Path(REPO, "Paper/incubator").is_dir():
        raise FileNotFoundError("This historical readout requires the unreleased Paper/incubator tree")
    all_md = glob.glob(os.path.join(REPO, "Paper/incubator/**/*.md"), recursive=True)
    skip = {
        os.path.join(REPO, "Paper/incubator/README.md"),
        os.path.join(REPO, "Paper/incubator/divergent-directions/README.md"),
    }
    files = sorted(p for p in all_md if p not in skip)
    return [os.path.relpath(p, REPO) for p in files]


def parse_edges(relpath):
    """Return list of (rel, to) tuples found in relpath's edges: frontmatter block."""
    full = os.path.join(REPO, relpath)
    with open(full) as f:
        text = f.read()
    lines = text.split("\n")
    triples = []
    in_edges = False
    for line in lines:
        if line.strip() == "edges:":
            in_edges = True
            continue
        if in_edges:
            # end of the edges: list: a line that is a top-level frontmatter
            # key/terminator, i.e. not a "- {...}" list item and not blank.
            stripped = line.strip()
            if stripped == "---":
                break
            if stripped and not stripped.startswith("-"):
                # some other top-level key started (shouldn't happen — edges:
                # is documented as the last key — but guard anyway)
                break
            m = EDGE_LINE_RE.search(line)
            if m:
                rel, to = m.group(1), m.group(2).strip()
                triples.append((rel, to))
    return triples


def build_corpus_triples(files):
    """Return list of (src, rel, to) over the given file list."""
    all_triples = []
    for f in files:
        for rel, to in parse_edges(f):
            all_triples.append((f, rel, to))
    return all_triples


if __name__ == "__main__":
    import tempfile

    assert Path(REPO) == Path(__file__).resolve().parents[3]
    with tempfile.TemporaryDirectory() as scratch:
        REPO = scratch
        try:
            in_scope_incubator_files()
        except FileNotFoundError:
            pass
        else:
            raise AssertionError("missing private input must not look like an empty graph")
        Path(scratch, "Paper/incubator").mkdir(parents=True)
        assert in_scope_incubator_files() == []
    print("graph readout input-boundary self-check passed")
