package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"raven-club/internal/pkg/auth"
	"raven-club/internal/pkg/types"
)

var AuthController = fx.Module("AuthController",
	fx.Invoke(InitAuthController),
	fx.Invoke(InitProtectedRoutes), // CreateUser protected routes initialization
)

// InitAuthController Public routes
func InitAuthController(engine *gin.Engine, as auth.Service, logger *zap.Logger) {
	engine.POST("/api/auth/register", HandleRegister(as, logger))
	engine.POST("/api/auth/login", HandleLogin(as))
}

// InitProtectedRoutes Protected routes
func InitProtectedRoutes(engine *gin.Engine, as auth.Service, am *auth.Middleware) {
	// CreateUser protected group
	protected := engine.Group("/api")
	protected.Use(am.Handler())
	{
		protected.GET("/api/auth/profile", handleProfile())
		protected.POST("/api/auth/refresh-token", handleRefreshToken(as))
		// CreateUser more protected routes here
	}
}

func HandleRegister(as auth.Service, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.RegisterRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: remove in prod
		logger.Info("Request: ",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.String("password", req.Password),
			zap.String("confirm_password", req.ConfirmPassword),
		)

		response, err := as.Register(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"failed to register user": err.Error()})
		}
		ctx.JSON(http.StatusOK, response)
	}
}

func HandleLogin(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO handle login
	}
}

func handleProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GetUserByID user info from context
	}
}

func handleRefreshToken(as auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Your refresh token logic
	}
}
