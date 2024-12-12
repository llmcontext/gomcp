package proxy

import "github.com/pterm/pterm"

type ProxyLogger struct {
}

func NewProxyLogger() *ProxyLogger {
	pterm.Debug.Prefix = pterm.Prefix{
		Text:  "DEBUG",
		Style: pterm.NewStyle(pterm.BgLightGreen, pterm.FgBlack),
	}
	pterm.Debug.MessageStyle = pterm.NewStyle(pterm.FgLightGreen)

	pterm.EnableDebugMessages()
	return &ProxyLogger{}
}

func (l *ProxyLogger) Info(message string) {
	pterm.Info.Println(message)
}

func (l *ProxyLogger) Error(message string) {
	pterm.Error.Println(message)
}

func (l *ProxyLogger) Debug(message string) {
	pterm.Debug.Println(message)
}
