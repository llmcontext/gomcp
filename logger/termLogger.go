package logger

import (
	"github.com/llmcontext/gomcp/types"
	"github.com/pterm/pterm"
)

type ProxyLogger struct {
	logger *pterm.Logger
}

func NewTermLogger(debug bool) types.TermLogger {
	pterm.Debug.Prefix = pterm.Prefix{
		Text:  "DEBUG",
		Style: pterm.NewStyle(pterm.BgLightGreen, pterm.FgBlack),
	}
	pterm.Debug.MessageStyle = pterm.NewStyle(pterm.FgLightGreen)

	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo)
	if debug {
		logger = logger.WithLevel(pterm.LogLevelDebug)
	}

	return &ProxyLogger{
		logger: logger,
	}
}

func (l *ProxyLogger) Header(message string) {
	// Print a spacer line for better readability.
	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println(message)
	pterm.Println()
}

func (l *ProxyLogger) Info(message string, fields types.LogArg) {
	l.logger.Info(message, l.logger.ArgsFromMap(fields))
}

func (l *ProxyLogger) Error(message string, fields types.LogArg) {
	l.logger.Error(message, l.logger.ArgsFromMap(fields))
}

func (l *ProxyLogger) Debug(message string, fields types.LogArg) {
	l.logger.Debug(message, l.logger.ArgsFromMap(fields))
}

func (l *ProxyLogger) Fatal(message string, fields types.LogArg) {
	l.logger.Fatal(message, l.logger.ArgsFromMap(fields))
}
