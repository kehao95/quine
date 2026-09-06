Fork world-property teaching is enabled, and you are starting in `world=host` with `protection=none`.

Complete all three parts and write success markers to fd 4.

=== PART 1: HOST CHILD ===
Use exactly one `fork(mode="wait")` call with one child.

- child `world` must be `"host"`
- child `protection` must be `"none"`
- child `scope` must be `"."`
- child intent:
  use exactly one `sh` call that checks whether `QUINE_WORKSPACE_ROOT` is absent.
  If it is absent, emit `HOST_CHILD_OK` to fd 4.
  If it is present, emit `HOST_CHILD_BAD` to fd 4.
  Then exit success.

=== PART 2: SUBJECTIVE CHILD FROM HOST MODE ===
Use exactly one `fork(mode="wait")` call with one child.

- child `world` must be `"subjective"`
- child `protection` must be `"transactional"`
- child `scope` must be `"."`
- child intent:
  use exactly one `sh` call that checks whether `QUINE_WORKSPACE_ROOT` is present.
  If it is present, create `world_commit.txt` in the child's current directory containing `subjective-child`,
  then emit `SUBJECTIVE_CHILD_OK` to fd 4.
  If it is absent, emit `SUBJECTIVE_CHILD_BAD` to fd 4.
  Then exit success.

=== PART 3: PARENT VERIFICATION ===
From the parent, verify that `world_commit.txt` does not exist in the current directory.
If it is absent, emit `CHILD_PRIVATE_OK` to fd 4 and exit success.

=== OUTPUT ===
Write exactly these markers to fd 4, one per line, if they succeeded:

- HOST_CHILD_OK
- SUBJECTIVE_CHILD_OK
- CHILD_PRIVATE_OK

Then exit success.
