---
id: sensor-binding-severance
type: failure-mode
status: deposited
applied: false
class: sensing-layer
symptom: a registered blocking validator goes silent when a restructure moves the surfaces it guards; its filter still matches something by accident, so every audit stays green while nothing is checked
evidence: [482675a52, 66e2d3773, 17c2650e9, 5c11cf31d, d4b5c96d8, c762154d0, 5026dc081]
candidate-solution: have each validator publish its own scan set and set-difference it against the hook filter that feeds it — a two-set comparison, not a heuristic judgement of whether a regex "looks right"
owner: lab
---

# sensor-binding-severance

## Symptom

A validator is correct. Its registration is correct. The surface it guards
exists. And it checks nothing, because the **binding** between the two — the
hook's `files` filter — still encodes a layout that a restructure moved.

Neither end changed. The middle was severed, silently, by a commit that touched
neither the validator nor the guarded surface.

The failure is invisible to inspection because the severed filter usually still
matches *something*. Any audit that asks "does this pattern match anything?"
gets a non-zero answer and reads green.

## Evidence

Six blocking pre-commit hooks in [`../../prek.toml`](../../prek.toml) shared two
`files` filters that encoded the pre-restructure layout. Discovered 2026-07-26
by a missionless probe (receipt: `5026dc081`, cycle 43); verified independently
by the integrator against `git show c521013e5:prek.toml` and
`git ls-tree -r c521013e5` rather than the working tree.

| filter | matched | surfaces it names |
|---|---:|---:|
| `^experiments/p[^/]+/README\.md$\|^experiments/p[^/]+/[^/]+/README\.md$` | **1** | 172 experiment READMEs |
| `^Paper/[^/]+/papers/[^/]+/README\.md$` | **1** | 38 dossier READMEs |

Hooks affected: `check-experiment-paper-feeds`, `check-paper-dossier-links`,
`check-paper-dossier-links-from-experiments`, `check-paper-evidence-backbone`,
`check-paper-evidence-backbone-from-experiments`,
`check-control-plane-frontmatter`.

**Both matches were coincidences.** `experiments/protocol/` happens to begin with
`p`; one essay dossier happens to sit a directory deeper than the rest. Each
filter therefore returned exactly one file — non-zero, and green.

**Blackout window.** `482675a52` (2026-06-21, retire publication lens roots) and
`66e2d3773` (2026-06-22, remove the legacy phase layer) → `c762154d0`
(2026-07-26). Roughly five weeks in which six blocking checks ran on every
commit and verified nothing.

**Drift that got through**, from the repair diff: one dossier with a malformed
reverse-link mapping and prose where its evidence backbone belonged, three
theory objects not naming the dossier declaring them, plus index skew. Modest
for five weeks across 210 surfaces — most of the visible link rot was caught by
a *different*, still-live validator. The counterfactual rejection set of these
six over the window is close to one dossier section, which is itself part of the
finding: a check can be inert for weeks without dramatic consequence, which is
exactly why nothing surfaced it.

**Two aggravating facts.**

- `check-validator-registration.sh` — the meta-sensor built for this class —
  reported `scripts=80 prek_entries=48 issues=0` throughout. It measures
  **registration completeness, not reachability**.
- `prek.toml` was edited on 2026-07-18 and 2026-07-19 (`17c2650e9`,
  `5c11cf31d`, `d4b5c96d8`) *during the proxy-enforcement census* — an audit
  whose explicit subject was blocking checks that manufacture false confidence.
  It deleted two experiment hooks and rebound a third, and passed these six.

And the rule was already written. [`../../scripts/README.md`](../../scripts/README.md)
§ Sensor Typology states that a validator "triggered against a file that does not
exist is inert, and inertness is a silent-failure state, not a pass." The mode
was **named, given a meta-sensor, and walked past during a census devoted to
checks** — and six instances still lived.

## Root cause

**The typology presupposes the binding.** The four sensor types — presence,
compliance, sync, inventory — all answer questions *about a population the
sensor is assumed to reach*. None of them observes whether the population is
reached at all. Reachability is prior to every type, and is watched by none.

**Why the accidental match matters.** A cleanly-zero filter would be visible: a
hook that matches no files is a plausible thing to audit for. What makes this
class durable is that a severed filter usually retains a coincidental match, so
the cheapest available check (`non-empty?`) returns the wrong answer. The
detection cost is not "is this set empty" but "does this set correspond to the
surfaces the hook's name claims" — a comparison against intent, not a property
of the set.

