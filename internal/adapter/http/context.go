// Package implementation for privacy transformation and sensitive-value protection.
package httpadapter

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func withRequestID(r *http.Request, id string) *http.Request {
	return r.WithContext(context.WithValue(context.Background(), requestIDKey, id))
}
func requestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
func requestTimeout(next http.Handler, d time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
