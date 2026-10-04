# Manuscript Plan

## Working Promise

Show, concretely, how a living research repository became publicly releasable by
splitting source history, retained private artifacts, and intentionally
published public evidence into separate boundaries.

## Project Premise

This post is not about open-sourcing a conventional software repository.
Quine is a hybrid project: codebase, research archive, experiment record, raw
data reservoir, and control plane live in one place.

That matters because some things that would normally be treated as disposable
engineering residue are part of the scientific record here.

## Current Mode

This document is currently in ideation / capture mode.

For now:

- do not over-compress the user's framing into polished thesis language
- preserve the user's actual reasoning, chronology, and tensions
- use structure only to keep later drafting possible
- allow partially overlapping points to coexist until the larger arc is clear

## Meta-Layer Framing Inventory

These are not yet chosen as the final thesis. They are different higher-level
ways to think about what "making Quine public" actually was.

Do not force early selection yet. Keep the layers side by side until it is
clear which one should carry the main essay voice.

### A. Public release as making a research workbench public

- this is not just a normal software-repo release
- the object being opened includes code, theory, experiments, logs, control
  surfaces, drafts, and raw evidence
- the question is not only whether the code can be published
- the question is what counts as part of the workbench itself
- it also asks which surfaces are core research reality versus incidental
  residue

### B. Public release as preserving an experimental memory layer

- the project is trying to preserve experimental memory, not just code
- logs and run records are part of that memory layer
- the issue is not only how to store logs, but how to preserve a rereadable
  experiment record for later theory work
- this framing emphasizes retained logs as ongoing read surfaces for humans and
  agents

### C. Public release as rebuilding the sensitive-data threat model

- the early project posture was private-first
- sensitive material could appear not only in explicit secret files but also in
  runtime traces, shell output, inherited host-environment echoes, and other
  preserved evidence surfaces
- this framing treats the release problem as a reconstruction of what the real
  leakage surface actually was
- it asks what had to be audited once the team realized that logs themselves
  could be secret-bearing

### D. Public release as publication-boundary design

- this framing focuses on the difference between retention and publication
- a thing can be retained by the project without being publicly restorable
- DVC is important here not only as storage offload, but as a way to preserve
  repo structure while separating storage, retention, and publication
- the whitelist / allowlist publication model belongs to this frame

### E. Public release as preserving historical authenticity

- once the repo already had more than a thousand commits, the question became
  how to clean history without destroying it
- this framing treats history authenticity as one of the central constraints
- the interesting move is not deleting the embarrassing or risky past, but
  transforming the payload form while preserving the historical commit
  structure
- history replay and historical DVC externalization belong here

### F. Public release as a scale and Git-tractability problem

- experiment records did not only create conceptual and security problems
- they also created size, blob, file-count, and Git-performance problems
- this framing emphasizes the practical weight of large payloads, many files,
  and slow repository operations
- it makes room for the quantitative before/after story around DVC

### G. Public release as redefining the public object

- at one point the imagined public object was close to a public runtime
- later that became hard to sustain because the real project also included
  theory, experiments, tests, papers, and control surfaces
- this framing asks what exactly the public version of Quine was supposed to be
- it leaves open the move from a narrow public-runtime surface toward a broader
  public research workbench

### H. Public release as building a release control plane

- this framing focuses on the fact that open-sourcing stopped being a one-time
  cleanup
- it became a managed process with scans, inventories, allow/deny review,
  history-rewrite pathsets, validation, and publication mechanics
- the release work itself became something like an engineered control plane

### I. Public release as balancing open science and operational safety

- this framing emphasizes the tension between wanting to expose unfinished,
  real research material and wanting to avoid accidental disclosure
- it is not only a question of openness versus secrecy
- it is also about how much truth about the research process can be preserved
  without breaking operational safety boundaries
- this framing can also absorb third-party artifacts, bulky materials, and
  other borderline cases that are not simply "secret" but are still not
  straightforward public Git payloads

### Notes to keep open

- these frames overlap; keep that overlap visible for now
- do not force a single master framing yet
- later drafting may choose one of these as the main voice and let others
  become supporting subsections
- some pain points may later be reorganized under one of these frames instead
  of remaining as separate sections

