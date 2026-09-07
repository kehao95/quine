# Local Path Exposure Inventory, 2026-05-24

This inventory is for the Quine full open-source release gate.

It records pattern counts and file paths only. It does not preserve matched line
contents, secret values, or machine-specific host configuration details beyond
the path patterns already being audited.

## Scan Scope

- current git-visible files from `git ls-files --cached --others --exclude-standard`
- fixed-string scans for user-home paths and known machine/service paths
- no history rewrite, no file edits, and no payload deletion

## Pattern Summary

| Pattern | Files | Occurrences | Release read |
|---------|-------|-------------|--------------|
| `/Users/kehao95` | 1088 | 9077 | mostly retained model tapes/logs; not public as raw git material |
| `/home/kehao95` | 33 | 183 | newer Linux run logs, docs, and metadata; normalize or classify |
| `/srv/storage` | 4 | 5 | artifact host/storage topology; review before public release |
| `/etc/nginx` | 2 | 5 | host service config/log path; remove or keep only as non-clickable historical evidence |
| `/etc/cloudflared` | 0 | 0 | no current git-visible hits |
| `/var/lib/cloudflared` | 0 | 0 | no current git-visible hits |
| `/home/` | 213 | 2472 | generic Linux absolute paths, including non-owner examples and run snapshots |
| `/Users/` | 1090 | 9085 | generic macOS absolute paths, mostly owner-home occurrences |

The most important result is shape, not individual line content: the dominant
exposure is raw model/run evidence, especially `tests/model/**/runs/**` tapes
and logs.

## Follow-Up Normalization

On 2026-05-24, active non-run docs/scripts outside this open-source register
were normalized from owner-machine absolute paths to repo-relative paths,
placeholders, or script-relative paths. The tracked DVC config was also reduced
to the public artifact remote; private storage remotes now belong in the
ignored `.dvc/config.local`.

The follow-up scan excluded `runs/**`, this open-source register, DVC cache, and
raw log/tape files. It found no remaining matches for:

```text
/Users/kehao95
/home/kehao95
/srv/storage
/etc/nginx
/etc/cloudflared
/var/lib/cloudflared
```

This does not clear raw runs, archived logs, or DVC payloads. Those remain under
the raw-run/DVC scan gate.

## Surface Concentration

| Surface | Files with path hits | Release read |
|---------|----------------------|--------------|
| `tests/model` | 1128 | raw evaluation evidence; needs DVC/sanitization policy before public release |
| `experiments/active/benchmarks` | 114 | benchmark/terminal-bench evidence; review before public release |
| `Paper/themes/autopoiesis-renewal` | 26 | theory/analysis docs; normalize active prose, preserve historical evidence carefully |
| `experiments/archive` | 10 | historical run snapshots; likely DVC/archive rather than git-visible public tree |
| `Paper/themes/runtime-ontology` | 6 | reference metadata and reports; review external/source provenance |
| `experiments/active/null-ecology` | 6 | experiment evidence; classify with run policy |
| `experiments/active/digital-red-queen` | 3 | host lineage scripts/docs; normalize machine paths |
| `experiments/active/autopoiesis-succession` | 2 | current experiment material; scan DVC payloads separately |
| `internal` | 2 | code references or tests; normalize if not intentional examples |

## Top Exposure Shapes

High-count `/Users/kehao95` hits are concentrated in:

```text
tests/model/l3-usage/exec/explicit-self-source-rebuild-handoff/runs/*/tape.jsonl
tests/model/l3-usage/exec/explicit-self-source-rebuild-handoff/runs/*/tape.log.yaml
tests/model/l3-usage/skills/explicit-frontmatter-refresh/runs/*/tape.jsonl
tests/model/l3-usage/skills/explicit-frontmatter-refresh/runs/*/tape.log.yaml
tests/model/l4-discovery/context-memory/inheritance-discovery/runs/*/tape.jsonl
tests/model/l4-discovery/context-memory/inheritance-discovery/runs/*/tape.log.yaml
```

High-count `/home/` hits include:

```text
tests/model/l3-usage/exec/explicit-self-source-rebuild-handoff/runs/20260419-184540-k2.5/workspace/rebuilt-quine
tests/model/l3-usage/exec/explicit-self-source-rebuild-handoff/runs/*/tape.jsonl
tests/model/l3-usage/exec/explicit-self-source-rebuild-handoff/runs/*/tape.log.yaml
tests/model/l3-usage/relation-recovery/explicit-outcome-semantics/runs/*/tape.jsonl
tests/model/l4-discovery/relation-recovery/status-discovery/runs/*/tape.log.yaml
```

Known machine/service path hits include:

```text
.dvc/config
experiments/protocol/dvc-externalization-rollout.md
experiments/active/benchmarks/bench-04-programbench/evidence/README.md
experiments/archive/ch5-emergent-behavior/2.1-babel/logs/run1.jsonl
```

The first three active/config hits above have been normalized after the baseline
scan. The archived log remains under the raw-run/DVC gate.

Some open-source risk-register docs intentionally mention audited path patterns.
Those self-referential hits are not release blockers by themselves, but they
should remain pattern-level, not operational host details.

## Required Handling

- Do not mechanically replace paths inside raw tapes/logs and call that public
  safe. Those files also need prompt, tool-output, and secret review.
- Treat `tests/model/**/runs/**`, archived experiment logs, and retained
  runtime snapshots as evidence payloads: summaries and pointers can live in
  git, raw evidence should go through DVC/public-payload review.
- Normalize active docs and scripts to placeholders such as `/path/to/workspace`
  or `$QUINE_WORKSPACE_ROOT` when the path is instructional rather than
  historical evidence.
- Keep private DVC storage-host topology in ignored local config, not tracked
  repo config or public docs.
- Keep host service paths such as nginx/cloudflared out of the public source
  tree unless they are intentionally documented as non-operational examples.

## Safe Next Actions

1. Build a candidate path-normalization patch for active docs/scripts only.
2. Keep run/tape path exposure under the raw-run/DVC gate instead of editing
   every log in place.
3. Scan the DVC payloads selected for public restore separately from git
   pointer files.
4. Re-run this inventory after provider-profile cleanup and DVC selection.
