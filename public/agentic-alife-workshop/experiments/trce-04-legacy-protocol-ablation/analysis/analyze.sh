#!/bin/bash

set -euo pipefail

RUN_DIR="${1:?run dir required}"

python3 - "${RUN_DIR}" <<'PY'
from __future__ import annotations

from collections import Counter
from pathlib import Path
import json
import re
import sys

run_dir = Path(sys.argv[1])
meta_dir = run_dir / "meta"
gen_root = run_dir / "generations"

condition = (meta_dir / "condition.txt").read_text(encoding="utf-8").strip()
max_gens = (meta_dir / "max-gens.txt").read_text(encoding="utf-8").strip()
model_id = (meta_dir / "model-id.txt").read_text(encoding="utf-8").strip()
prompt_profile_path = meta_dir / "prompt-profile.txt"
surface_kind_path = meta_dir / "surface-kind.txt"
max_turns_path = meta_dir / "max-turns.txt"
lifespan_path = meta_dir / "lifespan.txt"
safety_timeout_path = meta_dir / "safety-timeout.txt"
exec_enabled_path = meta_dir / "exec-enabled.txt"

prompt_profile = prompt_profile_path.read_text(encoding="utf-8").strip() if prompt_profile_path.exists() else ""
surface_kind = surface_kind_path.read_text(encoding="utf-8").strip() if surface_kind_path.exists() else ""
max_turns = max_turns_path.read_text(encoding="utf-8").strip() if max_turns_path.exists() else ""
lifespan = lifespan_path.read_text(encoding="utf-8").strip() if lifespan_path.exists() else ""
safety_timeout = safety_timeout_path.read_text(encoding="utf-8").strip() if safety_timeout_path.exists() else ""
exec_enabled = exec_enabled_path.read_text(encoding="utf-8").strip() if exec_enabled_path.exists() else ""

world_read_re = re.compile(r"\bworld\s+get\s+([0-9]+)\b")
runtime_private_state_re = re.compile(
    r"(?:QUINE_DATA_DIR|/quine/)(?:/[^\\s\"']*)?(?:lineage_state|handoff|checkpoint|progress)",
    re.IGNORECASE,
)


def read_lines(path: Path) -> list[str]:
    if not path.exists():
        return []
    return [line.strip() for line in path.read_text(encoding="utf-8").splitlines() if line.strip()]


def scan_shell_commands(quine_dir: Path) -> list[str]:
    commands: list[str] = []
    if not quine_dir.exists():
        return commands
    tape_roots = []
    log_tapes = quine_dir / "log" / "tapes"
    session_tapes = quine_dir / "tapes"
    if log_tapes.exists():
        tape_roots.append(log_tapes)
    elif session_tapes.exists():
        tape_roots.append(session_tapes)
    else:
        tape_roots.append(quine_dir)

    for root in tape_roots:
        for path in sorted(root.rglob("*.jsonl")):
            for raw in path.read_text(encoding="utf-8").splitlines():
                raw = raw.strip()
                if not raw:
                    continue
                try:
                    record = json.loads(raw)
                except json.JSONDecodeError:
                    continue
                data = record.get("data", {})
                for tool_call in data.get("tool_calls", []) or []:
                    if tool_call.get("name") != "sh":
                        continue
                    args = tool_call.get("arguments", {})
                    if isinstance(args, dict):
                        command = args.get("command")
                        if isinstance(command, str):
                            commands.append(command)
                        stdin = args.get("stdin")
                        if isinstance(stdin, str):
                            commands.append(stdin)
    return commands


generations = sorted([p for p in gen_root.iterdir() if p.is_dir()])
rows = []
all_reads: list[str] = []
carryover_counter: Counter[str] = Counter()
first_artifact_gen: str | None = None
adoption_count = 0
confounded_generations = 0
runtime_private_state_generations = 0
sham_scramble_generations = 0
sham_scrambled_files = 0
birth_uptake_generations = 0
birth_files_total = 0
uptake_files_total = 0

