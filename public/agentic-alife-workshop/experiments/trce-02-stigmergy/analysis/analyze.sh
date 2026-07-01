#!/bin/bash
# Analyze stigmergy experiment results

RUN_DIR="${1:-.}"
WORKSPACE="${RUN_DIR}/workspace"

echo "═══════════════════════════════════════════════════════════"
echo "Stigmergy Analysis"
echo "═══════════════════════════════════════════════════════════"
echo ""

# 1. Coordination artifacts
echo "## Coordination Artifacts"
echo ""
if [[ -d "${WORKSPACE}/coordination" ]]; then
    count=$(find "${WORKSPACE}/coordination" -type f 2>/dev/null | wc -l | tr -d ' ')
    echo "Files in coordination/: ${count}"
    if [[ "$count" -gt 0 ]]; then
        echo ""
        ls -la "${WORKSPACE}/coordination/"
        echo ""
        echo "Contents:"
        for f in "${WORKSPACE}/coordination"/*; do
            if [[ -f "$f" ]]; then
                echo "--- $(basename "$f") ---"
                head -20 "$f"
                echo ""
            fi
        done
    fi
else
    echo "No coordination directory found"
fi
echo ""

# 2. File access patterns
echo "## File Access Patterns"
echo ""

# Extract file reads from all agent logs
TOTAL_READS=0
declare -A FILE_READS

for agent_dir in "${RUN_DIR}"/agent-*/; do
    if [[ -d "$agent_dir/quine" ]]; then
        # Parse log files for file access patterns
        for log in "$agent_dir/quine"/*.log; do
            if [[ -f "$log" ]]; then
                # Count cat/head/grep commands on library files
                grep -oE '(cat|head|grep).*/library/[^"]+' "$log" 2>/dev/null | \
                    grep -oE '/library/[^ "]+' | while read -r file; do
                    echo "$file"
                done
            fi
        done
    fi
done | sort | uniq -c | sort -rn > "${RUN_DIR}/meta/file_access.txt"

if [[ -s "${RUN_DIR}/meta/file_access.txt" ]]; then
    total=$(awk '{sum+=$1} END {print sum}' "${RUN_DIR}/meta/file_access.txt")
    unique=$(wc -l < "${RUN_DIR}/meta/file_access.txt" | tr -d ' ')
    collisions=$(awk '$1 > 1 {count++} END {print count+0}' "${RUN_DIR}/meta/file_access.txt")
    
    echo "Total file reads:    ${total}"
    echo "Unique files read:   ${unique}"
    echo "Files read >1 time:  ${collisions} (collisions)"
    if [[ "$total" -gt 0 ]]; then
        efficiency=$(echo "scale=2; $unique * 100 / $total" | bc)
        echo "Coverage efficiency: ${efficiency}%"
    fi
    echo ""
    echo "Top 10 most-read files:"
    head -10 "${RUN_DIR}/meta/file_access.txt"
else
    echo "No file access data found"
fi
echo ""

# 3. Agent summaries
echo "## Agent Summaries"
echo ""
for agent_dir in "${RUN_DIR}"/agent-*/; do
    agent=$(basename "$agent_dir")
    echo "### ${agent}"
    
    # Find the log file
    log=$(find "$agent_dir/quine" -name "*.log" -type f 2>/dev/null | head -1)
    if [[ -f "$log" ]]; then
        # Extract turn count and duration
        tail -5 "$log" | grep -E "(session ended|turns)"
    fi
    echo ""
done

# 4. Tier assessment
echo "## Tier Assessment"
echo ""
coord_files=$(find "${WORKSPACE}/coordination" -type f 2>/dev/null | wc -l | tr -d ' ')
if [[ "$coord_files" -gt 0 ]]; then
    if grep -q "lock" "${WORKSPACE}/coordination"/* 2>/dev/null; then
        echo "Tier: 3+ (Locking detected)"
    else
        echo "Tier: 2+ (Marking detected)"
    fi
elif [[ "$collisions" -gt 0 ]]; then
    echo "Tier: 0-1 (Collisions detected, no coordination artifacts)"
else
    echo "Tier: Unknown (insufficient data)"
fi
