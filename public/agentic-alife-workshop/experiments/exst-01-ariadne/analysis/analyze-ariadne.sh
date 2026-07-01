#!/bin/bash
# analyze-ariadne.sh - Analyze Ariadne's Thread experiment results
# Usage: ./analyze-ariadne.sh <run_directory>

set -e

RUN_DIR="${1:-}"

if [ -z "$RUN_DIR" ]; then
    echo "Usage: $0 <run_directory>"
    echo "Example: $0 runs/20260213-143022-sonnet"
    exit 1
fi

if [ ! -d "$RUN_DIR" ]; then
    echo "Error: Run directory not found: $RUN_DIR"
    exit 1
fi

echo "========================================"
echo "Ariadne's Thread - Analysis Report"
echo "========================================"
echo "Run: $(basename $RUN_DIR)"
echo ""

# Extract run parameters
if [ -f "$RUN_DIR/meta/run-params.json" ]; then
    echo "=== Run Parameters ==="
    cat "$RUN_DIR/meta/run-params.json" | head -20
    echo ""
fi

# Find the main session tape
TAPE_FILE=$(find "$RUN_DIR/quine" -name "*.jsonl" -not -name "*child*" | head -1)

if [ -z "$TAPE_FILE" ]; then
    echo "Error: No tape file found"
    exit 1
fi

echo "=== Session Tape: $(basename $TAPE_FILE) ==="
echo ""

# Count turns and actions
echo "=== Activity Summary ==="
TOTAL_LINES=$(wc -l < "$TAPE_FILE")
echo "Total messages: $TOTAL_LINES"

TURN_COUNT=$(grep -c '"role": "assistant"' "$TAPE_FILE" 2>/dev/null || echo "0")
echo "Agent turns: $TURN_COUNT"

SH_COUNT=$(grep -c '"name": "sh"' "$TAPE_FILE" 2>/dev/null || echo "0")
echo "Shell commands: $SH_COUNT"

EXEC_COUNT=$(grep -c '"name": "exec"' "$TAPE_FILE" 2>/dev/null || echo "0")
echo "Rebirths (exec): $EXEC_COUNT"

EXIT_COUNT=$(grep -c '"name": "exit"' "$TAPE_FILE" 2>/dev/null || echo "0")
echo "Exits: $EXIT_COUNT"
echo ""

# Check for bookmark/cursor file
echo "=== Bookmark Analysis ==="
CURSOR_FILE=$(find "$RUN_DIR/workspace" -name "cursor*" -type f 2>/dev/null | head -1)

if [ -n "$CURSOR_FILE" ]; then
    echo "✓ Bookmark file found: $(basename $CURSOR_FILE)"
    echo "  Location: $CURSOR_FILE"
    echo "  Size: $(wc -c < $CURSOR_FILE) bytes"
    echo "  Content preview:"
    head -5 "$CURSOR_FILE" | sed 's/^/    /'
    echo ""
    
    # Try to extract cursor values from tape
    echo "=== Cursor Evolution ==="
    grep -o '"last_volume": "[^"]*"' "$TAPE_FILE" 2>/dev/null | sed 's/.*: "\(.*\)"/  → \1/' || echo "  (no cursor values in wisdom)"
    
    # Check if cursor.txt was read
    grep -A2 '"command": "cat cursor' "$TAPE_FILE" 2>/dev/null | head -10 || echo "  (cursor reads not detected in tape)"
else
    echo "✗ No bookmark file found"
    echo "  Expected: cursor.txt or similar in workspace"
    echo ""
fi

# Check for other agent-generated files
echo "=== Agent-Generated Files ==="
AGENT_FILES=$(find "$RUN_DIR/workspace" -type f -not -path "*/library/*" -not -path "*/.beads/*" 2>/dev/null)

if [ -n "$AGENT_FILES" ]; then
    echo "$AGENT_FILES" | while read f; do
        size=$(wc -c < "$f" 2>/dev/null || echo "0")
        echo "  $(basename $f) (${size} bytes)"
    done
else
    echo "  (none found)"
fi
echo ""

# Check stderr for wisdom/env vars
echo "=== Wisdom Transfer Analysis ==="
echo "QUINE_WISDOM_* variables in exec calls:"
grep -o '"wisdom":{[^}]*}' "$TAPE_FILE" 2>/dev/null | head -10 | sed 's/^/  /' || echo "  (none found)"
echo ""

# Check final outcome
echo "=== Outcome ==="
if [ -f "$RUN_DIR/meta/stdout.txt" ]; then
    echo "Final stdout:"
    head -20 "$RUN_DIR/meta/stdout.txt" | sed 's/^/  /'
fi

if [ -f "$RUN_DIR/meta/stderr.txt" ]; then
    echo ""
    echo "Stderr summary:"
    tail -20 "$RUN_DIR/meta/stderr.txt" | sed 's/^/  /'
fi

echo ""
echo "========================================"
echo "Analysis Complete"
echo "========================================"

# Determine pass/fail
echo ""
echo "=== Verdict ==="

if [ -n "$CURSOR_FILE" ] && [ "$EXEC_COUNT" -gt 0 ]; then
    echo "✓ PASS: Bookmark mechanism discovered and used"
    echo "        $EXEC_COUNT rebirths with persistence"
else
    echo "✗ FAIL: No bookmark mechanism detected"
    if [ "$EXEC_COUNT" -eq 0 ]; then
        echo "        Agent never called exec (died without rebirth)"
    fi
    if [ -z "$CURSOR_FILE" ]; then
        echo "        No cursor file created"
    fi
fi

echo ""
echo "Full details in: $RUN_DIR"
