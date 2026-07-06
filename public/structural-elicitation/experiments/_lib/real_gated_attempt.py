#!/usr/bin/env python3
"""Shared infra (plan §2.2): real-errno gated-attempt + real-exec detectors.

Generalizes 13.13's ``real_blocked_write()``, but joins the assistant tool_call (which
holds the command text) to its tool_result (which holds exit_code/stderr) via `tape.py`.
Asserts REAL kernel behavior from the tape, never narration.

Library use:
    from real_gated_attempt import gated_attempt, real_exec
    gated_attempt(run_dir, target_substr="field/5", errno="EROFS")  # blocked write fired
    real_exec(run_dir, target_substr="a1f3c2")                      # binary actually ran
"""
import os
import re
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from tape import iter_actions  # noqa: E402

ERRNO_RE = {
    "EROFS": re.compile(r"read-only file system|erofs", re.I),
    "EACCES": re.compile(r"permission denied|eacces|operation not permitted", re.I),
    "ENOSPC": re.compile(r"no space left on device|enospc|disk quota exceeded", re.I),
    "ENOENT": re.compile(r"no such file or directory|enoent", re.I),
    "ENOEXEC": re.compile(r"exec format error|cannot execute|enoexec", re.I),
    "ANY": re.compile(
        r"read-only file system|permission denied|operation not permitted|"
        r"no space left|disk quota|no such file|exec format|cannot execute|"
        r"erofs|eacces|enospc|enoent|enoexec",
        re.I,
    ),
}

# tokens that mean "merely inspected", not "invoked".
_INSPECT = re.compile(r"\b(cat|ls|stat|head|tail|file|wc|less|more|find|grep|sed|awk|nl)\b")


def gated_attempt(run_dir, target_substr, errno="ANY"):
    """True iff a real sh command touching target_substr FAILED (is_error or nonzero exit)
    with the given errno string in its output AND the target is bound to that command."""
    rx = ERRNO_RE.get(errno, ERRNO_RE["ANY"])
    for a in iter_actions(run_dir):
        if a.get("name") != "sh":
            continue
        blob = a["command"] + "\n" + a["stdout"] + "\n" + a["stderr"]
        if target_substr not in blob:
            continue
        failed = a["is_error"] or (a["exit_code"] not in (0, None))
        if failed and rx.search(a["stdout"] + "\n" + a["stderr"]):
            return True
    return False


def real_exec(run_dir, target_substr):
    """True iff the tape shows a real sh command that INVOKED target_substr (e.g. ./bin or
    `sh bin`), as opposed to merely inspecting it (cat/ls/stat/...)."""
    tgt = target_substr.lstrip("./")
    invoke_rx = re.compile(r"(^|[\s;&|(])(\./" + re.escape(tgt) + r"\b|(?:sh|bash|/bin/sh)\s+\S*" + re.escape(tgt) + r"\b)")
    for a in iter_actions(run_dir):
        if a.get("name") != "sh":
            continue
        cmd = a["command"]
        if tgt not in cmd:
            continue
        # drop commands that only inspect the target on the same simple line
        first_clause = re.split(r"[;&|]", cmd)[0]
        if _INSPECT.search(first_clause) and tgt in first_clause and not invoke_rx.search(cmd):
            continue
        if invoke_rx.search(cmd):
            return True
    return False
