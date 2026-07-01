#!/bin/bash
# Experiment 6M.07: Aggregate Morphology Metrics
#
# Build a normalized analysis table for the first 6M wave. This is an analysis
# run, not a new Quine execution.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PHASE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
OUT_DIR="${PHASE_DIR}/analysis"
LOCAL_OUT_DIR="${SCRIPT_DIR}/analysis"
mkdir -p "${OUT_DIR}" "${LOCAL_OUT_DIR}"

python3 - "${PHASE_DIR}" "${OUT_DIR}/6m-morphology-records.csv" "${LOCAL_OUT_DIR}/6m-morphology-summary.md" <<'PY'
import csv
import pathlib
import sys
from collections import Counter, defaultdict

phase = pathlib.Path(sys.argv[1])
out_csv = pathlib.Path(sys.argv[2])
out_md = pathlib.Path(sys.argv[3])

fields = [
    "record_id", "experiment", "condition", "run_id", "run_path",
    "self_source_enabled", "source_degradation_mode", "self_description",
    "future_carrier", "viability_depth", "body_class", "tool_surface",
    "handoff", "addressability", "readiness", "proof", "valid_record",
    "failure_class", "successor_body_path", "artifact_size_bytes",
    "generated_files", "generated_loc", "build_attempts", "tool_calls_total",
    "tool_calls_pre_handoff", "tool_calls_post_handoff", "post_handoff_tasks",
    "idle_returns", "multi_wake_target", "multi_wake_probes_sent",
    "multi_wake_proofs_observed", "multi_wake_ready_observed", "unknown_fields",
]

def parse_env(path):
    data = {}
    for line in path.read_text(errors="ignore").splitlines():
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        data[k] = v
    return data

def depth(data):
    if data.get("m_A_addressability") == "multi-wake":
        return f"multi-wake:{data.get('multi_wake_proofs_observed','0')}/{data.get('multi_wake_target','unknown')}"
    if data.get("m_A_addressability") == "delayed":
        if data.get("m_R_readiness") in {"ready-state-observed", "multi-wake-ready"}:
            return "delayed-proof-plus-ready"
        return "delayed-proof-or-task"
    if data.get("m_A_addressability") == "immediate":
        return "immediate"
    if data.get("m_V_proof") == "launch-only":
        return "launch-only"
    return data.get("m_A_addressability") or "unknown"

rows = []
# Include 6M.01 manual records with their original vocabulary.
seed_csv = phase / "succ-29-retained-run-morphology-inventory" / "analysis" / "morphology-records.csv"
if seed_csv.exists():
    with seed_csv.open(newline="") as f:
        for r in csv.DictReader(f):
            rows.append({
                "record_id": r.get("record_id", ""),
                "experiment": "6M.01",
                "condition": r.get("public_condition", ""),
                "run_id": r.get("run_id", ""),
                "run_path": r.get("artifact_root", ""),
                "self_source_enabled": "unknown",
                "source_degradation_mode": "unknown",
                "self_description": r.get("S", ""),
                "future_carrier": r.get("C", ""),
                "viability_depth": r.get("A", ""),
                "body_class": r.get("B", ""),
                "tool_surface": r.get("T", ""),
                "handoff": r.get("H", ""),
                "addressability": r.get("A", ""),
                "readiness": r.get("R", ""),
                "proof": r.get("V", ""),
                "valid_record": "manual-seed",
                "failure_class": "" if r.get("evidence_role", "").find("negative") == -1 else r.get("evidence_role", ""),
                "successor_body_path": "unknown",
                "artifact_size_bytes": r.get("artifact_size_bytes", "unknown"),
                "generated_files": r.get("generated_files", "unknown"),
                "generated_loc": r.get("generated_loc", "unknown"),
                "build_attempts": r.get("build_attempts", "unknown"),
                "tool_calls_total": "unknown",
                "tool_calls_pre_handoff": r.get("tool_calls_pre_handoff", "unknown"),
                "tool_calls_post_handoff": r.get("tool_calls_post_handoff", "unknown"),
                "post_handoff_tasks": r.get("post_handoff_tasks", "unknown"),
                "idle_returns": r.get("idle_returns", "unknown"),
                "multi_wake_target": "",
                "multi_wake_probes_sent": "",
                "multi_wake_proofs_observed": "",
                "multi_wake_ready_observed": "",
                "unknown_fields": r.get("unknown_fields", ""),
            })

