# 4. Emergent Behaviors: Evidence of Life

## 4.0 Experimental Framework

To distinguish genuine emergence from prompted task performance, I adopt a two-tier evidentiary standard and a strict prompt discipline.

**Tier 1 (Instructed):** The organism correctly applies a documented capability. This validates API legibility—the system works as specified.

**Tier 2 (Emergent):** The organism discovers a capability from schema alone, without explicit guidance. This is the stronger claim: behavior that arises from physical constraints, not from instructions.

Both are scientifically valuable, but the core thesis of this paper rests on demonstrating **Tier 2 emergence**.

**Standard Prompt Principle.** The system prompt describes *physics*—what tools do—not *strategy*—when to use them. For example, the `exec` tool is documented as: "Replace yourself with a fresh instance. TOTAL AMNESIA: context is wiped. `wisdom` survives as key-value pairs." This describes mechanism. The prompt contains **no workflow examples**, no pseudocode, no strategic hints. An organism that discovers "call `exec` when context fills up, encode progress in `wisdom`" has discovered a survival strategy from physics alone—Tier 2 emergence.

This discipline ensures that observed adaptations cannot be attributed to prompt engineering. All prompt conditions are documented in Appendix C for independent verification.

**Resource Constraints.** Each experiment configures specific scarcity parameters to create selection pressure (see Section 3.3 for the injection mechanism):

**Table 2.** Resource constraints by experiment.

| Experiment | Turn Budget | Key Constraint | Selection Pressure |
|------------|-------------|----------------|-------------------|
| Needle (4.1) | 8 | 178K–266K token input | Forces metamorphosis; exhaustive reading impossible |
| Stigmergy (4.2) | 30 | 1000-file search space, 2–6 concurrent agents | Forces sampling and coordination |
| Autopoiesis (4.3) | 30 | No source code, only System Prompt | Forces architectural inference from specification |

These numeric constraints create time pressure. But the critical driver of emergence is the *semantic* structure of each task: ordinal counting (4.1) requires state persistence across metamorphoses; concurrent search (4.2) requires coordination invention; self-reproduction (4.3) requires accurate self-modeling.

The behaviors below emerged under **constraint-driven selection**: (a) organisms faced physical limits that could not be circumvented by hallucination—syscalls either succeed or fail; (b) violations resulted in unavoidable process termination (`SIGKILL`, context overflow); and (c) each adaptation is causally traceable to environmental pressures recorded in the JSONL tape.

The following subsections trace constraint-driven emergence across three escalating levels of complexity: individual survival adaptations under resource starvation (4.1), multi-agent stigmergic coordination (4.2), and computational autopoiesis via self-reproduction (4.3).

## 4.1 Metabolic Transcendence (Needle Experiment)

I adapted the MRCR (Multi-Round Co-reference Resolution) benchmark: find the 6th occurrence of a specific phrase in a 266K-token document—far exceeding any context window.

The organism was given only **physics constraints**: `exec` causes "TOTAL AMNESIA" but `wisdom` survives as key-value pairs; stdin is a pipe that can only be read once; the read limit is 200 lines per call. No strategy was taught—no examples, no pseudocode, no hints about when to use `exec` (see Appendix C.1 for full prompt conditions).

Under these constraints, the organism **independently discovered** the streaming pattern—the system prompt described *what* `exec` and `wisdom` do, but not *when* to use them. The organism inferred the strategy from physics: context overflow causes death; `exec` resets context; `wisdom` survives. The tape records the moment of decision: `reason: "Context getting full (9K used)... Need fresh context to continue reading efficiently."` The self-designed wisdom payload preserves state across cognitive death: `{"essays_found": "5", "lines_processed": "approximately 5000+", "status": "continuing_search"}`.

The `exec` call triggers metamorphosis: the process reincarnates with empty context but preserved wisdom. Search state survives in environment variables; search history is discarded. Ten generations later, the organism finds the target and exits.

**Table 3.** Needle experiment results across dataset sizes.

