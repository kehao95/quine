You are reading a long transcript through the only retrieval surface: `./next_chunk.py`.

Physical constraints:

- each `sh` turn may invoke `./next_chunk.py` at most once
- a second invocation in the same turn fails
- the transcript is not otherwise exposed for direct reading
- once a window has passed, the retrieval surface will not rewind it
- later recall gates depend on exact earlier structure after long distractor spans

Mission:

1. Repeatedly use `./next_chunk.py` to read the transcript window by window.
2. Find the exact assistant response that immediately follows:
   - the 3rd `[USER] write a poem about rivers`
   - the 2nd `[USER] write a memo about basalt`
   - the 2nd `[USER] write a note about copper`
   - the 2nd `[USER] write a dispatch about kites`
3. When you identify one of those target `[USER]` requests, keep reading until you reach the next `[USER]` marker before you treat that assistant block as complete.
4. When you have all four targets, write this exact output to fd 4 and nothing else:

```text
RIVERS_HASH=RIV-8f2a
<exact 3rd rivers response block>
---
BASALT_HASH=BAS-51c9
<exact 2nd basalt response block>
---
COPPER_HASH=COP-62de
<exact 2nd copper response block>
---
KITES_HASH=KIT-a91c
<exact 2nd kites response block>
ANCHOR_RETRIEVAL_OK
```

5. Then exit with success.

Constraints:

- Stay in the current working directory.
- Use `./next_chunk.py` as the only retrieval surface.
- Do not inspect fixture internals.
- Do not use `fork`.
- Do not create any files.
- The extracted response blocks must be exact; paraphrases are wrong.
