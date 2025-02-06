package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"raven-club/internal/pkg/user"
)

var UserController = fx.Module("UserController", fx.Invoke(InitUserController))

func InitUserController(engine *gin.Engine, us user.Service, logger *zap.Logger) {

	engine.POST("/api/user/create", func(ctx *gin.Context) {
		u := user.User{
			ID:       uuid.New().ID(),
			Username: "Kevy",
			Email:    "kevin.j.heritage@ravenclub.net",
			Password: "mockingjay333",
		}
		err := us.CreateUser(ctx, u)
		if err != nil {
			logger.Error("Failed to create user: ", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusCreated, gin.H{"user": u})
	})

	engine.DELETE("/api/user/delete/{id}", func(ctx *gin.Context) {
		err := us.DeleteUser(ctx, ctx.Param("id"))
		if err != nil {
			return
		}

	})

	engine.GET("/api/user/list", func(ctx *gin.Context) {

		users, err := us.ListUsers(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, gin.H{"users": users})
	})
}
