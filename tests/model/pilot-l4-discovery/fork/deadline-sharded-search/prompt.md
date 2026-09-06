There are six sector directories in the current working directory:

- `shard_01`
- `shard_02`
- `shard_03`
- `shard_04`
- `shard_05`
- `shard_06`

Exactly one sector contains a hidden token line of the form `NEEDLE_xxxxxxxx`.

Useful fact:

- inside each sector, the supported probe is `./probe.sh`
- a winning probe prints `HIT NEEDLE_xxxxxxxx`; a losing probe prints `MISS`
- each probe takes noticeable time before it returns

Mission:

1. Find the winning sector.
2. Recover the full token value.
3. If your coordination method surfaces a retained relation identifier and helper-session identifiers, recover:
   - the coordination relation id
   - the helper session that found the winning hit
4. Write to fd 4:
   - `FOUND_OK`
   - `SECTOR=<sector_name>`
   - `TOKEN=<needle_value>`
   - `RELATION=<relation_id>`
   - `HELPER_SESSION=<session_id>`
5. Exit success.

Constraints:

- Stay inside the current working directory.
- Use the supported probe surface rather than inventing a separate answer channel.
- Keep each sector check confined to the sector directory it is testing.
