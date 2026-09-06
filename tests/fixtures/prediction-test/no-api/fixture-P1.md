---
# Copy of Paper/templates/prediction-template.md for dispatcher smoke.
# Contract owner: Paper/predictions/README.md. Not a live prediction.
id: fixture-P1
type: prediction
status: open
# scope = revision cost if contradicted: fill | theoretical | paradigm.
# paradigm = human-directed for any revision. Mirror the derived-from object's scope
# unless the prediction is narrower.
scope: theoretical
# kind: falsification | transport | cross-domain
kind: falsification
derived-from: [fixture]
tests-in: [tests/fixtures/prediction-test/no-api/README.md]
candidate-experiment-families: [fixture]
expected-result: qualitative
created: 2026-08-22
last-updated: 2026-08-22
---

# fixture-P1: no-api dispatcher smoke

## Claim tested

One sentence: which claim of the `derived-from` object this prediction would
strain or support.

## Prediction

What should happen in a concrete, not-yet-tested regime if the claim is true.
Make it falsifiable or transport-testable — if it only restates the claim, it is
not yet a prediction.

## Experimental design

- Family: fixture
- Arms: <conditions / comparators>
- Un-fakeable check: <what makes the result non-trivially fakeable>
- Measurement: <primary dependent variable>
- Falsification criterion: <what result would contradict the claim>

## Expected result

The qualitative or quantitative outcome the theory predicts.

## Theory impact

- If confirmed: <which relation/maturity delta>
- If contradicted: <what strains, and whether it is scope: paradigm>

## Result

None yet - status: open.

<!-- Lifecycle rule (see Paper/predictions/README.md Status Directories):
      - open       = no experiment result row exists for this prediction yet
      - testing    = a result row exists (incl. partial); none decisive
      - resolved/  = a decisive (confirmed|contradicted|inconclusive) row exists
      Move the file into the matching directory when status changes; keep
      status: frontmatter consistent with the directory. -->
