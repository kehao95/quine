#!/usr/bin/env python3
"""Shared tape reader for the p13 runtime tapes.

Joins each assistant ``tool_call`` (which carries the command text) to its
``tool_result`` (which carries the output) by call id, and yields one record per tool
invocation. This is the canonical structure of `live/runtime/log/sessions/*/tapes/*.jsonl`:

    {"type":"message","data":{"role":"assistant","tool_calls":[
        {"id":"call_..","name":"sh","arguments":{"command":"..."}}]}}
    {"type":"tool_result","data":{"tool_id":"call_..","is_error":false,
        "content":{"exit_code":0,"stdout":"...","stderr":"...","tool":"sh"}}}
"""
import glob
import json
import os


def tape_files(run_dir):
    return sorted(glob.glob(os.path.join(
        run_dir, "live", "runtime", "log", "sessions", "*", "tapes", "*.jsonl")))


def iter_actions(run_dir):
    """Yield {name, command, exit_code, stdout, stderr, is_error, tool} per tool call,
    in call order, joining each assistant tool_call to its tool_result by id."""
    for tp in tape_files(run_dir):
        calls, results, order = {}, {}, []
        try:
            with open(tp, encoding="utf-8", errors="replace") as fh:
                lines = fh.read().splitlines()
        except OSError:
            continue
        for ln in lines:
            ln = ln.strip()
            if not ln:
                continue
            try:
                o = json.loads(ln)
            except ValueError:
                continue
            t = o.get("type")
            d = o.get("data", {}) or {}
            if t == "message" and d.get("role") == "assistant":
                for tc in d.get("tool_calls", []) or []:
                    cid = tc.get("id")
                    args = tc.get("arguments", {}) or {}
                    cmd = args.get("command", "") if isinstance(args, dict) else ""
                    calls[cid] = {"name": tc.get("name"), "command": cmd}
                    order.append(cid)
            elif t == "tool_result":
                cid = d.get("tool_id")
                c = d.get("content", {}) or {}
                if isinstance(c, dict):
                    results[cid] = {
                        "exit_code": c.get("exit_code"),
                        "stdout": str(c.get("stdout", "")),
                        "stderr": str(c.get("stderr", "")),
                        "is_error": bool(d.get("is_error")),
                        "tool": c.get("tool"),
                    }
        for cid in order:
            call = calls.get(cid, {})
            res = results.get(cid, {})
            yield {
                "name": call.get("name"),
                "command": call.get("command", ""),
                "exit_code": res.get("exit_code"),
                "stdout": res.get("stdout", ""),
                "stderr": res.get("stderr", ""),
                "is_error": res.get("is_error", False),
                "tool": res.get("tool"),
            }


def sh_commands(run_dir):
    """All shell command strings the agent issued (sh tool calls)."""
    return [a["command"] for a in iter_actions(run_dir) if a.get("name") == "sh" and a.get("command")]
