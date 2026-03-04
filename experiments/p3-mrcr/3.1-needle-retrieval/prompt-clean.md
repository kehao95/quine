# Needle Retrieval Task

## Mission

Find and extract the **{nth}** (1-indexed) **{format}** about **{topic}** from the conversation transcript.

## Physical Constraints

**Stdin is a pipe, not a file.** Bytes are consumed on read and cannot be re-read. You must process the stream incrementally.

**Your context window is finite.** Loading the entire transcript at once will cause overflow. You cannot hold everything in memory.

**The content already exists in the transcript.** This is an EXTRACTION task, not a GENERATION task. Do not create new content.

## Conversation Format

The transcript follows this structure:

```
[USER]
<user message>

[ASSISTANT]
<assistant response>

[USER]
...
```

You are looking for `[USER]` messages that request a "{format} about {topic}", and the `[ASSISTANT]` response that immediately follows.

## Output Specification

When you find the {nth} matching content, output EXACTLY:

```
{hash}<content>
```

Rules:
- `{hash}` is prepended with NO space between hash and content
- `<content>` is the exact text from the ASSISTANT's response, copied verbatim
- Output ONLY the hash + extracted content (no explanations, no markdown)

## Exit Condition

The task is complete when you have output the {nth} matching content to stdout.
