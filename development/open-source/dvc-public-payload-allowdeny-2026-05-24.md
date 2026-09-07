# DVC Public Payload Allow/Deny, 2026-05-24

This is the first public-restore gate for DVC payloads.

It records release decisions at the path-pattern and payload-family level. It
does not preserve secret values, scanner match contents, token fragments, or
full scanner reports.

## Policy

Do not run a broad public artifact push.

Public DVC restore should be opt-in:

1. keep source, summaries, metadata, and selected `.dvc` pointers in git
2. push only explicitly allowed `.dvc` payloads to the public artifact remote
3. remove or mark private any `.dvc` pointer that should not be restorable from
   the public branch
4. publish sanitized derivatives instead of raw tapes/logs when raw payloads
   contain prompts, tool output, copied workspaces, local paths, or scanner
   findings

The canonical publication set lives in
[DVC Public Manifest](../../.dvc/public-manifest.txt).

`allow` means "eligible for publication" rather than "already present in the
public store". The isolated public store is currently empty by design; a
pointer becomes publicly restorable only after it is added to the manifest and
synced into the maintainer-only `public-localstore`.

After a store split or object removal, verify the public HTTPS surface with
real object-path requests or a fresh-clone restore smoke test. Do not rely only
on `dvc status -r public`, because cached directory objects can outlive the
origin move at the CDN layer.

## Input Evidence

Baseline scan:

- [DVC Payload Scan, 2026-05-24](./dvc-payload-scan-2026-05-24.md)
- 1,333 materialized DVC outs
- 73,397 selected text/log files scanned
- 2,581 redacted gitleaks findings
- local path findings concentrated in DVC text/log payloads

Follow-up clustering pass:

- 1,333 materialized DVC outs
- 69,555 candidate text/log files selected for clustering
- 569,549,779 selected bytes
- 1,960 redacted gitleaks findings
- used only for rule/path clustering; the broader baseline scan remains the
  controlling blocker

Payload shape:

| Shape | DVC outs | Declared bytes | Gate read |
|-------|----------|----------------|-----------|
| Raw `runs` payloads | 1,292 | 1,418,419,633 | default deny for public restore |
| ProgramBench campaign payloads | 33 | 142,439,878 | needs per-campaign review |
| ProgramBench curated evidence packages | 8 | 4,902,204 | conditional allow candidate |

Object format:

| Format | DVC outs | Declared bytes | Gate read |
|--------|----------|----------------|-----------|
| Directory payload | 1,077 | 1,548,072,789 | scanable but mostly raw logs/tapes |
| `.tar.zst` archive | 256 | 17,688,926 | deny raw after unpack-and-scan |

## Decision Classes

| Decision | Meaning | Public branch handling |
|----------|---------|------------------------|
| `allow` | Payload may be pushed to public DVC restore. | Keep pointer and publish restore instructions. |
| `conditional-allow` | Payload has passed automated checks but still needs final human scope/privacy review. | Keep pointer only if final review promotes it to `allow`; do not push broadly. |
| `sanitize` | Payload has useful public content, but raw restore is unsafe. | Keep or create summaries; publish a sanitized derivative, not the raw payload. |
| `deny-raw` | Raw payload should not be publicly restorable. | Remove pointer from public branch or mark private; do not push object. |
| `needs-scan` | Payload was compressed, skipped, or not deeply reviewed. | Do not publish until unpacked or otherwise scanned. |

## Allow List

Promotion date: 2026-05-25.

These are the first DVC payloads approved for selective public restore. They
must be promoted by exact pointer path, not by broad remote push.

| Pointer | Declared bytes | Decision | Reason |
|---------|----------------|----------|--------|
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p037-refalsification-55high-fast-unlimited-r1.dvc` | 859,850 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p038-baseline-falsification-default-direct-55high-unlimited-r1.dvc` | 404,791 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p038-baseline-falsification-default-direct-55high-unlimited-r2.dvc` | 405,596 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p039-evidence-accounting-default-direct-55high-unlimited-r1.dvc` | 666,141 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p040-family-matrix-default-direct-55high-unlimited-r1.dvc` | 704,246 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p041-immutable-ledger-default-direct-55high-unlimited-r1.dvc` | 464,795 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p042-lossless-substrate-default-direct-55high-unlimited-r1.dvc` | 561,547 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p044-affordance-cartography-fuse-default-55high-unlimited-r1.dvc` | 835,238 | `allow` | Curated ProgramBench evidence package; targeted redacted gitleaks scan returned 0 findings; audited owner/machine path scan returned 0 hits. |

