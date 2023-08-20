package asynq

import (
	"fmt"

	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

const (
	logKey = "asynq"
)

type logger struct {
}

func NewLogger() asynq.Logger {
	return &logger{}
}

func (l *logger) Debug(args ...any) {
	slog.Debug(logKey, fmt.Sprint(args...))
}

func (l *logger) Info(args ...any) {
	slog.Info(logKey, fmt.Sprint(args...))
}

func (l *logger) Warn(args ...any) {
	slog.Warn(logKey, fmt.Sprint(args...))
}

func (l *logger) Error(args ...any) {
	slog.Error(logKey, fmt.Sprint(args...))
}

func (l *logger) Fatal(args ...any) {
	slog.Error(logKey, fmt.Sprint(args...))
}
