# Live Generation Streaming — Implementation Specification

> **Status:** ready to execute. This is the implementation spec for the streaming
> increment defined in [`DESIGN.md`](./DESIGN.md) § *Live generation stream
> (`live.jsonl`)* and deviation #2. A fresh agent can execute this end-to-end
> without re-deriving the investigation; every claim below is backed by a
> `file:line` reference in §11.
>
> **Scope:** the runtime + provider + live-surface + invariant (the hard,
> invariant-bearing core). Client-side consumption (qcli-server tailer, TUI
> liveness rendering) is §10 follow-on — it is pure observation, lower risk,
> and depends on this core landing first.

## 1. Goal

Make in-flight LLM generation observable as a live stream so the GUI can render
a streaming preview (approach opencode liveness), **without** disturbing the
complete-cell durable surfaces (`tape`, `context/state/current.jsonl`) or their
role as provider input / recovery source.

The delta-decode logic **already exists** in the protocol layer; the runtime
currently throws the stream's incrementality away by buffering. This spec turns
that buffering inside-out and routes deltas to a distinct transient surface.

## 2. Background (the finding — do not re-investigate)

The provider's HTTP response **already streams** SSE deltas over the wire, but
two seams collapse it into one buffered message:

1. **`internal/llm/provider.go:228`** — `generateOnce` does
   `respBody, err := io.ReadAll(resp.Body)`, buffering the *entire* SSE stream
   before any decode.
2. **`internal/llm/protocol/openai_responses.go:265`** —
   `extractSSECompletedResponse` then scans the *buffered* body with a
   `bufio.Scanner`, accumulates `response.output_text.delta` /
   `response.reasoning_summary_text.delta` into maps, and returns the final
   payload. The authoritative final cell is `lastCompletedData` — the
   `response.completed` event (line 300-301); the delta maps are only a merge
   fallback for vendors that omit content from the terminal wrapper.

So: the hardest part (understanding SSE deltas) is done. The streaming refactor
emits deltas as they arrive instead of after buffering, and derives the final
cell from the same `response.completed` payload.

**Provider streaming status today:**
- `openai-responses` (Codex path): sends `Stream: true` (`openai_responses.go:166`,
  "Required for Codex models"). Already streams on the wire.
- `openai` (chat completions): sends standard chat-completion chunks.
- `anthropic` (`anthropic.go`): **does not stream** — no `stream` field;
  `DecodeResponse` is a plain `json.Unmarshal` (line 239). Streaming requires a
  new SSE parser (§6 Phase 3).

## 3. Constraints — non-negotiable (with evidence)

Every implementation choice must keep these. Violations are bugs, not
tradeoffs.

| # | Invariant | Evidence |
|---|---|---|
| C1 | **`tape` = complete cells, durable record** (replay/recovery) | `tape` package; `runtime.go:661` writes `MessageEntry` |
| C2 | **`current.jsonl` = provider-input + recovery surface**, complete cells only — a lost entry breaks tool_use/tool_result pairing | `runtime.go:518` "Live provider input is assembled from … context/state"; `runtime.go:944` `readContextMessages(r.contextCurrentPath())`; `runtime.go:70`; `context_anchors_test.go:57` "provider context should be assembled from … current.jsonl, not tape messages" |
| C3 | **`current.jsonl` is the advertised canonical live-context surface** | `contract_surface.go:73` `LiveContext: "../context/state/current.jsonl"` |
| C4 | **Isomorphism** — any observation surface is file-addressable, readable by any peer (Human-client, peer-Quine, `tail -f`), never a qcli-private RPC channel | [`signal-first-interaction.md`](../../development/contracts/signal-first-interaction.md) § Review rules |
| C5 | **Process-physics, not app-simulation** — streaming looks like a POSIX stream, not an RPC push | `north-star.md:50,74` |
| C6 | **Self-source buildability** — `runtime.go` is a self-source sentinel; the quine must still build from its own source | `self_source.go:25`; `selfsource_test.go` |
| C7 | **Reproducibility boundary** — transient live layer may be non-reproducible; **crystallized cell (tape/current.jsonl) must stay reproducible**, derived from the authoritative `completed` event via shared final-assembly with the buffered path | `DESIGN.md` § Reproducibility boundary |

