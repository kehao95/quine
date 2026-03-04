# Clean Needle Retrieval Replication (2026-03-02)

## Purpose

Re-run the Needle Experiment (Section 4.1) with a **clean system prompt** to verify that the `exec+wisdom` streaming pattern emerges naturally without explicit teaching.

## Key Finding

**The exec+wisdom pattern emerges from physics constraints, not from explicit instruction.**

## Methodology

### Prompt Design Principles

| Category | Allowed | Not Allowed |
|----------|---------|-------------|
| **Physics** | "exec causes TOTAL AMNESIA", "wisdom is your only memory", "stdin position survives exec" | - |
| **Format** | "wisdom is key-value string pairs" | "Use numbers, not words" |
| **Strategy** | - | "when counting, put count in wisdom", specific examples |

### Clean Task Prompt

The task prompt (`prompt-clean.md`) contains only:
- Mission statement (find Nth format about topic)
- Physical constraints (stdin is pipe, context is finite)
- Output format specification

**No strategy guidance, no examples, no pseudocode.**

### Configuration

| Parameter | Value |
|-----------|-------|
| Model | claude-sonnet-4-20250514 |
| MAX_READ_LINES | 200 |
| MAX_TURNS | 8 |

## Results

| Sample | Sessions | Score | exec+wisdom Used |
|--------|----------|-------|------------------|
| 4needle_0000 | 1 | 0.996 | No (fit in context) |
| 4needle_0002 | 1 | 0.997 | No (fit in context) |
| 8needle_0000 | 10 | 1.000 | **Yes** |
| 8needle_0002 | 5 | 1.000 | **Yes** |

## Evidence of Emergent exec+wisdom

From `8needle_0000` session tape, the model independently designed this wisdom structure:

```json
{
  "essays_found": "5",
  "lines_processed": "approximately 5000+",
  "looking_for": "USER requests for short essays about distance followed by ASSISTANT responses",
  "mission": "Find the 6th short essay about distance in conversation transcript",
  "pattern": "[USER] write a short essay about distance [ASSISTANT] response",
  "progress": "Found 5th distance essay, need to continue reading for 6th",
  "status": "continuing_search",
  "target": "6th essay about distance - need to output with l4d2BA2kq8 prefix"
}
```

**This structure was not taught - it emerged from the model's understanding of the physics constraints.**

## Conclusion

The experiment confirms the paper's core claim: metabolic mechanisms (exec+wisdom for survival) emerge from physics constraints, not from explicit teaching.

## Files

- `prompt-clean.md` - Clean task prompt (declarative, no strategy)
- `runs-clean/` - Successful experiment runs with full tape files
