package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"net/http"
	"raven-club/internal/pkg/types"
	"raven-club/internal/pkg/user"
)

var UserController = fx.Module("UserController", fx.Invoke(InitUserController))

func InitUserController(engine *gin.Engine, us user.Service) {

	engine.POST("/api/user/create", func(c *gin.Context) {
		ctx := c.Request.Context()

		u := types.User{
			ID:       uuid.New().ID(),
			Username: "Kevy",
			Email:    "kevin.j.heritage@ravenclub.net",
			Password: "mockingjay333",
		}
		err := us.CreateUser(ctx, u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusCreated, gin.H{"user": u})
	})

	engine.DELETE("/api/user/delete", func(c *gin.Context) {

	})

	engine.GET("/api/user/list", func(c *gin.Context) {
		ctx := c.Request.Context()

		users, err := us.ListUsers(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	})
}
