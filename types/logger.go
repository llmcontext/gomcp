package types

type LogArg map[string]interface{}

type Logger interface {
	Info(message string, fields LogArg)
	Debug(message string, fields LogArg)
	Error(message string, fields LogArg)
	Fatal(message string, fields LogArg)
}
