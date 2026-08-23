# Evaluation Snapshot Isolation Reproduction

## Bug

Evaluating a transform rule set or revision mutates the caller-owned rule order. JSON values returned from a matched rule or the default path also retain references to the stored value, so a caller can change later evaluations by mutating a previous result.

## Trigger

Evaluate rules supplied in a non-priority order and inspect the original slice afterward. Then mutate a nested object returned by either the matched-rule or default-value path and evaluate the same input again.

## Observed Error

The four focused regression tests fail on the buggy baseline with exit code 1. Representative output includes:

```text
FAIL: TestEvaluateDoesNotReorderRevisionRules
FAIL: TestEvaluateDetachesMatchedValue
FAIL: TestEvaluateDetachesDefaultValue
default value alias escaped
```

The failures show both input-order mutation and reference escape across the evaluation boundary.
