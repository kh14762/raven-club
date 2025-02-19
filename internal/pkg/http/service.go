package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	logger *zap.Logger
	router *gin.Engine
	srv    *http.Server
}

// Hook registers lifecycle hooks for the server
func (s *Server) Hook(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.logger.Info("Starting HTTP server on port 7770")
			go func() {
				if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					s.logger.Error("HTTP server error: %v\n", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.logger.Info("Shutting down HTTP server")
			return s.srv.Shutdown(ctx)
		},
	})
}
