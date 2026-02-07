package logger

import "log/slog"

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type SlogWrapper struct {
	l *slog.Logger
}

func NewSlogWrapper(l *slog.Logger) *SlogWrapper {
	return &SlogWrapper{l: l}
}

func (l *SlogWrapper) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *SlogWrapper) Warn(msg string, args ...any) {
	l.l.Warn(msg, args...)
}

func (l *SlogWrapper) Error(msg string, args ...any) {
	l.l.Error(msg, args...)
}
