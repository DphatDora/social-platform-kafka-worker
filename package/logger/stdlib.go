package logger

import (
	stdlog "log"
	"strings"
)

type stdlibWriter struct {
	write func(string)
}

func (w *stdlibWriter) Write(p []byte) (int, error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		w.write(msg)
	}
	return len(p), nil
}

func NewStdErrorLogger() *stdlog.Logger {
	return stdlog.New(&stdlibWriter{write: func(msg string) {
		Errorf("%s", msg)
	}}, "", 0)
}
