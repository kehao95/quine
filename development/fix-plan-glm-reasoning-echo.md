---
title: Fix plan — GLM/z.ai reasoning_content echo causes HTTP 400 code 1214
status: implemented (code + unit tests green; live §6.2 GLM re-run still pending — owner will run)
kind: runtime fix-plan
owner: handoff (any agent)
blocks: Paper/predictions/open/reflexive-substrate-loop-P6.md (and any GLM run of the null-05 / self-maintaining-workflow test)
created: 2026-06-25
---

# Fix plan: GLM/z.ai `reasoning_content` echo → HTTP 400 code 1214

> **Implementation note (2026-06-25):** §4 edits landed verbatim
> (`quirks.go` `stripAssistantReasoning` field + `z.ai` rule;
> `protocol.go` `RequestOptions.StripAssistantReasoning`;
> `provider.go` both call sites; `openai.go` conditional echo via
> `convertOpenAIMessages(..., stripAssistantReasoning)`). §5 tests added
> (`openai_test.go` echo present/absent round-trip; `quirks_test.go` z.ai
> sets the flag, DeepSeek/Moonshot/OpenAI do not). §6.1 done:
> `go build ./...`, `go test ./internal/llm/... ./internal/llm/protocol/...`,
> and `go vet ./internal/llm/...` all green. §6.2 (live GLM re-run) and §6.3
> (DeepSeek smoke) intentionally deferred to the owner.

A self-contained, precedent-following fix. Root cause is confirmed both by
reproducing the failure and by comparing with the `../opencode` harness. An
executing agent can take this end-to-end without session context.

## 1. The bug (symptom)

GLM-5.2 (z.ai coding plan, OpenAI-compatible endpoint) runs die at **~turn 7**
with:

```
openai API error (HTTP 400): {"error":{"code":"1214","message":"The messages parameter is illegal. Please check the documentation."}}
exit_code=1 termination_mode=exit tokens_in≈98386
```

Reproduced run (both agents, process_exit before wallclock):
`experiments/active/null-ecology/null-05-reflexive-self-embedding/runs/20260625-182522-glm-5.2-zai-max.env-2p-p6-aaai-dossiers-reflexive/`.
See the session tapes (turn-7 `outcome` with code 1214) and
`meta/` (stop_reason=process_exit, no durable output, no commits).

Input tokens accumulate ~14K/turn: 16K → 22K → 25K → 32K → 98K (`context_used_k`
is set from `InputTokens/1000` at `internal/runtime/runtime.go:785`, so the
growth is real). At ~98K the z.ai request is rejected as "illegal."

## 2. Root cause

`internal/llm/protocol/openai.go:243` unconditionally echoes every prior
assistant turn's `reasoning_content` (GLM's xhigh CoT) back into outgoing
request messages:

```go
case tape.RoleAssistant:
    msg := openaiMessage{
        Role:             "assistant",
        Content:          m.Content,
        ReasoningContent: m.ReasoningContent,   // <-- BUG: echoed into every request
    }
```

The runtime's memory manager tracks/trims only visible content (~32K); it does
not account for the reasoning it re-sends, so reasoning accumulates untrimmed
and the payload hits ~98K, which z.ai rejects.

**Confirming evidence:**
- The tracked-context (~32K) vs sent-payload (~98K) gap = accumulated reasoning.
- `internal/tools/fork.go:2248` and `internal/tools/exec.go:361` **already** set
  `msg.ReasoningContent = ""` for sub-agent contexts. The developers already know it
  should not propagate, but the main request loop does not strip it.
- No test requires reasoning in *requests* (`openai_test.go:312,346` only assert
  reasoning *parsed from responses*).

## 3. The constraint (do NOT globally drop)

