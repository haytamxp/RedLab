package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(l *zap.Logger) {
	Log = l
}