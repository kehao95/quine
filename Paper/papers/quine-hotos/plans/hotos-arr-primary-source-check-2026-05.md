# Deep Research Query: HotOS ARR Primary-Source Check

## Launch Status

- **2026-05-06:** Drafted after first HotOS-style Agentic Runtime Replacement
  manuscript draft.
- **Interaction id:** `v1_ChdtZV83YWYzMUI5Nkotc0FQLWEtVy1RMBIXbWVfN2FmMzFCOTZKLXNBUC1hLVctUTA`.
- **Model / mode:** launched with `dr run --detach --max`.
- **Metadata:** [`hotos-arr-primary-source-check-2026-05.meta.json`](hotos-arr-primary-source-check-2026-05.meta.json)
- **Status:** completed; report absorbed in `../references/`.
- **Expected report home:**
  `../references/hotos-arr-primary-source-check-report-2026-05.md`

## Context

Quine is a POSIX-native LLM agent that can construct or rebuild a successor
runtime body and cross real process-image replacement. The HotOS-style paper
frames this as **Agentic Runtime Replacement (ARR)**: the active agent
participates in constructing or selecting a successor body, while a bounded
runtime protocol admits, executes, and challenges that successor across native
process-image replacement. The criterion is post-replacement viability through
the future-work carrier: addressability, delayed challenge competence, and
return-to-ready.

The paper must not claim that `exec`, dynamic software update, checkpointing,
or self-modifying code are new. It needs primary-source citation anchors and
reviewer-trap distinctions.

## Core Research Question

Which systems literatures are the closest primary-source neighbors for ARR, and
how should the HotOS draft cite and distinguish them without overclaiming?

## Detailed Questions

1. What are the canonical sources for POSIX `exec` semantics and daemon re-exec
   / process-image replacement?
2. Which dynamic software update systems best represent externally authored and
   externally admitted runtime updates?
3. Which checkpoint/restore or process migration systems best represent
   snapshot-based continuity rather than constructed successor bodies?
4. Which actor hot-swap, Erlang/OTP, Akka, Kubernetes, or service restart
   mechanisms are likely reviewer comparisons, and what exactly is externally
   managed in those systems?
5. Which LLM agent / self-modifying agent papers are the closest neighbors, and
   do any perform native process-image replacement into an agent-constructed
   successor body?
6. Which observability / agent tracing systems can anchor the need for runtime
   audit evidence?
7. What claims in the current ARR framing should be weakened before HotOS-style
   submission?

## Relevance Filters

Prioritize primary systems papers, official documentation, or well-cited
technical reports. Include secondary surveys only when they point to primary
sources. Exclude generic code-generation or benchmark self-improvement papers
unless they directly involve generated code becoming the continuation runtime
body.

## Deliverables Requested

Return a concise citation pack with:

1. One-paragraph verdict on ARR novelty.
2. Nearest-neighbor table with columns: area, primary source, what it provides,
   who constructs successor/update, who admits it, handoff mechanism,
   post-handoff verification, distinction from ARR.
3. Ten primary-source citation anchors with stable URLs/DOIs where possible.
4. Reviewer attack list and safe response wording.
5. Suggested edits to the current HotOS abstract/thesis if any wording is too
   broad.

## Red Lines

Do not frame ARR as self-replication, recursive self-improvement, malware, or
proof of autonomous propagation. The claim is a bounded local systems
abstraction for process-native agent runtime replacement.