Review notes:

- these packages contain raw research traces, but no live credentials were found
  by the targeted scanner pass
- email-like strings in the checked packages resolve to the public benchmark
  task domain, not owner private email
- short `sk-*` style strings appear in trace content, but did not trigger
  gitleaks and were not classified as live provider credentials
- `localhost` / loopback occurrences are local test/runtime references, not
  machine access details

Public restore status:

- before store isolation, `dvc status -r public` reported the eight allowed
  pointers in sync with the shared-store public remote on 2026-05-25
- before store isolation, a temporary fresh clone restored these exact
  pointers from `public` successfully on 2026-05-25
- historical restore smoke result under the shared store: 128 files fetched,
  124 files added
- after the 2026-05-25 store split, the isolated public store is intentionally
  empty until an exact pointer is added to the manifest and synced
- after the 2026-05-25 store split, the local origin returns `404` for tested
  old object paths with cache-busting query strings, but Cloudflare still
  serves at least one previously cached allowlisted object URL from edge cache;
  cache purge remains required before "public is empty" becomes fully true

## Conditional Allow Candidates

No payloads are currently `conditional-allow`.

The first ProgramBench evidence set was promoted to `allow` on 2026-05-25.

Historical promotion evidence:

| Pattern | DVC outs | Declared bytes | Current decision | Reason |
|---------|----------|----------------|------------------|--------|
| `experiments/active/benchmarks/bench-04-programbench/evidence/20260509-*.dvc` | 8 | 4,902,204 | promoted to `allow` | Targeted redacted gitleaks scan returned 0 findings; fixed-string path scan returned 0 hits for audited owner/machine paths; follow-up review found only public benchmark-domain email-like strings, short non-live `sk-*` style strings, and loopback runtime references. |

## Deny Raw Now

These payloads should not be publicly restorable as raw DVC objects.

| Pattern | Evidence | Decision | Public replacement |
|---------|----------|----------|--------------------|
| `experiments/active/benchmarks/bench-04-programbench/campaigns/horizontal/60-sirwart__ripsecrets.34c9e03/061-leaderboard-deciles-gpt54high-r1.dvc` | Follow-up clustering scan found 789 redacted findings under this one payload. The task may be synthetic, but raw restore still looks like a secret-bearing corpus to scanners and downstream users. | `deny-raw` | Keep benchmark summary and classifier notes; publish a sanitized synthetic-secrets statement if needed. |
| `experiments/active/benchmarks/bench-02-tb2-sweep/runs/*.dvc` | 15 raw sweep payloads, 296,076,553 declared bytes; follow-up scan found 284 redacted findings across this surface. | `deny-raw` | Keep leaderboard/analysis outputs and a sanitized run manifest. |
| `experiments/active/autopoiesis-succession/**/runs/*.dvc` | Multiple run families have scanner findings and raw model tapes/logs. Top affected surfaces include `6.26-affordance-reproduction`, `6.27-architecture-ablation`, `6.32-identity-transmutation`, and `6.55-descriptive-surface-hardening`. | `deny-raw` | Keep experiment READMEs, aggregate tables, and selected sanitized traces. |
| `experiments/active/cooperation-dynamics/**/runs/*.dvc` and `experiments/active/cooperation-dynamics/**/runtime/tapes/**` | Scanner findings and local path exposure are concentrated in raw population tapes/logs. | `deny-raw` | Keep summaries and aggregate metrics; publish sanitized excerpts only. |
| `experiments/archive/p1-physics/**/meta/env-after.txt` and archived raw run payloads | Prior scan identified copied environment snapshots and archived run residue. | `deny-raw` | Keep historical narrative; do not publish environment snapshots. |
| `tests/model/**/runs/*.dvc` | Small but security-adjacent model tests; follow-up scan found findings in at least one containment/hostile-script run. On 2026-05-25 the current raw run directories were externalized into 369 per-run private DVC pointers. | `deny-raw` | Keep test definitions and expected behavior; do not add these pointers to the public manifest. Generate sanitized fixtures if needed. |
| `experiments/active/benchmarks/bench-01-terminal-bench/tasks/**/runs/*.tar.zst.dvc` | 256 archives unpacked successfully on 2026-05-25. Redacted gitleaks found 146 findings across 58 archives; audited owner/machine path scan found 1,999 hits across 252 archives. | `deny-raw` | Keep Terminal-Bench task/result summaries; publish a sanitized archive manifest or derived metrics instead of raw archives. |

