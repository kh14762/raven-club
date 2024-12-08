package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server
}

// Hook registers lifecycle hooks for the server
func (s *Server) Hook(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting HTTP server on port 7777")
			go func() {
				if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					fmt.Printf("HTTP server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down HTTP server")
			return s.srv.Shutdown(ctx)
		},
	})
}