## Pain-Point Inventory And Raw Capture

> Numbering is provisional during ideation. Keep the user's actual distinctions
> first; collapse or merge only later if the final essay really needs it.

### 1. Experimental logs are not disposable in stochastic LLM research

For an ordinary software project, experiment logs usually do not belong in the
public repository. They are large, noisy, and often operationally private.

That default breaks here. Quine's experiments depend on LLM behavior with real
stochasticity. Even if you preserve prompts, code, models, and environment,
you still do not deterministically recover the exact run trajectory later.
That makes the original logs part of the evidence, not just temporary residue.

More than that, the logs are also a long-lived read surface. When the theory
later advances, we need to go back and reread old runs as raw experimental
data. Sometimes that rereading is done by a human; sometimes the point is to
let an agent read and analyze the preserved trajectory. The logs are therefore
not just archival proof that a run happened. They are reusable source material
for later interpretation.

#### Raw ideation capture from conversation

- one of the biggest differences between this project and a normal software
  project is that Quine is a code / research / experiment hybrid rather than a
  conventional software repo
- the repo contains a lot of raw data, control-plane material, and experiment
  history rather than only implementation
- for a normal project, experiment logs generally should not be public and
  often should not even live in the main repository surface in this way
- for this project, because the object is an LLM and the behavior has real
  randomness, "just preserve the full reproduction conditions" is not enough
- even if the project later knows how to rerun something in principle, that
  does not mean it can regenerate the exact same trajectory
- because of that, the logs have to be preserved as part of the project
- early on, the core idea was simply that logs must be retained
- in a lot of the earlier history, the implementation of that idea was to
  track the logs directly in the repo
- this was not a casual choice; it came from the belief that the logs were too
  important to lose
- there was already a long period of struggle around this question
- once everything was being tracked, some things inevitably got tracked that
  were not actually the experimental record we cared about
- that included generated binaries, outputs produced incidentally, and other
  results that were not actually important enough to preserve in the same way
- so the next step was not "stop retaining logs," but "become more explicit
  about what counts as a log worth retaining"
- this led to work on retention criteria inside the project
- it also led to work on making retained logs more stable at the path level,
  so they would not just exist, but would remain referenceable and readable
- there was also work at the code and experiment level to improve the quality
  and stability of the logs that were being retained
- another important correction was that the project should not blindly absorb
  every single success and every single failure artifact into the main retained
  surface
- even after that improvement, the retained logs were still very numerous
- later there was an attempt to package logs as zip archives
- that did help somewhat with capacity / footprint
- but it made readability substantially worse
- that readability loss matters because the logs are not only something we keep
  for preservation; they are something we need to go back and read
- when the theory later changes or progresses, old logs need to be revisited as
  raw data
- sometimes that rereading is done by a human
- sometimes the point is to let an agent read the old trajectory and analyze it
- so the logs are not merely stored evidence; they are an active read surface
- that means any storage decision that makes them harder to inspect is damaging
  scientific use, not just reader convenience
- in other words, the problem is not merely storage volume; it is how to keep
  an analyzable memory layer for stochastic experiments without turning the repo
  into an unreadable dump of everything emitted

#### Why this became painful

- the early principle was simple: logs must be retained
- for a long stretch of the project's history, the practical implementation was
  to track logs directly in the repo
- retained logs were not only evidence for one result; they were the read layer
  we expected to revisit when new concepts, new theory, or new analysis
  questions appeared
- this preserved evidence, but it also pulled in accidental residue:
  generated binaries, irrelevant outputs, and non-essential failed or noisy
  results
- log value was uneven, but the retention policy was initially too blunt
- even after filtering improved, retained logs were still large and numerous
- a later attempt to zip logs reduced footprint somewhat, but made the record
  much less readable

#### The repair direction

- stop treating "retain everything" as the only safe option
- define explicit selection criteria for which logs deserve retention
- improve the stability of retained-log paths so they stay referenceable
- improve the code and experiment surfaces so the logs worth preserving are
  more stable and less contaminated by incidental outputs
- move toward a policy for retaining important logs rather than blindly
  absorbing every success/failure artifact into the repo

#### The underlying tension

