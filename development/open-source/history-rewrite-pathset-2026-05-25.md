# History Rewrite Pathset, 2026-05-25

This note turns the current open-source release blocker discussion into a
concrete first-pass public-history rewrite boundary.

It does **not** execute the rewrite. It defines the path families that should
be removed from the first public branch lineage once the current `lab` HEAD is
clean enough to serve as the post-rewrite tip.

## Purpose

The public-history rewrite now has two separate jobs:

1. shrink repo history by removing bulky generated payloads that do not need to
   stay in public git
2. remove secret-bearing, account-bearing, path-bearing, and raw-run history
   that should not survive into the public branch at all

This is intentionally separate from selective public DVC publication:

- private DVC keeps raw historical payloads when we still want them
- public git history keeps the control plane and sanitized source tree
- public DVC only gets exact allowlisted artifacts after independent review

## Current Inputs

Current branch at planning time: `lab`  
Current commit at planning time: `642fe0f46`

Current tracked HEAD blockers:

| Surface | Count | Why it matters now |
|---------|-------|--------------------|
| `tests/model/**/runs/**` | 0 tracked raw non-pointer files after 2026-05-25 follow-up; 369 per-run private `.dvc` pointers remain staged | current raw model-run files moved to private DVC pointers; history remains a rewrite target |
| `experiments/active/benchmarks/bench-01-terminal-bench/analysis/sweep-logs/**` | 0 tracked raw files after 2026-05-25 follow-up | current raw files moved to private DVC pointer; history remains a rewrite target |
| retained env/profile files | 12 tracked files after 2026-05-25 follow-up | `.env.gpt-5.4-codex-pool` was removed from HEAD |
| `experiments/active/benchmarks/bench-01-terminal-bench/tasks/**/runs/*.tar.zst.dvc` | 256 tracked pointers | raw Terminal-Bench archives are `deny-raw` |

History-shape inputs:

| Family | Count | Read |
|--------|-------|------|
| history tar/tar.zst paths | 527 distinct paths | raw archive history exists and should not survive public rewrite |
| history run/tape/session/log/workspace paths | 93712 distinct paths | generated run history is the dominant public-history risk family |
| history env/profile/settings paths | 91 distinct paths | old provider/local settings history still exists despite current HEAD cleanup |

Tooling status:

- `git filter-repo` is installed at `/usr/bin/git-filter-repo`
- `--paths-from-file` supports `literal:`, `glob:`, and `regex:` directives

## Canonical First-Pass Drop Set

Canonical machine-readable pathset:

- [history-rewrite-drop-paths.txt](./history-rewrite-drop-paths.txt)

Recommended command shape, to run in a fresh mirror or disposable clone rather
than this working repo:

```bash
git filter-repo --invert-paths \
  --paths-from-file development/open-source/history-rewrite-drop-paths.txt
```

The first-pass drop set intentionally targets high-confidence risk families:

1. `.env*`, nested `.env*`, and `profiles/*.env`
2. `.claude/**`, `.opencode/**`, and local settings copies
3. generated `runs/**`, `runtime/jobs/**`, `runtime/tapes/**`, and
   `workspace/**` payloads
4. `analysis/sweep-logs/**`
5. `meta/summary.env`, `meta/timeout-summary.env`,
   `meta/capability-surface.env`, and `meta/env-after.txt`
6. raw `*.tar` and `*.tar.zst`
7. ProgramBench campaign `.dvc` pointers that are still in `sanitize` / `deny`
   territory

## Why The Pathset Is Broad

The point of the first rewrite is not to preserve every private historical
moment in public form.

The public branch should preserve:

- source
- READMEs and protocol surfaces
- analysis notes
- small inventories
- selected safe `.dvc` pointers
- sanitized current-state artifacts that we explicitly choose to keep

It should **not** preserve:

- old provider profile history
- machine-local settings history
- raw run/tape/workspace history
- raw Terminal-Bench archive history
- bulky generated tarball history

If a safe public replacement matters, re-add it from the sanitized current tree
after the rewrite rather than keeping the unsafe private evolution.

## Current HEAD Preconditions Before Rewrite

The current tracked `lab` tip is no longer blocked by raw model-run files after
the 2026-05-25 `tests/model` DVC externalization. The public-history rewrite
still needs a fresh mirror or disposable clone because private history remains
unsafe even though the staged tip now scans clean.

Before claiming the public release boundary is externally true:

1. purge Cloudflare edge cache residue on the public artifact hostname
2. rerun secret/path scans on the committed tip and again on the rewritten
   public-history result

Recommended tip-first cleanup posture:

- keep `.dvc/public-manifest.txt` empty until public cache purge plus explicit
  first publication decision

Resolved follow-up decisions on 2026-05-25:

- `.env.gpt-5.4-codex-pool` was removed from HEAD
- Harbor sweep logs were moved from git-tracked raw files into
  `experiments/active/benchmarks/bench-01-terminal-bench/analysis/sweep-logs.dvc`
- the Harbor sweep-log DVC pointer is private-retention only and is not listed
  in `.dvc/public-manifest.txt`
- the archived path-bearing log under
  `experiments/archive/ch5-emergent-behavior/2.1-babel/logs/` is intentionally
  retained in HEAD per owner direction
- `tests/model/**/runs/**` was moved from raw tracked files into 369 per-run
  private DVC pointers
- the model-run DVC pointers were pushed to `private` and
  `lab-artifacts-gdrive-backup`
- `tests/model/**/runs/*.dvc` is private-retention only, remains `deny-raw`,
  and is not listed in `.dvc/public-manifest.txt`
- the staged-tree gitleaks scan after model-run externalization and residual
  false-positive cleanup found 0 findings

## What This Pathset Leaves For Later

This is a first pass, not the final size pass.

After the first rewrite, run a separate large-object review on a disposable
clone:

```bash
git filter-repo --analyze
```

Then decide whether a second pass is worth doing for non-sensitive but bulky
history such as:

- built binaries
- generated PDFs/PPTX previews
- oversized image exports
- any remaining historical blobs that are large but not risk-bearing

Keep this second pass separate so we do not mix leak prevention with purely
size-driven trimming.

## Execution Sequence

1. Push any raw historical payload that still needs preservation into private
   DVC.
2. Purge `dvc-artifacts.kehao.me` edge cache with a credential that can
   actually purge before claiming the public artifact store is externally
   empty.
3. Create a fresh mirror or disposable clone for rewrite work.
4. Run the first-pass `git filter-repo` command with
   [history-rewrite-drop-paths.txt](./history-rewrite-drop-paths.txt).
5. Re-apply the sanitized current tree if needed.
6. Re-run secret/path scans on the rewritten result.
7. Only then cut or force-update the eventual public branch.

## Current Decision

The recommended public-history stance is:

- aggressive about deleting unsafe generated history
- conservative about what gets reintroduced afterward
- keep raw history in private DVC when it still matters
- keep public git focused on sanitized source, control-plane state, and
  explicitly chosen evidence
