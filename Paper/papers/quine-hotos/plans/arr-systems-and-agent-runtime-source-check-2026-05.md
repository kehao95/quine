# Deep Research Query: ARR Systems and Agent Runtime Source Check

## Launch Status

- **2026-05-07:** Drafted after the three-paper maturation split and the first
  `6M` synthesis. This is a merged systems source-check for the HotOS ARR route
  and its necessary citation cleanup.
- **Interaction id:** `v1_ChdRRUQ5YWFDSUM0clFqTWNQOFByVnNRRRIXUUVEOWFhQ0lDNHJRak1jUDhQclZzUUU`.
- **Model / mode:** launched with `dr run --detach --max`.
- **Metadata:** [`arr-systems-and-agent-runtime-source-check-2026-05.meta.json`](arr-systems-and-agent-runtime-source-check-2026-05.meta.json)
- **Status:** completed; report absorbed in `../references/`.
- **Expected report home:**
  `../references/arr-systems-and-agent-runtime-source-check-report-2026-05.md`

## Context

Quine is a POSIX-native LLM agent that can inspect runtime-carried
self-description, construct or rebuild a successor execution body, and cross a
real process-image replacement boundary. The HotOS-style paper frames this as
**Agentic Runtime Replacement (ARR)**: a bounded transition in which the active
process-native agent participates in constructing or selecting a successor
runtime body, while a runtime protocol admits, executes, and challenges that
successor across native process-image replacement.

The claim must be systems-specific. ARR is not that `exec` is new, not that
software update is new, and not that LLM self-improvement is achieved. The
load-bearing criterion is **post-replacement viability through the future-work
carrier**: addressability after replacement, delayed challenge competence, and
return-to-ready. Quine's strongest current instance is control-surface closure:
a rebuilt successor crosses real `exec`, becomes reachable through `ctl/wake`,
handles delayed unseen work with runtime tools, and returns to ready.

The paper needs source-checked related-work anchors and reviewer-proof wording
for the nearest systems neighbors.

## Core Research Question

Which systems and LLM-agent literatures are the closest primary-source
neighbors for ARR, and how should a HotOS-style paper distinguish agent
construction/selection plus runtime admission/challenge from existing wrapper,
update, restart, migration, and agent-design work?

## Detailed Questions

1. What are the canonical primary sources for POSIX `exec` / process-image
   replacement semantics, including the fact that successful `exec` overlays
   the calling process image rather than returning to it?
2. What prior systems work best represents daemon re-exec, rolling restart,
   supervisor-managed restart, actor hot-swap, service hot-swap, and external
   orchestration of runtime bodies?
3. What dynamic software update systems best represent externally authored and
   externally admitted updates, and what do they preserve across update
   boundaries?
4. What checkpoint/restore and process migration systems best represent
   snapshot-based continuity rather than successor construction?
5. What LLM-agent frameworks or agentic-system design papers explicitly manage
   wrappers, scaffolds, code generation, self-modification, or automated agent
   design, and do any make the generated code the native continuation runtime
   after real process-image replacement?
6. What work on runtime audit, provenance, observability, or challenge-response
   verification can support the ARR audit story: artifact hash, exec evidence,
   PID/process lineage, fd/carrier map, delayed payload timestamp, proof, and
   return-to-ready marker?
7. What reviewer attacks are most likely, and what concise response wording
   should the HotOS paper use for each?
8. Which experimental controls are most important to design next: subprocess vs
   real exec, scripted re-exec baseline, post-handoff carrier audit, multi-wake
   readiness, or successor admission prototype?

## Relevance Filters

Prioritize primary systems papers, official specifications/documentation, and
well-cited technical reports. Include LLM-agent papers only when they affect
runtime-body ownership, generated scaffolds, self-modification, automated agent
design, or wrapper-managed continuity. De-emphasize generic code-generation,
benchmark self-improvement, and security-risk speculation unless directly
relevant to ARR distinctions.

## Deliverables Requested

Return a concise source-check report with:

1. One-paragraph verdict on the narrow novelty of ARR.
2. Nearest-neighbor table with columns: area, source, who constructs the
   successor/update, who admits it, handoff mechanism, continuity criterion,
   post-handoff verification, distinction from ARR.
3. Ten to fifteen primary-source citation anchors with stable URLs/DOIs where
   possible.
4. Reviewer attack list with safe response wording.
5. Suggested related-work paragraph outline for a five-page HotOS-style paper.
6. Experimental implications: which systems controls would most improve the
   ARR claim and why.

## Red Lines

Do not frame the claim as self-awareness, autonomous propagation, generic
self-replication, recursive self-improvement, or malware. Treat ARR as a
bounded local systems abstraction for process-native agent runtime replacement.

