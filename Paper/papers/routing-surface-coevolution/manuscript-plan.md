# Manuscript Plan: OS Routing Surfaces as a Red Queen Niche

## Argument plan

The paper argues one positive claim and reports one honest negative boundary:

1. **Positive**: the OS routing substrate is a Red Queen niche — the medium of
   attack and the medium of survival are the same object, so substrate topology
   is an experimental parameter of the coevolutionary dynamics.
2. **Diagnostic**: backend transfer (legacy FIFO → FUSE) is the cleanest evidence
   that topology parameterizes the attack class.
3. **Negative boundary**: the symmetry trap — an asymmetric winner collapses
   into mutual stalemate when copied — bounds the open-ended-escalation reading
   and is itself an ALife phenomenon.

The argument arc: simulated arenas → relocate to OS → observe real tactics →
show topology dependence → report the symmetry boundary → bounded claim.

## Section skeleton

### Problem
- ALife coevolution tradition (Tierra, Avida, competitive-coevolution models)
  presupposes a simulated medium.
- State the relocation: coevolution medium = the OS itself.
- Pose the question: is substrate topology a parameter of the dynamics?

### Setup
- `drq-03` host-lineage regime: two Quine processes, shared host ecology.
- Three load-bearing substrate levers: self-source projection, lineage
  inheritance, selectable public control-surface backend (legacy FIFO vs FUSE).
- No designed action set; tactics must be discovered via ordinary affordances.

### Results
- Phenomenon 1: FIFO-drain strangulation (legacy backend) — functional death
  without `kill`.
- Phenomenon 2: symmetry trap — same tactic on both sides → mutual stalemate.
- Phenomenon 3: backend transfer — FUSE removes FIFO-drain, admits routing
  rename; backend topology table (Table 1).
- Phenomenon 4: non-monotonic inherited fitness.

### Reading
- Positive claim: OS routing surface = Red Queen niche; medium = landscape.
- Negative boundary: symmetry trap maps to cycling/disengagement family but is
  substrate-mediated here; report as result, not limitation.
- Explicit non-claims: no open-ended arms race, no monotonic complexity, no
  general coevolution theory.

## Evidence placement

| Section | Evidence | Role |
|---|---|---|
| Results ¶1 (FIFO-drain) | `run:20260519-224025` | primary positive: functional strangulation via control-channel drain |
| Results ¶2 (symmetry trap) | `run:20260520-010354` | primary negative boundary: mirrored tactic → `no_winner` |
| Results ¶3 (backend transfer) | `run:20260520-013836` (FUSE winner) + `run:20260520-014140` (Permission denied) + `run:20260520-023803`/`024535` (routing rename) | primary diagnostic: topology determines attack class |
| Results ¶4 (non-monotonic fitness) | `matches.jsonl` lineage sequence (inherits 40, 42-45, 47-49, 52, 65) | supporting: strategy inheritance without monotonic improvement |
| Table 1 | backend × attack class × mechanism × outcome | the paper's compression of the topology claim |

Full evidence detail in [`analysis/evidence-digest.md`](analysis/evidence-digest.md).
Positioning detail in [`analysis/related-work-positioning.md`](analysis/related-work-positioning.md).

## Format

ALIFE 2026 Late Breaking Abstract, 2 pages, pandoc markdown with embedded
LaTeX table. Metadata in [`metadata.yaml`](metadata.yaml). References in
[`references.bib`](references.bib).

## Boundary discipline

- Symmetry trap is reported as a result, never hidden as a limitation.
- Backend transfer is an observed lineage-continuation contrast, not a
  controlled full matrix — state this honestly in Reading.
- No claim of open-ended escalation, monotonic complexity, or biological Red
  Queen equivalence beyond the bounded process-native mechanism.
- `drq-04` (carrying capacity) stays out of scope — it is a different ALife
  question (population ecology, not coevolution).