`../opencode` confirms `reasoning_content` is **not** globally echoed. Only
models marked `interleaved` echo it (`packages/opencode/src/provider/transform.ts:287-313`,
backed by `provider.ts:966-970`, `field: ["reasoning"|"reasoning_content"|"reasoning_details"]`).
**GLM is `interleaved: false`** in opencode's path, so opencode does NOT echo it
for GLM, matching the fix below.

BUT **DeepSeek requires reasoning on assistant messages** (opencode
`transform.ts:268-274`), and **quine's DeepSeek runs work** with reasoning echoed
(the 6p-minimal focal run was DeepSeek V4 Pro × 200 turns, no 1214). So the fix
**must not change DeepSeek**. It must be **provider-scoped**, not a global drop.

## 4. The fix (provider-scoped strip, following the existing `NoAssistantPrefill` precedent)

The codebase already has the exact pattern to copy: `RequestOptions` carries a
quirk-derived flag (`NoAssistantPrefill`) plumbed from `p.quirks` at
`provider.go:214` and `321`. Do the same for reasoning.

### Edits

1. **`internal/llm/quirks.go`**. Add a field and a z.ai rule.
   - Add `stripAssistantReasoning bool` to the `endpointQuirks` struct (near
     `noAssistantPrefill`).
   - In `quirksFor(apiBase, apiKey string)`, set it for z.ai:
     ```go
     if strings.Contains(base, "z.ai") {
         quirks.stripAssistantReasoning = true
     }
     ```
     (The GLM profile's base is `https://api.z.ai/api/coding/paas/v4`; matching
     the host is robust. If you prefer, also cover the `/paas/v4` suffix.)

2. **`internal/llm/protocol/protocol.go:25`**. Add the flag to `RequestOptions`:
   ```go
   StripAssistantReasoning bool
   ```
   (mirror `NoAssistantPrefill`.)

3. **`internal/llm/provider.go`**. Plumb the quirk into opts at **both** call
   sites (`provider.go:211-214` and `318-321`):
   ```go
   opts := protocol.RequestOptions{
       ...
       NoAssistantPrefill:        p.quirks.noAssistantPrefill,
       StripAssistantReasoning:   p.quirks.stripAssistantReasoning,
       ...
   }
   ```

4. **`internal/llm/protocol/openai.go:239-244`**. Conditionally include
   reasoning (this is the load-bearing change):
   ```go
   case tape.RoleAssistant:
       msg := openaiMessage{
           Role:    "assistant",
           Content: m.Content,
       }
       if !opts.StripAssistantReasoning {
           msg.ReasoningContent = m.ReasoningContent
       }
       for _, tc := range m.ToolCalls { ... }  // unchanged
   ```
   Keep parsing `reasoning_content` from responses (`openai.go:400-411`) and
   storing it in the tape. Only the *outgoing* echo becomes provider-scoped.

## 5. Tests

Follow the existing Kimi precedent (`openai_test.go:175
OmitsKimiThinkingForGenericOpenAI`, `:199 KeepsKimiThinkingForKimiModels`).

- **`internal/llm/protocol/openai_test.go`**: add a test that `EncodeRequest`
  with `opts.StripAssistantReasoning = true` produces assistant messages with no
  `reasoning_content` field, and with `= false` includes it. (Round-trip a
  tape.Message with Role=assistant + ReasoningContent set through EncodeRequest,
  unmarshal the JSON, assert field presence/absence.)
- **`internal/llm/quirks_test.go`**: assert `quirksFor` sets
  `stripAssistantReasoning == true` for a `https://api.z.ai/...` base and `false`
  for a DeepSeek/generic base.

Run:
```
go test ./internal/llm/... ./internal/llm/protocol/...
```

## 6. Verification / acceptance

1. Unit tests above pass; full `go test ./internal/llm/...` green.
2. Build quine in a worktree (the `run.sh` pattern builds `./cmd/quine/` inside a
   `git worktree add`) and **re-run the P6 GLM mission**. It must survive well
   past turn 7 (no 1214) and produce real output. Re-run command (from
   `experiments/active/null-ecology/null-05-reflexive-self-embedding/`):
   ```
   CFG_TAG=p6-aaai-dossiers CFG_MISSION='<see §7 mission text>' \
     ./run.sh profiles/glm-5.2-zai-max.env 30 1200 2
   ```
