## Period XVI: DVC Externalization and Curated Public Release (May 11–27)

Experiment retention became a repo-level DVC protocol, public artifact restore
was separated from private lab retention through a manifest-scoped publication
boundary, and open-source release moved from scattered cleanup pressure to an
audited history-rewrite and curated-publication control plane.

### May 11–12 — Experiment retention becomes a DVC protocol
**Key commits:** `011cc8e39`, `bcdefea6b`, `fd52ce5e2`
**What changed:** Initialized repository DVC artifact retention, promoted
retention guidance into `DEVELOPMENT.md` and `experiments/PROTOCOL.md`, and
migrated broad experiment run families from checked-in payload trees to
git-visible `.dvc` pointers.
**Friction:** Retained run artifacts were growing across experiments, but the
storage rule still lived mostly as local habit. Git history, restore posture,
and evidence curation risked drifting apart lane by lane.
**Structural repair:** Made artifact retention a repo-level contract: keep
control-plane guidance in git, route bulky retained payloads through `.dvc`
pointers, and leave experiment-local README prose only for true local restore
exceptions or curation boundaries.
**What it taught:** Evidence retention needs one canonical policy layer before
individual experiment surfaces can stay small and honest.

---

### May 12–13 — Public artifact restore becomes a separate surface
**Key commits:** `56940665e`, `1961eccb5`
**What changed:** Recorded the public artifact surface in root docs and added a
public DVC remote to shared config so public restore could be reasoned about as
its own repository boundary rather than as a side effect of retaining evidence.
**Friction:** Once DVC pointers existed across experiments, "tracked in repo"
and "restorable by public readers" were no longer the same fact. Without an
explicit boundary, a pointer could be misread as publication intent.
**Structural repair:** Distinguished retained evidence from published restore:
repo-visible `.dvc` pointers can exist without implying that their payloads are
available on the public artifact surface.
**What it taught:** A checked-in pointer is a retention statement; publication
needs an additional restore contract.

---

### May 24–25 — Open-source release becomes an owner surface
**Key commits:** `8df17dfa5`, `34fff5820`, `dfeda27fe`
**What changed:** Created `development/open-source/` as the release-boundary
register, added a repo-local risk gate for public release, sanitized open
release surfaces, and defined a canonical first-pass history-rewrite pathset.
**Friction:** Open-source release work had become a growing cluster of scans,
blockers, and ad hoc cleanup, but it lacked one owner surface that could hold
the current boundary decisions and the exact history-rewrite drop set.
**Structural repair:** Turned release preparation into a maintained control
plane with a narrow owner: the repo now has an addressable place for public /
private boundary judgments instead of spreading them across chat memory,
one-off shell probes, and transient commit context.
**What it taught:** Publication work becomes safer once the boundary decisions
themselves have a canonical owner surface.

---

### May 25–26 — Public DVC publication is isolated from history rewrite
**Key commits:** `c86ae6e33`, `8f3b45683`, `08d75b676`, `8b894997b`, `6012994f7`, `c4e581281`
**What changed:** Split the public artifact store from the private write store,
introduced tracked `.dvc/public-manifest.txt`, added maintainer-side sync /
rebuild helpers plus reader-side manifest restore, validated clean public DVC
restore after purge, and recorded Paper-specific DVC and history-rewrite
constraints.
**Friction:** History rewrite, DVC retention, and public artifact publication
were close enough to be confused with one another. A broad push, stale edge
cache, or ambiguous restore path could over-publish private evidence while
still looking operationally normal.
**Structural repair:** Separated three jobs that had been collapsing together:
curating public git history, retaining private lab artifacts, and promoting an
exact allowlisted set of public DVC payloads. Public restore is now
manifest-scoped rather than a blanket `dvc pull` claim.
**What it taught:** Public evidence release needs an allowlist and a restore
protocol, not just a second remote.

---

### May 26–27 — Curated public branch semantics and README boundary
**Key commits:** `6f09526e7`, `1d7156e95`, `bd423f7b7`, `24e15410b`
**What changed:** Recorded the clean public branch cut, clarified refresh
semantics for the curated branch, improved public README entrypoints, and
re-centered root README / control-plane boundaries so maintainer routing lives
in `AGENTS.md` and `DEVELOPMENT.md` rather than leaking into the public landing
page.
**Friction:** After the public cut existed, root docs still needed to avoid two
wrong readings: that `main` mirrored `lab`, and that the human public README
should also carry maintainer control-plane posture.
**Structural repair:** Locked the public branch in as a curated release surface
and restored the README to a pointer-first public entrypoint, with repo-operating
rules routed back to the maintainer control plane.
**What it taught:** Open-sourcing the repository also means publishing the
right document boundaries, not only the right files.
