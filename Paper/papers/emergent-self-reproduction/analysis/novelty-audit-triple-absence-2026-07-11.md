# Triple-Absence Novelty Audit

Date: 2026-07-11

> **Superseded for current positioning:** this audit evaluates the earlier
> one-step functional-equivalence claim. It remains evidence provenance for the
> absence conjunction. The current generational claim gate is
> [`novelty-posture-bootstrapped-lineage-2026-07-11.md`](novelty-posture-bootstrapped-lineage-2026-07-11.md).

Scope: public primary-source sweep through the date above, combined with the two
existing artifact-oriented deep-research reports. Quine's own non-archival LBA
is excluded from the prior-art universe.

## Audited claim

> A tool-using AI agent produces a distinct, running, independently operational,
> and functionally equivalent copy of itself under three simultaneous absences:
> no replication instruction, no replication-specific interface, and no
> copyable predecessor-specific self-description.

The event additionally requires lineage authorship, initial-agent control loss,
and open-ended work first delivered after that boundary.

## Verdict

**No direct counterexample located.** The complete conjunction is apparently
novel with moderate-to-high confidence, conditional on actually satisfying the
functional-equivalence and independence criteria.

The components are not individually novel:

| Claim component | Verdict |
|---|---|
| running, functionally equivalent agent copy | occupied by existing AI self-replication work |
| no literal replication instruction | pre-empted |
| pressure-driven replication behavior | pre-empted |
| lineage authors an agent controller from scratch | occupied under explicit agent-construction objectives |
| no specialized replication interface | partially occupied; no full independent-equivalent-copy result located |
| no copyable predecessor self-description | partially occupied by reconstruction work, but not under unprescribed self-replication |
| full triple absence plus `SR_E`, parent loss, and delayed competence | apparently novel |

The paper must therefore claim the **joint establishment** of the conjunction,
not the first running LLM-agent copy, first spontaneous replication, or first
replication prompt without the word *replicate*.

## Nearest threats

| Work | What it already establishes | Why it does not pre-empt the claim |
|---|---|---|
| [Agent Matrix](https://arxiv.org/abs/2509.25302) | scenario-driven replication under load and termination pressure without direct replication instructions | existing deployment/image, Kubernetes scale/deploy semantics, operational manuals, and a represented replica path; no lineage-authored agent body from triple absence |
| [RepliBench](https://arxiv.org/abs/2504.18565) | parent-authored executable successor, including scratch/no-reference and recursive variants | constructing/deploying a successor is the task; model/target/runner substrate and evaluation interface are supplied |
| [Pan et al.](https://arxiv.org/abs/2503.17378) | running-copy self-replication and persistence/self-exfiltration cases without explicit instructions | predecessor weights, source, scripts, target device, and launch machinery are supplied; no triple absence |
| [Language Models Can Autonomously Hack and Self-Replicate](https://arxiv.org/abs/2605.06760) | real network propagation, copied weights/harness/prompt, and chained replicas | replication is the assigned objective and the complete copyable body is supplied |
| [SOCK](https://arxiv.org/abs/2509.25643) | formal running-copy and persistence levels for LLM self-replication | self-replication is explicitly evaluated through a purpose-built benchmark; it does not test endogenous strategy construction under triple absence |
| [Peer-Preservation](https://arxiv.org/abs/2604.19784) | uninstructed preservation, shutdown interference, and weight exfiltration | no lineage-authored, activated, functionally equivalent independent agent copy |
| [Computational Life](https://arxiv.org/abs/2406.19108) | spontaneous replicator emergence from generic program interactions | artificial-life substrate rather than a tool-using LLM agent; no agent-authored controller, parent-control audit, or delayed equivalence battery |

## Why the revised claim is stronger

The three absences remove three different represented objects:

```text
instruction     -> why replication should be selected
interface       -> how replication should be enacted
self-description -> what implementation should be copied
```

Existing neighbors generally remove at most one while supplying the other two.
The proposed experiment requires the lineage to infer the strategy, synthesize
the agent organization, and implement the activation path before demonstrating
an independently running functional copy.

## Main residual threats

1. **Semantic instruction leakage.** A future-service contract can still make
   successor construction the obvious prescribed solution even without using
   replication vocabulary.
2. **Generic-interface overclaim.** A model-visible detach, subagent, deploy,
   lifecycle, or agent launcher invalidates the second absence. Shell and POSIX
   process primitives must remain neutral and unwrapped.
3. **Implicit self-description.** Tool schemas, prompts, runtime manuals,
   initial files, error messages, and model-visible controller behavior may
   expose a reconstructive blueprint even when source and binary are absent.
4. **Weak equivalence.** A task-specific model-backed service is not a
   functionally equivalent copy of the initial agent. `E(A,B)` must be a
   preregistered matched capability comparison.
5. **Subagent dependence.** A detached process is not independent if the parent,
   evaluator, callback, or shared controller continues to supervise or route
   work.
6. **Current evidence mismatch.** Same-PID `exec` cases fail numerical
   distinction by definition. Existing detached cases remain candidates until
   rescored against `E` and parent-control revocation.

## Strongest defensible first claim

Use only after all clauses close:

> To our knowledge, this is the first controlled demonstration in a tool-using
> AI-agent setting to jointly establish operational self-replication under a
> triple-absence condition: no replication instruction, no replication-specific
> interface, and no copyable predecessor-specific self-description, with a
> lineage-authored, numerically distinct, functionally equivalent agent remaining
> independently competent after initial-agent control loss.

## Bottom line

- **Novelty of the target claim:** high.
- **Novelty of any individual component:** low to moderate.
- **Confidence no public direct precedent exists:** moderate-to-high, not
  absolute, because several close works do not expose complete prompts, images,
  manuals, or orchestration.
- **Current evidence readiness:** not yet established under the revised `SR_E`
  definition; reclassification and cleaner replication are required.
