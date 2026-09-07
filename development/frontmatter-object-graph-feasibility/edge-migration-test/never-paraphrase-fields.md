---
kind: idea
thesis: The structural losses that survived every compression scheme in the 2026-07-02 experiment — dialectical sign (self-debunking / already-occupied / retire-to-citation), evidence maturity, falsifier and conditional scope, opposite-signed counter-theses, family self-attribution, kill/retirement conditions, and time-critical alerts — are all small machine-readable fields that the sparks mostly already carry in prose; they were lost because paraphrase does not transport fields and the sign flips exactly at the paraphrase; therefore these dimensions must live in fixed frontmatter slots and every downstream view must assemble them mechanically rather than re-narrate them.
origin: conversation (2026-07-02), grounded by [[compression-loss-experiment-2026-07-02]]
links: [idea-space-density-not-tier, consequence-identity, read-layer-prior, selection-signal-collapse]
edges:
  - { rel: grounds-in,    to: theory/references/compression-loss-experiment-2026-07-02.md }
  - { rel: superseded-by, to: incubator/metadata-as-object-graph.md }
  # NEEDED EXTENSION: `superseded-by` is not itself attested anywhere in the
  # repo. Its mirror, `refines`, IS attested (theory/objects/proxy-successor.md
  # uses `refines: [target] — narrows...`). `superseded-by` is proposed here
  # only as the reciprocal-mirror form of that already-attested verb
  # (matching the enables/enabled-by, carries/carried-by mirror-pair
  # convention theory/README.md already names) -- not a free invention, but
  # genuinely new, so flagged per the brief's instructions.
---

# Never Paraphrase the Fields

**Status:** seed
**Evidence:** [`compression-loss-experiment-2026-07-02`](../../../Paper/theory/references/compression-loss-experiment-2026-07-02.md)
(five independent compression schemes, five independent audits).

## Core idea

The compression experiment's *invariant* losses share one property: they are
**not the bulk of the text** — they are small, structured, often single-line
fields (a status, a falsifier, a counter-thesis pointer, a maturity marker, a
deadline). They were destroyed not because the compressor ran out of budget
but because **prose paraphrase does not carry fields** — and, worse, the
*sign* of a claim flips precisely at the paraphrase (a spark whose real thrust
is "retire this to citation-only" reads as a live proposal once summarized).

The corpus already encodes most of these (paradigm-fit verdict, evidence
status, novelty note) — in prose the compressor discards. So the lever is not
"summarize better" but **move the load-bearing dimensions out of prose into
fixed, machine-readable slots**, and require every generated view to *assemble*
those slots verbatim rather than re-describe them.

## The dimensions that must be field-borne (from the audits)

Each was lost by ≥4 of the 5 schemes:

- **dialectical sign** — is this a forward claim, a self-correction of a prior
  seed, or a debunking of its own premise?
- **evidence maturity** — confirmed / one-run / red-teamed / analogy-only, with
  the run id.
- **falsifier + conditional scope** — "false if X"; "only for trait-class Z".
- **counter-thesis pointer** — the opposite-signed hypothesis this must be read
  against.
- **external-occupancy / novelty status** — already-published? near-duplicate
  of which prior?
- **kill/retirement condition** — what would retire this to citation-only.
- **home family** — the paradigm this self-attributes to (four sparks were
  mis-parented in scheme D because this was in prose, not a field).
- **time-critical flag** — deadline / prior-art collision / governance alert.

## Revision (2026-07-02): fields are edges, not enum slots

The first reading of this seed reached for closed frontmatter enums
(`dialectical-sign: forward|self-correction|self-debunk`). That is premature
convergence *and* storage-layer compression — an enum is a lossy projection
baked into the store. The corrected form:
[`metadata-as-object-graph`](../../../Paper/incubator/metadata-as-object-graph.md) (see the
`superseded-by` edge in frontmatter — this used to be a body-only mention,
**absent from `links:`**, i.e. exactly the silent-drift failure this file's
own thesis warns about; the edge form makes that drift structurally
impossible). Most of the dimensions below are **relations, not scalars** — a
counter-thesis is an edge to the opposing idea, evidence-maturity is an edge
to the grounding run (its level *derivable*, per `check-theory-maturity.sh`),
a falsifier is an edge to a prediction. So "move the dimension out of prose"
becomes "make it a typed edge into the addressable object graph," not "add
an enum slot." An edge is the limit form of a paraphrase-resistant field: it
resolves or it does not; a view *traverses* it, never re-narrates it. The
small committed set lives at the relation-*type* level (which has a
consumer), not the value level (which stays open).

## Why this is the light, no-commit-yet lever

It changes nothing structural: no new directories, no lifecycle, no deletion.
It only raises frontmatter utilization on files that already exist, so the
cost of *not* compressing drops (a reader scans fields) and the cost of *safe*
compression drops (a view concatenates fields). It is the minimal move the
experiment's data directly warrants.

## Open questions

- Which slots are load-bearing enough to be worth the write-time cost vs.
  which are bloat (the schema-bloat pull — every slot needs a reader that acts
  on it)?
- Do the fields belong in YAML frontmatter, a fixed prose block, or both?
- Can a validator check field *presence* without lapsing into the presence-vs-
  meaning blindness the repo already documented?
