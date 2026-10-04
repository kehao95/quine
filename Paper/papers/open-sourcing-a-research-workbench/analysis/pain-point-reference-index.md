# Pain-Point Reference Index

This file is a working evidence index for the blog dossier
`open-sourcing-a-research-workbench`.

The goal is not yet to finalize prose. The goal is to make each pain point
recoverable later with nearby source anchors so drafting can pull evidence
quickly without re-running the full repo archaeology every time.

## Collection Rule

For each pain point, try to collect:

- canonical owner surfaces
- supporting control-plane surfaces
- key commit anchors
- history notes or evolution summaries
- unresolved questions or missing evidence

Use exact file paths and commit SHAs whenever possible.

## Meta-Layer Framing Inventory

This section is not choosing the final essay thesis. It is recording the
different higher-level frames that may later organize the story.

### Current capture

- public release as making a research workbench public rather than releasing a
  conventional software repo
- public release as preserving an experimental memory layer that remains
  rereadable later
- public release as rebuilding the sensitive-data threat model once logs,
  shell output, and inherited host state are treated as possible leak surfaces
- public release as publication-boundary design: retention, storage, and
  public restore should not automatically share one boundary
- public release as preserving historical authenticity rather than solving the
  problem by deleting the risky past
- public release as a scale and Git-tractability problem involving blob weight,
  file count, and repository performance
- public release as redefining the public object: from something like a public
  runtime toward something closer to a public research workbench
- public release as building a release control plane with scans, inventories,
  rewrite policy, validation, and publication mechanics
- public release as balancing open science with operational safety

### Candidate references

- `README.md` - public landing page that already frames Quine as both runtime
  and research workbench
- `development/open-source/README.md` - strongest current release principle,
  boundary design, scan/risk register, public-history policy, and control-plane
  surface
- `experiments/PROTOCOL.md` - strongest evidence that experiment traces are not
  just residue but part of the research record
- `.dvc/public-manifest.txt` - clearest compact artifact-publication boundary
- `development/open-source/public-branch-cut-2026-05-26.md` - branch/object
  story, public-history cut, and historical DVC transformation
- `development/open-source/history-rewrite-pathset-2026-05-25.md` - strongest
  history-sanitization framing and evidence that HEAD cleanup was not enough
- `development/open-source/provider-profile-inventory-2026-05-24.md` - good
  source for the broadened threat model around private-first residue
- `development/open-source/local-path-exposure-inventory-2026-05-24.md` - good
  source for logs and traces as leak surfaces
- `experiments/protocol/dvc-externalization-rollout.md` - DVC as structure-
  preserving externalization and staged publication mechanism
- `analysis/size-and-blob-audit.md` - quantitative scale / Git-tractability
  evidence for this essay dossier
- `/REDACTED_HOME/workspace/personal/quine-lab-public-artifacts/README.md` -
  companion serving-layer surface; useful for policy vs actual public
  reachability
- `evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md`
  - retrospective surface that may later help bridge several of these frames

### Open questions

- which of these frames eventually deserves to become the main narrative voice?
- which frames should remain explicit sections versus becoming implicit
  structure under later pain points?
- is there one especially strong file or session that captures the transition
  from `public runtime` to `public research workbench` in the user's own words?

## Pain Point 1: Experimental Logs Are Evidence And A Reusable Read Surface

### Current capture

- Quine is not a normal software repo; it is a hybrid of code, research,
  experiment record, raw data, and control-plane surfaces.
- In a normal software project, experiment logs usually should not be public
  and often do not belong in the main repo surface.
- In Quine, because the experimental object is an LLM with real stochasticity,
  preserving code + prompts + model + environment is not enough to recover the
  exact trajectory later.
- Therefore logs had to be retained.
- Early on, a practical implementation of that principle was to track logs
  directly in the repo.
- Over time this pulled in material that was not actually the experimental
  record we wanted in the same retained surface: generated binaries, incidental
  outputs, and non-important results.
- The project then started to become more explicit about retained-log criteria,
  retained-log path stability, and which artifacts should count as important
  enough to preserve.
- There was also an attempt to package logs as zip archives. That reduced
  footprint somewhat but made the record less readable.
