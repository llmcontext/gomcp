package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

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
	if config.File != "" && !config.IsFifo {
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
		// special case for fifo
		if config.IsFifo {
			if err := createFifoIfNotExists(config.File); err != nil {
				return fmt.Errorf("failed to create FIFO: %v", err)
			}
		}
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

func createFifoIfNotExists(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create parent directories if they don't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directories: %w", err)
		}

		// Create the FIFO file
		if err := syscall.Mkfifo(path, 0666); err != nil {
			return fmt.Errorf("failed to create FIFO: %w", err)
		}
	}
	return nil
}
