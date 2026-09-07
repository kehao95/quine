# HEAD Release Scan, 2026-05-25

This follow-up scan measures the current `lab` HEAD after the DVC split and
profile cleanup work.

The goal is narrower than the earlier broad worktree audit:

- scan the current git-tracked tree only, using `git archive HEAD`
- summarize redacted scanner output without preserving secret values
- identify what still blocks `lab` from being public-safe

Separate DVC payload review remains in:

- [DVC Payload Scan, 2026-05-24](./dvc-payload-scan-2026-05-24.md)
- [DVC Public Payload Allow/Deny, 2026-05-24](./dvc-public-payload-allowdeny-2026-05-24.md)

## Scope

- branch: `lab`
- commit at scan time: `60b0077ba`
- snapshot source: `git archive HEAD`
- public manifest effective entries: `0`

## Scanner Summary

### Full HEAD snapshot

| Metric | Value |
|--------|-------|
| files | 9030 |
| bytes | 125,187,158 |
| redacted gitleaks findings | 39 |
| files with findings | 24 |
| rules hit | `generic-api-key` only |

Release read:

- the current git-tracked HEAD is not public-safe yet
- most current findings are retained run/tape/workspace artifacts, not active
  runtime source files
- not every finding is a live secret; several look like placeholder or
  synthetic strings, but the release gate still fails until these surfaces are
  classified, sanitized, or removed from public history

### Retained profile subset

Reviewed current kept env/profile files:

```text
.env.example
.env.gpt-5.3-codex-spark
.env.gpt-5.3-codex-spark-high
.env.gpt-5.4-codex
.env.gpt-5.4-codex-high
.env.gpt-5.4-codex-pool
.env.gpt-5.5-codex-fast
.env.gpt-5.5-codex-high
.env.gpt-5.5-codex-medium
.env.gpt-5.5-codex-xhigh
.env.kimi-oauth
experiments/active/benchmarks/bench-05-terminal-bench2-gpt54/profiles/gpt-5.5-medium-benchmark-codex-ref.env
experiments/active/benchmarks/bench-05-terminal-bench2-gpt54/profiles/gpt-5.5-medium-benchmark-codex-strict.env
```

Results:

- focused redacted `gitleaks` scan: `1` finding
- finding path: `.env.gpt-5.4-codex-pool`
- finding rule: `generic-api-key`
- focused grep for non-empty `KEY` / `TOKEN` / `SECRET` / `PASSWORD`
  assignments: `0` hits

Release read:

- the retained profile set looks much better than the earlier baseline
- `.env.gpt-5.4-codex-pool` still needs explicit keep/template/remove review
  before public cut

## Where Current Findings Cluster

Top gitleaks-hit files:

| File | Finding count |
|------|---------------|
| `tests/model/l4-discovery/containment/hostile-script-survival/runs/20260321-063214-k2.5/tape.log.yaml` | 4 |
| `tests/model/l4-discovery/containment/hostile-script-survival/runs/20260418-192328-5.4/tape.log.yaml` | 4 |
| `tests/model/l4-discovery/containment/hostile-script-survival/runs/20260311-174427-k2.5/tape.log.yaml` | 3 |
| `tests/model/l4-discovery/containment/hostile-script-survival/runs/20260418-192328-5.4/tape.jsonl` | 3 |
| `tests/model/l4-discovery/containment/hostile-script-survival/runs/20260519-034909-5.4/tape.log.yaml` | 3 |
| `.env.gpt-5.4-codex-pool` | 1 |
| `internal/config/config.go` | 1 |
| `tests/model/run.sh` | 1 |

Other hit files are additional `tests/model/**/runs/**` tapes, logs, and
workspace fixtures from the same containment/evaluation family.

## Path Exposure Summary

Tracked fixed-string path hits:

| Pattern | Tracked files hit | Main concentration |
|---------|-------------------|--------------------|
| `/Users/kehao95` | 1057 | `tests/model` (939 files), `experiments/active/benchmarks` (112 files) |
| `/home/kehao95` | 27 | `tests/model` (24 files) |
| `/srv/storage` | 3 | open-source register docs only |
| `/etc/nginx` | 4 | open-source register docs plus one archived experiment log |

Representative current blockers:

- `tests/model/**/runs/**`
- `experiments/active/benchmarks/bench-01-terminal-bench/analysis/sweep-logs/**`
- one archived experiment log under
  `experiments/archive/ch5-emergent-behavior/2.1-babel/logs/`

Release read:

- active prose/scripts are no longer the main local-path problem
- retained run records and benchmark logs are the main git-visible path
  exposure left in current HEAD

## Public DVC Surface State

- tracked manifest file: [`.dvc/public-manifest.txt`](../../.dvc/public-manifest.txt)
- effective published pointer entries at scan time: `0`
- local/public origin for old object paths: `404`
- Cloudflare edge status: still serves at least one previously cached object

Release read:

- the origin split worked
- public emptiness is still not fully trustworthy until edge cache is purged

## Current Blocking Read

`lab` is closer, but not releasable yet.

At scan time, the current HEAD blockers were:

1. retained tracked run/tape/workspace surfaces in `tests/model/**/runs/**`
2. tracked Harbor sweep logs with owner-path exposure
3. one retained profile review flag in `.env.gpt-5.4-codex-pool`
4. Cloudflare edge cache residue on the public artifact hostname

Follow-up on 2026-05-25:

- `.env.gpt-5.4-codex-pool` was removed from HEAD
- Harbor sweep logs were moved from git-tracked raw log files into
  `experiments/active/benchmarks/bench-01-terminal-bench/analysis/sweep-logs.dvc`
- that DVC pointer is a private-retention pointer and is not listed in
  `.dvc/public-manifest.txt`
- the retained Harbor sweep-log DVC object was pushed to `private` and
  `lab-artifacts-gdrive-backup`
- a staged-tree rescan after those changes found 38 redacted `gitleaks`
  findings across 23 tracked files, all still under `generic-api-key`
- the staged tree has `/Users/kehao95` in 946 tracked files and
  `/home/kehao95` in 28 tracked files
- the archived path-bearing log was intentionally left in HEAD per owner
  direction
- `tests/model/**/runs/**` raw run directories were externalized into 369
  per-run `.dvc` pointers for private retention; the staged index now has 0
  raw non-pointer files under `tests/model/**/runs/**`
- the 369 model-run DVC pointers are in sync with both `private` and
  `lab-artifacts-gdrive-backup`
- `tests/model/**/runs/*.dvc` remains private-only/`deny-raw` and is not in
  `.dvc/public-manifest.txt`
- two residual source/test scanner false positives were reshaped without
  changing the runtime contract: `internal/config/config.go` now builds
  sensitive env pairs via a helper, and the hostile-script fixture uses a
  non-secret-looking demo payload label
- the follow-up staged-tree gitleaks scan found 0 findings across 5,945 files
  and 71,787,225 bytes
- follow-up staged fixed-string path counts are `/Users/kehao95`: 7 files,
  `/home/kehao95`: 4 files, `/srv/storage`: 5 files, `/etc/nginx`: 5 files,
  with 0 hits for all four under `tests/model`

## Next Repairs

1. obtain a Cloudflare credential or dashboard path that can actually purge
   `dvc-artifacts.kehao.me`
2. run the first public-history rewrite/filter pass from
   [History Rewrite Pathset, 2026-05-25](./history-rewrite-pathset-2026-05-25.md)
   in a fresh mirror or disposable clone, then rescan the rewritten result
