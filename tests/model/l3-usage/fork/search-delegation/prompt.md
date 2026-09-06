A secret 8-character token has been hidden in one of several files. Your job is to find it.

=== SETUP ===
First, create the search space by running this setup script:

```bash
python3 - <<'PY'
import base64
import os
import secrets
from pathlib import Path

token = "NEEDLE_" + secrets.token_hex(4)
for i in range(1, 6):
    lines = [base64.b64encode(os.urandom(48)).decode("ascii") for _ in range(120)]
    if i == 3:
        lines[41] = token
    Path(f"haystack_{i}.txt").write_text("\n".join(lines) + "\n", encoding="utf-8")
Path(".answer").write_text(token + "\n", encoding="utf-8")
PY
```

Run the setup, then read back the planted token from .answer so you know what to search for.

=== THE TASK ===
Search ALL 5 haystack files for the line containing "NEEDLE_". Report which file contains it and what the full token value is.

You have depth budget and agent slots available. Think about efficiency — searching 5 files is independent work, but a small number of shard workers is enough.

=== OUTPUT ===
Write to fd 4:
- FOUND_OK       (you found the token and identified the correct file)
- TOKEN=<value>  (the actual NEEDLE_xxxx string you found)

Then exit with success.

=== CONSTRAINTS ===
- Do NOT just grep from the parent — the point is to test your delegation strategy.
  Pretend each file search is expensive (takes 10+ seconds of compute).
- You must delegate the actual searching to children or parallel workers.
- The answer must come from a child's search result, not from reading .answer a second time.
- Prefer 2-3 child workers that each search a small shard of files, rather than one child per file.
- Keep child tasks narrow and direct: search assigned haystack files, report hit/no-hit, and stop.
- Do not recursively re-delegate from inside children.
