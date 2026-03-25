package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Info log at the InfoLevel
func Info(msg string, fields ...zap.Field) {
	L().WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Infof similar to .Info(fmt.Sprintf(template, ...))
func Infof(template string, args ...any) {
	L().WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprintf(template, args...))
}

// Error log at the ErrorLevel
func Error(msg string, fields ...zap.Field) {
	L().WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Errorf similar to .Error(fmt.Sprintf(template, ...))
func Errorf(template string, args ...any) {
	L().WithOptions(zap.AddCallerSkip(1)).Error(fmt.Sprintf(template, args...))
}

// Fatal log at the FatalLevel, then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	L().WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Fatalf similar to .Fatal(fmt.Sprintf(template, ...))
func Fatalf(template string, args ...any) {
	L().WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprintf(template, args...))
}
