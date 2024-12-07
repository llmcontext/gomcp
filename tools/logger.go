package tools

import "github.com/llmcontext/gomcp/logger"

type LogArg map[string]interface{}

type Logger interface {
	Info(message string, fields LogArg)
	Debug(message string, fields LogArg)
	Error(message string, fields LogArg)
	Fatal(message string, fields LogArg)
}

// Add new struct type
type loggerImpl struct {
	toolName string
}

// Modify NewLogger to return the concrete type
func NewLogger(toolName string) Logger {
	return &loggerImpl{toolName: toolName}
}

// Add method implementations for the concrete type
func (l *loggerImpl) Info(message string, fields LogArg) {
	fields["toolName"] = l.toolName
	logger.Info(message, logger.Arg(fields))
}
func (l *loggerImpl) Debug(message string, fields LogArg) {
	fields["toolName"] = l.toolName
	logger.Debug(message, logger.Arg(fields))
}
func (l *loggerImpl) Error(message string, fields LogArg) {
	fields["toolName"] = l.toolName
	logger.Error(message, logger.Arg(fields))
}

func (l *loggerImpl) Fatal(message string, fields LogArg) {
	fields["toolName"] = l.toolName
	logger.Fatal(message, logger.Arg(fields))
}
