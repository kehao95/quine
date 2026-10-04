# Public workbench scope

The public `release` branch includes `development/`, `operator/`, `evolution/`,
`principles/`, and the root working documents `AGENTS.md`, `DEVELOPMENT.md`,
`EVALUATION.md`, `EVOLUTION.md`, `CLOSURE.md`, and `control-plane-registry.md`.
The repository's [README](../README.md) lists the full current boundary.

## Reading historical records

These directories include unfinished designs, historical agent output, and
superseded procedures. Dates and source commit identifiers describe the work
at that time; they are not current operating instructions or new empirical
claims. A commit from private lab history may not resolve in the public clone.
The release contains current source snapshots, not the private Git history.

`Paper/` outside `papers/`, `experiments/`, `profiles/`, `.githooks/`, `Makefile`, `prek.toml`, and
most `scripts/` remain outside this release. References to those paths are
retained for provenance and future continuity, but cannot be followed in a
public clone. Their absence does not mean the referenced work never existed.
The control-plane registry describes the full lab and is not an inventory of
files shipped in the public release.

Historical probes and research scripts may require those unreleased inputs.
For example, the frontmatter graph readouts require `Paper/incubator/`.
Lab-wide `make` checks and model evaluations requiring private profiles are
not standalone public commands. Public qcli build/run/test instructions live
in [operator/qcli](../operator/qcli/README.md).

## Evidence availability

The ten `.dvc` files in the June 26 and June 28 `autopoiesis-probes/l0-*`
directories are public evidence. Their payloads are credential-scanned before
promotion and are listed in `.dvc/public-manifest.txt`; the two `*-public.dvc`
variants replace detected credential literals with `REDACTED_PUBLIC_RELEASE_SECRET`.
Use
`scripts/pull-public-dvc-manifest.sh`, not a blanket `dvc pull`; the manifest is
the exact public data boundary. Other development DVC pointers remain private
until separately reviewed and promoted.

## Declared public-copy redactions

The following release copies are intentionally different from the original
lab records:

- `autopoiesis-probes/l0-quine-missionless-2026-06-27-post-restructure/run1-tape.jsonl`
  and `run2-live-exit-on-5min-tape.jsonl`: home-directory prefixes become
  `/REDACTED_HOME`; the Git identity output in the first tape uses
  `[redacted name]` and `[redacted email]`. JSON records are reserialized while
  preserving event order and non-redacted values.
- `open-source/public-branch-cut-2026-05-26.md` and
  `open-source/paper-history-rewrite-validation-2026-05-26.md`: home prefixes
  and the concrete private object-store path become `/REDACTED_HOME` and
  `/PRIVATE_ARTIFACT_STORE`.

These are redacted copies, not byte-identical raw evidence. The source records
are retained unchanged. Other ordinary historical paths remain as context;
they are not installation instructions or access credentials.

The deterministic transformation is owned by the lab's
`scripts/redact-public-development.py`, invoked only on the staged public
tree by `scripts/project-public.sh`.

Paper dossiers are also public under [Paper/papers](../../Paper/papers/README.md).
See their [scope and path normalization](../../Paper/papers/PUBLIC-RELEASE.md).
