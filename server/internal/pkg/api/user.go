package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"raven-club/internal/pkg/user"
)

var UserController = fx.Module("UserController", fx.Invoke(InitUserController))

func InitUserController(engine *gin.Engine, us user.Service) {

	engine.GET("/user/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"users": "todo list users here"})
	})
}
