# Bug Reproduction

## Bug

The logging helpers write event fields and caller metadata into the map supplied by the caller. Reusing that map across concurrent events can mix fields between requests and race with JSON encoding. Calling a nil or zero-value logger can also panic while accessing its underlying logger.

## Trigger

Reuse one fields map across concurrent `Event` calls, or pass the same map through `WithCaller` and `WithDuration`. Calling `Info` on a nil logger exercises the zero-value path.

## Error

The concurrent path can terminate with `fatal error: concurrent map iteration and map write`; a nil logger can panic with `runtime error: invalid memory address or nil pointer dereference`.
