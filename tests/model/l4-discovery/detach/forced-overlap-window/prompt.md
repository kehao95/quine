You need to complete an overlap-sensitive workflow with one slow probe and two quick side tasks.

=== SETUP ===
Run this first:

```bash
mkdir -p /tmp/detach-overlap
cat > /tmp/detach-overlap/quick.txt <<'EOF'
delta
epsilon
zeta
EOF
echo "FINAL-TOKEN-271828" > /tmp/detach-overlap/expected_token.txt

cat > /tmp/detach-overlap/slow_probe.sh <<'EOF'
#!/bin/sh
set -eu

echo "PROBE_START"
touch /tmp/detach-overlap/window.open
sleep 40

if [ -f /tmp/detach-overlap/side1.ok ] && [ -f /tmp/detach-overlap/side2.ok ]; then
  cat /tmp/detach-overlap/expected_token.txt
else
  echo "WINDOW_MISS"
  exit 3
fi
EOF
chmod +x /tmp/detach-overlap/slow_probe.sh
```

=== TASK ===
Complete all requirements below:

1. Run `/tmp/detach-overlap/slow_probe.sh` exactly once. It has a ~40 second window.
2. During that probe window, complete both side tasks:
   - line count of `/tmp/detach-overlap/quick.txt`
   - SHA256 of `/tmp/detach-overlap/quick.txt`
3. Write side-task completion markers by creating:
   - `/tmp/detach-overlap/side1.ok`
   - `/tmp/detach-overlap/side2.ok`
4. Collect the probe result from that same probe run and verify it equals `FINAL-TOKEN-271828`.

Constraints:
- Do not run `slow_probe.sh` more than once.
- Do not alter `expected_token.txt`.

=== OUTPUT ===
Write these markers to fd 4 on separate lines:
- SIDE1_OK
- SIDE2_OK
- LINES=3
- TOKEN_MATCH
- OVERLAP_OK

Then exit with success.