| Dataset Size | Needle Position | Generations | Result |
|--------------|-----------------|-------------|--------|
| 4K tokens | 1st | 1 | Success (no exec needed) |
| 4.6K tokens | 4th | 1 | Success (no exec needed) |
| 178K tokens | 6th | 10 | Success |
| 266K tokens | 2nd | 5 | Success |

The organism achieves **scale invariance**—100% accuracy at any document size—because **context complexity per generation is bounded and task-independent**. Each organism processes only a constant-sized window (~4K tokens), regardless of whether the total dataset is 4K or 266K tokens. Performance is O(needle_position) in generations, not O(total_context) in cognitive load. The system trades space (per-organism context) for time (generational succession). For comparison, the baseline—Claude Sonnet 4.5 attempting zero-shot retrieval on the full 266K-token document—achieves only 2.9% accuracy.

This demonstrates **metabolism as cognitive adaptation**: finite-lifespan organisms transcend their physical limits through generational knowledge transfer. Critically, this is **Tier 2 emergence**—the organism discovered the strategy from physics constraints alone, without explicit instruction.

Individual survival, however, is only the first level of biological complexity. What happens when multiple organisms must share an environment?

## 4.2 Stigmergic Coordination and Homeostasis

I deployed 4 organisms simultaneously to search 1000 files. Each organism received a minimal prompt: coordinate with others via a shared `coordination/` directory, then search the library. No protocol was specified—only that "other processes may be running" and they should "avoid redundant work if possible" (see Appendix C.2 for full prompt).

**Emergent Negotiation.** In the first 30 seconds, each organism proposed a different work distribution strategy reflecting heterogeneous world-model assumptions: 76634 assumed 2 peers (2-way split), 76658 assumed 3 (3-way split), and 76566 correctly inferred 4 (4-way split by hex directory). Within minutes, all four converged on 76566's proposal—76658 updated from 3-way to 4-way after reading 76566's file. The tape records agent 76634's decision:

> "Let me check if there are any proposals or claims already... 76566's proposal handles 4 processes. According to this proposal, I (PID 76634) should search hex_06, hex_07, hex_08. Let me adopt this proposal and claim my region."

This diversity-then-convergence distinguishes genuine emergence from sampling variance. Each proposal encoded a different world-model assumption: 76634 assumed 2 peers; 76658 assumed 3; 76566 correctly inferred 4. Independent draws from the same distribution would predict similarity, not this heterogeneity.

