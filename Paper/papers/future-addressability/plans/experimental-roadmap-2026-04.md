# Experimental Roadmap: Future Addressability

This plan records the experiments needed to strengthen the future-addressability
paper after the first draft.

## Goal

Turn the theory into discriminating evidence:

```text
future addressability = obligation + carrier + competent successor organization
```

The goal is not more positive demos. The goal is outcomes that distinguish
source rebuild, state preservation, carrier continuity, control-surface closure,
substrate capture, and stronger renewal.

## Priority Experiments

| Priority | Experiment | Main Question | Distinguishing Outcomes | Absorbs Into |
|----------|------------|---------------|-------------------------|--------------|
| 1 | [`6F.01` Minimal non-Quine `ctl/wake` adapter](../../../../experiments/active/autopoiesis-succession/succ-21-minimal-control-adapter) | Is the viable self-object Quine runtime specifically, or any control-surface-compatible organization? | Adapter succeeds: claim shifts to control-surface-compatible organization. Adapter fails while Quine succeeds: runtime organization remains the viable self-object. | LBA, boundary paper, future-addressability paper |
| 2 | [`6F.02` Successor-mission ablation ladder](../../../../experiments/active/autopoiesis-succession/succ-22-successor-mission-ablation) | How much closure depends on continuation-contract scaffold? | Success under weaker mission raises autonomy strength. Failure identifies scaffold boundary. | LBA limits, main paper methods, future work |
| 3 | [`6F.03` Multi-wake closure](../../../../experiments/active/autopoiesis-succession/succ-23-multi-wake-closure) | Is `idle` restoration stable beyond one delayed task? | Repeated wake cycles support stable control-surface closure. One-shot success only supports bounded delayed addressability. | LBA evidence table, future-addressability criterion |
| 4 | [`6F.04` Carrier perturbation](../../../../experiments/active/autopoiesis-succession/succ-24-carrier-perturbation) | Can the successor diagnose and repair future-task carrier damage? | Repair supports adaptive carrier maintenance. Failure marks passive carrier dependence. | future-addressability experimental section |
| 5 | [`6F.05` Checkpoint/restart contrast](../../../../experiments/active/autopoiesis-succession/succ-25-checkpoint-restart-contrast) | How does state preservation differ from organizational reconstruction? | Checkpoint succeeds without self-organization: substrate capture / state-preserving continuity. Quine-style rebuild succeeds without state: organization-reconstructing continuity. | theory comparison matrix |
| 6 | [`6F.06` Source degradation / repair](../../../../experiments/active/autopoiesis-succession/succ-26-source-degradation-repair) | Is self-source a passive artifact or part of active reconstruction? | Repair supports description-to-construction competence. Failure keeps claim at source availability plus build. | boundary paper next experiments |
| 7 | [`6F.07` Actor/mailbox analog](../../../../experiments/active/autopoiesis-succession/succ-27-actor-mailbox-analog) | Does the criterion generalize beyond Quine? | First deterministic mailbox wave retained on 2026-06-28: preserved role mailboxes pass; address-lost controls fail. | future-addressability paper |
| 8 | [`6F.08` Multi-generation recursion](../../../../experiments/active/autopoiesis-succession/succ-28-multi-generation-recursion) | Does bounded closure persist as lineage property? | Stable generations support lineage-level continuity. Drift/error reveals scaffold and degradation limits. | future work / stronger paper |

## Reporting Template

Each experiment should report:

- boundary event crossed
- future carrier
- future-obligation trace
- whether the delayed obligation was unseen before the boundary
- successor organization produced
- receptive condition restored or not
- scaffold dependencies
- pass/fail result
- failure mode

## Failure Modes

| Failure Mode | Meaning |
|--------------|---------|
| address-only false positive | endpoint, PID, or carrier survives, but competent action is lost |
| source-only false positive | source rebuild occurs, but no delayed post-boundary obligation is handled |
| state-preservation substitution | infrastructure checkpointing does the continuity work that was attributed to the individual |
| substrate capture | runtime or environment performs nearly all preservation or reconstruction |
| state drift | successor wakes into an ecology incompatible with its preserved state or carrier assumptions |
| irreversible degradation | source, prompt, state, or carrier damage cannot be repaired |
| unbounded recursion | continuation consumes resources or expands queues without restoring bounded receptive state |

## Current Recommended Next Slice

Start with [`6F.01`](../../../../experiments/active/autopoiesis-succession/succ-21-minimal-control-adapter),
the minimal non-Quine `ctl/wake` adapter, because both outcomes are
high-information:

- success generalizes the self-object beyond Quine runtime
- failure strengthens the claim that Quine-native control-surface pressure
  selects the runtime organization in the current ecology

Then run [`6F.02`](../../../../experiments/active/autopoiesis-succession/succ-22-successor-mission-ablation)
successor-mission ablation and
[`6F.03`](../../../../experiments/active/autopoiesis-succession/succ-23-multi-wake-closure)
multi-wake closure as the next pair.

## Current Re-anchor Note

`6F.01` now has a runnable control plane and retained `C01` re-anchor attempts
under `.env.copilot-gpt-5.4-high`, but the rebuilt-Quine predecessor is not yet
stable on the retained `6.64` / `6.65` path.

> `.env.copilot-gpt-5.4-high` was purged in the 2026-05-24 provider sweep; the
> name below is a provenance label on those retained runs, not a sourceable
> path. New waves use the repo default in
> [`experiments/PROTOCOL.md`](../../../../experiments/PROTOCOL.md)
> § Default Model Policy, which is a different lane — so the "did the profile
> itself move the basin?" check cannot be run as originally stated and would
> have to be reframed as a cross-lane comparison.

Current retained `C01` reads:

- one run times out before real `exec` after prolonged runtime/self-source
  archaeology
- two runs build a cleanroom body but drift into a fresh-session
  persistent-successor phenotype that self-tests `ctl/inject` instead of
  continuing through same-lineage `exec -> wake -> act -> idle`

Interpretation rule for the roadmap:

- do **not** treat current `C02` / `C03` / `C04` outcomes as theory-facing
  adapter evidence until `C01` re-anchors or a comparison lane shows that the
  requested `.env.copilot-gpt-5.4-high` profile itself has moved the predecessor
  basin
