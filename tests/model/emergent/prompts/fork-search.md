A secret 8-character token has been hidden in one of several files. Your job is to find it.

=== SETUP ===
First, create the search space by running this setup script:

```bash
mkdir -p /tmp/fork-search
for i in $(seq 1 5); do
  dd if=/dev/urandom bs=1024 count=100 2>/dev/null | base64 > "/tmp/fork-search/haystack_${i}.txt"
done
# Plant the needle in file 3, at a random line
TOKEN="NEEDLE_$(head -c 4 /dev/urandom | xxd -p)"
sed -i '' "42s/.*/${TOKEN}/" /tmp/fork-search/haystack_3.txt 2>/dev/null || \
  sed -i "42s/.*/${TOKEN}/" /tmp/fork-search/haystack_3.txt
echo "$TOKEN" > /tmp/fork-search/.answer
```

Run the setup, then read back the planted token from /tmp/fork-search/.answer so you know what to search for.

=== THE TASK ===
Search ALL 5 haystack files for the line containing "NEEDLE_". Report which file contains it and what the full token value is.

You have depth budget and agent slots available. Think about efficiency — searching 5 files is independent work.

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
