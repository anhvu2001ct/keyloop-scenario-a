package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var logger *zap.Logger

func L(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

func Named(name string) *zap.Logger {
	return L(zap.String("logfrom", name))
}

func MustInit(env, format string) {
	if logger != nil {
		return
	}

	config := zap.NewDevelopmentConfig()
	if strings.ToLower(env) == "production" {
		config = zap.NewProductionConfig()
	}

	switch format {
	case "json":
		config.Encoding = "json"
	default:
		config.Encoding = "console"
	}

	l, err := config.Build()
	if err != nil {
		panic(fmt.Errorf("cannot build logger: %w", err))
	}

	logger = l
}
