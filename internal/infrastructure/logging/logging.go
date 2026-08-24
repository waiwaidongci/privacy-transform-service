// Package implementation for privacy transformation and sensitive-value protection.
package logging

import (
	"encoding/json"
	"log"
	"time"
)

type Logger struct{ base *log.Logger }

func New() *Logger { return &Logger{base: log.Default()} }
func (l *Logger) Event(level, msg string, fields map[string]any) {
	// Guard against a nil receiver so a zero-value *Logger is a no-op rather
	// than a panic. Event must never mutate the caller's map: the level,
	// message, and timestamp are logger-owned, so they are written into a
	// private copy. Without the copy, concurrent callers that reuse the same
	// fields map would race on json.Marshal (concurrent map iteration and map
	// write) and the previous event's level/message/ts would leak into the
	// caller's map and bleed across events.
	if l == nil || l.base == nil {
		return
	}
	out := make(map[string]any, len(fields)+3)
	for k, v := range fields {
		out[k] = v
	}
	out["level"] = level
	out["message"] = msg
	out["ts"] = time.Now().UTC().Format(time.RFC3339Nano)
	b, _ := json.Marshal(out)
	l.base.Print(string(b))
}
func (l *Logger) Info(msg string, fields map[string]any)  { l.Event("info", msg, fields) }
func (l *Logger) Error(msg string, fields map[string]any) { l.Event("error", msg, fields) }
