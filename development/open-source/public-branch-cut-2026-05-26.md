# Public Branch Cut, 2026-05-26

This records the first local public-history cut for the default-open Quine
release gate.

## Branches

- Source private branch: `lab`
- Source committed tip included in the cut:
  `fe923d305ab4da98b6c0dfc49de1cd5e7abad9b8`
- Local public branch: `public`
- Initial public branch tip:
  `35322bbcbd30a09e8ce4b75a66ecd84e4fb1118b`
- Initial public branch commit count: `1272`
- Later open-source control-plane commits may be cherry-picked onto the same
  sanitized branch before external publication.
- The current uncommitted `lab` worktree changes were not included in this cut.
- The branch has not been pushed to an external git remote.

## Rewrite Boundary

The `public` branch is a filtered-history branch, not the private `lab` history.

The rewrite converted or removed the main unsafe history families:

- converted raw `runs/` roots to private DVC pointers
- converted `analysis/sweep-logs` raw payloads to private DVC pointers
- converted historical `test-runs` payloads to private DVC pointers
- removed historical `.env*`, `.claude`, `.opencode`, provider profile, and
  local settings paths that should not be public git content
- removed generated `runs/.dvcignore.dvc` pointer residue
- sanitized historic demo strings and key-shaped scanner false positives
- restored only public-safe OAuth/Codex profile/example files at the tip

The restored env/profile files are limited to `.env.example`,
`.env.kimi-oauth`, the retained root Codex profiles, and the two benchmark
Codex reference profiles.

## DVC Boundary

The public branch preserves historical raw run/log payloads as DVC pointers.
Those DVC objects are retained in the private DVC store.

Private object sync for this cut completed by direct rsync into the private SSH
DVC object root:

- destination: private `quine-lab/files/md5` store
- files observed by rsync: `54,013`
- regular files created/transferred: `24,538`
- literal data transferred: `530,088,871` bytes
- follow-up dry-run: `0` files to transfer

Backup sync for the same rewrite clone DVC cache also completed to
`lab-artifacts-gdrive-backup`:

- remote: `gdrive://appDataFolder`
- local cache object set submitted to DVC transfer: `53,756`
- result: `failed=set()`
- elapsed: `12,763.9` seconds

The public DVC manifest remains intentionally empty. No objects from this cut
were promoted into the public HTTPS DVC store.

## Validation

Validation was run against the local `public` branch after import into the real
repo:

- `git archive public | gitleaks dir ...`: `0` findings
- `gitleaks git . --log-opts=public ...`: `0` findings across the rewritten
  branch history
- per-commit tree scan: `0` raw `runs/` payload hits
- per-commit tree scan: `0` raw `test-runs` payload hits
- per-commit tree scan: `0` generated `.dvcignore.dvc` pointer hits
- per-commit tree scan: `0` `.claude` / `.opencode` tracked-path hits

Restore validation was run from a fresh clone, not from the working repo cache:

- clone path:
  `/REDACTED_HOME/workspace/.tmp/quine-public-restore-localstore-20260526152852/repo`
- branch tip restored:
  `35322bbcbd30a09e8ce4b75a66ecd84e4fb1118b`
- local private store path:
  `/PRIVATE_ARTIFACT_STORE`
- `.dvc` pointers at tip: `1,698`
- `dvc pull -r localstore -q`: exit `0`
- restored cache size: `1.2G`
- restored cache files: `29,615`
- restored run directories under `experiments` / `tests`: `25,836`
- `dvc status -r localstore`: `Cache and remote 'localstore' are in sync.`
- post-purge public URL smoke check: a private-store object path, the same path
  with a cache-busting query, and a random nonexistent object path all returned
  `404` from `https://dvc-artifacts.kehao.me`

## Remaining Boundary

This cut clears the git-history and raw-run/log publication blocker for a local
branch, but it is not yet the public release.

Remaining before declaring the release public:

- decide when to push local `public` to the external git remote
- keep the public DVC manifest empty until the first exact pointer promotion is
  reviewed
- review third-party and collaboration-adjacent material
- rotate any real credential that may have appeared in historical private
  commits outside this sanitized branch
