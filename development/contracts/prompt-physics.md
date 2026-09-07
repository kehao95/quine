# Prompt Physics Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md) § Navigation And Routing.

When a change affects the runtime prompt or other model-facing control surfaces, prefer declarative world description over procedural coaching.

Interpretation rules:

- describe what the runtime is, what remains stable, and what degrades under pressure
- prefer ontology and constraints over step-by-step behavioral advice
- assess procedural guidance by whether it helps the current task and accurately describes available operations

Prompt disclosure control-plane rules — per-knob names and values are owned by [`env-controls.md`](../../Paper/core/registries/env-controls.md); what follows is the shape those knobs must keep:

- base process identity ("quine is a running process") stays always visible; do not add knobs that can hide it
- use one explicit mode env per disclosure family rather than one boolean per sentence, and let each control one whole block rather than one sentence
- a runtime behavior gate is not a prompt disclosure mirror: if a knob changes what the runtime *does*, it is a capability env and belongs in that category, not in the `QUINE_PROMPT_*` family
- capability disclosure should follow actual runtime truth; do not add prompt-only mirrors for feature gates that already exist as runtime capability envs
- model-facing gates may be enforced by prompt/schema omission alone when the hidden capability remains harmless, internal, or legacy-compatible if called directly; do not add runtime rejection merely for symmetry unless the capability is unsafe, externally reachable, or the env registry marks it as a runtime acceptance gate
- pressure / scarcity disclosure should be derived from configured limits and memory thresholds, not from a separate prompt-profile switch
- metaphor/style controls are orthogonal: they may style visible content, but they must not reintroduce hidden disclosure

For memory-facing runtime surfaces in particular — the primitives' own design constraints are owned by [`memory.md`](../../Paper/core/primitives/memory.md) § Design Constraints; what belongs here is how those surfaces are *described*:

- raw context should be described as exploratory frontier rather than as a checklist buffer
- anchors should be described as the stable form of finished understanding
- prefer topology and ratio descriptions over warning semantics when the runtime is only exposing information rather than enforcing a hard limit
- prefer local entropy minimization language over summary-task language: already-resolved structure should drift toward compressed, stable form without becoming a competing mission

Progressive disclosure rules:

- the system prompt should describe ontology, active constraints, and primary signals first; secondary topology or debug detail should appear only when it changes action selection
- schema descriptions should state tool contracts and parameter semantics, not serve as miniature runbooks
- when the runtime prompt already owns a global surface map (for example peer-control file names), schema text should name only the tool-local contract and defer the full map to the prompt
- if the same fact appears in both system prompt and schema, one surface should carry the primary statement and the other should only add tool-specific delta
- keep shared runtime descriptions focused; place task-specific procedures with the task that uses them