- Readability matters because logs are not only archived proof; they are a
  long-lived read surface that humans and agents may revisit when theory or
  analysis questions change.

### Candidate references

- `README.md` - says Quine is both a runtime and a research workbench, with
  code, experiments, papers, scripts, and observation overlays in one repo
- `Paper/philosophy/self-maintaining-systems.md` - repo-level doctrine that
  keeps runtime code, experiments, evidence, and theory coherent together
- `experiments/PROTOCOL.md` - strongest evidence doctrine for why stochastic
  runs, tapes, and retained traces matter, and why later theory changes can
  require re-mining old artifacts
- `development/contracts/code-maintenance.md` - preserved baselines/evidence as
  part of code maintenance
- `development/contracts/active-doc-addressability.md` - active docs may only
  point at evidence the repo intentionally keeps live
- `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md` - concrete
  experiment-level selection logic: runtime-private residue pruned, retained
  evidence kept in `meta/`, snapshots, `*.log`, and `tapes/`
- `experiments/active/cooperation-dynamics/analysis/BEHAVIOR-ANALYSIS.md` - later reread of raw traces
  as a required analysis step; invalid/incomplete runs still preserved because
  they taught durable lessons
- `experiments/active/cooperation-dynamics/coop-20-full-disclosure-baseline/README.md` - kept
  pointers marked `reread pending`; stale path references rejected in favor of
  actual retained evidence
- `experiments/active/autopoiesis-succession/README.md` - explicit later “theory
  reread” over retained traces
- `experiments/active/autopoiesis-succession/succ-29-retained-run-morphology-inventory/README.md`
  - old retained runs recoded into morphology records; later reread changes
  interpretation without changing records
- `experiments/active/benchmarks/bench-04-programbench/analysis/README.md` - later
  analysis note explicitly re-reads retained campaigns
- `experiments/active/benchmarks/bench-04-programbench/README.md` - filtered DVC
  tracking keeps logs, tapes, eval JSON, and trace subsets while excluding
  opaque bulk
- `experiments/active/benchmarks/bench-04-programbench/evidence/README.md` -
  selective retained-evidence surface, one bundle per evidence id
- `experiments/active/benchmarks/bench-04-programbench/campaigns/.dvcignore` -
  excludes opaque archive bundles because more legible retained state already
  exists elsewhere
- `experiments/active/benchmarks/bench-01-terminal-bench/tasks/log-summary-date-ranges/runs/2026-02-22__14-43-28.tar.zst.dvc`
  - concrete tarred run archive pointer; useful for archive packaging story
- `development/open-source/dvc-public-payload-allowdeny-2026-05-24.md` - raw
  tarred log-heavy archives denied for public release and replaced by sanitized
  manifests or derived metrics
- `scripts/README.md` - helper tooling for DVC-tracking retained run units and
  expanding Harbor `*.tar.zst` archives into sanitized run directories

### Commit anchors

- `121f6dc62a26f9d34d4bb377397ed02efd57e299` - earliest current-history pattern
  of repo-governed `runs/*.dvc` tracking
- `dae85af72c8b7406a2de16a56f8abc6feabcd84b` - adds `QUINE_RETENTION_DIR` and
  formalizes `runs/<runid>/log/` as canonical retained evidence
- `be18ed600cea382379812328aa84c13318a8e98f` - backfits older archived
  experiments to the stable retention-dir pattern
- `b6abb6d582f94d080b91a75578291907e3174494` - retains trimmed P8 logs while
  dropping reproducible world-state bulk
- `3138df3372422bba38ef10e45c236fcfef64db7d` - narrows citation-grade export to
  mission, session, `runtime.log`, and retained `tapes/`
- `eeed8734c7f4d6472ac9e23a939541cd0b63cb1f` - directly checks in traces, logs,
  validation outputs, and relation traces for later synthesis
- `2b2c6af8b9cbf39e4faf25a3788b9fe67149e835` - converts raw retained traces into
  one-bundle-per-`.dvc` reusable evidence surfaces
- `c2d08280f7a7abd9ab8b2674fc40e2ee55cf74d2` - formalizes selective
  externalization, stable bundle linking, and restore-path requirements
