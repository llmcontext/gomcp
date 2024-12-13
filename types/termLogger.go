package types

// type TermLogArg map[string]interface{}

type TermLogger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
	Fatal(message string)
}
