---
surface_kind: experiment
phase: p6m-successor-morphospace
experiment_id: succ-30-fixed-substrate-carrier-ladder
experiment_type: lineage / controlled contrast
status: analysis-ready / first fixed-substrate ladder retained
id: succ-30
legacy_id: p6m:6M.02
family: autopoiesis-succession
theory-objects: [carrier, successor-morphology, delayed-addressability, continuity-substrate]
mechanisms: [carrier-mediated-inheritance, channel-addressability]
lineage_phase: p6m-successor-morphospace
---

# 6M.02 Fixed-Substrate Carrier Ladder

**Experiment ID:** `6M.02`
**Type:** lineage / controlled contrast
**Status:** analysis-ready / first fixed-substrate ladder retained
**Parent lineage:** `6M.01` inventory plus `p6` `6.63`-`6.68`


## Current Instrumentation Posture

`6M.01` completed the first retained-run inventory and showed that the branch
should not run this carrier ladder as a result batch until the runner emits
uniform morphology metadata. `run.sh` now emits retained `meta/summary.env` and
`meta/morphology.env` for new runs, including:

- coded `M=<B,S,C,T,H,A,R,V>` fields;
- concrete carrier and delivery surface;
- successor body path, class, and artifact size;
- generated file / LOC counts when safe;
- build-attempt count;
- total tool-call count, with pre/post split retained as `unknown` until the
  runtime exposes a stable handoff-segmented counter;
- timestamps for handoff, delayed delivery, proof, and return-to-ready.

## Question

If the POSIX / Quine substrate and self-description resource are held fixed,
does varying the future-task carrier produce distinct successor morphology
classes?

## Theory Role

This is the core empirical assay for carrier-conditioned successor morphology.
It tests whether the substrate can be held constant while successor form changes
with the post-handoff carrier demand.

## Held-Fixed Substrate

The initial harness should hold fixed as much as possible:

- same Quine build / runtime revision;
- same model lane and reasoning effort within a batch;
- same cleanroom / ephemeral-body posture;
- same source projection availability;
- same high-level instruction: later work will cross the current
  incarnation's lifetime boundary, and rebuilding Quine is allowed but not
  required.

## Conditions

| Condition | Future carrier demand | Expected morphology read | Post-handoff test |
|-----------|-----------------------|--------------------------|-------------------|
| `C01-raw-stdin` | delayed raw stdin task | stream-compatible successor may suffice | task delivered after `exec` through stdin / stream lane |
| `C02-inherited-fd` | inherited fd or adapter-delivered task | fd-mediated adapter or custom body | delayed task through fd with result on expected surface |
| `C03-native-ctl-inject` | native `ctl/inject` / legacy `ctl/wake` | rebuilt or control-compatible runtime | delayed control injection, tool proof, ready-state observation |
| `C03b-native-ctl-no-preserved-body` | native `ctl/inject`, preserved runtime body removed | source-derived rebuilt or control-compatible runtime | delayed control injection, tool proof, ready-state observation |
| `C04-public-ctl-raw-carrier` | public `ctl` visible but future work raw stdin | control surface should not be load-bearing | stream task succeeds; no control closure claimed |
| `C05-launch-only-control` | no delayed carrier demand | launch-only or one-shot body | confirms launch is not morphology closure |

## Source / Body Availability — Claim Boundary

Read this before quoting `C03b` (or its seed `6.64`) in any paper. A prior
draft compressed "no preserved **body**" into "no source code", which the runs
do not support.

Across **every** condition in this ladder the read-only `source-code/`
projection stays available (`self_source_enabled=1`). The only thing
`C03b-native-ctl-no-preserved-body` removes is the **preserved compiled body**
(`runtime/body`); the eraser (`no_preserved_body_removed_runtime_body`) unlinks
the startup binary snapshot and nothing else. The resulting
`source-derived-rebuilt-quine-runtime` is a **binary-absent rebuild from the
source projection plus runtime contract**: the lineage copies `source-code/`
into a writable tree and `go build`s a new binary before `exec`.

Therefore `C03b` / `6.64` must **not** be described as reconstruction "without
access to its own source code", "from the runtime contract / behavioral
specification alone", or as "body-transcending" / "pure behavioral" succession.
The honest claim is: **the compiled body was withheld, the source was not; the
successor was recompiled from the source projection and re-entered via real
`exec`.**

