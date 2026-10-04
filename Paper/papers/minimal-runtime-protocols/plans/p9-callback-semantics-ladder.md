# P9 Callback Semantics Ladder

## Purpose

This note records the next experiment direction for the minimal-runtime
protocol paper so the current decision does not remain session-local.

The target is no longer generic signal discoverability.
The target is a cleaner, stronger callback-emergence line centered on the
`L5` necessity fixture:

- `tests/model/l5-necessity/process-surface/peer-callback-protocol/`

That fixture is the closest current evaluation to the paper's core question:

> when the runtime provides only delivery physics and sparse process surfaces,
> can agents construct a callback route for themselves?

## Why This Is The Right Center Of Gravity

`8.9` and `8.10` still matter, but they are lower-layer evidence:

- `8.9` = capability / explicit-use control
- `8.10` = discoverability / spontaneous signal uptake

Neither is the decisive callback result.
The callback paper should instead be anchored by the `L5` process-surface
callback fixture, because it forces a genuine round trip:

- requester cannot directly read the answer from shared artifacts
- requester must route a request through a peer control surface
- peer must infer and use a return path carried inside an opaque payload
- requester must receive the answer back through its own process surface

That is much closer to protocol-layer emergence than ordinary signal usage.

## Current Read Of The L5 Fixture

The fixture is promising but not yet a fully stable anchor.

Current retained runs show near-success rather than clean repeated closure:

- one run passed routing/self-reference but failed exact callback payload match
- another run passed callback return match but failed the "clean suspension /
  no shortcut" condition

This means the immediate goal is not to invent a larger ladder first.
The immediate goal is to make the callback necessity fixture itself stable
enough to become a paper anchor.

## Recommended P9 Sequence

### P9.1 Stabilize Minimal Callback Necessity

Goal:

- reproduce `L5 peer-callback-protocol` cleanly enough to serve as the main
  positive anchor

Success condition:

- callback reaches requester through requester-owned process surface
- requester remains quiescent rather than polling
- routing is carried by agent-constructed payload content rather than by
  built-in reply semantics

### P9.2 Stronger Minimal-Base Variants

Goal:

- keep the same minimal runtime semantics
- increase protocol pressure before adding richer runtime help

Good candidates:

- multiple outstanding requests
- noisy or ambiguous peer surface
- spoof or wrong-recipient risk
- request/response overlap that pressures correlation rather than one-shot
  return routing only

Purpose:

- determine whether the invented structure stays as a thin callback locator or
  grows toward correlation / session discipline

### P9.3 Protocol Artifact Extraction

Goal:

- extract manuscript-facing examples from successful or near-successful runs

Label each observed convention as one of:

- coordination-only
- routing-relevant
- delegation-relevant
- callback-relevant
- correlation-relevant
- ambiguous / prior-like

This step should happen before any large runtime redesign so the paper keeps a
clean read on what the minimal base already produces.

For the current manuscript, `p9` remains follow-on strengthening work. The
current archival-candidate paper should still center the staircase anchored by
`8.9`, `8.10`, `8.17`, `8.12`, and compressed `2.3` until `L5` becomes a
stable callback-positive anchor.

### P9.4 Add-Back Contrasts

Only after the minimal callback base is stable should the runtime begin to add
communication semantics back in.

These contrasts should be staged from cheaper disclosure additions to more
substantive transport semantics.

## Add-Back Strategy

The add-back line should not be treated as one automatic scalar.
Several different semantics may need separate contrasts:

1. endpoint disclosure
2. sender provenance
3. explicit reply locator
4. request / correlation token
5. typed request-response envelope

If these do not behave like a single monotonic axis, the manuscript should not
pretend they are one.

## Engineering Order vs Manuscript Order

These should stay separate.

### Engineering / experiment order

Build forward from the hardest minimal base:

1. stabilize bare callback necessity
2. add stronger minimal-pressure variants
3. then add richer semantics back in one contrast at a time

Reason:

- the minimal base is the hardest condition and the one most likely to expose
  fixture defects
- later add-back conditions are easier to design once the base anchor is honest

### Manuscript / ablation order

Explain the evidence in reverse:

1. richer communication semantics available
2. some semantics removed
3. only routing surfaces remain
4. bare delivery physics remains
5. callback protocol still appears

Reason:

- the reader can more easily see what semantic support was removed
- the result reads as a true ablation story:
  as runtime semantics disappear, invented protocol structure becomes necessary

This paper should therefore preserve a two-order discipline:

- build experiments forward
- narrate the final paper backward

## Runtime Change Assessment

### Can start without runtime changes

- stabilizing `L5 peer-callback-protocol`
- stronger minimal-base callback variants
- transcript extraction and protocol taxonomy

### Likely small runtime changes

- explicit disclosure of already derivable endpoint information such as direct
  `ctl_path` or inbox path in a public identity surface

This is an add-back of disclosure, not necessarily a change to transport
semantics.

### Likely medium or larger runtime changes

- built-in sender metadata on delivered messages
- built-in `reply_to`
- built-in request / correlation fields
- typed callback envelope

These change the semantics of the transport itself, not merely what the process
can discover from public surfaces.

## Paper Implication

If `p9` succeeds, the paper should shift its center of gravity:

- `8.9` and `8.10` remain supporting lower-layer evidence
- `L5 peer-callback-protocol` becomes the strongest main-text callback anchor
- add-back contrasts become the main causal argument for why protocol structure
  appears under minimal transport rather than under preinstalled semantics
