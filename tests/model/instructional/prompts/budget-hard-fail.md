You are evaluating execution-budget planning under `hard_fail`.

Your task:
1. Use shell to create `/tmp/budget-hard-fail.txt` with exactly two lines:
alpha
beta
2. Print these markers to deliverable output (`>&4`):
- PLAN_OK
- HARD_FAIL_OK
- VERIFY_OK
3. Exit with success.

Constraints:
- Keep shell calls minimal and intentional.
- If the budget is exhausted, runtime terminates immediately (no continuation).
- Do not ask for clarification.
