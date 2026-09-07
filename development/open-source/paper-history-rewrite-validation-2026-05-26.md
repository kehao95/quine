# Paper History Rewrite Validation, 2026-05-26

This records the Paper-specific history rewrite candidate for the open-source
release gate.

## Candidate

- Rewrite workspace:
  `/REDACTED_HOME/workspace/.tmp/quine-paper-history-rewrite-20260526202229/repo`
- Imported local branch:
  `main-paper-dvc-candidate`
- Source branch:
  `main`
- Source commit at rewrite start:
  `0f74742d078863bc52b6939ffc5fd414ac0ccadd`
- Later source commit cherry-picked after rewrite:
  `a7e491d3d7ae2b1c654ad1b152841ca00e54e358`

The candidate branch is local only. It has not been pushed to `origin`,
`public`, or `backup`.

## Rewrite Result

Final Paper private-DVC manifest:
`/REDACTED_HOME/workspace/.tmp/quine-paper-history-rewrite-20260526202229/private-paper-dvc-manifest.tsv`

| Metric | Value |
|--------|-------|
| private-DVC manifest records | 153 |
| unique private-DVC objects | 116 |
| generated private-DVC object bytes | 102,489,920 |
| HEAD Paper `.dvc` pointers | 96 |
| public DVC manifest entries | 0 |
| clone pack size before Paper rewrite | 141.85 MiB |
| clone pack size after Paper rewrite | 54.46 MiB |
| reachable Paper blob bytes after rewrite | 53,002,279 |

The second filter pass broadened external-source PDF matching from the current
`Paper/essays/_meta/external-sources/**/*.pdf` path to historical
`Paper/**/_meta/external-sources/**/*.pdf`, because earlier history used
`Paper/TRACK-blogs/_meta/external-sources/`.

The third filter pass added `.gitignore` exceptions for `Paper/*/output/*.dvc`
so legacy `Paper/Quine:*/output/*.dvc` pointers remain visible to DVC even
when the original output directory is ignored.

## Validation

| Check | Result |
|-------|--------|
| denied raw Paper blob scan | `0` matching raw blobs |
| HEAD archive `gitleaks dir` | no leaks found; ~22.51 MB scanned |
| full-history `gitleaks git --log-opts=main` | no leaks found; 1,248 commits / ~44.09 MB scanned |
| private localstore object presence | all 116 manifest objects present with expected sizes |
| private localstore rsync dry-run | `0` files to transfer |
| public DVC manifest | unchanged / empty |

The attempted `dvc status -r localstore --all-commits` check was stopped after
several minutes without output. It was replaced by direct manifest-object
presence and size validation against the local private DVC object store plus
historical restore smoke tests.

## Restore Smoke

Restore validation used a fresh clone of the rewritten candidate and pulled
from the local private store at `/PRIVATE_ARTIFACT_STORE`.

| Case | Commit | Restored output | Size |
|------|--------|-----------------|------|
| current external-source PDF | `a65188cb439d` | `Paper/essays/_meta/external-sources/slate-references/01-rlm-recursive-language-models-paper.pdf` | 8,250,587 |
| historical external-source PDF | `d28b5e74340b` | `Paper/TRACK-blogs/_meta/external-sources/slate-references/07-altera-project-sid-piano-architecture.pdf` | 21,335,354 |
| current presentation note | `a65188cb439d` | `Paper/archive/track-alife-presentations/materials/randall-beer-research-invitation/slides.md` | 5,111 |
| current archive PDF | `485c4b0c56f3` | `Paper/archive/rejected-2024/autopoietic-intelligence-via-recursive-posix-processes/artifacts/Quine_Autopoietic_Intelligence_via_Recursive_POSIX_Processes.pdf` | 2,678,711 |
| current graph payload | `95f88a479caf` | `Paper/core/graphs/HostGuest.svg` | 34,315 |
| historical legacy paper output | `4c1208272fea` | `Paper/Quine: Autopoietic Intelligence via Recursive POSIX Processes/output/paper.docx` | 13,401,186 |

## Backup Remote

Primary private DVC retention is complete in the local private store.

Google Drive backup is not closed for this Paper rewrite. The command
`dvc push -r lab-artifacts-gdrive-backup --all-commits -q` was stopped after
more than 10 minutes with no output. The next pass should either let that
command run to completion or create a narrower backup-push path for the 116
manifest objects.
