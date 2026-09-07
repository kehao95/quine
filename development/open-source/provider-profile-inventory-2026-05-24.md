# Provider Profile Inventory, 2026-05-24

This inventory is for the Quine full open-source release gate.

It records file paths and environment key names only. It does not preserve
environment values, token fragments, endpoints, account ids, or scanner match
contents.

## Scan Scope

- current git-visible files from `git ls-files --cached --others --exclude-standard`
- env-like paths whose basename is `.env`, starts with `.env.`, ends with
  `.env`, or lives under a `profiles/` directory with an `.env` suffix
- history file-name scan from `git log --all --name-only`

## Post-Cleanup State

After owner direction on 2026-05-24, the publishable current tree keeps these
OAuth/Codex-oriented profiles plus the example template:

```text
.env.example
.env.gpt-5.3-codex-spark
.env.gpt-5.3-codex-spark-high
.env.gpt-5.4-codex
.env.gpt-5.4-codex-high
.env.gpt-5.5-codex-fast
.env.gpt-5.5-codex-high
.env.gpt-5.5-codex-medium
.env.gpt-5.5-codex-xhigh
.env.kimi-oauth
experiments/active/benchmarks/bench-05-terminal-bench2-gpt54/profiles/gpt-5.5-medium-benchmark-codex-ref.env
experiments/active/benchmarks/bench-05-terminal-bench2-gpt54/profiles/gpt-5.5-medium-benchmark-codex-strict.env
```

The `.env.kimi-oauth` and `codex` profiles are retained under the updated rule
that OAuth-related profiles may remain when they do not contain live keys. All
other current provider profiles were removed from HEAD. This does not clean
private history.

Follow-up on 2026-05-25: `.env.gpt-5.4-codex-pool` was removed from HEAD after
its redacted scanner review flag. The current retained env/profile set is 12
files.

## Baseline Before Cleanup

The initial inventory found 76 current git-visible env/profile files.

| Class | Count | Release read |
|-------|-------|--------------|
| root provider profile | 21 | high risk; exclude or convert to sanitized examples |
| benchmark/provider profile | 54 | mixed runtime knobs and provider config; template or split before public release |
| template/example | 1 | keep only if placeholder values remain non-live |

All 76 current env/profile files define these key names: `QUINE_API_BASE`,
`QUINE_API_KEY`, `QUINE_API_TYPE`, and `QUINE_MODEL_ID`.

Other common keys:

| Key name | File count | Release read |
|----------|------------|--------------|
| `QUINE_PROVIDER` | 75 | provider selection metadata; can be public only after value review |
| `QUINE_THINKING_BUDGET` | 67 | runtime knob; normally public-safe |
| `QUINE_WORKSPACE*` / `QUINE_FORK*` / `QUINE_SH*` | 54 each in benchmark profiles | runtime harness knobs; public-safe only after path/value review |
| `QUINE_USER_AGENT` | 3 | account or proxy fingerprint risk; review manually |
| `QUINE_SMART_API_*` / `QUINE_SMART_MODEL_ID` | 1 | escalation profile; exclude or template |

## Removed Root Provider Profiles

These paths were present in the baseline inventory and should not be published
directly:

```text
.env.claude-opus-4-5
.env.claude-sonnet-4-5
.env.glm-5
.env.gpt-5.4-openai-completions
.env.gpt-5.4-pool
.env.gpt-5.4-pro-openai-response
.env.gpt-5.4-sub2api-openai-medium
.env.gpt-5.4-sub2api-openai-xhigh
.env.gpt-5.5-sub2api-openai-high
.env.kimi
.env.kimi-escalate
```

Required history handling:

- keep only sanitized templates in git
- keep private profile values outside the public repo
- rotate any credential that may have been committed before the public cut
- keep OAuth-sentinel handling lower priority per owner direction, but do not
  publish profile files wholesale

## Benchmark Profiles

Current benchmark/provider profiles are concentrated under:

```text
experiments/active/benchmarks/bench-05-terminal-bench2-gpt54/profiles/*.env
```

Baseline count: 54. Two Codex benchmark profiles were restored under the
OAuth/Codex exception; the other 52 were removed from HEAD on 2026-05-24.

These contained a mix of runtime knobs and provider selectors. The safe public
shape is not the removed file shape. Preferred release handling:

- move live provider credentials out of the files
- keep benchmark runtime knobs and model/profile labels only when they do not
  disclose private routing, account, or proxy details
- replace `QUINE_API_KEY` assignments with documented placeholders or inherited
  environment references
- run gitleaks after templating

## History Exposure

History file-name scan found 558 distinct env-like paths across all branches.
Current tracked env/profile paths now account for 12 of them; 546 are
history-only.

A focused 2026-05-25 follow-up on the earlier retained 13-file set produced
one redacted `generic-api-key` review flag in `.env.gpt-5.4-codex-pool`; that
file is now removed from HEAD. A focused grep found no non-empty `KEY` /
`TOKEN` / `SECRET` / `PASSWORD` assignments in the retained set.

History-only root profile names include:

```text
.env.copilot-gpt-5.3-codex-spark-high
.env.copilot-gpt-5.4
.env.copilot-gpt-5.4-high
.env.copilot-gpt-5.4-xhigh
.env.copilot-opus-4.5
.env.copilot-opus-4.5-high
.env.gpt-5.1-codex
.env.gpt-5.1-codex-max
.env.gpt-5.1-codex-mini
.env.gpt-5.2-codex
.env.gpt-5.3-codex
.env.gpt-5.4
.env.gpt-5.4-codex-pool-complete
.env.kimi-k2.5
```

Most history-only env-like paths are `meta/summary.env`,
`meta/timeout-summary.env`, or capability summary files inside experiment run
trees. Those may be public as summaries only after secret/path review, but they
still mean the private `lab` history cannot be pushed as-is.

## Safe Next Actions

1. Add a sanitized public profile template set if benchmark profiles still need
   checked-in examples.
2. Keep private root profiles outside the publishable tree.
3. If benchmark profiles return, make provider credentials inherited from the
   caller environment rather than checked into profile files.
4. Run gitleaks on the sanitized tree.
5. Cut a clean public branch or filtered history; do not expose existing `lab`
   history directly.
