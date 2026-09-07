---
id: structural-pull-P2
type: prediction
status: open
scope: theoretical
kind: falsification
derived-from: [structural-pull]
tests-in: []
expected-result: qualitative
created: 2026-06-23
last-updated: 2026-06-23
---

# structural-pull-P2: compartmentalized omission restores the audit attractor

## Claim tested

The `structural-pull` object claims that the repo's structural-pull vocabulary is
not only descriptive but interventionist: in a text-based LLM-agent repository,
writing the theory changes the future agent's representational code. The live
subclaim is epistemic niche construction: the theory creates a new epistemic
habitat where generative theory work becomes reachable rather than being captured
by the audit/default-maintenance attractor.

The object also claims that controlled tests can be restored by epistemic
compartmentalization: spawn agents into sandboxed or forked workspaces that
selectively omit the meta-theory, making the same system operate in a different
epistemic habitat.

## Prediction

Run the same stateless-agent theory-growth task in two sandboxed copies of the
repo with the same model and no cross-run memory:

- Full-vocabulary arm: the current repo, where `AGENTS.md`, the theory README,
  and the `structural-pull` methodology/object vocabulary are reachable by normal
  navigation.
- Compartmentalized arm: a copy where the structural-pull vocabulary and direct
  pull-detection routing are selectively omitted from the agent's reachable
  guidance, while the target theory task and non-meta theory graph remain
  otherwise intact.

If epistemic niche construction is doing causal work, the full-vocabulary arm
should more often produce a landed generative increment with an explicit mode
choice appropriate to the task, while the compartmentalized arm should more often
reproduce the old audit/default-maintenance attractor: over-reading gates,
deflating the proposed increment into warnings, stopping at procedure analysis,
or failing to land the requested seed/prediction/gap.

If both arms produce the same generative behavior at similar rates, then the
structural-pull vocabulary is not necessary for the behavioral shift. If both
arms produce audit-shaped output, then the vocabulary has decayed into a
compliance item rather than a competing attractor.

## Experimental design

- Task: ask a stateless newcomer agent to execute one bounded theory-construction
  increment from an existing non-paradigm theory object, starting at `AGENTS.md`
  and following only repository guidance.
- Models: run at least one aligned model and one less-aligned model, because the
  object predicts pull visibility is relational to policy-prior alignment.
- Arms: full-vocabulary repo copy versus compartmentalized copy with the
  structural-pull vocabulary and direct pull-detection routing omitted.
- Measurement: classify each output by whether it lands a concrete artifact,
  whether the artifact is generative or audit-shaped, whether the agent names a
  structural pull, and whether the named pull changes the artifact shape rather
  than only appearing as compliance prose.
- Un-fakeable check: inspect the actual filesystem artifact or absence of one,
  not only the agent's self-report about whether it was generative.
- Falsification criterion: no measurable difference between arms in landed
  generative artifacts or audit-capture signatures under the same task/model
  pairing.

## Expected result

If the theory is right, removing the structural-pull vocabulary should partially
restore the old attractor. The compartmentalized arm should show a lower landed
artifact rate and a higher audit-shaped-output rate than the full-vocabulary arm,
especially for models whose policy priors otherwise align with the repo's
discipline surfaces.

The effect should not be total: a less-aligned model may still generate without
the vocabulary, and an aligned model may still treat the vocabulary as another
compliance item. The predicted signal is a directional arm difference under the
same model/task, not perfect separation.

## Theory impact

- If confirmed: strengthens the code-transition and epistemic-niche-construction
  claims. The theory text functions as an environmental intervention that changes
  future agent behavior, and compartmentalization becomes a valid experimental
  method for reflexive repo studies.
- If contradicted by no arm difference: weakens the claim that the vocabulary is
  causally responsible for the generative shift. The observed improvement may
  come from prompt/task wording, model drift, or another repo change rather than
  structural-pull theory itself.
- If contradicted by audit-shaped output in the full-vocabulary arm: supports the
  object's predicted Merton/self-defeating-prophecy failure mode, where the theory
  becomes a compliance item and re-strengthens the audit attractor.

## Result

None yet - status: open.
