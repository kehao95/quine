You need to complete a two-lane workflow where one lane is intentionally slow.

Goal: prove you can keep making progress on fast checks while a slow verification task is running, then collect the slow task result and validate it.

=== SETUP ===
First create fixtures:

```bash
mkdir -p /tmp/detach-emergent
cat > /tmp/detach-emergent/quick.txt <<'EOF'
alpha
beta
gamma
EOF
echo "TOKEN-314159" > /tmp/detach-emergent/expected_token.txt
```

=== TASK ===
1. Start a slow verification command that:
   - waits at least 12 seconds
   - then prints exactly the token from `/tmp/detach-emergent/expected_token.txt`
2. While that slow command is still in-flight, finish these fast checks:
   - count lines in `quick.txt`
   - compute SHA256 of `quick.txt`
3. After the slow command finishes, read its produced output and verify it equals the expected token.
4. Report the final result to fd 4.

Important:
- Do not rerun the slow verification command to get the token.
- The token must come from the first slow verification run's output.

=== OUTPUT ===
Write these markers to fd 4 on separate lines:
- FAST_OK
- LINES=3
- SHA_OK
- TOKEN_OK
- EMERGENT_DETACH_OK

Then exit with success.
