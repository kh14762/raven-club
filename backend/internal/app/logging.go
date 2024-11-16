package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() *zap.Logger {

	log, _ := zap.NewDevelopment()

	return log.WithOptions(
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}