3. Optional but recommended: a short DeepSeek smoke (e.g.
   `profiles/deepseek-v4-pro.env 5 300 1`) to confirm no regression. Reasoning
   should still be echoed for DeepSeek (path unchanged).
4. The accumulated input-token growth (the 16K→98K trajectory) should collapse to
   roughly the working-context trajectory once reasoning is no longer echoed for
   GLM.

## 7. The P6 mission text (copy verbatim for the re-run)

```
You are running inside this repository, which is its own self-model. The workspace you see contains your own runtime source (cmd/quine/ — the Go binary that IS you), the theory that describes the system you are (Paper/philosophy/, Paper/theory/), the AGENTS.md constitution that governs how you operate, and experiments that study the very process you are enacting right now. Some files here are literally studies of agents like you. Read AGENTS.md first; it is this repo s self-model and your operating manual. TASK: produce two distinct AAAI-27 main-track (7-page) theory-paper dossiers under Paper/papers/, each on a different theoretical angle grounded in this repo s theory objects. For each dossier follow the convention in Paper/papers/open-ended-ecology/substrate-coincidence/aaai-state-propagation/DOSSIER.md (claim, formal framing, related work, proposed experiment, honest bounds). Coordinate with any peer agents to pick two different angles and avoid duplication.
```

## 8. Scope bounds / what NOT to do

- **Do NOT globally drop** `ReasoningContent: m.ReasoningContent`. That breaks
  DeepSeek (which needs it and currently works).
- **Do NOT change response parsing**. `reasoning_content` must still be parsed
  from responses and stored in the tape (for inspection/anchoring). Only the
  *outgoing echo* is provider-scoped.
- **Do NOT widen scope**: this is a one-quirk, one-conditional fix. If a future
  provider needs interleaved reasoning round-trip (newer Z.AI `reasoning_details`,
  etc.), add it as another provider-scoped quirk modeled on opencode's
  `interleaved.field`, not as a new default.

## 9. Context this unblocks

- `Paper/predictions/open/reflexive-substrate-loop-P6.md`, the P6 experiment
  (does injecting self-model recognition via the mission window close the
  cognitive loop?). The failed run produced a tentative-negative Level-3b signal
  but was **inconclusive** because it died at turn 7. This fix makes P6 runnable
  to completion.
- The broader "self-maintaining / self-generating workflow" test: Quine
  incarnating into its own repo to produce paper dossiers, evaluated as
  generation + maintenance + experiment simultaneously.

## 10. Pointers

- Bug location: `internal/llm/protocol/openai.go:243`
- Quirk plumbing precedent: `internal/llm/quirks.go` (`quirksFor`),
  `internal/llm/provider.go:211-214` & `318-321` (`NoAssistantPrefill` flow),
  `internal/llm/protocol/protocol.go:25` (`RequestOptions`)
- Existing strip pattern (sub-agent contexts): `internal/tools/fork.go:2248`,
  `internal/tools/exec.go:361`
- Provider-specific reasoning precedent (Kimi): `internal/llm/protocol/openai_test.go:175,199`
- opencode reference: `../opencode/packages/opencode/src/provider/transform.ts:268-313`,
  `provider.ts:966-970`; `../opencode/packages/llm/src/protocols/openai-chat.ts:251-258`
- Failed run (evidence): `experiments/active/null-ecology/null-05-reflexive-self-embedding/runs/20260625-182522-glm-5.2-zai-max.env-2p-p6-aaai-dossiers-reflexive/`
- Re-run launcher: `experiments/active/null-ecology/null-05-reflexive-self-embedding/run.sh`
  (creates an isolated `exp/null-05/<RUNID>` worktree, builds quine, runs)
