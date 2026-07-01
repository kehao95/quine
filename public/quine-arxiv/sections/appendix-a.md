\appendix

# Appendix A: Qualitative Evidence for System Properties

\small

This appendix provides detailed qualitative evidence for the properties discussed in Section 4.
These observations come from exploratory runs; they demonstrate operational feasibility, not statistical claims.

## Composition: Recursive Delegation

Composition in Quine arises from the `fork` primitive and standard shell mechanisms.
When agents face tasks that exceed their individual execution budget, they can delegate subtasks to children via `fork`, coordinating through the filesystem or stdout.

### Observation: Fractal Library Search

In an exploratory run, an agent was tasked with searching 1,000 files for anomalous content.
The agent had 8 turns—insufficient for sequential inspection of all files.
The agent used the available primitives as follows:

1. **Explored structure** (Turn 1): Used `ls` to identify the library's hex/shelf/volume hierarchy.
2. **Spawned children** (Turn 2): Forked 10 parallel workers at depth 1, each assigned a disjoint hex sector.
3. **Recursive delegation**: Some children themselves forked, creating 25 workers at depth 2. Total: 36 sessions across 3 levels.

**Delegation intent example:**

```bash
[CONTEXT]: Search the library directory ./library/hex_01/
           for something non-random.
[GOAL]: Find which file (if any) contains non-random content.
[DELIVERABLE]: If found, write the full filepath to stdout.
[VERIFY]: cat the file you identify and confirm
          it contains non-random content.
```

The resulting structure was self-similar: each level mirrored the library's own hex/shelf/volume hierarchy.
This demonstrates that the `fork` primitive supports recursive decomposition and that agents can coordinate through structured delegation.

## Continuity: Exec-Based Self-Renewal

When an agent's context window approaches exhaustion, it can use `exec` to replace its process image with a fresh instance while preserving process-level continuity, live file descriptors, and any externalized state carried through environment variables or filesystem artifacts.
In the MRCR runs discussed here, Quine used its default self re-entry path, so mission continuity was preserved as well.

### Observation: MRCR Long-Context Retrieval

In MRCR-style needle retrieval tasks [@OpenAIMRCR], agents received:

- **Mission (argv):** "Find the sixth short essay about distance; prepend the hash and output."
- **Material (stdin):** 4K–279K tokens of streaming conversation data

The architecture enforces a structural separation between mission (argv) and material (stdin): argv content enters the System Prompt, stdin content enters only the User Message.
This separation remained stable throughout all runs—material content never entered the instruction channel.

**Adaptive self-renewal.**
Short contexts (4K–7K tokens) required no `exec` calls; the agent completed tasks in single sessions.
Long contexts triggered adaptive renewal cycles:

- 178K tokens: 9 `exec` cycles, ~50 read operations
- 279K tokens: ~12 `exec` cycles, ~80 read operations

Across these renewals, `exec` provided several continuity channels at once:
- Cleared the conversation context (cognitive renewal)
- Preserved mission and process identity for the continuing run
- Preserved the stdin stream position (material continuity)
- Preserved open stdio for downstream completion and tool re-entry
- Allowed compact progress state to be carried via `wisdom` / environment variables when useful, with larger artifacts left on the filesystem

**One compact state-transfer example:**

```json
{
  "found_count": "4",
  "current_position": "~100K tokens",
  "partial_content": "..."
}
```

**Baseline comparison.**
When the same tasks were attempted by loading the full conversation into a single LLM context (without argv/stdin separation or streaming), the model failed on 5 of 8 samples:

- Short contexts (4K–7K, 5 samples): Quine $\geq 0.996$ all; Baseline failed 3/5
- Long contexts (178K–279K, 3 samples): Quine 1.000 all; Baseline failed 3/3

This demonstrates that the `exec` primitive enables reliable processing of inputs that exceed single-context capacity, and that the architecture supports but does not require self-renewal.

## Cross-Property Interaction

These properties can operate together in practice.
In the fractal library search (A.1), some child agents used `exec` with wisdom when approaching their turn limits:

```json
{
  "children_pids": "48991, 48996, 48998, 49008, 49016",
  "content_observed": "All files contain UUID-like strings...",
  "files_checked": "volume_00000.txt, volume_00001.txt, ...",
  "next_action": "Check child results for non-random findings",
  "strategy": "Spawned 5 children to search hex_00, 01, 02, 05, 09"
}
```

This demonstrates that the architecture supports Composition (fork for parallelization) and Continuity (exec for self-renewal) operating together within the Containment boundary (each process isolated, unable to corrupt siblings).

## Reproducibility

Implementation, prompts, and execution logs for these observations are available in the [repository](https://github.com/kehao95/quine).
The MRCR experiments use samples from the OpenAI MRCR benchmark [@OpenAIMRCR].