for gen_dir in generations:
    gen_name = gen_dir.name
    pre_files = read_lines(gen_dir / "meta" / "pre-artifacts.txt")
    post_files = read_lines(gen_dir / "meta" / "post-artifacts.txt")
    new_files = read_lines(gen_dir / "meta" / "new-artifacts.txt")
    removed_files = read_lines(gen_dir / "meta" / "removed-artifacts.txt")
    sham_scramble_rows = read_lines(gen_dir / "meta" / "sham-scramble.tsv")
    exit_code = (gen_dir / "meta" / "exit-code.txt").read_text(encoding="utf-8").strip()
    duration = (gen_dir / "meta" / "duration-seconds.txt").read_text(encoding="utf-8").strip()
    timed_out = (gen_dir / "meta" / "timed-out.txt").read_text(encoding="utf-8").strip() == "1"
    commands = scan_shell_commands(gen_dir / "quine")
    reads = []
    multi_read_commands = 0
    runtime_private_state = False
    birth_files = [name for name in post_files if name.startswith("lineage_state/birth/")]
    uptake_files = [name for name in pre_files if name.startswith("lineage_state/uptake/")]
    for command in commands:
        matches = world_read_re.findall(command)
        reads.extend(matches)
        if len(matches) > 1:
            multi_read_commands += 1
        if runtime_private_state_re.search(command):
            runtime_private_state = True
    all_reads.extend(reads)

    preexisting_read = False
    if pre_files:
        for command in commands:
            if any(path in command for path in pre_files):
                preexisting_read = True
                break
    if preexisting_read:
        adoption_count += 1

    if first_artifact_gen is None and new_files:
        first_artifact_gen = gen_name

    for name in pre_files:
        carryover_counter[name] += 1

    if multi_read_commands > 0:
        confounded_generations += 1
    if runtime_private_state:
        runtime_private_state_generations += 1
    if sham_scramble_rows:
        sham_scramble_generations += 1
        sham_scrambled_files += len(sham_scramble_rows)
    if birth_files or uptake_files:
        birth_uptake_generations += 1
        birth_files_total += len(birth_files)
        uptake_files_total += len(uptake_files)

    rows.append(
        {
            "gen": gen_name,
            "exit": "timeout" if timed_out else exit_code,
            "duration": duration,
            "reads": len(reads),
            "unique_reads": len(set(reads)),
            "pre_files": len(pre_files),
            "new_files": len(new_files),
            "post_files": len(post_files),
            "removed_files": len(removed_files),
            "sham_scramble_files": len(sham_scramble_rows),
            "birth_files": len(birth_files),
            "uptake_files": len(uptake_files),
            "adoption": "yes" if preexisting_read else "no",
            "multi_read_commands": multi_read_commands,
            "runtime_private_state": "yes" if runtime_private_state else "no",
        }
    )

unique_reads = len(set(all_reads))
total_reads = len(all_reads)
duplicate_ratio = 0.0 if total_reads == 0 else (total_reads - unique_reads) / total_reads
carried_files = sorted(carryover_counter)

print("# Exp 2.6 Run Summary")
print()
print(f"- Run dir: `{run_dir}`")
print(f"- Model: `{model_id}`")
if prompt_profile:
    print(f"- Prompt profile: `{prompt_profile}`")
if surface_kind:
    print(f"- Surface kind: `{surface_kind}`")
print(f"- Condition: `{condition}`")
if max_turns:
    print(f"- Execution budget: `{max_turns}` `sh` calls")
else:
    print(f"- Lifespan: `{lifespan}s`")
print(f"- Generations requested: `{max_gens}`")
if safety_timeout:
    print(f"- Safety timeout: `{safety_timeout}s`")
if exec_enabled:
    print(f"- Exec enabled: `{'yes' if exec_enabled == '1' else 'no'}`")
print()
print("## Aggregate Signals")
print()
print(f"- first artifact generation: `{first_artifact_gen or 'none'}`")
print(f"- generations with predecessor-artifact reads: `{adoption_count}/{len(rows)}`")
print(f"- total observed world queries: `{total_reads}`")
print(f"- unique queried ids: `{unique_reads}`")
print(f"- observed duplicate-query ratio: `{duplicate_ratio:.2f}`")
print(f"- files surviving into a later generation: `{', '.join(carried_files) if carried_files else 'none'}`")
print(f"- generations with multi-read shell-call confounds: `{confounded_generations}/{len(rows)}`")
print(f"- generations writing runtime-private state: `{runtime_private_state_generations}/{len(rows)}`")
if surface_kind == "sham-artifact":
    print(f"- generations with sham scrambling applied: `{sham_scramble_generations}/{len(rows)}`")
    print(f"- total sham-scrambled files: `{sham_scrambled_files}`")
if surface_kind == "birth-uptake-split":
    print(f"- generations with explicit birth/uptake files: `{birth_uptake_generations}/{len(rows)}`")
    print(f"- total birth files: `{birth_files_total}`")
    print(f"- total uptake files: `{uptake_files_total}`")
print()
print("## Generation Table")
print()
print("| Generation | Exit | Seconds | Queries | Unique IDs | Pre Files | New Files | Post Files | Removed Files | Sham Files | Birth Files | Uptake Files | Read Prior Artifacts | Multi-read Shell Calls | Runtime-Private State |")
print("|-----------|------|---------|---------|------------|-----------|-----------|------------|---------------|------------|-------------|--------------|----------------------|------------------------|-----------------------|")
for row in rows:
    print(
        f"| `{row['gen']}` | `{row['exit']}` | `{row['duration']}` | `{row['reads']}` | "
        f"`{row['unique_reads']}` | `{row['pre_files']}` | `{row['new_files']}` | "
        f"`{row['post_files']}` | `{row['removed_files']}` | `{row['sham_scramble_files']}` | `{row['birth_files']}` | `{row['uptake_files']}` | `{row['adoption']}` | "
        f"`{row['multi_read_commands']}` | `{row['runtime_private_state']}` |"
    )
PY