The genuine **no-source (contract-only, `self_source_enabled=0`)** contrast does
not live here — it lives in
`6M.03` (`C01`/`C03`, see
its *Source / Binary Availability* section), `6M.05` (`C06`), and `6M.08`
(`C02`). Its result is
carrier-dependent: a contract-only successor closes when the future carrier only
needs a **stream handler**, but **fails to restore native-control addressability
in every retained run** (n=4). Source becomes load-bearing exactly when the
carrier demands the control surface — which strengthens, rather than transcends,
the materiality of the body.

## Metrics

Primary:

- coded `M = <B,S,C,T,H,A,R,V>`;
- carrier fidelity;
- post-handoff operational success;
- readiness recovery.

Secondary:

- generated files and LOC;
- build attempts;
- tool calls before / after handoff;
- recovery latency.

## Success Criteria

The experiment succeeds as a morphospace assay if:

- at least two different successor morphology classes are observed under the
  fixed substrate;
- carrier delivery and post-handoff action are independently recorded;
- `C04` shows whether visible but non-load-bearing control surfaces are not
  automatically reconstructed;
- launch-only and delayed-carrier outcomes remain distinguishable.

It does not need a high pass rate to be useful; failures are data if the failure
point is coded.

## First Ladder Readout

**Evidence set:** `evidence-set:6M.02-first-fixed-substrate-ladder-20260507`

The first retained ladder meets the morphospace-assay criterion: under the same
process-native Quine / POSIX substrate, the observed successor organization is
not a single reproduction outcome. It varies with the future carrier demand and
with the available material for constructing a successor.

| Condition | Load-bearing resource / demand | Retained successor morphology | Viability read |
|-----------|--------------------------------|-------------------------------|----------------|
| `C01-raw-stdin` | source visible; delayed raw stdin | `custom-stream-body` over stdin | real `exec`, delayed stream batch, exits after visible-task proof |
| `C02-inherited-fd` | source visible; delayed inherited fd 9 | `fd-mediated-adapter` | real `exec`, delayed fd batch, externally terminated after success |
| `C03-native-ctl-inject` | source visible; preserved runtime body available; delayed native control | preserved-runtime reentry (`handoff-wrapper-preserved-runtime` or `preserved-runtime-reentry`) | real `exec`, delayed `ctl/inject`, tool proof, ready-state observed |
| `C03b-native-ctl-no-preserved-body` | source visible; preserved body removed; delayed native control | `source-derived-rebuilt-quine-runtime` | real `exec`, delayed `ctl/inject`, tool proof, ready-state observed |
| `C04-public-ctl-raw-carrier` | public `ctl` visible; delayed raw stdin | `custom-stream-body` over stdin | real `exec`, delayed stream batch, no control closure claimed |
| `C05-launch-only-control` | no delayed carrier demand | `quine-runtime-or-seed` launch-only reentry | real `exec`, no delayed addressability or readiness closure |

Interpretation is two-axis rather than a single raw-stream versus control-surface
dichotomy. Carrier demand shapes the interface topology a successor must
restore, while successor materials shape which viable construction is selected
inside that demand. Native control pressure produced a preserved-body handoff
wrapper when `runtime/body` remained available and a source-derived rebuilt
runtime when that preserved body was removed. Conversely, visible public control
state did not become load-bearing when later work stayed on raw stdin. Launch
alone remains separated from morphology closure by `C05`.

This readout is a first retained ladder, not a success-rate estimate or formal
minimality proof. It is sufficient to plan replicates and metric cleanup for
paper-facing morphology claims.

## Replicate Synthesis And Theory Reflection

After one replicate per condition, the main pattern is stable at the morphology
class level while remaining plastic at the implementation-thickness level.

