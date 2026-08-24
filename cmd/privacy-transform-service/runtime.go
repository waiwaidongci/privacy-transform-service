package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/config"
)

var ErrRuntimeNotStarted = errors.New("service runtime not started")

type runtimeDeps struct {
	load  func(context.Context) (config.Config, error)
	build func(context.Context, config.Config) (http.Handler, error)
	serve func(context.Context, string, http.Handler) error
}

func defaultRuntimeDeps() runtimeDeps {
	return runtimeDeps{load: defaultLoadConfig, build: defaultBuildHandler, serve: defaultServe}
}

func runRuntime(ctx context.Context, deps runtimeDeps) error {
	cfg, err := loadRuntimeConfig(ctx, deps)
	if err != nil {
		return err
	}
	handler, err := buildRuntimeHandler(ctx, deps, cfg)
	if err != nil {
		return err
	}
	return serveRuntime(ctx, deps, cfg.HTTPAddr, handler)
}
