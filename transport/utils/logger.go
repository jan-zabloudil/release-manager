package utils

import (
	"log/slog"

	logx "go.strv.io/net/logger"
)

type ServerLogger struct {
	Logger *slog.Logger
}

func NewServerLogger(caller string) ServerLogger {
	return ServerLogger{Logger: slog.Default().With("caller", caller)}
}

func (l ServerLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l ServerLogger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l ServerLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l ServerLogger) Error(msg string, err error) {
	l.Logger.Error(msg, "error", err)
}

func (l ServerLogger) With(fields ...logx.Field) logx.ServerLogger {
	f := make([]any, 0, len(fields))
	for _, field := range fields {
		f = append(f, field.Key, field.Value)
	}
	return ServerLogger{Logger: l.Logger.With(f...)}
}