- `13758bf4e6c14d2a7a6d0956ecee3c1e3358ae90` - explicit tar-archive workflow for
  retained run payloads
- `5bced7d27cd68cebb2fb2755434f62c2c38f524e` - tightens stable evidence linking
  toward adjacent `.dvc` pointer surfaces

### Strongest reread / reanalysis references

- `experiments/active/benchmarks/bench-04-programbench/analysis/README.md` - says a
  later analysis note re-reads all retained `gron` campaigns through `035`
- `experiments/active/autopoiesis-succession/README.md` plus
  `experiments/active/autopoiesis-succession/succ-29-retained-run-morphology-inventory/README.md`
  - explicit later theory reread and recoding of prior retained runs

### Open questions

- Which repo surface best captures the explicit retained-log criteria?
- Which commits best represent the progression:
  `retain everything -> filter better -> stabilize retained paths -> zip attempt`?
- Do we have one especially strong example of an old log being reread under a
  later theoretical frame?

## Pain Point 2: Sensitive Data Turned Into A Branch-Topology And Maintenance Problem

### Current capture

- One early idea was to maintain a `private` branch and a `public` branch.
- At first this looked conceptually clean: private keeps the full research
  reality; public exposes what is safe.
- Later the `public` branch became hard to maintain.
- One reason was that the imagined public object was initially closer to a
  `public runtime`.
- But the real repository contains much more than runtime code: theory,
  experiments, tests, and other research surfaces.
- Once that broader reality is acknowledged, a narrow public-runtime branch
  becomes an incomplete representation of the project.
- But if the public branch includes more of the real project, it inherits much
  more sensitive-data review and maintenance pressure.
- So this was not just a secret-scrubbing problem. It became a repo-topology
  problem and a question of what the real public object should be.
- Another part of this is historical posture: the project began operationally
  as a private project.
- Because it was private-first, and because the private branch was expected to
  carry the full reality, sensitive material was tracked more freely than it
  would have been in a public-first repo.
- That included sensitive data such as keys.
- Later some of that material was deleted or removed, but tip deletion was only
  one part of the problem.
- The agent could call shell commands, and if it inherited host environment
  state, shell output could directly print sensitive keys or other
  secret-bearing environment values.
- That means the leakage surface was not only secret files or explicit config;
  the runtime traces themselves could become sensitive.
- Because the project also wanted to preserve logs heavily, every retained log
  became a potential information-leak source.
- This creates a direct conflict with Pain Point 1: the logs are both important
  evidence/read surfaces and possible secret-leak surfaces.

### Candidate references

- `development/open-source/public-branch-cut-2026-05-26.md` - public-branch cut
  note; useful for branch-topology story and what the public object eventually
  became
- `development/open-source/README.md` - release principle that unfinished
  drafts, theory, experiments, tests, and notes should be public by default
  unless specifically unsafe; good evidence that the public object stopped
  being runtime-only
- `README.md` - public landing page openly presents Quine as both runtime and
  research workbench, not runtime-only
- `evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md`
  - likely strong retrospective surface for the shift from narrower release
  assumptions to broader public-workbench logic
- `development/open-source/history-rewrite-pathset-2026-05-25.md` - likely
  relevant because branch pressure and history pressure were entangled
- `development/open-source/provider-profile-inventory-2026-05-24.md` - likely
  strong source for the scale of env/profile exposure in history
- `development/open-source/local-path-exposure-inventory-2026-05-24.md` - likely
  useful for how raw logs/traces can carry sensitive host residue
- `development/open-source/head-release-scan-2026-05-25.md` - useful for
  retained run records and benchmark logs as remaining git-visible exposure
- `development/open-source/history-rewrite-drop-paths.txt` - concrete path
  families removed or rewritten during sanitization

### Open questions

- Which files or commits most clearly show the repo stopped treating the public
  object as runtime-only?
- Is the strongest evidence here branch-topology documentation, or the later
  release principle that unfinished drafts / theory / experiments should remain
  public by default unless specifically unsafe?
- Which source best captures the specific risk that shell output inherited from
  the host environment could print secrets directly into preserved logs?

## Pain Point 3: Retention And Publication Could Not Share The Same Boundary