**Why the census could not see it.** `proxy-enforcement-collapse`'s detection
heuristic is: *can the pass condition be satisfied by editing only the checked
text?* That is a question about a criterion's **semantics**. A validator that
never runs has no pass condition in play at all — it is not a proxy, it is
absent. The two failure modes are orthogonal, and an audit tuned to one is blind
to the other by construction.

### The dual of proxy-enforcement-collapse

|  | [`proxy-enforcement-collapse`](./proxy-enforcement-collapse.md) | this object |
|---|---|---|
| what is wrong | the criterion measures a proxy for the latent property | the criterion is fine; it reaches nothing |
| population effect | converges on the cheapest passing shape | drifts freely, unobserved |
| "green" means | the tokens are present | nothing was checked |
| detection question | can the pass condition be met by editing only the checked text? | does the hook's file set correspond to the surfaces it names? |
| felt symptom | uniformity ("AI-generated" tone) | none — that is the problem |

Together they bracket the sensing layer: one check watches the wrong property,
the other watches nothing. Neither is visible to the other's audit.

## Candidate solution (deposited — applied: false)

**Direction, not a build.** Have each validator publish the set of surfaces it
believes it guards, and set-difference that against the file set the hook filter
actually feeds it. A discrepancy is then a **fact about two computed sets**, not
a judgement about whether a regex looks correct — which keeps it on the
ground-truth side of `proxy-enforcement-collapse`'s requirement test.

The instance-level repair already exists: `c762154d0` rebinds the two filters
(1 → 170 and 1 → 95 matched files) and carries a negative test — deleting one
Experiment Feeds line now fails pre-commit, where before it was silently
skipped.

## Why not applied yet

Two different reasons, and they should not be collapsed.

1. **The class-level fix is deliberately withheld.** The discovering probe
   declined to ship a heuristic reachability sensor, and named why: the only
   cheaply available ground truth was the one that had just been proven
   foolable. `proxy-enforcement-collapse`'s meta-root warns that any repair of
   that mode which adds watcher structures re-instantiates the mode one level
   up — and a reachability *checker* is exactly a watcher over watchers. The
   set-difference form above may escape that objection because it compares two
   computed sets rather than inspecting a pattern, but it has not been argued
   through, and shipping it unargued is the reflex the sibling object exists to
   catch.
2. **The instance fix is merely unmerged.** `c762154d0` sits on the branch
   `worktree-agent-a14ecf90ec0e19179`, deferred by human decision while
   concurrent sessions held `prek.toml` and two validator scripts. This is a
   scheduling state, not a design position. **Until it merges, the six hooks on
   `main` still run green on every commit while checking nothing** — recorded in
   [`../../development/status/README.md`](../../development/status/README.md)
   § Blockers.

An agent attempting this must separate the two: merging the branch closes the
instance and changes nothing about the class.

## Theory links

- [`proxy-enforcement-collapse.md`](./proxy-enforcement-collapse.md) — the dual;
  see the table above. Its repair grammar (delete > rebind > amend, never a
  meta-check) applies here, and `c762154d0` is a **rebind**.
- [`validation-attraction.md`](./validation-attraction.md) — the run-side sibling.
  A relevant qualification it supplies: the probe that found this *did* open by
  running the repo-wide gate, which is that object's documented attractor. The
  difference is that the gate came back red and was traced to a defect. The
  failure mode is validation without a drift hypothesis, not validation.
- [`../../Paper/theory/objects/transmissible-pattern.md`](../../Paper/theory/objects/transmissible-pattern.md)
  — the theory-side generalization deposited from the same session: a binding is
  the cheapest of {criterion, surface, binding} to describe and the most
  expensive to verify, which makes it the **structurally preferred hiding place
  in any system that enforces through indirection**. Registered as a `SEED` on
  [`../../Paper/theory/objects/narrative-claim.md`](../../Paper/theory/objects/narrative-claim.md).
- [`../patterns.md`](../patterns.md) — reverse-link row; nearest neighbours are
  *Externalize-without-reach (dormant rule)* (a fix that never acquires a
  trigger) and *Declared ownership without compliance pressure*.
- [`../../scripts/README.md`](../../scripts/README.md) § Sensor Typology — owns
  the four types and the inertness rule this object shows to be unenforced.
