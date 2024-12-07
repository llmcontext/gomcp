package tools

import (
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/types"
)

// Add new struct type
type loggerImpl struct {
	toolName string
}

// Modify NewLogger to return the concrete type
func NewLogger(toolName string) types.Logger {
	return &loggerImpl{toolName: toolName}
}

// Add method implementations for the concrete type
func (l *loggerImpl) Info(message string, fields types.LogArg) {
	fields["toolName"] = l.toolName
	logger.Info(message, logger.Arg(fields))
}
func (l *loggerImpl) Debug(message string, fields types.LogArg) {
	fields["toolName"] = l.toolName
	logger.Debug(message, logger.Arg(fields))
}
func (l *loggerImpl) Error(message string, fields types.LogArg) {
	fields["toolName"] = l.toolName
	logger.Error(message, logger.Arg(fields))
}

func (l *loggerImpl) Fatal(message string, fields types.LogArg) {
	fields["toolName"] = l.toolName
	logger.Fatal(message, logger.Arg(fields))
}