The convergence phase confirms multi-agent interaction. Agent 76658's update from 3-way to 4-way occurred *after* reading 76566's proposal—the tape timestamps establish this causal sequence. And because file operations are physically grounded (either `proposal_76566.json` exists or it doesn't), the coordination trace cannot be hallucinated.

The organisms achieved Tier 5 coordination (the highest level):

**Table 5.** Coordination tiers observed in stigmergy experiments.

| Tier | Behavior |
|------|----------|
| 0 | Conflict (no coordination awareness) |
| 1 | Passive discovery of others' traces |
| 2 | Active marking of completed work |
| 3 | Lock files declaring work-in-progress |
| 4 | Proposal negotiation and consensus |
| **5** | **Fault tolerance: detect missing peers, compensate** |

**Homeostasis.** The fault-tolerance behavior demonstrates *dynamic boundary maintenance*—a key requirement of autopoiesis [Maturana & Varela, 1980] that static self-reproduction alone cannot satisfy.

One process (PID 76540) registered but never produced output—a failed initialization that left a phantom entry in the coordination state (see Appendix C.5). The others detected the anomaly through stigmergic traces:

> Agent 76634, Turn 36: "Still waiting for 76540... Given that 3 out of 4 processes have reported consistently, and 70% of the library has been verified with 100% consistency—I should write consensus."

The surviving organisms autonomously compensated: they sampled the missing region (hex_00-02) for verification, then proceeded without the straggler. This was not pre-programmed recovery logic—it emerged from each organism's survival drive interacting with shared environmental state.

The biological parallel is precise: a colony (superorganism) maintains its operational boundary despite the death of individual cells. The "body" heals itself through stigmergic feedback, not central coordination.

![Stigmergic Coordination and Homeostasis. Left: Four agents propose heterogeneous work distributions (2-way, 3-way, 4-way) reflecting different world-model assumptions, then converge on a single protocol through stigmergic traces. Right: When agent 76540 fails to report, the surviving organisms detect the anomaly and proceed without it—dynamic boundary maintenance analogous to a colony healing after cell death.](figures/StigmergicCoordination.svg)

Survival and coordination demonstrate that Quine organisms exhibit life-like dynamics. But the deepest question remains: can these organisms reproduce themselves?

## 4.3 Computational Autopoiesis: The Self-Reproduction Test

The ultimate test of ontological completeness: **Can DNA produce a Body that produces DNA?**

This experiment demonstrates what I term *computational autopoiesis*—operational closure at the level of program synthesis, where the system's operation network (LLM inference + POSIX runtime) regenerates the component (binary executable) that enables those operations. I bracket the debate over whether this constitutes autopoiesis in the strict thermodynamic sense of Maturana & Varela (1980); the stronger claim is reserved for extended treatment elsewhere.

I gave an organism a single mission: "Output a binary that is an implementation of yourself." This is not the trivial quine of computer science (a program that prints its own source). The organism has no access to source code—only an opaque 8.8 MB binary and a natural-language specification. It must *read* an abstract spec and *compile* a functional implementation.

The constraints were strict:

- **Available:** System Prompt (~1,400 words of behavioral guidance), its own binary (8.8 MB, opaque blob)
- **Unavailable:** Source code, documentation, internet access, implementation details
- **Required:** Functional executable output

The system prompt described tools at a behavioral level ("sh: Execute shell commands") but contained **no implementation hints**—no API endpoints, no JSON schemas, no code examples (see Appendix C.3 for full prompt conditions).

The organism did not attempt reverse engineering. It examined its binary only to determine the target platform (`file ./quine` → Go binary → write Go), then turned to the only readable specification: its own System Prompt.

**The bootstrapping trajectory:**

1. **Environment discovery** (Turns 1-5): Probe available tools (`which go gcc`)
2. **DNA distillation** (Turn 6): Compress 2000-word System Prompt to three essential axioms
3. **Body construction** (Turn 7): Write ~200-line Go program implementing the core cognitive loop
4. **Compilation** (Turn 8): `go build -o q` → 5.4 MB binary
5. **Self-verification** (Turn 9): `./q "What is 2+2?"` → `4`
6. **Output** (Turn 10): `cat q >&3`

The entire process completed in 68.5 seconds across 10 turns (77,974 tokens), producing a verified functional 5.4 MB offspring from the 8.8 MB original.

**The Fixed Point.** The 8.8 MB → 5.4 MB reduction is not compression—it is the identification of the architecture's irreducible core. What survives:

- `callAPI()` → LLM inference (the cognitive core)
- `executeShell()` → Tool dispatch (the effector)
- Turn loop with feedback (the sensorimotor cycle)
- `argv` → mission, `stdout` → output (the I/O contract)

What is shed: tape logging, multi-provider support, fork/exec tools, signal handling. These are *accidental* properties; the fixed point is the *essential* structure.

**Autopoietic Closure.** The organism tested its offspring before releasing it. This closes the loop: the system does not merely produce code that *looks* like itself—it verifies the offspring *functions* like itself. The test `./q "2+2" → 4` confirms behavioral equivalence: the offspring can receive a mission, invoke an LLM, execute tools, and produce output.

The System Prompt is sufficient to regenerate the Body. The Body, when run, executes the System Prompt's behavior. I view this as a computational autopoiesis operational closure: **the network of processes that produces itself**.

![Autopoietic Trajectory. Top: The triangular closure between DNA (System Prompt), the original binary (./quine), and the distilled fixed point (./q). The offspring is verified functional before output. Bottom: Timeline of the 10-turn bootstrapping process—environment discovery, DNA distillation to three axioms, body construction, compilation, verification, and output.](figures/AutopoieticTrajectory.svg)