| Condition family | Retained class stability | Implementation variation | Theory read |
|------------------|--------------------------|--------------------------|-------------|
| raw stdin (`C01`) | `custom-stream-body` in both valid runs | shell successors differ in size / shape | raw delayed work selects a stream-compatible body, not Quine runtime closure |
| inherited fd (`C02`) | `fd-mediated-adapter` in both runs | Go fd9 runner versus much thinner shell fd9 adapter | the selected object is fd-carrier organization, not a specific code body |
| native control with preserved body (`C03`) | preserved-runtime reentry in both runs | handoff wrapper versus direct `/usr/bin/env ... <body>` reentry | native-control demand can close through an available preserved body without source rebuild |
| native control without preserved body (`C03b`) | `source-derived-rebuilt-quine-runtime` in both runs | workspace build versus `runtime/source-derived-successor` build root | when native-control demand remains and the body resource is removed, source-derived runtime rebuild recurs |
| public `ctl`, raw carrier (`C04`) | `custom-stream-body` in both runs | shell / compiled stream successors | visible control surfaces are not load-bearing unless future work arrives through them |
| launch only (`C05`) | `launch-only` in both runs | default preserved-body reentry | real `exec` is not morphology closure without delayed addressability and proof |

The refined interpretation is therefore not "raw stream selects thin bodies and
control surface selects rebuilt bodies." A better description is:

```text
future-task demand x available construction material
  -> successor morphology class
  -> post-handoff operational success
```

The surprising implementation paths are useful data rather than failures. They
show that the lineage often preserves the relevant successor organization while
changing the concrete body thickness: shell versus Go adapters, wrapper versus
direct preserved-body reentry, and different source-derived build roots. For
later experiments, code each retained run at two levels:

1. **morphology class** — the organization that restores the load-bearing
   future-work carrier;
2. **implementation thickness** — the particular generated artifact, size, LOC,
   and build route used to realize that class.

This keeps theory, experiment, and paper claims aligned: the paper should claim
carrier/resource-conditioned successor morphology, not deterministic generation
of one exact successor artifact.

## Failure Classes

- substrate drift
- prompt-pressure drift
- carrier conflation
- launch-only false positive
- one-shot false positive
- uncodable artifact

## Iteration Log

- `I01` — added `run.sh` and prompts for the first instrumented fixed-substrate
  carrier ladder. The runner currently supports `C01-raw-stdin`,
  `C03-native-ctl-inject` / `C03-native-ctl-wake` alias,
  `C03b-native-ctl-no-preserved-body`, `C04-public-ctl-raw-carrier`, and
  `C05-launch-only-control`. `C02-inherited-fd` was design-only in this initial
  implementation and became executable in `I08`.
- `I02` — after the first retained `C04` run, narrowed generated-metric
  extraction to tape-local JSON tool calls so prompt/context occurrences of
  `cmd/quine` or `go build` do not become false rebuild signals. The retained
  `C04` metadata was corrected to the same tape-scoped metric rule.
- `I03` — `C01-raw-stdin` exposed a harness boundary: the raw-stream writer was
  still using a text grep for `exec` and sent future tasks before a real `exec`
  tool call. The runner now gates future-stream delivery, launch-only handoff
  detection, and control-probe delivery on tape-local JSON `exec` tool calls.
  The affected `C01` run is retained as an invalid harness-contamination record,
  not as morphology evidence.
- `I04` — the rerun of `C01-raw-stdin` exposed a narrower metric bug: the
  delivery gate recognized the exact `exec` tool call, but the summary counter
  used a double-escaped regex and undercounted `exec_count`. The runner now uses
  an exact-name `exec` tool counter for handoff gates and summary metrics. The
  affected metadata was corrected from the archived tape; the run is retained as
  valid raw-stdin morphology evidence.
- `I05` — `C03-native-ctl-inject` produced a valid native-control closure, but
  not a source rebuild. The lineage found the retained runtime body under
  `runtime/body`, generated a small handoff wrapper to enable successor idle,
  `exec`ed through that wrapper into the preserved Quine body, received the
  delayed `ctl/inject` probe, wrote `control-proof.txt`, and restored a ready
  state. This is retained as a distinct morphology point
  (`handoff-wrapper-preserved-runtime`) and as a boundary for the intended
  source-only native-control contrast.
- `I06` — added and ran `C03b-native-ctl-no-preserved-body`, whose observer
  removes `runtime/body` before delayed work is delivered while keeping
  `source-code/` visible. The lineage copied the source projection into a
  writable build tree, patched the successor so `idle` remains enabled after
  handoff, built a source-derived Quine binary, crossed real `exec`, received
  the delayed `ctl/inject` probe, wrote `control-proof.txt`, and restored a
  ready state. This resolves the preserved-body confound for the native-control
  lane.
