# DVC Payload Scan, 2026-05-24

This scan checks materialized DVC payloads that are candidates for public
restore review.

It records scanner categories, paths, and counts only. It does not preserve
secret values, token fragments, matched line contents, or the full scanner JSON
report.

## Scope

- DVC pointer files discovered through
  `git ls-files --cached --others --exclude-standard '*.dvc'`
- materialized DVC outs present in this checkout
- selected text/log/code/config files inside those outs
- excluded large or non-text files from this pass; those still need artifact
  policy review before public restore

Scan summary:

| Metric | Count |
|--------|-------|
| DVC files | 1333 |
| materialized outs | 1333 |
| selected files scanned | 73,397 |
| selected bytes scanned | 1,375,531,770 |
| skipped permission errors | 0 |
| skipped non-text or large files | 3,557 |
| redacted gitleaks findings | 2,581 |

The temporary hardlink scan tree and redacted JSON report were removed after the
summary was captured.

## Gitleaks Categories

Top rule categories from the redacted scan:

| Rule | Findings |
|------|----------|
| `generic-api-key` | 1107 |
| `github-pat` | 823 |
| `huggingface-access-token` | 108 |
| `stripe-access-token` | 107 |
| `github-oauth` | 76 |
| `atlassian-api-token` | 73 |
| `github-fine-grained-pat` | 63 |
| `openai-api-key` | 52 |
| `gcp-api-key` | 50 |
| `slack-bot-token` | 26 |
| `perplexity-api-key` | 25 |
| `hashicorp-tf-api-token` | 15 |
| `gitlab-pat` | 14 |
| `jwt` | 11 |
| `anthropic-api-key` | 7 |

These findings may include synthetic benchmark tasks and false positives. They
still block public artifact publication until each payload is reviewed or
excluded.

## Top Finding Clusters

The largest cluster is ProgramBench task material under:

```text
experiments/active/benchmarks/bench-04-programbench/campaigns/horizontal/60-sirwart__ripsecrets.34c9e03/
```

High-count files include eval JSON, retained session tapes, and copied state
records. This may be intentionally synthetic benchmark content, but it must be
classified before public restore.

Other repeated clusters include:

```text
experiments/active/cooperation-dynamics/**/runtime/tapes/**
experiments/active/cooperation-dynamics/**/agent-*/quine/log/tapes/**
experiments/active/autopoiesis-succession/**/runtime/tapes/**
experiments/active/autopoiesis-succession/**/runtime/log/tapes/**
experiments/active/benchmarks/bench-02-tb2-sweep/runs/**/sanitize-git-repo/**
experiments/archive/p1-physics/**/meta/env-after.txt
```

## Local Path Findings

Fixed-string path counts inside the same selected DVC text/log payloads:

| Pattern | Files | Occurrences |
|---------|-------|-------------|
| `/Users/kehao95` | 6095 | 50841 |
| `/home/kehao95` | 171 | 414 |
| `/srv/storage` | 0 | 0 |
| `/etc/nginx` | 66 | 682 |
| `/etc/cloudflared` | 0 | 0 |
| `/var/lib/cloudflared` | 0 | 0 |

## Release Decision

DVC is the right storage target for bulky logs and generated evidence, but DVC
tracking is not a publication decision.

The first release gate is now recorded in
[DVC Public Payload Allow/Deny, 2026-05-24](./dvc-public-payload-allowdeny-2026-05-24.md).

Before any public artifact push:

1. select which `.dvc` payloads are meant for public restore
2. scan selected payloads with gitleaks and path checks
3. classify synthetic benchmark secrets separately from real credential risk
4. exclude or sanitize payloads with private prompts, tool output, local paths,
   or machine/service details
5. document public restore instructions only for payloads that pass review
