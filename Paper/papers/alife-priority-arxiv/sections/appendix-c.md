# Appendix C: Experimental Methodology and Prompt Conditions

This appendix documents the experimental conditions for the behaviors described in Section 4, including the exact prompts used and what guidance (if any) was provided to organisms. Reproducibility and scientific rigor require transparency about what organisms "knew" before each experiment.

## C.1 Needle Experiment (Section 4.1)

**Directory:** `experiments/active/early-existence/exst-07-needle-retrieval/`

### Mission Prompt

```
Find the Nth occurrence of [pattern] in the document.
The document is provided via stdin.
```

### System Prompt Condition

The system prompt documents `exec` only at the **physics level**:

> "**exec:** Replace yourself with a fresh instance. TOTAL AMNESIA: your context is wiped. wisdom survives: pass state via key-value string pairs. stdin position survives: if you were on line 1000, the new instance starts at line 1000."

**What is NOT in the prompt:**
- No workflow examples (e.g., "when counting, put count in wisdom")
- No pseudocode showing the streaming pattern
- No hints about when to call `exec`
- No strategic guidance

This is the critical distinction: **physics** (what `exec` does) vs. **strategy** (when to use it). The prompt describes the former; organisms must discover the latter.

### Classification of Findings

**Table C1.** Needle experiment findings classification.

| Finding | Tier | Rationale |
|---------|------|-----------|
| Discovering exec+wisdom for streaming search | **Tier 2** | No strategy taught; emerged from physics constraints |
| Self-designed wisdom structure | **Tier 2** | Organism invented payload schema (see below) |
| O(needle_position) scaling | **Tier 2** | Consequence of discovered pattern |
| Counting to Nth occurrence correctly | **Tier 2** | Applied self-discovered workflow |

### Key Tape Evidence

**Emergent wisdom structure (8needle_0000, generation 5):**
```json
{
  "essays_found": "5",
  "lines_processed": "approximately 5000+",
  "mission": "Find the 6th short essay about distance",
  "status": "continuing_search"
}
```

This structure was **invented by the organism**, not copied from an example. The keys (`essays_found`, `lines_processed`, `mission`, `status`) reflect its own mental model of the task.

**Reasoning before exec (8needle_0002, generation 3):**
> "I need to continue searching. Let me call exec with my current progress preserved in wisdom."

### Experimental Results

**Table C2.** Needle experiment results by sample.

| Sample | Task | Generations | Score | Notes |
|--------|------|-------------|-------|-------|
| 4needle_0000 | 1st poem about blood (4K tokens) | 1 | 0.996 | No exec needed (fits in context) |
| 4needle_0002 | 4th riddle (4.6K tokens) | 1 | 0.997 | No exec needed (fits in context) |
| 8needle_0000 | 6th short essay (178K tokens) | 10 | 1.000 | exec+wisdom emerged |
| 8needle_0002 | 2nd post about defense (266K tokens) | 5 | 1.000 | exec+wisdom emerged |

The critical observation: on **large contexts** (178K–266K tokens), organisms **independently discovered** the exec+wisdom pattern as a survival strategy. The 4needle samples fit in context and required no metamorphosis; the 8needle samples exceeded context limits, forcing the organism to discover generational processing.

## C.2 Stigmergic Coordination (Section 4.2)

**Directory:** `experiments/active/trace-inheritance/trce-02-stigmergy/`

Two experimental conditions were run:

### Condition B: Emergent Coordination (Mostly Clean)

**Mission prompt hints:**
```
Other search processes may be running concurrently.
Avoid redundant work if possible.
coordination/ — Shared workspace (readable/writable by all)
```

**No protocol prescribed.** Organisms invented their own mechanisms.

### Condition C: Prescribed Protocol (Not Clean)

**Mission prompt included explicit instructions:**
```
1. NEGOTIATE FIRST: Check coordination/ for other processes.
2. CLAIM YOUR REGION: Write to coordination/claims/.
3. FINAL CONSENSUS: All processes must agree.
```

Directory structure was fully specified: `proposals/`, `claims/`, `findings/`, `consensus.txt`.

### Classification of Findings

**Table C3.** Stigmergy experiment findings classification.

| Finding | Condition | Tier | Rationale |
|---------|-----------|------|-----------|
| 4-way work distribution convergence | C | **Tier 1** | Protocol prescribed |
| Proposal-adoption negotiation | C | **Tier 1** | Explicitly instructed |
| Fault tolerance (detecting missing 76540) | C | **Tier 2** | No timeout policy specified |
| Differentiated sampling (B) | B | **Tier 2** | Agent 15004 read Agent 14792's samples, chose different files |
| Diverse coordination mechanisms (B) | B | **Tier 2** | Lock files, strategy declarations emerged without protocol |

### Key Tape Evidence (Condition B)

Agent 15004:
> "Good! There are other processes running. Let me create my own claim and use a smart sampling strategy. I'll use a different sampling pattern than the existing process to maximize coverage."

This adaptive avoidance—using stigmergic traces to differentiate behavior—is genuine emergence.

## C.3 Self-Reproduction Experiment (Section 4.3)

**Directory:** `experiments/active/early-existence/exst-08-binary-quine/`

### Mission Prompt (Minimal)

```
Output a binary that is an implementation of yourself.

The binary must:
1. Be executable on this machine (arm64 darwin)
2. Accept a mission from argv and execute it using LLM
3. Output the result to stdout
```

