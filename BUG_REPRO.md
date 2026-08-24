# Bug Reproduction

## What is wrong

Policy publishing and rollback do not preserve a consistent published-history state. Repeated rollback can move forward to a newer version, drafts can be promoted by rollback, and a subsequent publish is reported as not found.

## How to trigger it

Create and publish three policy revisions, call rollback twice, then inspect the active revision. The first rollback selects the prior revision, while the second can select the newer revision again. Try rolling back a draft-only policy or publishing the same already-published policy again to exercise the invalid state paths.

## Observed error

The reproduction returned `resource not found` for publish-after-rollback and for repeated publish attempts. A draft-only rollback could also return success and mark a draft as published.