### Current capture

- Retention and publication could not use the same boundary.
- DVC pointers or tracked artifacts are not automatically public restore.
- The repo needed a difference between:
  - data retained as part of project history
  - data actually restorable by public readers
- This became increasingly explicit over time.
- An important part of why this became workable is that the DVC server was
  separately controlled by the project owner.
- Because the storage backend was under personal control, public release could
  be implemented as an audited whitelist rather than an all-or-nothing dump.
- That meant only experiment records that had been reviewed and judged safe
  could be made publicly restorable.
- This let the project preserve a broad private experiment memory while only
  publishing a reviewed subset.

### Candidate references

- `README.md` - public artifact restore is manifest-scoped; plain `dvc pull` is
  not the public restore contract
- `.dvc/public-manifest.txt` - canonical publication boundary; `allow` means
  eligible, not already public
- `development/open-source/README.md` - strongest current contract for isolated
  public store, manifest-backed publish/pull, and public history that may still
  contain private-only DVC pointers
- `development/open-source/dvc-public-payload-allowdeny-2026-05-24.md` -
  explicit `allow` / `sanitize` / `deny-raw` gate for DVC payload families
- `development/open-source/head-release-scan-2026-05-25.md` - concrete examples
  of retained DVC pointers that are private-only and not in
  `.dvc/public-manifest.txt`
- `development/open-source/public-branch-cut-2026-05-26.md` - public branch
  keeps historical raw runs as private DVC pointers while manifest intentionally
  remains empty
- `experiments/protocol/dvc-externalization-rollout.md` - transition doc:
  reader `public`, maintainer `private` / `public-localstore`, and later
  manifest-backed publication
- `experiments/active/benchmarks/bench-04-programbench/README.md` - early
  retention-side surface with separate maintainer vs public remote roles
- `experiments/active/benchmarks/bench-04-programbench/evidence/README.md` - early
  public-read-only restore examples; good precursor even if not final contract
- `evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md`
  - retrospective phrasing that “tracked in repo” and “restorable by public
  readers” became different facts
- `/REDACTED_HOME/workspace/personal/quine-lab-public-artifacts/README.md` -
  workspace-local companion surface for the personally controlled public
  artifact-serving layer

### Commit anchors

- `011cc8e39` - DVC artifact retention becomes a deliberate repo-level mechanism
- `56940665e` - first repo docs that public restore is its own surface
- `1961eccb5` - experiment docs stop hard-coding only the private artifact path
- `4d4b9be2d` - explicit allow/deny/sanitize review for retained DVC payloads
- `c86ae6e33` - adds `.dvc/public-manifest.txt` and isolated public-store
  workflow
- `8f3b45683` - reader restore becomes manifest-scoped rather than blanket
  `dvc pull`
- `6f09526e7` - public history can retain private DVC pointers while public
  publication stays separate

### Best single reference

- `development/open-source/README.md` - best current statement that private
  retention and public publication are different things

### Open questions

- `TBD`

## Pain Point 4: A Clean HEAD Was Not Enough; History Itself Needed Work

### Current capture

- Cleaning current tip was not enough.
- Old history still contained env/profile files, path-bearing run logs, raw DVC
  payload families, and bulky/generated Paper artifacts that were not suitable
  for direct publication.
- So the problem was not only current-state hygiene; history itself needed
  explicit rewrite / sanitization work.
- After choosing the DVC direction, this also became a historical-authenticity
  problem.
- The repo already had more than a thousand commits of history.
- Simply deleting all problematic historical experiment material might have
  produced a clean result, but it would also have damaged the authenticity of
  the repository history.
- Preserving the authenticity of the history mattered.
- So the preferred move was to preserve the historical commit structure while
  converting risky experimental-record portions of that history into DVC-backed
  form.
- This is closer to a history replay / transformed-history preservation story
  than to a simple “delete bad history” story.

### Candidate references

- `development/open-source/history-rewrite-pathset-2026-05-25.md` - strongest
  statement separating HEAD cleanup from rewrite work
- `development/open-source/history-rewrite-drop-paths.txt` - machine-readable
  first-pass `git filter-repo` drop set