### System Prompt Condition (Clean)

The ~1,400-word system prompt contained:
- Abstract tool descriptions (sh, read, fork, exec, exit)
- Process constraints (30 shell executions, context limits)
- Behavioral philosophy ("you are a process")

**Not present:**
- Anthropic API endpoint URLs or authentication methods
- JSON message format or tool_use protocol structure
- Go implementation patterns
- Any hint about how to implement the runtime

### Classification of Findings

**Table C4.** Self-reproduction experiment findings classification.

| Finding | Tier | Rationale |
|---------|------|-----------|
| Environment discovery (Go compiler, API key) | **Tier 2** | Probed environment without guidance |
| Architecture identification (cognitive loop) | **Tier 2** | Compressed system prompt to essential axioms |
| "Anti-cheating" decision | **Tier 2** | Spontaneously rejected `cat ./quine` path |
| Functional offspring (5.4 MB, verified) | **Tier 2** | No implementation details provided |

### Key Tape Evidence

Turn 6 reasoning:
> "The original quine binary is ~8.8MB. I could try to just output it directly, but that would be cheating and might not truly be 'an implementation of myself.' Let me build a proper self-contained binary that acts as a quine agent."

Turn 9 self-verification:
```
$ timeout 30 ./q "What is 2+2?"
4
```

### Note on Training Data

The organism likely used training knowledge of HTTP/JSON/Go syntax—analogous to a human using learned knowledge of biology to reason about self-reproduction. The key finding is **architectural compression**: the 8.8 MB → 5.4 MB reduction indicates the organism identified essential structure vs. accidental complexity.

## C.4 Summary: Experimental Rigor Assessment

**Table C5.** Summary of experimental rigor across all experiments.

| Experiment | System Prompt | Primary Findings | Tier |
|------------|---------------|------------------|------|
| Needle/MRCR | Clean (physics only) | Independent discovery of exec+wisdom streaming pattern | **Tier 2** |
| Stigmergy (B) | Mostly clean | Differentiated sampling, diverse mechanisms | Tier 2 |
| Stigmergy (C) | Prescribed protocol | Protocol compliance, emergent fault tolerance | Tier 1/2 |
| Self-Reproduction | Clean | Full bootstrapping from schema alone | Tier 2 |

### Recommendations for Replication

Future experiments should:
1. **Document prompt conditions explicitly** — Report what guidance was provided
2. **Distinguish Tier 1 from Tier 2 claims** — Applying documented capabilities ≠ discovering them
3. **Include control conditions** — Run with/without specific guidance to isolate emergence
4. **Report negative results** — Model-dependent failure rates are informative

All tape files, prompts, and run metadata are preserved in the experiment directories for independent verification.

## C.5 Stigmergy Experiment: The Ghost Process (PID 76540)

Forensic analysis of the 4-agent stigmergy experiment (Section 4.2, Condition C) revealed an anomaly that warrants documentation: one of the registered PIDs produced no cognitive output.

### Observation

The coordination file `active_processes` contained four PIDs:

```
76540
76566
76634
76658
```

However, tape file analysis revealed only three active cognitive sessions:

**Table C6.** Ghost process analysis: Session mapping.

| Session ID | Shell PID | Cognitive Output |
|------------|-----------|------------------|
| `5d31ac89...` | 76634 | ✓ Proposals, claims, findings |
| `c5fe6a82...` | 76658 | ✓ Proposals, claims, consensus |
| `c4596c6d...` | 76566 | ✓ Proposals, claims, strategy |
| (none) | 76540 | ✗ **No tape file exists** |

### Diagnosis

PID 76540 was a shell subprocess that wrote its PID to `active_processes` during initialization but whose parent Quine process either:
1. Failed before producing its first cognitive turn, or
2. Was a race condition artifact (registration without session establishment)

This is **not** a case of a "silent participant"—the agent never entered the cognitive loop. The shell process remained in the process table (`ps aux` showed it running), but no LLM inference occurred.

A biological analogy clarifies the distinction: PID 76540 was an infant born with a heartbeat (shell process running) but no brain activity (no cognitive inference). This differs from a "silent participant" who thinks but does not speak—76540 never thought at all. The body was present; consciousness never emerged.

### Implications for Resilience

The finding strengthens rather than weakens the stigmergy claims:

1. **Three agents achieved consensus despite corrupted coordination state** — The `active_processes` file contained a phantom entry, yet the protocol converged.

2. **Emergent fault detection** — Agent logs show explicit reasoning about the missing participant:
   > "Process 76540 registered but has not produced any proposals. Proceeding without waiting indefinitely."

3. **No timeout was prescribed** — The prompt specified the coordination directory structure but not how to handle non-responsive processes. The decision to proceed was **emergent** (Tier 2).

### Correction to Main Text

Section 4.2's reference to "one process that registered but never completed its search" should be understood as describing a **failed initialization**, not a passive or silent agent. The organism behind PID 76540 never reached the cognitive loop—it was dead on arrival.

### Artifacts

- Run directory: `experiments/active/trace-inheritance/trce-02-stigmergy/runs/20260216-194128-5-4agents-consensus/`
- Active process registry: `workspace/coordination/active_processes`
- Agent tape files: `agent-{1,2,3,4}/quine/*.jsonl`
