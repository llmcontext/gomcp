package types

import "context"

type LogArg map[string]interface{}

type Logger interface {
	Info(message string, fields LogArg)
	Debug(message string, fields LogArg)
	Error(message string, fields LogArg)
	Fatal(message string, fields LogArg)
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// loggerContextKey is the key used to store the logger in the context
var loggerKey = contextKey("logger")

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) Logger {
	logger := ctx.Value(loggerKey)
	if logger == nil {
		return nil
	}
	return logger.(Logger)
}