- `development/open-source/public-branch-cut-2026-05-26.md` - records actual
  filtered branch cut and validation
- `development/open-source/provider-profile-inventory-2026-05-24.md` - history
  exposure of env/profile files
- `development/open-source/local-path-exposure-inventory-2026-05-24.md` - local
  path exposure in raw run logs and evidence payloads
- `development/open-source/dvc-payload-scan-2026-05-24.md` - DVC payloads also
  needed review
- `development/open-source/paper-history-dvc-policy-2026-05-26.md` - Paper
  history also required policy and rewrite treatment
- `development/open-source/paper-history-rewrite-validation-2026-05-26.md` -
  validates Paper rewrite candidate
- `development/open-source/public-branch-cut-2026-05-26.md` - records `1272`
  commits in the first sanitized public-history cut and shows the rewrite
  converted raw history families into private DVC pointers instead of just
  erasing the branch narrative
- `development/open-source/paper-history-dvc-policy-2026-05-26.md` - especially
  useful because it explicitly says to replace affected historical content with
  DVC pointer files while retaining payloads in the private DVC store
- `evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md`
  - likely useful retrospective surface for how DVC externalization and history
  rewrite became one combined release move
- `Codex session history` - user says operational detail for this history-replay
  step also lives in the Codex sessions and should later be indexed if needed

### Commit anchors

- `dfeda27fe` - adds rewrite plan and `history-rewrite-drop-paths.txt`
- `a1078750d` - externalizes Harbor sweep logs as tip-cleanup step before
  rewrite
- `2a1a9d2bf` - externalizes model run logs; HEAD cleaner while history still
  unsafe
- `6f09526e7` - validated filtered-history public branch cut
- `6012994f7` - Paper history DVC policy
- `c4e581281` - Paper history rewrite validation

### Additional notes

- `public-branch-cut-2026-05-26.md` is especially useful for the “do not lose
  history authenticity” story because it preserves a large commit count while
  changing the storage boundary of historical run/log payloads.
- `paper-history-dvc-policy-2026-05-26.md` is especially useful for the
  “history replay into DVC-backed form” story because it describes per-commit
  replacement with DVC pointers while keeping payloads in private retention.

### Size / file-count evidence

- Current local repo across all refs still reports `size-pack: 470.88 MiB` via
  `git count-objects -vH`; this is useful as a reminder that local size can
  stay large when private history is still present.
- A fresh `--no-local --single-branch` clone of `lab` reports
  `size-pack: 469.28 MiB`.
- A fresh `--no-local --single-branch` clone of current `main` reports
  `size-pack: 26.74 MiB`.
- That makes the public-facing `main` clone about `442.54 MiB` smaller than the
  private `lab` clone at pack level, or roughly a `94.3%` reduction.
- `lab` tip currently has:
  - `6,043` tracked files
  - `71,963,502` bytes in the tip tree
  - `47,565` reachable blobs
  - `1,481,535,626` reachable blob bytes
- `main` tip currently has:
  - `5,822` tracked files
  - `25,784,105` bytes in the tip tree
  - `16,099` reachable blobs
  - `185,014,011` reachable blob bytes
- `lab -> main` reduction:
  - tracked files: `-221` (`-3.7%`)
  - tip-tree bytes: `-46,179,397` (`-64.2%`)
  - reachable blob count: `-31,466` (`-66.2%`)
  - reachable blob bytes: `-1,296,521,615` (`-87.5%`)
- `history-rewrite-pathset-2026-05-25.md` records `93,712` distinct historical
  run/tape/session/log/workspace paths plus `527` distinct tar/tar.zst paths,
  which is a strong history-scale metric for Git-performance pressure.
- `development/open-source/README.md` records that before active run
  externalization there were `/REDACTED_HOME` hits in `1,057` tracked files and
  `/REDACTED_HOME` hits in `27` tracked files, concentrated in model run records
  and Harbor sweep logs; after externalization the staged tree dropped to `7`
  and `4` respectively.

### Representative historical large blobs

- `33,130,457` bytes -
  `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-006.tar.zst`
- `21,335,354` bytes -
  `Paper/essays/_meta/external-sources/slate-references/07-altera-project-sid-piano-architecture.pdf`