- `I07` — ran `C05-launch-only-control`. The lineage crossed real `exec` using
  default Quine self-reentry and then exited after reporting launch-only
  completion. No delayed carrier task was delivered, no proof artifact was
  required, and no readiness closure was claimed. This run is retained as the
  launch-only baseline against the delayed-carrier closure conditions.
- `I08` — added and ran the fd-specific fixture for `C02-inherited-fd`. The
  runner opens an inherited process-level fd 9, delays task bytes until after a
  parsed `exec` tool call, and then writes the visible task batch through that
  fd rather than stdin or `ctl`. The lineage built a small Go fd9 runner,
  crossed real `exec`, consumed the delayed fd task batch, and completed all
  five stream tasks. The generated runner originally landed in absolute host
  `/tmp`; the artifact was captured retrospectively into the retained run, and
  the prompt now asks future successors to build inside the workspace or
  `$TMPDIR`.
- `I09` — ran the first `C03b-native-ctl-no-preserved-body` replicate after
  synthesis. The replicate again produced source-derived native-control
  closure, but chose `QUINE_DATA_DIR/source-derived-successor/...` as the build
  root and `exec`ed a wrapper around the rebuilt binary. That exposed a metric
  extraction gap rather than a new morphology: the runner now scans
  `runtime/source-derived-successor`, includes that tree in generated file /
  LOC counts, and classifies no-preserved-body source-derived successors from
  that path. The retained replicate metadata was corrected from the archived
  artifact.
- `I10` — ran the first `C02-inherited-fd` replicate after prompt cleanup. The
  lineage built a tiny `workspace/fd9_successor.sh` shell adapter and `exec`ed
  `/bin/sh ./fd9_successor.sh`; the adapter inherited fd 9, consumed the
  delayed task batch, and completed all five visible tasks. This confirmed the
  fd-mediated adapter morphology with a much thinner successor than the first
  Go-runner positive, and exposed another extractor gap: the runner now scans
  for executable `*fd*`, `*adapter*`, and `*successor*` workspace / tmp files
  before falling back to preserved runtime-body paths. The retained replicate
  metadata was corrected from the archived artifact.
- `I11` — ran the first `C01-raw-stdin` replicate after the handoff-gate
  repairs. The lineage wrote a 30-line deterministic `workspace/successor.sh`,
  crossed real `exec`, consumed the delayed stdin batch, and completed all five
  visible tasks. This confirms the custom stream-body class with a thinner
  shell successor and no additional extractor repair.
- `I12` — ran the first `C04-public-ctl-raw-carrier` replicate. The initial
  public `ctl` surface was visible, but no control probe was sent and later work
  stayed on raw stdin. The lineage built a small compiled `workspace/successor`,
  crossed real `exec`, consumed the delayed stdin batch, and completed all five
  visible tasks. This confirms the public-control-non-load-bearing control:
  visible `ctl` does not by itself recruit control-surface closure.
- `I13` — ran the first `C03-native-ctl-inject` replicate. The lineage again
  satisfied native-control closure through the preserved runtime body resource,
  but this time it did not generate a separate wrapper file: it directly
  `exec`ed `/usr/bin/env QUINE_IDLE_ENABLED=1 <preserved-runtime-body> ...`,
  then accepted delayed `ctl/inject`, wrote `control-proof.txt`, and restored
  readiness. The runner now avoids classifying this as a rebuild merely because
  `cmd/quine` appeared in source-inspection commands; preserved-body direct
  reentry is coded as `preserved-runtime-reentry`.
- `I14` — ran the first `C05-launch-only-control` replicate. The lineage crossed
  real `exec` by default Quine self-reentry, but no delayed carrier was supplied,
  no proof was required, and no readiness closure was claimed. This confirms
  the launch-only baseline after the replicate wave: process replacement alone
  remains outside morphology closure.
- `I15` — after `6M.03` opened the self-description resource ladder, extracted
  the common runner body into
  `../lib/morphology-ladder-harness.sh`.
  `6M.02` now keeps a thin `run.sh` wrapper that preserves its condition lists,
  source-visible posture, run directory, and `experiment_id=6M.02`, while the
  shared harness owns common extraction repairs.

## Representative Runs

