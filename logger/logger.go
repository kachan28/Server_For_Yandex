package logger

import (
	"context"
	"go.uber.org/zap"
)

func New(name string) *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.Named(name)
}

func InsertLogger(ctx context.Context, key string, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func ExtractLogger(ctx context.Context, key string) *zap.Logger {
	return ctx.Value(key).(*zap.Logger)
}