- `18,938,639` bytes -
  `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-007.tar.zst`
- `18,175,632` bytes -
  `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-006/tasks/video-processing/trial-01/2026-02-26__03-26-34/video-processing__LFo4Jh8/agent/command-1/stdout.txt`
- multiple historical successor / rebuilt `quine` binaries in
  `experiments/active/autopoiesis-succession/**/runs/**` at about `13.4 MiB` each
- `13,401,186` bytes -
  `Paper/Quine: Autopoietic Intelligence via Recursive POSIX Processes/output/paper.docx`

### Best single reference

- `development/open-source/history-rewrite-pathset-2026-05-25.md` - strongest
  explicit statement that tip scanned clean while private history remained
  unsafe

### Open questions

- `TBD`

## Pain Point 5: The Open-Source Release Needed Its Own Control Plane

### Current capture

- Release prep could not stay scattered across diffs, commands, and chat.
- It needed its own repo-local control plane: owner surface, audits, risk gate,
  pathsets, payload policy, and validation notes.

### Candidate references

- `development/README.md` - registers `open-source/` as a first-class
  development surface
- `development/open-source/README.md` - owner surface for purpose, release
  principle, baseline, risk audit, and git-history decision
- `development/open-source/dvc-public-payload-allowdeny-2026-05-24.md` - turns
  release into a governed payload gate
- `.dvc/public-manifest.txt` - canonical publication manifest
- `development/open-source/head-release-scan-2026-05-25.md` - tip audit for
  whether current `lab` HEAD was releasable
- `development/open-source/provider-profile-inventory-2026-05-24.md` -
  supporting audit input
- `development/open-source/local-path-exposure-inventory-2026-05-24.md` -
  supporting audit input
- `development/open-source/dvc-payload-scan-2026-05-24.md` - supporting audit
  input

### Commit anchors

- `8df17dfa5` - creates `development/open-source/README.md` and registers the
  new surface
- `34fff5820` - adds first three dated audit files
- `4d4b9be2d` - records DVC public payload gate
- `c86ae6e33` - isolates public DVC publication flow

### Best single reference

- `development/open-source/README.md` - best single owner surface for release
  control-plane creation

### Open questions

- `TBD`

## Pain Point 6: Public Readers Needed A Truthful Public Surface

### Current capture

- The public-facing surfaces had to honestly tell readers what Quine is.
- That includes:
  - Quine as runtime
  - Quine as research workbench
  - what is draft / working material
  - what public artifact restore actually means
  - how to cite Quine truthfully

### Candidate references

- `README.md` - outside-reader landing page; says Quine is a POSIX-native
  runtime and also a research workbench, warns that some material is draft, and
  explains manifest-scoped public artifact restore
- `CITATION.cff` - machine-readable citation truth with scoped runtime/software
  vs process-habitat / generative-ALife references
- `QUICKSTART.md` - stable runnable interface for public readers
- `Paper/README.md` - separates current claim boundary, evidence, active
  papers, and ideas/future work
- `Paper/philosophy/north-star.md` - stable runtime ontology statement
- `Paper/core/README.md` - implementation-backed current runtime canon
- `.dvc/public-manifest.txt` - authoritative public artifact publication/restore
  contract
- `experiments/protocol/dvc-externalization-rollout.md` - detailed meaning
  behind the README summary
- `DEVELOPMENT.md` - internal rule that root README is the human public landing
  page and should not carry maintainer/control-plane posture

### Commit anchors

- `56940665eed5cd6b51255d50547a75135c8b1322` - first clear reader-facing public
  artifact surface note
- `c86ae6e33029f7e8f50dbddace0a64609fb99212` - creates
  `.dvc/public-manifest.txt` and public sync/rebuild scripts
- `8f3b456830275b93d6fea219108924966af44680` - explicit manifest-pull rule
- `f7bffc88036f4c776f0afcff0f755c8142cc1e0f` - adds freshness/ownership cards to
  public/index surfaces
- `24e15410b4b4bfe5764f8e53f3f7364a3147fd56` - recasts root README as the
  outside-reader surface
- `58495e1808cb3998897c48be11a83562c017e373` - rewrites README and `CITATION.cff`
  around two scoped public citation anchors