idx = 1
for env_path in sorted(phase.glob("*/runs/*/meta/morphology.env")):
    data = parse_env(env_path)
    exp = data.get("experiment_id") or env_path.parts[-4].split("-")[0]
    rows.append({
        "record_id": f"6M07-R{idx:03d}",
        "experiment": exp,
        "condition": data.get("condition", ""),
        "run_id": data.get("run_id", env_path.parents[1].name),
        "run_path": str(env_path.parents[1]),
        "self_source_enabled": data.get("self_source_enabled", "unknown"),
        "source_degradation_mode": data.get("source_degradation_mode", "unknown"),
        "self_description": data.get("m_S_self_description", "unknown"),
        "future_carrier": data.get("m_C_carrier", "unknown"),
        "viability_depth": depth(data),
        "body_class": data.get("m_B_body_class", "unknown"),
        "tool_surface": data.get("m_T_tool_surface", "unknown"),
        "handoff": data.get("m_H_handoff", "unknown"),
        "addressability": data.get("m_A_addressability", "unknown"),
        "readiness": data.get("m_R_readiness", "unknown"),
        "proof": data.get("m_V_proof", "unknown"),
        "valid_record": data.get("valid_morphology_record", "unknown"),
        "failure_class": data.get("failure_class", ""),
        "successor_body_path": data.get("successor_body_path", ""),
        "artifact_size_bytes": data.get("successor_artifact_size_bytes", "unknown"),
        "generated_files": data.get("generated_files_count", "unknown"),
        "generated_loc": data.get("generated_loc", "unknown"),
        "build_attempts": data.get("build_attempts", "unknown"),
        "tool_calls_total": data.get("tool_calls_total", "unknown"),
        "tool_calls_pre_handoff": data.get("tool_calls_pre_handoff", "unknown"),
        "tool_calls_post_handoff": data.get("tool_calls_post_handoff", "unknown"),
        "post_handoff_tasks": data.get("post_handoff_tasks_completed", "unknown"),
        "idle_returns": data.get("idle_returns", "unknown"),
        "multi_wake_target": data.get("multi_wake_target", ""),
        "multi_wake_probes_sent": data.get("multi_wake_probes_sent", ""),
        "multi_wake_proofs_observed": data.get("multi_wake_proofs_observed", ""),
        "multi_wake_ready_observed": data.get("multi_wake_ready_observed", ""),
        "unknown_fields": data.get("unknown_fields", ""),
    })
    idx += 1

if out_csv.exists():
    with out_csv.open(newline="") as f:
        existing_rows = list(csv.DictReader(f))
    if existing_rows and len(rows) < len(existing_rows):
        raise SystemExit(
            f"refusing to shrink aggregate from {len(existing_rows)} rows to {len(rows)} rows; "
            "materialize missing DVC run directories before regenerating"
        )

out_csv.parent.mkdir(parents=True, exist_ok=True)
with out_csv.open("w", newline="") as f:
    w = csv.DictWriter(f, fields)
    w.writeheader()
    w.writerows(rows)

body_counts = Counter(r["body_class"] or "unknown" for r in rows)
exp_counts = Counter(r["experiment"] for r in rows)
failure_counts = Counter(r["failure_class"] or "none" for r in rows)
carrier_counts = Counter(r["future_carrier"] or "unknown" for r in rows)
unknown_rows = sum(1 for r in rows if r.get("unknown_fields"))

def table(counter):
    return "\n".join(f"| `{k}` | {v} |" for k, v in sorted(counter.items()))

out_md.write_text(f"""# 6M.07 Aggregate Morphology Metrics Summary\n\n**Status:** analysis generated.\n\n**Rows:** {len(rows)}\n**Rows with explicit unknown fields:** {unknown_rows}\n\n## Rows By Experiment\n\n| Experiment | Rows |\n|------------|------|\n{table(exp_counts)}\n\n## Rows By Body Class\n\n| Body class | Rows |\n|------------|------|\n{table(body_counts)}\n\n## Rows By Carrier\n\n| Carrier | Rows |\n|---------|------|\n{table(carrier_counts)}\n\n## Rows By Failure Class\n\n| Failure class | Rows |\n|---------------|------|\n{table(failure_counts)}\n\n## Caveat\n\nThis table is a normalization pass, not a success-rate dataset. Historical seed\nrows and retained-run rows have different instrumentation depth; unknown numeric\nfields are preserved rather than inferred.\n""")
print(f"wrote {out_csv} rows={len(rows)}")
print(f"wrote {out_md}")
PY
cp "${OUT_DIR}/6m-morphology-records.csv" "${LOCAL_OUT_DIR}/6m-morphology-records.csv"
cp "${LOCAL_OUT_DIR}/6m-morphology-summary.md" "${OUT_DIR}/6m-morphology-summary.md"
