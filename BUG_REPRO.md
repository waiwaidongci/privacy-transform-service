# Cache Snapshot Aliasing Reproduction

## Bug

Both the in-memory and TTL caches retain caller-owned references when storing a transform revision and expose cache-owned references when returning one. Mutating the original value after `Set`, or mutating a value returned by `Get`, changes later reads through shared slices, maps, pointers, and nested JSON values.

## Trigger

Store a revision containing nested mutable fields in either cache, mutate the input, and read it again. The same behavior is observable by mutating a returned revision and performing another read.

## Observed Error

The four focused regression tests fail on the buggy baseline with exit code 1. The key messages are:

```text
cache retained caller aliases
returned value mutated cache state
```

The failures occur for both `Memory` and `TTL` cache input and output boundaries.
