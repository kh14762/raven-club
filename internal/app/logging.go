package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() *zap.Logger {

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := zap.NewDevelopment()

	return logger.WithOptions(
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}
