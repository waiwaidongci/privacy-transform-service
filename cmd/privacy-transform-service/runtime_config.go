package main

import (
	"context"
	"fmt"

	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/config"
)

func defaultLoadConfig(context.Context) (config.Config, error) { return config.Load(), nil }

func loadRuntimeConfig(_ context.Context, deps runtimeDeps) (config.Config, error) {
	cfg, err := deps.load(context.Background())
	if err != nil {
		return config.Config{}, fmt.Errorf("load runtime config: %w", err)
	}
	return cfg, nil
}
