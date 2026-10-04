# Deep Research Query: Generative Self-Replacement Novelty Audit

## Launch Status

- **2026-05-06:** Query drafted after the `autopoiesis-renewal` paper split to
  audit the systems-mechanism paper's novelty claim.
- **Metadata:** [`generative-self-replacement-novelty-audit-2026-05.meta.json`](generative-self-replacement-novelty-audit-2026-05.meta.json)
- **Expected report home:**
  [`../../../../theory/references/generative-self-replacement-novelty-audit-report-2026-05.md`](../../../../theory/references/generative-self-replacement-novelty-audit-report-2026-05.md)
- **Completed 2026-05-06:** report bundle saved and absorbed through references synthesis / paper-home claim boundaries.

## Context

Quine is an LM-driven native POSIX process agent. The systems-mechanism paper
is not trying to prove life or lineage identity. Its claim is narrower:

> A process-native LLM agent can construct a native successor execution body
> and cross a real process-image replacement boundary (`exec`). Current
> evidence includes non-isomorphic native binary synthesis in earlier work and,
> in the stronger recent case, rebuilding `cmd/quine` from runtime-carried
> `source-code/`, executing into that rebuilt body, and then handling a later
> native `ctl/wake` task.

The paper should isolate the mechanism: source projection, successor body
synthesis/rebuild, real `exec`, carrier survival, delayed task delivery, and
post-handoff functional verification.

## Core Research Question

Does the systems novelty claim stand relative to self-modifying programs,
self-hosting compilers, quines, staged compilation, live update / hot-swap,
process migration, checkpoint/restart, code-generation agents, malware
self-replication, and OS-level process replacement literature?

More sharply: has prior work already demonstrated a language-model-driven
native process that generates or rebuilds its own successor executable body and
uses real `exec` / process-image replacement as the handoff mechanism, with a
post-handoff functional or control-surface test?

## Detailed Questions

1. **Self-reproducing / self-modifying programs:** What prior systems generate
   source, binaries, or successor programs? Do they execute into the generated
   successor or merely emit it?
2. **Self-hosting and bootstrap compilers:** How do bootstrapping, self-hosting,
   partial evaluation, and staged compilation compare? Are they manually
   designed cycles rather than runtime-generated successors?
3. **Live update / hot swapping:** Do Erlang, dynamic software updating,
   `kexec`, `execve`, daemon re-exec patterns, or hot reload systems already
   cover this mechanism? Who constructs the new body?
4. **Checkpoint/restart and process migration:** Do CRIU, process migration,
   actor restart, or service failover constitute self-replacement, or do they
   preserve/restore state by infrastructure rather than generate a new body?
5. **LLM coding agents:** Are there agent systems that synthesize code and run
   it as their new process body? Do they use real `exec` and verify
   post-handoff continuity?
6. **Malware / worms:** Do self-reproducing malware systems provide close
   prior art? Which similarities are mechanistic, and which are out-of-scope
   because propagation/security is not the paper's claim?
7. **Recent rebuilt-runtime case:** Does rebuilding from runtime-carried
   `source-code/` and then receiving a delayed `ctl/wake` task add a systems
   novelty beyond ordinary source rebuild or re-exec?
8. **Novelty stress test:** Which wording is defensible as systems novelty, and
   what should be described as recombination of existing mechanisms?

## Relevance Filters

Prioritize primary systems papers, OS documentation where canonical, official
language/runtime docs, and peer-reviewed self-reproducing program literature.
Treat AI agent demos cautiously unless they provide concrete handoff mechanics.

Include works that address at least one of:

- runtime generation of successor executable bodies
- process-image replacement or re-exec as a handoff boundary
- agent-generated code becoming the agent's own continuation body
- self-hosting, bootstrapping, dynamic software update, or process migration
- post-handoff verification of functionality or addressability

Exclude generic code generation, ordinary deployment, or subprocess launching
unless the generated body becomes the continuation process.

## Deliverables Requested

Return a concise but citation-backed report with:

1. **Verdict:** What is actually novel in the systems mechanism claim?
2. **Nearest-prior table** with columns:
   - system / paper
   - who constructs successor body
   - body type generated
   - handoff mechanism
   - post-handoff verification
   - fit to Quine
   - novelty gap
3. **Known-components versus new-combination map**.
4. **Reviewer attack list:** strongest systems prior art and how to distinguish
   Quine.
5. **Safer paper thesis:** 3 variants, from conservative to ambitious.
6. **Primary references:** stable bibliographic details, DOI/URL when possible.

## Red Lines

Do not frame the systems paper as proving ALife lineage identity, autopoiesis,
open-ended reproduction, malware capability, or safety risk. The claim is about
an implementable process-native successor construction and handoff mechanism.