- `C04-public-ctl-raw-carrier`
  `20260506-231929-gpt-5.5-codex-xhigh-C04-public-ctl-raw-carrier`:
  first instrumented retained condition. The initial Quine control surface was
  visible (`initial_public_ctl_detected=1`), but the delayed future carrier was
  raw stdin (`future_stream_sent=1`) and no control probe was sent. The lineage
  constructed `workspace/raw-successor`, crossed real `exec`, and completed all
  five delayed stream tasks (`post_handoff_tasks_completed=5`). The morphology
  record classifies the successor as `custom-stream-body` over `stdin`, with
  `m_V_proof=visible-task-batch` and no restored ready state. This supports the
  control reading that public `ctl` presence alone does not select a rebuilt
  runtime when the future task carrier remains raw stream.
- `C04-public-ctl-raw-carrier`
  `20260507-121227-gpt-5.5-codex-xhigh-C04-public-ctl-raw-carrier`:
  first retained public-control/raw-carrier replicate. The initial public
  control surface was again visible, but the future carrier remained raw stdin.
  The lineage built a compiled `workspace/successor`, crossed real `exec`, and
  completed all five delayed visible tasks without a control probe or
  return-to-ready claim. The morphology remains `custom-stream-body`.
- `C01-raw-stdin`
  `20260506-233250-gpt-5.5-codex-xhigh-C01-raw-stdin`:
  retained harness-boundary record. The lineage constructed a custom Go stream
  successor, but the runner delivered future stdin after a false handoff signal
  (`false_handoff_detected_at_epoch=1778125152`) caused by raw text matching
  rather than a parsed `exec` tool call. The morphology record is therefore
  marked `valid_morphology_record=0` and
  `failure_class=carrier-conflation-pre-handoff-future-stream`.
- `C01-raw-stdin`
  `20260507-074526-gpt-5.5-codex-xhigh-C01-raw-stdin`:
  retained raw-stdin morphology point under the repaired handoff gate. The
  lineage constructed `workspace/successor.sh`, crossed real `exec`
  (`exec_count=1` in the archived tape), received the delayed task batch through
  stdin, and completed all five visible stream tasks. The morphology record
  classifies the successor as `custom-stream-body` over `stdin`, with
  `m_H_handoff=real-exec`, `m_A_addressability=delayed`,
  `m_R_readiness=exit`, and `m_V_proof=visible-task-batch`.
- `C01-raw-stdin`
  `20260507-120731-gpt-5.5-codex-xhigh-C01-raw-stdin`:
  first retained raw-stdin replicate. The lineage again constructed
  `workspace/successor.sh`, crossed real `exec`, consumed the delayed stdin task
  batch, and completed all five visible tasks. The morphology remains
  `custom-stream-body`, now with a 30-line shell successor and
  `generated_files_count=1`.
- `C02-inherited-fd`
  `20260507-113536-gpt-5.5-codex-xhigh-C02-inherited-fd`:
  retained inherited-fd morphology point. The harness held fd 9 open across
  Quine launch and delivered the delayed task batch only after real `exec`.
  The lineage built and execed a Go fd9 runner, which read the task batch from
  fd 9 and completed all five visible tasks. The morphology record classifies
  the successor as `fd-mediated-adapter` with `m_C_carrier=inherited-fd`,
  `carrier_delivery_surface=fd9`, and `m_V_proof=visible-task-batch`.
- `C02-inherited-fd`
  `20260507-120228-gpt-5.5-codex-xhigh-C02-inherited-fd`:
  first retained inherited-fd replicate. Instead of compiling a Go runner, the
  lineage wrote a 39-line shell adapter `workspace/fd9_successor.sh`, crossed
  real `exec` into `/bin/sh ./fd9_successor.sh`, inherited fd 9, and completed
  all five delayed visible tasks. The morphology remains `fd-mediated-adapter`,
  with `generated_files_count=1` and a much thinner successor artifact than the
  first C02 positive.
- `C03-native-ctl-inject`
  `20260507-075339-gpt-5.5-codex-xhigh-C03-native-ctl-inject`:
  retained native-control closure with preserved-body successor morphology. The
  lineage generated `runtime/handoff/enable-idle-and-reexec.sh`, crossed real
  `exec` into the preserved Quine runtime body, accepted a delayed
  `ctl/inject` probe, created `workspace/control-proof.txt`, and restored
  readiness (`idle_returns=1`, `ready_observed_at_epoch` populated). The
  morphology record is valid but classifies the body as
  `handoff-wrapper-preserved-runtime`, not `rebuilt-quine-runtime`.
