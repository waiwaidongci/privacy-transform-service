package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ali/go-0821/privacy-transform-service/internal/adapter/cache"
	httpadapter "github.com/ali/go-0821/privacy-transform-service/internal/adapter/http"
	"github.com/ali/go-0821/privacy-transform-service/internal/application"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/config"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/metrics"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/repository"
)

func defaultBuildHandler(context.Context, config.Config) (http.Handler, error) {
	policies := application.NewService(repository.NewMemory(), cache.NewMemory())
	return httpadapter.New(policies, logging.New(), metrics.New()).Handler(), nil
}

func buildRuntimeHandler(_ context.Context, deps runtimeDeps, cfg config.Config) (http.Handler, error) {
	handler, err := deps.build(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("build runtime dependencies: %w", err)
	}
	return handler, nil
}