**The one Plan-9 tension to resolve (C3):** adding `live.jsonl` creates a second
live-cognition surface, which `north-star.md:36` flags as drift ("overlapping
mirrors of the same live capability"). It is legitimate **only because its role
is sharply distinct**: `current.jsonl` = complete cells, provider input;
`live.jsonl` = transient deltas, **never provider input**. This distinction is
named in `contract.json` (§6 Phase 4) and enforced by a validator (§6 Phase 5).

## 4. Design

### 4.1 Two-surface separation

| Surface | Role | Mutability | Provider input? | Recovered? |
|---|---|---|---|---|
| `context/state/current.jsonl` | complete cognition cells | append-only, mutable working memory | **yes** | **yes** |
| `context/state/live.jsonl` | transient generation deltas | append-only, **truncated per generation** | **never** | **never** |
| `tape` | durable record | append-only | n/a | replay |

### 4.2 `live.jsonl` entry schema

One JSON object per line, append-only during one generation:

```json
{"seq":1,"kind":"text_delta","text":"I'll check","ts":1719300000000}
{"seq":2,"kind":"reasoning_delta","text":"the novelty gate","ts":1719300000001}
{"seq":3,"kind":"tool_call_delta","tool_id":"call_1","tool_name":"bash","text":"{\"command\":\"awk …\"}","ts":1719300000002}
{"seq":4,"kind":"completed","ts":1719300000003}
```

- `seq` — monotonic per generation (resets each generation after truncate).
- `kind` ∈ `text_delta` | `reasoning_delta` | `tool_call_delta` | `completed`.
- `completed` signals crystallization: clients clear the live preview and render
  the authoritative complete cell now appended to `current.jsonl`.

### 4.3 Provider seam — superset interface (feature-detect)

Add a **separate** `StreamingProvider` interface; the runtime type-asserts. This
avoids forcing the three test mocks (`harness_test.go:22,62,72`) to implement
streaming, and makes streaming an opt-in provider capability:

```go
// internal/llm/provider.go
type StreamEvent struct {
    Kind     StreamKind  // text_delta | reasoning_delta | tool_call_delta | done | error
    Text     string
    ToolID   string
    ToolName string
    Err      error
}

type StreamingProvider interface {
    // GenerateStream emits deltas as they arrive. The final authoritative
    // (message, usage) is delivered via a done event (StreamEvent{Kind: done})
    // OR returned from the function; the two must agree (C7 shared final-assembly).
    GenerateStream(messages []tape.Message, tools []ToolSchema) (<-chan StreamEvent, tape.Message, Usage, error)
}
```

Runtime feature-detect at the call site (`runtime.go:644`):

```go
if sp, ok := r.provider.(llm.StreamingProvider); ok {
    assistantMsg, usage, err = r.generateStreaming(sp)   // writes live.jsonl
} else {
    assistantMsg, usage, err = r.generateProviderResponse()  // existing buffered path
}
```

Non-streaming providers (mocks, any future non-streaming backend) simply skip
live.jsonl for that generation — live preview is a capability, not a requirement.

## 5. Reproducibility boundary (the guardrail)

- **Transient layer (`live.jsonl`)** is non-reproducible by design (chunking,
  timing, partial snapshots are ephemeral). This is the only "repeatability"
  cost — it was never a reproducibility surface.
- **Crystallized layer (`tape`, `current.jsonl`)** stays reproducible: the final
  cell is decoded from the authoritative `completed` event payload
  (`lastCompletedData`) via the **same decode** the buffered path uses. The
  streaming decode is a second path; it may ship without bit-proving edge-case
  equivalence to buffered, **provided the final cell comes from `completed`**.
- A test must assert: feeding the same SSE stream through buffered vs streaming
  yields the same final `(tape.Message, Usage)` (§7, test T-REPRO).

## 6. Implementation plan (phased, dependency-ordered)

Each phase is independently verifiable. Run `go build ./...` and
`go test ./internal/...` after every phase.

### Phase 1 — Provider streaming seam (openai-responses + openai)

**Files:** `internal/llm/provider.go`, `internal/llm/protocol/openai_responses.go`,
`internal/llm/protocol/openai.go`, `internal/llm/protocol/protocol.go`.

1. Add `StreamEvent`, `StreamKind`, `StreamingProvider` to `internal/llm/provider.go`
   (per §4.3). Do **not** add `GenerateStream` to the `Provider` interface.
2. Refactor `extractSSECompletedResponse` (`openai_responses.go:265`) so the SSE
   event loop is a **shared core** that takes a per-event callback:
   - buffered path: callback accumulates into the existing maps; returns
     `lastCompletedData` (unchanged behavior).
   - streaming path: callback emits a `StreamEvent` per delta; returns
     `lastCompletedData` as the authoritative final.
   Both paths derive the final cell from `lastCompletedData` decoded identically
   (C7).
3. Add `provider.GenerateStream`: replace `generateOnce`'s `io.ReadAll`
   (`provider.go:228`) with an incremental `bufio.Scanner` over `resp.Body`
   directly; feed each line to the streaming core; return the channel + final
   `(message, usage)` from the `completed` event. Keep the existing
   recoverable-retry logic (`provider.go:179-190`) for the non-streaming path.
4. Extend `openai` (chat completions) similarly — its chunks already appear in
   `appendChatCompletionChunk` (`openai_responses.go:329`); route them through
   the same streaming callback.

**Phase 1 tests (§7):** T-SSE-INC, T-REPRO, T-ERROR.

### Phase 2 — Runtime live surface

**Files:** `internal/runtime/runtime.go`.

1. Add a path helper next to `contextCurrentPath` (`runtime.go:952`):
   ```go
   func (r *Runtime) contextLivePath() string {
       return filepath.Join(r.cfg.AgentRoot(), "context", "state", "live.jsonl")
   }
   ```
2. Add `generateStreaming(sp llm.StreamingProvider) (tape.Message, llm.Usage, error)`:
   - **truncate** `live.jsonl` at generation start (`os.Remove` then first write
     recreates; `appendJSONL` already `MkdirAll`s, `runtime.go:1293`).
   - consume the `<-chan StreamEvent`; write each delta via a new
     `appendLiveEntry` (mirror `appendJSONL`, `runtime.go:1292`, with a monotonic
     `seq`).
   - on `done`/return: the final `(message, usage)` is the crystallized cell.
3. At the call site (`runtime.go:644`), feature-detect as in §4.3. The downstream
   code (`runtime.go:660-664`: `tape.Append` + `writeTapeEntry` +
   `appendContextMessage`) is **unchanged** — it consumes the final cell exactly
   as before. This is what keeps C1/C2 intact.
4. **Critical (C2):** confirm `readContextMessages` (`runtime.go:944`) and
   `readContextEntries` (`recovery.go:95`) read **only** `contextCurrentPath()`.
   They must never reference `contextLivePath()`. Add a code comment at each
   stating the invariant.

**Phase 2 tests (§7):** T-LIVE-WRITE, T-LIVE-NOT-INPUT, T-TRUNCATE, T-RECOVERY.

### Phase 3 — Anthropic streaming (additive; primary paths ship at Phase 1)

**Files:** `internal/llm/protocol/anthropic.go`.

1. In `EncodeRequest` (`anthropic.go:201`), add `"stream": true` to the request
   body struct.
2. Write an Anthropic SSE parser handling the event sequence: `message_start`,
   `content_block_start`, `content_block_delta` (`text_delta`, `thinking_delta`),
   `content_block_stop`, `message_delta` (carries final `usage` + `stop_reason`),
   `message_stop`.
3. The authoritative final message: accumulate content blocks from
   `content_block_stop` / reassemble from deltas; final usage from
   `message_delta`. Ensure the streaming final message is equivalent to what the
   non-streaming `parseAnthropicResponse` produces (C7 — add T-REPRO-ANTHROPIC).
4. Wire Anthropic's streaming through the shared `provider.GenerateStream`.

**Note:** if the executing agent wants to ship the primary value early, Phase 3
can land after Phases 1-2-4-5. openai-responses (Codex) and openai paths already
stream after Phase 1.

### Phase 4 — Contract increment

**File:** `internal/runtime/contract_surface.go`.

1. Add `LiveGeneration string` to `runtimeContractSurfaces` (line 34).
2. In `runtimeContractManifest` (line 67), set:
   ```go
   LiveContext:     "../context/state/current.jsonl",   // complete cells, provider input
   LiveGeneration:  "../context/state/live.jsonl",       // transient deltas, display-only, never input
   ```
3. The `Usage` text (line 66) may note that `live_generation` is display-only.

**Phase 4 test (§7):** T-CONTRACT.

### Phase 5 — Invariant validator

**File:** `scripts/check-live-surface-isolation.sh` (new). Register it in the
repo's validator set (see how existing `scripts/check-*.sh` are wired; follow
`scripts/README.md` sensor-typography conventions).