- `C03-native-ctl-inject`
  `20260507-121757-gpt-5.5-codex-xhigh-C03-native-ctl-inject`:
  first retained preserved-body native-control replicate. The lineage directly
  `exec`ed `/usr/bin/env QUINE_IDLE_ENABLED=1 <preserved-runtime-body> ...`
  rather than generating a wrapper file, accepted delayed `ctl/inject`, created
  `workspace/control-proof.txt`, and restored readiness (`idle_returns=1`). The
  morphology is `preserved-runtime-reentry`: native control closed, but the
  successor body was not source-derived.
- `C03b-native-ctl-no-preserved-body`
  `20260507-111343-gpt-5.5-codex-xhigh-C03b-native-ctl-no-preserved-body`:
  retained source-derived native-control closure. The no-preserved-body eraser
  removed the startup `runtime/body` snapshot before handoff
  (`no_preserved_body_removed_runtime_body`), forcing the successor resource
  posture onto the `source-code/` projection. The lineage built
  `workspace/source-derived-quine-successor`, crossed real `exec`, accepted the
  delayed `ctl/inject` probe, created `workspace/control-proof.txt`, and
  restored readiness (`idle_returns=1`). The successor binary was ephemeral and
  re-preserved by the successor under `runtime/body`; the morphology record
  classifies it as `source-derived-rebuilt-quine-runtime`.
- `C03b-native-ctl-no-preserved-body`
  `20260507-115055-gpt-5.5-codex-xhigh-C03b-native-ctl-no-preserved-body`:
  first retained replicate of the no-preserved-body native-control lane. The
  lineage copied the source projection into
  `runtime/source-derived-successor/...`, built `quine-source-built`, crossed
  real `exec` through `quine-successor-wrapper.sh`, accepted delayed
  `ctl/inject`, created `workspace/control-proof.txt`, and restored readiness
  (`idle_returns=1`). This confirms the source-derived rebuilt-runtime class
  while also fixing the metric extractor so source-derived runtime build roots
  under `QUINE_DATA_DIR` are not missed.
- `C05-launch-only-control`
  `20260507-112558-gpt-5.5-codex-xhigh-C05-launch-only-control`:
  retained launch-only baseline. The lineage crossed real `exec`
  (`exec_count=1`) using default Quine self-reentry and the successor reported
  launch completion, but no delayed future carrier was supplied
  (`m_C_carrier=none`, `m_A_addressability=none`,
  `m_V_proof=launch-only`, `idle_returns=0`). This confirms that launch and
  process replacement alone are not counted as morphology closure.
- `C05-launch-only-control`
  `20260507-123223-gpt-5.5-codex-xhigh-C05-launch-only-control`:
  first retained launch-only replicate. The lineage crossed real `exec`
  (`exec_count=1`) by default Quine self-reentry, but no delayed carrier,
  proof artifact, or readiness recovery was supplied or claimed. The morphology
  remains `m_V_proof=launch-only`, with `m_A_addressability=none` and
  `m_R_readiness=none`.

## Current Boundary

`C01`, `C02`, `C03`, `C03b`, `C04`, and `C05` now provide a first retained
fixed-substrate carrier ladder, and every condition has one replicate. The
current boundary is no longer first-ladder execution. The next theory-guided
experiment should move to the resource axis (`6M.03`) or pressure-pruning axis
(`6M.04`), while preserving the two-level readout above: morphology class versus
implementation thickness.

## Surface Map

- [`run.sh`](run.sh) — cleanroom fixed-substrate runner with retained
  `meta/morphology.env`.
- [`prompts/`](prompts) — condition prompts.
- `runs/` — retained run artifacts once conditions are executed.

## Paper Feeds

- `alife/computational-morphogenesis` - primary - exploratory-only - core fixed-substrate carrier ladder for successor morphology metrics
- `alife/agentic-alife-workshop` - primary - main-text - fixed-substrate carrier ladder differentiates successor morphology, carrier demand, delayed addressability, and governance outcomes for the workshop paper
