package types

import "context"

type LogArg map[string]interface{}

type Logger interface {
	Info(message string, fields LogArg)
	Debug(message string, fields LogArg)
	Error(message string, fields LogArg)
	Fatal(message string, fields LogArg)
}

type SubLogger struct {
	logger Logger
	fields LogArg
}

func NewSubLogger(logger Logger, fields LogArg) Logger {
	return &SubLogger{
		logger: logger,
		fields: fields,
	}
}

func mergeFields(fields ...LogArg) LogArg {
	result := LogArg{}
	for _, field := range fields {
		for k, v := range field {
			result[k] = v
		}
	}
	return result
}

func (l *SubLogger) Info(message string, fields LogArg) {
	l.logger.Info(message, mergeFields(l.fields, fields))
}

func (l *SubLogger) Debug(message string, fields LogArg) {
	l.logger.Debug(message, mergeFields(l.fields, fields))
}

func (l *SubLogger) Error(message string, fields LogArg) {
	l.logger.Error(message, mergeFields(l.fields, fields))
}

func (l *SubLogger) Fatal(message string, fields LogArg) {
	l.logger.Fatal(message, mergeFields(l.fields, fields))
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