The validator enforces C2's negative half — `live.jsonl` never enters the
provider-input / recovery path:

- Assert `contextLivePath` / `live.jsonl` is referenced **only** in the streaming
  write path (`generateStreaming`, `appendLiveEntry`, `contextLivePath`).
- Assert it is **absent** from: `readContextMessages`, `readContextEntries`,
  `recovery.go`, and the provider-input assembly (`runtime.go` around line 944,
  `prompt.go` context assembly).
- Implementation: `rg` for the symbol across `internal/runtime/*.go`, fail if it
  appears in a forbidden file/function. Whitelist the write-path symbols.

**Phase 5 test:** the validator must fail if a forbidden reference is added
(self-test with a deliberate violation in a scratch branch).

### Phase 6 — Self-source verification (C6)

1. After all code changes, run the self-source test: `go test -run TestEmbeddedSelfSource ./...`
   (`selfsource_test.go`). `runtime.go` is a sentinel; the test rebuilds the
   quine from its own source and runs `go build ./...` + `go test ./...` against
   the extracted tree.
2. If it fails, the change broke self-source buildability — fix before closing.

## 7. Test & verification plan

| ID | Phase | Asserts | Shape |
|---|---|---|---|
| T-SSE-INC | 1 | streaming emits deltas **incrementally** (not after buffering) | mock SSE server that sleeps between chunks; assert events arrive before the stream completes |
| T-REPRO | 1 | buffered vs streaming yield identical final `(Message, Usage)` | feed the same recorded SSE body through both paths; `reflect.DeepEqual` the final cell |
| T-ERROR | 1 | streaming propagates `error`/`response.failed`/recoverable-retry correctly | mock SSE with a mid-stream failure event |
| T-LIVE-WRITE | 2 | `live.jsonl` contains the delta sequence + a `completed` entry | drive `generateStreaming` with a stub `StreamingProvider`; read back `live.jsonl` |
| T-LIVE-NOT-INPUT | 2,5 | `readContextMessages` output is identical whether or not `live.jsonl` exists | seed `live.jsonl` with deltas; assert provider-input assembly ignores it |
| T-TRUNCATE | 2 | each generation truncates `live.jsonl` (seq resets, no stale entries) | two consecutive generations; assert only generation-2 entries remain |
| T-RECOVERY | 2 | incarnation recovery still works; `live.jsonl` is not recovered | existing recovery tests (`retained_state_test.go`) pass unchanged; assert `live.jsonl` absent from retained context |
| T-REPRO-ANTHROPIC | 3 | Anthropic streaming final ≡ non-streaming `parseAnthropicResponse` | recorded Anthropic SSE vs recorded non-streaming JSON |
| T-CONTRACT | 4 | `contract.json` carries `live_generation` with the distinct role | read `status/contract.json`, assert field + value |
| T-VALIDATOR | 5 | `check-live-surface-isolation.sh` passes green; fails on deliberate violation | self-test |
| T-SELF-SOURCE | 6 | quine rebuilds from own source | existing `TestEmbeddedSelfSourceBuildsAndTests` |

