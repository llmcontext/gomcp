package logger

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

type Arg map[string]interface{}

func InitLogger(config config.LoggingInfo, debug bool) error {
	cfg := zap.NewProductionConfig()

	// Configure encoder to use ISO 8601 timestamp format
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// disable caller to avoid extra noise in the logs (always logger.go anyway)
	cfg.DisableCaller = true

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
			return fmt.Errorf("failed to parse logging level: %v", err)
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

	zapLog = zap.Must(cfg.Build())
	defer zapLog.Sync()

	return nil
}

func Info(message string, fields Arg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zapLog.Info(message, zapFields...)
}

func Debug(message string, fields Arg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zapLog.Debug(message, zapFields...)
}

func Error(message string, fields Arg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zapLog.Error(message, zapFields...)
}

func Fatal(message string, fields Arg) {
	zapFields := []zap.Field{}
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	zapLog.Fatal(message, zapFields...)
}
