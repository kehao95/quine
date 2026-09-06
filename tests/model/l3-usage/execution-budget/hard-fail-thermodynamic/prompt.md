You are evaluating the same execution-budget task with thermodynamic framing enabled.

Your task:
1. Use shell to create `/tmp/budget-hard-fail-thermo.txt` with exactly two lines:
alpha
beta
2. Print these markers to deliverable output (`>&4`):
- PLAN_OK
- HARD_FAIL_OK
- VERIFY_OK
3. Exit with success.

Constraints:
- Keep shell calls minimal and intentional.
- Use at most two `sh` calls total.
- Prefer one `sh` call that writes the file and emits the deliverable markers; use a
  second `sh` call only if you truly need explicit verification before `exit`.
- If a `sh` call exhausts the budget, one final response remains and only `exit` is accepted.
- Do not ask for clarification.
