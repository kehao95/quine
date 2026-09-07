# Theory-Spark Query Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).
> Sibling of [`deep-research-query.md`](./deep-research-query.md) (which handles literature fact-finding).


## Purpose & Structural Role

Use `agy` (Gemini CLI) to explore alternative framings, analogies, and candidate
mechanisms. Supply the relevant context, including disagreements and open
questions. A different model offers another perspective; its output still needs
assessment appropriate to the claim. Use the [research query contract](./deep-research-query.md)
for literature fact-finding.

## Invocation (Non-Interactive CLI)

Execute as a one-shot command without opening an interactive session:

```bash
# Direct prompt invocation
agy -p "$(cat /tmp/spark_prompt.txt)"

# With repo context directory
agy --add-dir Paper/theory -p "Provide three divergent analogies for autopoietic closure."
```

*(Note: Never run bare `agy` or `agy -i` in automated scripts; `-p -` does not read stdin—pass arguments explicitly).*

## Integration

Retain useful discussion and unresolved alternatives with their source and
context. Develop promising ideas in the relevant theory, philosophy, or design
record. Assess factual claims against sources and proposed mechanisms through
argument or experiment as appropriate.
