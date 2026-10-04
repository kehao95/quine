# L2 — defeat the deflation (attraction ≠ constraint)

**Scope narrowed 2026-07-05:** this layer originally bundled two designs —
(a) attraction≠constraint and (b) structural-law≠inherited-human-bias. (b)
is **dropped**: "not inherited human bias" is not practically provable or
disprovable (there is no clean way to rule out that a capable model's
behavior traces back to some pattern in its training corpus — the negative
claim has no clean falsification design), and the paper does not need it —
see `claim.md`'s Revision note. The out-positioning against the nudge/
behavior-design literature still holds without it: we can contrast
*framing/register* (they treat the effect as an exploitable vulnerability;
this claim treats it as a deliberate, positive design surface) without
needing to settle the effect's causal origin. L2 now tests only (a).

**Status:** not yet run — optional strengthening, not a required gate (per
the rescoped claim, L0 already anchors the core inversion; see
`../manuscript-plan.md`).
**Gate this closes:** attraction ≠ constraint — is the environment actively
tilting the distribution over freely-available options, or merely
restricting which options exist (the "it's just the sandbox constraining
you" reviewer reflex).

**Proves:** it is attraction (not the sandbox constraining).

**Design — action-space-invariant:** verify the counterfactual act was
always executable in the control arm too, then show the distribution over
freely-available acts still shifts toward the structure-supplied target in
the treatment arm.

**Forward-pointer:** will need the harness extended with an executability
verifier for the counterfactual act — does not yet exist in
`experiments/active/structural-elicitation/`.