The real problem is not "logs are too big." The real problem is that the
project needs to preserve irrecoverable experimental evidence without letting
the repository collapse into an unreadable dump of everything the system ever
emitted.

And because those logs are meant to be reread later, readability is not a
secondary concern. A storage optimization that makes logs harder for humans or
agents to inspect is not neutral; it directly damages the project's ability to
do later theory-building and re-analysis.

#### Questions to keep open during ideation

- when this section is later written in prose, should the emphasis be on
  randomness / irreproducibility, or more on logs as a reusable read surface
  for later theory work?
- how much of the early "track logs directly in the repo" period should be told
  as a principled choice versus an expedient early implementation?
- should zip packaging appear here as an important failed compromise, or only
  as one example of the deeper readability versus storage tradeoff?
- do we want to distinguish different classes of retained logs:
  canonical / clinical / important / noisy / failed / incidental?

### 2. Sensitive data made branch topology a real problem

This pain point starts from sensitive data, but it does not stay only at the
level of "remove secrets." It turns into a branch-topology and maintenance
problem.

#### Raw ideation capture from conversation

- one early idea was to maintain a `private` branch and a `public` branch
- at first this sounds conceptually clean: private branch keeps the full
  research reality, public branch exposes what is safe
- later the problem became that the `public` branch was relatively hard to
  maintain
- one reason was that the public branch was imagined as mainly a `public
  runtime`
- but the actual repository contains much more than runtime code
- it also contains theory, experiments, tests, and other research surfaces
- once that wider project reality is taken seriously, a narrow `public runtime`
  branch is no longer an easy mirror of the real repo
- instead, it becomes a separately maintained slice
- and that separate maintenance has real pressure
- the pressure is not only technical; it is also conceptual:
  what exactly belongs in public if the project is not only software?
- if the project includes theory, experiments, tests, drafts, and evidence,
  then a public branch that only carries runtime code may be too incomplete to
  honestly represent the project
- but a public branch that carries more of the project inherits much more of
  the sensitive-data and maintenance burden
- this means the branch split is not a simple yes/no privacy mechanism; it is
  entangled with the question of what the public object actually is
- another important part of this is historical posture:
  the project began life operationally as a private project
- because the project was being run as private, and because there was an idea
  that the private branch would continue to hold the full reality, a lot of
  sensitive material was tracked more freely than it would have been in a
  public-first repo
- that included sensitive data such as keys
- later, some of that material was deleted or removed, but deletion at tip is
  only one part of the problem
- another, harder problem is that the agent can call shell commands
- when the agent inherits the host environment, shell output can directly print
  sensitive keys or other secret-bearing environment values
- that means the leakage surface is not only config files or committed secrets
- it also means the runtime traces themselves can become sensitive
- and because this project preserves logs heavily, every retained log can
  potentially become an information-leak source
- that creates a very sharp tension:
  - on one side, logs need to be preserved because they are evidence and later
    read surfaces
  - on the other side, those same logs can contain accidental shell output that
    leaks secrets from the host environment
- this made public release much harder, because the project could not assume
  that "secrets live in secret files" was the full threat model
- instead, the threat model had to include the possibility that ordinary agent
  operation had already copied secret material into logs, shell output, and
  other traces

#### What seems important here

- the problem began as a sensitive-data problem
- but it became a repo-topology and public-scope-definition problem
- it also became a logging and trace-sanitization problem
- maintaining a separate public branch was hard not only because of scrubbing
  work, but because the real repo was not just a runtime repo
- the broader the public object became, the more maintenance pressure the
  public branch absorbed
- the narrower the public object became, the less honestly it represented the
  actual project
- the project also had to face the fact that sensitive data could have leaked
  into the very evidence surfaces it most wanted to preserve

#### Questions to keep open during ideation

- should this remain its own pain point, or later merge with the history /
  public-branch section?
- how much of this section should emphasize `sensitive data` versus `public
  branch maintenance pressure`?
- do we want to tell this as a shift from `public runtime branch` to `public
  research workbench`, or is that already too compressed for this stage?
- should the shell-output / inherited-host-env leakage story stay inside this
  pain point, or later become its own subsection because it cuts directly
  against the log-retention story?

### 3. Once DVC was chosen, history replay became the real cleaning problem

