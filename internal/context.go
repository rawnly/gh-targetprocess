package internal

import (
	"context"

	"github.com/rawnly/gh-targetprocess/internal/config"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
)

type ContextKey string

const (
	targetProcessKey ContextKey = "targetprocess"
	configKey        ContextKey = "config"
)

func GetTargetProcess(ctx context.Context) *targetprocess.Client {
	if tp, ok := ctx.Value(targetProcessKey).(*targetprocess.Client); ok {
		return tp
	}
	return nil
}

func GetConfig(ctx context.Context) *config.Config {
	if cfg, ok := ctx.Value(configKey).(*config.Config); ok {
		return cfg
	}
	return nil
}

func SetTargetProcess(ctx context.Context, tp *targetprocess.Client) context.Context {
	return context.WithValue(ctx, targetProcessKey, tp)
}

func SetConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func InitContext(ctx context.Context, cfg *config.Config, tp *targetprocess.Client) context.Context {
	ctx = SetConfig(ctx, cfg)
	ctx = SetTargetProcess(ctx, tp)
	return ctx
}
