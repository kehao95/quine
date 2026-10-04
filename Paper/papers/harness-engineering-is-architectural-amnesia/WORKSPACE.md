# Workspace: Harness Engineering Is Architectural Amnesia

> Session workspace for drafting "Harness Engineering Is Architectural Amnesia"

## Current Focus

Tightening the promoted manuscript as the single canonical prose surface.

Release strategy now lives in [`./plans/release-plan.md`](plans/release-plan.md).

## Active Manuscript Surface

| Surface | Status | Location |
|---------|--------|----------|
| manuscript | active | `manuscript.md` |
| release plan | active | `plans/release-plan.md` |
| X thread draft | active | `plans/x-thread-draft.md` |

## Key Assets

- **The Isomorphism Table:** Agent concept ↔ UNIX primitive ↔ Framework reinvention
- **The Quine snippet:** `./quine "..." > q && ./q "..."`
- **Potential origin story hook:** Compressed "discovery of Quine" narrative (optional)
- **Comparative external source:** Slate as an independent engineering rediscovery of runtime problems from a different implementation path; preserved at [`../../essays/_meta/external-sources/randomlabs-slate-2026-03-09.md`](../../essays/_meta/external-sources/randomlabs-slate-2026-03-09.md)
- **Reference triage:** ranked shortlist of which Slate references actually matter for this manifesto at [`../../essays/_meta/external-sources/slate-references/TRIAGE-harness-engineering.md`](../../essays/_meta/external-sources/slate-references/TRIAGE-harness-engineering.md)
- **Plans surface:** `plans/` now holds the release plan and X-thread draft

## Session Log

- 2026-03-13: Created dossier structure (README, manuscript-plan, WORKSPACE). Framework captured.
- 2026-03-13: Wrote first full draft in `sections/01-draft.md`.
- 2026-03-13: Deleted the draft to return to outline-first iteration.
- 2026-03-13: Wrote a new full draft after tightening the manifesto outline.
- 2026-03-13: Promoted the draft into `manuscript.md` and added the public GitHub link.
- 2026-03-13: Collapsed the blog workspace onto a single-manuscript control plane.
- 2026-03-13: Split launch/distribution strategy into `plans/release-plan.md`.
- 2026-03-13: Recorded current companion-post recommendation: follow `harness-engineering-is-architectural-amnesia` with an implementation walkthrough framed as `quine-as-posix`, while keeping `the-autopoietic-repository` as a parallel explanatory lane.
- 2026-03-16: Added Random Labs' Slate report as external source material because it validates the runtime-problem diagnosis from a different implementation direction, while still leaving room for the essay's core accusation that the industry is building the substrate at the wrong layer.
- 2026-03-16: Triaged the 25 Slate references and identified the high-value subset for this manifesto: Terminus, Devin/Cognition, Manus, Amp handoff, and Context Rot.
- 2026-03-16: Broadened the steelman definition of harness engineering using LangChain's official harness docs and blog so the dossier now treats harnesses as reliability, control, observability, and governance stacks rather than only memory-management or subagent machinery.
- 2026-03-16: Added Claude Code and Codex as productized harness witnesses so the comparative landscape no longer over-indexes on open-framework examples alone.

## Open Decisions

- [ ] Include origin story as cold open or defer to companion post?
- [ ] Keep the McIlroy paragraph or cut it for speed?
- [ ] Decide how much to soften the LinkedIn version after X feedback.
- [ ] Name full series arc or just tease?
