package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Module("HttpServer",
	fx.Provide(NewHttpServer),
	fx.Invoke(func(lc fx.Lifecycle, server *Server) {
		server.Hook(lc)
	}),
)

// NewHttpServer creates a new HTTP server instance
func NewHttpServer(router *gin.Engine) *Server {
	return &Server{
		router: router,
		srv: &http.Server{
			Addr:    ":7777",
			Handler: router,
		},
	}
}
