package types

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
