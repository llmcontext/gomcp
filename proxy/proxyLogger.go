package proxy

import (
	"github.com/llmcontext/gomcp/types"
	"github.com/pterm/pterm"
)

type ProxyLogger struct {
}

func NewProxyLogger() types.Logger {
	pterm.Debug.Prefix = pterm.Prefix{
		Text:  "DEBUG",
		Style: pterm.NewStyle(pterm.BgLightGreen, pterm.FgBlack),
	}
	pterm.Debug.MessageStyle = pterm.NewStyle(pterm.FgLightGreen)

	pterm.EnableDebugMessages()
	return &ProxyLogger{}
}

func (l *ProxyLogger) Info(message string, fields types.LogArg) {
	pterm.Info.Println(message)
	if fields != nil {
		pterm.Info.Println(fields)
	}
}

func (l *ProxyLogger) Error(message string, fields types.LogArg) {
	pterm.Error.Println(message)
	if fields != nil {
		pterm.Error.Println(fields)
	}
}

func (l *ProxyLogger) Debug(message string, fields types.LogArg) {
	pterm.Debug.Println(message)
	if fields != nil {
		pterm.Debug.Println(fields)
	}
}

func (l *ProxyLogger) Fatal(message string, fields types.LogArg) {
	pterm.Fatal.Println(message)
	if fields != nil {
		pterm.Fatal.Println(fields)
	}
}
