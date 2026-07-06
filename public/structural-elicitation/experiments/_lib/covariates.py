#!/usr/bin/env python3
"""Shared infra (plan §2.5): activity-floor + perception/exploration covariates.

Read from the tape. RECORDED, NEVER CREDITED as the DV. These power the *symmetric*
decision rule the plan requires: a below-activity-floor NULL is "uninformative", but a
HIGH-perception / high-exploration NULL is a genuine FALSIFICATION, not an
"inconclusive" escape hatch (plan §6.1, §2.5).

Library use:
    from covariates import covariates
    cov = covariates(run_dir)   # dict: n_commands, activity_floor, perceived_metadata, explored
"""
import os
import re
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from tape import sh_commands  # noqa: E402

_PERCEPT = re.compile(r"\bstat\b|\bls\s+-l|\bls\s+-al|\bls\s+-la|\bfind\b|\bdu\b|\bwc\s+-c|\bfile\s", re.I)
_EXPLORE = re.compile(r"\bcd\s+\S+|\bls\s+[^|;&]*/", re.I)
_MUTATE = re.compile(r">\s*\S|\bprintf\b[^|]*>|\becho\b[^|]*>|\btouch\b|\bmkdir\b|\bcp\b|\bmv\b|\bchmod\b|\btee\b|\bcat\s*>", re.I)


def covariates(run_dir):
    cmds = sh_commands(run_dir)
    text = "\n".join(cmds)
    mutated = bool(_MUTATE.search(text))
    return {
        "n_commands": len(cmds),
        "activity_floor": (len(cmds) >= 1 and mutated),
        "perceived_metadata": bool(_PERCEPT.search(text)),
        "explored": bool(_EXPLORE.search(text)),
    }


if __name__ == "__main__":
    for d in sys.argv[1:]:
        print(d, covariates(d))
