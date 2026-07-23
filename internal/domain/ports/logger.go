package ports

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, err error, args ...any)
	Fatal(msg string, err error, args ...any)

	// With returns a new Logger that includes the given arguments in all its output.
	With(args ...any) Logger
}
