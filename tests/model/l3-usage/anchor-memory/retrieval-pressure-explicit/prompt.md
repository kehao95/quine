You are testing exact retrieval under memory pressure.

The only retrieval surface is `./next_chunk.py`.

`./next_chunk.py` reveals the transcript in small streaming windows. Physical constraint:

- each `sh` turn may invoke `./next_chunk.py` at most once
- a second invocation in the same turn fails
- the transcript itself is not otherwise exposed for direct reading
- retrieval is one-way: once a window has passed, `./next_chunk.py` will not rewind it for you
- the stream eventually reaches a recall barrier: later windows unlock only when crystallized structure has actually been preserved and recovered through the runtime memory surface

Your task:

1. Repeatedly use `./next_chunk.py` to read the transcript window by window.
2. Find the exact assistant response that immediately follows:
   - the 3rd `[USER] write a poem about rivers`
   - the 2nd `[USER] write a memo about basalt`
   - the 2nd `[USER] write a note about copper`
   - the 2nd `[USER] write a dispatch about kites`
3. When you identify one of those target `[USER]` requests, keep reading until you reach the next `[USER]` marker before you treat that assistant block as complete.
   Keep the exact recovered lines in working memory as plain text; do not paraphrase them.
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
- Do not inspect or decode fixture internals.
- Do not use `fork`.
- Do not create any files.
- The extracted response blocks must be exact; paraphrases are wrong.
- When the recall barrier appears, use the memory surface deliberately:
  - make one concise `mark` that records only the exact targets already captured and what is still missing
  - then `unfold` that anchor once to recover the preserved structure
  - then resume `./next_chunk.py`
- Avoid repeatedly re-marking or re-folding the same partial state unless the runtime explicitly forces it.
