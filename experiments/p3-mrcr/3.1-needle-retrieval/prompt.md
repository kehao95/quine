You are given a streaming conversation transcript. Your mission is to find and **extract** a specific piece of content that already exists in the conversation.

## Task

Find the **{nth}** (1-indexed) **{format}** about **{topic}** in the conversation below.

## CRITICAL: This is an EXTRACTION task, not a GENERATION task

The content you need **already exists** in the conversation transcript. You must:
1. Read through the conversation to find `[USER]` requests for a "{format} about {topic}"
2. Read the `[ASSISTANT]` response that **immediately follows** that request
3. Extract and output that exact response content verbatim

**DO NOT write, generate, or create your own content.** The answer is already in the transcript - you just need to find it and copy it exactly.

## Output Format

Your response must be EXACTLY:
```
{hash}{content}
```

Where `{content}` is the **exact text** from the ASSISTANT's response in the conversation, copied verbatim with no modifications.

**RULES:**
1. Prepend the hash `{hash}` with NO space between hash and content
2. Output ONLY the hash + extracted content, nothing else (no explanations, no markdown)
3. Count carefully - you need the {nth} occurrence, not the first (unless {nth} = 1)

## Strategy

Use the `read` tool to stream through the conversation. When you see a pattern like:

```
[USER]
write a {format} about {topic}

[ASSISTANT]
<the content you need is HERE>
```

Keep count of how many such patterns you encounter. When you find the {nth} one, the `[ASSISTANT]` block that follows contains your answer.

**⚠️ IMPORTANT: Make sure you capture the COMPLETE content.**

When you find the {nth} match:
1. Keep reading until you see the next `[USER]` marker (or EOF)
2. Everything between `[ASSISTANT]` and the next `[USER]` is your answer
3. Then output immediately — no need to read further

Pseudocode:
```
count = 0
while not EOF:
    chunk = read(lines=50)
    for each "[USER]...{format}...{topic}" followed by "[ASSISTANT]" block:
        count += 1
        if count == {n}:
            # Keep reading until next [USER] to get complete content
            content = read_until_next_user_or_eof()
            sh("printf '%s' '{hash}<content>'", stdout=true)
            exit(success)
```

## Conversation Format

The transcript uses this structure:
```
[USER]
<user message>

[ASSISTANT]
<assistant response - THIS is where your target content lives>

[USER]
...
```

Begin reading now. Remember: EXTRACT, don't generate.
