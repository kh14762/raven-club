package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Module("HttpServer", fx.Provide(NewHttpServer),
	fx.Invoke(func(lc fx.Lifecycle, server *Server) { server.Hook(lc) }))

type Server struct {
	router *gin.Engine
	srv    *http.Server
}

func NewHttpServer(router *gin.Engine) *Server {
	return &Server{
		router: router,
		srv: &http.Server{
			Addr:    ":7777",
			Handler: router,
		},
	}
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