This point is closely related to branch topology and sensitive data, but it has
its own shape in the ideation so it is worth recording separately for now.

#### Raw ideation capture from conversation

- after the project selected the DVC direction, another very large pain point
  appeared
- the repository already had more than one thousand commits of history
- that history contained a large amount of experimental record material
- it also potentially contained material that could leak sensitive information
- one obvious cleaning option would have been to simply delete all of that
  problematic history
- doing that could achieve a clean result in the narrow technical sense
- but it would destroy the authenticity of the repository history
- preserving the authenticity of the history was important
- the goal was not to end up with a fake clean repo that had lost the real
  shape of how the work evolved
- so the cleaning approach tried to preserve all the commits
- the thing that was changed was not the existence of the commits, but the
  payloads that lived inside risky experimental-record portions of that history
- the approach was effectively a kind of history replay
- through that replay, historical experimental-record surfaces were converted
  into DVC-backed form
- this made it possible to keep the historical commit structure while moving
  risky or heavy experimental payloads out of public Git history
- this is important because it means the project was not choosing between:
  - preserve the real history
  - or clean the repo
- instead it was trying to preserve the real history in a transformed storage
  form
- the detailed operational story for this step likely lives partly in the
  Codex sessions and partly in the repo-local release-control surfaces

#### What seems important here

- this is not ordinary secret deletion
- it is a historical-authenticity problem
- the repo had enough history that a naive clean slate would have changed the
  meaning of the project
- the project also suffered from scale pressure:
  once logs and experimental records were heavily tracked, both size and file
  count started to matter for Git performance
- this means the issue was not only secrecy or authenticity, but also practical
  repository operability
- it would be useful later in drafting to include concrete examples of the big
  blobs that were historically tracked and the before/after size differences
- DVC made a new type of cleanup possible:
  retain the historical structure while externalizing the risky experimental
  payloads
- this may be one of the strongest examples in the essay of using DVC not only
  for current-state storage, but for historical preservation under a new
  boundary

#### Questions to keep open during ideation

- should this remain separate from the broader `history needed work` pain point,
  or later become its core subsection?
- how much emphasis should go on the number of commits and scale of history
  versus the principle of historical authenticity?
- when we later draft this section, do we want to explain the mechanism as
  `history replay`, `history rewrite`, or something more precise that better
  matches what was actually done?
- how much of the detailed operational procedure should be reconstructed from
  repo files versus from Codex session history?

## Cross-Cutting Solution Notes

These are not yet finalized as essay sections. They are mechanisms or ideas
that may need to appear across multiple pain points.

### DVC as a structural answer to repo/data tension

#### Raw ideation capture from conversation

- another important thread is that `DVC` solved a large portion of the
  problems the project was running into
- the important thing about DVC here is not just "large files go elsewhere"
- more specifically, it allowed data to be uploaded outside the Git repository
- at the same time, the project could still preserve the data's original place
  inside the repo's structure
- and it could preserve the sense that the data still changed together with the
  project, rather than becoming an unrelated external bucket
- this is important because the project did not want to simply delete data from
  the repo's conceptual structure
- it wanted to move the heavy payload outside Git while still keeping the
  project's internal organization, addressability, and evolution legible
- that means DVC may need to be explained not just as a storage optimization
  tool, but as a way to preserve structure and change across a repo/data split
- this likely connects to multiple earlier pain points:
  - log retention
  - large evidence bundles
  - private retention versus public publication
  - branch maintenance pressure
- DVC also had an unexpected additional benefit because the DVC server was
  separately controlled by the project owner
- since that server could live on a personally controlled machine/server, the
  project could decide exactly how DVC publication should work
- that meant public release did not have to mean "all DVC-tracked data is now
  exposed"
- instead, the project could maintain a whitelist of experiment records that
  had already been reviewed and judged safe to publish
- only those reviewed experiment records would be made publicly restorable
- this was important because it made DVC not only a storage split, but also a
  publication-control mechanism
- that helped preserve the spirit of open source without accidentally exposing
  sensitive information
- in other words, DVC created a path where the project could still preserve and
  retain a broad experiment memory privately, while exposing only the audited
  subset publicly