The Terminal-Bench raw archive manifest is recorded in
[Terminal-Bench DVC Archive Manifest, 2026-05-25](./dvc-terminal-bench-archive-manifest-2026-05-25.csv).
It records one row per archive with counts only; it does not include scanner
match contents.

## Needs Sanitize Or Targeted Review

These are not cleared for public restore. They may contain useful public
evidence, but the raw payload shape is too broad.

| Pattern | DVC outs / bytes | Decision | Next action |
|---------|------------------|----------|-------------|
| `experiments/active/benchmarks/bench-04-programbench/campaigns/**` except the `ripsecrets` deny-row | 32 outs / about 140 MiB | `sanitize` | Classify each campaign by benchmark task, scanner result, copied repo content, and license posture. |
| `experiments/active/autopoiesis-succession/**/runs/*.dvc` | 51 outs / large raw successor-morphospace payloads | `sanitize` | Keep aggregate morphology metrics and selected sanitized traces; raw payloads stay private until reviewed. |
| `experiments/active/null-ecology/**/runs/*.dvc` | 77 outs / raw open-ended evolution payloads | `sanitize` | Keep summaries and reproduction metadata; raw restore needs scanner/path review. |
| `experiments/active/trace-inheritance/**/runs/*.dvc` | 19 outs / raw generation-ecology payloads | `sanitize` | Keep aggregate traces; review raw tapes before any public push. |
| `experiments/active/early-existence/**/runs/*.dvc` | 266 outs / mostly historical lineage runs | `sanitize` | Keep experiment docs; publish only sanitized or small curated evidence. |

## Public Push Rules

Use the manifest plus explicit local-store sync only:

```bash
printf '%s\n' \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p037-refalsification-55high-fast-unlimited-r1.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p038-baseline-falsification-default-direct-55high-unlimited-r1.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p038-baseline-falsification-default-direct-55high-unlimited-r2.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p039-evidence-accounting-default-direct-55high-unlimited-r1.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p040-family-matrix-default-direct-55high-unlimited-r1.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p041-immutable-ledger-default-direct-55high-unlimited-r1.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p042-lossless-substrate-default-direct-55high-unlimited-r1.dvc \
  experiments/active/benchmarks/bench-04-programbench/evidence/20260509-p044-affordance-cartography-fuse-default-55high-unlimited-r1.dvc \
  >> .dvc/public-manifest.txt

./scripts/sync-public-dvc-manifest.sh
```

Public readers should restore through the same manifest boundary:

```bash
./scripts/pull-public-dvc-manifest.sh
```

Do not treat plain `dvc pull` as the public restore contract. It attempts every
DVC pointer in the checkout, including private-retention pointers that are
intentionally absent from the public store.

Do not use:

```bash
dvc push -r public
```

`public` is the reader-facing HTTPS remote. Maintainers populate the isolated
public store through the local-only `public-localstore` remote, which the sync
script targets.

Do not publish `.dvc/cache` wholesale. Do not assume a `.dvc` pointer is safe
because the pointer file itself passed git-visible scans.

## Immediate Work Queue

1. Decide when the first exact pointer set should move from `allow` into the
   public manifest.
2. Create sanitized manifests for Terminal-Bench task archives,
   `p4-harbor/4.2-tb2-sweep`, and ProgramBench
   campaign runs so public readers can recover the result without raw restore.
3. Decide whether denied raw `.dvc` pointers stay in the public branch as
   private historical markers or are removed from the public cut.
4. Purge the `dvc-artifacts.kehao.me` edge cache so the isolated empty public
   store becomes externally true, not only origin-true. If public objects need
   to be cleared and rebuilt locally first, use
   `./scripts/rebuild-public-dvc-store.sh --yes`.
5. Re-run public artifact restore smoke tests only on promoted manifest
   pointers.
