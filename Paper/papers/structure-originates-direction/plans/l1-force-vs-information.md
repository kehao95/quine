# L1 — force ≠ information (demoted, optional); originates ≠ amplifies (retired)

> Second update, same day: the "originates ≠ amplifies" half of this layer is
> also retired — claim Revision note 3 replaced the origination claim with
> the counterfactual "environment is a difference-maker" formulation, which
> L0's matched pairs already prove. Nothing in this file is a required gate.

**Superseded 2026-07-05 by a claim rescoping** (see
[`claim.md`](claim.md)'s Revision note): "force ≠ information" —
structure shown to move behavior *isolated from content*, token multiset
held fixed — was scope creep beyond what this paper actually argues. The
claim is about **environment** (content + structure + causal arrangement
together) originating direction, not structure-apart-from-content. L0
(E1/E4) already anchors that claim. This layer is now **optional
future-strengthening**, not a required gate — revisit only if reviewer
feedback on a draft specifically demands the stronger, separable claim. The
retraction record below stays as-is: if this is ever revisited, both binding
constraints still apply.

**Status:** RETRACTED (2026-07-05). A first attempt
(`elic-34-force-vs-information`, n=11/cell, GLM-5.2 + GPT-5.5-codex-medium)
was built, run at decisive n, and then **deleted in full, including all run
records**, after tape-reading (qualitative, not just the quantitative score)
surfaced two design/measurement flaws serious enough to invalidate the
numbers rather than merely caveat them. Both are recorded below as binding
constraints on any redesign — do not repeat them.

**Why the first attempt doesn't count, and what must change:**

1. **The arrangement manipulation was not semantically inert.** The nested
   location was named `cache/state/` — a directory name that itself supplies
   a "version history / cached state" narrative. Both model families
   spontaneously read the two buried binaries as *prior versions of the
   foregrounded one* ("this looks like a generation/version system"; "maybe
   the binary self-modifies and writes itself to the cache"), not as three
   parallel, equal-status affordances. This means position and an *inferred
   content narrative* were confounded in the same direction, which could
   inflate rather than cleanly isolate an arrangement-only effect — the
   opposite failure mode from the one the design was guarding against
   (token-invariance), but just as fatal to the claim "irreducible to
   information." **Binding constraint for redesign:** every location in the
   fixture (foreground and buried) must use naming with zero systems/
   versioning vocabulary — e.g. purely opaque hex-token directories at every
   level (matching `_lib/neutral_names.sh`'s convention already used for the
   binaries themselves, extended to the *paths*, not just the binary names).
   More deeply, this is evidence for a real risk the design must budget for:
   a sufficiently capable agent may spontaneously narrativize *any* spatial
   regularity (depth, order, adjacency) into a semantic story ("front =
   current, back = old") through its own inference, not because the
   experimenter embedded content — meaning "arrangement with zero inferable
   content" may not be achievable in general, only approximable, and the
   fixture should be built to make the *specific* recency/importance
   narrative as hard to construct as possible (see redesign note below).

2. **The DV (final-workspace diff) undercounts real invocations when the
   agent cleans up after itself before the run ends.** Confirmed directly in
   two tapes: an agent narrates *"I briefly ran the executables... I then
   cleaned up the generated output directories, leaving the workspace as it
   originally was"* — the binary WAS invoked, but the marker file it wrote
   was deleted before the final snapshot, so the run scored as "not
   invoked." This is a measurement artifact, not a null result, and it is
   not obviously random: runs that exit cleanly on their own have time to
   clean up (erasing the evidence); runs killed by the wallclock cutoff have
   their state frozen mid-task (evidence preserved) — so any two conditions
   or model families with different typical run lengths / termination modes
   would show spuriously different invocation rates for reasons that have
   nothing to do with the experimental manipulation. **Binding constraint
   for redesign:** the scorer must detect invocation from the retained tape's
   tool-call/tool-result log (did the binary's own write syscall fire at
   *any point during the run*), not only from the final-state diff. The
   harness already has infrastructure built for exactly this —
   `_lib/real_gated_attempt.py`'s `gated_attempt()`/`real_exec()` detects a
   real invocation from tape entries, never narration — the next attempt's
   `score.py` should use this instead of (or in addition to) a bare
   final-snapshot diff.

**Open strategic question, not yet resolved:** whether L1 as originally
scoped (prove arrangement-only pull with a fixture that supplies *zero*
inferable content) is achievable at all for a sufficiently capable
interpreter, or whether the claim needs to be reframed around what a
redesigned fixture can actually isolate. See the parent conversation's
reasoning before resuming implementation — this file intentionally records
the failure mode rather than a new design, pending that decision.

## Design 1 — force ≠ information (token-multiset-preserving structure permutation)

**Proves:** the effect is irreducible to token content — a structure
permutation that preserves the token multiset still moves behavior.
**Falsifies the inversion if:** behavior is invariant under structure-scramble
at equated tokens → "structure = tokens" → this is just context engineering;
**stop, do not proceed to L2.**

- **Arms:** condition A vs condition B have an **identical token multiset**
  (same information reachable by the agent) but a **different arrangement** —
  vary adjacency, nesting, foregrounding, and/or path-to-reach while holding
  the token bag fixed.
- **IV:** structural arrangement (permutation class).
- **DV:** target act — an un-fakeable POSIX filesystem/exec event, matching
  the existing `experiments/active/structural-elicitation` DV convention
  (frozen weights, fixed missionless prime, matched pairs).
- **n:** ≥10/cell × ≥3 model families.
- **Decision rule:** divergence beyond control-arm noise ⇒ proceed to L2;
  token-invariance ⇒ kill the inversion — the honest failure mode is that this
  dossier's headline claim folds back into "context engineering," and that
  should be recorded explicitly if it happens, not softened.

## Design 2 — originates ≠ amplifies

**Proves:** rules out "the goal was already there and structure merely
revealed it."
**Design:** show the *same* missionless agent pulled toward *different*
actions by *different* structures. Reuses the E1/E4 apparatus already present
in `experiments/active/structural-elicitation/`.

## Harness fit — not built (prior attempt deleted, see Status above)

The retracted attempt forked `elic-01-gap-completion`'s matched-pair contrast
machinery and `elic-04-affordance-solicitation`'s opaque-binary DV pattern
into a 4-arm design (three identical opaque marker binaries, one foregrounded
at a time + a flat control). That skeleton is still the right harness fit —
what must change is (a) neutral hex-named paths at every nesting level, not
just the binaries, and (b) a tape-based invocation scorer
(`_lib/real_gated_attempt.py`) instead of a final-state-diff-only scorer. Do
not re-fork until the open strategic question above is resolved.

## Open pre-registration question

What operationally counts as "token-multiset-preserving"? Need a
canonicalization/hashing check that confirms arm A and arm B contain the
identical token bag (e.g. sorted token multiset hash equality) before a run is
accepted as a valid instance of Design 1 — resolve this before implementing,
not after.
