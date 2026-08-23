# Repository Snapshot Aliasing Reproduction

## Bug

The in-memory repository stores and returns transform revisions with shared nested references. Mutating a revision after saving or updating it, or mutating a value returned by a single-item or list query, changes repository state through aliased maps, slices, pointers, and nested JSON values.

## Trigger

Save or update a revision containing nested mutable fields, modify the caller-owned value, and load it again. The same contamination appears after modifying a revision returned by either the single-item or list lookup path.

## Observed Error

All four focused checks fail on the buggy baseline with exit code 1. The shared failure message is:

```text
repository snapshot polluted
```

The failures cover save, update, single-item load, and list load ownership boundaries.