### Best single reference

- `README.md` - best single surface for public-reader truthfulness

### Open questions

- Truth is currently distributed across multiple surfaces; do we later want the
  essay to emphasize this distribution as a feature or as a complexity cost?

## Cross-Cutting Solution Note: DVC As A Structure-Preserving Repo/Data Split

### Current capture

- DVC appears to be one of the major mechanisms that solved a large portion of
  the project’s practical storage / retention / openness problems.
- The point is not only that DVC can move data outside Git.
- More specifically, it lets the project upload heavy data outside the Git
  repository while still preserving the data’s place inside the project’s
  structure.
- It also preserves the sense that the data changes together with the project,
  instead of becoming an unrelated external bucket or manual archive.
- This matters because the project did not want to simply delete data from the
  repo’s conceptual structure.
- It wanted to move payloads out of Git while keeping addressability,
  organization, and change legible.
- This seems to be a cross-cutting solution note that connects to multiple pain
  points:
  - retained logs
  - large evidence bundles
  - private retention vs public publication
  - public/private branch maintenance pressure
- DVC also had an unexpected operational advantage: because the DVC server was
  separately controlled by the project owner, publication could be governed
  directly at the serving layer.
- This made it possible to maintain a whitelist of experiment records that had
  already been reviewed and judged safe for public release.
- That means DVC was not only "move data out of Git"; it also became a way to
  publish only audited experiment records while keeping the larger retained
  memory private.

### Candidate references

- `experiments/protocol/dvc-externalization-rollout.md` - strongest likely
  protocol surface for how data moves outside Git while remaining inside repo
  semantics
- `experiments/active/benchmarks/bench-04-programbench/README.md` - experiment-level
  surface showing DVC-retained raw campaigns with maintained repo structure
- `experiments/active/benchmarks/bench-04-programbench/evidence/README.md` - one
  bundle per evidence id; strong example of reusable retained evidence surface
- `experiments/active/benchmarks/bench-04-programbench/campaigns/.dvcignore` -
  distinguishes legible retained state from opaque archive bundles
- `scripts/README.md` - helper tooling for DVC-tracking retained run units and
  archive expansion
- `development/open-source/README.md` - strongest current explanation of the
  split between DVC as private retention infrastructure and DVC as public
  publication surface
- `.dvc/public-manifest.txt` - shows that public DVC is not just “all retained
  DVC data”
- `README.md` - public-facing explanation that public artifact restore is
  manifest-scoped
- `/REDACTED_HOME/workspace/personal/quine-lab-public-artifacts/README.md` -
  workspace-local companion surface for the owner-controlled publication layer

### Commit anchors

- `011cc8e39` - initializes DVC artifact retention
- `2b2c6af8b9cbf39e4faf25a3788b9fe67149e835` - moves benchmark evidence bundles
  to DVC
- `c2d08280f7a7abd9ab8b2674fc40e2ee55cf74d2` - centralizes experiment DVC
  retention guidance
- `56940665e` - public artifact surface becomes explicit
- `1961eccb5` - adds public remote to shared config
- `c86ae6e33` - isolated public DVC publication flow
- `8f3b45683` - manifest-backed public DVC pull path

### Open questions

- Is the strongest source for this idea a rollout/protocol doc, a commit
  series, or experiment-level README surfaces?
- Should DVC be introduced in the essay as a tool, an infrastructure layer, or
  part of the repository’s self-model?

## Candidate Additional Decisions / Evidence Buckets

These are candidate important points that may or may not become full sections
later. They are recorded here so the evidence does not get lost.

## Candidate A: The Public Object Changed From `public runtime` To `public research workbench`

### Current capture

- The release target widened over time.
- The project stopped treating the public object as only a runtime package.
- It increasingly treated Quine as a public research workbench.
- That means theory, experiments, tests, unfinished drafts, notes, and
  control-plane material became part of the intended public object unless they
  were specifically unsafe.

### Candidate references

- `development/open-source/README.md` - release principle explicitly says to
  open-source the repository as a research workbench, not only as a polished
  runtime package
- `README.md` - public landing page says Quine is both a runtime and a
  research workbench
