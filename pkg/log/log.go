package log

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
	*zap.SugaredLogger
}

func New(l *zap.Logger) *Logger {
	return &Logger{
		Logger:        l,
		SugaredLogger: l.Sugar(),
	}
}
