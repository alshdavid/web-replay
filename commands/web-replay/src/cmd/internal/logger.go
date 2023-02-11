package internal_serve

import "fmt"

type ILogger interface {
	Printf(format string, a ...any)
	Println(a ...any)
}

type Logger struct {
	logLevel string
}

func NewLogger(logLevel string) *Logger {
	return &Logger{
		logLevel: logLevel,
	}
}

func (l *Logger) Printf(format string, a ...any) {
	if l.logLevel == "silent" {
		return
	}
	fmt.Printf(format, a...)
}

func (l *Logger) Println(a ...any) {
	if l.logLevel == "silent" {
		return
	}
	fmt.Println(a...)
}
