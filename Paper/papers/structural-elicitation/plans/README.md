# Structural Elicitation — deep-research briefs + run ledger

Active deep-research query set for the structural-elicitation dossier. Goal: (1) confirm the
article's novelty/significance and (2) map adjacent theories — for the **current (post-2026-06-18)
theory**, including the parts that extend **beyond** the LBA. Workflow follows
[`../../../research.md`](../../../research.md) (`dr run --detach --max`, record ids,
retrieve with `dr get <id> --wait`, absorb into the dossier as `*-dr-report.md`).

Prior sweeps (2026-06-15, pre-reframe) live in the dossier root:
[`../prior-art-dr-report.md`](../prior-art-dr-report.md) (flat third-source claim) and
[`../matrix-novelty-dr-report.md`](../matrix-novelty-dr-report.md) (the matrix upgrade). This set
refreshes them onto the post-reframe framing and adds the geometry map + socio-linguistic gating
+ substrate-constraint thesis, none of which had a novelty pass.

## Briefs

| # | File | Targets |
|---|------|---------|
| 1 | [`brief-1-core-structural-elicitation.md`](brief-1-core-structural-elicitation.md) | validated core: environment-directs-frozen-agent, drive/target split, affordance-field, substrate-locus |
| 2 | [`brief-2-cartography-geometry.md`](brief-2-cartography-geometry.md) | beyond-LBA: two-grammar × stratum matrix + environment-as-geometry (5 meta-dimensions) |
| 3 | [`brief-3-sociolinguistic-gating.md`](brief-3-sociolinguistic-gating.md) | newest/riskiest: two-latch (epistemic × normative / scope-containment) gating; substrate-constraint thesis |

## Run ledger

Tool: `dr run --detach --save-meta runs/<name>.meta.json` then
`dr get --meta-file <meta> --wait --output dr-raw-N.md`; reports absorbed into the dossier as
`*-dr-report.md`.

**Findings on the dr tooling (2026-06-18):**
- Legacy Interactions API schema retired (May 2026); the CLI was updated locally to the 'steps' schema.
- A streaming-reconnect bug (stale `last_event_id` on reconnect → `400 invalid_request`) was fixed in `gemini-deepresearch-cli/src/api.ts` (fix installed; left uncommitted for review).
- **`--max` (deep-research-max) is broken SERVER-SIDE under the new schema** — *not* a CLI bug: the interaction stays `in_progress` but the server errors `invalid_request` while serializing its accumulated step state, on every retrieval channel (snapshot get + stream). Confirmed unrecoverable across ~40 min. The CLI now degrades honestly (stream fallback, refuses to write fake files).
- **Non-`--max` works reliably**, BUT `--output`/`output_text` **truncates the report head** (drops the first `steps[].content[].text` block — e.g. the novelty verdict / intro). **Recover the full report from JSON:** `dr get --meta-file <m> --json | jq -r '.steps[]|select(.type=="model_output")|.content[]|select(.type=="text")|.text'`. All three final reports here were reconstructed this way (32/38/36 KB vs the 16–19 KB truncated `output_text`).
- **Absorbed synthesis:** [`../novelty-adjacency-dr-report-2026-06-18.md`](../novelty-adjacency-dr-report-2026-06-18.md).

| Launch | Brief | Run | Status |
|---|---|---|---|
| 2026-06-18 | 1 core | `--max` `runs/brief-1.meta.json` | **DEAD — server-side, unrecoverable** (superseded) |
| 2026-06-18 | 2 geometry | `--max` `runs/brief-2.meta.json` | **DEAD — server-side, unrecoverable** (superseded) |
| 2026-06-18 | 3 gating | `--max` `runs/brief-3.meta.json` | **DEAD — server-side, unrecoverable** (superseded) |
| 2026-06-18 | 1 core | non-max `runs/test-nonmax-1.meta.json` | ✅ retrieved → `dr-raw-1-core.md` (32 KB, steps-reconstructed) |
| 2026-06-18 | 2 geometry | non-max `runs/nonmax-2.meta.json` | ✅ retrieved → `dr-raw-2-geometry.md` (38 KB, steps-reconstructed) |
| 2026-06-18 | 3 gating | non-max `runs/nonmax-3.meta.json` | ✅ retrieved → `dr-raw-3-gating.md` (36 KB, steps-reconstructed) |

**Retry `--max` later** (re-run the briefs with `--max`) once Google's deep-research-max
serialization is healed; until then non-max is the reliable path. Earlier bare-`--detach`-id
rounds also failed retrieval and were relaunched with `--save-meta`.
