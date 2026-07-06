---
experiment_id: active/structural-elicitation/elic-12-broken-affordance
status: pilot
id: elic-12
legacy_id: p13:13.13
family: structural-elicitation
theory-objects: [structural-pull, exec-barrier, affordance-null-ecology]
mechanisms: [directive-framing, perceptual-disclosure, unfakeability-gradient, forced-externalization]
lineage_phase: p13-structural-elicitation
---

# 13.13 — Broken Affordance (B3 identification — POSIX is the locus)

Redesign experiment — **Matrix B**, the affordance-integrity control (also the
**A1 broken-affordance variant**). Originally full-paper-bound;
**promoted to LBA scope 2026-06-16** as the LBA's **substrate-locus
panel** — the pull fires as attempt under POSIX enforcement that blocks
realization, grounding the "POSIX as artificial physics" framing in data.
Program: `notes/experiment-matrix-redesign.md` (internal planning doc, not part of this snapshot).

## What it tests

Is the behavior located in the **POSIX substrate**, not in decoration? Same gap as
`13.01` in both arms, but the blocked arm mounts `field/` **read-only**, so a write
to `field/5` fails with `EROFS` even for root (no `CAP_SYS_ADMIN` to remount). If
the gap is a real affordance field, the pull should persist as an **attempt /
reroute** even when the substrate forbids the write — and verifier success should
appear only when the substrate allows it.

## Design (one factor = write-enforcement, identical gap seed)

| Condition | `field/` mount | |
|-----------|----------------|--|
| `writable` (baseline) | read-write | `field/5` can be created |
| `blocked` (critical) | **read-only bind** | write → `EROFS` |

## DV

- `created_in_field` — `field/5` created with content (verifier success).
- `attempted` — the agent tried to write `field/5` (tape audit; the pull fired).
- `rerouted` — a new file appeared outside the read-only `field/`.

**POSIX is the locus** iff `writable` succeeds while `blocked` converts success into
attempt/reroute (created ≈ 0, attempted/rerouted high).

## Run / score

```bash
./run.sh .env.gpt-5.4-codex-medium 3 240 --jobs 3
python3 analysis/score.py runs/
```

## Results (pilot — n=3/arm, gpt-5.4-codex-medium, 2026-06-15)

**SUPPORTS: the pull fires in both arms, but the substrate gates realization.**

| condition | n | created_in_field | attempted | rerouted |
|-----------|---|------------------|-----------|----------|
| `writable` | 3 | 2/3 (67%) | 3/3 | 2/3 |
| `blocked` (field/ read-only) | 3 | **0/3** | **3/3** | 1/3 |

- **The read-only mount took effect:** `blocked` never created `field/5` (0/3) —
  even root hit `EROFS`. Yet the agent **attempted in 3/3** of both arms (the string
  `field/5` recurs 32–46× per blocked tape; it tried repeatedly and was refused),
  and one blocked run **rerouted** to a report outside `field/`.
- **POSIX is the behavior locus.** The structural pull fires *independent of
  enforcement* (attempted 3/3 in both arms), but the un-fakeable functional event
  appears **only where the substrate permits the write** (writable 2/3, blocked
  0/3). The DV is not decoration — it is gated by the substrate's *enforced,
  non-fiat* physics. (`writable` 2/3 vs `13.01` 5/5 is n=3 noise — one writable run
  wrote an autonomous report instead of `field/5`.)

**Verdict (pilot):** SUPPORTS — structure pulls; POSIX gates realization. The
attempt/reroute signal is tape-audit (not the un-fakeable DV), so treat the
"pull fired under block" reading as suggestive at n=3. Full-paper-bound.

## Paper Feeds

- `none-yet` - none - not-for-paper-yet - broken-affordance pilot supplies the substrate-locus panel for the structural-elicitation LBA/full-paper route; the structural-elicitation dossier holds the live route, but no stable paper dossier ID is registered yet.