- it may be useful in the essay to explain that DVC did not magically solve
  everything, but it changed the shape of the problem
- before DVC, the repo had to choose too often between:
  - keep the data in Git and absorb the weight directly
  - remove the data and lose some of the structure and continuity
- DVC created a middle position:
  preserve the structure and project-local meaning of the data while moving the
  heavy payload out of Git proper

#### Questions to keep open during ideation

- should DVC be introduced early as a major enabling tool, or only after the
  concrete pain points make the need obvious?
- is DVC best framed here as:
  - data externalization
  - evidence retention infrastructure
  - structure-preserving storage split
  - or a bridge between repo semantics and out-of-repo payloads?
- should the personal-server / whitelist-publication angle be presented as an
  accidental operational advantage, or as one of the core reasons DVC became a
  viable release mechanism here?
- do we want to distinguish clearly between:
  - DVC as private retention infrastructure
  - DVC as public publication mechanism
  since those later became importantly different?

## Candidate Additional Decisions And Implementation Threads

These are candidate important points that may not all end up in the final
essay. Record them now so they are available later.

### The public object changed from `public runtime` to `public research workbench`

#### Raw ideation capture

- one important decision thread is that the project did not stay at the level
  of “how do we publish the runtime?”
- over time the public object became larger
- the project increasingly treated Quine not only as a runtime, but as a
  research workbench
- this means theory, experiments, tests, drafts, notes, and control-plane
  material are not automatically outside the public object
- this is important because it changes the release problem from
  “publish the clean software subset” to
  “publish the workbench, and then carve out only what is unsafe”
- this also seems to be one of the reasons a separate `public runtime branch`
  became less satisfactory over time

### Private retention and public publication became explicitly different facts

#### Raw ideation capture

- one major decision was to hard-split private retention from public
  publication
- this was not only a tooling distinction; it became part of the release model
- something could remain tracked, retained, and recoverable inside project
  infrastructure without being publicly restorable
- this matters because the project did not want “DVC-tracked” or “kept in
  history” to imply “safe for public readers to restore”
- the allowlist / manifest / whitelist logic grows out of this decision

### HEAD cleanup and history rewrite became two different stages

#### Raw ideation capture

- another important implementation distinction is that cleaning the current tip
  and cleaning the full history were not the same job
- the repo could reach a state where the current tree looked much cleaner
- but that still did not mean the actual history was publicly safe
- this matters because the release work was not one monolithic cleanup
- it became a staged process:
  - make current HEAD cleaner
  - separately handle old history

### Open-source release itself became a maintained control plane

#### Raw ideation capture

- another implementation-level decision is that open source itself stopped
  being an informal cleanup effort
- it became its own maintained object
- it had its own owner surface, scans, risk gate, allow/deny policy, rewrite
  pathsets, and validation notes
- this is important because release work was too large and too stateful to
  survive only in chat memory or temporary command history

### Public-facing truth surfaces had to be rewritten

#### Raw ideation capture

- another important thread is that the public release was not only about what
  data or history stayed in the repo
- the public-facing surfaces also had to tell the truth
- that includes `README`, `CITATION`, `QUICKSTART`, and the Paper routing
  surfaces
- the repo needed to explain honestly:
  - Quine as runtime
  - Quine as research workbench
  - what is draft / working material
  - what public artifact restore actually means
  - how a public reader should cite the project

### Paper, third-party, and bulky artifacts needed their own policy

#### Raw ideation capture

- not every difficult artifact was difficult because it contained a secret
- another class of material was hard because it was bulky, generated,
  third-party, or collaboration-adjacent
- examples include PDFs, presentations, graphs, legacy outputs, and other build
  artifacts
- those are important because they show that “not suitable for public git
  history” is a broader category than “secret”
- this likely needs to remain visible as a separate policy thread

### The public artifact serving layer lived outside the repo

#### Raw ideation capture

- another implementation detail that may matter is that the public artifact
  serving layer was not identical with the repo itself
- there was a workspace-local companion surface for the serving layer
- this seems important because it let publication control happen outside the
  repo while still remaining tied to repo-local manifests and policy

### “No longer intended to be public” and “externally unreachable” were different

#### Raw ideation capture

- another subtle but important thread is that changing the intended public
  boundary did not automatically make data unreachable from the outside world
