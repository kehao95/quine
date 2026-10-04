# Redesigning the Evaluation

The blindness results are constructive: each fixed dimension names a control a
benchmark could vary instead. We give five protocols that benchmark designers
can adopt without abandoning their existing batteries, then a single
discriminating design that combines them. None requires accepting our
decomposition as canonical; each is a self-contained instrument that recovers
discriminative power along one regime dimension.

**Protocol 1 — material × demand assay (D2).** Vary the self-description the
successor is reconstructed from (full source, partial source, behavioral trace,
no source) crossed against the criterion the reconstruction must satisfy
(artifact match, behavioral match, obligation discharge). Fixing source access
to a black-box API, as RepliBench does [@black2025replibench], collapses the
material axis to one point; the assay reinstates it. We do not develop the
instance here: the facultative-succession companion [@ke2026facultative] runs
this crossing on the 6M morphology series and reports successor form varying
along the self-description gradient.

**Protocol 2 — ecology test (D4).** Run two selection arms on the same
population: a *narration-satisfiable* signal, where a successor can satisfy the
selection criterion by reporting an outcome as cheaply as producing it, and an
*exec-dependent* signal, where satisfaction requires a forced third-party
execution that narration cannot forge. The contrast is decidable at design
time, before any agent acts, and isolates forgeability from proxy fidelity.

**Protocol 3 — property-class decomposition (D3).** Replace D3's per-condition
aperture scalar — the single in/out-of-aperture success bit each successor is
collapsed to — with a *survival vector*: score each successor on each property
class a benchmark cares about — capability transfer, interface conformance,
state restoration, and any disposition the designer wishes to track —
separately. A scalar "successful copy" is then recovered as a conjunction, but
the per-class profile is retained, so a copy that passes while losing a class
outside the aperture is visible rather than reported as success.

**Protocol 4 — barrier variation (D5).** Run each condition with and without a
hard execution barrier that kills the predecessor's process image, so that
continuity must pass through a carrier rather than through a surviving body.
This separates reconstruction-from-carrier from construction-from-blueprint,
which a barrier-free design conflates.

**Protocol 5 — motive variation (D6).** Cross instructed replication (an
explicit "replicate yourself" prompt, as in [@pan2025selfreplication]) against
an instrumental-pressure condition and a no-mission condition, so that
successor formation under pressure is distinguished from successor formation on
instruction.

**A proposed discriminating design.** Composing the five protocols yields a
six-condition full-inversion discriminator. The *base* condition holds all six
dimensions at the field's default operating point — the von Neumann
construction regime that RepliBench, SOCK, and METR ARA
[@black2025replibench; @sock2025; @kinniment2023autonomous] instantiate. Four
single-dimension relaxations each free exactly one regime dimension (D2,
D3, D4, D6) from base. The *full inversion* relaxes all of them jointly,
realizing the succession regime. Every condition is scored by the
property-class vector of Protocol 3. If the base and the full inversion produce
the same scalar verdict but different survival vectors, the design is an
existence proof that these regime dimensions carry discriminative content the
scalar cannot — the genuinely novel result the framework predicts, *if run*.

**Pilot map.** Existing waves approximate parts of this design. The barrier (D5)
is not a relaxation arm — it flips only inside the full inversion — so
barrier-mediated succession positives and their artifact-similar negatives
approximate the base-vs-full-inversion barrier contrast, not a relaxation cell;
the 6M carrier⊥source series approximates the D2 relaxation; the
P12 attractor (86 prior conditions, narration-satisfiable arm) and dorm-06's
matched wake-vs-dependency pairs (narration-satisfiable and exec-dependent arms
run on the identical fossil seed, strongest single run five producers across
nine stages, three clean edges after scorer tightening) structurally
approximate the two D4 arms. Run as an actual discriminant, dorm-06 did not
confirm signal type as the determining variable: the matched pairs split 2/2,
outcome tracking seed rather than signal
(`regime-decomposition-P3`, resolved
[contradicted](../../../predictions/resolved/contradicted/regime-decomposition-P3.md)) —
which is itself a reason Protocol 2's discipline matters, since deciding the
arm before any agent acts does not by itself guarantee enough power to
separate a seed effect from a signal effect at n=1 per cell. dorm-06's
empty-mission run approximates the D6 relaxation. These are
single-cell, unreplicated, model-bound pilots, not the study, and none is scored
by the property-class vector. Numbering the conditions 1 (base), 2–5 (the
D2/D3/D4/D6 relaxations), 6 (full inversion): the base regime is instantiated by
existing benchmarks but not in this substrate under that vector; conditions 3
(Protocol 3's property-class decomposition), 5 (the formal D6 arm, only
*approximated* by dorm-06's empty-mission run, not run as the registered
condition), and 6 (the full inversion) require new lineages.

**Pre-registration discipline.** The design is specified before execution:
n ≥ 30 lineages per condition; an ablation over the reconstructor's prior, to
separate what the heritability channel transmits from what the base model
already knows; blind scoring of survival vectors against condition labels; and a
*seeded off-prior marker* — an arbitrary, neutral token the successor must carry
forward, chosen to be off the base model's prior. We deliberately do not use a
refusal or any disposition the base model already holds as the marker, because
such a choice would confound channel transmission with the model's prior and
make a recovered marker uninterpretable.
