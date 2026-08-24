package main

import (
	"context"
	"fmt"
	"net/http"

	httpadapter "github.com/ali/go-0821/privacy-transform-service/internal/adapter/http"
)

func defaultServe(ctx context.Context, addr string, handler http.Handler) error {
	return httpadapter.Serve(ctx, addr, handler)
}

func serveRuntime(_ context.Context, deps runtimeDeps, addr string, handler http.Handler) error {
	if err := deps.serve(context.Background(), addr, handler); err != nil {
		return fmt.Errorf("serve runtime: %w", err)
	}
	return nil
}