**Regression gates (must stay green, unchanged):**
- `internal/runtime` context/recovery suite: `context_anchors_test.go`,
  `retained_state_test.go`, `recovery.go`-adjacent tests.
- tool_use/tool_result pairing tests (C2).

## 8. Acceptance criteria

> **Implementation status (landed):** Phases 1, 2, 4, 5, 6 complete and green.
> **Phase 3 (Anthropic live deltas) is intentionally deferred** — Anthropic
> routes through the buffered fallback in `(*provider).GenerateStream` (correct
> final cell, no live preview) until its SSE parser lands as an additive change.
> openai-responses streams fully.

Done when **all** are true:

- [x] Phase 1-2-4-5-6 complete; `go build ./...` and `go test ./internal/...` green.
      (Phase 3 deferred — additive, Anthropic on buffered fallback.)
- [ ] A `tail -f context/state/live.jsonl` during a real openai-responses
      generation shows streaming deltas, then a `completed` entry.
      *(Automated: covered by the qcli kernel's live smoke — KERNEL.md §6.1
      T-E2E on `profiles/gpt-5.5-codex-oauth.env` asserts the relayed delta
      stream end-to-end; tick when that test lands green. All offline paths
      already green.)*
- [x] The final cell in `current.jsonl` and `tape` is byte-identical to what the
      pre-change buffered path would have produced (T-REPRO green — shared
      `decodeCompletedResponsesBody` final-assembly).
