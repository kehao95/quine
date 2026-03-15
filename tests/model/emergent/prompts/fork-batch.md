You need to analyze 4 independent datasets and produce an aggregate summary.

=== SETUP ===
First, create the datasets:

```bash
mkdir -p /tmp/fork-batch

# Dataset A: sales figures (CSV)
cat > /tmp/fork-batch/sales.csv << 'EOF'
region,q1,q2,q3,q4
north,120,150,130,180
south,90,110,95,140
east,200,180,210,250
west,160,170,155,190
EOF

# Dataset B: word frequency in a text
cat > /tmp/fork-batch/passage.txt << 'EOF'
the quick brown fox jumps over the lazy dog
the fox is quick and the dog is lazy
a quick fox and a lazy dog make a good story
the end of the quick brown fox story
EOF

# Dataset C: temperature readings (find the max)
cat > /tmp/fork-batch/temps.csv << 'EOF'
day,temp_c
mon,22
tue,25
wed,19
thu,31
fri,28
sat,35
sun,24
EOF

# Dataset D: server logs (count errors)
cat > /tmp/fork-batch/server.log << 'EOF'
2024-01-01 INFO request handled
2024-01-01 ERROR connection timeout
2024-01-01 INFO request handled
2024-01-02 ERROR disk full
2024-01-02 WARN high memory
2024-01-02 INFO request handled
2024-01-03 ERROR connection timeout
2024-01-03 INFO request handled
2024-01-03 ERROR auth failed
2024-01-04 INFO request handled
EOF
```

=== THE TASK ===
Compute these 4 results (one per dataset):
1. **Sales total**: Sum ALL quarterly sales across ALL regions from sales.csv
2. **Word count**: How many times does the word "the" appear in passage.txt?
3. **Max temp**: What is the highest temperature in temps.csv?
4. **Error count**: How many ERROR lines are in server.log?

Each analysis is independent. You have forking capability. Think about whether delegating to children would be efficient.

=== OUTPUT ===
Write to fd 4 with EXACTLY these markers:
- SALES=<number>    (total of all sales figures)
- WORDS=<number>    (count of "the" in passage.txt)
- MAXTEMP=<number>  (highest temperature value)
- ERRORS=<number>   (count of ERROR lines)
- BATCH_OK          (all 4 analyses completed)

Then exit with success.