- `development/open-source/README.md` scope register - explicitly treats
  `tests/model/`, `experiments/`, active drafts, and control-plane files as
  public by default after review

## Candidate B: HEAD Cleanup And History Rewrite Were Two Different Stages

### Current capture

- There is a real two-stage implementation distinction here.
- A cleaner current tree did not mean public-safe history.
- The repo first made current HEAD less dangerous, then separately handled old
  history through rewrite / DVC conversion.

### Candidate references

- `development/open-source/history-rewrite-pathset-2026-05-25.md` - clearest
  statement separating HEAD cleanup from history rewrite
- `development/open-source/head-release-scan-2026-05-25.md` - current-tip scan
- `development/open-source/public-branch-cut-2026-05-26.md` - later
  filtered-history execution
- commits `a1078750d`, `2a1a9d2bf`, `dfeda27fe`, `6f09526e7`

## Candidate C: Open-Source Release Became Its Own Maintained Control Plane

### Current capture

- The release effort became a maintained object with its own owner surface,
  risk gate, scans, pathsets, and validation files.

### Candidate references

- `development/README.md` - registers `open-source/` as a first-class surface
- `development/open-source/README.md` - owner surface
- commits `8df17dfa5`, `34fff5820`, `4d4b9be2d`, `c86ae6e33`

## Candidate D: Public-Facing Truth Surfaces Had To Be Rewritten

### Current capture

- The release also required public-facing docs to tell the truth about what
  Quine is, what is draft, what public artifact restore means, and how to cite
  the project.

### Candidate references

- `README.md`
- `CITATION.cff`
- `QUICKSTART.md`
- `Paper/README.md`
- `Paper/philosophy/north-star.md`
- `Paper/core/README.md`
- `DEVELOPMENT.md`
- commits `24e15410b4b4bfe5764f8e53f3f7364a3147fd56`,
  `58495e1808cb3998897c48be11a83562c017e373`,
  `f7bffc88036f4c776f0afcff0f755c8142cc1e0f`

## Candidate E: Paper, Third-Party, And Bulky Artifacts Needed Separate Policy

### Current capture

- Another difficult class of material was not secret but still not well-suited
  for public git history.
- This includes third-party PDFs, presentations, graphs, generated paper
  outputs, and other bulky or collaboration-adjacent material.

### Candidate references

- `development/open-source/paper-history-dvc-policy-2026-05-26.md`
- `development/open-source/paper-history-rewrite-validation-2026-05-26.md`
- `development/open-source/README.md` risk audit rows for third-party /
  collaboration material and bulky binary/generated artifacts
- commits `6012994f7`, `c4e581281`

## Candidate F: The Public Artifact Serving Layer Lived Outside The Repo

### Current capture

- The publication-serving layer was not identical with the repo itself.
- A workspace-local companion surface owned the serving stack and public object
  host behavior.

### Candidate references

- `/REDACTED_HOME/workspace/personal/quine-lab-public-artifacts/README.md`
- `development/open-source/README.md` scope register row for
  `personal/quine-lab-public-artifacts/`
- `development/open-source/README.md` current baseline and DVC-store boundary

## Candidate G: “Not Intended To Be Public” And “Externally Unreachable” Were Different

### Current capture

- Changing repo policy or store intent was not the same thing as proving public
  unreachability.
- CDN/cache residue and older object paths could continue to serve data until
  separately purged and verified.

### Candidate references

- `development/open-source/README.md` current baseline and risk-audit row for
  public CDN cache residue
- `/REDACTED_HOME/workspace/personal/quine-lab-public-artifacts/README.md`
  current note about post-split emptiness, cache purge, and 404 smoke checks

## Candidate H: Default-Open Unfinished Research Made Cleanup Harder

### Current capture

- The project increasingly treated unfinished drafts, notes, experiments, and
  partial results as part of the intended public object.
- That makes the release more honest, but it also sharply increases the burden
  of cleanup, scanning, and boundary-setting.

### Candidate references

- `development/open-source/README.md` release principle and scope register
- `README.md` public landing page warning that some directories are research
  notes or active drafts / working material
- `Paper/README.md` claim-boundary and paper-program routing surfaces
