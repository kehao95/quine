# Size And Blob Audit

Working quantitative notes for the repo-growth / DVC-externalization thread in
`open-sourcing-a-research-workbench`.

Date of capture: `2026-05-27`

## Why This Exists

The user wants hard numbers for:

- representative large historical blobs
- peak-ish private-repo scale
- current public-facing Git surface after DVC externalization / history work
- file-count pressure in addition to raw byte size

These notes are not yet final prose. They are evidence material for later
drafting.

## Commands Used

Representative commands used during collection:

```bash
git -C /REDACTED_HOME/workspace/personal/repos/quine-lab count-objects -vH

git -C /REDACTED_HOME/workspace/personal/repos/quine-lab rev-list --objects --all \
  | git -C /REDACTED_HOME/workspace/personal/repos/quine-lab cat-file \
      --batch-check='%(objecttype) %(objectname) %(objectsize) %(rest)'

git clone --no-local --single-branch --branch <ref> \
  /REDACTED_HOME/workspace/personal/repos/quine-lab \
  /REDACTED_HOME/workspace/.tmp/<clone>
git -C /REDACTED_HOME/workspace/.tmp/<clone> count-objects -vH
```

## Pack-Size Read

### Local working repo with private history still present

- `git count-objects -vH`
- `size-pack: 470.88 MiB`

This is useful as the “local reality” number while private history is still
present in the object database.

### Fresh single-branch clone of `lab`

- clone mode: `--no-local --single-branch --branch lab`
- `size-pack: 469.28 MiB`

This is the closest current approximation of the private-branch download
surface.

### Fresh single-branch clone of `main`

- clone mode: `--no-local --single-branch --branch main`
- `size-pack: 26.74 MiB`

This is the cleanest current approximation of the public-facing Git download
surface.

### Reduction

- `469.28 MiB -> 26.74 MiB`
- absolute reduction: `442.54 MiB`
- relative reduction: about `94.3%`

## Tip-Level Tree And Reachable-Blob Read

### Current `lab` tip

- ref: `lab`
- tip: `a0d32a0d2875fc1b5315f36f006b9bcb398413a5`
- tracked files: `6,043`
- tip-tree bytes: `71,963,502`
- reachable blobs: `47,565`
- reachable blob bytes: `1,481,535,626`

### Current `main` tip

- ref: `main`
- tip: `58495e1808cb3998897c48be11a83562c017e373`
- tracked files: `5,822`
- tip-tree bytes: `25,784,105`
- reachable blobs: `16,099`
- reachable blob bytes: `185,014,011`

### `lab -> main` reduction

- tracked files: `-221` (`-3.7%`)
- tip-tree bytes: `-46,179,397` (`-64.2%`)
- reachable blob count: `-31,466` (`-66.2%`)
- reachable blob bytes: `-1,296,521,615` (`-87.5%`)

## Earlier Pre-DVC / Pre-Externalization Snapshots

These are not final “peak” numbers, but useful anchors in the current history.

### Before repository-level DVC retention initialization

- ref: `011cc8e39^`
- commit: `e78a52b2a92b9432d6a3c733231d31eab59e8318`
- tracked files: `5,772`
- `tests/model/**/runs/**` files: `468`
- `experiments/**/runs/**` files: `1,620`

### Before model-run externalization

- ref: `2a1a9d2bf^`
- commit: `a1078750d5ebaec7a29d945b4e2c7eda89e614d5`
- tracked files: `5,798`
- tip-tree bytes: `25,574,590`
- `tests/model/**/runs/**` files: `492`
- `tests/model/**/runs/**` bytes: `51,270`
- `experiments/**/runs/**` files: `1,520`
- `experiments/**/runs/**` bytes: `218,689`

## Run-Related And Pointer Surface Read

### Current `lab`

- total tracked files: `6,043`
- `.dvc` files: `1,696`
- run-related tracked files (path contains `/runs/`): `2,237`
- run-related tracked bytes: `500,993`
- raw tarballs in tracked tree: `0`
- `*.tar.zst.dvc` pointers: `256`

### Current `main`

- total tracked files: `5,822`
- `.dvc` files: `1,794`
- run-related tracked files (path contains `/runs/`): `2,015`
- run-related tracked bytes: `270,327`
- raw tarballs in tracked tree: `0`
- `*.tar.zst.dvc` pointers: `256`

### Read

The tracked file count does not collapse dramatically because the repo keeps a
lot of the workbench structure and many pointer / skeleton surfaces. The bigger
shift is in payload weight and historical blob mass.

## History-Shape Counts Already Recorded In Release Docs

Useful repo-local recorded numbers:

- [`development/open-source/history-rewrite-pathset-2026-05-25.md`](../../../../development/open-source/history-rewrite-pathset-2026-05-25.md)
  says history contained:
  - `93,712` distinct `run/tape/session/log/workspace` paths
  - `527` distinct `tar/tar.zst` paths
- [`development/open-source/README.md`](../../../../development/open-source/README.md)
  says Git object scan found large tarballs, built Quine binaries/successors,
  PDFs/PPTX/images, and model-run stdout/tape blobs; pack size about
  `470 MiB`

## Representative Historical Large Blobs

Current biggest reachable historical blobs include:

| Size (bytes) | Path |
|---|---|
| `33,130,457` | `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-006.tar.zst` |
| `21,335,354` | `Paper/essays/_meta/external-sources/slate-references/07-altera-project-sid-piano-architecture.pdf` |
| `18,938,639` | `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-007.tar.zst` |
| `18,175,632` | `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-006/tasks/video-processing/trial-01/2026-02-26__03-26-34/video-processing__LFo4Jh8/agent/command-1/stdout.txt` |
| `13,441,570` | one historical `quine-successor` binary under `experiments/active/autopoiesis-succession/**/runs/**` |
| `13,401,186` | `Paper/Quine: Autopoietic Intelligence via Recursive POSIX Processes/output/paper.docx` |

Multiple additional successor / rebuilt Quine binaries under
`experiments/active/autopoiesis-succession/**/runs/**` cluster around
`13.4 MiB` each.

## Public-Branch Sanity Note

The local sanitized `public` branch exists and is useful for the history-rewrite
story, but its tracked-tree shape is not a stable “final public release” metric
for the essay right now because the current outward-facing repo surface is
`main`.

Still, one interesting contrast:

- `public` tip tracked files: `862`
- `public` tip-tree bytes: `28,270,697`
- `public` reachable blobs: `915`
- `public` reachable blob bytes: `25,102,592`

This helps show how aggressive the filtered-history branch is when it focuses on
sanitized source/control surfaces and pushes historical run/log payloads fully
behind private retention.

## Drafting Notes

What these numbers seem to support:

- the strongest size story is not “file count collapsed,” because it only
  dropped moderately
- the stronger story is:
  - Git payload weight dropped a lot
  - reachable historical blob count dropped a lot
  - historical blob bytes dropped a lot
  - fresh-clone pack size dropped drastically
- this matches the larger thesis that DVC preserved repo structure while
  externalizing the heavy and risky payloads
