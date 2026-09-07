# Open Source Release

Current expansion work: [four-directory audit and release plan, 2026-09-07](./four-directory-release-audit-2026-09-07.md)
covers `development/`, `operator/`, `evolution/`, and `principles/`. The active
release mechanism remains [the projection contract](../../Paper/_design/migrations/public-branch-redesign.md);
older branch names and audit counts below describe earlier release work.

This surface owns Quine's open-source scope register.

It records what can move from private `main` and private operational
surfaces into the public `release` branch of `kehao95/quine`. It is not a
replacement for the dual-repository rules in
[`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).

## Release Principle

Quine should move toward a full open-source release by default.

The default boundary is:

- open source the repository as a research workbench, not only as a polished
  runtime package
- include unfinished drafts, papers, experiment notes, protocols, and lab
  control-plane material unless they contain secrets or private third-party
  material
- exclude credentials, tokens, private provider profiles, account identifiers
  that create access risk, machine-local access details, and private
  collaboration material that was not created for publication
- prefer redaction, templating, or path normalization over keeping whole
  surfaces private when only a small secret-bearing part is unsafe

This replaces the earlier "curated subset by default" posture for this release
line. Curation is still needed, but its purpose is leak prevention, legal /
consent review, size control, and claim clarity rather than hiding unfinished
research state.

## Current Baseline

- Work on private `main`; publish additive projection commits to `public/release`.
- The current release includes the runtime, tests, published portfolio,
  `development/`, `operator/`, `evolution/`, `principles/`, qcli's entrypoint,
  and selected root working documents.
- [Public scope and redactions](../PUBLIC-RELEASE.md) describes historical
  references and the ten development DVC pointers whose payloads remain private.
- Only `.dvc/public-manifest.txt` controls public data restore. This expansion
  does not add data to the manifest or synchronize private payloads.
- The source of truth is the Git-tracked lab tree, with declared public-copy
  transformations; private Git history and embedded self-source remain excluded.
- `Paper/`, `experiments/`, `profiles/`, hooks, and most developer tooling remain
  outside this release. The register below expresses longer-term scope intent,
  not proof that those paths have already been released.

## Scope Register

| Surface | Current read | Default action |
|---------|--------------|----------------|
| `cmd/quine/`, `internal/`, runtime-facing `Makefile` and Go module files | Core runtime implementation is the public center of the project. | Open after normal build/test and secret review. |
| `README.md`, `QUICKSTART.md`, `EVALUATION.md`, `LICENSE`, `CITATION.cff` | Public-facing entry docs are already release material. | Open, keeping claims aligned with current evidence. |
| `tests/runtime/` and deterministic validation helpers | Runtime contract tests support the public systems claim. | Open; project tracked files only. |
| `tests/model/` prompts, expected outputs, harness code, and retained run records | Evaluation behavior is part of the research workbench. | Open; restore run data only from the credential-scanned public manifest. |
| `experiments/` protocols, summaries, notebooks, experiment notes, and evidence metadata | Experiment history is part of the public research record, including partial or negative results. | Open by default; move bulky artifacts through DVC or documented public restore paths. |
| `Paper/papers/runtime-substrate/quine-arxiv/` and published systems-paper material | Stable public anchor. | Open. |
| Active paper drafts, workshop plans, and partial manuscripts | Unfinished writing should be visible as part of the research process. | Open by default unless it contains private collaboration material or confidential review context. |
| Research invitations, meeting prep, and correspondence-adjacent notes | Mixed public/private collaboration surface. | Review for consent and private third-party material before release. |
| `.env*`, provider profiles, local auth config, and machine-specific paths | High leak risk. | Do not open directly; replace with templates or documented examples. |
| Agent control-plane files such as `AGENTS.md`, `DEVELOPMENT.md`, `CLOSURE.md`, `EVALUATION.md` | The control-plane practice is part of the project and should be inspectable. | Open by default after secret/path/private-surface review. |
| `.claude/`, `.opencode/`, temporary diagnostics, scratch previews, and local tool state | Mostly local operating residue. | Exclude unless a specific sanitized artifact is intentionally promoted. |
| Public DVC artifact host config under `personal/quine-lab-public-artifacts/` | Workspace-local serving surface, not repo source. | Keep outside repo; publish only public restore instructions and selected object references. |

## Discussion Queue

- Define the secret/private-material rubric precisely enough that future agents
  can scan without repeatedly asking whether unfinished work is allowed.
- Keep the release gate green: projected-tree and DVC secret scan, build/test
  pass, public restore smoke, and claim/readme review.
- Decide how to handle private collaboration notes, third-party correspondence,
  reviewer context, and unpublished material from people other than the repo
  owner.
- Decide whether any huge generated run directories should stay out of git but
  be available through DVC/public artifact restore.

## Historical Risk Audits (May–September 2026)

The findings below retain their recorded dates; they are not a fresh audit of
today's remote history, storage, or CDN. Current expansion evidence is in the
[September 7 audit](./four-directory-release-audit-2026-09-07.md).

This pass used path-only, key-name-only, and redacted scanner output where
possible, so this file does not preserve any secret values.

Scanner status: `gitleaks` is installed at `$HOME/.local/bin/gitleaks`. The
2026-05-24 redacted worktree scan found 53 current findings, concentrated in
tracked `.env.*` profiles, model-run tapes/logs/workspaces, and config copies
inside retained run snapshots. A 2026-05-25 `git archive HEAD` rescan of the
then-current `lab` tree found 39 redacted findings across 24 tracked files, all
under the `generic-api-key` rule; the findings were concentrated in retained
`tests/model/.../runs/...` tapes/logs/workspaces plus one retained profile
review flag in `.env.gpt-5.4-codex-pool`. The pool profile was removed from
HEAD later on 2026-05-25, and a staged-tree rescan after the pool removal and
Harbor sweep-log DVC externalization found 38 redacted findings across 23
tracked files. After the `tests/model/**/runs/**` externalization and two
source/test false-positive reshapes, a staged-tree gitleaks scan found 0
findings. A whole-directory scan also showed that local `.dvc/cache` payloads
can contain scanner findings; public DVC payloads need their own review before
publication, not only git pointer review.

On 2026-05-26, the local `public` branch was imported at
`35322bbcbd30a09e8ce4b75a66ecd84e4fb1118b`. That branch rewrites the private
`lab` history into a DVC-backed public history, keeps raw run/log payloads in
private DVC, restores only the public-safe OAuth/Codex profile subset, and
passes both HEAD-archive and full-history gitleaks scans with 0 findings. See
[Public Branch Cut, 2026-05-26](./public-branch-cut-2026-05-26.md).

Detailed inventories:

- [Provider Profile Inventory, 2026-05-24](./provider-profile-inventory-2026-05-24.md)
- [Local Path Exposure Inventory, 2026-05-24](./local-path-exposure-inventory-2026-05-24.md)
- [DVC Payload Scan, 2026-05-24](./dvc-payload-scan-2026-05-24.md)
- [DVC Public Payload Allow/Deny, 2026-05-24](./dvc-public-payload-allowdeny-2026-05-24.md)
- [DVC Public Manifest](../../.dvc/public-manifest.txt)
- [Terminal-Bench DVC Archive Manifest, 2026-05-25](./dvc-terminal-bench-archive-manifest-2026-05-25.csv)
- [HEAD Release Scan, 2026-05-25](./head-release-scan-2026-05-25.md)
- [History Rewrite Pathset, 2026-05-25](./history-rewrite-pathset-2026-05-25.md)
- [Public Branch Cut, 2026-05-26](./public-branch-cut-2026-05-26.md)
- [Paper History DVC Policy, 2026-05-26](./paper-history-dvc-policy-2026-05-26.md)
- [Paper History Rewrite Validation, 2026-05-26](./paper-history-rewrite-validation-2026-05-26.md)

### Blocking Before Public Push

| Risk | Evidence from scan | Required action |
|------|--------------------|-----------------|
| Tracked provider profiles | The 2026-05-24 inventory found 76 current env/profile files. After owner direction, 63 provider profiles were removed from HEAD. At the time of that scan the tree kept `.env.example`, `.env.kimi-oauth`, eight root `codex` profiles, and two benchmark `codex` profiles. The earlier retained 13-file subset produced one redacted `generic-api-key` review flag in `.env.gpt-5.4-codex-pool`; that file was removed from HEAD on 2026-05-25. Since then `.env.kimi-oauth` has been retired in favor of repo-tracked `profiles/*.env` provider presets. A focused key-assignment grep found no non-empty `KEY` / `TOKEN` / `SECRET` / `PASSWORD` assignments in the remaining retained set. `.claude/settings.local.json` was removed from HEAD on 2026-05-24 and `.claude/` is now ignored. | Keep `.env.example` as the only non-provider public template. Repo-tracked provider profiles may remain when they do not contain live keys. Rotate any real credential that may have been committed. |
| Git history contains provider profiles | Private `lab` history includes old `.env.*` and provider profile files across March-May 2026, including files later deleted from HEAD. The local `public` branch cut removes those historical paths and restores only the public-safe OAuth/Codex profile subset at the tip. | Do not publish existing private `lab` history as-is. Publish the sanitized `public` branch after final review and remote push decision. Rotate any real credential that may have appeared in private history. |
| Raw run/tape history | Private history contains many `runs/`, `tape.jsonl`, `tape.log.yaml`, runtime logs, `session.json`, workspace snapshots, and older `test-runs` payloads. The local `public` branch cut converts these historical payload families to private DVC pointers; full-history tree scan found 0 raw `runs/` payload hits and 0 raw `test-runs` payload hits. | Git-history raw run/log blocker is handled on the local `public` branch. Raw DVC payloads remain private and are not in the public manifest unless separately reviewed and promoted. |
| Local path and machine exposure | Initial scan found `/Users/kehao95` in many files, `/home/kehao95` in dozens, `/srv/storage` in a few, and `/etc/nginx` in at least one repo file. The 2026-05-25 HEAD follow-up still found `/Users/kehao95` in 1,057 tracked files and `/home/kehao95` in 27 tracked files, concentrated in `tests/model/**` run records plus `experiments/active/benchmarks/bench-01-terminal-bench/analysis/sweep-logs/**`. Harbor sweep logs were moved out of git raw-file tracking and into a private DVC pointer on 2026-05-25; they are not in the public DVC manifest. After the model-run externalization, the staged tree has `/Users/kehao95` in 7 tracked files, `/home/kehao95` in 4 tracked files, `/srv/storage` in 5 tracked files, and `/etc/nginx` in 5 tracked files; all audited `tests/model` path counts are 0. Remaining path strings are in this open-source audit register plus archived experiment logs retained per owner direction. DVC text/log payloads still contain many local path occurrences. | Current HEAD path blocker from active model/Harbor logs is handled. Keep raw DVC payloads private unless explicitly reviewed; private history still needs rewrite/filter for old path-bearing logs and run records. |
| Public/private DVC store boundary | The `public` HTTPS surface and private write store have separate object roots. Publication is controlled by `.dvc/public-manifest.txt`; the 2026-09-06 rebuild produced an exact 29,073-object public store for 645 reviewed pointers, with zero missing or extra objects and full MD5 verification. Public readers use `./scripts/pull-public-dvc-manifest.sh`, not plain `dvc pull`, so private-retention pointers can remain in public history without implying public restore access. | Keep public synchronization manifest-bounded and credential-scanned; never mirror the private store wholesale into the public root. |
| Public CDN cache residue | The 2026-09-06 origin rebuild removed all 928 reviewed historical unsafe objects from the public root, but Cloudflare still returned 865 old content objects and 30 directory indexes from edge cache. An exact purge request was rejected because the installed token lacks Cache Purge permission. | Treat this as an active publication blocker until the exact URL list is purged with a Zone Cache Purge token and all 928 URLs are rechecked. |
| Third-party / collaboration material | Presentations and meeting-prep surfaces under `randall-beer-research-invitation`, research invitations, contact material, external-source PDFs, Escher images, and template/support assets need license/consent review. | Discussion item. Keep default-open for owner-authored drafts, but decide case-by-case for third-party PDFs/images/correspondence-adjacent notes. |
| Bulky binary and generated artifacts | Git object scan found large tarballs, built Quine binaries/successors, PDFs/PPTX/images, and model-run stdout/tape blobs; pack size is about 470 MiB. | Needs handling: move generated binaries and bulky raw outputs out of git where possible; keep source, summaries, and `.dvc` pointers in git. |

### Not Blocking By Itself

- Unfinished drafts, manuscripts, paper plans, experiment notes, and negative
  or partial results are not private merely because they are unfinished.
- Research-control documents are public by default unless they expose secrets,
  private machine access, or private third-party collaboration context.
- Synthetic test fixtures that look like emails or credentials can be public
  after confirming they are generated examples and not real personal data.

### Git History Decision

The current private `lab` history remains unsafe for direct public publication.

The public release path is now the sanitized local `public` branch, not the
private `lab` history. Its initial filtered-history cut was
`35322bbcbd30a09e8ce4b75a66ecd84e4fb1118b`; later open-source control-plane
commits can be cherry-picked onto that sanitized branch before external push.
That branch preserves useful project history while replacing historical raw
run/log payloads with private DVC pointers and removing secret-bearing local
configuration history. It should receive final review before any external git
push, but the git-history rewrite itself is no longer only a planned step.

## Open

- Convert this default-open scope register into a concrete release checklist.
- Add a repo-local scanner command/config/baseline around the installed
  `gitleaks`, covering both git-visible files and any DVC payloads intended for
  public restore.
- Create sanitized manifests for Terminal-Bench task archives and other raw DVC
  families that remain `deny-raw`, so public readers can recover results
  without restoring private tapes/logs.
- Review and push the local sanitized `public` branch to the external git
  remote when ready; keep the private `lab` history unpublished.
