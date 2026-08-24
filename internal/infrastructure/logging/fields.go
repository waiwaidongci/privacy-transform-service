// Package implementation for privacy transformation and sensitive-value protection.
package logging

import (
	"fmt"
	"runtime"
	"time"
)

// Fields assembles key/value pairs into a new map. The returned map is owned
// by the caller; no other logging helper mutates it.
func Fields(pairs ...any) map[string]any {
	out := map[string]any{}
	for i := 0; i+1 < len(pairs); i += 2 {
		out[fmt.Sprint(pairs[i])] = pairs[i+1]
	}
	return out
}

// WithCaller returns a copy of fields with the caller location appended. The
// input map is not mutated, so concurrent goroutines that share a fields map
// cannot clobber each other's caller entry.
func WithCaller(fields map[string]any) map[string]any {
	out := clone(fields)
	if _, file, line, ok := runtime.Caller(1); ok {
		out["caller"] = fmt.Sprintf("%s:%d", file, line)
	}
	return out
}

// WithDuration returns a copy of fields with the elapsed duration appended. The
// input map is not mutated, for the same reason as WithCaller.
func WithDuration(fields map[string]any, start time.Time) map[string]any {
	out := clone(fields)
	out["duration_ms"] = time.Since(start).Milliseconds()
	return out
}

func clone(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	return out
}
