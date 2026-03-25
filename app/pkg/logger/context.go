package logger

import (
	"context"

	"go.uber.org/zap"
)

type keyType struct{}

var contextKey = keyType{}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return nil
	}

	l, ok := ctx.Value(contextKey).(*zap.Logger)
	if !ok {
		return nil
	}

	return l
}

func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, contextKey, l)
}
