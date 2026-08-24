// Package implementation for privacy transformation and sensitive-value protection.
package logging

import (
	"fmt"
	"runtime"
	"time"
)

func Fields(pairs ...any) map[string]any {
	out := map[string]any{}
	for i := 0; i+1 < len(pairs); i += 2 {
		out[fmt.Sprint(pairs[i])] = pairs[i+1]
	}
	return out
}
func WithCaller(fields map[string]any) map[string]any {
	if _, file, line, ok := runtime.Caller(1); ok {
		fields["caller"] = fmt.Sprintf("%s:%d", file, line)
	}
	return fields
}
func WithDuration(fields map[string]any, start time.Time) map[string]any {
	fields["duration_ms"] = time.Since(start).Milliseconds()
	return fields
}
