package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func Launch() {
	logger := logger()
	fx.New(Module,
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(func() *gin.Engine { return gin.Default() }),
		fx.Provide(func() *zap.Logger { return logger }),
	).Run()
}
