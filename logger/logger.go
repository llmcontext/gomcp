package logger

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerImpl struct {
	zapLog *zap.Logger
}

func NewLogger(config *config.LoggingInfo, debug bool) (types.Logger, error) {
	cfg := zap.NewProductionConfig()

	// Configure encoder to use ISO 8601 timestamp format
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// disable caller to avoid extra noise in the logs (always logger.go anyway)
	cfg.DisableCaller = true

	// if file path is not absolute, we assume it is a relative path
	// to the default hub configuration directory
	// delete output file if present
	if config.File != "" {
		if _, err := os.Stat(config.File); err == nil {
			os.Remove(config.File)
		}
	}

	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else if config.Level != "" {
		level, err := zapcore.ParseLevel(config.Level)
		if err != nil {
			return nil, fmt.Errorf("failed to parse logging level: %v", err)
		}
		cfg.Level = zap.NewAtomicLevelAt(level)
	}

	if config.File != "" {
		cfg.OutputPaths = []string{
			config.File,
		}
	} else {
		cfg.OutputPaths = []string{}
	}
	if config.WithStderr {
		cfg.OutputPaths = append(cfg.OutputPaths, "stderr")
	}

	zapLog := zap.Must(cfg.Build())
	defer zapLog.Sync()

	return &LoggerImpl{
		zapLog: zapLog,
	}, nil
}

func (l *LoggerImpl) Info(message string, fields types.LogArg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	l.zapLog.Info(message, zapFields...)
}

func (l *LoggerImpl) Debug(message string, fields types.LogArg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	l.zapLog.Debug(message, zapFields...)
}

func (l *LoggerImpl) Error(message string, fields types.LogArg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	l.zapLog.Error(message, zapFields...)
}

func (l *LoggerImpl) Fatal(message string, fields types.LogArg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	l.zapLog.Fatal(message, zapFields...)
}