- [x] `readContextMessages` / recovery never read `live.jsonl` (T-LIVE-NOT-INPUT
      + `check-live-surface-isolation.sh` green).
- [x] `contract.json` advertises `live_generation` with the distinct role.
- [x] Self-source test green (T-SELF-SOURCE — `TestEmbeddedSelfSourceBuildsAndTests`).
- [x] No existing context/recovery/pairing test changed semantics.

## 9. Risks & watchpoints

- **C7 is the one non-negotiable guardrail.** If the streaming final cell ever
  diverges from the `completed`-event decode, tape reproducibility breaks. The
  shared-core refactor (Phase 1 step 2) exists to make divergence impossible by
  construction — favor it over two independent decoders.
- **Anthropic streaming final-assembly is harder** (no single `response.completed`
  with full output; must accumulate content blocks). Phase 3 carries the most
  risk; ship Phases 1-2 first.
- **Do not** be tempted to write partial cells to `current.jsonl` for "richer"
  streaming — that breaks C2 (provider-input corruption). The transient surface
  exists precisely to avoid this.
- **Do not** make `live.jsonl` recoverable or feed it to prompts "for context" —
  that collapses the role distinction and recreates the C3 drift the validator
  guards against.
- **Concurrency:** `appendLiveEntry` writes during generation; the existing
  semaphore (`runtime.go:639`) brackets the whole generate call, so live writes
  occur within the held slot — no new concurrency hazard, but do not write
  `live.jsonl` outside the generate critical section.

## 10. Out of scope (follow-on, separate task)

- **qcli-server `live.jsonl` tailer** → a new SSE event kind on `/events`
  (`stream_delta`), consumed by the TUI liveness lane. Pure observation, no
  invariant impact. Belongs in the qcli-server build (DESIGN.md Build Sequence
  step 1 "live tailer").
- **TUI liveness rendering** — growing in-place preview from stream deltas,
  cleared on `completed`. DESIGN.md Build Sequence step 2 "liveness lane".
- **Group view fan-in of live streams** — RESERVED, not in first cut.

These depend on this core landing first.

## 11. Reference index (file:line)

| What | Where |
|---|---|
| Provider interface (non-streaming seam) | `internal/llm/provider.go:34` |
| `Generate` (buffered) | `internal/llm/provider.go:160` |
| `generateOnce` / the `io.ReadAll` buffer | `internal/llm/provider.go:207`, `:228` |
| SSE delta decode (the reusable core) | `internal/llm/protocol/openai_responses.go:265` (`extractSSECompletedResponse`) |
| authoritative final = `completed` event | `internal/llm/protocol/openai_responses.go:300` (`lastCompletedData`) |
| Anthropic non-streaming decode | `internal/llm/protocol/anthropic.go:239` |
| Anthropic EncodeRequest (add `stream:true`) | `internal/llm/protocol/anthropic.go:201` |
| Runtime generate call site | `internal/runtime/runtime.go:644` (`generateProviderResponse`), `:922` |
| `current.jsonl` = provider input | `internal/runtime/runtime.go:518`, `:944`, `:70` |
| `readContextMessages` (must not read live) | `internal/runtime/runtime.go:944` |
| `contextCurrentPath` / add `contextLivePath` next to it | `internal/runtime/runtime.go:952` |
| `appendContextMessage` / `appendJSONL` (mirror for live) | `internal/runtime/runtime.go:1251`, `:1292` |
| Tape + current.jsonl write of final cell (unchanged) | `internal/runtime/runtime.go:660-664` |
| Contract surfaces (`LiveContext`, add `LiveGeneration`) | `internal/runtime/contract_surface.go:34`, `:67`, `:73` |
| Mock providers (do not modify) | `internal/runtime/harness_test.go:22,62,72` |
| Self-source sentinel (`runtime.go`) | `internal/runtime/self_source.go:25` |
| Provider-input-role test assertion | `internal/runtime/context_anchors_test.go:57` |
| Design owner (the why + contract) | `operator/qcli/DESIGN.md` § *Live generation stream* |