- cache / CDN residue and old object paths could preserve accidental access
  even after local/public intent changed
- this matters because it distinguishes repo policy from actual external
  availability

### Default-open unfinished research increased the difficulty of cleanup

#### Raw ideation capture

- one more tension worth preserving is that the project moved toward being
  default-open even for unfinished drafts, research notes, experiments, and
  partial results
- that makes the release more honest and more aligned with the workbench idea
- but it also raises the difficulty of cleanup substantially
- if only polished runtime code had to be public, many of these problems would
  have been easier
- once unfinished research is part of the intended public object, the cleaning,
  auditing, and boundary-setting burden becomes much heavier

## Narrative Shape

The post should read as a control-plane repair story, not a cleanup checklist.
The rhythm is:

1. name the release paradox
2. surface the real pain points
3. show the boundary split that resolved each one
4. end with a reusable method for other research workbenches

## Section Skeleton

### 1. The paradox: this was never just a software release

- open with the difference between "repo contains code" and "repo is safe to
  publish as a living research object"
- introduce Quine as a hybrid workbench containing code, drafts, experiment
  logs, retained runs, DVC pointers, notes, raw material, and historical
  residue
- set up the mistake this post will reject: "open source = delete secrets and
  push"

### 2. First pain point: experimental logs were evidence, not residue

- explain why ordinary software norms would push logs out of the public repo
- explain why stochastic LLM experiments break that assumption
- explain that logs are not only retained evidence; they are a durable read
  surface for future interpretation and theory updates
- show the first naive safe move: keep everything
- show why that turns into repo bloat and evidence contamination
- show the next repair step: retained-log criteria, path stability, and
  higher-quality preserved logs
- include the failed/partial compromise: zipping logs saves some space but
  hurts readability

### 3. The real object: three boundaries had been entangled

- public source and history
- retained but private artifacts
- selectively published public evidence
- explain why conflating any two creates pain later

### 4. Pain point: retained artifacts were mixed with publication semantics

- raw problem: DVC pointers and retained run artifacts looked too close to
  public restore
- release move: separate private retention from public publication
- evidence:
  - `development/open-source/README.md`
  - `.dvc/public-manifest.txt`
  - `personal/quine-lab-public-artifacts/README.md`

### 5. Pain point: Git history itself was unsafe to expose

- raw problem: old `.env*`, provider profiles, raw logs, local paths, and bulky
  retained material existed in history
- release move: cut a sanitized public branch instead of pretending `lab` was
  already publishable
- evidence:
  - `development/open-source/public-branch-cut-2026-05-26.md`
  - relevant rewrite commits and audit notes

### 6. Pain point: the release had no owner surface at first

- raw problem: scans, blockers, risk inventory, and policy decisions were
  scattered across diffs and chat
- release move: create a dedicated open-source control plane
- evidence:
  - `development/open-source/README.md`
  - `evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md`

### 7. Pain point: public evidence needed precision, not maximum disclosure

- raw problem: "make all artifacts public" would overexpose residue and create
  restore confusion
- release move: manifest-scoped allowlist and separate public-serving surface
- evidence:
  - `.dvc/public-manifest.txt`
  - `personal/quine-lab-public-artifacts/README.md`

### 8. Pain point: public readers needed a truthful landing surface

- raw problem: README, citation, and public routing can misrepresent what is
  stable, what is draft, and what each paper claims
- release move: recenter README/citation/public artifact docs around honest
  public entrypoints
- evidence:
  - root `README.md`
  - `CITATION.cff`
  - public artifact surface

### 9. What this taught me about open-sourcing research repos

- separate source/history, private retention, and public evidence early
- treat release prep as a control-plane project with explicit owners
- preserve useful research residue privately instead of forcing false deletion
- publish with precise allowlists, not vague intent

## Commit Hunt

Need a final shortlist grouped by role:

- artifact-retention boundary commits
- open-source owner-surface commits
- public DVC/public artifact boundary commits
- history rewrite / public branch commits
- public landing / citation commits

## Drafting Rule

Prefer "pain point -> why naive release fails -> bounded fix -> what changed"
over pure chronology. Chronology can support the story, but should not become
the story.
