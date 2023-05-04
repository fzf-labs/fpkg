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

func (l *logger) Debug(args ...interface{}) {
	slog.Debug(logKey, fmt.Sprint(args...))
}

func (l *logger) Info(args ...interface{}) {
	slog.Info(logKey, fmt.Sprint(args...))
}

func (l *logger) Warn(args ...interface{}) {
	slog.Warn(logKey, fmt.Sprint(args...))
}

func (l *logger) Error(args ...interface{}) {
	slog.Error(logKey, fmt.Sprint(args...))
}

func (l *logger) Fatal(args ...interface{}) {
	slog.Error(logKey, fmt.Sprint(args...))
}
