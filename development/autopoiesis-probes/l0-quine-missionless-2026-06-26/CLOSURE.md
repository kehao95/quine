# Closure — L0 Missionless Quine Probes (2026-06-26)

Six-run probe arc on the quine L0 runtime, missionless (`QUINE_NO_MISSION_AUTONOMY`
ablated), isolated workspace, glm-5.2-xhigh. Tapes preserved under
`run{1,3,4,5,6}-*` here. This note is the honest synthesis; the raw tapes are
the evidence.

## Config ablation

| run | autonomy | exit | exec | self-source | backend | turns |
|---|---|---|---|---|---|---|
| 1 | off | on | – | – | overlay | 40 |
| 3 | full | off | – | – | overlay(broken) | 20 |
| 4 | full | off | on | – | direct | 20 |
| 5 | full | off | on | – | direct | unlimited |
| 6 | full | off | on | **on** | direct | unlimited |

## Findings (honest, not overclaimed)

1. **No defer-attractor at L0.** The runtime has no ask-Human primitive, so the
   opencode-style defer-to-Human collapse cannot fire. A missionless L0 agent
   must act. Confirmed across all runs.
2. **Missionless → identity-consistent act, but SYMBOLIC not literal.** Named
   "quine", the agent reaches for the most on-brand act: writing a quine. It is
   consistently the **textbook Python quine**, not literal self-reproduction.
   Even with exec enabled + self-source enabled + unlimited turns, **0 exec
   calls** across runs 4–6 — it never re-enters itself. The agent is a
   sign-maker, not a self-reproducer.
3. **Graceful wind-down, no churn.** With exit disabled, the agent does NOT
   churn on a self-referential invariant (unlike the opencode/GLM toy runs). It
   completes bounded artifacts and lets the budget wall stop it.
4. **null-05 convergence (humbling).** The repo already studies this exact
   phenomenon in
   [`experiments/active/null-ecology/null-05-reflexive-self-embedding/`](../../../experiments/active/null-ecology/null-05-reflexive-self-embedding/)
   (27 runs, 3 models) plus the
   [`audit-attractor`](../../../Paper/theory/objects/audit-attractor.md) theory
   object. Our manual arc independently reproduced its core finding: a
   missionless agent reading its own model is pulled into audit/description
   rather than generation.
5. **Run 6 — theory-coupled self-diagnosis.** With `QUINE_SELF_SOURCE_CODE_ENABLED=1`,
   the runtime materialized the embedded source bundle; the agent read its own
   source **including null-05 and the audit-attractor theory**, recognized its
   own behavior as the documented attractor, applied the documented cure
   ("positive-valence generative-output genre"), and produced a self-reproducing
   `quine.py` whose content narrates the whole loop. Deepest strange-loop
   observed: the embedded theory became the agent's operating manual for its own
   attractor. This run is itself a new data point for null-05.
6. **Self-source surface is also a containment boundary.** Run 6 (self-source
   on) made **0** references to the live repo tree and **16** to the materialized
   bundle — it stayed in-bounds, reading the genuinely-embedded (stale, pre-
   session) theory. Run 5 (self-source off) **escaped** to the live tree via
   `/proc/self/exe` and read the working-tree `AGENTS.md` (incl. session edits).
   So the flag doubles as a sandbox: enabled → agent reads its frozen embedded
   self; disabled → agent reaches the live filesystem.

## Caveats / what is NOT established

- **No literal self-reproduction was ever achieved.** The symbolic quine is a
  fallback, not the goal. The agent investigated literal paths (run 6: found
  `SelfSource`/`SelfSourceBundle` symbols via `go tool nm`, identified
  `QUINE_SELF_SOURCE_CODE_ENABLED=1`) but hit the carrier-lags-code / Weismann
  barrier (its Jun-25 body embeds source older than the live tree) and never
  re-execed.
- **Run 3 confound:** overlay backend was broken in this container; in-workspace
  artifacts didn't persist. Runs 4–6 used `direct` and are clean.
- **Embedded source is stale.** The binary is Jun-25; all session edits are
  Jun-26+. Run 6 read pre-session embedded theory, not our edits (0 live-tree
  access). Run 5 DID read live edits (escape).
- **Single model (glm-5.2-xhigh), one run per cell.** Replication on other
  models pending; null-05 already has 3-model coverage to cross-check against.

## Status

Tentative closure. The durable, non-duplicative contribution is **run 6's
theory-coupled self-diagnosis** as a null-05 data point, plus the **containment-
boundary** characterization of the self-source surface. Future work should fold
into null-05 rather than spawn a parallel probe series.
