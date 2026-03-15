# Git Hooks

This directory holds repo-tracked Git hooks for Quine.

Do not treat `.git/hooks/` as canonical policy. Hook logic that matters to the project should live here, be version-controlled, and be installed through `core.hooksPath`.

## Install

```bash
make install-githooks
```

That sets:

```bash
git config core.hooksPath .githooks
```

## Current Hooks

| Hook | Purpose |
|------|---------|
| `pre-commit` | protect model-scenario inventory/scorer coverage and public-surface doc drift |

## Intention Contract

Each tracked hook should state its intention in the hook file itself, not only here.

That intention should make four things explicit:

- what invariant the hook protects
- what files it is allowed to touch or validate
- what canonical script or control-plane source it depends on
- how a future agent should repair it if the implementation breaks

This is deliberate: a broken hook should be repairable from repo context alone, without needing chat history.

## Scope

The current pre-commit hook is intentionally narrow:

- if you stage `tests/model/README.md`, `run.sh`, `scenarios.toml`, or a prompt, it checks registry/prompt/scorer synchronization
- if you stage root public docs, it checks for stale public-branch wording and stale test-path references

It does not yet force repo-wide completion. The hook protects touched files while the broader rollout is still in `pending` / `partial`.

## Repair Rule

If the pre-commit hook fails unexpectedly, the default repair strategy is:

1. keep the policy intent intact
2. repair the hook or its called script
3. only relax scope if the canonical policy documents changed

Do not "fix" a broken hook by deleting it or by silently bypassing the underlying control-plane rule.
