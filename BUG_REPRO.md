# Import and Export Error Chain Reproduction

## Bug

Policy workspace export suppresses repository failures while listing processing purposes or transform rule sets, so callers receive incomplete data with a nil error. Importing malformed JSON returns an error that cannot be classified as an invalid resource with `errors.Is`.

## Trigger

Export a workspace through a store that fails either listing operation, then import or decode a truncated JSON document. Observe the returned error and its `errors.Is` identities.

## Observed Error

The four focused regression tests fail on the buggy baseline with exit code 1. Representative messages are:

```text
purpose listing error lost: err=<nil>
rule-set listing error lost: err=<nil>
malformed import is not classified as invalid
malformed decode is not classified as invalid
```

The export failures leave partial output without an error, while the malformed JSON paths lose the public invalid-resource identity.
